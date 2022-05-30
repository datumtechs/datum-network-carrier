package signsuite

import (
	"fmt"
	"github.com/datumtechs/datum-network-carrier/common"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"github.com/datumtechs/datum-network-carrier/params"
)

func Sender (userType carriertypespb.UserType, hash common.Hash, sig []byte) (string, error) {
	switch userType {
	case carriertypespb.UserType_User_1: // PlatON
		signer := NewLatSigner(params.CarrierConfig().BlockChainIdCache[userType])
		return signer.Sender(hash, sig)
	case carriertypespb.UserType_User_2: // Alaya
		return "", nil
	case carriertypespb.UserType_User_3: // Ethereum
		return "", nil
	default:
		return "", fmt.Errorf("unknown userType")
	}
}
