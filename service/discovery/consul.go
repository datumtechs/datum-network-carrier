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

const (
	// service Id separator
	ConsulServiceIdSeparator = "_"
	// service Id prefix
	CarrierConsulServiceIdPrefix        = "carrierService"
	JobNodeConsulServiceIdPrefix        = "jobService"
	DataNodeConsulServiceIdPrefix        = "dataService"

	// tag
	JobNodeConsulServiceTag             = "job"
	DataNodeConsulServiceTag            = "data"

	// expression
	ConsulServiceTagFuzzQueryExpression = "%s in Tags"
	ConsulServiceIdEquelQueryExpression = `ID=="%s"`
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
		log.Fatalf("Failed to check consul server ip: %s, consul server port %s", consulIp, consulPort)
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
func (ca *ConnectConsul) IsExistSerivces(tags []string) bool {
	for _, tag := range tags {
		filterReq := fmt.Sprintf(ConsulServiceTagFuzzQueryExpression, tag)
		exits, err := ca.QueryServiceInfoByFilter(filterReq)
		if nil != err {
			log.WithError(err).Fatalf("Failed to query carrier service from discovery center when check discovery service whether exist, filter req: %s", filterReq)
		}
		if len(exits) != 0 {
			return true
		}
	}
	return false
}
func (ca *ConnectConsul) RegisterService2DiscoveryCenter() error {
	if "" == ca.Id {
		ca.Id = strings.Join([]string{CarrierConsulServiceIdPrefix, ca.ServiceIP, ca.ServicePort}, ConsulServiceIdSeparator)
	}
	servicePort, _ := strconv.Atoi(ca.ServicePort)
	if ca.IsExistSerivces(ca.Tags) {
		log.Fatalf("Alerady exits carrier service on consule server, it must be unique")
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

	log.Infof("Start register carrier service, %v", ca)
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
	return ca.Client.Agent().ServicesWithFilter(filterStr)
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

func (ca *ConnectConsul) DeregisterService2DiscoveryCenter() error {
	serviceId := strings.Join([]string{CarrierConsulServiceIdPrefix, ca.ServiceIP, ca.ServicePort}, ConsulServiceIdSeparator)
	filterReq := fmt.Sprintf(ConsulServiceIdEquelQueryExpression, serviceId)
	result, err := ca.QueryServiceInfoByFilter(filterReq)
	if err != nil {
		log.WithError(err).Errorf("Failed to query carrier service from discovery center before deregister discovery service, filter req: %s", filterReq)
		return err
	}
	log.Infof("Start deregister carrier service, serviceId: %s", serviceId)
	if len(result) != 0 {
		if err := ca.Client.Agent().ServiceDeregister(serviceId); err != nil {
			log.WithError(err).Errorf("Failed to deregister carrier service, serviceId: %s", serviceId)
			return err
		}
	}
	return nil
}

func (ca *ConnectConsul) QueryDataNodeServices() (map[string]*api.AgentService, error) {
	return ca.QueryServiceInfoByFilter(fmt.Sprintf(ConsulServiceTagFuzzQueryExpression, DataNodeConsulServiceTag))
}

func (ca *ConnectConsul) QueryJobNodeServices() (map[string]*api.AgentService, error) {
	return ca.QueryServiceInfoByFilter(fmt.Sprintf(ConsulServiceTagFuzzQueryExpression, JobNodeConsulServiceTag))
}
