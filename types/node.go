package types

type YarnNodeInfo struct {
	NodeType       string                `json:"nodeType"`
	NodeId         string                `json:"nodeId"`
	InternalIp     string                `json:"internalIp"`
	ExternalIp     string                `json:"externalIp"`
	InternalPort   string                `json:"internalPort"`
	ExternalPort   string                `json:"externalPort"`
	IdentityType   string                `json:"identityType"`
	IdentityId     string                `json:"identityId"`
	Name           string                `json:"name"`
	TotalMem       uint64                `json:"totalMem"`
	UsedMem        uint64                `json:"usedMem"`
	TotalProcessor uint64                `json:"totalProcessor"`
	UsedProcessor  uint64                `json:"usedProcessor"`
	TotalBandwidth uint64                `json:"totalBandwidth"`
	UsedBandwidth  uint64                `json:"usedBandwidth"`
	Peers          *RegisteredNodeDetail `json:"peers"`
	SeedPeers      *SeedNodeInfo         `json:"seedPeers"`
	State          string                `json:"state"`
}

type YarnRegisteredNodeDetail struct {
	JobNodes  []*YarnRegisteredJobNode  `json:"jobNodes"`
	DataNodes []*YarnRegisteredDataNode `json:"DataNodes"`
}

type YarnRegisteredJobNode struct {
	Id            string         `json:"id"`
	InternalIp    string         `json:"internalIp"`
	ExternalIp    string         `json:"externalIp"`
	InternalPort  string         `json:"internalPort"`
	ExternalPort  string         `json:"externalPort"`
	ResourceUsage *ResourceUsage `json:"resourceUsage"`
	Duration      uint64         `json:"duration"`
	Task          struct {
		Count   uint32   `json:"count"`
		TaskIds []string `json:"taskIds"`
	} `json:"task"`
}

type YarnRegisteredDataNode struct {
	Id            string         `json:"id"`
	InternalIp    string         `json:"internalIp"`
	ExternalIp    string         `json:"externalIp"`
	InternalPort  string         `json:"internalPort"`
	ExternalPort  string         `json:"externalPort"`
	ResourceUsage *ResourceUsage `json:"resourceUsage"`
	Duration      uint64         `json:"duration"`
	Delta         struct {
		FileCount     uint32 `json:"fileCount"`
		FileTotalSize uint32 `json:"fileTotalSize"`
	} `json:"delta"`
}
