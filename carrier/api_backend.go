package carrier

import (
	"errors"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	"github.com/RosettaFlow/Carrier-Go/grpclient"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
	"strings"
)

// CarrierAPIBackend implements rpc.Backend for Carrier
type CarrierAPIBackend struct {
	carrier *Service
}

func NewCarrierAPIBackend(carrier *Service) *CarrierAPIBackend {
	return &CarrierAPIBackend{carrier: carrier}
}

func (s *CarrierAPIBackend) SendMsg(msg types.Msg) error {
	return s.carrier.mempool.Add(msg)
}

// system (the yarn node self info)
func (s *CarrierAPIBackend) GetNodeInfo() (*pb.YarnNodeInfo, error) {
	jobNodes, err := s.carrier.carrierDB.GetRegisterNodeList(pb.PrefixTypeJobNode)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.Errorf("Failed to get all `job nodes`, on GetNodeInfo(), err: {%s}", err)
		return nil, err
	}
	dataNodes, err := s.carrier.carrierDB.GetRegisterNodeList(pb.PrefixTypeDataNode)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.Errorf("Failed to get all `data nodes, on GetNodeInfo(), err: {%s}", err)
		return nil, err
	}
	jobsLen := len(jobNodes)
	datasLen := len(dataNodes)
	length := jobsLen + datasLen
	registerNodes := make([]*pb.YarnRegisteredPeer, length)
	if len(jobNodes) != 0 {
		for i, v := range jobNodes {
			n := &pb.YarnRegisteredPeer{
				NodeType:   pb.NodeType_NodeType_JobNode,
				NodeDetail: v,
			}
			registerNodes[i] = n
		}
	}
	if len(dataNodes) != 0 {
		for i, v := range dataNodes {
			n := &pb.YarnRegisteredPeer{
				NodeType:   pb.NodeType_NodeType_DataNode,
				NodeDetail: v,
			}
			/*n.RegisteredNodeInfo = &pb.YarnRegisteredPeerDetail{
				Id:           v.Id,
				InternalIp:   v.InternalIp,
				InternalPort: v.InternalPort,
				ExternalIp:   v.ExternalIp,
				ExternalPort: v.ExternalPort,
				ConnState:    v.ConnState,
			}*/
			registerNodes[jobsLen+i] = n
		}
	}

	identity, err := s.carrier.carrierDB.GetIdentity()
	if nil != err {
		log.Warnf("Failed to get identity, on GetNodeInfo(), err: {%s}", err)
		//return nil, fmt.Errorf("query local identity failed, %s", err)
	}
	var identityId string
	var nodeId string
	var nodeName string
	if nil != identity {
		identityId = identity.IdentityId
		nodeId = identity.NodeId
		nodeName = identity.NodeName
	}

	seedNodes, err := s.carrier.carrierDB.GetSeedNodeList()
	return &pb.YarnNodeInfo{
		NodeType:     pb.NodeType_NodeType_YarnNode,
		NodeId:       nodeId,
		InternalIp:   "",                               //
		ExternalIp:   "",                               //
		InternalPort: "",                               //
		ExternalPort: "",                               //
		//TODO: 需要更改
		//IdentityType: types.IDENTITY_TYPE_DID.String(), // 默认先是 DID
		IdentityId:   identityId,
		Name:         nodeName,
		Peers:        registerNodes,
		SeedPeers:    seedNodes,
		//TODO: 需要更改
		//State:        types.YARN_STATE_ACTIVE.String(),
	}, nil
}

func (s *CarrierAPIBackend) GetRegisteredPeers() ([]*pb.YarnRegisteredPeer, error) {
	// all dataNodes on yarnNode
	dataNodes, err := s.carrier.carrierDB.GetRegisterNodeList(pb.PrefixTypeDataNode)
	if nil != err {
		return nil, err
	}

	// all jobNodes on yarnNode
	jobNodes, err := s.carrier.carrierDB.GetRegisterNodeList(pb.PrefixTypeJobNode)
	if nil != err {
		return nil, err
	}

	result := make([]*pb.YarnRegisteredPeer, 0)

	// 处理计算节点
	for _, v := range jobNodes {
		var duration uint64
		node, has := s.carrier.resourceClientSet.QueryJobNodeClient(v.Id)
		if has {
			duration = uint64(node.RunningDuration())
		}
		v.TaskCount, _ = s.carrier.carrierDB.GetRunningTaskCountOnJobNode(v.Id)
		v.TaskIdList, _ = s.carrier.carrierDB.GetJobNodeRunningTaskIdList(v.Id)
		v.Duration = duration
		registeredPeer := &pb.YarnRegisteredPeer{
			NodeType:   pb.NodeType_NodeType_JobNode,
			NodeDetail: v,
		}
		result = append(result, registeredPeer)
	}

	// 处理数据节点
	for _, v := range dataNodes {
		var duration uint64
		node, has := s.carrier.resourceClientSet.QueryDataNodeClient(v.Id)
		if has {
			duration = uint64(node.RunningDuration())
		}
		v.Duration = duration // ms
		v.FileCount = 0
		v.FileTotalSize = 0
		registeredPeer := &pb.YarnRegisteredPeer{
			NodeType:   pb.NodeType_NodeType_DataNode,
			NodeDetail: v,
		}
		result = append(result, registeredPeer)
		//TODO: 需要更改
		//n := &types.YarnRegisteredDataNode{
		//	Id:           v.Id,
		//	InternalIp:   v.InternalIp,
		//	ExternalIp:   v.ExternalIp,
		//	InternalPort: v.InternalPort,
		//	ExternalPort: v.ExternalPort,
		//	//ResourceUsage:  &types.ResourceUsage{},
		//	Duration: duration, // ms
		//}
		//n.Delta.FileCount = 0
		//n.Delta.FileTotalSize = 0
		//dns[i] = n
	}
	return result, nil
}

func (s *CarrierAPIBackend) SetSeedNode(seed *pb.SeedPeer) (pb.ConnState, error) {
	//TODO: current node need to connect with seed node.(delay processing)
	return s.carrier.carrierDB.SetSeedNode(seed)
}

func (s *CarrierAPIBackend) DeleteSeedNode(id string) error {
	return s.carrier.carrierDB.DeleteSeedNode(id)
}

func (s *CarrierAPIBackend) GetSeedNode(id string) (*pb.SeedPeer, error) {
	return s.carrier.carrierDB.GetSeedNode(id)
}

func (s *CarrierAPIBackend) GetSeedNodeList() ([]*pb.SeedPeer, error) {
	return s.carrier.carrierDB.GetSeedNodeList()
}

func (s *CarrierAPIBackend) SetRegisterNode(typ pb.RegisteredNodeType, node *pb.YarnRegisteredPeerDetail) (pb.ConnState, error) {
	switch typ {
	case pb.PrefixTypeDataNode, pb.PrefixTypeJobNode:
	default:
		return pb.ConnState_ConnState_UnConnected, errors.New("invalid nodeType")
	}
	if typ == pb.PrefixTypeJobNode {
		client, err := grpclient.NewJobNodeClientWithConn(s.carrier.ctx, fmt.Sprintf("%s:%s", node.InternalIp, node.InternalPort), node.Id)
		if err != nil {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("connect new jobNode failed, %s", err)
		}
		s.carrier.resourceClientSet.StoreJobNodeClient(node.Id, client)
	}

	if typ == pb.PrefixTypeDataNode {
		client, err := grpclient.NewDataNodeClientWithConn(s.carrier.ctx, fmt.Sprintf("%s:%s", node.InternalIp, node.InternalPort), node.Id)
		if err != nil {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("connect new dataNode failed, %s", err)
		}
		s.carrier.resourceClientSet.StoreDataNodeClient(node.Id, client)

		// add data resource  (disk)  todo 后续 需要根据 真实的 dataNode 上报自身的 disk 信息
		err = s.carrier.carrierDB.StoreDataResourceTable(types.NewDataResourceTable(node.Id, types.DefaultDisk, 0))
		if err != nil {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("store disk summary of new dataNode failed, %s", err)
		}
	}
	node.ConnState = pb.ConnState_ConnState_Connected
	_, err := s.carrier.carrierDB.SetRegisterNode(typ, node)
	if err != nil {
		return pb.ConnState_ConnState_UnConnected, fmt.Errorf("Store registerNode to db failed, %s", err)
	}
	return pb.ConnState_ConnState_Connected, nil
}

func (s *CarrierAPIBackend) UpdateRegisterNode(typ pb.RegisteredNodeType, node *pb.YarnRegisteredPeerDetail) (pb.ConnState, error) {

	switch typ {
	case pb.PrefixTypeDataNode, pb.PrefixTypeJobNode:
	default:
		return pb.ConnState_ConnState_UnConnected, errors.New("invalid nodeType")
	}
	if typ == pb.PrefixTypeJobNode {

		// 算力已发布的jobNode不可以直接删除
		resourceTable, err := s.carrier.carrierDB.QueryLocalResourceTable(node.Id)
		if rawdb.IsNoDBNotFoundErr(err) {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("query local power resource on old jobNode failed, %s", err)
		}

		if nil != resourceTable {
			log.Debugf("still have the published computing power information on old jobNode on UpdateRegisterNode, %s", resourceTable.String())
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("still have the published computing power information on old jobNode failed, input jobNodeId: {%s}, old jobNodeId: {%s}, old powerId: {%s}",
				node.Id, resourceTable.GetNodeId(), resourceTable.GetPowerId())
		}

		// 先校验 jobNode 上是否有正在执行的 task
		runningTaskCount, err := s.carrier.carrierDB.GetRunningTaskCountOnJobNode(node.Id)
		if rawdb.IsNoDBNotFoundErr(err) {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("query local running taskCount on old jobNode failed, %s", err)
		}
		if runningTaskCount > 0 {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("the old jobNode have been running {%d} task current, don't remove it", runningTaskCount)
		}

		if client, ok := s.carrier.resourceClientSet.QueryJobNodeClient(node.Id); ok {
			// remove old client instanse
			client.Close()
			s.carrier.resourceClientSet.RemoveJobNodeClient(node.Id)
		}

		// generate new client
		client, err := grpclient.NewJobNodeClientWithConn(s.carrier.ctx, fmt.Sprintf("%s:%s", node.InternalIp, node.InternalPort), node.Id)
		if err != nil {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("connect new jobNode failed, %s", err)
		}
		s.carrier.resourceClientSet.StoreJobNodeClient(node.Id, client)

	}

	if typ == pb.PrefixTypeDataNode {

		// 先校验 dataNode 上是否已被 使用
		dataNodeTable, err := s.carrier.carrierDB.QueryDataResourceTable(node.Id)
		if rawdb.IsNoDBNotFoundErr(err) {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("query disk used summary on old dataNode failed, %s", err)
		}
		if dataNodeTable.IsNotEmpty() && dataNodeTable.IsUsed() {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("the disk of old dataNode was used, don't remove it, totalDisk: {%d byte}, usedDisk: {%d byte}, remainDisk: {%d byte}",
				dataNodeTable.GetTotalDisk(), dataNodeTable.GetUsedDisk(), dataNodeTable.RemainDisk())
		}

		if client, ok := s.carrier.resourceClientSet.QueryDataNodeClient(node.Id); ok {
			// remove old client instanse
			client.Close()
			s.carrier.resourceClientSet.RemoveDataNodeClient(node.Id)
		}

		// remove old data resource  (disk)  todo 后续 需要根据 真实的 dataNode 上报自身的 disk 信息
		if err := s.carrier.carrierDB.RemoveDataResourceTable(node.Id); rawdb.IsNoDBNotFoundErr(err) {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("remove disk summary of old dataNode, %s", err)
		}

		client, err := grpclient.NewDataNodeClientWithConn(s.carrier.ctx, fmt.Sprintf("%s:%s", node.InternalIp, node.InternalPort), node.Id)
		if err != nil {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("connect new dataNode failed, %s", err)
		}
		s.carrier.resourceClientSet.StoreDataNodeClient(node.Id, client)

		// add new data resource  (disk)  todo 后续 需要根据 真实的 dataNode 上报自身的 disk 信息
		err = s.carrier.carrierDB.StoreDataResourceTable(types.NewDataResourceTable(node.Id, types.DefaultDisk, 0))
		if err != nil {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("store disk summary of new dataNode failed, %s", err)
		}

	}

	// remove  old jobNode from db
	if err := s.carrier.carrierDB.DeleteRegisterNode(typ, node.Id); nil != err {
		return pb.ConnState_ConnState_UnConnected, fmt.Errorf("remove old registerNode from db failed, %s", err)
	}

	// add new node to db
	node.ConnState = pb.ConnState_ConnState_Connected
	_, err := s.carrier.carrierDB.SetRegisterNode(typ, node)
	if err != nil {
		return pb.ConnState_ConnState_UnConnected, fmt.Errorf("Store new registerNode to db failed, %s", err)
	}
	return pb.ConnState_ConnState_Connected, nil
}

func (s *CarrierAPIBackend) DeleteRegisterNode(typ pb.RegisteredNodeType, id string) error {
	switch typ {
	case pb.PrefixTypeDataNode, pb.PrefixTypeJobNode:
	default:
		return errors.New("invalid nodeType")
	}
	if typ == pb.PrefixTypeJobNode {

		// 算力已发布的jobNode不可以直接删除
		resourceTable, err := s.carrier.carrierDB.QueryLocalResourceTable(id)
		if rawdb.IsNoDBNotFoundErr(err) {
			return fmt.Errorf("query local power resource on old jobNode failed, %s", err)
		}

		if nil != resourceTable {
			log.Debugf("still have the published computing power information on old jobNode on DeleteRegisterNode, %s", resourceTable.String())
			return fmt.Errorf("still have the published computing power information on old jobNode failed,input jobNodeId: {%s}, old jobNodeId: {%s}, old powerId: {%s}",
				id, resourceTable.GetNodeId(), resourceTable.GetPowerId())
		}

		// 先校验 jobNode 上是否有正在执行的 task
		runningTaskCount, err := s.carrier.carrierDB.GetRunningTaskCountOnJobNode(id)
		if rawdb.IsNoDBNotFoundErr(err) {
			return fmt.Errorf("query local running taskCount on old jobNode failed, %s", err)
		}
		if runningTaskCount > 0 {
			return fmt.Errorf("the old jobNode have been running {%d} task current, don't remove it", runningTaskCount)
		}

		if client, ok := s.carrier.resourceClientSet.QueryJobNodeClient(id); ok {
			client.Close()
			s.carrier.resourceClientSet.RemoveJobNodeClient(id)
		}
	}

	if typ == pb.PrefixTypeDataNode {

		// 先校验 dataNode 上是否已被 使用
		dataNodeTable, err := s.carrier.carrierDB.QueryDataResourceTable(id)
		if rawdb.IsNoDBNotFoundErr(err) {
			return fmt.Errorf("query disk used summary on old dataNode failed, %s", err)
		}
		if dataNodeTable.IsNotEmpty() && dataNodeTable.IsUsed() {
			return fmt.Errorf("the disk of old dataNode was used, don't remove it, totalDisk: {%d byte}, usedDisk: {%d byte}, remainDisk: {%d byte}",
				dataNodeTable.GetTotalDisk(), dataNodeTable.GetUsedDisk(), dataNodeTable.RemainDisk())
		}

		if client, ok := s.carrier.resourceClientSet.QueryDataNodeClient(id); ok {
			client.Close()
			s.carrier.resourceClientSet.RemoveDataNodeClient(id)
		}
		// remove data resource  (disk)  todo 后续 需要根据 真实的 dataNode 上报自身的 disk 信息
		if err := s.carrier.carrierDB.RemoveDataResourceTable(id); rawdb.IsNoDBNotFoundErr(err) {
			return fmt.Errorf("remove disk summary of old registerNode, %s", err)
		}
	}
	return s.carrier.carrierDB.DeleteRegisterNode(typ, id)
}

func (s *CarrierAPIBackend) GetRegisterNode(typ pb.RegisteredNodeType, id string) (*pb.YarnRegisteredPeerDetail, error) {
	return s.carrier.carrierDB.GetRegisterNode(typ, id)
}

func (s *CarrierAPIBackend) GetRegisterNodeList(typ pb.RegisteredNodeType) ([]*pb.YarnRegisteredPeerDetail, error) {
	nodeList, err := s.carrier.carrierDB.GetRegisterNodeList(typ)
	if nil != err {
		return nil, err
	}

	// 需要处理 计算服务 信息
	if typ == pb.PrefixTypeJobNode {
		for i, jobNode := range nodeList {

			client, ok := s.carrier.resourceClientSet.QueryJobNodeClient(jobNode.Id)
			if !ok {
				jobNode.ConnState = pb.ConnState_ConnState_UnConnected
			}
			if !client.IsConnected() {
				jobNode.ConnState = pb.ConnState_ConnState_UnConnected
			} else {
				jobNode.ConnState = pb.ConnState_ConnState_Connected
			}

			table, err := s.carrier.carrierDB.QueryLocalResourceTable(jobNode.Id)
			if nil != err {
				continue
			}
			if nil != table {
				jobNode.ConnState = pb.ConnState_ConnState_Enabled
			}
			taskCount, err := s.carrier.carrierDB.GetRunningTaskCountOnJobNode(jobNode.Id)
			if nil != err {
				continue
			}
			if taskCount > 0 {
				jobNode.ConnState = pb.ConnState_ConnState_Occupied
			}
			nodeList[i] = jobNode
		}
	}

	return nodeList, nil
}

func (s *CarrierAPIBackend) SendTaskEvent(event *libTypes.TaskEvent) error {
	// return s.carrier.resourceManager.SendTaskEvent(evengine)
	return s.carrier.taskManager.SendTaskEvent(event)
}

// metadata api
func (s *CarrierAPIBackend) GetMetadataDetail(identityId, metaDataId string) (*pb.GetMetadataDetailResponse, error) {
	metadata, err := s.carrier.carrierDB.GetMetadataByDataId(metaDataId)
	if metadata == nil {
		return nil, errors.New("not found metadata by special Id")
	}
	return types.NewOrgMetadataInfoFromMetadata(metadata), err
}

// GetMetadataDetailList returns a list of all metadata details in the network.
func (s *CarrierAPIBackend) GetMetadataDetailList() ([]*pb.GetMetadataDetailResponse, error) {
	metadataArray, err := s.carrier.carrierDB.GetMetadataList()
	return types.NewOrgMetadataInfoArrayFromMetadataArray(metadataArray), err
}

func (s *CarrierAPIBackend) GetMetadataDetailListByOwner(identityId string) ([]*pb.GetMetadataDetailResponse, error) {
	log.WithField("identityId", identityId).Debug("Invoke: GetMetadataDetailListByOwner executing...")
	metadataList, err := s.GetMetadataDetailList()
	if err != nil {
		return nil, err
	}
	result := make([]*pb.GetMetadataDetailResponse, 0)
	for _, metadata := range metadataList {
		if metadata.Owner.IdentityId == identityId {
			result = append(result, metadata)
		}
	}
	log.Debugf("Query metaData list, identityId: {%s}, len: {%d}", identityId, len(result))
	return result, nil
}

// power api
func (s *CarrierAPIBackend) GetPowerTotalDetailList() ([]*pb.GetPowerTotalDetailResponse, error) {
	log.Debug("Invoke: GetPowerTotalDetailList executing...")
	resourceList, err := s.carrier.carrierDB.GetResourceList()
	if err != nil {
		return nil, err
	}
	log.Debugf("Query all org's power list, len: {%d}", len(resourceList))
	powerList := make([]*pb.GetPowerTotalDetailResponse, 0, resourceList.Len())
	for _, resource := range resourceList.To() {
		powerList = append(powerList, &pb.GetPowerTotalDetailResponse{
			Owner: &apipb.Organization{
				NodeName:   resource.GetNodeName(),
				NodeId:     resource.GetNodeId(),
				IdentityId: resource.GetIdentityId(),
			},
			Power: &libTypes.PowerUsageDetail{
				TotalTaskCount:   0,
				CurrentTaskCount: 0,
				Tasks:            make([]*libTypes.PowerTask, 0),
				Information: &libTypes.ResourceUsageOverview{
					TotalMem:       resource.GetTotalMem(),
					UsedMem:        resource.GetUsedMem(),
					TotalProcessor: uint32(resource.GetTotalProcessor()),
					UsedProcessor:  uint32(resource.GetUsedProcessor()),
					TotalBandwidth: resource.GetTotalBandwidth(),
					UsedBandwidth:  resource.GetUsedBandwidth(),
				},
				State: resource.GetState(),
			},
		})
	}
	return powerList, nil
}

func (s *CarrierAPIBackend) GetPowerSingleDetailList() ([]*pb.GetPowerSingleDetailResponse, error) {
	log.Debug("Invoke:GetPowerSingleDetailList executing...")
	// query local resource list from db.
	machineList, err := s.carrier.carrierDB.GetLocalResourceList()
	if err != nil {
		return nil, err
	}
	log.Debugf("Invoke:GetPowerSingleDetailList, call GetLocalResourceList, machineList: %s", machineList.String())

	// query used of power for local task. : taskId -> {taskId, jobNodeId, slotCount}
	localTaskPowerUsedList, err := s.carrier.carrierDB.QueryLocalTaskPowerUseds()
	if err != nil {
		return nil, err
	}
	log.Debugf("Invoke:GetPowerSingleDetailList, call QueryLocalTaskPowerUseds, localTaskPowerUsedList: %s",
		utilLocalTaskPowerUsedArrString(localTaskPowerUsedList))

	slotUnit, err := s.carrier.carrierDB.QueryNodeResourceSlotUnit()
	if err != nil {
		return nil, err
	}
	log.Debugf("Invoke:GetPowerSingleDetailList, call QueryNodeResourceSlotUnit, slotUint: %s",
		slotUnit.String())

	// 收集 本地所有的 jonNode 上的 powerUsed 数组
	validLocalTaskPowerUsedMap := make(map[string][]*types.LocalTaskPowerUsed, 0)
	for _, taskPowerUsed := range localTaskPowerUsedList {
		// condition: jobNode
		usedArr, ok := validLocalTaskPowerUsedMap[taskPowerUsed.GetNodeId()]
		if ok {
			usedArr = append(usedArr, taskPowerUsed)
		} else {
			usedArr = []*types.LocalTaskPowerUsed{taskPowerUsed}
		}
		validLocalTaskPowerUsedMap[taskPowerUsed.GetNodeId()] = usedArr
	}

	log.Debugf("Invoke:GetPowerSingleDetailList, make validLocalTaskPowerUsedMap, validLocalTaskPowerUsedMap: %s",
		utilLocalTaskPowerUsedMapString(validLocalTaskPowerUsedMap))

	//// 收集 本地所有 计算资源的 powerUsed 数组
	//validLocalTaskPowerUsedMap := make(map[string][]*types.LocalTaskPowerUsed, 0)
	//for _, jobNode := range machineList {
	//	// condition: jobNode
	//	if usedArr, ok := localTaskPowerUsedTmp[jobNode.GetJobNodeId()]; ok {
	//		validLocalTaskPowerUsedMap[jobNode.GetJobNodeId()] = usedArr
	//	}
	//}

	readElement := func(jobNodeId string, taskId string) uint64 {
		if usedArr, ok := validLocalTaskPowerUsedMap[jobNodeId]; ok {
			for _, powerUsed := range usedArr {
				if jobNodeId == powerUsed.GetNodeId() && powerUsed.GetTaskId() == taskId {
					return powerUsed.GetSlotCount()
				}
			}
		}
		return 0
	}

	buildPowerTaskList := func(jobNodeId string) []*libTypes.PowerTask {
		powerTaskList := make([]*libTypes.PowerTask, 0)

		// 逐个 处理 jobNodeId 上的 task 信息
		if usedArr, ok := validLocalTaskPowerUsedMap[jobNodeId]; ok {
			for _, powerUsed := range usedArr {
				taskId := powerUsed.GetTaskId()
				task, err := s.carrier.carrierDB.GetLocalTask(taskId)
				if err != nil {
					log.Errorf("Failed to query local task on GetPowerSingleDetailList, taskId: {%s}, err: {%s}", taskId, err)
					continue
				}

				// 封装任务 摘要 ...
				powerTask := &libTypes.PowerTask{
					TaskId:   taskId,
					TaskName: task.TaskData().TaskName,
					Owner: &apipb.Organization{
						NodeName:   task.TaskData().GetNodeName(),
						NodeId:     task.TaskData().GetNodeId(),
						IdentityId: task.TaskData().GetIdentityId(),
					},
					Receivers: make([]*apipb.Organization, 0),
					OperationCost: &apipb.TaskResourceCostDeclare{
						CostProcessor: task.TaskData().GetOperationCost().GetCostProcessor(),
						CostMem:       task.TaskData().GetOperationCost().GetCostMem(),
						CostBandwidth: task.TaskData().GetOperationCost().GetCostBandwidth(),
						Duration:      task.TaskData().GetOperationCost().GetDuration(),
					},
					OperationSpend: nil, // 下面单独计算 任务资源使用 实况 ...
					CreateAt:       task.TaskData().CreateAt,
				}
				// 组装 数据参与方
				//for _, dataSupplier := range task.TaskPB().DataSupplier {
				//	// 协作方, 需要过滤掉自己
				//	if task.TaskPB().GetNodeId() != dataSupplier.MemberInfo.IdentityId {
				//		powerTask.Patners = append(powerTask.Patners, &apipb.Organization{
				//			NodeName:   dataSupplier.MemberInfo.GetNodeName(),
				//			NodeId:     dataSupplier.MemberInfo.GetNodeId(),
				//			IdentityId: dataSupplier.MemberInfo.GetIdentityId(),
				//		})
				//	}
				//
				//}
				// 组装结果接收方
				for _, receiver := range task.TaskData().Receivers {
					powerTask.Receivers = append(powerTask.Receivers, &apipb.Organization{
						NodeName:   receiver.GetNodeName(),
						NodeId:     receiver.GetNodeId(),
						IdentityId: receiver.GetIdentityId(),
					})
				}

				// 计算任务使用实况 ...
				slotCount := readElement(jobNodeId, powerTask.TaskId)
				powerTask.OperationSpend = &apipb.TaskResourceCostDeclare{
					CostProcessor: uint32(slotUnit.Processor * slotCount),
					CostMem:       slotUnit.Mem * slotCount,
					CostBandwidth: slotUnit.Bandwidth * slotCount,
					Duration:      task.TaskData().GetOperationCost().GetDuration(),
				}
				powerTaskList = append(powerTaskList, powerTask)
			}
		}
		return powerTaskList
	}

	// 计算 jobNodeId 上的 task 数量
	taskCount := func(jobNodeId string) int {
		return len(validLocalTaskPowerUsedMap[jobNodeId])
	}

	resourceList := machineList.To()
	// 逐个处理当前
	result := make([]*pb.GetPowerSingleDetailResponse, len(resourceList))
	for i, resource := range resourceList {
		nodePowerDetail := &pb.GetPowerSingleDetailResponse{
			Owner: &apipb.Organization{
				NodeName:   resource.GetNodeName(),
				NodeId:     resource.GetNodeId(),
				IdentityId: resource.GetIdentityId(),
			},
			Power: &libTypes.PowerUsageDetail{
				//JobNodeId:        resource.GetJobNodeId(),
				//PowerId:          resource.DataId,
				TotalTaskCount:   uint32(taskCount(resource.GetJobNodeId())),
				CurrentTaskCount: uint32(taskCount(resource.GetJobNodeId())),
				Tasks:            make([]*libTypes.PowerTask, 0),
				Information: &libTypes.ResourceUsageOverview{
					TotalMem:       resource.GetTotalMem(),
					UsedMem:        resource.GetUsedMem(),
					TotalProcessor: uint32(resource.GetTotalProcessor()),
					UsedProcessor:  uint32(resource.GetUsedProcessor()),
					TotalBandwidth: resource.GetTotalBandwidth(),
					UsedBandwidth:  resource.GetUsedBandwidth(),
				},
				State: resource.GetState(),
			},
		}
		powerTaskArray := buildPowerTaskList(resource.GetJobNodeId())
		nodePowerDetail.Power.Tasks = powerTaskArray
		result[i] = nodePowerDetail
	}
	return result, nil
}

// identity api
func (s *CarrierAPIBackend) ApplyIdentityJoin(identity *types.Identity) error {
	//TODO: 申请身份标识时，相关数据需要进行本地存储，然后进行网络发布
	return s.carrier.carrierDB.InsertIdentity(identity)
}

func (s *CarrierAPIBackend) RevokeIdentityJoin(identity *types.Identity) error {
	return s.carrier.carrierDB.RevokeIdentity(identity)
}

func (s *CarrierAPIBackend) GetNodeIdentity() (*types.Identity, error) {
	nodeAlias, err := s.carrier.carrierDB.GetIdentity()
	if nil != err {
		return nil, err
	}
	return types.NewIdentity(&libTypes.IdentityPB{
		IdentityId: nodeAlias.IdentityId,
		NodeId:     nodeAlias.NodeId,
		NodeName:   nodeAlias.NodeName,
	}), err
}

func (s *CarrierAPIBackend) GetIdentityList() ([]*types.Identity, error) {
	return s.carrier.carrierDB.GetIdentityList()
}

func (s *CarrierAPIBackend) GetMetadataAuthorityList(identityId string, lastUpdate uint64) (types.MetadataAuthArray, error) {
	return s.carrier.carrierDB.GetMetadataAuthorityList(identityId, lastUpdate)
}

// task api
func (s *CarrierAPIBackend) GetTaskDetailList() ([]*pb.TaskDetailShow, error) {
	// the task is executing.
	localTaskArray, err := s.carrier.carrierDB.GetLocalTaskList()

	if rawdb.IsNoDBNotFoundErr(err) {
		return nil, err
	}
	localIdentityId, err := s.carrier.carrierDB.GetIdentityId()
	if err != nil {
		return nil, fmt.Errorf("query local identityId failed, %s", err)
	}

	// the task has been executed.
	networkTaskList, err := s.carrier.carrierDB.GetTaskListByIdentityId(localIdentityId)
	if rawdb.IsNoDBNotFoundErr(err) {
		return nil, err
	}

	makeTaskViewFn := func(task *types.Task) *pb.TaskDetailShow {
		// task 发起方
		if task.TaskData().GetIdentityId() == localIdentityId {
			return types.NewTaskDetailShowFromTaskData(task, apipb.TaskRole_TaskRole_Sender)
		}

		// task 参与方
		for _, dataSupplier := range task.TaskData().DataSuppliers {
			if dataSupplier.Organization.IdentityId == localIdentityId {
				return types.NewTaskDetailShowFromTaskData(task, apipb.TaskRole_TaskRole_DataSupplier)
			}
		}

		// 算力提供方
		for _, powerSupplier := range task.TaskData().PowerSuppliers {
			if powerSupplier.Organization.IdentityId == localIdentityId {
				return types.NewTaskDetailShowFromTaskData(task, apipb.TaskRole_TaskRole_PowerSupplier)
			}
		}

		// 数据接收方
		for _, receiver := range task.TaskData().Receivers {
			if receiver.IdentityId == localIdentityId {
				return types.NewTaskDetailShowFromTaskData(task, apipb.TaskRole_TaskRole_Receiver)
			}
		}
		return nil
	}

	result := make([]*pb.TaskDetailShow, 0)
	for _, task := range localTaskArray {
		if taskView := makeTaskViewFn(task); nil != taskView {
			result = append(result, taskView)
		}
	}

	for _, networkTask := range networkTaskList {
		if taskView := makeTaskViewFn(networkTask); nil != taskView {
			result = append(result, taskView)
		}
	}

	return result, err
}

func (s *CarrierAPIBackend) GetTaskEventList(taskId string) ([]*pb.TaskEventShow, error) {

	identity, err := s.carrier.carrierDB.GetIdentity()
	if nil != err {
		return nil, err
	}

	// 先查出 task 在本地的 eventList
	localEventList, err := s.carrier.carrierDB.GetTaskEventList(taskId)
	if rawdb.IsNoDBNotFoundErr(err) {
		return nil, err
	}

	evenList := make([]*pb.TaskEventShow, len(localEventList))
	for i, e := range localEventList {
		evenList[i] = &pb.TaskEventShow{
			TaskId:   e.TaskId,
			Type:     e.Type,
			CreateAt: e.CreateAt,
			Content:  e.Content,
			Owner: &apipb.Organization{
				NodeName:   identity.NodeName,
				NodeId:     identity.NodeId,
				IdentityId: identity.IdentityId,
			},
		}
	}

	taskEvent, err := s.carrier.carrierDB.GetTaskEventListByTaskId(taskId)
	if nil != err {
		return nil, err
	}
	evenList = append(evenList, types.NewTaskEventFromAPIEvent(taskEvent)...)
	return evenList, nil
}

func (s *CarrierAPIBackend) GetTaskEventListByTaskIds(taskIds []string) ([]*pb.TaskEventShow, error) {

	identity, err := s.carrier.carrierDB.GetIdentity()
	if nil != err {
		return nil, err
	}

	evenList := make([]*pb.TaskEventShow, 0)

	// 先查出 task 在本地的 eventList
	for _, taskId := range taskIds {
		localEventList, err := s.carrier.carrierDB.GetTaskEventList(taskId)
		if rawdb.IsNoDBNotFoundErr(err) {
			return nil, err
		}
		if rawdb.IsDBNotFoundErr(err) {
			continue
		}
		for _, e := range localEventList {
			evenList = append(evenList, &pb.TaskEventShow{
				TaskId:   e.TaskId,
				Type:     e.Type,
				CreateAt: e.CreateAt,
				Content:  e.Content,
				Owner: &apipb.Organization{
					NodeName:   identity.NodeName,
					NodeId:     identity.NodeId,
					IdentityId: identity.IdentityId,
				},
			})
		}
	}

	taskEvent, err := s.carrier.carrierDB.GetTaskEventListByTaskIds(taskIds)
	if nil != err {
		return nil, err
	}
	evenList = append(evenList, types.NewTaskEventFromAPIEvent(taskEvent)...)
	return evenList, nil
}

// about DataResourceTable
func (s *CarrierAPIBackend) StoreDataResourceTable(dataResourceTable *types.DataResourceTable) error {
	return s.carrier.carrierDB.StoreDataResourceTable(dataResourceTable)
}

func (s *CarrierAPIBackend) StoreDataResourceTables(dataResourceTables []*types.DataResourceTable) error {
	return s.carrier.carrierDB.StoreDataResourceTables(dataResourceTables)
}

func (s *CarrierAPIBackend) RemoveDataResourceTable(nodeId string) error {
	return s.carrier.carrierDB.RemoveDataResourceTable(nodeId)
}

func (s *CarrierAPIBackend) QueryDataResourceTable(nodeId string) (*types.DataResourceTable, error) {
	return s.carrier.carrierDB.QueryDataResourceTable(nodeId)
}

func (s *CarrierAPIBackend) QueryDataResourceTables() ([]*types.DataResourceTable, error) {
	return s.carrier.carrierDB.QueryDataResourceTables()
}

// about DataResourceFileUpload
func (s *CarrierAPIBackend) StoreDataResourceFileUpload(dataResourceDataUsed *types.DataResourceFileUpload) error {
	return s.carrier.carrierDB.StoreDataResourceFileUpload(dataResourceDataUsed)
}

func (s *CarrierAPIBackend) StoreDataResourceFileUploads(dataResourceDataUseds []*types.DataResourceFileUpload) error {
	return s.carrier.carrierDB.StoreDataResourceFileUploads(dataResourceDataUseds)
}

func (s *CarrierAPIBackend) RemoveDataResourceFileUpload(originId string) error {
	return s.carrier.carrierDB.RemoveDataResourceFileUpload(originId)
}

func (s *CarrierAPIBackend) QueryDataResourceFileUpload(originId string) (*types.DataResourceFileUpload, error) {
	return s.carrier.carrierDB.QueryDataResourceFileUpload(originId)
}

func (s *CarrierAPIBackend) QueryDataResourceDataUseds() ([]*types.DataResourceFileUpload, error) {
	return s.carrier.carrierDB.QueryDataResourceFileUploads()
}

func utilLocalTaskPowerUsedArrString(used []*types.LocalTaskPowerUsed) string {
	arr := make([]string, len(used))
	for i, u := range used {
		arr[i] = u.String()
	}
	if len(arr) != 0 {
		return "[" + strings.Join(arr, ",") + "]"
	}
	return "[]"
}

func utilLocalTaskPowerUsedMapString(taskPowerUsedMap map[string][]*types.LocalTaskPowerUsed) string {
	arr := make([]string, 0)
	for jobNodeId, useds := range taskPowerUsedMap {
		arr = append(arr, fmt.Sprintf(`{"%s": %s}`, jobNodeId, utilLocalTaskPowerUsedArrString(useds)))
	}
	if len(arr) != 0 {
		return "[" + strings.Join(arr, ",") + "]"
	}
	return "[]"
}
