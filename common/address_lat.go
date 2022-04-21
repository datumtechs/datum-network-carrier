package common

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/btcsuite/btcutil/bech32"

	"github.com/Metisnetwork/Metis-Carrier/common/bech32util"
	"github.com/Metisnetwork/Metis-Carrier/common/hexutil"
	"github.com/Metisnetwork/Metis-Carrier/crypto/sha3"
)

const (
	DefaultPlatONAddressHRP = "lat"
)


var currentPlatONAddressHRP string

func GetAddressHRP() string {
	if currentPlatONAddressHRP == "" {
		return DefaultPlatONAddressHRP
	}
	return currentPlatONAddressHRP
}

func SetAddressHRP(s string) error {
	if s == "" {
		s = DefaultPlatONAddressHRP
	}
	if len(s) != 3 {
		return errors.New("the length of addressHRP must be 3")
	}
	log.Info("the address hrp has been set", "hrp", s)
	currentPlatONAddressHRP = s
	return nil
}

func CheckAddressHRP(s string) bool {
	if currentPlatONAddressHRP != "" && s != currentPlatONAddressHRP {
		return false
	}
	return true
}

/////////// Address

// PlatONAddress represents the 20 byte address of an Ethereum account.
type PlatONAddress [AddressLength]byte

// BytesToAddress returns PlatONAddress with value b.
// If b is larger than len(h), b will be cropped from the left.
func BytesToPlatONAddress(b []byte) PlatONAddress {
	var a PlatONAddress
	a.SetBytes(b)
	return a
}

// BigToAddress returns Address with byte values of b.
// If b is larger than len(h), b will be cropped from the left.
func BigToPlatONAddress(b *big.Int) PlatONAddress { return BytesToPlatONAddress(b.Bytes()) }

// Deprecated: address to string is use bech32 now
// HexToAddress returns Address with byte values of s.
// If s is larger than len(h), s will be cropped from the left.
func HexToPlatONAddress(s string) PlatONAddress { return BytesToPlatONAddress(FromHex(s)) }

// MustBech32ToAddress returns Address with byte values of s.
// If s is Decode fail, it will return zero address.
func MustBech32ToPlatONAddress(s string) PlatONAddress {
	add, err := Bech32ToPlatONAddress(s)
	if err != nil {
		log.Error("must Bech32ToPlatONAddress fail", "err", err)
		panic(err)
	}
	return add
}

// MustBech32ToAddress returns Address with byte values of s.
// If s is Decode fail, it will return zero address.
func Bech32ToPlatONAddress(s string) (PlatONAddress, error) {
	hrpDecode, converted, err := bech32util.DecodeAndConvert(s)
	if err != nil {
		return PlatONAddress{}, err
	}
	if !CheckAddressHRP(hrpDecode) {
		return PlatONAddress{}, fmt.Errorf("the lat address hrp not compare right,input:%s", s)
	}

	if currentPlatONAddressHRP == "" {
		log.Warn("the lat address hrp not set yet", "input", s)
	} else if currentPlatONAddressHRP != hrpDecode {
		log.Warn("the lat address not compare current net", "want", currentPlatONAddressHRP, "input", s)
	}
	var a PlatONAddress
	a.SetBytes(converted)
	return a, nil
}

// Bech32ToPlatONAddressWithoutCheckHrp returns Address with byte values of s.
// If s is Decode fail, it will return zero address.
func Bech32ToPlatONAddressWithoutCheckHrp(s string) PlatONAddress {
	_, converted, err := bech32util.DecodeAndConvert(s)
	if err != nil {
		log.Error(" hrp lat address decode  fail", "err", err)
		panic(err)
	}
	var a PlatONAddress
	a.SetBytes(converted)
	return a
}

// Deprecated: address to string is use bech32 now
// IsHexAddress verifies whether a string can represent a valid hex-encoded
// Ethereum address or not.
func IsHexPlatONAddress(s string) bool {
	if has0xPrefix(s) {
		s = s[2:]
	}
	return len(s) == 2*AddressLength && isHex(s)
}

func IsBech32PlatONAddress(s string) bool {
	hrp, _, err := bech32.Decode(s)
	if err != nil {
		return false
	}
	if !CheckAddressHRP(hrp) {
		return false
	}
	return true
}

// Bytes gets the string representation of the underlying address.
func (a PlatONAddress) Bytes() []byte { return a[:] }

// Big converts an address to a big integer.
func (a PlatONAddress) Big() *big.Int { return new(big.Int).SetBytes(a[:]) }

// Hash converts an address to a hash by left-padding it with zeros.
func (a PlatONAddress) Hash() Hash { return BytesToHash(a[:]) }

// Deprecated: address to string is use bech32 now
// Hex returns an EIP55-compliant hex string representation of the address.it's use for node address
func (a PlatONAddress) Hex() string {
	return "0x" + a.HexWithNoPrefix()
}

// Deprecated: address to string is use bech32 now
func (a PlatONAddress) HexWithNoPrefix() string {
	unchecksummed := hex.EncodeToString(a[:])
	sha := sha3.NewKeccak256()
	sha.Write([]byte(unchecksummed))
	hash := sha.Sum(nil)

	result := []byte(unchecksummed)
	for i := 0; i < len(result); i++ {
		hashByte := hash[i/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if result[i] > '9' && hashByte > 7 {
			result[i] -= 32
		}
	}
	return string(result)
}

// String implements fmt.Stringer.
func (a PlatONAddress) String() string {
	return a.Bech32()
}

func (a PlatONAddress) Bech32() string {
	return a.Bech32WithHRP(GetAddressHRP())
}

func (a PlatONAddress) Bech32WithHRP(hrp string) string {
	if v, err := bech32util.ConvertAndEncode(hrp, a.Bytes()); err != nil {
		log.Error("lat address can't ConvertAndEncode to string", "err", err, "add", a.Bytes())
		return ""
	} else {
		return v
	}
}

// Format implements fmt.Formatter, forcing the byte slice to be formatted as is,
// without going through the stringer interface used for logging.
func (a PlatONAddress) Format(s fmt.State, c rune) {
	switch string(c) {
	case "s":
		fmt.Fprintf(s, "%"+string(c), a.String())
	default:
		fmt.Fprintf(s, "%"+string(c), a[:])
	}
}

// SetBytes sets the address to the value of b.
// If b is larger than len(a) it will panic.
func (a *PlatONAddress) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// MarshalText returns the hex representation of a.
func (a PlatONAddress) MarshalText() ([]byte, error) {
	v, err := bech32util.ConvertAndEncode(GetAddressHRP(), a.Bytes())
	if err != nil {
		return nil, err
	}
	return []byte(v), nil
}

// UnmarshalText parses a hash in hex syntax.
func (a *PlatONAddress) UnmarshalText(input []byte) error {
	hrpDecode, converted, err := bech32util.DecodeAndConvert(string(input))
	if err != nil {
		return err
	}
	if !CheckAddressHRP(hrpDecode) {
		return fmt.Errorf("the address not compare current net,want %v,have %v", GetAddressHRP(), string(input))
	}
	a.SetBytes(converted)
	return nil
}

// UnmarshalJSON parses a hash in hex syntax.
func (a *PlatONAddress) UnmarshalJSON(input []byte) error {
	if !isLatString(input) {
		return &json.UnmarshalTypeError{Value: "non-string", Type: addressT}
	}
	hrpDecode, v, err := bech32util.DecodeAndConvert(string(input[1 : len(input)-1]))
	if err != nil {
		return &json.UnmarshalTypeError{Value: err.Error(), Type: addressT}
	}
	if !CheckAddressHRP(hrpDecode) {
		return &json.UnmarshalTypeError{Value: fmt.Sprintf("hrpDecode not compare the current net,want %v,have %v", GetAddressHRP(), hrpDecode), Type: addressT}
	}
	a.SetBytes(v)
	return nil
}

func isLatString(input []byte) bool {
	return len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"'
}

// Scan implements Scanner for database/sql.
func (a *PlatONAddress) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into Address", src)
	}
	if len(srcB) != AddressLength {
		return fmt.Errorf("can't scan []byte of len %d into Address, want %d", len(srcB), AddressLength)
	}
	copy(a[:], srcB)
	return nil
}

// Value implements valuer for database/sql.
func (a PlatONAddress) Value() (driver.Value, error) {
	return a[:], nil
}

// UnprefixedAddress allows marshaling an Address without 0x prefix.
type UnprefixedPlatONAddress PlatONAddress

// UnmarshalText decodes the address from hex. The 0x prefix is optional.
func (a *UnprefixedPlatONAddress) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedUnprefixedText("UnprefixedAddress", input, a[:])
}

// MarshalText encodes the address as hex.
func (a UnprefixedPlatONAddress) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(a[:])), nil
}

// MixedcaseAddress retains the original string, which may or may not be
// correctly checksummed
type MixedcasePlatONAddress struct {
	addr     PlatONAddress
	original string
}

// NewMixedcaseAddress constructor (mainly for testing)
func NewMixedcasePlatONAddress(addr PlatONAddress) MixedcasePlatONAddress {
	return MixedcasePlatONAddress{addr: addr, original: addr.Hex()}
}

// NewMixedcasePlatONAddressFromString is mainly meant for unit-testing
func NewMixedcasePlatONAddressFromString(hexaddr string) (*MixedcasePlatONAddress, error) {
	if !IsHexAddress(hexaddr) {
		return nil, fmt.Errorf("Invalid address")
	}
	a := FromHex(hexaddr)
	return &MixedcasePlatONAddress{addr: BytesToPlatONAddress(a), original: hexaddr}, nil
}

// UnmarshalJSON parses MixedcasePlatONAddress
func (ma *MixedcasePlatONAddress) UnmarshalJSON(input []byte) error {
	if err := hexutil.UnmarshalFixedJSON(addressT, input, ma.addr[:]); err != nil {
		return err
	}
	return json.Unmarshal(input, &ma.original)
}

// MarshalJSON marshals the original value
func (ma *MixedcasePlatONAddress) MarshalJSON() ([]byte, error) {
	if strings.HasPrefix(ma.original, "0x") || strings.HasPrefix(ma.original, "0X") {
		return json.Marshal(fmt.Sprintf("0x%s", ma.original[2:]))
	}
	return json.Marshal(fmt.Sprintf("0x%s", ma.original))
}

// Address returns the address
func (ma *MixedcasePlatONAddress) Address() PlatONAddress {
	return ma.addr
}

// String implements fmt.Stringer
func (ma *MixedcasePlatONAddress) String() string {
	if ma.ValidChecksum() {
		return fmt.Sprintf("%s [chksum ok]", ma.original)
	}
	return fmt.Sprintf("%s [chksum INVALID]", ma.original)
}

// ValidChecksum returns true if the address has valid checksum
func (ma *MixedcasePlatONAddress) ValidChecksum() bool {
	return ma.original == ma.addr.Hex()
}

// Original returns the mixed-case input string
func (ma *MixedcasePlatONAddress) Original() string {
	return ma.original
}

// BytesToAddress returns Address with value b.
// If b is larger than len(h), b will be cropped from the left.
func BytesToNodePlatONAddress(b []byte) NodePlatONAddress {
	var a NodePlatONAddress
	a.SetBytes(b)
	return a
}

// HexToNodeAddress returns NodeAddress with byte values of s.
// If s is larger than len(h), s will be cropped from the left.
func HexToNodePlatONAddress(s string) NodePlatONAddress { return NodePlatONAddress(BytesToPlatONAddress(FromHex(s))) }

type NodePlatONAddress PlatONAddress

// Bytes gets the string representation of the underlying address.
func (a NodePlatONAddress) Bytes() []byte { return a[:] }

// Big converts an address to a big integer.
func (a NodePlatONAddress) Big() *big.Int { return new(big.Int).SetBytes(a[:]) }

// Hash converts an address to a hash by left-padding it with zeros.
func (a NodePlatONAddress) Hash() Hash { return BytesToHash(a[:]) }

// Hex returns an EIP55-compliant hex string representation of the address.
func (a NodePlatONAddress) Hex() string {
	unchecksummed := hex.EncodeToString(a[:])
	sha := sha3.NewKeccak256()
	sha.Write([]byte(unchecksummed))
	hash := sha.Sum(nil)

	result := []byte(unchecksummed)
	for i := 0; i < len(result); i++ {
		hashByte := hash[i/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if result[i] > '9' && hashByte > 7 {
			result[i] -= 32
		}
	}
	return "0x" + string(result)
}

func (a NodePlatONAddress) HexWithNoPrefix() string {
	unchecksummed := hex.EncodeToString(a[:])
	sha := sha3.NewKeccak256()
	sha.Write([]byte(unchecksummed))
	hash := sha.Sum(nil)

	result := []byte(unchecksummed)
	for i := 0; i < len(result); i++ {
		hashByte := hash[i/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if result[i] > '9' && hashByte > 7 {
			result[i] -= 32
		}
	}
	return string(result)
}

// String implements fmt.Stringer.
func (a NodePlatONAddress) String() string {
	return a.Hex()
}

// Format implements fmt.Formatter, forcing the byte slice to be formatted as is,
// without going through the stringer interface used for logging.
func (a NodePlatONAddress) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%"+string(c), a[:])
}

// SetBytes sets the address to the value of b.
// If b is larger than len(a) it will panic.
func (a *NodePlatONAddress) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// MarshalText returns the hex representation of a.
func (a NodePlatONAddress) MarshalText() ([]byte, error) {
	return hexutil.Bytes(a[:]).MarshalText()
}

// UnmarshalText parses a hash in hex syntax.
func (a *NodePlatONAddress) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("Address", input, a[:])
}

// UnmarshalJSON parses a hash in hex syntax.
func (a *NodePlatONAddress) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(addressT, input, a[:])
}

// Scan implements Scanner for database/sql.
func (a *NodePlatONAddress) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into Address", src)
	}
	if len(srcB) != AddressLength {
		return fmt.Errorf("can't scan []byte of len %d into Address, want %d", len(srcB), AddressLength)
	}
	copy(a[:], srcB)
	return nil
}

// Value implements valuer for database/sql.
func (a NodePlatONAddress) Value() (driver.Value, error) {
	return a[:], nil
}
