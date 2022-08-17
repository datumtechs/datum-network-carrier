package signsuite

import (
	"fmt"
	"github.com/datumtechs/datum-network-carrier/common"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("prefix", "signSuite")

func Sender(userType commonconstantpb.UserType, hash common.Hash, sig []byte) (string, string, error) {

	switch userType {
	case commonconstantpb.UserType_User_1: // PlatON
		return recoverEIP712(sig, eip712TypeDataForSign(hash.Hex()))
	case commonconstantpb.UserType_User_2: // Alaya
		return "", "", nil
	case commonconstantpb.UserType_User_3: // Ethereum
		return "", "", nil
	default:
		return "", "", fmt.Errorf("unknown userType")
	}
}
