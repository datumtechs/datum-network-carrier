package debug

import (
	"context"
	carrierrpcdebugpbv1 "github.com/datumtechs/datum-network-carrier/pb/carrier/rpc/debug/v1"
	"github.com/datumtechs/datum-network-carrier/p2p"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/golang/protobuf/ptypes/empty"
	golog "github.com/ipfs/go-log/v2"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
)

// Server defines a server implementation of the gRPC Debug service,
// providing RPC endpoints for runtime debugging of a node, this server is
// gated behind the feature flag --enable-debug-rpc-endpoints.
type Server struct {
	PeerManager        p2p.PeerManager
	PeersFetcher       p2p.PeersProvider
	DebugAPI           DebugBackend
}

// SetLoggingLevel of a beacon node according to a request type,
// either INFO, DEBUG, or TRACE.
func (ds *Server) SetLoggingLevel(_ context.Context, req *carrierrpcdebugpbv1.LoggingLevelRequest) (*empty.Empty, error) {
	var verbosity string
	switch req.Level {
	case carrierrpcdebugpbv1.LoggingLevelRequest_INFO:
		verbosity = "info"
	case carrierrpcdebugpbv1.LoggingLevelRequest_DEBUG:
		verbosity = "debug"
	case carrierrpcdebugpbv1.LoggingLevelRequest_TRACE:
		verbosity = "trace"
	default:
		return nil, status.Error(codes.InvalidArgument, "Expected valid verbosity level as argument")
	}
	level, err := logrus.ParseLevel(verbosity)
	if err != nil {
		return nil, status.Error(codes.Internal, "Could not parse verbosity level")
	}
	logrus.SetLevel(level)
	if level == logrus.TraceLevel {
		// Libp2p specific logging.
		golog.SetAllLoggers(golog.LevelDebug)
		// Geth specific logging.
		glogger := gethlog.NewGlogHandler(gethlog.StreamHandler(os.Stderr, gethlog.TerminalFormat(true)))
		glogger.Verbosity(gethlog.LvlTrace)
		gethlog.Root().SetHandler(glogger)
	}
	return &empty.Empty{}, nil
}
