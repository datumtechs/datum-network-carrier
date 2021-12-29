package carrier

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
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

			if err := s.refreshResourceNodes(); nil != err {
				log.WithError(err).Errorf("Failed to call service.refreshResourceNodes() on service.loop()")
			}

		case <-s.quit:
			log.Info("Stopped carrier service ...")
			return
		}
	}
}

func (s *Service) initServicesWithDiscoveryCenter() error {

	if nil == s.carrierDB {
		return fmt.Errorf("init services with discorvery center failed, carrier local db is nil")
	}

	if nil == s.consulManager {
		return fmt.Errorf("init services with discorvery center failed, consul manager is nil")
	}

	// 1. fetch datacenter ip&port config from discorvery center
	datacenterIpAndPort, err := s.consulManager.GetKV(discovery.DataCenterConsulServiceAddressKey, nil)
	if nil == err {
		return fmt.Errorf("query datacenter IP and PORT KVconfig failed from discovery center, %s", err)
	}
	if "" == datacenterIpAndPort {
		return fmt.Errorf("datacenter IP and PORT KVconfig is empty from discovery center")
	}
	configArr := strings.Split(datacenterIpAndPort, discovery.ConsulServiceIdSeparator)

	if len(configArr) != 2 {
		return fmt.Errorf("datacenter IP and PORT lack one on KVconfig from discovery center")
	}

	datacenterIP, datacenterPortStr := configArr[0], configArr[1]

	datacenterPort, err := strconv.Atoi(datacenterPortStr)
	if nil == err {
		return fmt.Errorf("invalid datacenter Port from discovery center, %s", err)
	}

	if err := s.carrierDB.SetConfig(&params.CarrierChainConfig{
		GrpcUrl: datacenterIP,
		Port:    uint64(datacenterPort),
	}); nil != err {
		return fmt.Errorf("connot setConfig of carrierDB, %s", err)
	}

	// 2. register carrier service to discorvery center
	if err := s.consulManager.RegisterService2DiscoveryCenter(); nil != err {
		return fmt.Errorf("register discovery service to discovery center failed, %s", err)
	}

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

	// 2. fetch via external ip&port config from discorvery center
	viaExternalIpAndPort, err := s.consulManager.GetKV(discovery.ViaNodeConsulServiceExternalAddressKey, nil)
	if nil == err {
		return fmt.Errorf("query via external IP and PORT KVconfig from discovery center failed, %s", err)
	}
	if "" == viaExternalIpAndPort {
		return fmt.Errorf("via external IP and PORT KVconfig is empty from discovery center")
	}
	configArr := strings.Split(viaExternalIpAndPort, discovery.ConsulServiceIdSeparator)

	if len(configArr) != 2 {
		return fmt.Errorf("via external IP and PORT lack one on KVconfig from discovery center")
	}

	viaExternalIP, viaExternalPort := configArr[0], configArr[1]

	storeLocalResourceFn := func(identity *apicommonpb.Organization, jobNodeId string, jobNodeStatus *computesvc.GetStatusReply) error {
		// store into local db
		if err := s.carrierDB.InsertLocalResource(types.NewLocalResource(&libtypes.LocalResourcePB{
			IdentityId: identity.GetIdentityId(),
			NodeId:     identity.GetNodeId(),
			NodeName:   identity.GetNodeName(),
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
			log.WithError(err).Errorf("Failed to store power to local on service.refreshResourceNodes(), jobNodeId: {%s}",
				jobNodeId)
			return err
		}
		return nil
	}

	// ##########################
	// ##########################
	// ABOUT JOBNODE SERVICES
	// ##########################
	// ##########################

	jobNodeServices, err := s.consulManager.QueryJobNodeServices()
	if nil != err {
		log.WithError(err).Warnf("query jobNodeServices failed from discovery center on service.refreshResourceNodes()")
	} else {

		jobNodeCache := make(map[string]*pb.YarnRegisteredPeerDetail, 0)
		// load stored jobNode
		jobNodeList, err := s.carrierDB.QueryRegisterNodeList(pb.PrefixTypeJobNode)
		if nil != err && rawdb.IsNoDBNotFoundErr(err) {
			log.WithError(err).Warnf("query jobNodes failed from local db on service.refreshResourceNodes()")
		} else {
			for _, node := range jobNodeList {
				jobNodeCache[node.GetId()] = node
			}

			for _, jobNodeService := range jobNodeServices {
				if node, ok := jobNodeCache[jobNodeService.ID]; !ok { // add new registered jobNode service
					client, err := grpclient.NewJobNodeClient(s.ctx, fmt.Sprintf("%s:%d", jobNodeService.Address, jobNodeService.Port), jobNodeService.ID)
					if err != nil {
						log.WithError(err).Errorf("Failed to connect new jobNode on service.refreshResourceNodes(), jobNodeServiceId: {%s}, jobNodeService: {%s:%d}",
							jobNodeService.ID, jobNodeService.Address, jobNodeService.Port)
						continue
					}
					jobNodeStatus, err := client.GetStatus()
					if err != nil {
						log.WithError(err).Errorf("Failed to connect jobNode to query status on service.refreshResourceNodes(), jobNodeServiceId: {%s}, jobNodeService: {%s:%d}",
							jobNodeService.ID, jobNodeService.Address, jobNodeService.Port)
						continue
					}
					// 1. add local jobNode resource
					// add resource usage first, but not own power now (mem, proccessor, bandwidth)
					if err = storeLocalResourceFn(identity, jobNodeService.ID, jobNodeStatus); nil != err {
						log.WithError(err).Errorf("Failed to store jobNode local resource on service.refreshResourceNodes(), jobNodeServiceId: {%s}, jobNodeService: {%s:%d}",
							jobNodeService.ID, jobNodeService.Address, jobNodeService.Port)
						continue
					}

					// 2. add rpc client
					s.resourceClientSet.StoreJobNodeClient(jobNodeService.ID, client)

					// 3. add local jobNode info
					// build new jobNode info that was need to store local db
					node = &pb.YarnRegisteredPeerDetail{
						InternalIp:   jobNodeService.Address,
						InternalPort: strconv.Itoa(jobNodeService.Port),
						ExternalIp:   viaExternalIP,
						ExternalPort: viaExternalPort,
						ConnState:    pb.ConnState_ConnState_Connected,
					}
					node.Id = strings.Join([]string{discovery.JobNodeConsulServiceIdPrefix, jobNodeService.Address,
						strconv.Itoa(jobNodeService.Port)}, discovery.ConsulServiceIdSeparator)
					// 4. store jobNode ip port into local db
					if err = s.carrierDB.SetRegisterNode(pb.PrefixTypeJobNode, node); err != nil {
						log.WithError(err).Errorf("Failed to store registerNode into local db on service.refreshResourceNodes(), jobNodeServiceId: {%s}, jobNodeService: {%s:%d}",
							jobNodeService.ID, jobNodeService.Address, jobNodeService.Port)
					}

				} else {
					// check the  via external ip and port comparing old infomation,
					// if it is, update the some things about jobNode.
					if node.GetExternalIp() != viaExternalIP || node.GetExternalPort() != viaExternalPort {
						// update jobNode info that was need to store local db
						node.ExternalIp = viaExternalIP
						node.ExternalPort = viaExternalPort
						// 1. update local jobNode info
						// update jobNode ip port into local db
						if err = s.carrierDB.SetRegisterNode(pb.PrefixTypeJobNode, node); err != nil {
							log.WithError(err).Errorf("Failed to update jobNode into local db on service.refreshResourceNodes(), jobNodeServiceId: {%s}, jobNodeService: {%s:%d}",
								jobNodeService.ID, jobNodeService.Address, jobNodeService.Port)
						}
					}
					delete(jobNodeCache, jobNodeService.ID)
				}
			}

			// delete old deregistered jobNode service
			for jobNodeId, _ := range jobNodeCache {

				// The published jobNode cannot be updated directly
				resourceTable, err := s.carrierDB.QueryLocalResourceTable(jobNodeId)
				if rawdb.IsNoDBNotFoundErr(err) {
					log.WithError(err).Errorf("Failed to query local power resource on old jobNode on service.refreshResourceNodes()")
					continue
				}
				if nil != resourceTable {
					log.Warnf("still have the published computing power information by the jobNode, that need revoke power short circuit on service.refreshResourceNodes(), %s",
						resourceTable.String())
					// ##############################
					// A. remove power about jobNode
					// ##############################

					// 1. remove jobNodeId and powerId mapping
					if err := s.carrierDB.RemoveJobNodeIdByPowerId(resourceTable.GetPowerId()); nil != err {
						log.WithError(err).Errorf("Failed to call RemoveJobNodeIdByPowerId() on service.refreshResourceNodes() with revoke power, powerId: {%s}, jobNodeId: {%s}",
							resourceTable.GetPowerId(), jobNodeId)
						continue
					}

					// 2. remove local jobNode resource table
					if err := s.carrierDB.RemoveLocalResourceTable(jobNodeId); nil != err {
						log.WithError(err).Errorf("Failed to call RemoveLocalResourceTable() on service.refreshResourceNodes() with revoke power, powerId: {%s}, jobNodeId: {%s}",
							resourceTable.GetPowerId(), jobNodeId)
						continue
					}

					// 3. revoke power about jobNode from global
					if err := s.carrierDB.RevokeResource(types.NewResource(&libtypes.ResourcePB{
						IdentityId: identity.GetIdentityId(),
						NodeId:     identity.GetNodeId(),
						NodeName:   identity.GetNodeName(),
						DataId:     resourceTable.GetPowerId(),
						// the status of data, N means normal, D means deleted.
						DataStatus: apicommonpb.DataStatus_DataStatus_Deleted,
						// resource status, eg: create/release/revoke
						State:    apicommonpb.PowerState_PowerState_Revoked,
						UpdateAt: timeutils.UnixMsecUint64(),
					})); nil != err {
						log.WithError(err).Errorf("Failed to remove dataCenter resource on service.refreshResourceNodes() with revoke power, powerId: {%s}, jobNodeId: {%s}",
							resourceTable.GetPowerId(), jobNodeId)
						continue
					}
				}

				// ##############################
				// B. remove except power things about jobNode
				// ##############################

				// 1. remove all running task
				taskIdList, _ := s.carrierDB.QueryJobNodeRunningTaskIdList(jobNodeId)
				for _, taskId := range taskIdList {
					s.carrierDB.RemoveJobNodeTaskIdAllPartyIds(jobNodeId, taskId)
				}
				// 2. remove local jobNode reource
				// remove jobNode local resource
				if err = s.carrierDB.RemoveLocalResource(jobNodeId); nil != err {
					log.WithError(err).Errorf("Failed to remove jobNode local resource on service.refreshResourceNodes(), jobNodeId: {%s}",
						jobNodeId)
					continue
				}
				// 3. remove rpc client
				if client, ok := s.resourceClientSet.QueryJobNodeClient(jobNodeId); ok {
					client.Close()
					s.resourceClientSet.RemoveJobNodeClient(jobNodeId)
				}
				// 4. remove local jobNode info
				if err = s.carrierDB.DeleteRegisterNode(pb.PrefixTypeJobNode, jobNodeId); err != nil {
					log.WithError(err).Errorf("Failed to remove jobNode into local db on service.refreshResourceNodes(), jobNodeId: {%s}",
						jobNodeId)
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

		dataNodeCache := make(map[string]*pb.YarnRegisteredPeerDetail, 0)
		// load stored dataNode
		dataNodeList, err := s.carrierDB.QueryRegisterNodeList(pb.PrefixTypeDataNode)
		if nil != err && rawdb.IsNoDBNotFoundErr(err) {
			log.WithError(err).Warnf("query dataNodes from local db failed on service.refreshResourceNodes()")
		} else if nil == err {
			for _, node := range dataNodeList {
				dataNodeCache[node.GetId()] = node
			}

			for _, dataNodeService := range dataNodeServices {
				if node, ok := dataNodeCache[dataNodeService.ID]; !ok { // add new registered dataNode service
					client, err := grpclient.NewDataNodeClient(s.ctx, fmt.Sprintf("%s:%d", dataNodeService.Address, dataNodeService.Port), dataNodeService.ID)
					if err != nil {
						log.WithError(err).Errorf("Failed to connect new dataNode on service.refreshResourceNodes(), dataNodeServiceId: {%s}, dataNodeService: {%s:%d}",
							dataNodeService.ID, dataNodeService.Address, dataNodeService.Port)
						continue
					}
					dataNodeStatus, err := client.GetStatus()
					if err != nil {
						log.WithError(err).Errorf("Failed to connect jobNode to query status on service.refreshResourceNodes(), dataNodeServiceId: {%s}, dataNodeService: {%s:%d}",
							dataNodeService.ID, dataNodeService.Address, dataNodeService.Port)
						continue
					}
					// 1. add data resource  (disk)
					err = s.carrierDB.StoreDataResourceTable(types.NewDataResourceTable(node.GetId(), dataNodeStatus.GetTotalDisk(), dataNodeStatus.GetUsedDisk()))
					if err != nil {
						log.WithError(err).Errorf("Failed to store disk summary of new dataNode")
						continue
					}

					// 2. add rpc client
					s.resourceClientSet.StoreDataNodeClient(dataNodeService.ID, client)

					// 3. add local dataNode info
					// build new dataNode info that was need to store local db
					node = &pb.YarnRegisteredPeerDetail{
						InternalIp:   dataNodeService.Address,
						InternalPort: strconv.Itoa(dataNodeService.Port),
						ExternalIp:   viaExternalIP,
						ExternalPort: viaExternalPort,
						ConnState:    pb.ConnState_ConnState_Connected,
					}
					node.Id = strings.Join([]string{discovery.DataNodeConsulServiceIdPrefix, dataNodeService.Address,
						strconv.Itoa(dataNodeService.Port)}, discovery.ConsulServiceIdSeparator)
					// 4. store dataNode ip port into local db
					if err = s.carrierDB.SetRegisterNode(pb.PrefixTypeDataNode, node); err != nil {
						log.WithError(err).Errorf("Failed to store dataNode into local db on service.refreshResourceNodes(), dataNodeServiceId: {%s}, dataNodeService: {%s:%d}",
							dataNodeService.ID, dataNodeService.Address, dataNodeService.Port)
					}

				} else {
					// check the  via external ip and port comparing old infomation,
					// if it is, update the some things about dataNode.
					if node.GetExternalIp() != viaExternalIP || node.GetExternalPort() != viaExternalPort {
						// update dataNode info that was need to store local db
						node.ExternalIp = viaExternalIP
						node.ExternalPort = viaExternalPort
						// 1. update local dataNode info
						// update dataNode ip port into local db
						if err = s.carrierDB.SetRegisterNode(pb.PrefixTypeDataNode, node); err != nil {
							log.WithError(err).Errorf("Failed to update dataNode into local db on service.refreshResourceNodes(), dataNodeServiceId: {%s}, dataNodeService: {%s:%d}",
								dataNodeService.ID, dataNodeService.Address, dataNodeService.Port)
						}
					}
					delete(dataNodeCache, dataNodeService.ID)
				}
			}

			// delete old deregistered dataNode service
			for dataNodeId, _ := range dataNodeCache {

				// 1. remove data resource  (disk)
				if err := s.carrierDB.RemoveDataResourceTable(dataNodeId); rawdb.IsNoDBNotFoundErr(err) {
					log.WithError(err).Errorf("Failed to remove disk summary of old dataNode on service.refreshResourceNodes()")
					continue
				}
				// 2. remove rpc client
				if client, ok := s.resourceClientSet.QueryDataNodeClient(dataNodeId); ok {
					client.Close()
					s.resourceClientSet.RemoveDataNodeClient(dataNodeId)
				}
				// 3. remove local dataNode info
				if err = s.carrierDB.DeleteRegisterNode(pb.PrefixTypeDataNode, dataNodeId); err != nil {
					log.WithError(err).Errorf("Failed to remove dataNode into local db on service.refreshResourceNodes(), dataNodeId: {%s}",
						dataNodeId)
				}
			}
		}
	}
	return nil
}
