package twopc

import (
	"crypto/ecdsa"
	"github.com/Metisnetwork/Metis-Carrier/p2p"
)

type OptionConfig struct {
	NodePriKey *ecdsa.PrivateKey `json:"-"`
	NodeID     p2p.NodeID        `json:"nodeID"`
}

type Config struct {
	Option           *OptionConfig `json:"option"`
	PeerMsgQueueSize uint64        `json:"peerMsgQueueSize"`
	ConsensusStateFile string	   `json:"consensusStateFile"`
	DefaultConsensusWal string	   `json:"defaultConsensusWal"`
	DatabaseHandles int `json:"databaseHandles"`
	DatabaseCache   int `json:"databaseCache"`
}
