package addressutil

import (
	"fmt"
	"github.com/Metisnetwork/Metis-Carrier/common/bech32util"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/ethereum/go-ethereum/common"

	log "github.com/sirupsen/logrus"
)

//func HexToAddress
const (
	AddressLength = 20
	LatAddressHRP = "lat"
	AtpAddressHRP = "atp"
)

func IsBech32Address(s string) bool {
	hrp, _, err := bech32.Decode(s)
	if err != nil {
		return false
	}
	if !CheckAddressHRP(hrp) {
		return false
	}
	return true
}

func IsHexAddress(s string) bool {
	if hasHexPrefix(s) {
		s = s[2:]
	}
	return len(s) == 2*AddressLength && isHex(s)
}

// isHex validates whether each byte is valid hexadecimal string.
func isHex(str string) bool {
	if len(str)%2 != 0 {
		return false
	}
	for _, c := range []byte(str) {
		if !isHexCharacter(c) {
			return false
		}
	}
	return true
}

// isHexCharacter returns bool of c being a valid hexadecimal.
func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

// hasHexPrefix validates str begins with '0x' or '0X'.
func hasHexPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

func CheckAddressHRP(hrp string) bool {
	if LatAddressHRP == hrp || AtpAddressHRP == hrp {
		return true
	}
	return false
}

func MustBech32ToAddress(s string) common.Address {
	add, err := Bech32ToAddress(s)
	if err != nil {
		log.Error("must Bech32ToAddress fail", "err", err)
		panic(err)
	}
	return add
}
func Bech32ToAddress(s string) (common.Address, error) {
	hrpDecode, converted, err := bech32util.DecodeAndConvert(s)
	if err != nil {
		return common.Address{}, err
	}
	//support lat/atp
	if !CheckAddressHRP(hrpDecode) {
		return common.Address{}, fmt.Errorf("the address hrp not compare right,input:%s", s)
	}

	var a common.Address
	a.SetBytes(converted)
	return a, nil
}

func ConvertEthAddress(s string) common.Address {
	if IsHexAddress(s) {
		return common.HexToAddress(s)
	} else if IsBech32Address(s) {
		return MustBech32ToAddress(s)
	} else {
		log.Errorf("not supported address: %s", s)
		panic("not supported address")
	}
}
