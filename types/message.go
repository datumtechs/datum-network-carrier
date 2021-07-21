package types

import (
	"encoding/json"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/rlputil"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	"github.com/RosettaFlow/Carrier-Go/lib/types"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"sync/atomic"
	"time"
)

const (
	PREFIX_POWER_ID    = "power:"
	PREFIX_METADATA_ID = "metadata:"
	PREFIX_TASK_ID     = "task:"
)
const (
	MSG_IDENTITY        = "identityMsg"
	MSG_IDENTITY_REVOKE = "identityRevokeMsg"
	MSG_POWER           = "powerMsg"
	MSG_POWER_REVOKE    = "powerRevokeMsg"
	MSG_METADATA        = "metaDataMsg"
	MSG_METADATA_REVOKE = "metaDataRevokeMsg"
	MSG_TASK            = "taskMsg"
)

type MsgType string

type Msg interface {
	Marshal() ([]byte, error)
	Unmarshal(b []byte) error
	String() string
	MsgType() string
}

// ------------------- SeedNode -------------------
type SeedNodeMsg struct {
	InternalIp   string `json:"internalIp"`
	InternalPort string `json:"internalPort"`
}

// ------------------- JobNode AND DataNode -------------------
type RegisteredNodeMsg struct {
	InternalIp   string `json:"internalIp"`
	InternalPort string `json:"internalPort"`
	ExternalIp   string `json:"externalIp"`
	ExternalPort string `json:"externalPort"`
}

// ------------------- identity -------------------
type IdentityMsg struct {
	*NodeAlias
	CreateAt uint64 `json:"createAt"`
}
type IdentityRevokeMsg struct {
	CreateAt uint64 `json:"createAt"`
}

type IdentityMsgs []*IdentityMsg
type IdentityRevokeMsgs []*IdentityRevokeMsg

func (msg *IdentityMsg) ToDataCenter() *Identity {
	return NewIdentity(&libTypes.IdentityData{
		NodeName: msg.Name,
		NodeId:   msg.NodeId,
		Identity: msg.IdentityId,
	})
}
func (msg *IdentityMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *IdentityMsg) Unmarshal(b []byte) error { return nil }
func (msg *IdentityMsg) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *IdentityMsg) MsgType() string                { return MSG_IDENTITY }
func (msg *IdentityMsg) OwnerName() string              { return msg.Name }
func (msg *IdentityMsg) OwnerNodeId() string            { return msg.NodeId }
func (msg *IdentityMsg) OwnerIdentityId() string        { return msg.IdentityId }
func (msg *IdentityMsg) MsgCreateAt() uint64            { return msg.CreateAt }
func (msg *IdentityRevokeMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *IdentityRevokeMsg) Unmarshal(b []byte) error { return nil }
func (msg *IdentityRevokeMsg) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *IdentityRevokeMsg) MsgType() string { return MSG_IDENTITY_REVOKE }

// ------------------- power -------------------

type PowerMsg struct {
	// This is only used when marshaling to JSON.
	PowerId   string `json:"powerId"`
	JobNodeId string `json:"jobNodeId"`
	CreateAt  uint64 `json:"createAt"`
	// caches
	hash atomic.Value
}

func NewPowerMessageFromRequest(req *pb.PublishPowerRequest) *PowerMsg {
	msg := &PowerMsg{
		JobNodeId: req.JobNodeId,
		CreateAt: uint64(time.Now().UnixNano()),
	}
	msg.GetPowerId()
	return msg
}

type powerData struct {
	*NodeAlias
	JobNodeId   string `json:"jobNodeId"`
	Information struct {
		Mem       uint64 `json:"mem,omitempty"`
		Processor uint64 `json:"processor,omitempty"`
		Bandwidth uint64 `json:"bandwidth,omitempty"`
	} `json:"information"`
	CreateAt uint64 `json:"createAt"`
}
type PowerRevokeMsg struct {
	*NodeAlias
	PowerId  string `json:"powerId"`
	CreateAt uint64 `json:"createAt"`
}

func NewPowerRevokeMessageFromRequest(req *pb.RevokePowerRequest) *PowerRevokeMsg {
	return &PowerRevokeMsg{
		PowerId: req.PowerId,
		CreateAt: uint64(time.Now().UnixNano()),
	}
}

type PowerMsgs []*PowerMsg
type PowerRevokeMsgs []*PowerRevokeMsg

//func (msg *PowerMsg) ToLocal() *LocalResource {
//	return NewLocalResource(&libTypes.LocalResourceData{
//		Identity:  msg.OwnerIdentityId(),
//		NodeId:    msg.OwnerNodeId(),
//		NodeName:  msg.OwnerName(),
//		JobNodeId: msg.JobNodeId(),
//		DataId:    msg.PowerId,
//		// the status of data, N means normal, D means deleted.
//		DataStatus: DataStatusNormal.String(),
//		// resource status, eg: create/release/revoke
//		State: PowerStateRelease.String(),
//		// unit: byte
//		TotalMem: msg.Memory(),
//		// unit: byte
//		UsedMem: 0,
//		// number of cpu cores.
//		TotalProcessor: msg.Processor(),
//		UsedProcessor:  0,
//		// unit: byte
//		TotalBandWidth: msg.Bandwidth(),
//		UsedBandWidth:  0,
//	})
//}
//func (msg *PowerMsg) ToDataCenter() *Resource {
//	return NewResource(&libTypes.ResourceData{
//		Identity: msg.OwnerIdentityId(),
//		NodeId:   msg.OwnerNodeId(),
//		NodeName: msg.OwnerName(),
//		DataId:   msg.PowerId,
//		// the status of data, N means normal, D means deleted.
//		DataStatus: DataStatusNormal.String(),
//		// resource status, eg: create/release/revoke
//		State: PowerStateRelease.String(),
//		// unit: byte
//		TotalMem: msg.Memory(),
//		// unit: byte
//		UsedMem: 0,
//		// number of cpu cores.
//		TotalProcessor: msg.Processor(),
//		UsedProcessor:  0,
//		// unit: byte
//		TotalBandWidth: msg.Bandwidth(),
//		UsedBandWidth:  0,
//	})
//}
func (msg *PowerMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *PowerMsg) Unmarshal(b []byte) error { return nil }
func (msg *PowerMsg) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *PowerMsg) MsgType() string { return MSG_POWER }

func (msg *PowerMsg) GetJobNodeId() string    { return msg.JobNodeId }
func (msg *PowerMsg) GetCreateAt() uint64        { return msg.CreateAt }
func (msg *PowerMsg) GetPowerId() string {
	if "" != msg.PowerId {
		return msg.PowerId
	}
	msg.PowerId = PREFIX_POWER_ID + msg.HashByCreateTime().Hex()
	return msg.PowerId
}

func (msg *PowerMsg) Hash() common.Hash {
	if hash := msg.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlputil.RlpHash([]interface{}{
		msg.JobNodeId,
	})
	msg.hash.Store(v)
	return v
}

func (msg *PowerMsg) HashByCreateTime() common.Hash {
	if hash := msg.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlputil.RlpHash([]interface{}{
		msg.JobNodeId,
		msg.CreateAt,
		//uint64(time.Now().UnixNano()),
	})
	msg.hash.Store(v)
	return v
}

//func (msg *PowerRevokeMsg) ToDataCenter() *Resource {
//	return NewResource(&libTypes.ResourceData{
//		Identity: msg.IdentityId,
//		NodeId:   msg.NodeId,
//		NodeName: msg.Name,
//		DataId:   msg.PowerId,
//		// the status of data, N means normal, D means deleted.
//		DataStatus: DataStatusDeleted.String(),
//		// resource status, eg: create/release/revoke
//		State: PowerStateRevoke.String(),
//		// unit: byte
//		TotalMem: 0,
//		// unit: byte
//		UsedMem: 0,
//		// number of cpu cores.
//		TotalProcessor: 0,
//		UsedProcessor:  0,
//		// unit: byte
//		TotalBandWidth: 0,
//		UsedBandWidth:  0,
//	})
//}
func (msg *PowerRevokeMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *PowerRevokeMsg) Unmarshal(b []byte) error { return nil }
func (msg *PowerRevokeMsg) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *PowerRevokeMsg) MsgType() string { return MSG_POWER_REVOKE }

// Len returns the length of s.
func (s PowerMsgs) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s PowerMsgs) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s PowerMsgs) Less(i, j int) bool { return s[i].CreateAt < s[j].CreateAt }

// ------------------- metaData -------------------

type MetaDataMsg struct {
	MetaDataId string `json:"metaDataId"`
	Data       *metadataData
	// caches
	hash atomic.Value
}

func NewMetaDataMessageFromRequest(req *pb.PublishMetaDataRequest) *MetaDataMsg {
	return &MetaDataMsg{
		Data: &metadataData{
			NodeAlias: &NodeAlias{
				Name:       req.Owner.Name,
				NodeId:     req.Owner.NodeId,
				IdentityId: req.Owner.IdentityId,
			},
			Information: struct {
				MetaDataSummary *MetaDataSummary    `json:"metaDataSummary"`
				ColumnMetas     []*types.ColumnMeta `json:"columnMetas"`
			}{
				MetaDataSummary: &MetaDataSummary{
					MetaDataId: req.Information.MetaDataSummary.MetaDataId,
					OriginId:   req.Information.MetaDataSummary.OriginId,
					TableName:  req.Information.MetaDataSummary.TableName,
					Desc:       req.Information.MetaDataSummary.Desc,
					FilePath:   req.Information.MetaDataSummary.FilePath,
					Rows:       req.Information.MetaDataSummary.Rows,
					Columns:    req.Information.MetaDataSummary.Columns,
					Size:       req.Information.MetaDataSummary.Size_,
					FileType:   req.Information.MetaDataSummary.FileType,
					HasTitle:   req.Information.MetaDataSummary.HasTitle,
					State:      req.Information.MetaDataSummary.State,
				},
				ColumnMetas: make([]*types.ColumnMeta, 0),
			},
			CreateAt: uint64(time.Now().UnixNano()),
		},
	}
}

type metadataData struct {
	*NodeAlias
	Information struct {
		MetaDataSummary *MetaDataSummary    `json:"metaDataSummary"`
		ColumnMetas     []*types.ColumnMeta `json:"columnMetas"`
	} `json:"information"`
	CreateAt uint64 `json:"createAt"`
}

type MetaDataSummary struct {
	MetaDataId string `json:"metaDataId,omitempty"`
	OriginId   string `json:"originId,omitempty"`
	TableName  string `json:"tableName,omitempty"`
	Desc       string `json:"desc,omitempty"`
	FilePath   string `json:"filePath,omitempty"`
	Rows       uint32 `json:"rows,omitempty"`
	Columns    uint32 `json:"columns,omitempty"`
	Size       uint32 `json:"size,omitempty"`
	FileType   string `json:"fileType,omitempty"`
	HasTitle   bool   `json:"hasTitle,omitempty"`
	State      string `json:"state,omitempty"`
}

//type ColumnMeta struct {
//	Cindex   uint64 `json:"cindex,omitempty"`
//	Cname    string `json:"cname,omitempty"`
//	Ctype    string `json:"ctype,omitempty"`
//	Csize    uint32 `json:"csize,omitempty"`
//	Ccomment string `json:"ccomment,omitempty"`
//}
type MetaDataRevokeMsg struct {
	*NodeAlias
	MetaDataId string `json:"metaDataId"`
	CreateAt   uint64 `json:"createAt"`
}

func NewMetadataRevokeMessageFromRequest(req *pb.RevokeMetaDataRequest) *MetaDataRevokeMsg {
	return &MetaDataRevokeMsg{
		NodeAlias: &NodeAlias{
			Name:       req.Owner.Name,
			NodeId:     req.Owner.NodeId,
			IdentityId: req.Owner.IdentityId,
		},
		MetaDataId: req.MetaDataId,
		CreateAt: uint64(time.Now().UnixNano()),
	}
}

type MetaDataMsgs []*MetaDataMsg
type MetaDataRevokeMsgs []*MetaDataRevokeMsg

func (msg *MetaDataMsg) ToDataCenter() *Metadata {
	return NewMetadata(&libTypes.MetaData{
		Identity:       msg.OwnerIdentityId(),
		NodeId:         msg.OwnerNodeId(),
		NodeName:       msg.OwnerName(),
		DataId:         msg.MetaDataId,
		OriginId:       msg.OriginId(),
		TableName:      msg.TableName(),
		FilePath:       msg.FilePath(),
		FileType:       msg.FileType(),
		Desc:           msg.Desc(),
		Rows:           uint64(msg.Rows()),
		Columns:        uint64(msg.Columns()),
		Size_:          uint64(msg.Size()),
		HasTitleRow:    msg.HasTitle(),
		ColumnMetaList: msg.ColumnMetas(),
		// the status of data, N means normal, D means deleted.
		DataStatus: DataStatusNormal.String(),
		// metaData status, eg: create/release/revoke
		State: MetaDataStateRelease.String(),
	})
}
func (msg *MetaDataMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *MetaDataMsg) Unmarshal(b []byte) error { return nil }
func (msg *MetaDataMsg) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *MetaDataMsg) MsgType() string { return MSG_METADATA }
func (msg *MetaDataMsg) Onwer() *NodeAlias {
	return &NodeAlias{
		Name:       msg.Data.Name,
		NodeId:     msg.Data.NodeId,
		IdentityId: msg.Data.IdentityId,
	}
}
func (msg *MetaDataMsg) OwnerName() string       { return msg.Data.Name }
func (msg *MetaDataMsg) OwnerNodeId() string     { return msg.Data.NodeId }
func (msg *MetaDataMsg) OwnerIdentityId() string { return msg.Data.IdentityId }
func (msg *MetaDataMsg) MetaDataSummary() *MetaDataSummary {
	return msg.Data.Information.MetaDataSummary
}
func (msg *MetaDataMsg) OriginId() string                 { return msg.Data.Information.MetaDataSummary.OriginId }
func (msg *MetaDataMsg) TableName() string                { return msg.Data.Information.MetaDataSummary.TableName }
func (msg *MetaDataMsg) Desc() string                     { return msg.Data.Information.MetaDataSummary.Desc }
func (msg *MetaDataMsg) FilePath() string                 { return msg.Data.Information.MetaDataSummary.FilePath }
func (msg *MetaDataMsg) Rows() uint32                     { return msg.Data.Information.MetaDataSummary.Rows }
func (msg *MetaDataMsg) Columns() uint32                  { return msg.Data.Information.MetaDataSummary.Columns }
func (msg *MetaDataMsg) Size() uint32                     { return msg.Data.Information.MetaDataSummary.Size }
func (msg *MetaDataMsg) FileType() string                 { return msg.Data.Information.MetaDataSummary.FileType }
func (msg *MetaDataMsg) HasTitle() bool                   { return msg.Data.Information.MetaDataSummary.HasTitle }
func (msg *MetaDataMsg) State() string                    { return msg.Data.Information.MetaDataSummary.State }
func (msg *MetaDataMsg) ColumnMetas() []*types.ColumnMeta { return msg.Data.Information.ColumnMetas }
func (msg *MetaDataMsg) CreateAt() uint64                 { return msg.Data.CreateAt }
func (msg *MetaDataMsg) GetMetaDataId() string {
	if "" != msg.MetaDataId {
		return msg.MetaDataId
	}
	msg.MetaDataId = PREFIX_METADATA_ID + msg.Hash().Hex()
	return msg.MetaDataId
}
func (msg *MetaDataMsg) Hash() common.Hash {
	if hash := msg.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlputil.RlpHash(msg.Data)
	msg.hash.Store(v)
	return v
}

func (msg *MetaDataRevokeMsg) ToDataCenter() *Metadata {
	return NewMetadata(&libTypes.MetaData{
		Identity: msg.IdentityId,
		NodeId:   msg.NodeId,
		NodeName: msg.Name,
		DataId:   msg.MetaDataId,
		// the status of data, N means normal, D means deleted.
		DataStatus: DataStatusDeleted.String(),
		// metaData status, eg: create/release/revoke
		State: MetaDataStateRevoke.String(),
	})
}
func (msg *MetaDataRevokeMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *MetaDataRevokeMsg) Unmarshal(b []byte) error { return nil }
func (msg *MetaDataRevokeMsg) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *MetaDataRevokeMsg) MsgType() string { return MSG_METADATA_REVOKE }

// Len returns the length of s.
func (s MetaDataMsgs) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s MetaDataMsgs) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s MetaDataMsgs) Less(i, j int) bool { return s[i].Data.CreateAt < s[j].Data.CreateAt }

// ------------------- task -------------------

type TaskBullet struct {
	UnschedTask *UnSchedTaskWrap
	Starve      bool
	Term        uint32
	Resched     uint32
}

func NewTaskBulletByTaskMsg(msg *TaskMsg) *TaskBullet {
	return &TaskBullet{
		UnschedTask: NewUnSchedTaskWrap(msg.Data, msg.PowerPartyIds),
	}
}

func (b *TaskBullet) IncreaseResched() { b.Resched++ }
func (b *TaskBullet) DecreaseResched() {
	if b.Resched > 0 {
		b.Resched--
	}
}
func (b *TaskBullet) IncreaseTerm() { b.Term++ }
func (b *TaskBullet) DecreaseTerm() {
	if b.Term > 0 {
		b.Term--
	}
}

type TaskBullets []*TaskBullet

func (h TaskBullets) Len() int           { return len(h) }
func (h TaskBullets) Less(i, j int) bool { return h[i].Term > h[j].Term } // term:  a.3 > c.2 > b.1,  So order is: a c b
func (h TaskBullets) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *TaskBullets) Push(x interface{}) {
	*h = append(*h, x.(*TaskBullet))
}

func (h *TaskBullets) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
func (h *TaskBullets) IncreaseTerm() {
	for i := range *h {
		(*h)[i].IncreaseTerm()
	}
}
func (h *TaskBullets) DecreaseTerm() {
	for i := range *h {
		(*h)[i].DecreaseTerm()
	}
}

type UnSchedTaskWrap struct {
	Data          *Task
	PowerPartyIds []string `json:"powerPartyIds"`
}

func NewUnSchedTaskWrap(task *Task, powerPartyIds []string) *UnSchedTaskWrap {
	return &UnSchedTaskWrap{
		Data:          task,
		PowerPartyIds: powerPartyIds,
	}
}

type TaskMsg struct {
	TaskId        string `json:"taskId"`
	Data          *Task
	PowerPartyIds []string `json:"powerPartyIds"`
	// caches
	hash atomic.Value
}

func NewTaskMessageFromRequest(req *pb.PublishTaskDeclareRequest) *TaskMsg {

	return &TaskMsg{
		TaskId:        "",
		PowerPartyIds: req.PowerPartyIds,
		Data: NewTask(&libTypes.TaskData{
			TaskId:     "",
			TaskName:   req.TaskName,
			PartyId:    req.Owner.PartyId,
			Identity:   req.Owner.IdentityId,
			NodeId:     req.Owner.NodeId,
			NodeName:   req.Owner.Name,
			DataId:     "",
			DataStatus: DataStatusNormal.String(),
			State:      TaskStatePending.String(),
			Reason:     "",
			EventCount: 0,
			Desc:       "",
			CreateAt:   uint64(time.Now().UnixNano()),
			EndAt:      0,
			StartAt:    0,
			AlgoSupplier: &libTypes.OrganizationData{
				PartyId:  req.Owner.PartyId,
				Identity: req.Owner.IdentityId,
				NodeId:   req.Owner.NodeId,
				NodeName: req.Owner.Name,
			},
			TaskResource: &libTypes.TaskResourceData{
				CostProcessor: uint32(req.OperationCost.CostProcessor),
				CostMem:       req.OperationCost.CostMem,
				CostBandwidth: req.OperationCost.CostBandwidth,
				Duration:      req.OperationCost.Duration,
			},
		}),
	}
}
func ConvertTaskMsgToTaskWithPowers(task *Task, powers []*libTypes.TaskResourceSupplierData) *Task {
	task.SetResourceSupplierArr(powers)

	if len(powers) == 0 {
		return task
	}

	// 组装 选出来的, powerSuppliers 到 receivers 中
	privors := make([]*libTypes.OrganizationData, len(powers))
	for i, supplier := range powers {
		privors[i] = supplier.Organization
	}

	for i, _ := range task.TaskData().Receivers {
		receiver := task.TaskData().Receivers[i]
		receiver.Provider = privors
		task.TaskData().Receivers[i] = receiver
	}

	return task
}

//type taskdata struct {
//	TaskName              string                `json:"taskName"`
//	Owner                 *TaskSupplier         `json:"owner"`
//	Partners              []*TaskSupplier       `json:"partners"`
//	PowerPartyIds         []string              `json:"powerPartyIds"`
//	Receivers             []*TaskResultReceiver `json:"receivers"`
//	CalculateContractCode string                `json:"calculateContractCode"`
//	DataSplitContractCode string                `json:"dataSplitContractCode"`
//	ContractExtraParams   string                `json:"contractExtraParams"`
//	OperationCost         *TaskOperationCost    `json:"spend"`
//	CreateAt              uint64                `json:"createAt"`
//}

type TaskMsgs []*TaskMsg

func (msg *TaskMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *TaskMsg) Unmarshal(b []byte) error { return nil }
func (msg *TaskMsg) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *TaskMsg) MsgType() string { return MSG_TASK }
func (msg *TaskMsg) Owner() *libTypes.OrganizationData {
	return &libTypes.OrganizationData{
		PartyId:  msg.Data.data.PartyId,
		NodeName: msg.Data.data.NodeName,
		NodeId:   msg.Data.data.NodeId,
		Identity: msg.Data.data.Identity,
	}
}
func (msg *TaskMsg) OwnerName() string       { return msg.Data.data.NodeName }
func (msg *TaskMsg) OwnerNodeId() string     { return msg.Data.data.NodeId }
func (msg *TaskMsg) OwnerIdentityId() string { return msg.Data.data.Identity }
func (msg *TaskMsg) OwnerPartyId() string    { return msg.Data.data.PartyId }
func (msg *TaskMsg) TaskName() string        { return msg.Data.data.TaskName }

func (msg *TaskMsg) TaskMetadataSuppliers() []*libTypes.OrganizationData {
	partners := make([]*libTypes.OrganizationData, len(msg.Data.data.MetadataSupplier))
	for i, v := range msg.Data.data.MetadataSupplier {
		partners[i] = &libTypes.OrganizationData{
			PartyId:  v.Organization.PartyId,
			NodeName: v.Organization.NodeName,
			NodeId:   v.Organization.NodeId,
			Identity: v.Organization.Identity,
		}
	}
	return partners
}
func (msg *TaskMsg) TaskMetadataSupplierDatas() []*libTypes.TaskMetadataSupplierData {
	return msg.Data.data.MetadataSupplier
}

func (msg *TaskMsg) TaskResourceSuppliers() []*libTypes.OrganizationData {
	powers := make([]*libTypes.OrganizationData, len(msg.Data.data.ResourceSupplier))
	for i, v := range msg.Data.data.ResourceSupplier {
		powers[i] = &libTypes.OrganizationData{
			PartyId:  v.Organization.PartyId,
			NodeName: v.Organization.NodeName,
			NodeId:   v.Organization.NodeId,
			Identity: v.Organization.Identity,
		}
	}
	return powers
}
func (msg *TaskMsg) TaskResourceSupplierDatas() []*libTypes.TaskResourceSupplierData {
	return msg.Data.data.ResourceSupplier
}
func (msg *TaskMsg) GetPowerPartyIds() []string { return msg.PowerPartyIds }
func (msg *TaskMsg) GetReceivers() []*libTypes.OrganizationData {
	receivers := make([]*libTypes.OrganizationData, len(msg.Data.data.Receivers))
	for i, v := range msg.Data.data.Receivers {
		receivers[i] = &libTypes.OrganizationData{
			PartyId:  v.Receiver.PartyId,
			NodeName: v.Receiver.NodeName,
			NodeId:   v.Receiver.NodeId,
			Identity: v.Receiver.Identity,
		}
	}
	return receivers
}
func (msg *TaskMsg) TaskResultReceiverDatas() []*libTypes.TaskResultReceiverData {
	return msg.Data.data.Receivers
}
func (msg *TaskMsg) CalculateContractCode() string             { return msg.Data.data.CalculateContractCode }
func (msg *TaskMsg) DataSplitContractCode() string             { return msg.Data.data.DataSplitContractCode }
func (msg *TaskMsg) ContractExtraParams() string               { return msg.Data.data.ContractExtraParams }
func (msg *TaskMsg) OperationCost() *libTypes.TaskResourceData { return msg.Data.data.TaskResource }
func (msg *TaskMsg) CreateAt() uint64                          { return msg.Data.data.CreateAt }
func (msg *TaskMsg) SetTaskId() string {
	if "" != msg.TaskId {
		return msg.TaskId
	}
	msg.TaskId = PREFIX_TASK_ID + msg.Hash().Hex()
	return msg.TaskId
}
func (msg *TaskMsg) Hash() common.Hash {
	if hash := msg.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlputil.RlpHash(msg.Data)
	msg.hash.Store(v)
	return v
}

// Len returns the length of s.
func (s TaskMsgs) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s TaskMsgs) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s TaskMsgs) Less(i, j int) bool { return s[i].Data.data.CreateAt < s[j].Data.data.CreateAt }

type TaskSupplier struct {
	*TaskNodeAlias
	MetaData *SupplierMetaData `json:"metaData"`
}

type SupplierMetaData struct {
	MetaDataId      string   `json:"metaDataId"`
	ColumnIndexList []uint64 `json:"columnIndexList"`
}

type TaskResultReceiver struct {
	*TaskNodeAlias
	Providers []*TaskNodeAlias `json:"providers"`
}

type TaskOperationCost struct {
	Processor uint64 `json:"processor"`
	Mem       uint64 `json:"mem"`
	Bandwidth uint64 `json:"bandwidth"`
	Duration  uint64 `json:"duration"`
}

func ConvertTaskOperationCostToPB(cost *TaskOperationCost) *pb.TaskOperationCostDeclare {
	return &pb.TaskOperationCostDeclare{
		CostMem:       cost.Mem,
		CostProcessor: cost.Processor,
		CostBandwidth: cost.Bandwidth,
		Duration:      cost.Duration,
	}
}
func ConvertTaskOperationCostFromPB(cost *pb.TaskOperationCostDeclare) *TaskOperationCost {
	return &TaskOperationCost{
		Mem:       cost.CostMem,
		Processor: cost.CostProcessor,
		Bandwidth: cost.CostBandwidth,
		Duration:  cost.Duration,
	}
}

// ------------------- data using authorize -------------------

type DataAuthorizationApply struct {
	PoposalHash string     `json:"poposalHash"`
	Proposer    *NodeAlias `json:"proposer"`
	Approver    *NodeAlias `json:"approver"`
	Apply       struct {
		MetaId       string `json:"metaId"`
		UseCount     uint64 `json:"useCount"`
		UseStartTime string `json:"useStartTime"`
		UseEndTime   string `json:"useEndTime"`
	} `json:"apply"`
}

type DataAuthorizationConfirm struct {
	PoposalHash string     `json:"poposalHash"`
	Proposer    *NodeAlias `json:"proposer"`
	Approver    *NodeAlias `json:"approver"`
	Approve     struct {
		Vote   uint16 `json:"vote"`
		Reason string `json:"reason"`
	} `json:"approve"`
}

// ------------------- common -------------------

type NodeAlias struct {
	Name       string `json:"name"`
	NodeId     string `json:"nodeId"`
	IdentityId string `json:"identityId"`
}

type TaskNodeAlias struct {
	PartyId    string `json:"partyId"`
	Name       string `json:"name"`
	NodeId     string `json:"nodeId"`
	IdentityId string `json:"identityId"`
}

func ConvertNodeAliasToPB(alias *NodeAlias) *pb.OrganizationIdentityInfo {
	return &pb.OrganizationIdentityInfo{
		Name:       alias.Name,
		NodeId:     alias.NodeId,
		IdentityId: alias.IdentityId,
	}
}

func ConvertTaskNodeAliasToPB(alias *TaskNodeAlias) *pb.TaskOrganizationIdentityInfo {
	return &pb.TaskOrganizationIdentityInfo{
		PartyId:    alias.PartyId,
		Name:       alias.Name,
		NodeId:     alias.NodeId,
		IdentityId: alias.IdentityId,
	}
}

func ConvertNodeAliasFromPB(org *pb.OrganizationIdentityInfo) *NodeAlias {
	return &NodeAlias{
		Name:       org.Name,
		NodeId:     org.NodeId,
		IdentityId: org.IdentityId,
	}
}

func ConvertTaskNodeAliasFromPB(org *pb.TaskOrganizationIdentityInfo) *TaskNodeAlias {
	return &TaskNodeAlias{
		PartyId:    org.PartyId,
		Name:       org.Name,
		NodeId:     org.NodeId,
		IdentityId: org.IdentityId,
	}
}

func ConvertNodeAliasArrToPB(aliases []*NodeAlias) []*pb.OrganizationIdentityInfo {
	orgs := make([]*pb.OrganizationIdentityInfo, len(aliases))
	for i, a := range aliases {
		org := ConvertNodeAliasToPB(a)
		orgs[i] = org
	}
	return orgs
}

func ConvertTaskNodeAliasArrToPB(aliases []*TaskNodeAlias) []*pb.TaskOrganizationIdentityInfo {
	orgs := make([]*pb.TaskOrganizationIdentityInfo, len(aliases))
	for i, a := range aliases {
		org := ConvertTaskNodeAliasToPB(a)
		orgs[i] = org
	}
	return orgs
}

func ConvertNodeAliasArrFromPB(orgs []*pb.OrganizationIdentityInfo) []*NodeAlias {
	aliases := make([]*NodeAlias, len(orgs))
	for i, o := range orgs {
		alias := ConvertNodeAliasFromPB(o)
		aliases[i] = alias
	}
	return aliases
}

func ConvertTaskNodeAliasArrFromPB(orgs []*pb.TaskOrganizationIdentityInfo) []*TaskNodeAlias {
	aliases := make([]*TaskNodeAlias, len(orgs))
	for i, o := range orgs {
		alias := ConvertTaskNodeAliasFromPB(o)
		aliases[i] = alias
	}
	return aliases
}

func (n *NodeAlias) GetNodeName() string       { return n.Name }
func (n *NodeAlias) GetNodeIdStr() string      { return n.NodeId }
func (n *NodeAlias) GetNodeIdentityId() string { return n.IdentityId }

type ResourceUsage struct {
	TotalMem       uint64 `json:"totalMem"`
	UsedMem        uint64 `json:"usedMem"`
	TotalProcessor uint64 `json:"totalProcessor"`
	UsedProcessor  uint64 `json:"usedProcessor"`
	TotalBandwidth uint64 `json:"totalBandwidth"`
	UsedBandwidth  uint64 `json:"usedBandwidth"`
}

func ConvertResourceUsageToPB(usage *ResourceUsage) *pb.ResourceUsedDetailShow {
	return &pb.ResourceUsedDetailShow{
		TotalMem:       usage.TotalMem,
		UsedMem:        usage.UsedMem,
		TotalProcessor: usage.TotalProcessor,
		UsedProcessor:  usage.UsedProcessor,
		TotalBandwidth: usage.TotalBandwidth,
		UsedBandwidth:  usage.UsedBandwidth,
	}
}
func ConvertResourceUsageFromPB(usage *pb.ResourceUsedDetailShow) *ResourceUsage {
	return &ResourceUsage{
		TotalMem:       usage.TotalMem,
		UsedMem:        usage.UsedMem,
		TotalProcessor: usage.TotalProcessor,
		UsedProcessor:  usage.UsedProcessor,
		TotalBandwidth: usage.TotalBandwidth,
		UsedBandwidth:  usage.UsedBandwidth,
	}
}

/**
Example:
{
  "party_id": "p0",
  "data_party": {
      "input_file": "../data/bank_predict_data.csv",
      "id_column_name": "CLIENT_ID"
    },
  "dynamic_parameter": {
    "model_restore_party": "p0",
    "train_task_id": "task_id"
  }
}
*/
type FighterTaskReadyGoReqContractCfg struct {
	PartyId   string `json:"party_id"`
	DataParty struct {
		InputFile    string `json:"input_file"`
		IdColumnName string `json:"id_column_name"` // 目前 默认只会用一列, 后面再拓展 ..
	} `json:"data_party"`
	DynamicParameter map[string]interface{} `json:"dynamic_parameter"`
}
