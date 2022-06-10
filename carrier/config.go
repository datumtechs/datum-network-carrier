package carrier

import (
	"github.com/datumtechs/datum-network-carrier/carrierdb"
	"github.com/datumtechs/datum-network-carrier/p2p"
	"github.com/datumtechs/datum-network-carrier/params"
)

// DefaultConfig contains default settings for use on the Carrier main.
var DefaultConfig = Config{
	DatabaseCache: 768,
	DiscoverServiceConfig: &params.DiscoverServiceConfig{
		DiscoveryServerIP: "",
		DiscoveryServerPort: 0,
		DiscoveryServiceId: "",
		DiscoveryServiceName: "carrierService",
		DiscoveryServerTags: []string{"carrier"},
		DiscoveryServiceHealthCheckInterval: 3000,
		DiscoveryServiceHealthCheckDeregister: 10000,
	},
	TaskManagerConfig: &params.TaskManagerConfig{
		MetadataConsumeOption:          0,
		NeedReplayScheduleTaskChanSize: 600,
		NeedExecuteTaskChanSize:        600,
		TaskConsResultChanSize:         600,
	},
}

//go:generate gencodec -type Config -formats toml -out gen_config.go

type Config struct {
	CarrierDB carrierdb.CarrierDB
	P2P       p2p.P2P
	// Database options
	DatabaseHandles       int
	DatabaseCache         int
	DefaultConsensusWal   string
	DiscoverServiceConfig *params.DiscoverServiceConfig
	TaskManagerConfig     *params.TaskManagerConfig
	//Token20PayConfig    // TODO 将 token  添加到这里 不可以么
}
