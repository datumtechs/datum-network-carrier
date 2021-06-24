package p2p

import (
	"github.com/RosettaFlow/Carrier-Go/params"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"time"
)

// ENR key used for eth2-related fork data.
var eth2ENRKey = params.CarrierNetworkConfig().ETH2Key

// ForkDigest returns the current fork digest of
// the node.
func (s *Service) forkDigest() ([4]byte, error) {
	//if s.currentForkDigest != [4]byte{} {
	//	return s.currentForkDigest, nil
	//}
	//fd, err := p2putils.CreateForkDigest(s.genesisTime, s.genesisValidatorsRoot)
	//if err != nil {
	//	s.currentForkDigest = fd
	//}
	//return fd, err
	return [4]byte{0x1, 0x1, 0x0, 0x0}, nil
}

// Adds a fork entry as an ENR record under the eth2EnrKey for
// the local node. The fork entry is an ssz-encoded enrForkID type
// which takes into account the current fork version from the current
// epoch to create a fork digest, the next fork version,
// and the next fork epoch.
func addForkEntry(
	node *enode.LocalNode,
	genesisTime time.Time,
	genesisValidatorsRoot []byte,
) (*enode.LocalNode, error) {
	//digest, err := p2putils.CreateForkDigest(genesisTime, genesisValidatorsRoot)
	//if err != nil {
	//	return nil, err
	//}
	//currentSlot := helpers.SlotsSince(genesisTime)
	//currentEpoch := helpers.SlotToEpoch(currentSlot)
	//if timeutils.Now().Before(genesisTime) {
	//	currentEpoch = 0
	//}
	//fork, err := p2putils.Fork(currentEpoch)
	//if err != nil {
	//	return nil, err
	//}
	//
	//nextForkEpoch := params.BeaconConfig().NextForkEpoch
	//nextForkVersion := params.BeaconConfig().NextForkVersion
	//// Set to the current fork version if our next fork is not planned.
	//if nextForkEpoch == math.MaxUint64 {
	//	nextForkVersion = fork.CurrentVersion
	//}
	//enrForkID := &pb.ENRForkID{
	//	CurrentForkDigest: digest[:],
	//	NextForkVersion:   nextForkVersion,
	//	NextForkEpoch:     nextForkEpoch,
	//}
	//enc, err := enrForkID.MarshalSSZ()
	//if err != nil {
	//	return nil, err
	//}
	//forkEntry := enr.WithEntry(eth2ENRKey, enc)
	//node.Set(forkEntry)
	return node, nil
}
