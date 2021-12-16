package discovery

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/health/grpc_health_v1"

	"time"
)

var log = logrus.WithField("prefix", "discovery-consul")

type HealthCheck struct{}

type ConsulService struct {
	IP         string //carrier grpc ip
	Port       int    //carrier grpc port
	Tags       []string
	Name       string
	Id         string
	Interval   int
	Deregister int
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
func (ca *ConnectConsul) RegisterDiscoveryService() error {
	// ca is consul "IP:PORT"

	interval := time.Duration(ca.Interval) * time.Second
	deregister := time.Duration(ca.Deregister) * time.Minute

	//register consul
	reg := &api.AgentServiceRegistration{
		ID:      ca.Id,
		Name:    ca.Name,
		Tags:    ca.Tags,
		Port:    ca.Port,
		Address: ca.IP,
		Check: &api.AgentServiceCheck{
			Interval:                       interval.String(),
			GRPC:                           fmt.Sprintf("%v:%v/%v", ca.IP, ca.Port, ca.Name),
			DeregisterCriticalServiceAfter: deregister.String(),
		},
	}

	log.Info("Begin RegisterDiscoveryService", ca)
	if err := ca.Client.Agent().ServiceRegister(reg); err != nil {
		log.Errorf("Service Register error %v\n", err)
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
