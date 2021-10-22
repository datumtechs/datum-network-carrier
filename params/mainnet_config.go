package params

import (
	"time"
)

// MainnetConfig returns the configuration to be used in the main network.
func MainnetConfig() *CarrierChainConfig {
	return mainnetCarrierConfig
}

func UseMainnetConfig() {
	carrierConfig = MainnetConfig()
}

// the default config for main network.
var mainnetCarrierConfig = &CarrierChainConfig{

}

var mainnetNetworkConfig = &NetworkConfig{
	GossipMaxSize:          1 << 20, // 1 MiB
	MaxChunkSize:           1 << 20, // 1 MiB
	AttestationSubnetCount: 64,
	//AttestationPropagationSlotRange: 32,
	MaxRequestBlocks:            1 << 10, // 1024
	TtfbTimeout:                 5 * time.Second,
	RespTimeout:                 10 * time.Second,
	MaximumGossipClockDisparity: 500 * time.Millisecond,
	MessageDomainInvalidSnappy:  [4]byte{00, 00, 00, 00},
	MessageDomainValidSnappy:    [4]byte{01, 00, 00, 00},
	ETH2Key:                     "eth2",
	AttSubnetKey:                "attnets",
	MinimumPeersInSubnet:        4,
	MinimumPeersInSubnetSearch:  20,
	ContractDeploymentBlock:     11184524, // Note: contract was deployed in block 11052984 but no transactions were sent until 11184524.
	BootstrapNodes: []string{
		// Teku team's bootnode
		//"enr:-KG4QOtcP9X1FbIMOe17QNMKqDxCpm14jcX5tiOE4_TyMrFqbmhPZHK_ZPG2Gxb1GE2xdtodOfx9-cgvNtxnRyHEmC0ghGV0aDKQ9aX9QgAAAAD__________4JpZIJ2NIJpcIQDE8KdiXNlY3AyNTZrMaEDhpehBDbZjM_L9ek699Y7vhUJ-eAdMyQW_Fil522Y0fODdGNwgiMog3VkcIIjKA",
		//"enr:-KG4QDyytgmE4f7AnvW-ZaUOIi9i79qX4JwjRAiXBZCU65wOfBu-3Nb5I7b_Rmg3KCOcZM_C3y5pg7EBU5XGrcLTduQEhGV0aDKQ9aX9QgAAAAD__________4JpZIJ2NIJpcIQ2_DUbiXNlY3AyNTZrMaEDKnz_-ps3UUOfHWVYaskI5kWYO_vtYMGYCQRAR3gHDouDdGNwgiMog3VkcIIjKA",
		//"enr:-Jy4QAy3ph-KqpQNCGLFryppazc5ccMPc0ieFEozz0Uww1_WEt9susMjFyG_v55_My_0dV4w88aIeUNLeaR3Xu4HsoEBh2F0dG5ldHOIAAAAAAAAAACCaWSCdjSCaXCEfwAAAYlzZWNwMjU2azGhArg67qYwMcW9bKPA71p0KxurmyBX50lDhaVMqceS6HZzg3RjcIJBlYN1ZHCCgjU",
	},
}
