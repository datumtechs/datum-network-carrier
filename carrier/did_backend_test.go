package carrier

import (
	"context"
	"github.com/datumtechs/datum-network-carrier/grpclient"
	"github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	"github.com/datumtechs/datum-network-carrier/service/discovery"
	"strconv"
	"testing"
)

func Test_DialRemoteCarrier(t *testing.T) {

	ctx, cancelFn := context.WithTimeout(context.Background(), grpclient.DefaultGrpcDialTimeout)
	defer cancelFn()

	issuerUrl := "192.168.9.154:10033"
	conn, err := grpclient.DialContext(ctx, issuerUrl, true)
	if err != nil {
		t.Fatalf("cannot dial up url: %+v", err)
	}
	defer conn.Close()

	client := api.NewVcServiceClient(conn)
	t.Log(client)

}

func Test_ApplyLocalVC(t *testing.T) {
	/*s := CarrierAPIBackend{}
	s.ApplyVCLocal()*/
}

func Test_findAdminService(t *testing.T) {
	consulManager := newConsulManager()
	//查找本地admin服务端口
	adminService, err := consulManager.QueryAdminService()
	if err != nil {
		t.Fatalf("cannot find local admin RPC service: %+v", err)
		return
	}
	if adminService == nil {
		t.Fatalf("cannot find local admin RPC service: %+v", err)
		return
	}
	log.Debugf("adminService info:%+v", adminService)

	ctx, cancelFn := context.WithTimeout(context.Background(), grpclient.DefaultGrpcDialTimeout)
	defer cancelFn()

	adminServiceUrl := adminService.Address + ":" + strconv.Itoa(adminService.Port)
	conn, err := grpclient.DialContext(ctx, adminServiceUrl, true)
	defer conn.Close()

	if err != nil {
		t.Fatalf("failed to dial admin service: %+v", err)
		return
	}
	t.Log(conn)
}

func newConsulManager() *discovery.ConnectConsul {
	consulManager := discovery.NewConsulClient(&discovery.ConsulService{
		ServiceIP:   "192.168.9.156",
		ServicePort: "8500",
		Tags:        DefaultConfig.DiscoverServiceConfig.DiscoveryServerTags,
		Name:        DefaultConfig.DiscoverServiceConfig.DiscoveryServiceName,
		Id:          DefaultConfig.DiscoverServiceConfig.DiscoveryServiceId,
		Interval:    DefaultConfig.DiscoverServiceConfig.DiscoveryServiceHealthCheckInterval,
		Deregister:  DefaultConfig.DiscoverServiceConfig.DiscoveryServiceHealthCheckDeregister,
	},
		"192.168.9.156",
		8500,
	)

	return consulManager
}
