package types

//// Resource usage of a single node host machine (Local storage)
//type NodeResourceUsagePower struct {
//	NodeId        string         `json:"nodeId"` // The node host machine id
//	PowerId       string         `json:"powerId"`
//	ResourceUsage *ResourceUsage `json:"resourceUsage"`
//	GetState         string         `json:"state"`
//}
//
//type NodeResourceUsagePowerRes struct {
//	Owner *NodeAlias              `json:"owner"` // organization identity
//	Usage *NodeResourceUsagePower `json:"usage"`
//}
//
//// Resource usage of all node host machines in a single organization (Decentralized storage)
//type OrgResourceUsagePower struct {
//	Owner         *NodeAlias     `json:"owner"` // organization identity
//	ResourceUsage *ResourceUsage `json:"resourceUsage"`
//}
//
//type OrgResourcePowerAndTaskCount struct {
//	Power     *OrgResourceUsagePower `json:"power"`
//	TaskCount uint32                 `json:"taskCount"`
//}

//type PowerTotalDetail struct {
//	TotalTaskCount   uint32         `json:"totalTaskCount"`
//	CurrentTaskCount uint32         `json:"currentTaskCount"`
//	Tasks            []*PowerTask   `json:"tasks"`
//	ResourceUsage    *ResourceUsage `json:"resourceUsage"`
//	GetState            string         `json:"state"`
//}

//type PowerSingleDetail struct {
//	JobNodeId        string         `json:"jobNodeId"`
//	PowerId          string         `json:"powerId"`
//	TotalTaskCount   uint32         `json:"totalTaskCount"`
//	CurrentTaskCount uint32         `json:"currentTaskCount"`
//	Tasks            []*PowerTask   `json:"tasks"`
//	ResourceUsage    *ResourceUsage `json:"resourceUsage"`
//	GetState            string         `json:"state"`
//}

//type PowerTask struct {
//	GetTaskId         string             `json:"taskId"`
//	TaskName       string             `json:"taskName"`
//	Owner          *NodeAlias         `json:"owner"`
//	Patners        []*NodeAlias       `json:"patners"`
//	Receivers      []*NodeAlias       `json:"receivers"`
//	OperationCost  *TaskOperationCost `json:"operationCost"`
//	OperationSpend *TaskOperationCost `json:"operationSpend"`
//	GetCreateAt       uint64             `json:"createAt"`
//}

//func ConvertPowerTaskToPB(task *PowerTask) *libTypes.PowerTask {
//	return &libTypes.PowerTask{
//		GetTaskId:         task.GetTaskId,
//		TaskName:       task.TaskName,
//		Owner:          ConvertNodeAliasToPB(task.Owner),
//		Patners:        ConvertNodeAliasArrToPB(task.Patners),
//		Receivers:      ConvertNodeAliasArrToPB(task.Receivers),
//		OperationCost:  ConvertTaskOperationCostToPB(task.OperationCost),
//		OperationSpend: ConvertTaskOperationCostToPB(task.OperationSpend),
//		GetCreateAt:       task.GetCreateAt,
//	}
//}
//func ConvertPowerTaskFromPB(task *libTypes.PowerTask) *PowerTask {
//	return &PowerTask{
//		GetTaskId:         task.GetTaskId,
//		TaskName:       task.TaskName,
//		Owner:          ConvertNodeAliasFromPB(task.Owner),
//		Patners:        ConvertNodeAliasArrFromPB(task.Patners),
//		Receivers:      ConvertNodeAliasArrFromPB(task.Receivers),
//		OperationCost:  ConvertTaskOperationCostFromPB(task.OperationCost),
//		OperationSpend: ConvertTaskOperationCostFromPB(task.OperationSpend),
//		GetCreateAt:       task.GetCreateAt,
//	}
//}

//func ConvertPowerTaskArrToPB(tasks []*PowerTask) []*libTypes.PowerTask {
//
//	arr := make([]*libTypes.PowerTask, len(tasks))
//	for i, task := range tasks {
//		t := &libTypes.PowerTask{
//			GetTaskId:         task.GetTaskId,
//			TaskName:       task.TaskName,
//			Owner:          ConvertNodeAliasToPB(task.Owner),
//			Patners:        ConvertNodeAliasArrToPB(task.Patners),
//			Receivers:      ConvertNodeAliasArrToPB(task.Receivers),
//			OperationCost:  ConvertTaskOperationCostToPB(task.OperationCost),
//			OperationSpend: ConvertTaskOperationCostToPB(task.OperationSpend),
//			GetCreateAt:       task.GetCreateAt,
//		}
//		arr[i] = t
//	}
//	return arr
//}
//func ConvertPowerTaskArrFromPB(tasks []*libTypes.PowerTask) []*PowerTask {
//	arr := make([]*PowerTask, len(tasks))
//	for i, task := range tasks {
//		t := &PowerTask{
//			GetTaskId:         task.GetTaskId,
//			TaskName:       task.TaskName,
//			Owner:          ConvertNodeAliasFromPB(task.Owner),
//			Patners:        ConvertNodeAliasArrFromPB(task.Patners),
//			Receivers:      ConvertNodeAliasArrFromPB(task.Receivers),
//			OperationCost:  ConvertTaskOperationCostFromPB(task.OperationCost),
//			OperationSpend: ConvertTaskOperationCostFromPB(task.OperationSpend),
//			GetCreateAt:       task.GetCreateAt,
//		}
//		arr[i] = t
//	}
//	return arr
//}

//type OrgPowerDetail struct {
//	Owner       *NodeAlias        `json:"owner"`
//	PowerDetail *PowerTotalDetail `json:"powerDetail"`
//}

//type NodePowerDetail struct {
//	Owner       *NodeAlias         `json:"owner"`
//	PowerDetail *PowerSingleDetail `json:"powerDetail"`
//}
