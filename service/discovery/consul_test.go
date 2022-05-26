package discovery

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"gotest.tools/assert"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &HelloReply{Message: "Hello " + in.Name}, nil
}

func TestConsulClientFlow(t *testing.T) {

	var (
		consulIp           = "10.1.1.10"
		consulPort         = 8500
		grpcIp             = "192.168.21.187"
		grpcPort           = "9999"
		carrierServiceName = "carrier"
		configKey          = "datum-network/via_ip_port"
		configValue        = "192.168.10.111_12345"
	)

	ctx, cancelFn := context.WithCancel(context.Background())
	go func(cancelFn context.CancelFunc) {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
		if err != nil {
			t.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		RegisterGreeterServer(s, &server{})
		grpc_health_v1.RegisterHealthServer(s, &HealthCheck{})

		cancelFn()

		if err := s.Serve(lis); err != nil {
			t.Fatalf("failed to serve: %v", err)
		}
	}(cancelFn)
	<-ctx.Done()

	time.Sleep(1 * time.Second)
	// REGISTER GRPC SERVER TO CONSUL
	serverInfo := &ConsulService{
		ServiceIP:   grpcIp,
		ServicePort: grpcPort,
		Tags:        []string{carrierServiceName},
		Name:        carrierServiceName,
		Id:          strings.Join([]string{"carrierService", grpcIp, grpcPort}, "_"),
		Interval:    2,
		Deregister:  3,
	}
	conn := NewConsulClient(serverInfo, consulIp, consulPort)

	err := conn.RegisterService2DiscoveryCenter()
	assert.NilError(t, err, "register service succeed")

	// QUERY SERVICE
	result, err := conn.QueryServiceInfoByFilter(fmt.Sprintf(ConsulServiceTagFuzzQueryExpression, carrierServiceName))
	assert.NilError(t, err, "query services succeed")
	var has bool
	for _, value := range result {
		if value.Address == grpcIp && strconv.Itoa(value.Port) == grpcPort {
			has = true
		}
	}
	assert.Equal(t, has, true, "find carrier service succeed")

	//KV PUT
	err = conn.PutKV(configKey, []byte(configValue))
	assert.NilError(t, err, "put config succeed")
	//KV GET
	value, err := conn.GetKV(configKey, nil)
	assert.NilError(t, err, "query config succeed")
	assert.Equal(t, value, configValue, "find config succeed")
	assert.NilError(t, conn.DeregisterService2DiscoveryCenter(), "deregister service succeed")
}
