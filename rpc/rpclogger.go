package rpc
import "github.com/sirupsen/logrus"
type RPCdebugLogger struct {
	logger logrus.Logger
}

// NewRPCdebugLogger  retrn gRPC LoggerV2 interface implements
func NewRPCdebugLogger(logger *logrus.Logger) *RPCdebugLogger {
	return &RPCdebugLogger{
		logger: *logger,
	}
}

// Info returns
func (rdl *RPCdebugLogger) Info(args ...interface{}) {
	rdl.logger.Info(args...)
}

// Infoln returns
func (rdl *RPCdebugLogger) Infoln(args ...interface{}) {
	rdl.logger.Info(args...)
}

// Infof returns
func (rdl *RPCdebugLogger) Infof(format string, args ...interface{}) {
	rdl.logger.Infof(format, args...)
}

// Warning returns
func (rdl *RPCdebugLogger) Warning(args ...interface{}) {
	rdl.logger.Warn(args...)
}

// Warningln returns
func (rdl *RPCdebugLogger) Warningln(args ...interface{}) {
	rdl.logger.Warn(args...)
}

// Warningf returns
func (rdl *RPCdebugLogger) Warningf(format string, args ...interface{}) {
	rdl.logger.Warnf(format, args...)
}

// Error returns
func (rdl *RPCdebugLogger) Error(args ...interface{}) {
	rdl.logger.Error(args...)
}

// Errorln returns
func (rdl *RPCdebugLogger) Errorln(args ...interface{}) {
	rdl.logger.Error(args...)
}

// Errorf returns
func (rdl *RPCdebugLogger) Errorf(format string, args ...interface{}) {
	rdl.logger.Errorf(format, args...)
}

// Fatal returns
func (rdl *RPCdebugLogger) Fatal(args ...interface{}) {
	rdl.logger.Fatal(args...)
}

// Fatalln returns
func (rdl *RPCdebugLogger) Fatalln(args ...interface{}) {
	rdl.logger.Fatal(args...)
}

// Fatalf logs to fatal level
func (rdl *RPCdebugLogger) Fatalf(format string, args ...interface{}) {
	rdl.logger.Fatalf(format, args...)
}

// V reports whether verbosity level l is at least the requested verbose level.
func (rdl *RPCdebugLogger) V(v int) bool {
	return false
}
