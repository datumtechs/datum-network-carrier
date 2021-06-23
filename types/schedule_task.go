package types

type ScheduleTask struct {
	TaskId   			  string 							`json:"TaskId"`
	TaskName              string                			`json:"taskName"`
	Owner                 *ScheduleTaskDataSupplier 		`json:"owner"`
	Partners              []*ScheduleTaskDataSupplier		`json:"partners"`
	PowerSuppliers 		  []*ScheduleTaskPowerSupplier 		`json:"powerSuppliers"`
	Receivers             []*ScheduleTaskResultReceiver 	`json:"receivers"`
	CalculateContractCode string                			`json:"calculateContractCode"`
	DataSplitContractCode string                			`json:"dataSplitContractCode"`
	OperationCost         *TaskOperationCost    			`json:"spend"`
}

type ScheduleTaskDataSupplier struct {
	*NodeAlias
	MetaData *SupplierMetaData `json:"metaData"`
}
type ScheduleTaskPowerSupplier struct {
	*NodeAlias
}

type ScheduleSupplierMetaData struct {
	MetaId          string   `json:"metaId"`
	ColumnIndexList []uint64 `json:"columnIndexList"`
}

type ScheduleTaskResultReceiver struct {
	*NodeAlias
	Providers []*NodeAlias `json:"providers"`
}

// Task consensus result
type ScheduleResult struct {
	TaskId 		string 			`json:"taskId"`
	Done 	    bool 			`json:"done"`
}

//type ScheduleNodeAlias struct {
//	Name         string `json:"name"`
//	NodeId       string `json:"nodeId"`
//	IdentityId   string `json:"identityId"`
//	PeerIp   string `json:"peerIp"`
//	PeerPort string `json:"peerPort"`
//}