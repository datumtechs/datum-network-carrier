package types

// ------------------- SeedNode -------------------
// ------------------- JobNode -------------------
// ------------------- DataNode -------------------

// ------------------- identity -------------------
type IdentityMsg struct {
	*NodeAlias
	CreateAt uint64 `json:"createAt"`
}

// ------------------- power -------------------

type PowerMsg struct {
	PowerId string `json:"powerId"`
	*NodeAlias
	JobNodeId   string `json:"jobNodeId"`
	Information struct {
		Mem       string `json:"mem,omitempty"`
		Processor string `json:"processor,omitempty"`
		Bandwidth string `json:"bandwidth,omitempty"`
	} `json:"information"`
	CreateAt uint64 `json:"createAt"`
}

func (msg *PowerMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *PowerMsg) Unmarshal(b []byte) error { return nil }
func (msg *PowerMsg) String() string           { return "" }
func (msg *PowerMsg) MsgType() string          { return "" }

type PowerMsgs []*PowerMsg

// Len returns the length of s.
func (s PowerMsgs) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s PowerMsgs) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// ------------------- metaData -------------------

type MetaDataMsg struct {
	MetaDataId string `json:"metaDataId"`
	*NodeAlias
	Information struct {
		MetaSummary struct {
			MetaDataId string `json:"metaDataId,omitempty"`
			OriginId   string `json:"originId,omitempty"`
			TableName  string `json:"tableName,omitempty"`
			Desc       string `json:"desc,omitempty"`
			FilePath   string `json:"filePath,omitempty"`
			Rows       uint32 `json:"rows,omitempty"`
			Columns    uint32 `json:"columns,omitempty"`
			Size_      string `json:"size,omitempty"`
			FileType   string `json:"fileType,omitempty"`
			HasTitle   bool   `json:"hasTitle,omitempty"`
			State      string `json:"state,omitempty"`
		} `json:"metaSummary"`
		ColumnMeta []*struct {
			Cindex   uint64 `json:"cindex,omitempty"`
			Cname    string `json:"cname,omitempty"`
			Ctype    string `json:"ctype,omitempty"`
			Csize    uint64 `json:"csize,omitempty"`
			Ccomment string `json:"ccomment,omitempty"`
		} `json:"columnMeta"`
	} `json:"information"`
	CreateAt uint64 `json:"createAt"`
}

func (msg *MetaDataMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *MetaDataMsg) Unmarshal(b []byte) error { return nil }
func (msg *MetaDataMsg) String() string           { return "" }
func (msg *MetaDataMsg) MsgType() string          { return "" }

type MetaDataMsgs []*MetaDataMsg

// Len returns the length of s.
func (s MetaDataMsgs) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s MetaDataMsgs) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// ------------------- task -------------------

type TaskMsg struct {
	TaskId                string                `json:"taskId"`
	TaskName              string                `json:"taskName"`
	Owner                 *TaskSupplier         `json:"owner"`
	Partners              []*TaskSupplier       `json:"partners"`
	Receivers             []*TaskResultReceiver `json:"receivers"`
	CalculateContractCode string                `json:"calculateContractCode"`
	DataSplitContractCode string                `json:"dataSplitContractCode"`
	OperationCost         *TaskOperationCost    `json:"spend"`
	CreateAt              uint64                `json:"createAt"`
}

func (msg *TaskMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *TaskMsg) Unmarshal(b []byte) error { return nil }
func (msg *TaskMsg) String() string           { return "" }
func (msg *TaskMsg) MsgType() string          { return "" }

type TaskMsgs []*TaskMsg

// Len returns the length of s.
func (s TaskMsgs) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s TaskMsgs) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type TaskSupplier struct {
	*NodeAlias
	MetaData *SupplierMetaData `json:"metaData"`
}

type SupplierMetaData struct {
	MetaId          string   `json:"metaId"`
	ColumnIndexList []uint32 `json:"columnIndexList"`
}

type TaskResultReceiver struct {
	*NodeAlias
	Providers []*NodeAlias `json:"providers"`
}

type TaskOperationCost struct {
	Processor uint64 `json:"processor"`
	Mem       uint64 `json:"mem"`
	Bandwidth uint64 `json:"bandwidth"`
	Duration  uint64 `json:"duration"`
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
