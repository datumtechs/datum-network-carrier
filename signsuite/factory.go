package signsuite

import (
	"fmt"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/params"
)

func Sender (userType libtypes.UserType, sig []byte) (string, error) {
	switch userType {
	case libtypes.UserType_User_1: // PlatON
		signer := NewLatSigner(params.CarrierConfig().BlockChainIdCache[userType])
		return signer.Sender(sig)
	case libtypes.UserType_User_2: // Alaya
		return "", nil
	case libtypes.UserType_User_3: // Ethereum
		return "", nil
	default:
		return "", fmt.Errorf("unknown userType")
	}
}
