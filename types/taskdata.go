package types

import (
	"bytes"
	"github.com/RosettaFlow/Carrier-Go/common"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"io"
	"sync/atomic"
)

type Task struct {
	data *libtypes.TaskPB

	// caches
	hash atomic.Value
	size atomic.Value
}

func NewTask(data *libtypes.TaskPB) *Task {
	return &Task{data: data}
}

func (m *Task) EncodePb(w io.Writer) error {
	if m.data == nil {
		m.data = new(libtypes.TaskPB)
	}
	data, err := m.data.Marshal()
	if err == nil {
		w.Write(data)
	}
	return err
}

func (m *Task) DecodePb(data []byte) error {
	m.size.Store(common.StorageSize(len(data)))
	return m.data.Unmarshal(data)
}

func (m *Task) Hash() common.Hash {
	if hash := m.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	buffer := new(bytes.Buffer)
	m.EncodePb(buffer)
	v := protoBufHash(buffer.Bytes())
	m.hash.Store(v)
	return v
}

func (m *Task) GetTaskId() string {
	return m.data.TaskId
}

func (m *Task) GetTaskData() *libtypes.TaskPB {
	return m.data
}

func (m *Task) GetTaskSender() *apicommonpb.TaskOrganization {
	return &apicommonpb.TaskOrganization {
		PartyId: m.data.GetPartyId(),
		NodeName: m.data.GetNodeName(),
		NodeId: m.data.GetNodeId(),
		IdentityId: m.data.GetIdentityId(),
	}
}


func (m *Task) SetEventList(eventList []*libtypes.TaskEvent) {
	/*eventArr := make([]*libtypes.TaskEvent, len(eventList))
	for i, ev := range eventList {
		eventArr[i] = &libtypes.TaskEvent{
			GetTaskId:     ev.GetTaskId,
			Type:       ev.Type,
			GetCreateAt:   ev.GetCreateAt,
			Content:    ev.Content,
			IdentityId: ev.IdentityId,
		}
	}*/
	m.data.TaskEvents = eventList
}
func (m *Task) SetMetadataSupplierArr(arr []*libtypes.TaskDataSupplier) {
	m.data.DataSuppliers = arr
}
func (m *Task) SetResourceSupplierArr(arr []*libtypes.TaskPowerSupplier) {
	m.data.PowerSuppliers = arr
}
func (m *Task) SetReceivers(arr []*apicommonpb.TaskOrganization) {
	m.data.Receivers = arr
}

// TaskDataArray is a Transaction slice type for basic sorting.
type TaskDataArray []*Task

// Len returns the length of s.
func (s TaskDataArray) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s TaskDataArray) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s TaskDataArray) GetPb(i int) []byte {
	buffer := new(bytes.Buffer)
	s[i].EncodePb(buffer)
	return buffer.Bytes()
}

func NewTaskDataArray(metaData []*libtypes.TaskPB) TaskDataArray {
	var s TaskDataArray
	for _, v := range metaData {
		s = append(s, NewTask(v))
	}
	return s
}

func (s TaskDataArray) To() []*libtypes.TaskPB {
	arr := make([]*libtypes.TaskPB, 0, s.Len())
	for _, v := range s {
		arr = append(arr, v.data)
	}
	return arr
}

//type TaskDetailShow struct {
//	GetTaskId        string                   `json:"taskId"`
//	TaskName      string                   `json:"taskName"`
//	Role          string                   `json:"role"`
//	Owner         *TaskNodeAlias           `json:"owner"`
//	AlgoSupplier  *TaskNodeAlias           `json:"algoSupplier"`
//	DataSupplier  []*TaskDataSupplierShow  `json:"dataSupplier"`
//	PowerSupplier []*TaskPowerSupplierShow `json:"powerSupplier"`
//	Receivers     []*TaskNodeAlias         `json:"receivers"`
//	GetCreateAt      uint64                   `json:"createat"`
//	StartAt       uint64                   `json:"startAt"`
//	EndAt         uint64                   `json:"endAt"`
//	GetState         string                   `json:"state"`
//	OperationCost *TaskOperationCost       `json:"operationCost"`
//}

//func ConvertTaskDetailShowToPB(task *TaskDetailShow) *pb.TaskDetailShow {
//	return &pb.TaskDetailShow {
//		GetTaskId:        task.GetTaskId,
//		TaskName:      task.TaskName,
//		Owner:         ConvertTaskNodeAliasToPB(task.Owner),
//		AlgoSupplier:  ConvertTaskNodeAliasToPB(task.AlgoSupplier),
//		DataSupplier:  ConvertTaskDataSupplierShowArrToPB(task.DataSupplier),
//		PowerSupplier: ConvertTaskPowerSupplierShowArrToPB(task.PowerSupplier),
//		//TODO: 结构问题，待确认
//		//Receivers:     ConvertTaskNodeAliasArrToPB(task.Receivers),
//		GetCreateAt:      task.GetCreateAt,
//		StartAt:       task.StartAt,
//		EndAt:         task.EndAt,
//		GetState:         task.GetState,
//		OperationCost: ConvertTaskOperationCostToPB(task.OperationCost),
//	}
//}
//func ConvertTaskDetailShowFromPB(task *pb.TaskDetailShow) *TaskDetailShow {
//	return &TaskDetailShow{
//		GetTaskId:        task.GetTaskId,
//		TaskName:      task.TaskName,
//		Owner:         ConvertTaskNodeAliasFromPB(task.Owner),
//		AlgoSupplier:  ConvertTaskNodeAliasFromPB(task.AlgoSupplier),
//		DataSupplier:  ConvertTaskDataSupplierShowArrFromPB(task.DataSupplier),
//		PowerSupplier: ConvertTaskPowerSupplierShowArrFromPB(task.PowerSupplier),
//		Receivers:     ConvertTaskNodeAliasArrFromPB(task.Receivers),
//		GetCreateAt:      task.GetCreateAt,
//		StartAt:       task.StartAt,
//		EndAt:         task.EndAt,
//		GetState:         task.GetState,
//		OperationCost: ConvertTaskOperationCostFromPB(task.OperationCost),
//	}
//}

//type TaskDataSupplierShow struct {
//	MemberInfo   *TaskNodeAlias `json:"memberInfo"`
//	MetadataId   string         `json:"metaId"`
//	MetadataName string         `json:"metaName"`
//}

//func ConvertTaskDataSupplierShowToPB(dataSupplier *TaskDataSupplierShow) *pb.TaskDataSupplierShow {
//	return &pb.TaskDataSupplierShow{
//		MemberInfo:   ConvertTaskNodeAliasToPB(dataSupplier.MemberInfo),
//		MetadataId:   dataSupplier.MetadataId,
//		MetadataName: dataSupplier.MetadataName,
//	}
//}
//func ConvertTaskDataSupplierShowFromPB(dataSupplier *pb.TaskDataSupplierShow) *TaskDataSupplierShow {
//	return &TaskDataSupplierShow{
//		MemberInfo:   ConvertTaskNodeAliasFromPB(dataSupplier.MemberInfo),
//		MetadataId:   dataSupplier.MetadataId,
//		MetadataName: dataSupplier.MetadataName,
//	}
//}
//func ConvertTaskDataSupplierShowArrToPB(dataSuppliers []*TaskDataSupplierShow) []*pb.TaskDataSupplierShow {
//	arr := make([]*pb.TaskDataSupplierShow, len(dataSuppliers))
//	for i, dataSupplier := range dataSuppliers {
//
//		supplier := &pb.TaskDataSupplierShow{
//			MemberInfo:   ConvertTaskNodeAliasToPB(dataSupplier.MemberInfo),
//			MetadataId:   dataSupplier.MetadataId,
//			MetadataName: dataSupplier.MetadataName,
//		}
//		arr[i] = supplier
//	}
//
//	return arr
//}
//func ConvertTaskDataSupplierShowArrFromPB(dataSuppliers []*pb.TaskDataSupplierShow) []*TaskDataSupplierShow {
//	arr := make([]*TaskDataSupplierShow, len(dataSuppliers))
//	for i, dataSupplier := range dataSuppliers {
//
//		supplier := &TaskDataSupplierShow{
//			MemberInfo:   ConvertTaskNodeAliasFromPB(dataSupplier.MemberInfo),
//			MetadataId:   dataSupplier.MetadataId,
//			MetadataName: dataSupplier.MetadataName,
//		}
//		arr[i] = supplier
//	}
//
//	return arr
//}

//type TaskPowerSupplierShow struct {
//	MemberInfo    *TaskNodeAlias `json:"memberInfo"`
//	ResourceUsage *ResourceUsage `json:"resourceUsage"`
//}

//func ConvertTaskPowerSupplierShowShowToPB(powerSupplier *TaskPowerSupplierShow) *pb.TaskPowerSupplierShow {
//	return &pb.TaskPowerSupplierShow{
//		MemberInfo: ConvertTaskNodeAliasToPB(powerSupplier.MemberInfo),
//		PowerInfo:  ConvertResourceUsageToPB(powerSupplier.ResourceUsage),
//	}
//}
//func ConvertTaskPowerSupplierShowFromPB(powerSupplier *pb.TaskPowerSupplierShow) *TaskPowerSupplierShow {
//	return &TaskPowerSupplierShow{
//		MemberInfo:    ConvertTaskNodeAliasFromPB(powerSupplier.MemberInfo),
//		ResourceUsage: ConvertResourceUsageFromPB(powerSupplier.PowerInfo),
//	}
//}
//func ConvertTaskPowerSupplierShowArrToPB(powerSuppliers []*TaskPowerSupplierShow) []*pb.TaskPowerSupplierShow {
//	arr := make([]*pb.TaskPowerSupplierShow, len(powerSuppliers))
//	for i, powerSupplier := range powerSuppliers {
//
//		supplier := &pb.TaskPowerSupplierShow{
//			MemberInfo: ConvertTaskNodeAliasToPB(powerSupplier.MemberInfo),
//			PowerInfo:  ConvertResourceUsageToPB(powerSupplier.ResourceUsage),
//		}
//		arr[i] = supplier
//	}
//
//	return arr
//}
//func ConvertTaskPowerSupplierShowArrFromPB(powerSuppliers []*pb.TaskPowerSupplierShow) []*TaskPowerSupplierShow {
//	arr := make([]*TaskPowerSupplierShow, len(powerSuppliers))
//	for i, powerSupplier := range powerSuppliers {
//
//		supplier := &TaskPowerSupplierShow{
//			MemberInfo:    ConvertTaskNodeAliasFromPB(powerSupplier.MemberInfo),
//			ResourceUsage: ConvertResourceUsageFromPB(powerSupplier.PowerInfo),
//		}
//		arr[i] = supplier
//	}
//
//	return arr
//}
