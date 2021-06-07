package types

// ------------------- task -------------------

type TaskMessage struct {
	TaskHash              string                `json:"taskHash"`
	Owner                 *TaskParticipant      `json:"owner"`
	Partners              []*TaskParticipant    `json:"partners"`
	Receivers             []*TaskResultReceiver `json:"receivers"`
	CalculateContractCode string                `json:"calculateContractCode"`
	DataSplitContractCode string                `json:"dataSplitContractCode"`
	Spend                 *TaskSpend            `json:"spend"`
	Extra                 string                `json:"extra"`
}

type TaskParticipant struct {
	NodeId     string               `json:"nodeId"`
	IdentityId string               `json:"identityId"`
	Alias      string               `json:"alias"`
	MetaData   *ParticipantMetaData `json:"metaData"`
}

type ParticipantMetaData struct {
	MetaId     string   `json:"metaId"`
	FilePath   string   `json:"filePath"`
	CindexList []uint32 `json:"cindexList"`
}

type TaskResultReceiver struct {
	NodeId     string       `json:"nodeId"`
	IdentityId string       `json:"identityId"`
	Providers  []*NodeAlias `json:"providers"`
}

type TaskSpend struct {
	Processor uint64 `json:"processor"`
	Mem       string `json:"mem"`
	Bandwidth string `json:"bandwidth"`
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
	NodeId     string `json:"nodeId"`
	IdentityId string `json:"identityId"`
	Alias      string `json:"alias"`
}

// ------------------- power -------------------

type PowerMsg struct {
	NodeId      string `json:"nodeId"`
	IdentityId  string `json:"identityId"`
	Information struct {
		InternalIp   string `json:"internalIp"`
		InternalPort string `json:"internalPort"`
		ExternalIp   string `json:"externalIp"`
		ExternalPort string `json:"externalPort"`
	} `json:"information"`
}


type DataNodeMsg struct {

}


// ------------------- metaData -------------------

type MetaDataMsg struct {

}