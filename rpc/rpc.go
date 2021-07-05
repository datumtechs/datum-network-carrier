package rpc

import (
	"github.com/RosettaFlow/Carrier-Go/common"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	"google.golang.org/grpc"
	"log"
	"net"
)

type RpcApiServer struct {
	Cfg 				*common.RpcConfig
	GrpcSvr 			*grpc.Server
	b             		Backend
}


func New(cfg *common.RpcConfig, b Backend) *RpcApiServer {
	return &RpcApiServer{Cfg: cfg, b: b}
}
func (svr *RpcApiServer) Start() error {

	lis, err := net.Listen(svr.Cfg.Protocol, svr.Cfg.Ip+":"+svr.Cfg.Port)
	if err != nil {
		log.Fatalf("Failed to start rpc server, listen failed: %v", err)
	}

	var opts []grpc.ServerOption

	s := grpc.NewServer(opts...)

	pb.RegisterYarnServiceServer(s, &yarnServiceServer{b: svr.b})
	pb.RegisterMetaDataServiceServer(s, &metaDataServiceServer{b: svr.b})
	pb.RegisterPowerServiceServer(s, &powerServiceServer{b: svr.b})
	pb.RegisterAuthServiceServer(s, &authServiceServer{b: svr.b})
	pb.RegisterTaskServiceServer(s, &taskServiceServer{b: svr.b})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to start rpc server, serve failed: %v", err)
	}
	svr.GrpcSvr = s
	return nil
}

func (svr *RpcApiServer) Stop () error {
	svr.GrpcSvr.Stop()
	log.Println("Stop grpc server ...")
	return nil
}

func (svr *RpcApiServer) Status () error {

	return nil
}

func (svr *RpcApiServer) GracefulStop () {
	svr.GrpcSvr.GracefulStop()
	log.Println("Stop grpc server graceful ...")
}

func (svr *RpcApiServer) GetServiceInfo() map[string]grpc.ServiceInfo {
	return svr.GrpcSvr.GetServiceInfo()
}


