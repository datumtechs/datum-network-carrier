package types

type RegisterNodeConnectStatus int32

const (
	CONNECTED    RegisterNodeConnectStatus = 0
	NONCONNECTED RegisterNodeConnectStatus = -1
	ENABLECOMPUTERESOURCE RegisterNodeConnectStatus = 1

)

type SeedNodeInfo struct {
	Id           string                    `json:"id"`
	InternalIp   string                    `json:"internalIp"`
	InternalPort string                    `json:"internalPort"`
	ConnState    RegisterNodeConnectStatus `json:"connState"`
}

type RegisteredNodeInfo struct {
	Id           string                    `json:"id"`
	InternalIp   string                    `json:"internalIp"`
	InternalPort string                    `json:"internalPort"`
	ExternalIp   string                    `json:"externalIp"`
	ExternalPort string                    `json:"externalPort"`
	ConnState    RegisterNodeConnectStatus `json:"connState"`
}
