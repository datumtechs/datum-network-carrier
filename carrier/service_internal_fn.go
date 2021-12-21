package carrier

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	"github.com/RosettaFlow/Carrier-Go/grpclient"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/lib/fighter/computesvc"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/params"
	"github.com/RosettaFlow/Carrier-Go/service/discovery"
	"github.com/RosettaFlow/Carrier-Go/types"
	"strconv"
	"strings"
	"time"
)

const (
	defaultRefreshResourceNodesInternal = 30 * time.Millisecond
)

func (s *Service) loop() {

	refreshResourceNodesTicker := time.NewTicker(defaultRefreshResourceNodesInternal) // 30 ms

	for {

		select {

		case <-refreshResourceNodesTicker.C:

			s.refreshResourceNodes()

		case <-s.quit:
			log.Info("Stopped carrier service ...")
			return
		}
		}
	}
}

func (s *Service) initServicesWithDiscoveryCenter() error {

	if nil == s.carrierDB {
		return fmt.Errorf("Failed to init services with discorvery center, carrier local db is nil")
	}

	if nil == s.consulManager {
		return fmt.Errorf("Failed to init services with discorvery center, consul manager is nil")
	}

	// First: register carrier service to discorvery center
	if err := s.consulManager.RegisterService2DiscoveryCenter(); nil != err {
		return fmt.Errorf("Failed to register discovery service to discovery center, %s", err)
	}

	// second: fetch datacenter and via and all jobnodes and datanodes from discorvery center
	datacenterIP, err := s.consulManager.GetKV("datacenterIP", nil)
	if nil == err {
		return fmt.Errorf("Failed to query datacenter IP from discovery center, %s", err)
	}
	datacenterPortStr, err := s.consulManager.GetKV("datacenterPort", nil)
	if nil == err {
		return fmt.Errorf("Failed to query datacenter Port from discovery center, %s", err)
	}
	datacenterPort, err := strconv.Atoi(datacenterPortStr)
	if nil == err {
		return fmt.Errorf("invalid datacenter Port from discovery center, %s", err)
	}

	if err := s.carrierDB.SetConfig( &params.CarrierChainConfig{
		GrpcUrl: datacenterIP,
		Port:    uint64(datacenterPort),
	}); nil != err {
		return fmt.Errorf("connot setConfig of carrierDB, %s", err)
	}

	return nil
}


func (s *Service) refreshResourceNodes() error {

	if nil == s.carrierDB {
		return fmt.Errorf("Failed to refresh internal resource nodes, carrier local db is nil")
	}

	if nil == s.consulManager {
		return fmt.Errorf("Failed to refresh internal resource nodes, consul manager is nil")
	}

	viaExternalIP, err := s.consulManager.GetKV("viaExternalIP", nil)
	if nil == err {
		return fmt.Errorf("Failed to query via external IP from discovery center, %s", err)
	}
	viaExternalPort, err := s.consulManager.GetKV("viaExternalPort", nil)
	if nil == err {
		return fmt.Errorf("Failed to query via external Port from discovery center, %s", err)
	}


	// QUERY SERVICE
	result, _ := conn.QueryServiceInfoByFilter("carrier in Tags")
	for _, value := range result {
		fmt.Println(value.Address, value.Port)
	}


	client, err := grpclient.NewJobNodeClientWithConn(s.ctx, fmt.Sprintf("%s:%s", node.GetInternalIp(), node.GetInternalPort()), node.GetId())
	if err != nil {
		return pb.ConnState_ConnState_UnConnected, fmt.Errorf("connect new jobNode failed, %s", err)
	}





	storeLocalResourceFn := func (jobNodeId string, jobNodeStatus *computesvc.GetStatusReply) error {

		// store into local db
		if err := s.carrierDB.InsertLocalResource(types.NewLocalResource(&libtypes.LocalResourcePB{
		//IdentityId: identity.GetIdentityId(),
		//NodeId:     identity.GetNodeId(),
		//NodeName:   identity.GetNodeName(),
		JobNodeId:  jobNodeId,
		DataId:     "", // can not own powerId now, because power have not publish
		// the status of data, N means normal, D means deleted.
		DataStatus: apicommonpb.DataStatus_DataStatus_Normal,
		// resource status, eg: create/release/revoke
		State: apicommonpb.PowerState_PowerState_Created,
		// unit: byte
		TotalMem: jobNodeStatus.GetTotalMemory(),
		UsedMem:  0,
		// number of cpu cores.
		TotalProcessor: jobNodeStatus.GetTotalCpu(),
		UsedProcessor:  0,
		// unit: byte
		TotalBandwidth: jobNodeStatus.GetTotalBandwidth(),
		UsedBandwidth:  0,
		TotalDisk:      jobNodeStatus.GetTotalDisk(),
		UsedDisk:       0,
	})); nil != err {
		log.WithError(err).Errorf("Failed to store power to local on MessageHandler with broadcast, jobNodeId: {%s}",
		jobNodeId)
		return err
	}
		return nil
	}













	jobNodeCache := make(map[string]struct{}, 0)

	// load stored jobNode and dataNode
	jobNodeList, err := s.carrierDB.QueryRegisterNodeList(pb.PrefixTypeJobNode)
	if nil != err && rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Warnf("query jobNodes from local db failed")
	} else if nil == err {
		for _, node := range jobNodeList {
			jobNodeCache[node.GetId()] = struct{}{}
		}
	}

	jobNodeServices, err := s.consulManager.QueryJobNodeServices()
	if nil != err {
		log.WithError(err).Warnf("query jobNodeServices from discovery center failed")
	} else {
		for _, jobNodeService := range jobNodeServices {
			if _, ok := jobNodeCache[jobNodeService.ID]; !ok { // add new registered jobNode service
				client, err := grpclient.NewJobNodeClient(s.ctx, fmt.Sprintf("%s:%s", jobNodeService.Address, jobNodeService.Port), jobNodeService.ID)
				if err != nil {
					log.WithError(err).Errorf("connect new jobNode failed, jobNodeServiceId: {%s}, jobNodeService: {%s:%s}",
						jobNodeService.ID, jobNodeService.Address, jobNodeService.Port)
					continue
				}
				s.resourceClientSet.StoreJobNodeClient(jobNodeService.ID, client)

				jobNodeStatus, err := client.GetStatus()
				if err != nil {
					log.WithError(err).Errorf("connect jobNode query status failed, jobNodeServiceId: {%s}, jobNodeService: {%s:%s}",
						jobNodeService.ID, jobNodeService.Address, jobNodeService.Port)
					continue
				}
				// add resource usage first, but not own power now (mem, proccessor, bandwidth)
				if err = storeLocalResourceFn(jobNodeService.ID, jobNodeStatus); nil != err {
					log.WithError(err).Errorf("store jobNode local resource failed, jobNodeServiceId: {%s}, jobNodeService: {%s:%s}",
						jobNodeService.ID, jobNodeService.Address, jobNodeService.Port)
					continue
				}

				// build jobNode info that was need to store local db
				node := &pb.YarnRegisteredPeerDetail{
					InternalIp:   jobNodeService.Address,
					InternalPort: strconv.Itoa(jobNodeService.Port),
					ExternalIp:   viaExternalIP,
					ExternalPort: viaExternalPort,
					ConnState:    pb.ConnState_ConnState_UnConnected,
				}
				node.Id = strings.Join([]string{discovery.JobNodeConsulServiceIdPrefix, jobNodeService.Address, strconv.Itoa(jobNodeService.Port)}, discovery.ConsulServiceIdSeparator)

				// store jobNode ip port into local db
				if err = s.carrierDB.SetRegisterNode(pb.PrefixTypeJobNode, node); err != nil {
					log.WithError(err).Errorf("Store registerNode to db failed, jobNodeServiceId: {%s}, jobNodeService: {%s:%s}",
						jobNodeService.ID, jobNodeService.Address, jobNodeService.Port)
				}
			} else {
				delete(jobNodeCache, jobNodeService.ID)
			}
		}
	}

	if len(jobNodeCache) != 0 { // delete old deregistered jobNode service
		for jobNodeId, _ := range jobNodeCache {
			// The published jobNode cannot be updated directly
			resourceTable, err := s.carrierDB.QueryLocalResourceTable(jobNodeId)
			if rawdb.IsNoDBNotFoundErr(err) {
				log.WithError(err).Errorf("query local power resource on old jobNode failed")
				continue
			}

			if nil != resourceTable {
				log.Warnf("still have the published computing power information on old jobNode on RemoveRegisterNode, %s", resourceTable.String())
				continue
			}

			if client, ok := s.resourceClientSet.QueryJobNodeClient(jobNodeId); ok {
				client.Close()
				s.resourceClientSet.RemoveJobNodeClient(jobNodeId)
			}

			// remove jobNode local resource
			if err = s.carrierDB.RemoveLocalResource(jobNodeId); nil != err {
				log.WithError(err).Errorf("remove jobNode local resource failed")
			}
		}
	}


	if err == nil {
		for _, node := range jobNodeList {
			client, err := grpclient.NewJobNodeClient(ctx, fmt.Sprintf("%s:%s", node.GetInternalIp(), node.GetInternalPort()), node.GetId())
			if err == nil {
				s.resourceClientSet.StoreJobNodeClient(node.GetId(), client)
			}
		}
	}
	dataNodeList, err := s.carrierDB.QueryRegisterNodeList(pb.PrefixTypeDataNode)
	if err == nil {
		for _, node := range dataNodeList {
			client, err := grpclient.NewDataNodeClient(ctx, fmt.Sprintf("%s:%s", node.GetInternalIp(), node.GetInternalPort()), node.GetId())
			if err == nil {
				s.resourceClientSet.StoreDataNodeClient(node.GetId(), client)
			}
		}
	}



	return nil
}

