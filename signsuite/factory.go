package signsuite

import (
	"fmt"
	"github.com/datumtechs/datum-network-carrier/common"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/datumtechs/datum-network-carrier/params"
)

func Sender (userType commonconstantpb.UserType, hash common.Hash, sig []byte) (string, error) {
	switch userType {
	case commonconstantpb.UserType_User_1: // PlatON
		signer := NewLatSigner(params.CarrierConfig().BlockChainIdCache[userType])
		return signer.Sender(hash, sig)
	case commonconstantpb.UserType_User_2: // Alaya
		return "", nil
	case commonconstantpb.UserType_User_3: // Ethereum
		return "", nil
	default:
		return "", fmt.Errorf("unknown userType")
	}
}
