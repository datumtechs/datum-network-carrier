package signsuite

import (
	"fmt"
	"github.com/datumtechs/datum-network-carrier/common"
	libtypes "github.com/datumtechs/datum-network-carrier/lib/types"
	"github.com/datumtechs/datum-network-carrier/params"
)

func Sender (userType libtypes.UserType, hash common.Hash, sig []byte) (string, error) {
	switch userType {
	case libtypes.UserType_User_1: // PlatON
		signer := NewLatSigner(params.CarrierConfig().BlockChainIdCache[userType])
		return signer.Sender(hash, sig)
	case libtypes.UserType_User_2: // Alaya
		return "", nil
	case libtypes.UserType_User_3: // Ethereum
		return "", nil
	default:
		return "", fmt.Errorf("unknown userType")
	}
}
