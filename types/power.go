package types

import (
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
)

//// Resource usage of a single node host machine (Local storage)
//type NodeResourceUsagePower struct {
//	NodeId        string         `json:"nodeId"` // The node host machine id
//	PowerId       string         `json:"powerId"`
//	ResourceUsage *ResourceUsage `json:"resourceUsage"`
//	State         string         `json:"state"`
//}
//
//type NodeResourceUsagePowerRes struct {
//	Onwer *NodeAlias              `json:"owner"` // organization identity
//	Usage *NodeResourceUsagePower `json:"usage"`
//}
//
//// Resource usage of all node host machines in a single organization (Decentralized storage)
//type OrgResourceUsagePower struct {
//	Onwer         *NodeAlias     `json:"owner"` // organization identity
//	ResourceUsage *ResourceUsage `json:"resourceUsage"`
//}
//
//type OrgResourcePowerAndTaskCount struct {
//	Power     *OrgResourceUsagePower `json:"power"`
//	TaskCount uint32                 `json:"taskCount"`
//}

type PowerTotalDetail struct {
	TotalTaskCount   uint32         `json:"totalTaskCount"`
	CurrentTaskCount uint32         `json:"currentTaskCount"`
	Tasks            []*PowerTask   `json:"tasks"`
	ResourceUsage    *ResourceUsage `json:"resourceUsage"`
	State            string         `json:"state"`
}

type PowerSingleDetail struct {
	JobNodeId        string         `json:"jobNodeId"`
	PowerId          string         `json:"powerId"`
	TotalTaskCount   uint32         `json:"totalTaskCount"`
	CurrentTaskCount uint32         `json:"currentTaskCount"`
	Tasks            []*PowerTask   `json:"tasks"`
	ResourceUsage    *ResourceUsage `json:"resourceUsage"`
	State            string         `json:"state"`
}

type PowerTask struct {
	TaskId         string             `json:"taskId"`
	TaskName       string             `json:"taskName"`
	Owner          *NodeAlias         `json:"owner"`
	Patners        []*NodeAlias       `json:"patners"`
	Receivers      []*NodeAlias       `json:"receivers"`
	OperationCost  *TaskOperationCost `json:"operationCost"`
	OperationSpend *TaskOperationCost `json:"operationSpend"`
	CreateAt       uint64             `json:"createAt"`
}

//func ConvertPowerTaskToPB(task *PowerTask) *libTypes.PowerTask {
//	return &libTypes.PowerTask{
//		TaskId:         task.TaskId,
//		TaskName:       task.TaskName,
//		Owner:          ConvertNodeAliasToPB(task.Owner),
//		Patners:        ConvertNodeAliasArrToPB(task.Patners),
//		Receivers:      ConvertNodeAliasArrToPB(task.Receivers),
//		OperationCost:  ConvertTaskOperationCostToPB(task.OperationCost),
//		OperationSpend: ConvertTaskOperationCostToPB(task.OperationSpend),
//		CreateAt:       task.CreateAt,
//	}
//}
//func ConvertPowerTaskFromPB(task *libTypes.PowerTask) *PowerTask {
//	return &PowerTask{
//		TaskId:         task.TaskId,
//		TaskName:       task.TaskName,
//		Owner:          ConvertNodeAliasFromPB(task.Owner),
//		Patners:        ConvertNodeAliasArrFromPB(task.Patners),
//		Receivers:      ConvertNodeAliasArrFromPB(task.Receivers),
//		OperationCost:  ConvertTaskOperationCostFromPB(task.OperationCost),
//		OperationSpend: ConvertTaskOperationCostFromPB(task.OperationSpend),
//		CreateAt:       task.CreateAt,
//	}
//}

func ConvertPowerTaskArrToPB(tasks []*PowerTask) []*libTypes.PowerTask {

	arr := make([]*libTypes.PowerTask, len(tasks))
	for i, task := range tasks {
		t := &libTypes.PowerTask{
			TaskId:         task.TaskId,
			TaskName:       task.TaskName,
			Owner:          ConvertNodeAliasToPB(task.Owner),
			Patners:        ConvertNodeAliasArrToPB(task.Patners),
			Receivers:      ConvertNodeAliasArrToPB(task.Receivers),
			OperationCost:  ConvertTaskOperationCostToPB(task.OperationCost),
			OperationSpend: ConvertTaskOperationCostToPB(task.OperationSpend),
			CreateAt:       task.CreateAt,
		}
		arr[i] = t
	}
	return arr
}
//func ConvertPowerTaskArrFromPB(tasks []*libTypes.PowerTask) []*PowerTask {
//	arr := make([]*PowerTask, len(tasks))
//	for i, task := range tasks {
//		t := &PowerTask{
//			TaskId:         task.TaskId,
//			TaskName:       task.TaskName,
//			Owner:          ConvertNodeAliasFromPB(task.Owner),
//			Patners:        ConvertNodeAliasArrFromPB(task.Patners),
//			Receivers:      ConvertNodeAliasArrFromPB(task.Receivers),
//			OperationCost:  ConvertTaskOperationCostFromPB(task.OperationCost),
//			OperationSpend: ConvertTaskOperationCostFromPB(task.OperationSpend),
//			CreateAt:       task.CreateAt,
//		}
//		arr[i] = t
//	}
//	return arr
//}

type OrgPowerDetail struct {
	Owner       *NodeAlias        `json:"owner"`
	PowerDetail *PowerTotalDetail `json:"powerDetail"`
}

type NodePowerDetail struct {
	Owner       *NodeAlias         `json:"owner"`
	PowerDetail *PowerSingleDetail `json:"powerDetail"`
}
