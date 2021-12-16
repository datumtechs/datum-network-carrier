package discovery

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/health/grpc_health_v1"
	"strconv"
	"strings"

	"time"
)

var log = logrus.WithField("prefix", "discovery-consul")

type HealthCheck struct{}

type ConsulService struct {
	ServiceIP   string   //carrier grpc server ip
	ServicePort string   //carrier grpc server port
	Tags        []string // carrier service []tags on discovery center
	Name        string   // carrier service name on discovery center
	Id          string   // carrier service id on discovery center
	Interval    int      // Health check interval between service discovery center and this service (default: 3s)
	Deregister  int      // When the service leaves, the service discovery center removes the service information (default: 1min)
}
type ConnectConsul struct {
	Client *api.Client
	ConsulService
}

func (h *HealthCheck) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	/*Implement the health check interface, where the health status is directly returned,
	and there can also be more complex health check strategies, such as returning based on server load.
	grpc health check detail see https://github.com/grpc/grpc/blob/master/doc/health-checking.md
	*/
	log.Debugf("discovery consul health checking")
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (h *HealthCheck) Watch(req *grpc_health_v1.HealthCheckRequest, w grpc_health_v1.Health_WatchServer) error {
	return nil
}

func New(consulSvr *ConsulService, consulIp, consulPort string) *ConnectConsul {
	if "" == consulIp || "" == consulPort {
		log.Fatalf("consulIp is %s,consulPort is %s", consulIp, consulPort)
	}

	// consulIp is consul server ip,consulPort is consul server address
	consulConfig := api.DefaultConfig()
	consulConfig.Address = fmt.Sprintf("%s:%s", consulIp, consulPort)
	client, err := api.NewClient(consulConfig)
	if err != nil {
		log.Errorf("NewClient error\n%v", err)
		return nil
	}
	connectCon := &ConnectConsul{
		Client:        client,
		ConsulService: *consulSvr,
	}
	return connectCon
}
func (ca *ConnectConsul) CheckDoesItExistCarrierRegister(tags []string) bool {
	for _, tag := range tags {
		exits, err := ca.QueryServiceInfoByFilter(fmt.Sprintf("%s in Tags", tag))
		if err == nil && len(exits) != 0 {
			return true
		}
	}
	return false
}
func (ca *ConnectConsul) RegisterDiscoveryService() error {
	// ca is consul "IP:PORT"
	servicePort, _ := strconv.Atoi(ca.ServicePort)
	if "" == ca.Id {
		ca.Id = strings.Join([]string{"carrierService", ca.ServiceIP, ca.ServicePort}, "_")
	}
	checkResult := ca.CheckDoesItExistCarrierRegister(ca.Tags)
	if checkResult {
		log.Fatalf("alerady exits carrier register! there can only be one carrier.")
	}
	//register consul
	reg := &api.AgentServiceRegistration{
		ID:      ca.Id,
		Name:    ca.Name,
		Tags:    ca.Tags,
		Port:    servicePort,
		Address: ca.ServiceIP,
		Check: &api.AgentServiceCheck{
			Interval:                       (time.Duration(ca.Interval) * time.Millisecond).String(),
			GRPC:                           fmt.Sprintf("%s:%d/%s", ca.ServiceIP, servicePort, ca.Name),
			DeregisterCriticalServiceAfter: (time.Duration(ca.Deregister) * time.Millisecond).String(),
		},
	}

	log.Infof("Start RegisterDiscoveryService, %v", ca)
	if err := ca.Client.Agent().ServiceRegister(reg); err != nil {
		log.WithError(err).Errorf("Failed to register carrier service to discovery center")
		return err
	}
	return nil
}

func (ca *ConnectConsul) QueryServiceInfoByFilter(filterStr string) (map[string]*api.AgentService, error) {
	/*
	   filter_str detail please see https://www.consul.io/api-docs/agent/service#ns
	   filter_str example:
	                     by service name query filter_str is Service==jobNode
	                     by service tag  query filter_str is via in Tags
	                     by service id   query filter_str is ID=="jobNode_192.168.21.188_50053
	*/
	data, err := ca.Client.Agent().ServicesWithFilter(filterStr)
	return data, err
}
func (ca *ConnectConsul) GetKV(key string, q *api.QueryOptions) (string, error) {
	kv := ca.Client.KV()
	pair, _, err := kv.Get(key, q)
	if err != nil {
		return "", err
	}
	return string(pair.Value), nil
}

func (ca *ConnectConsul) PutKV(key string, value []byte) error {
	p := &api.KVPair{Key: key, Value: value}
	kv := ca.Client.KV()
	_, err := kv.Put(p, nil)
	if err != nil {
		return err
	}
	return nil
}

func (ca *ConnectConsul) Stop() error {
	serviceId := strings.Join([]string{"carrierService", ca.ServiceIP, ca.ServicePort}, "_")
	log.Infof("Stop service,it serive id is %s", serviceId)
	result, err := ca.QueryServiceInfoByFilter(fmt.Sprintf("ID==\"%s\"", serviceId))
	if err != nil {
		return err
	}
	if len(result) != 0 {
		if err := ca.Client.Agent().ServiceDeregister(serviceId); err != nil {
			return err
		}
	}
	return nil
}
