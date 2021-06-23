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
	*ScheduleNodeAlias
	MetaData *SupplierMetaData `json:"metaData"`
}
type ScheduleTaskPowerSupplier struct {
	*ScheduleNodeAlias
}

type ScheduleSupplierMetaData struct {
	MetaId          string   `json:"metaId"`
	ColumnIndexList []uint64 `json:"columnIndexList"`
}

type ScheduleTaskResultReceiver struct {
	*ScheduleNodeAlias
	Providers []*ScheduleNodeAlias `json:"providers"`
}

// Task consensus result
type ScheduleResult struct {
	TaskId 		string 			`json:"taskId"`
	Done 	    bool 			`json:"done"`
}