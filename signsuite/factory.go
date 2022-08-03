package signsuite

import (
	"fmt"
	"github.com/datumtechs/datum-network-carrier/common"
	"github.com/datumtechs/datum-network-carrier/params"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/datumtechs/datum-network-carrier/types"
	"math/big"
)

func Sender(userType commonconstantpb.UserType, hash common.Hash, sig []byte) (string, error) {

	switch userType {
	case commonconstantpb.UserType_User_1: // PlatON
		fmt.Println(1233)
		return SignFunc(params.CarrierConfig().BlockChainIdCache[userType])(hash, sig)
	case commonconstantpb.UserType_User_2: // Alaya
		return "", nil
	case commonconstantpb.UserType_User_3: // Ethereum
		return "", nil
	default:
		return "", fmt.Errorf("unknown userType")
	}
}

func SignFunc(chainId *big.Int) types.LatSingFunc {
	signer := NewLatSigner(chainId)
	return signer.Sender
}
