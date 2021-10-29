package gateway

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	rpcpb "github.com/RosettaFlow/Carrier-Go/lib/rpc/v1"
	rpcapipb "github.com/RosettaFlow/Carrier-Go/lib/api"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"net"
	"net/http"
	"time"
)

var _ common.Service = (*Gateway)(nil)

// Gateway is the gRPC gateway to serve HTTP JSON traffic as a proxy and forward
// it to the carrier gRPC server.
type Gateway struct {
	conn                    *grpc.ClientConn
	ctx                     context.Context
	cancel                  context.CancelFunc
	gatewayAddr             string
	remoteAddr              string
	remoteCert              string
	server                  *http.Server
	mux                     *http.ServeMux
	allowedOrigins          []string
	startFailure            error
	enableDebugRPCEndpoints bool
	maxCallRecvMsgSize      uint64
}

// New returns a new gateway server which translates HTTP into gRPC.
// Accepts a context and optional http.ServeMux.
func New(
	ctx context.Context,
	remoteAddress,
	remoteCert,
	gatewayAddress string,
	mux *http.ServeMux,
	allowedOrigins []string,
	enableDebugRPCEndpoints bool,
	maxCallRecvMsgSize uint64,
) *Gateway {
	if mux == nil {
		mux = http.NewServeMux()
	}

	return &Gateway{
		remoteAddr:              remoteAddress,
		remoteCert:              remoteCert,
		gatewayAddr:             gatewayAddress,
		ctx:                     ctx,
		mux:                     mux,
		allowedOrigins:          allowedOrigins,
		enableDebugRPCEndpoints: enableDebugRPCEndpoints,
		maxCallRecvMsgSize:      maxCallRecvMsgSize,
	}
}

// Start the gateway service. This serves the HTTP JSON traffic on the specified port.
func (g *Gateway) Start() error {
	ctx, cancel := context.WithCancel(g.ctx)
	g.cancel = cancel

	log.WithField("address", g.gatewayAddr).Info("Starting JSON-HTTP API")

	conn, err := g.dial(ctx, "tcp", g.remoteAddr)
	if err != nil {
		log.WithError(err).Error("Failed to connect to gRPC server")
		g.startFailure = err
		return err
	}

	g.conn = conn

	gwmux := gwruntime.NewServeMux(
		gwruntime.WithMarshalerOption(
			gwruntime.MIMEWildcard,
			&gwruntime.JSONPb{OrigName: false, EmitDefaults: true},
		),
	)
	handlers := []func(context.Context, *gwruntime.ServeMux, *grpc.ClientConn) error {
		rpcapipb.RegisterAuthServiceHandler,
		rpcapipb.RegisterMetadataServiceHandler,
		rpcapipb.RegisterPowerServiceHandler,
		rpcapipb.RegisterYarnServiceHandler,
		rpcapipb.RegisterTaskServiceHandler,
	}
	if g.enableDebugRPCEndpoints {
		handlers = append(handlers, rpcpb.RegisterDebugHandler)
	}
	for _, f := range handlers {
		if err := f(ctx, gwmux, conn); err != nil {
			log.WithError(err).Error("Failed to start gateway")
			g.startFailure = err
			return err
		}
	}

	g.mux.Handle("/", gwmux)

	g.server = &http.Server{
		Addr:    g.gatewayAddr,
		Handler: newCorsHandler(g.mux, g.allowedOrigins),
	}
	go func() {
		if err := g.server.ListenAndServe(); err != http.ErrServerClosed {
			log.WithError(err).Error("Failed to listen and serve")
			g.startFailure = err
			return
		}
	}()
	return nil
}

// Status of grpc gateway. Returns an error if this service is unhealthy.
func (g *Gateway) Status() error {
	if g.startFailure != nil {
		return g.startFailure
	}
	if s := g.conn.GetState(); s != connectivity.Ready {
		return fmt.Errorf("grpc server is %s", s)
	}
	return nil
}

// Stop the gateway with a graceful shutdown.
func (g *Gateway) Stop() error {
	if g.server != nil {
		shutdownCtx, shutdownCancel := context.WithTimeout(g.ctx, 2*time.Second)
		defer shutdownCancel()
		if err := g.server.Shutdown(shutdownCtx); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				log.Warn("Existing connections terminated")
			} else {
				log.WithError(err).Error("Failed to gracefully shut down server")
			}
		}
	}

	if g.cancel != nil {
		g.cancel()
	}

	return nil
}

// dial the gRPC server.
func (g *Gateway) dial(ctx context.Context, network, addr string) (*grpc.ClientConn, error) {
	switch network {
	case "tcp":
		return g.dialTCP(ctx, addr)
	case "unix":
		return g.dialUnix(ctx, addr)
	default:
		return nil, fmt.Errorf("unsupported network type %q", network)
	}
}

// dialTCP creates a client connection via TCP.
// "addr" must be a valid TCP address with a port number.
func (g *Gateway) dialTCP(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	security := grpc.WithInsecure()
	if len(g.remoteCert) > 0 {
		creds, err := credentials.NewClientTLSFromFile(g.remoteCert, "")
		if err != nil {
			return nil, err
		}
		security = grpc.WithTransportCredentials(creds)
		log.Info("Gateway use TLS to dial.")
	}
	opts := []grpc.DialOption{
		security,
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(int(g.maxCallRecvMsgSize))),
	}

	return grpc.DialContext(
		ctx,
		addr,
		opts...,
	)
}

// dialUnix creates a client connection via a unix domain socket.
// "addr" must be a valid path to the socket.
func (g *Gateway) dialUnix(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	d := func(addr string, timeout time.Duration) (net.Conn, error) {
		return net.DialTimeout("unix", addr, timeout)
	}
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDialer(d),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(int(g.maxCallRecvMsgSize))),
	}
	return grpc.DialContext(ctx, addr, opts...)
}

