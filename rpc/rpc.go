package rpc

import (
	rpcapi "github.com/RosettaFlow/Carrier-Go/lib/api"
	"google.golang.org/grpc"
	"log"
	"net"
)

type RpcServer struct {

}

func (svr *RpcServer) Start() {

	lis, err := net.Listen("tcp", "port")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	s := grpc.NewServer(opts...)


	rpcapi.RegisterYarnServiceServer(s, &rpcapi.YarnServiceServer{})

	fmt.Printf("Start Server :%s ... \n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
}