package carrier

import (
	"github.com/RosettaFlow/Carrier-Go/event"
	pbtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
)

// CarrierAPIBackend implements rpc.Backend for Carrier
type CarrierAPIBackend struct {
	carrier *Service
}

func (s *CarrierAPIBackend) SendMsg(msg types.Msg) error {
	return s.carrier.mempool.Add(msg)
}

// system (the yarn node self info)
func (s *CarrierAPIBackend) GetNodeInfo() (*types.YarnNodeInfo, error) {
	jobNodes, err := s.carrier.datachain.GetRegisterNodeList(types.PREFIX_TYPE_JOBNODE)
	if nil != err {
		log.Error("Failed to get all job nodes, on GetNodeInfo(), err:", err)
		return nil, err
	}
	dataNodes, err := s.carrier.datachain.GetRegisterNodeList(types.PREFIX_TYPE_DATANODE)
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
	name, err := s.carrier.datachain.GetYarnName()
	if nil != err {
		log.Error("Failed to get yarn nodeName, on GetNodeInfo(), err:", err)
		return nil, err
	}
	identity, err := s.carrier.datachain.GetIdentity()
	if nil != err {
		log.Error("Failed to get identity, on GetNodeInfo(), err:", err)
		return nil, err
	}
	seedNodes, err := s.carrier.datachain.GetSeedNodeList()
	return &types.YarnNodeInfo{
		NodeType:     types.PREFIX_TYPE_YARNNODE.String(),
		NodeId:       "", //TODO 读配置
		InternalIp:   "", // TODO 读配置
		ExternalIp:   "", // TODO 读p2p
		InternalPort: "", // TODO 读配置
		ExternalPort: "", //TODO 读p2p
		IdentityType: "", // TODO 读配置
		IdentityId:   identity, // TODO 读接口
		Name:         name,
		Peers:        registerNodes,
		SeedPeers:    seedNodes,
		State:        "", // TODO 读系统状态
	}, nil
}
func (s *CarrierAPIBackend) GetRegisteredPeers() (*types.YarnRegisteredNodeDetail, error) {
	// all dataNodes on yarnNode
	dataNodes, err := s.carrier.datachain.GetRegisterNodeList(types.PREFIX_TYPE_DATANODE)
	if nil != err {
		return nil, err
	}
	// all jobNodes on yarnNode
	jobNodes, err := s.carrier.datachain.GetRegisterNodeList(types.PREFIX_TYPE_JOBNODE)
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
		n.Task.Count = s.carrier.datachain.GetRunningTaskCountOnJobNode(v.Id)
		n.Task.TaskIds = s.carrier.datachain.GetJobNodeRunningTaskIdList(v.Id)
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
	return s.carrier.datachain.SetSeedNode(seed)
}

func (s *CarrierAPIBackend) DeleteSeedNode(id string) error {
	return s.carrier.datachain.DeleteSeedNode(id)
}

func (s *CarrierAPIBackend) GetSeedNode(id string) (*types.SeedNodeInfo, error) {
	return s.carrier.datachain.GetSeedNode(id)
}

func (s *CarrierAPIBackend) GetSeedNodeList() ([]*types.SeedNodeInfo, error) {
	return s.carrier.datachain.GetSeedNodeList()
}

func (s *CarrierAPIBackend) SetRegisterNode(typ types.RegisteredNodeType, node *types.RegisteredNodeInfo) (types.NodeConnStatus, error) {
	return s.carrier.datachain.SetRegisterNode(typ, node)
}

func (s *CarrierAPIBackend) DeleteRegisterNode(typ types.RegisteredNodeType, id string) error {
	return s.carrier.datachain.DeleteRegisterNode(typ, id)
}

func (s *CarrierAPIBackend) GetRegisterNode(typ types.RegisteredNodeType, id string) (*types.RegisteredNodeInfo, error) {
	return s.carrier.datachain.GetRegisterNode(typ, id)
}

func (s *CarrierAPIBackend) GetRegisterNodeList(typ types.RegisteredNodeType) ([]*types.RegisteredNodeInfo, error) {
	return s.carrier.datachain.GetRegisterNodeList(typ)
}

func (s *CarrierAPIBackend) SendTaskEvent(event *event.TaskEvent) error {
	return s.carrier.resourceManager.SendTaskEvent(event)
}

// identity api
func (s *CarrierAPIBackend) ApplyIdentityJoin(identity *types.Identity) error {
	return s.carrier.dataCenter.InsertIdentity(identity)
}

func (s *CarrierAPIBackend) RevokeIdentityJoin(identity *types.Identity) error  {
	return s.carrier.dataCenter.RevokeIdentity(identity)
}

// power api
func (s *CarrierAPIBackend) GetPowerTotalSummaryList() ([]*types.OrgResourcePowerAndTaskCount, error) {
	// todo: to be determined.
	return nil, nil
}

func (s *CarrierAPIBackend) GetPowerSingleSummaryList() ([]*types.NodeResourceUsagePowerRes, error) {
	return nil, nil
}

func (s *CarrierAPIBackend) GetPowerTotalSummaryByState(state string) ([]*types.OrgResourcePowerAndTaskCount, error) {
	return nil, nil
}

func (s *CarrierAPIBackend) GetPowerSingleSummaryByState(state string) ([]*types.NodeResourceUsagePowerRes, error) {
	return nil, nil
}

func (s *CarrierAPIBackend) GetPowerTotalSummaryByOwner(identityId string) (*types.OrgResourcePowerAndTaskCount, error) {
	return nil, nil
}

func (s *CarrierAPIBackend) GetPowerSingleSummaryByOwner(identityId string) ([]*types.NodeResourceUsagePowerRes, error) {
	return nil, nil
}

func (s *CarrierAPIBackend) GetPowerSingleDetail(identityId, powerId string) (*types.OrgPowerTaskDetail, error) {
	return nil, nil
}

// metadata api
func (s *CarrierAPIBackend) GetMetaDataSummaryList() ([]*types.OrgMetaDataSummary, error) {
	return nil, nil
}
func (s *CarrierAPIBackend) GetMetaDataSummaryByState(state string) ([]*types.OrgMetaDataSummary, error) {
	return nil, nil
}
func (s *CarrierAPIBackend) GetMetaDataSummaryByOwner(identityId string) ([]*types.OrgMetaDataSummary, error) {
	return nil, nil
}
func (s *CarrierAPIBackend) GetMetaDataDetail(identityId, metaDataId string) ([]types.OrgMetaDataInfo, error) {
	return nil, nil
}

// task api
func (s *CarrierAPIBackend) GetTaskSummaryList() ([]*types.Task, error)       { return nil, nil }
func (s *CarrierAPIBackend) GetTaskJoinSummaryList() ([]*types.Task, error)   { return nil, nil }
func (s *CarrierAPIBackend) GetTaskDetail(taskId string) (*types.Task, error) { return nil, nil }
func (s *CarrierAPIBackend) GetTaskEventList(taskId string) ([]*pbtypes.EventData, error) {
	return nil, nil
}
