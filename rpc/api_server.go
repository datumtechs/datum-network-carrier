package rpc


import (
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
)

type YarnServiceServer {
pb.UnimplementedYarnServiceServer



}


type routeGuideServer struct {
	pb.UnimplementedYarnServiceServer
	savedFeatures []*pb.Feature // read-only after initialized

	mu         sync.Mutex // protects routeNotes
	routeNotes map[string][]*pb.RouteNote
}