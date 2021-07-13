package carrier

import (
	"errors"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/grpclient"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
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
	if nil != err {
		log.Error("Failed to get all job nodes, on GetNodeInfo(), err:", err)
		return nil, err
	}
	dataNodes, err := s.carrier.carrierDB.GetRegisterNodeList(types.PREFIX_TYPE_DATANODE)
	if nil != err {
		log.Error("Failed to get all data nodes, on GetNodeInfo(), err:", err)
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
	name, err := s.carrier.carrierDB.GetYarnName()
	if nil != err {
		log.Error("Failed to get yarn nodeName, on GetNodeInfo(), err:", err)
		return nil, err
	}
	identity, err := s.carrier.carrierDB.GetIdentityId()
	if nil != err {
		log.Error("Failed to get identity, on GetNodeInfo(), err:", err)
		return nil, err
	}
	seedNodes, err := s.carrier.carrierDB.GetSeedNodeList()
	return &types.YarnNodeInfo{
		NodeType:     types.PREFIX_TYPE_YARNNODE.String(),
		NodeId:       "",       //TODO 读配置
		InternalIp:   "",       // TODO 读配置
		ExternalIp:   "",       // TODO 读p2p
		InternalPort: "",       // TODO 读配置
		ExternalPort: "",       //TODO 读p2p
		IdentityType: "",       // TODO 读配置
		IdentityId:   identity, // TODO 读接口
		Name:         name,
		Peers:        registerNodes,
		SeedPeers:    seedNodes,
		State:        "", // TODO 读系统状态
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
		n := &types.YarnRegisteredJobNode{
			Id:           v.Id,
			InternalIp:   v.InternalIp,
			ExternalIp:   v.ExternalIp,
			InternalPort: v.InternalPort,
			ExternalPort: v.ExternalPort,
			//ResourceUsage:  &types.ResourceUsage{},
			Duration: 0, // TODO 添加运行时长 ...
		}
		n.Task.Count, _ = s.carrier.carrierDB.GetRunningTaskCountOnJobNode(v.Id)
		n.Task.TaskIds, _ = s.carrier.carrierDB.GetJobNodeRunningTaskIdList(v.Id)
		jns[i] = n
	}
	dns := make([]*types.YarnRegisteredDataNode, len(jobNodes))
	for i, v := range dataNodes {
		n := &types.YarnRegisteredDataNode{
			Id:           v.Id,
			InternalIp:   v.InternalIp,
			ExternalIp:   v.ExternalIp,
			InternalPort: v.InternalPort,
			ExternalPort: v.ExternalPort,
			//ResourceUsage:  &types.ResourceUsage{},
			Duration: 0, // TODO 添加运行时长 ...
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
	case types.PREFIX_TYPE_DATANODE:
	case types.PREFIX_TYPE_JOBNODE:
	default:
		return types.NONCONNECTED, errors.New("invalid nodeType")
	}
	if typ == types.PREFIX_TYPE_JOBNODE {
		client, err := grpclient.NewJobNodeClientWithConn(s.carrier.ctx, fmt.Sprintf("%s:%s", node.ExternalIp, node.ExternalPort), node.Id)
		if err != nil {
			return types.NONCONNECTED, err
		}
		s.carrier.resourceClientSet.StoreJobNodeClient(node.Id, client)
	}
	if typ == types.PREFIX_TYPE_DATANODE {
		client, err := grpclient.NewDataNodeClientWithConn(s.carrier.ctx, fmt.Sprintf("%s:%s", node.InternalIp, node.InternalPort), node.Id)
		if err != nil {
			return types.NONCONNECTED, err
		}
		s.carrier.resourceClientSet.StoreDataNodeClient(node.Id, client)
	}
	_, err := s.carrier.carrierDB.SetRegisterNode(typ, node)
	if err != nil {
		return types.NONCONNECTED, err
	}
	node.ConnState = types.CONNECTED
	return types.CONNECTED, nil
}

func (s *CarrierAPIBackend) DeleteRegisterNode(typ types.RegisteredNodeType, id string) error {
	if typ == types.PREFIX_TYPE_JOBNODE {
		if client, ok := s.carrier.resourceClientSet.QueryJobNodeClient(id); ok {
			client.Close()
			s.carrier.resourceClientSet.RemoveJobNodeClient(id)
		}
	}
	if typ == types.PREFIX_TYPE_DATANODE {
		if client, ok := s.carrier.resourceClientSet.QueryDataNodeClient(id); ok {
			client.Close()
			s.carrier.resourceClientSet.RemoveDataNodeClient(id)
		}
	}
	return s.carrier.carrierDB.DeleteRegisterNode(typ, id)
}

func (s *CarrierAPIBackend) GetRegisterNode(typ types.RegisteredNodeType, id string) (*types.RegisteredNodeInfo, error) {
	return s.carrier.carrierDB.GetRegisterNode(typ, id)
}

func (s *CarrierAPIBackend) GetRegisterNodeList(typ types.RegisteredNodeType) ([]*types.RegisteredNodeInfo, error) {
	return s.carrier.carrierDB.GetRegisterNodeList(typ)
}

func (s *CarrierAPIBackend) SendTaskEvent(event *types.TaskEventInfo) error {
	// return s.carrier.resourceManager.SendTaskEvent(evengine)
	return s.carrier.taskManager.SendTaskEvent(event)
}

// metadata api
func (s *CarrierAPIBackend) GetMetaDataDetail(identityId, metaDataId string) (*types.OrgMetaDataInfo, error) {
	metadata, err := s.carrier.carrierDB.GetMetadataByDataId(metaDataId)
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
	return result, nil
}

// power api
func (s *CarrierAPIBackend) GetPowerTotalDetailList() ([]*types.OrgPowerDetail, error) {
	log.Debug("Invoke: GetPowerTotalDetailList executing...")
	resourceList, err := s.carrier.carrierDB.GetResourceList()
	if err != nil {
		return nil, err
	}
	powerList := make([]*types.OrgPowerDetail, 0, resourceList.Len())
	for _, resource := range resourceList.To() {
		powerList = append(powerList, &types.OrgPowerDetail{
			Owner:       &types.NodeAlias{
				Name:       resource.GetNodeName(),
				NodeId:     resource.GetNodeId(),
				IdentityId: resource.GetIdentity(),
			},
			PowerDetail: &types.PowerTotalDetail{
				TotalTaskCount:   0,
				CurrentTaskCount: 0,
				Tasks:            make([]*types.PowerTask, 0),
				ResourceUsage:    &types.ResourceUsage{
					TotalMem:       resource.GetTotalMem(),
					UsedMem:        resource.GetUsedMem(),
					TotalProcessor: resource.GetTotalProcessor(),
					UsedProcessor:  resource.GetUsedProcessor(),
					TotalBandwidth: resource.GetTotalBandWidth(),
					UsedBandwidth:  resource.GetUsedBandWidth(),
				},
				State:            resource.GetState(),
			},
		})
	}
	return powerList, nil
}

func (s *CarrierAPIBackend) GetPowerSingleDetailList() ([]*types.NodePowerDetail, error) {
	// 本机存储的所有计算节点的信息
	return nil, nil // TODO 未完成,  需要查自己参与过的任务信息
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
	if err != nil {
		return nil, err
	}
	localIdentityId, err := s.carrier.carrierDB.GetIdentityId()
	if err != nil {
		return nil, err
	}
	result := make(types.TaskDataArray, 0)
	result = append(result, localTaskArray...)
	// the task has been executed.
	networkTaskList, err := s.carrier.carrierDB.GetTaskList()
	for _, networkTask := range networkTaskList {
		if networkTask.TaskData().GetIdentity() != localIdentityId {
			continue
		}
		result = append(result, networkTask)
	}
	return types.NewTaskDetailShowArrayFromTaskDataArray(result), err
}

func (s *CarrierAPIBackend) GetTaskEventList(taskId string) ([]*types.TaskEvent, error) {
	taskEvent, err := s.carrier.carrierDB.GetTaskEventListByTaskId(taskId)
	return types.NewTaskEventFromAPIEvent(taskEvent), err
}

// about DataResourceTable
func (s *CarrierAPIBackend) StoreDataResourceTable(dataResourceTable *types.DataResourceTable) error {
	return nil
}

func (s *CarrierAPIBackend) StoreDataResourceTables(dataResourceTables []*types.DataResourceTable) error {
	return nil
}

func (s *CarrierAPIBackend) RemoveDataResourceTable(nodeId string) error {
	return nil
}

func (s *CarrierAPIBackend) QueryDataResourceTable(nodeId string) (*types.DataResourceTable, error) {
	return nil, nil
}

func (s *CarrierAPIBackend) QueryDataResourceTables() ([]*types.DataResourceTable, error) {
	return nil, nil
}

// about DataResourceDataUsed
func (s *CarrierAPIBackend) StoreDataResourceDataUsed(dataResourceDataUsed *types.DataResourceDataUsed) error {
	return nil
}

func (s *CarrierAPIBackend) StoreDataResourceDataUseds(dataResourceDataUseds []*types.DataResourceDataUsed) error {
	return nil
}

func (s *CarrierAPIBackend) RemoveDataResourceDataUsed(originId string) error {
	return nil
}

func (s *CarrierAPIBackend) QueryDataResourceDataUsed(originId string) (*types.DataResourceDataUsed, error) {
	return nil, nil
}

func (s *CarrierAPIBackend) QueryDataResourceDataUseds() ([]*types.DataResourceDataUsed, error) {
	return nil, nil
}
