package p2p

import statefeed "github.com/RosettaFlow/Carrier-Go/common/feed/state"

// Config for the p2p service. These parameters are set from application
// level flags to initialize the p2p service.
type Config struct {
	NoDiscovery         bool
	EnableUPnP          bool
	DisableDiscv5       bool
	StaticPeers         []string
	BootstrapNodeAddr   []string
	Discv5BootStrapAddr []string
	RelayNodeAddr       string
	LocalIP             string
	HostAddress         string
	HostDNS             string
	PrivateKey          string
	DataDir             string
	MetaDataDir         string
	TCPPort             uint
	UDPPort             uint
	MaxPeers            uint
	AllowListCIDR       string
	DenyListCIDR        []string
	//TODO: need to update..
	StateNotifier       statefeed.Notifier
	//DB                  db.ReadOnlyDatabase
}
