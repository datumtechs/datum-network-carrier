package types

// Resource usage of a single node host machine (Local storage)
type NodeResourceUsagePower struct {
	NodeId        string         `json:"nodeId"` // The node host machine id
	PowerId       string         `json:"powerId"`
	ResourceUsage *ResourceUsage `json:"resourceUsage"`
	State         string         `json:"state"`
}

type NodeResourceUsagePowerRes struct {
	Onwer *NodeAlias              `json:"owner"` // organization identity
	Usage *NodeResourceUsagePower `json:"usage"`
}

// Resource usage of all node host machines in a single organization (Decentralized storage)
type OrgResourceUsagePower struct {
	Onwer         *NodeAlias     `json:"owner"` // organization identity
	ResourceUsage *ResourceUsage `json:"resourceUsage"`
}

type OrgResourcePowerAndTaskCount struct {
	Power     *OrgResourceUsagePower `json:"power"`
	TaskCount uint32                 `json:"taskCount"`
}

type PowerTaskDetail struct {
	TotalTaskCount   uint32       `json:"total_task_count,omitempty"`
	CurrentTaskCount uint32       `json:"current_task_count,omitempty"`
	Tasks            []*PowerTask `json:"tasks,omitempty"`
}

type PowerTask struct {
	TaskId         string             `json:"task_id,omitempty"`
	Owner          *NodeAlias         `json:"owner,omitempty"`
	Patners        []*NodeAlias       `json:"patners,omitempty"`
	Receivers      []*NodeAlias       `json:"receivers,omitempty"`
	OperationCost  *TaskOperationCost `json:"operation_cost,omitempty"`
	OperationSpend *TaskOperationCost `json:"operation_spend,omitempty"`
}

type OrgPowerTaskDetail struct {
	Owner 		  *NodeAlias  				`json:"owner"`
	PowerDetail   *PowerTaskDetail			`json:"powerDetail"`
}

type OrgMetaDataSummary struct {
	Owner 				*NodeAlias 				`json:"owner"`
	MetaDataSummary 	*MetaDataSummary 		`json:"metaDataSummary"`
}

type OrgMetaDataInfo struct {
	Owner 				*NodeAlias 				`json:"owner"`
	MetaData 			*MetaDataInfo  			`json:"metaData"`
}

type MetaDataInfo struct {
	MetaDataSummary *MetaDataSummary `json:"metaDataSummary"`
	ColumnMetas     []*ColumnMeta    `json:"columnMetas"`
}