package types

const (
	IDENTITY_TYPE_CA  = "CA"
	IDENTITY_TYPE_DID = "DID"
)

type NetworkMsgLocationSymbol bool

func (nmls NetworkMsgLocationSymbol) String() string {
	switch nmls {
	case LocalNetworkMsg:
		return "local network msg"
	case RemoteNetworkMsg:
		return "remote network msg"
	default:
		return "Unknown network location symbol"

	}
}

const (
	LocalNetworkMsg  NetworkMsgLocationSymbol = true
	RemoteNetworkMsg NetworkMsgLocationSymbol = false
)