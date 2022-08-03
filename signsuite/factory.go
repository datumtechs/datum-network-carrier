package signsuite

import (
	"fmt"
	"github.com/datumtechs/datum-network-carrier/common"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/ethereum/go-ethereum/crypto"
)

func Sender(userType commonconstantpb.UserType, hash common.Hash, sig []byte) (string, string, error) {

	switch userType {
	case commonconstantpb.UserType_User_1: // PlatON
		return signatureUserTypeIsUser1(hash, sig)
	case commonconstantpb.UserType_User_2: // Alaya
		return "", "", nil
	case commonconstantpb.UserType_User_3: // Ethereum
		return "", "", nil
	default:
		return "", "", fmt.Errorf("unknown userType")
	}
}

func signatureUserTypeIsUser1(hash common.Hash, signature []byte) (string, string, error) {
	if len(signature) != 65 {
		return common.Address{}.String(), "", fmt.Errorf("signature must be 65 bytes long")
	}
	if signature[64] != 27 && signature[64] != 28 {
		return common.Address{}.String(), "", fmt.Errorf("invalid signature (V is not 27 or 28)")
	}
	signature[64] -= 27 // Transform yellow paper V from 27/28 to 0/1

	rpk, err := crypto.SigToPub(hash.Bytes(), signature)
	if err != nil {
		return common.Address{}.String(), "", err
	}
	publicKey := string(crypto.FromECDSAPub(rpk))
	return crypto.PubkeyToAddress(*rpk).String(), publicKey, nil
}
