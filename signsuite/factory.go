package signsuite

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/datumtechs/datum-network-carrier/common"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/datumtechs/datum-network-carrier/signsuite/eip712"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
	"strings"
)

var log = logrus.WithField("prefix", "signSuite")

func Sender(userType commonconstantpb.UserType, hash common.Hash, sig []byte) (string, string, error) {

	switch userType {
	case commonconstantpb.UserType_User_1: // PlatON
		return RecoverEIP712(sig, eip712TypeDataForSign(hash.Hex()))
	case commonconstantpb.UserType_User_2: // Alaya
		return "", "", nil
	case commonconstantpb.UserType_User_3: // Ethereum
		return "", "", nil
	default:
		return "", "", fmt.Errorf("unknown userType")
	}
}

// RecoverEIP712 recovers the public key for eip712 signed data.
func RecoverEIP712(signature []byte, data *eip712.TypedData) (string, string, error) {
	if len(signature) != 65 {
		return "", "", errors.New("invalid length")
	}
	// Convert to btcec input format with 'recovery id' v at the beginning.
	btcsig := make([]byte, 65)
	btcsig[0] = signature[64]
	copy(btcsig[1:], signature)

	rawData, err := eip712.EncodeForSigning(data)
	if err != nil {
		return "", "", err
	}

	sighash, err := LegacyKeccak256(rawData)
	if err != nil {
		return "", "", err
	}

	p, _, err := btcec.RecoverCompact(btcec.S256(), btcsig, sighash)
	if err != nil {
		return "", "", err
	}
	pk := (*ecdsa.PublicKey)(p)
	address := crypto.PubkeyToAddress(*pk).String()
	publicKeyHexString := hex.EncodeToString(crypto.FromECDSAPub(pk))
	return strings.ToLower(address), publicKeyHexString, err
}

func LegacyKeccak256(data []byte) ([]byte, error) {
	var err error
	hasher := sha3.NewLegacyKeccak256()
	_, err = hasher.Write(data)
	if err != nil {
		return nil, err
	}
	return hasher.Sum(nil), err
}

func eip712TypeDataForSign(contents string) *eip712.TypedData {
	log.Debugf("eip712TypeDataForSign contents detail is:%s", contents)
	return &eip712.TypedData{
		Domain: eip712.TypedDataDomain{
			Name: "Datum",
		},
		Message: eip712.TypedDataMessage{
			"contents": contents,
		},
		PrimaryType: "sign",
		Types: eip712.Types{
			"EIP712Domain": {
				{
					Name: "name",
					Type: "string",
				},
			},
			"sign": {
				{
					Name: "contents",
					Type: "string",
				},
			},
		},
	}
}
