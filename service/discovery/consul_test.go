package discovery

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

type server struct{}

const (
	consulIp   = "192.168.235.155"
	consulPort = "8500"
	grpcIp     = "192.168.21.188"
	grpcPort   = "9999"
)

func (s *server) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &HelloReply{Message: "Hello " + in.Name}, nil
}

func TestNew(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		RegisterGreeterServer(s, &server{})
		grpc_health_v1.RegisterHealthServer(s, &HealthCheck{})
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
		wg.Done()
	}()
	// REGISTER GRPC SERVER TO CONSUL
	serverInfo := &ConsulService{
		ServiceIP:   grpcIp,
		ServicePort: grpcPort,
		Tags:        []string{"carrier"},
		Name:        "carrier",
		Id:          strings.Join([]string{"carrierService", grpcIp, grpcPort}, "_"),
		Interval:    2,
		Deregister:  3,
	}
	conn := New(serverInfo, consulIp, consulPort)
	err := conn.RegisterDiscoveryService()
	if err != nil {
		fmt.Println("register successful")
	}
	// QUERY SERVICE
	result, _ := conn.QueryServiceInfoByFilter("carrier in Tags")
	for _, value := range result {
		fmt.Println(value.Address, value.Port)
	}

	//KV PUT
	err = conn.PutKV("metis/via_ip_port", []byte("192.168.10.111_12345"))
	if err != nil {
		panic(err)
	}
	//KV GET
	value, err := conn.GetKV("metis/via_ip_port", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
	time.Sleep(10 * time.Second)

	_ = conn.DeregisterDiscoveryService()
	wg.Wait()
}
