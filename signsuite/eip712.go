package signsuite // recoverEIP712 recovers the public key for eip712 signed data.
import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec"
	"github.com/datumtechs/datum-network-carrier/signsuite/eip712"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"
	"strings"
)

func recoverEIP712(signature []byte, data *eip712.TypedData) (string, string, error) {
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

	sighash, err := legacyKeccak256(rawData)
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

func legacyKeccak256(data []byte) ([]byte, error) {
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
