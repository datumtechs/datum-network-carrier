package carrier

import (
	"errors"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	"github.com/RosettaFlow/Carrier-Go/grpclient"
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
func (s *CarrierAPIBackend) GetNodeInfo() (*types.YarnNodeInfo, error) {
	jobNodes, err := s.carrier.carrierDB.GetRegisterNodeList(types.PREFIX_TYPE_JOBNODE)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.Errorf("Failed to get all `job nodes`, on GetNodeInfo(), err: {%s}", err)
		return nil, err
	}
	dataNodes, err := s.carrier.carrierDB.GetRegisterNodeList(types.PREFIX_TYPE_DATANODE)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.Errorf("Failed to get all `data nodes, on GetNodeInfo(), err: {%s}", err)
		return nil, err
	}
	jobsLen := len(jobNodes)
	datasLen := len(dataNodes)
	length := jobsLen + datasLen
	registerNodes := make([]*types.RegisteredNodeDetail, length)
	if len(jobNodes) != 0 {
		for i, v := range jobNodes {
			n := &types.RegisteredNodeDetail{
				NodeType: types.PREFIX_TYPE_JOBNODE.String(),
			}
			n.RegisteredNodeInfo = &types.RegisteredNodeInfo{
				Id:           v.Id,
				InternalIp:   v.InternalIp,
				InternalPort: v.InternalPort,
				ExternalIp:   v.ExternalIp,
				ExternalPort: v.ExternalPort,
				ConnState:    v.ConnState,
			}
			registerNodes[i] = n
		}
	}
	if len(dataNodes) != 0 {
		for i, v := range dataNodes {
			n := &types.RegisteredNodeDetail{
				NodeType: types.PREFIX_TYPE_DATANODE.String(),
			}
			n.RegisteredNodeInfo = &types.RegisteredNodeInfo{
				Id:           v.Id,
				InternalIp:   v.InternalIp,
				InternalPort: v.InternalPort,
				ExternalIp:   v.ExternalIp,
				ExternalPort: v.ExternalPort,
				ConnState:    v.ConnState,
			}
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
		nodeName = identity.Name
	}

	seedNodes, err := s.carrier.carrierDB.GetSeedNodeList()
	return &types.YarnNodeInfo{
		NodeType:     types.PREFIX_TYPE_YARNNODE.String(),
		NodeId:       nodeId,
		InternalIp:   "",                             //
		ExternalIp:   "",                             //
		InternalPort: "",                             //
		ExternalPort: "",                             //
		IdentityType: types.IdentityTypeDID.String(), // 默认先是 DID
		IdentityId:   identityId,
		Name:         nodeName,
		Peers:        registerNodes,
		SeedPeers:    seedNodes,
		State:        types.YarnStateActive.String(),
	}, nil
}

func (s *CarrierAPIBackend) GetRegisteredPeers() (*types.YarnRegisteredNodeDetail, error) {
	// all dataNodes on yarnNode
	dataNodes, err := s.carrier.carrierDB.GetRegisterNodeList(types.PREFIX_TYPE_DATANODE)
	if nil != err {
		return nil, err
	}
	// all jobNodes on yarnNode
	jobNodes, err := s.carrier.carrierDB.GetRegisterNodeList(types.PREFIX_TYPE_JOBNODE)
	if nil != err {
		return nil, err
	}
	jns := make([]*types.YarnRegisteredJobNode, len(jobNodes))
	for i, v := range jobNodes {

		var duration uint64
		node, has := s.carrier.resourceClientSet.QueryJobNodeClient(v.Id)
		if has {
			duration = uint64(node.RunningDuration())
		}

		n := &types.YarnRegisteredJobNode{
			Id:           v.Id,
			InternalIp:   v.InternalIp,
			ExternalIp:   v.ExternalIp,
			InternalPort: v.InternalPort,
			ExternalPort: v.ExternalPort,
			//ResourceUsage:  &types.ResourceUsage{},
			Duration: duration, // ms
		}

		n.Task.Count, _ = s.carrier.carrierDB.GetRunningTaskCountOnJobNode(v.Id)
		n.Task.TaskIds, _ = s.carrier.carrierDB.GetJobNodeRunningTaskIdList(v.Id)
		jns[i] = n
	}
	dns := make([]*types.YarnRegisteredDataNode, len(jobNodes))
	for i, v := range dataNodes {

		var duration uint64
		node, has := s.carrier.resourceClientSet.QueryDataNodeClient(v.Id)
		if has {
			duration = uint64(node.RunningDuration())
		}

		n := &types.YarnRegisteredDataNode{
			Id:           v.Id,
			InternalIp:   v.InternalIp,
			ExternalIp:   v.ExternalIp,
			InternalPort: v.InternalPort,
			ExternalPort: v.ExternalPort,
			//ResourceUsage:  &types.ResourceUsage{},
			Duration: duration, // ms
		}
		n.Delta.FileCount = 0
		n.Delta.FileTotalSize = 0
		dns[i] = n
	}
	return &types.YarnRegisteredNodeDetail{
		JobNodes:  jns,
		DataNodes: dns,
	}, nil
}

func (s *CarrierAPIBackend) SetSeedNode(seed *types.SeedNodeInfo) (types.NodeConnStatus, error) {
	//TODO: current node need to connect with seed node.(delay processing)
	return s.carrier.carrierDB.SetSeedNode(seed)
}

func (s *CarrierAPIBackend) DeleteSeedNode(id string) error {
	return s.carrier.carrierDB.DeleteSeedNode(id)
}

func (s *CarrierAPIBackend) GetSeedNode(id string) (*types.SeedNodeInfo, error) {
	return s.carrier.carrierDB.GetSeedNode(id)
}

func (s *CarrierAPIBackend) GetSeedNodeList() ([]*types.SeedNodeInfo, error) {
	return s.carrier.carrierDB.GetSeedNodeList()
}

func (s *CarrierAPIBackend) SetRegisterNode(typ types.RegisteredNodeType, node *types.RegisteredNodeInfo) (types.NodeConnStatus, error) {
	switch typ {
	case types.PREFIX_TYPE_DATANODE, types.PREFIX_TYPE_JOBNODE:
	default:
		return types.NONCONNECTED, errors.New("invalid nodeType")
	}
	if typ == types.PREFIX_TYPE_JOBNODE {
		client, err := grpclient.NewJobNodeClientWithConn(s.carrier.ctx, fmt.Sprintf("%s:%s", node.InternalIp, node.InternalPort), node.Id)
		if err != nil {
			return types.NONCONNECTED, fmt.Errorf("connect new jobNode failed, %s", err)
		}
		s.carrier.resourceClientSet.StoreJobNodeClient(node.Id, client)
	}

	if typ == types.PREFIX_TYPE_DATANODE {
		client, err := grpclient.NewDataNodeClientWithConn(s.carrier.ctx, fmt.Sprintf("%s:%s", node.InternalIp, node.InternalPort), node.Id)
		if err != nil {
			return types.NONCONNECTED, fmt.Errorf("connect new dataNode failed, %s", err)
		}
		s.carrier.resourceClientSet.StoreDataNodeClient(node.Id, client)

		// add data resource  (disk)  todo 后续 需要根据 真实的 dataNode 上报自身的 disk 信息
		err = s.carrier.carrierDB.StoreDataResourceTable(types.NewDataResourceTable(node.Id, types.DefaultDisk, 0))
		if err != nil {
			return types.NONCONNECTED, fmt.Errorf("store disk summary of new dataNode failed, %s", err)
		}
	}
	node.ConnState = types.CONNECTED
	_, err := s.carrier.carrierDB.SetRegisterNode(typ, node)
	if err != nil {
		return types.NONCONNECTED, fmt.Errorf("Store registerNode to db failed, %s", err)
	}
	return types.CONNECTED, nil
}

func (s *CarrierAPIBackend) UpdateRegisterNode(typ types.RegisteredNodeType, node *types.RegisteredNodeInfo) (types.NodeConnStatus, error) {

	switch typ {
	case types.PREFIX_TYPE_DATANODE, types.PREFIX_TYPE_JOBNODE:
	default:
		return types.NONCONNECTED, errors.New("invalid nodeType")
	}
	if typ == types.PREFIX_TYPE_JOBNODE {

		// 算力已发布的jobNode不可以直接删除
		resourceTable, err := s.carrier.carrierDB.QueryLocalResourceTable(node.Id)
		if rawdb.IsNoDBNotFoundErr(err) {
			return types.NONCONNECTED, fmt.Errorf("query local power resource on old jobNode failed, %s", err)
		}

		if nil != resourceTable {
			log.Debugf("still have the published computing power information on old jobNode on UpdateRegisterNode, %s", resourceTable.String())
			return types.NONCONNECTED, fmt.Errorf("still have the published computing power information on old jobNode failed, input jobNodeId: {%s}, old jobNodeId: {%s}, old powerId: {%s}",
				node.Id, resourceTable.GetNodeId(), resourceTable.GetPowerId())
		}

		// 先校验 jobNode 上是否有正在执行的 task
		runningTaskCount, err := s.carrier.carrierDB.GetRunningTaskCountOnJobNode(node.Id)
		if rawdb.IsNoDBNotFoundErr(err) {
			return types.NONCONNECTED, fmt.Errorf("query local running taskCount on old jobNode failed, %s", err)
		}
		if runningTaskCount > 0 {
			return types.NONCONNECTED, fmt.Errorf("the old jobNode have been running {%d} task current, don't remove it", runningTaskCount)
		}

		if client, ok := s.carrier.resourceClientSet.QueryJobNodeClient(node.Id); ok {
			// remove old client instanse
			client.Close()
			s.carrier.resourceClientSet.RemoveJobNodeClient(node.Id)
		}

		// generate new client
		client, err := grpclient.NewJobNodeClientWithConn(s.carrier.ctx, fmt.Sprintf("%s:%s", node.InternalIp, node.InternalPort), node.Id)
		if err != nil {
			return types.NONCONNECTED, fmt.Errorf("connect new jobNode failed, %s", err)
		}
		s.carrier.resourceClientSet.StoreJobNodeClient(node.Id, client)

	}

	if typ == types.PREFIX_TYPE_DATANODE {

		// 先校验 dataNode 上是否已被 使用
		dataNodeTable, err := s.carrier.carrierDB.QueryDataResourceTable(node.Id)
		if rawdb.IsNoDBNotFoundErr(err) {
			return types.NONCONNECTED, fmt.Errorf("query disk used summary on old dataNode failed, %s", err)
		}
		if dataNodeTable.IsNotEmpty() && dataNodeTable.IsUsed() {
			return types.NONCONNECTED, fmt.Errorf("the disk of old dataNode was used, don't remove it, totalDisk: {%d byte}, usedDisk: {%d byte}, remainDisk: {%d byte}",
				dataNodeTable.GetTotalDisk(), dataNodeTable.GetUsedDisk(), dataNodeTable.RemainDisk())
		}

		if client, ok := s.carrier.resourceClientSet.QueryDataNodeClient(node.Id); ok {
			// remove old client instanse
			client.Close()
			s.carrier.resourceClientSet.RemoveDataNodeClient(node.Id)
		}

		// remove old data resource  (disk)  todo 后续 需要根据 真实的 dataNode 上报自身的 disk 信息
		if err := s.carrier.carrierDB.RemoveDataResourceTable(node.Id); rawdb.IsNoDBNotFoundErr(err) {
			return types.NONCONNECTED, fmt.Errorf("remove disk summary of old dataNode, %s", err)
		}

		client, err := grpclient.NewDataNodeClientWithConn(s.carrier.ctx, fmt.Sprintf("%s:%s", node.InternalIp, node.InternalPort), node.Id)
		if err != nil {
			return types.NONCONNECTED, fmt.Errorf("connect new dataNode failed, %s", err)
		}
		s.carrier.resourceClientSet.StoreDataNodeClient(node.Id, client)

		// add new data resource  (disk)  todo 后续 需要根据 真实的 dataNode 上报自身的 disk 信息
		err = s.carrier.carrierDB.StoreDataResourceTable(types.NewDataResourceTable(node.Id, types.DefaultDisk, 0))
		if err != nil {
			return types.NONCONNECTED, fmt.Errorf("store disk summary of new dataNode failed, %s", err)
		}

	}

	// remove  old jobNode from db
	if err := s.carrier.carrierDB.DeleteRegisterNode(typ, node.Id); nil != err {
		return types.NONCONNECTED, fmt.Errorf("remove old registerNode from db failed, %s", err)
	}

	// add new node to db
	node.ConnState = types.CONNECTED
	_, err := s.carrier.carrierDB.SetRegisterNode(typ, node)
	if err != nil {
		return types.NONCONNECTED, fmt.Errorf("Store new registerNode to db failed, %s", err)
	}
	return types.CONNECTED, nil
}

func (s *CarrierAPIBackend) DeleteRegisterNode(typ types.RegisteredNodeType, id string) error {
	switch typ {
	case types.PREFIX_TYPE_DATANODE, types.PREFIX_TYPE_JOBNODE:
	default:
		return errors.New("invalid nodeType")
	}
	if typ == types.PREFIX_TYPE_JOBNODE {

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

	if typ == types.PREFIX_TYPE_DATANODE {

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

func (s *CarrierAPIBackend) GetRegisterNode(typ types.RegisteredNodeType, id string) (*types.RegisteredNodeInfo, error) {
	return s.carrier.carrierDB.GetRegisterNode(typ, id)
}

func (s *CarrierAPIBackend) GetRegisterNodeList(typ types.RegisteredNodeType) ([]*types.RegisteredNodeInfo, error) {
	nodeList, err := s.carrier.carrierDB.GetRegisterNodeList(typ)
	if nil != err {
		return nil, err
	}

	// 需要处理 计算服务 信息
	if typ == types.PREFIX_TYPE_JOBNODE {
		for i, jobNode := range nodeList {

			client, ok := s.carrier.resourceClientSet.QueryJobNodeClient(jobNode.Id)
			if !ok {
				jobNode.ConnState = types.NONCONNECTED
			}
			if !client.IsConnected() {
				jobNode.ConnState = types.NONCONNECTED
			} else {
				jobNode.ConnState = types.CONNECTED
			}

			table, err := s.carrier.carrierDB.QueryLocalResourceTable(jobNode.Id)
			if nil != err {
				continue
			}
			if nil != table {
				jobNode.ConnState = types.ENABLE_POWER
			}
			taskCount, err := s.carrier.carrierDB.GetRunningTaskCountOnJobNode(jobNode.Id)
			if nil != err {
				continue
			}
			if taskCount > 0 {
				jobNode.ConnState = types.BUSY_POWER
			}
			nodeList[i] = jobNode
		}
	}

	return nodeList, nil
}

func (s *CarrierAPIBackend) SendTaskEvent(event *types.TaskEventInfo) error {
	// return s.carrier.resourceManager.SendTaskEvent(evengine)
	return s.carrier.taskManager.SendTaskEvent(event)
}

// metadata api
func (s *CarrierAPIBackend) GetMetaDataDetail(identityId, metaDataId string) (*types.OrgMetaDataInfo, error) {
	metadata, err := s.carrier.carrierDB.GetMetadataByDataId(metaDataId)
	if metadata == nil {
		return nil, errors.New("not found metadata by special Id")
	}
	return types.NewOrgMetaDataInfoFromMetadata(metadata), err
}

// GetMetaDataDetailList returns a list of all metadata details in the network.
func (s *CarrierAPIBackend) GetMetaDataDetailList() ([]*types.OrgMetaDataInfo, error) {
	metadataArray, err := s.carrier.carrierDB.GetMetadataList()
	return types.NewOrgMetaDataInfoArrayFromMetadataArray(metadataArray), err
}

func (s *CarrierAPIBackend) GetMetaDataDetailListByOwner(identityId string) ([]*types.OrgMetaDataInfo, error) {
	log.WithField("identityId", identityId).Debug("Invoke: GetMetaDataDetailListByOwner executing...")
	metadataList, err := s.GetMetaDataDetailList()
	if err != nil {
		return nil, err
	}
	result := make([]*types.OrgMetaDataInfo, 0)
	for _, metadata := range metadataList {
		if metadata.Owner.IdentityId == identityId {
			result = append(result, metadata)
		}
	}
	log.Debugf("Query metaData list, identityId: {%s}, len: {%d}", identityId, len(result))
	return result, nil
}

// power api
func (s *CarrierAPIBackend) GetPowerTotalDetailList() ([]*types.OrgPowerDetail, error) {
	log.Debug("Invoke: GetPowerTotalDetailList executing...")
	resourceList, err := s.carrier.carrierDB.GetResourceList()
	if err != nil {
		return nil, err
	}
	log.Debugf("Query all org's power list, len: {%d}", len(resourceList))
	powerList := make([]*types.OrgPowerDetail, 0, resourceList.Len())
	for _, resource := range resourceList.To() {
		powerList = append(powerList, &types.OrgPowerDetail{
			Owner: &types.NodeAlias{
				Name:       resource.GetNodeName(),
				NodeId:     resource.GetNodeId(),
				IdentityId: resource.GetIdentity(),
			},
			PowerDetail: &types.PowerTotalDetail{
				TotalTaskCount:   0,
				CurrentTaskCount: 0,
				Tasks:            make([]*types.PowerTask, 0),
				ResourceUsage: &types.ResourceUsage{
					TotalMem:       resource.GetTotalMem(),
					UsedMem:        resource.GetUsedMem(),
					TotalProcessor: resource.GetTotalProcessor(),
					UsedProcessor:  resource.GetUsedProcessor(),
					TotalBandwidth: resource.GetTotalBandWidth(),
					UsedBandwidth:  resource.GetUsedBandWidth(),
				},
				State: resource.GetState(),
			},
		})
	}
	return powerList, nil
}

func (s *CarrierAPIBackend) GetPowerSingleDetailList() ([]*types.NodePowerDetail, error) {
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

	buildPowerTaskList := func(jobNodeId string) []*types.PowerTask {
		powerTaskList := make([]*types.PowerTask, 0)

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
				powerTask := &types.PowerTask{
					TaskId:   taskId,
					TaskName: task.TaskData().TaskName,
					Owner: &types.NodeAlias{
						Name:       task.TaskData().GetNodeName(),
						NodeId:     task.TaskData().GetNodeId(),
						IdentityId: task.TaskData().GetIdentity(),
					},
					Patners:   make([]*types.NodeAlias, 0),
					Receivers: make([]*types.NodeAlias, 0),
					OperationCost: &types.TaskOperationCost{
						Processor: uint64(task.TaskData().GetTaskResource().GetCostProcessor()),
						Mem:       uint64(task.TaskData().GetTaskResource().GetCostMem()),
						Bandwidth: uint64(task.TaskData().GetTaskResource().GetCostBandwidth()),
						Duration:  uint64(task.TaskData().GetTaskResource().GetDuration()),
					},
					OperationSpend: nil, // 下面单独计算 任务资源使用 实况 ...
					CreateAt:       task.TaskData().CreateAt,
				}
				// 组装 数据参与方
				for _, dataSupplier := range task.TaskData().MetadataSupplier {
					// 协作方, 需要过滤掉自己
					if task.TaskData().GetNodeId() != dataSupplier.Organization.Identity {
						powerTask.Patners = append(powerTask.Patners, &types.NodeAlias{
							Name:       dataSupplier.Organization.GetNodeName(),
							NodeId:     dataSupplier.Organization.GetNodeId(),
							IdentityId: dataSupplier.Organization.GetIdentity(),
						})
					}

				}
				// 组装结果接收方
				for _, receiver := range task.TaskData().Receivers {
					powerTask.Receivers = append(powerTask.Receivers, &types.NodeAlias{
						Name:       receiver.Receiver.GetNodeName(),
						NodeId:     receiver.Receiver.GetNodeId(),
						IdentityId: receiver.Receiver.GetIdentity(),
					})
				}

				// 计算任务使用实况 ...
				slotCount := readElement(jobNodeId, powerTask.TaskId)
				powerTask.OperationSpend = &types.TaskOperationCost{
					Processor: slotUnit.Processor * slotCount,
					Mem:       slotUnit.Mem * slotCount,
					Bandwidth: slotUnit.Bandwidth * slotCount,
					Duration:  task.TaskData().GetTaskResource().GetDuration(),
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
	result := make([]*types.NodePowerDetail, len(resourceList))
	for i, resource := range resourceList {
		nodePowerDetail := &types.NodePowerDetail{
			Owner: &types.NodeAlias{
				Name:       resource.GetNodeName(),
				NodeId:     resource.GetNodeId(),
				IdentityId: resource.GetIdentity(),
			},
			PowerDetail: &types.PowerSingleDetail{
				JobNodeId:        resource.GetJobNodeId(),
				PowerId:          resource.DataId,
				TotalTaskCount:   uint32(taskCount(resource.GetJobNodeId())),
				CurrentTaskCount: uint32(taskCount(resource.GetJobNodeId())),
				Tasks:            make([]*types.PowerTask, 0),
				ResourceUsage: &types.ResourceUsage{
					TotalMem:       resource.GetTotalMem(),
					UsedMem:        resource.GetUsedMem(),
					TotalProcessor: resource.GetTotalProcessor(),
					UsedProcessor:  resource.GetUsedProcessor(),
					TotalBandwidth: resource.GetTotalBandWidth(),
					UsedBandwidth:  resource.GetUsedBandWidth(),
				},
				State: resource.GetState(),
			},
		}
		powerTaskArray := buildPowerTaskList(resource.GetJobNodeId())
		nodePowerDetail.PowerDetail.Tasks = powerTaskArray
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
	return types.NewIdentity(&libTypes.IdentityData{
		Identity: nodeAlias.IdentityId,
		NodeId:   nodeAlias.NodeId,
		NodeName: nodeAlias.Name,
	}), err
}

func (s *CarrierAPIBackend) GetIdentityList() ([]*types.Identity, error) {
	return s.carrier.carrierDB.GetIdentityList()
}

// task api
func (s *CarrierAPIBackend) GetTaskDetailList() ([]*types.TaskDetailShow, error) {
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

	makeTaskViewFn := func(task *types.Task) *types.TaskDetailShow {
		// task 发起方
		if task.TaskData().GetIdentity() == localIdentityId {
			return types.NewTaskDetailShowFromTaskData(task, types.TaskRoleOwner.String())
		}

		// task 参与方
		for _, dataSupplier := range task.TaskData().MetadataSupplier {
			if dataSupplier.Organization.Identity == localIdentityId {
				return types.NewTaskDetailShowFromTaskData(task, types.TaskRoleDataSupplier.String())
			}
		}

		// 算力提供方
		for _, powerSupplier := range task.TaskData().ResourceSupplier {
			if powerSupplier.Organization.Identity == localIdentityId {
				return types.NewTaskDetailShowFromTaskData(task, types.TaskRolePowerSupplier.String())
			}
		}

		// 数据接收方
		for _, receiver := range task.TaskData().Receivers {
			if receiver.Receiver.Identity == localIdentityId {
				return types.NewTaskDetailShowFromTaskData(task, types.TaskRoleReceiver.String())
			}
		}
		return nil
	}

	result := make([]*types.TaskDetailShow, 0)
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

func (s *CarrierAPIBackend) GetTaskEventList(taskId string) ([]*types.TaskEvent, error) {

	identity, err := s.carrier.carrierDB.GetIdentity()
	if nil != err {
		return nil, err
	}

	// 先查出 task 在本地的 eventList
	localEventList, err := s.carrier.carrierDB.GetTaskEventList(taskId)
	if rawdb.IsNoDBNotFoundErr(err) {
		return nil, err
	}

	evenList := make([]*types.TaskEvent, len(localEventList))
	for i, e := range localEventList {
		evenList[i] = &types.TaskEvent{
			TaskId:   e.TaskId,
			Type:     e.Type,
			CreateAt: e.CreateTime,
			Content:  e.Content,
			Owner: &types.NodeAlias{
				Name:       identity.Name,
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

func (s *CarrierAPIBackend) GetTaskEventListByTaskIds(taskIds []string) ([]*types.TaskEvent, error) {

	identity, err := s.carrier.carrierDB.GetIdentity()
	if nil != err {
		return nil, err
	}

	evenList := make([]*types.TaskEvent, 0)

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
			evenList = append(evenList, &types.TaskEvent{
				TaskId:   e.TaskId,
				Type:     e.Type,
				CreateAt: e.CreateTime,
				Content:  e.Content,
				Owner: &types.NodeAlias{
					Name:       identity.Name,
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
		return "[" +  strings.Join(arr, ",") + "]"
	}
	return "[]"
}

func utilLocalTaskPowerUsedMapString(taskPowerUsedMap map[string][]*types.LocalTaskPowerUsed) string {
	arr := make([]string, 0)
	for jobNodeId, useds := range taskPowerUsedMap {
		arr = append(arr, fmt.Sprintf(`{"%s": %s}`, jobNodeId, utilLocalTaskPowerUsedArrString(useds)))
	}
	if len(arr) != 0 {
		return "[" +  strings.Join(arr, ",") + "]"
	}
	return "[]"
}