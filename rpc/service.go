package rpc

import (
	"context"
	"fmt"
	"github.com/datumtechs/datum-network-carrier/common"
	statefeed "github.com/datumtechs/datum-network-carrier/common/feed/state"
	"github.com/datumtechs/datum-network-carrier/common/traceutil"
	"github.com/datumtechs/datum-network-carrier/p2p"
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	carrierrpcdebugpbv1 "github.com/datumtechs/datum-network-carrier/pb/carrier/rpc/debug/v1"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"github.com/datumtechs/datum-network-carrier/rpc/backend"
	"github.com/datumtechs/datum-network-carrier/rpc/backend/auth"
	"github.com/datumtechs/datum-network-carrier/rpc/backend/did"
	"github.com/datumtechs/datum-network-carrier/rpc/backend/metadata"
	"github.com/datumtechs/datum-network-carrier/rpc/backend/power"
	"github.com/datumtechs/datum-network-carrier/rpc/backend/task"
	"github.com/datumtechs/datum-network-carrier/rpc/backend/workflow"
	"github.com/datumtechs/datum-network-carrier/rpc/backend/yarn"
	"github.com/datumtechs/datum-network-carrier/rpc/debug"
	health_check "github.com/datumtechs/datum-network-carrier/service/discovery"
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
	rpcMetadata "google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	rpcPeer "google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
	"net"
	"strings"
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
	Host                          string
	Port                          string
	CertFlag                      string
	KeyFlag                       string
	EnableDebugRPCEndpoints       bool
	Broadcaster                   p2p.Broadcaster
	PeersFetcher                  p2p.PeersProvider
	PeerManager                   p2p.PeerManager
	MetadataProvider              p2p.MetadataProvider
	StateNotifier                 statefeed.Notifier
	BackendAPI                    backend.Backend
	DebugAPI                      debug.DebugBackend
	MaxMsgSize                    int
	MaxSendMsgSize                int
	EnableGrpcGateWayPrivateCheck bool
	PrivateIPCache                map[string]struct{} // {"ip1":{},"ip2":{}}
	CarrierPrivateIP              *common.CarrierPrivateIP
	NotCheckPrivateIP             bool
	PublicRpcService              map[int]struct{}
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
	carrierapipb.RegisterYarnServiceServer(s.grpcServer, &yarn.Server{B: s.cfg.BackendAPI, RpcSvrIp: s.cfg.Host, RpcSvrPort: s.cfg.Port})
	carrierapipb.RegisterMetadataServiceServer(s.grpcServer, &metadata.Server{B: s.cfg.BackendAPI})
	carrierapipb.RegisterPowerServiceServer(s.grpcServer, &power.Server{B: s.cfg.BackendAPI})
	carrierapipb.RegisterAuthServiceServer(s.grpcServer, &auth.Server{B: s.cfg.BackendAPI})
	carrierapipb.RegisterTaskServiceServer(s.grpcServer, &task.Server{B: s.cfg.BackendAPI})
	carrierapipb.RegisterDIDServiceServer(s.grpcServer, &did.Server{B: s.cfg.BackendAPI})
	carrierapipb.RegisterVcServiceServer(s.grpcServer, &did.Server{B: s.cfg.BackendAPI})
	carrierapipb.RegisterProposalServiceServer(s.grpcServer, &did.Server{B: s.cfg.BackendAPI})
	carrierapipb.RegisterWorkFlowServiceServer(s.grpcServer, &workflow.Server{B: s.cfg.BackendAPI})
	service_discover_health.RegisterHealthServer(s.grpcServer, &health_check.HealthCheck{})

	if s.cfg.EnableDebugRPCEndpoints {
		log.Info("Enabled debug gRPC endpoints")
		debugServer := &debug.Server{
			PeerManager:  s.cfg.PeerManager,
			PeersFetcher: s.cfg.PeersFetcher,
			DebugAPI:     s.cfg.DebugAPI,
		}
		carrierrpcdebugpbv1.RegisterDebugServer(s.grpcServer, debugServer)

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
	un *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	s.logNewClientConnection(ctx)
	return s.checkRpcServicePrivate(ctx, req, un, handler)
}

func (s *Service) checkRpcServicePrivate(
	ctx context.Context,
	req interface{},
	un *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	if s.cfg.NotCheckPrivateIP {
		return handler(ctx, req)
	}

	fullName := un.FullMethod
	fullNameCode := common.ServiceMethodControlCode[fullName]
	if _, ok := s.cfg.PublicRpcService[fullNameCode]; ok {
		return handler(ctx, req)
	}

	var clientIP string
	in, ok := rpcMetadata.FromIncomingContext(ctx)
	if !ok {
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequirePrivateIP.ErrCode(), Msg: "CheckRequestIpIsPrivate FromIncomingContext parsing failed"}, nil
	}
	// use x-forwarded-for to check if the client has gone through the proxy
	// in.Get("x-forwarded-for") remoteAddress len is 0 ,then no use proxy
	xForwardedFor := in.Get("x-forwarded-for")
	if len(xForwardedFor) != 0 {
		if !s.cfg.EnableGrpcGateWayPrivateCheck {
			// request is proxy not check private
			return handler(ctx, req)
		}
		ipList := strings.Split(xForwardedFor[0], ",")
		clientIP = ipList[len(ipList)-1]
		clientIP = strings.ReplaceAll(clientIP, " ", "")
	} else {
		pr, ok := rpcPeer.FromContext(ctx)
		if !ok {
			return &carriertypespb.SimpleResponse{Status: backend.ErrRequirePrivateIP.ErrCode(), Msg: "CheckRequestIpIsPrivate FromContext parsing failed"}, nil
		}
		clientIPAndPort := strings.Split(pr.Addr.String(), ":")
		clientIP = clientIPAndPort[0]
	}
	s.cfg.CarrierPrivateIP.PrivateIPCacheCacheLock.RLock()
	defer s.cfg.CarrierPrivateIP.PrivateIPCacheCacheLock.RUnlock()
	if _, ok := s.cfg.PrivateIPCache[clientIP]; !ok {
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequirePrivateIP.ErrCode(), Msg: fmt.Sprintf("%s does not have permission to access %s", clientIP, un.FullMethod)}, nil
	}
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
