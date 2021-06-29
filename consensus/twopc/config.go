package twopc

import (
	"crypto/ecdsa"
	"github.com/RosettaFlow/Carrier-Go/p2p"
)

type OptionConfig struct {
	NodePriKey *ecdsa.PrivateKey `json:"-"`
	NodeID     p2p.NodeID        `json:"nodeID"`
}

type Config struct {
	Option *OptionConfig `json:"option"`
}
