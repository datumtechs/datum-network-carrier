package carrier

import (
	"fmt"
	"github.com/Metisnetwork/Metis-Carrier/core/rawdb"
	pb "github.com/Metisnetwork/Metis-Carrier/lib/api"
	"github.com/Metisnetwork/Metis-Carrier/service/discovery"
	"github.com/Metisnetwork/Metis-Carrier/types"
	"strconv"
	"strings"
	"time"
)

const (
	defaultRefreshResourceNodesInternal = 6 * time.Second
)

func (s *Service) loop() {

	refreshResourceNodesTicker := time.NewTicker(defaultRefreshResourceNodesInternal) // 6 s

	for {

		select {

		case <-refreshResourceNodesTicker.C:

			if err := s.refreshResourceNodes(); nil != err {
				log.WithError(err).Errorf("Failed to call service.refreshResourceNodes() on service.loop()")
			}

		case <-s.quit:
			log.Info("Stopped carrier service ...")
			refreshResourceNodesTicker.Stop()
			return
		}
	}
}

func (s *Service) initServicesWithDiscoveryCenter() error {

	log.Infof("Start init Services with discoveryCenter...")

	if nil == s.carrierDB {
		return fmt.Errorf("init services with discorvery center failed, carrier local db is nil")
	}

	if nil == s.consulManager {
		return fmt.Errorf("init services with discorvery center failed, consul manager is nil")
	}

	// 1. fetch datacenter ip&port config from discorvery center
	datacenterIpAndPort, err := s.consulManager.GetKV(discovery.DataCenterConsulServiceAddressKey, nil)
	if nil != err {
		return fmt.Errorf("query datacenter IP and PORT KVconfig failed from discovery center, %s", err)
	}
	if "" == datacenterIpAndPort {
		return fmt.Errorf("datacenter IP and PORT KVconfig is empty from discovery center")
	}
	configArr := strings.Split(datacenterIpAndPort, discovery.ConsulServiceIdSeparator)

	// datacenter address config in consul server
	// 	key: metis/dataCenter_ip_port
	//  value: pi_port
	if len(configArr) != 2 {
		return fmt.Errorf("datacenter IP and PORT lack one on KVconfig from discovery center")
	}

	datacenterIP, datacenterPortStr := configArr[0], configArr[1]

	datacenterPort, err := strconv.Atoi(datacenterPortStr)
	if nil != err {
		return fmt.Errorf("invalid datacenter Port from discovery center, %s", err)
	}

	if err := s.carrierDB.SetConfig(&types.CarrierChainConfig{
		GrpcUrl: datacenterIP,
		Port:    uint64(datacenterPort),
	}); nil != err {
		return fmt.Errorf("connot setConfig of carrierDB, %s", err)
	}

	// 2. register carrier service to discorvery center
	if err := s.consulManager.RegisterService2DiscoveryCenter(); nil != err {
		return fmt.Errorf("register discovery service to discovery center failed, %s", err)
	}
	log.Infof("Succeed init Services with discoveryCenter...")
	return nil
}

func (s *Service) refreshResourceNodes() error {

	if nil == s.carrierDB {
		return fmt.Errorf("refresh internal resource nodes failed, carrier local db is nil")
	}

	if nil == s.consulManager {
		return fmt.Errorf("refresh internal resource nodes failed, consul manager is nil")
	}

	// 1. check whether the current organization is connected to the network.
	//    if it is not connected to the network, without subsequent processing and short circuit.
	identity, err := s.carrierDB.QueryIdentity()
	if nil != err {
		//log.WithError(err).Warnf("Failed to query local identity")
		return nil
	}

	// 2. fetch taskGateWay ip&port config from discorvery center
	taskGateWayIpAndPort, err := s.consulManager.GetKV(discovery.TaskGateWayConsulServiceAddressKey, nil)
	if nil != err {
		return fmt.Errorf("query taskGateWay IP and PORT KVconfig from discovery center failed, %s", err)
	}
	if "" == taskGateWayIpAndPort {
		return fmt.Errorf("taskGateWay IP and PORT KVconfig is empty from discovery center")
	}
	configArr := strings.Split(taskGateWayIpAndPort, discovery.ConsulServiceIdSeparator)

	// taskGateWay address config in consul server
	// 	key: metis/glacier2_ip_port
	//  value: pi_port
	if len(configArr) != 2 {
		return fmt.Errorf("taskGateWay IP and PORT lack one on KVconfig from discovery center")
	}

	taskGateWaylIP, taskGateWaylPort := configArr[0], configArr[1]

	// ##########################
	// ##########################
	// ABOUT JOBNODE SERVICES
	// ##########################
	// ##########################

	jobNodeServices, err := s.consulManager.QueryJobNodeServices()
	if nil != err {
		log.WithError(err).Warnf("query jobNodeServices failed from discovery center on service.refreshResourceNodes()")
	} else {

		jobNodeDBCache := make(map[string]*pb.YarnRegisteredPeerDetail, 0)
		// load stored jobNode
		jobNodeList, err := s.carrierDB.QueryRegisterNodeList(pb.PrefixTypeJobNode)
		if nil != err && rawdb.IsNoDBNotFoundErr(err) {
			log.WithError(err).Warnf("query jobNodes failed from local db on service.refreshResourceNodes()")
		} else {
			for _, node := range jobNodeList {
				jobNodeDBCache[node.GetId()] = node
			}

			for _, jobNodeService := range jobNodeServices {

				if node, ok := jobNodeDBCache[jobNodeService.ID]; !ok { // add new registered jobNode service

					if err = s.resourceManager.AddDiscoveryJobNodeResource(
						identity, jobNodeService.ID, jobNodeService.Address, strconv.Itoa(jobNodeService.Port),
						taskGateWaylIP, taskGateWaylPort); nil != err {
						log.WithError(err).Errorf("Failed to store a new jobNode on service.refreshResourceNodes(), jobNodeServiceId: {%s}, jobNodeService: {%s:%d}",
							jobNodeService.ID, jobNodeService.Address, jobNodeService.Port)
						continue
					}

				} else {

					if err = s.resourceManager.UpdateDiscoveryJobNodeResource(
						identity, jobNodeService.ID, jobNodeService.Address, strconv.Itoa(jobNodeService.Port),
						taskGateWaylIP, taskGateWaylPort, node); nil != err {
						log.WithError(err).Errorf("Failed to update a old jobNode on service.refreshResourceNodes(), jobNodeServiceId: {%s}, jobNodeService: {%s:%d}",
							jobNodeService.ID, jobNodeService.Address, jobNodeService.Port)
						continue
					}

					delete(jobNodeDBCache, jobNodeService.ID)
				}
			}

			// delete old deregistered jobNode service
			for jobNodeId, node := range jobNodeDBCache {

				if err = s.resourceManager.RemoveDiscoveryJobNodeResource(
					identity, jobNodeId, node.GetInternalIp(), node.GetInternalPort(),
					taskGateWaylIP, taskGateWaylPort, node); nil != err {
					log.WithError(err).Errorf("Failed to removed a old jobNode on service.refreshResourceNodes(), jobNodeServiceId: {%s}, jobNodeService: {%s:%s}",
						jobNodeId, node.GetInternalIp(), node.GetInternalPort())
					continue
				}

			}
		}
	}

	// ##########################
	// ##########################
	// ABOUT DATANODE SERVICES
	// ##########################
	// ##########################
	dataNodeServices, err := s.consulManager.QueryDataNodeServices()
	if nil != err {
		log.WithError(err).Warnf("query dataNodeServices from discovery center failed on service.refreshResourceNodes()")
	} else {

		dataNodeDBCache := make(map[string]*pb.YarnRegisteredPeerDetail, 0)
		// load stored dataNode
		dataNodeList, err := s.carrierDB.QueryRegisterNodeList(pb.PrefixTypeDataNode)
		if nil != err && rawdb.IsNoDBNotFoundErr(err) {
			log.WithError(err).Warnf("query dataNodes from local db failed on service.refreshResourceNodes()")
		} else {
			for _, node := range dataNodeList {
				dataNodeDBCache[node.GetId()] = node
			}

			for _, dataNodeService := range dataNodeServices {
				if node, ok := dataNodeDBCache[dataNodeService.ID]; !ok { // add new registered dataNode service

					if err = s.resourceManager.AddDiscoveryDataNodeResource(
						identity, dataNodeService.ID, dataNodeService.Address, strconv.Itoa(dataNodeService.Port),
						taskGateWaylIP, taskGateWaylPort); nil != err {
						log.WithError(err).Errorf("Failed to store a new dataNode on service.refreshResourceNodes(), dataNodeServiceId: {%s}, dataNodeService: {%s:%d}",
							dataNodeService.ID, dataNodeService.Address, dataNodeService.Port)
						continue
					}

				} else {

					if err = s.resourceManager.UpdateDiscoveryDataNodeResource(
						identity, dataNodeService.ID, dataNodeService.Address, strconv.Itoa(dataNodeService.Port),
						taskGateWaylIP, taskGateWaylPort, node); nil != err {
						log.WithError(err).Errorf("Failed to update a old dataNode on service.refreshResourceNodes(), dataNodeServiceId: {%s}, dataNodeService: {%s:%d}",
							dataNodeService.ID, dataNodeService.Address, dataNodeService.Port)
						continue
					}

					delete(dataNodeDBCache, dataNodeService.ID)
				}
			}

			// delete old deregistered dataNode service
			for dataNodeId, node := range dataNodeDBCache {

				if err = s.resourceManager.RemoveDiscoveryDataNodeResource(
					identity, dataNodeId, node.GetInternalIp(), node.GetInternalPort(),
					taskGateWaylIP, taskGateWaylPort, node); nil != err {
					log.WithError(err).Errorf("Failed to removed a old dataNode on service.refreshResourceNodes(), dataNodeServiceId: {%s}, dataNodeService: {%s:%s}",
						dataNodeId, node.GetInternalIp(), node.GetInternalPort())
					continue
				}

			}
		}
	}
	return nil
}