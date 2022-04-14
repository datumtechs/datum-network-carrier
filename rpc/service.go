package rpc

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	statefeed "github.com/RosettaFlow/Carrier-Go/common/feed/state"
	"github.com/RosettaFlow/Carrier-Go/common/traceutil"
	"github.com/RosettaFlow/Carrier-Go/handler"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	pbrpc "github.com/RosettaFlow/Carrier-Go/lib/rpc/v1"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend/auth"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend/metadata"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend/power"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend/task"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend/yarn"
	"github.com/RosettaFlow/Carrier-Go/rpc/debug"
	health_check "github.com/RosettaFlow/Carrier-Go/service/discovery"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	service_discover_health "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
	"net"
	"sync"
)

var _ common.Service = (*Service)(nil)

// Service defining an RPC server for a carrier server.
type Service struct {
	cfg                  *Config
	ctx                  context.Context
	cancel               context.CancelFunc
	listener             net.Listener
	grpcServer           *grpc.Server
	credentialError      error
	connectedRPCClients  map[net.Addr]bool
	clientConnectionLock sync.Mutex
}

// Config options for the beacon node RPC server.
type Config struct {
	Host                    string
	Port                    string
	CertFlag                string
	KeyFlag                 string
	EnableDebugRPCEndpoints bool
	Broadcaster             p2p.Broadcaster
	PeersFetcher            p2p.PeersProvider
	PeerManager             p2p.PeerManager
	MetadataProvider        p2p.MetadataProvider
	StateNotifier           statefeed.Notifier
	BackendAPI              backend.Backend
	TwoPcAPI				handler.Engine
	MaxMsgSize              int
	MaxSendMsgSize          int
}

// NewService instantiates a new RPC service instance that will
// be registered into a running carrier server.
func NewService(ctx context.Context, cfg *Config) *Service {
	ctx, cancel := context.WithCancel(ctx)
	return &Service{
		cfg:                 cfg,
		ctx:                 ctx,
		cancel:              cancel,
		connectedRPCClients: make(map[net.Addr]bool),
	}
}

// Start the gRPC server.
func (s *Service) Start() error {
	address := fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Errorf("Could not listen to port in grpcSvr.Start() %s: %v", address, err)
		return err
	}
	s.listener = lis
	log.WithField("address", address).Info("gRPC server listening on port")

	opts := []grpc.ServerOption{
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
		grpc.StreamInterceptor(middleware.ChainStreamServer(
			recovery.StreamServerInterceptor(
				recovery.WithRecoveryHandlerContext(traceutil.RecoveryHandlerFunc),
			),
			grpc_prometheus.StreamServerInterceptor,
			grpc_opentracing.StreamServerInterceptor(),
			s.validatorStreamConnectionInterceptor,
		)),
		grpc.UnaryInterceptor(middleware.ChainUnaryServer(
			recovery.UnaryServerInterceptor(
				recovery.WithRecoveryHandlerContext(traceutil.RecoveryHandlerFunc),
			),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_opentracing.UnaryServerInterceptor(),
			s.validatorUnaryConnectionInterceptor,
		)),
		grpc.MaxRecvMsgSize(s.cfg.MaxMsgSize), // 1 << 23 == 1024*1024*8 == 8388608 == 8mb
		grpc.MaxSendMsgSize(s.cfg.MaxSendMsgSize),
	}
	grpc_prometheus.EnableHandlingTimeHistogram()
	if s.cfg.CertFlag != "" && s.cfg.KeyFlag != "" {
		creds, err := credentials.NewServerTLSFromFile(s.cfg.CertFlag, s.cfg.KeyFlag)
		if err != nil {
			log.WithError(err).Fatal("Could not load TLS keys")
			return err
		}
		opts = append(opts, grpc.Creds(creds))
	} else {
		log.Warn("You are using an insecure gRPC server.")
	}
	// create grpc server
	s.grpcServer = grpc.NewServer(opts...)

	// init server instance and register server.
	pb.RegisterYarnServiceServer(s.grpcServer, &yarn.Server{ B: s.cfg.BackendAPI, RpcSvrIp: s.cfg.Host, RpcSvrPort: s.cfg.Port})
	pb.RegisterMetadataServiceServer(s.grpcServer, &metadata.Server{ B: s.cfg.BackendAPI })
	pb.RegisterPowerServiceServer(s.grpcServer, &power.Server{ B: s.cfg.BackendAPI })
	pb.RegisterAuthServiceServer(s.grpcServer, &auth.Server{ B: s.cfg.BackendAPI })
	pb.RegisterTaskServiceServer(s.grpcServer, &task.Server{ B: s.cfg.BackendAPI })
	service_discover_health.RegisterHealthServer(s.grpcServer, &health_check.HealthCheck{})

	if s.cfg.EnableDebugRPCEndpoints {
		log.Info("Enabled debug gRPC endpoints")
		debugServer := &debug.Server{
			PeerManager:  s.cfg.PeerManager,
			PeersFetcher: s.cfg.PeersFetcher,
			ConsensusStateInfo: s.cfg.TwoPcAPI.GetConsensusStateInfo(),
		}
		pbrpc.RegisterDebugServer(s.grpcServer, debugServer)

		// add logrus instanse to gRPC loggerV2 interface implements.
		grpclog.SetLoggerV2(NewRPCdebugLogger(log.Logger))

		//// open grpc trace
		//go func() {
		//	trace.AuthRequest = func(req *http.Request) (any, sensitive bool) {
		//		return true, true
		//	}
		//	go http.ListenAndServe(":50051", nil)
		//	grpclog.Println("Trace listen on 50051")
		//}()
	}
	// Register reflection service on gRPC server.
	reflection.Register(s.grpcServer)

	go func() {
		if s.listener != nil {
			if err := s.grpcServer.Serve(s.listener); err != nil {
				log.Errorf("Could not serve gRPC: %v", err)
			}
		}
	}()
	log.Info("Started grpcServer ...")
	return nil
}

// Stop the service.
func (s *Service) Stop() error {
	s.cancel()
	if s.listener != nil {
		s.grpcServer.GracefulStop()
		log.Debug("Initiated graceful stop of gRPC server")
	}
	return nil
}

// Status returns nil or credentialError
func (s *Service) Status() error {
	if s.credentialError != nil {
		return s.credentialError
	}
	return nil
}

// Stream interceptor for new validator client connections to the beacon node.
func (s *Service) validatorStreamConnectionInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	_ *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	s.logNewClientConnection(ss.Context())
	return handler(srv, ss)
}

// Unary interceptor for new validator client connections to the beacon node.
func (s *Service) validatorUnaryConnectionInterceptor(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	s.logNewClientConnection(ctx)
	return handler(ctx, req)
}

func (s *Service) logNewClientConnection(ctx context.Context) {
	/*if featureconfig.Get().DisableGRPCConnectionLogs {
		return
	}*/
	if clientInfo, ok := peer.FromContext(ctx); ok {
		// Check if we have not yet observed this grpc client connection
		// in the running beacon node.
		s.clientConnectionLock.Lock()
		defer s.clientConnectionLock.Unlock()
		if !s.connectedRPCClients[clientInfo.Addr] {
			log.WithFields(logrus.Fields{
				"addr": clientInfo.Addr.String(),
			}).Debug("New gRPC client connected to carrier node")
			s.connectedRPCClients[clientInfo.Addr] = true
		}
	}
}
