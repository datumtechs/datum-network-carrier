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

	"github.com/datumtechs/datum-network-carrier/common/bech32util"
	"github.com/datumtechs/datum-network-carrier/common/hexutil"
	"github.com/datumtechs/datum-network-carrier/crypto/sha3"
)

const (
	DefaultAlayaAddressHRP = "atp"
)

var currentAlayaAddressHRP string

func GetAlayaAddressHRP() string {
	if currentAlayaAddressHRP == "" {
		return DefaultAlayaAddressHRP
	}
	return currentAlayaAddressHRP
}

func SetAlayaAddressHRP(s string) error {
	if s == "" {
		s = DefaultAlayaAddressHRP
	}
	if len(s) != 3 {
		return errors.New("the length of AlayaAddressHRP must be 3")
	}
	log.Info("the AlayaAddress hrp has been set", "hrp", s)
	currentAlayaAddressHRP = s
	return nil
}

func CheckAlayaAddressHRP(s string) bool {
	if currentAlayaAddressHRP != "" && s != currentAlayaAddressHRP {
		return false
	}
	return true
}

/////////// AlayaAddress

// AlayaAddress represents the 20 byte AlayaAddress of an Ethereum account.
type AlayaAddress [AddressLength]byte

// BytesToAlayaAddress returns AlayaAddress with value b.
// If b is larger than len(h), b will be cropped from the left.
func BytesToAlayaAddress(b []byte) AlayaAddress {
	var a AlayaAddress
	a.SetBytes(b)
	return a
}

// BigToAlayaAddress returns AlayaAddress with byte values of b.
// If b is larger than len(h), b will be cropped from the left.
func BigToAlayaAddress(b *big.Int) AlayaAddress { return BytesToAlayaAddress(b.Bytes()) }

// Deprecated: AlayaAddress to string is use bech32 now
// HexToAlayaAddress returns AlayaAddress with byte values of s.
// If s is larger than len(h), s will be cropped from the left.
func HexToAlayaAddress(s string) AlayaAddress { return BytesToAlayaAddress(FromHex(s)) }

// MustBech32ToAlayaAddress returns AlayaAddress with byte values of s.
// If s is Decode fail, it will return zero AlayaAddress.
func MustBech32ToAlayaAddress(s string) AlayaAddress {
	add, err := Bech32ToAlayaAddress(s)
	if err != nil {
		log.Error("must Bech32ToAlayaAddress fail", "err", err)
		panic(err)
	}
	return add
}

// MustBech32ToAlayaAddress returns AlayaAddress with byte values of s.
// If s is Decode fail, it will return zero AlayaAddress.
func Bech32ToAlayaAddress(s string) (AlayaAddress, error) {
	hrpDecode, converted, err := bech32util.DecodeAndConvert(s)
	if err != nil {
		return AlayaAddress{}, err
	}
	if !CheckAlayaAddressHRP(hrpDecode) {
		return AlayaAddress{}, fmt.Errorf("the AlayaAddress hrp not compare right,input:%s", s)
	}

	if currentAlayaAddressHRP == "" {
		log.Warn("the AlayaAddress hrp not set yet", "input", s)
	} else if currentAlayaAddressHRP != hrpDecode {
		log.Warn("the AlayaAddress not compare current net", "want", currentAlayaAddressHRP, "input", s)
	}
	var a AlayaAddress
	a.SetBytes(converted)
	return a, nil
}

// Bech32ToAlayaAddressWithoutCheckHrp returns AlayaAddress with byte values of s.
// If s is Decode fail, it will return zero AlayaAddress.
func Bech32ToAlayaAddressWithoutCheckHrp(s string) AlayaAddress {
	_, converted, err := bech32util.DecodeAndConvert(s)
	if err != nil {
		log.Error(" hrp AlayaAddress decode  fail", "err", err)
		panic(err)
	}
	var a AlayaAddress
	a.SetBytes(converted)
	return a
}

// Deprecated: AlayaAddress to string is use bech32 now
// IsHexAlayaAddress verifies whether a string can represent a valid hex-encoded
// Ethereum AlayaAddress or not.
func IsHexAlayaAddress(s string) bool {
	if has0xPrefix(s) {
		s = s[2:]
	}
	return len(s) == 2*AddressLength && isHex(s)
}

func IsBech32AlayaAddress(s string) bool {
	hrp, _, err := bech32.Decode(s)
	if err != nil {
		return false
	}
	if !CheckAlayaAddressHRP(hrp) {
		return false
	}
	return true
}

// Bytes gets the string representation of the underlying AlayaAddress.
func (a AlayaAddress) Bytes() []byte { return a[:] }

// Big converts an AlayaAddress to a big integer.
func (a AlayaAddress) Big() *big.Int { return new(big.Int).SetBytes(a[:]) }

// Hash converts an AlayaAddress to a hash by left-padding it with zeros.
func (a AlayaAddress) Hash() Hash { return BytesToHash(a[:]) }

// Deprecated: AlayaAddress to string is use bech32 now
// Hex returns an EIP55-compliant hex string representation of the AlayaAddress.it's use for node AlayaAddress
func (a AlayaAddress) Hex() string {
	return "0x" + a.HexWithNoPrefix()
}

// Deprecated: AlayaAddress to string is use bech32 now
func (a AlayaAddress) HexWithNoPrefix() string {
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
func (a AlayaAddress) String() string {
	return a.Bech32()
}

func (a AlayaAddress) Bech32() string {
	return a.Bech32WithHRP(GetAlayaAddressHRP())
}

func (a AlayaAddress) Bech32WithHRP(hrp string) string {
	if v, err := bech32util.ConvertAndEncode(hrp, a.Bytes()); err != nil {
		log.Error("AlayaAddress can't ConvertAndEncode to string", "err", err, "add", a.Bytes())
		return ""
	} else {
		return v
	}
}

// Format implements fmt.Formatter, forcing the byte slice to be formatted as is,
// without going through the stringer interface used for logging.
func (a AlayaAddress) Format(s fmt.State, c rune) {
	switch string(c) {
	case "s":
		fmt.Fprintf(s, "%"+string(c), a.String())
	default:
		fmt.Fprintf(s, "%"+string(c), a[:])
	}
}

// SetBytes sets the AlayaAddress to the value of b.
// If b is larger than len(a) it will panic.
func (a *AlayaAddress) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// MarshalText returns the hex representation of a.
func (a AlayaAddress) MarshalText() ([]byte, error) {
	v, err := bech32util.ConvertAndEncode(GetAlayaAddressHRP(), a.Bytes())
	if err != nil {
		return nil, err
	}
	return []byte(v), nil
}

// UnmarshalText parses a hash in hex syntax.
func (a *AlayaAddress) UnmarshalText(input []byte) error {
	hrpDecode, converted, err := bech32util.DecodeAndConvert(string(input))
	if err != nil {
		return err
	}
	if !CheckAlayaAddressHRP(hrpDecode) {
		return fmt.Errorf("the AlayaAddress not compare current net,want %v,have %v", GetAlayaAddressHRP(), string(input))
	}
	a.SetBytes(converted)
	return nil
}

// UnmarshalJSON parses a hash in hex syntax.
func (a *AlayaAddress) UnmarshalJSON(input []byte) error {
	if !isAtpString(input) {
		return &json.UnmarshalTypeError{Value: "non-string", Type: addressT}
	}
	hrpDecode, v, err := bech32util.DecodeAndConvert(string(input[1 : len(input)-1]))
	if err != nil {
		return &json.UnmarshalTypeError{Value: err.Error(), Type: addressT}
	}
	if !CheckAlayaAddressHRP(hrpDecode) {
		return &json.UnmarshalTypeError{Value: fmt.Sprintf("hrpDecode not compare the current net,want %v,have %v", GetAlayaAddressHRP(), hrpDecode), Type: addressT}
	}
	a.SetBytes(v)
	return nil
}

func isAtpString(input []byte) bool {
	return len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"'
}

// Scan implements Scanner for database/sql.
func (a *AlayaAddress) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into AlayaAddress", src)
	}
	if len(srcB) != AddressLength {
		return fmt.Errorf("can't scan []byte of len %d into AlayaAddress, want %d", len(srcB), AddressLength)
	}
	copy(a[:], srcB)
	return nil
}

// Value implements valuer for database/sql.
func (a AlayaAddress) Value() (driver.Value, error) {
	return a[:], nil
}

// UnprefixedAlayaAddress allows marshaling an AlayaAddress without 0x prefix.
type UnprefixedAlayaAddress AlayaAddress

// UnmarshalText decodes the AlayaAddress from hex. The 0x prefix is optional.
func (a *UnprefixedAlayaAddress) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedUnprefixedText("UnprefixedAlayaAddress", input, a[:])
}

// MarshalText encodes the AlayaAddress as hex.
func (a UnprefixedAlayaAddress) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(a[:])), nil
}

// MixedcaseAlayaAddress retains the original string, which may or may not be
// correctly checksummed
type MixedcaseAlayaAddress struct {
	addr     AlayaAddress
	original string
}

// NewMixedcaseAlayaAddress constructor (mainly for testing)
func NewMixedcaseAlayaAddress(addr AlayaAddress) MixedcaseAlayaAddress {
	return MixedcaseAlayaAddress{addr: addr, original: addr.Hex()}
}

// NewMixedcaseAlayaAddressFromString is mainly meant for unit-testing
func NewMixedcaseAlayaAddressFromString(hexaddr string) (*MixedcaseAlayaAddress, error) {
	if !IsHexAlayaAddress(hexaddr) {
		return nil, fmt.Errorf("Invalid AlayaAddress")
	}
	a := FromHex(hexaddr)
	return &MixedcaseAlayaAddress{addr: BytesToAlayaAddress(a), original: hexaddr}, nil
}

// UnmarshalJSON parses MixedcaseAlayaAddress
func (ma *MixedcaseAlayaAddress) UnmarshalJSON(input []byte) error {
	if err := hexutil.UnmarshalFixedJSON(addressT, input, ma.addr[:]); err != nil {
		return err
	}
	return json.Unmarshal(input, &ma.original)
}

// MarshalJSON marshals the original value
func (ma *MixedcaseAlayaAddress) MarshalJSON() ([]byte, error) {
	if strings.HasPrefix(ma.original, "0x") || strings.HasPrefix(ma.original, "0X") {
		return json.Marshal(fmt.Sprintf("0x%s", ma.original[2:]))
	}
	return json.Marshal(fmt.Sprintf("0x%s", ma.original))
}

// AlayaAddress returns the AlayaAddress
func (ma *MixedcaseAlayaAddress) AlayaAddress() AlayaAddress {
	return ma.addr
}

// String implements fmt.Stringer
func (ma *MixedcaseAlayaAddress) String() string {
	if ma.ValidChecksum() {
		return fmt.Sprintf("%s [chksum ok]", ma.original)
	}
	return fmt.Sprintf("%s [chksum INVALID]", ma.original)
}

// ValidChecksum returns true if the AlayaAddress has valid checksum
func (ma *MixedcaseAlayaAddress) ValidChecksum() bool {
	return ma.original == ma.addr.Hex()
}

// Original returns the mixed-case input string
func (ma *MixedcaseAlayaAddress) Original() string {
	return ma.original
}

// BytesToAlayaAddress returns AlayaAddress with value b.
// If b is larger than len(h), b will be cropped from the left.
func BytesToNodeAlayaAddress(b []byte) NodeAlayaAddress {
	var a NodeAlayaAddress
	a.SetBytes(b)
	return a
}

// HexToNodeAlayaAddress returns NodeAlayaAddress with byte values of s.
// If s is larger than len(h), s will be cropped from the left.
func HexToNodeAlayaAddress(s string) NodeAlayaAddress { return NodeAlayaAddress(BytesToAlayaAddress(FromHex(s))) }

type NodeAlayaAddress AlayaAddress

// Bytes gets the string representation of the underlying AlayaAddress.
func (a NodeAlayaAddress) Bytes() []byte { return a[:] }

// Big converts an AlayaAddress to a big integer.
func (a NodeAlayaAddress) Big() *big.Int { return new(big.Int).SetBytes(a[:]) }

// Hash converts an AlayaAddress to a hash by left-padding it with zeros.
func (a NodeAlayaAddress) Hash() Hash { return BytesToHash(a[:]) }

// Hex returns an EIP55-compliant hex string representation of the AlayaAddress.
func (a NodeAlayaAddress) Hex() string {
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

func (a NodeAlayaAddress) HexWithNoPrefix() string {
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
func (a NodeAlayaAddress) String() string {
	return a.Hex()
}

// Format implements fmt.Formatter, forcing the byte slice to be formatted as is,
// without going through the stringer interface used for logging.
func (a NodeAlayaAddress) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%"+string(c), a[:])
}

// SetBytes sets the AlayaAddress to the value of b.
// If b is larger than len(a) it will panic.
func (a *NodeAlayaAddress) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// MarshalText returns the hex representation of a.
func (a NodeAlayaAddress) MarshalText() ([]byte, error) {
	return hexutil.Bytes(a[:]).MarshalText()
}

// UnmarshalText parses a hash in hex syntax.
func (a *NodeAlayaAddress) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("AlayaAddress", input, a[:])
}

// UnmarshalJSON parses a hash in hex syntax.
func (a *NodeAlayaAddress) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(addressT, input, a[:])
}

// Scan implements Scanner for database/sql.
func (a *NodeAlayaAddress) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into AlayaAddress", src)
	}
	if len(srcB) != AddressLength {
		return fmt.Errorf("can't scan []byte of len %d into AlayaAddress, want %d", len(srcB), AddressLength)
	}
	copy(a[:], srcB)
	return nil
}

// Value implements valuer for database/sql.
func (a NodeAlayaAddress) Value() (driver.Value, error) {
	return a[:], nil
}
