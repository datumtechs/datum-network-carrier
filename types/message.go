package types

import (
	"encoding/json"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/rlputil"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/lib/types"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"strings"
	"sync/atomic"
)

const (
	PREFIX_POWER_ID    = "power:"
	PREFIX_METADATA_ID = "metadata:"
	PREFIX_TASK_ID     = "task:"

	MSG_IDENTITY        = "identityMsg"
	MSG_IDENTITY_REVOKE = "identityRevokeMsg"
	MSG_POWER           = "powerMsg"
	MSG_POWER_REVOKE    = "powerRevokeMsg"
	MSG_METADATA        = "metaDataMsg"
	MSG_METADATA_REVOKE = "metaDataRevokeMsg"
	MSG_TASK            = "taskMsg"
)

type MessageType string

type Msg interface {
	Marshal() ([]byte, error)
	Unmarshal(b []byte) error
	String() string
	MsgType() string
}



// ------------------- identity -------------------
type IdentityMsg struct {
	*apipb.Organization
	CreateAt uint64 `json:"createAt"`
}

type IdentityRevokeMsg struct {
	CreateAt uint64 `json:"createAt"`
}

type IdentityMsgArr []*IdentityMsg
type IdentityRevokeMsgArr []*IdentityRevokeMsg

func (msg *IdentityMsg) ToDataCenter() *Identity {
	return NewIdentity(&libTypes.IdentityPB{
		NodeName:   msg.NodeName,
		NodeId:     msg.NodeId,
		IdentityId: msg.IdentityId,
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
func (msg *IdentityMsg) OwnerName() string              { return msg.NodeName }
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
		CreateAt:  uint64(timeutils.UnixMsec()),
	}
	msg.SetPowerId()
	return msg
}

type PowerRevokeMsg struct {
	*apipb.Organization
	PowerId  string `json:"powerId"`
	CreateAt uint64 `json:"createAt"`
}

func NewPowerRevokeMessageFromRequest(req *pb.RevokePowerRequest) *PowerRevokeMsg {
	return &PowerRevokeMsg{
		PowerId:  req.PowerId,
		CreateAt: uint64(timeutils.UnixMsec()),
	}
}

type PowerMsgArr []*PowerMsg
type PowerRevokeMsgArr []*PowerRevokeMsg

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

func (msg *PowerMsg) GetJobNodeId() string { return msg.JobNodeId }
func (msg *PowerMsg) GetCreateAt() uint64  { return msg.CreateAt }
func (msg *PowerMsg) SetPowerId() string {
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

	return rlputil.RlpHash([]interface{}{
		msg.JobNodeId,
		//msg.CreateAt,
		uint64(timeutils.UnixMsec()),
	})
}

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
func (s PowerMsgArr) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s PowerMsgArr) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s PowerMsgArr) Less(i, j int) bool { return s[i].CreateAt < s[j].CreateAt }

// ------------------- metaData -------------------

type MetadataMsg struct {
	MetadataId string `json:"metaDataId"`
	Data       *metadataData
	// caches
	hash atomic.Value
}

func NewMetadataMessageFromRequest(req *pb.PublishMetadataRequest) *MetadataMsg {
	return &MetadataMsg{
		Data: &metadataData{
			Information: struct {
				MetadataSummary *libTypes.MetadataSummary  `json:"metadataSummary"`
				ColumnMetas     []*libTypes.MetadataColumn `json:"columnMetas"`
			}{
				MetadataSummary: &libTypes.MetadataSummary{
					MetadataId: req.Information.MetadataSummary.MetadataId,
					OriginId:   req.Information.MetadataSummary.OriginId,
					TableName:  req.Information.MetadataSummary.TableName,
					Desc:       req.Information.MetadataSummary.Desc,
					FilePath:   req.Information.MetadataSummary.FilePath,
					Rows:       req.Information.MetadataSummary.Rows,
					Columns:    req.Information.MetadataSummary.Columns,
					Size_:      req.Information.MetadataSummary.Size_,
					FileType:   req.Information.MetadataSummary.FileType,
					HasTitle:   req.Information.MetadataSummary.HasTitle,
					State:      req.Information.MetadataSummary.State,
				},
				ColumnMetas: make([]*types.MetadataColumn, 0),
			},
			CreateAt: uint64(timeutils.UnixMsec()),
		},
	}
}

type metadataData struct {
	*apipb.Organization
	Information struct {
		MetadataSummary *libTypes.MetadataSummary `json:"metadataSummary"`
		ColumnMetas     []*types.MetadataColumn   `json:"columnMetas"`
	} `json:"information"`
	CreateAt uint64 `json:"createAt"`
}

type MetadataRevokeMsg struct {
	*apipb.Organization
	MetadataId string `json:"metadataId"`
	CreateAt   uint64 `json:"createAt"`
}

func NewMetadataRevokeMessageFromRequest(req *pb.RevokeMetadataRequest) *MetadataRevokeMsg {
	return &MetadataRevokeMsg{
		MetadataId: req.MetadataId,
		CreateAt:   uint64(timeutils.UnixMsec()),
	}
}

type MetadataMsgArr []*MetadataMsg
type MetadataRevokeMsgArr []*MetadataRevokeMsg

func (msg *MetadataMsg) ToDataCenter() *Metadata {
	return NewMetadata(&libTypes.MetadataPB{
		IdentityId:      msg.OwnerIdentityId(),
		NodeId:          msg.OwnerNodeId(),
		NodeName:        msg.OwnerName(),
		DataId:          msg.MetadataId,
		OriginId:        msg.OriginId(),
		TableName:       msg.TableName(),
		FilePath:        msg.FilePath(),
		FileType:        msg.FileType(),
		Desc:            msg.Desc(),
		Rows:            msg.Rows(),
		Columns:         msg.Columns(),
		Size_:           uint64(msg.Size()),
		HasTitle:        msg.HasTitle(),
		MetadataColumns: msg.ColumnMetas(),
		// the status of data, N means normal, D means deleted.
		DataStatus: apipb.DataStatus_DataStatus_Normal,
		// metaData status, eg: create/release/revoke
		State: apipb.MetadataState_MetadataState_Released,
	})
}
func (msg *MetadataMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *MetadataMsg) Unmarshal(b []byte) error { return nil }
func (msg *MetadataMsg) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *MetadataMsg) MsgType() string { return MSG_METADATA }
func (msg *MetadataMsg) Owner() *apipb.Organization {
	return &apipb.Organization{
		NodeName:   msg.Data.NodeName,
		NodeId:     msg.Data.NodeId,
		IdentityId: msg.Data.IdentityId,
	}
}
func (msg *MetadataMsg) OwnerName() string       { return msg.Data.NodeName }
func (msg *MetadataMsg) OwnerNodeId() string     { return msg.Data.NodeId }
func (msg *MetadataMsg) OwnerIdentityId() string { return msg.Data.IdentityId }
func (msg *MetadataMsg) MetaDataSummary() *libTypes.MetadataSummary {
	return msg.Data.Information.MetadataSummary
}
func (msg *MetadataMsg) OriginId() string  { return msg.Data.Information.MetadataSummary.OriginId }
func (msg *MetadataMsg) TableName() string { return msg.Data.Information.MetadataSummary.TableName }
func (msg *MetadataMsg) Desc() string      { return msg.Data.Information.MetadataSummary.Desc }
func (msg *MetadataMsg) FilePath() string  { return msg.Data.Information.MetadataSummary.FilePath }
func (msg *MetadataMsg) Rows() uint32      { return msg.Data.Information.MetadataSummary.Rows }
func (msg *MetadataMsg) Columns() uint32   { return msg.Data.Information.MetadataSummary.Columns }
func (msg *MetadataMsg) Size() uint64      { return msg.Data.Information.MetadataSummary.Size_ }
func (msg *MetadataMsg) FileType() apipb.OriginFileType {
	return msg.Data.Information.MetadataSummary.FileType
}
func (msg *MetadataMsg) HasTitle() bool                       { return msg.Data.Information.MetadataSummary.HasTitle }
func (msg *MetadataMsg) State() apipb.MetadataState           { return msg.Data.Information.MetadataSummary.State }
func (msg *MetadataMsg) ColumnMetas() []*types.MetadataColumn { return msg.Data.Information.ColumnMetas }
func (msg *MetadataMsg) CreateAt() uint64                     { return msg.Data.CreateAt }
func (msg *MetadataMsg) SetMetadataId() string {
	if "" != msg.MetadataId {
		return msg.MetadataId
	}
	msg.MetadataId = PREFIX_METADATA_ID + msg.HashByCreateTime().Hex()
	return msg.MetadataId
}
func (msg *MetadataMsg) Hash() common.Hash {
	if hash := msg.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlputil.RlpHash(msg.Data)
	msg.hash.Store(v)
	return v
}

func (msg *MetadataMsg) HashByCreateTime() common.Hash {
	return rlputil.RlpHash([]interface{}{
		msg.Data.Information.MetadataSummary.OriginId,
		uint64(timeutils.UnixMsec()),
	})
}

func (msg *MetadataRevokeMsg) ToDataCenter() *Metadata {
	return NewMetadata(&libTypes.MetadataPB{
		IdentityId: msg.IdentityId,
		NodeId:     msg.NodeId,
		NodeName:   msg.NodeName,
		DataId:     msg.MetadataId,
		// the status of data, N means normal, D means deleted.
		DataStatus: apipb.DataStatus_DataStatus_Deleted,
		// metaData status, eg: create/release/revoke
		State: apipb.MetadataState_MetadataState_Revoked,
	})
}
func (msg *MetadataRevokeMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *MetadataRevokeMsg) Unmarshal(b []byte) error { return nil }
func (msg *MetadataRevokeMsg) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *MetadataRevokeMsg) MsgType() string { return MSG_METADATA_REVOKE }

// Len returns the length of s.
func (s MetadataMsgArr) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s MetadataMsgArr) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s MetadataMsgArr) Less(i, j int) bool { return s[i].Data.CreateAt < s[j].Data.CreateAt }

// ------------------- task -------------------

type TaskBullet struct {
	TaskId  string
	Starve  bool
	Term    uint32
	Resched uint32
}

func NewTaskBullet(taskId string) *TaskBullet {
	return &TaskBullet{
		TaskId: taskId,
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

type TaskMsg struct {
	Data          *Task
	PowerPartyIds []string `json:"powerPartyIds"`
	// caches
	hash atomic.Value
}

func NewTaskMessageFromRequest(req *pb.PublishTaskDeclareRequest) *TaskMsg {

	return &TaskMsg{
		PowerPartyIds: req.PowerPartyIds,
		Data: NewTask(&libTypes.TaskPB{

			TaskId:     "",
			TaskName:   req.TaskName,
			PartyId:    req.Sender.PartyId,
			IdentityId: req.Sender.IdentityId,
			NodeId:     req.Sender.NodeId,
			NodeName:   req.Sender.NodeName,
			DataId:     "",
			DataStatus: apipb.DataStatus_DataStatus_Normal,
			State:      apipb.TaskState_TaskState_Pending,
			Reason:     "",
			EventCount: 0,
			Desc:       "",
			CreateAt:   uint64(timeutils.UnixMsec()),
			EndAt:      0,
			StartAt:    0,
			//TODO: 缺失，算法提供者信息
			/*AlgoSupplier: &apipb.TaskOrganization{
				PartyId:  req.Owner.PartyId,
				Identity: req.Owner.IdentityId,
				NodeId:   req.Owner.NodeId,
				NodeName: req.Owner.Name,
			},*/
			OperationCost:         req.OperationCost,
			CalculateContractCode: req.CalculateContractCode,
			DataSplitContractCode: req.DataSplitContractCode,
			ContractExtraParams:   req.ContractExtraParams,
		}),
	}
}
func ConvertTaskMsgToTaskWithPowers(task *Task, powers []*libTypes.TaskPowerSupplier) *Task {
	task.SetResourceSupplierArr(powers)
	return task
}

type TaskMsgArr []*TaskMsg

func (msg *TaskMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *TaskMsg) Unmarshal(b []byte) error { return nil }
func (msg *TaskMsg) String() string {
	return fmt.Sprintf(`{"taskId": %s, "powerPartyIds": %s, "task": %s}`,
		msg.Data.GetTaskId(), "["+strings.Join(msg.PowerPartyIds, ",")+"]", msg.Data.GetTaskData().String())
}
func (msg *TaskMsg) MsgType() string { return MSG_TASK }
func (msg *TaskMsg) Owner() *apipb.TaskOrganization {
	return &apipb.TaskOrganization{
		PartyId:    msg.Data.GetTaskData().GetPartyId(),
		NodeName:   msg.Data.GetTaskData().GetNodeName(),
		NodeId:     msg.Data.GetTaskData().GetNodeId(),
		IdentityId: msg.Data.GetTaskData().GetIdentityId(),
	}
}
func (msg *TaskMsg) OwnerName() string       { return msg.Data.GetTaskData().GetNodeName() }
func (msg *TaskMsg) OwnerNodeId() string     { return msg.Data.GetTaskData().GetNodeId() }
func (msg *TaskMsg) OwnerIdentityId() string { return msg.Data.GetTaskData().GetIdentityId() }
func (msg *TaskMsg) OwnerPartyId() string    { return msg.Data.GetTaskData().GetPartyId() }
func (msg *TaskMsg) TaskId() string          { return msg.Data.GetTaskData().GetTaskId() }
func (msg *TaskMsg) TaskName() string        { return msg.Data.GetTaskData().GetTaskName() }

func (msg *TaskMsg) TaskMetadataSuppliers() []*apipb.TaskOrganization {

	partners := make([]*apipb.TaskOrganization, len(msg.Data.GetTaskData().GetDataSuppliers()))
	for i, v := range msg.Data.GetTaskData().GetDataSuppliers() {

		partners[i] = &apipb.TaskOrganization{
			PartyId:    v.GetOrganization().GetPartyId(),
			NodeName:   v.GetOrganization().GetNodeName(),
			NodeId:     v.GetOrganization().GetNodeId(),
			IdentityId: v.GetOrganization().GetIdentityId(),
		}
	}
	return partners
}
func (msg *TaskMsg) TaskMetadataSupplierDatas() []*libTypes.TaskDataSupplier {

	return msg.Data.GetTaskData().GetDataSuppliers()
}

func (msg *TaskMsg) TaskResourceSuppliers() []*apipb.TaskOrganization {
	powers := make([]*apipb.TaskOrganization, len(msg.Data.GetTaskData().GetPowerSuppliers()))
	for i, v := range msg.Data.GetTaskData().GetPowerSuppliers() {
		powers[i] = &apipb.TaskOrganization{
			PartyId:    v.GetOrganization().GetPartyId(),
			NodeName:   v.GetOrganization().GetNodeName(),
			NodeId:     v.GetOrganization().GetNodeId(),
			IdentityId: v.GetOrganization().GetIdentityId(),
		}
	}
	return powers
}
func (msg *TaskMsg) TaskResourceSupplierDatas() []*libTypes.TaskPowerSupplier {
	return msg.Data.GetTaskData().GetPowerSuppliers()
}
func (msg *TaskMsg) GetPowerPartyIds() []string { return msg.PowerPartyIds }
func (msg *TaskMsg) GetReceivers() []*apipb.TaskOrganization {
	return msg.Data.GetTaskData().GetReceivers()
}

func (msg *TaskMsg) CalculateContractCode() string {
	return msg.Data.GetTaskData().GetCalculateContractCode()
}
func (msg *TaskMsg) DataSplitContractCode() string {
	return msg.Data.GetTaskData().GetDataSplitContractCode()
}
func (msg *TaskMsg) ContractExtraParams() string { return msg.Data.GetTaskData().GetContractExtraParams() }
func (msg *TaskMsg) OperationCost() *apipb.TaskResourceCostDeclare {
	return msg.Data.GetTaskData().GetOperationCost()
}
func (msg *TaskMsg) CreateAt() uint64 { return msg.Data.GetTaskData().GetCreateAt() }
func (msg *TaskMsg) SetTaskId() string {
	if "" != msg.Data.GetTaskId() {
		return msg.Data.GetTaskId()
	}
	msg.Data.GetTaskData().TaskId = PREFIX_TASK_ID + msg.HashByCreateTime().Hex()
	return msg.Data.GetTaskId()
}
func (msg *TaskMsg) Hash() common.Hash {
	if hash := msg.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlputil.RlpHash(msg.Data.GetTaskData())
	msg.hash.Store(v)
	return v
}

func (msg *TaskMsg) HashByCreateTime() common.Hash {
	return rlputil.RlpHash([]interface{}{
		msg.Data.GetTaskData().GetIdentityId(),
		msg.Data.GetTaskData().GetPartyId(),
		msg.Data.GetTaskData().GetTaskName(),
		//msg.Data.GetTaskData().CreateAt,
		uint64(timeutils.UnixMsec()),
	})
}

// Len returns the length of s.
func (s TaskMsgArr) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s TaskMsgArr) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s TaskMsgArr) Less(i, j int) bool {
	return s[i].Data.GetTaskData().GetCreateAt() < s[j].Data.GetTaskData().GetCreateAt()
}

//type TaskSupplier struct {
//	*TaskNodeAlias
//	MetadataPB *SupplierMetaData `json:"metaData"`
//}

//type SupplierMetaData struct {
//	MetaDataId      string   `json:"metaDataId"`
//	ColumnIndexList []uint64 `json:"columnIndexList"`
//}

//type TaskResultReceiver struct {
//	*TaskNodeAlias
//	Providers []*TaskNodeAlias `json:"providers"`
//}

//func ConvertTaskOperationCostToPB(cost *TaskOperationCost) *apipb.TaskResourceCostDeclare {
//	return &apipb.TaskResourceCostDeclare{
//		CostMem:       cost.Mem,
//		CostProcessor: uint32(cost.Processor),
//		CostBandwidth: cost.Bandwidth,
//		Duration:      cost.Duration,
//	}
//}
//func ConvertTaskOperationCostFromPB(cost *apipb.TaskResourceCostDeclare) *TaskOperationCost {
//	return &TaskOperationCost{
//		Mem:       cost.CostMem,
//		Processor: uint64(cost.CostProcessor),
//		Bandwidth: cost.CostBandwidth,
//		Duration:  cost.Duration,
//	}
//}

// ------------------- data using authorize -------------------

//type DataAuthorizationApply struct {
//	PoposalHash string     `json:"poposalHash"`
//	Proposer    *NodeAlias `json:"proposer"`
//	Approver    *NodeAlias `json:"approver"`
//	Apply       struct {
//		MetaId       string `json:"metaId"`
//		UseCount     uint64 `json:"useCount"`
//		UseStartTime string `json:"useStartTime"`
//		UseEndTime   string `json:"useEndTime"`
//	} `json:"apply"`
//}

//type DataAuthorizationConfirm struct {
//	PoposalHash string     `json:"poposalHash"`
//	Proposer    *NodeAlias `json:"proposer"`
//	Approver    *NodeAlias `json:"approver"`
//	Approve     struct {
//		Vote   uint16 `json:"vote"`
//		Reason string `json:"reason"`
//	} `json:"approve"`
//}

// ------------------- common -------------------

//type NodeAlias struct {
//	Name       string `json:"name"`
//	NodeId     string `json:"nodeId"`
//	IdentityId string `json:"identityId"`
//}

//type TaskNodeAlias struct {
//	PartyId    string `json:"partyId"`
//	Name       string `json:"name"`
//	NodeId     string `json:"nodeId"`
//	IdentityId string `json:"identityId"`
//}

//func (tna *TaskNodeAlias) String() string {
//	return fmt.Sprintf(`{"partyId": %s, "name": %s, "nodeId": %s, "identityId": %s}`, tna.PartyId, tna.Name, tna.NodeId, tna.IdentityId)
//}

//func ConvertNodeAliasToPB(alias *NodeAlias) *apipb.Organization {
//	return &apipb.Organization{
//		NodeName:   alias.Name,
//		NodeId:     alias.NodeId,
//		IdentityId: alias.IdentityId,
//	}
//}

//func ConvertTaskNodeAliasToPB(alias *TaskNodeAlias) *apipb.TaskOrganization {
//	return &apipb.TaskOrganization{
//		PartyId:    alias.PartyId,
//		NodeName:   alias.Name,
//		NodeId:     alias.NodeId,
//		IdentityId: alias.IdentityId,
//	}
//}

//func ConvertNodeAliasFromPB(org *apipb.Organization) *NodeAlias {
//	return &NodeAlias{
//		Name:       org.NodeName,
//		NodeId:     org.NodeId,
//		IdentityId: org.IdentityId,
//	}
//}

//func ConvertTaskNodeAliasFromPB(org *apipb.TaskOrganization) *TaskNodeAlias {
//	return &TaskNodeAlias{
//		PartyId:    org.PartyId,
//		Name:       org.NodeName,
//		NodeId:     org.NodeId,
//		IdentityId: org.IdentityId,
//	}
//}

//func ConvertNodeAliasArrToPB(aliases []*NodeAlias) []*apipb.Organization {
//	orgs := make([]*apipb.Organization, len(aliases))
//	for i, a := range aliases {
//		org := ConvertNodeAliasToPB(a)
//		orgs[i] = org
//	}
//	return orgs
//}

//func ConvertTaskNodeAliasArrToPB(aliases []*TaskNodeAlias) []*apipb.TaskOrganization {
//	orgs := make([]*apipb.TaskOrganization, len(aliases))
//	for i, a := range aliases {
//		org := ConvertTaskNodeAliasToPB(a)
//		orgs[i] = org
//	}
//	return orgs
//}

//func ConvertNodeAliasArrFromPB(orgs []*apipb.Organization) []*NodeAlias {
//	aliases := make([]*NodeAlias, len(orgs))
//	for i, o := range orgs {
//		alias := ConvertNodeAliasFromPB(o)
//		aliases[i] = alias
//	}
//	return aliases
//}

//func ConvertTaskNodeAliasArrFromPB(orgs []*apipb.TaskOrganization) []*TaskNodeAlias {
//	aliases := make([]*TaskNodeAlias, len(orgs))
//	for i, o := range orgs {
//		alias := ConvertTaskNodeAliasFromPB(o)
//		aliases[i] = alias
//	}
//	return aliases
//}

//func (n *NodeAlias) GetNodeName() string       { return n.Name }
//func (n *NodeAlias) GetNodeIdStr() string      { return n.NodeId }
//func (n *NodeAlias) GetNodeIdentityId() string { return n.IdentityId }

//type ResourceUsage struct {
//	TotalMem       uint64 `json:"totalMem"`
//	UsedMem        uint64 `json:"usedMem"`
//	TotalProcessor uint64 `json:"totalProcessor"`
//	UsedProcessor  uint64 `json:"usedProcessor"`
//	TotalBandwidth uint64 `json:"totalBandwidth"`
//	UsedBandwidth  uint64 `json:"usedBandwidth"`
//}

//func ConvertResourceUsageToPB(usage *ResourceUsage) *libTypes.ResourceUsageOverview {
//	return &libTypes.ResourceUsageOverview{
//		TotalMem:       usage.TotalMem,
//		UsedMem:        usage.UsedMem,
//		TotalProcessor: uint32(usage.TotalProcessor),
//		UsedProcessor:  uint32(usage.UsedProcessor),
//		TotalBandwidth: usage.TotalBandwidth,
//		UsedBandwidth:  usage.UsedBandwidth,
//	}
//}
//func ConvertResourceUsageFromPB(usage *libTypes.ResourceUsageOverview) *ResourceUsage {
//	return &ResourceUsage{
//		TotalMem:       usage.TotalMem,
//		UsedMem:        usage.UsedMem,
//		TotalProcessor: uint64(usage.TotalProcessor),
//		UsedProcessor:  uint64(usage.UsedProcessor),
//		TotalBandwidth: usage.TotalBandwidth,
//		UsedBandwidth:  usage.UsedBandwidth,
//	}
//}

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
