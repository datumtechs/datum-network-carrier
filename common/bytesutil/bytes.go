// Package bytesutil defines helper methods for converting integers to byte slices.
package bytesutil

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math"
	"math/bits"
	"regexp"

	"github.com/ethereum/go-ethereum/common/hexutil"
	types "github.com/prysmaticlabs/eth2-types"
)

// ToBytes returns integer x to bytes in little-endian format at the specified length.
// Spec defines similar method uint_to_bytes(n: uint) -> bytes, which is equivalent to ToBytes(n, 8).
func ToBytes(x uint64, length int) []byte {
	makeLength := length
	if length < 8 {
		makeLength = 8
	}
	bytes := make([]byte, makeLength)
	binary.LittleEndian.PutUint64(bytes, x)
	return bytes[:length]
}

// Bytes1 returns integer x to bytes in little-endian format, x.to_bytes(1, 'little').
func Bytes1(x uint64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, x)
	return bytes[:1]
}

// Bytes2 returns integer x to bytes in little-endian format, x.to_bytes(2, 'little').
func Bytes2(x uint64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, x)
	return bytes[:2]
}

// Bytes3 returns integer x to bytes in little-endian format, x.to_bytes(3, 'little').
func Bytes3(x uint64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, x)
	return bytes[:3]
}

// Bytes4 returns integer x to bytes in little-endian format, x.to_bytes(4, 'little').
func Bytes4(x uint64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, x)
	return bytes[:4]
}

// Bytes8 returns integer x to bytes in little-endian format, x.to_bytes(8, 'little').
func Bytes8(x uint64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, x)
	return bytes
}

// Bytes32 returns integer x to bytes in little-endian format, x.to_bytes(32, 'little').
func Bytes32(x uint64) []byte {
	bytes := make([]byte, 32)
	binary.LittleEndian.PutUint64(bytes, x)
	return bytes
}

// FromBytes4 returns an integer which is stored in the little-endian format(4, 'little')
// from a byte array.
func FromBytes4(x []byte) uint64 {
	empty4bytes := make([]byte, 4)
	return binary.LittleEndian.Uint64(append(x[:4], empty4bytes...))
}

// FromBytes8 returns an integer which is stored in the little-endian format(8, 'little')
// from a byte array.
func FromBytes8(x []byte) uint64 {
	return binary.LittleEndian.Uint64(x)
}

// ToBytes4 is a convenience method for converting a byte slice to a fix
// sized 4 byte array. This method will truncate the input if it is larger
// than 4 bytes.
func ToBytes4(x []byte) [4]byte {
	var y [4]byte
	copy(y[:], x)
	return y
}

// ToBytes32 is a convenience method for converting a byte slice to a fix
// sized 32 byte array. This method will truncate the input if it is larger
// than 32 bytes.
func ToBytes32(x []byte) [32]byte {
	var y [32]byte
	copy(y[:], x)
	return y
}

// ToBytes48 is a convenience method for converting a byte slice to a fix
// sized 48 byte array. This method will truncate the input if it is larger
// than 48 bytes.
func ToBytes48(x []byte) [48]byte {
	var y [48]byte
	copy(y[:], x)
	return y
}

// ToBytes64 is a convenience method for converting a byte slice to a fix
// sized 64 byte array. This method will truncate the input if it is larger
// than 64 bytes.
func ToBytes64(x []byte) [64]byte {
	var y [64]byte
	copy(y[:], x)
	return y
}

// ToBool is a convenience method for converting a byte to a bool.
// This method will use the first bit of the 0 byte to generate the returned value.
func ToBool(x byte) bool {
	return x&1 == 1
}

// FromBytes2 returns an integer which is stored in the little-endian format(2, 'little')
// from a byte array.
func FromBytes2(x []byte) uint16 {
	return binary.LittleEndian.Uint16(x[:2])
}

// FromBool is a convenience method for converting a bool to a byte.
// This method will use the first bit to generate the returned value.
func FromBool(x bool) byte {
	if x {
		return 1
	}
	return 0
}

// FromBytes48 is a convenience method for converting a fixed-size byte array
// to a byte slice.
func FromBytes48(x [48]byte) []byte {
	return x[:]
}

// FromBytes48Array is a convenience method for converting an array of
// fixed-size byte arrays to an array of byte slices.
func FromBytes48Array(x [][48]byte) [][]byte {
	y := make([][]byte, len(x))
	for i := range x {
		y[i] = x[i][:]
	}
	return y
}

// Trunc truncates the byte slices to 6 bytes.
func Trunc(x []byte) []byte {
	if len(x) > 6 {
		return x[:6]
	}
	return x
}

// ToLowInt64 returns the lowest 8 bytes interpreted as little endian.
func ToLowInt64(x []byte) int64 {
	if len(x) > 8 {
		x = x[:8]
	}
	return int64(binary.LittleEndian.Uint64(x))
}

// SafeCopyBytes will copy and return a non-nil byte array, otherwise it returns nil.
func SafeCopyBytes(cp []byte) []byte {
	if cp != nil {
		copied := make([]byte, len(cp))
		copy(copied, cp)
		return copied
	}
	return nil
}

// Copy2dBytes will copy and return a non-nil 2d byte array, otherwise it returns nil.
func Copy2dBytes(ary [][]byte) [][]byte {
	if ary != nil {
		copied := make([][]byte, len(ary))
		for i, a := range ary {
			copied[i] = SafeCopyBytes(a)
		}
		return copied
	}
	return nil
}

// ReverseBytes32Slice will reverse the provided slice's order.
func ReverseBytes32Slice(arr [][32]byte) [][32]byte {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

// PadTo pads a byte slice to the given size. If the byte slice is larger than the given size, the
// original slice is returned.
func PadTo(b []byte, size int) []byte {
	if len(b) > size {
		return b
	}
	return append(b, make([]byte, size-len(b))...)
}

// SetBit sets the index `i` of bitlist `b` to 1.
// It grows and returns a longer bitlist with 1 set
// if index `i` is out of range.
func SetBit(b []byte, i int) []byte {
	if i >= len(b)*8 {
		h := (i + (8 - i%8)) / 8
		b = append(b, make([]byte, h-len(b))...)
	}

	bit := uint8(1 << (i % 8))
	b[i/8] |= bit
	return b
}

// ClearBit clears the index `i` of bitlist `b`.
// Returns the original bitlist if the index `i`
// is out of range.
func ClearBit(b []byte, i int) []byte {
	if i >= len(b)*8 {
		return b
	}

	bit := uint8(1 << (i % 8))
	b[i/8] &^= bit
	return b
}

// MakeEmptyBitlists returns an empty bitlist with
// input size `i`.
func MakeEmptyBitlists(i int) []byte {
	return make([]byte, (i+(8-i%8))/8)
}

// HighestBitIndex returns the index of the highest
// bit set from bitlist `b`.
func HighestBitIndex(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, errors.New("input list can't be empty or nil")
	}

	for i := len(b) - 1; i >= 0; i-- {
		if b[i] == 0 {
			continue
		}
		return bits.Len8(b[i]) + (i * 8), nil
	}

	return 0, nil
}

// HighestBitIndexAt returns the index of the highest
// bit set from bitlist `b` that is at `index` (inclusive).
func HighestBitIndexAt(b []byte, index int) (int, error) {
	bLength := len(b)
	if b == nil || bLength == 0 {
		return 0, errors.New("input list can't be empty or nil")
	}

	start := index / 8
	if start >= bLength {
		start = bLength - 1
	}

	mask := byte(1<<(index%8) - 1)
	for i := start; i >= 0; i-- {
		if index/8 > i {
			mask = 0xff
		}
		masked := b[i] & mask
		minBitsMasked := bits.Len8(masked)
		if b[i] == 0 || (minBitsMasked == 0 && index/8 <= i) {
			continue
		}

		return minBitsMasked + (i * 8), nil
	}

	return 0, nil
}

// Uint64ToBytesLittleEndian conversion.
func Uint64ToBytesLittleEndian(i uint64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, i)
	return buf
}

// Uint64ToBytesBigEndian conversion.
func Uint64ToBytesBigEndian(i uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}

// BytesToUint64BigEndian conversion. Returns 0 if empty bytes or byte slice with length less
// than 8.
func BytesToUint64BigEndian(b []byte) uint64 {
	if len(b) < 8 { // This will panic otherwise.
		return 0
	}
	return binary.BigEndian.Uint64(b)
}

// EpochToBytesLittleEndian conversion.
func EpochToBytesLittleEndian(i types.Epoch) []byte {
	return Uint64ToBytesLittleEndian(uint64(i))
}

// EpochToBytesBigEndian conversion.
func EpochToBytesBigEndian(i types.Epoch) []byte {
	return Uint64ToBytesBigEndian(uint64(i))
}

// BytesToEpochBigEndian conversion.
func BytesToEpochBigEndian(b []byte) types.Epoch {
	return types.Epoch(BytesToUint64BigEndian(b))
}

// SlotToBytesBigEndian conversion.
func SlotToBytesBigEndian(i types.Slot) []byte {
	return Uint64ToBytesBigEndian(uint64(i))
}

// BytesToSlotBigEndian conversion.
func BytesToSlotBigEndian(b []byte) types.Slot {
	return types.Slot(BytesToUint64BigEndian(b))
}

// IsBytes32Hex checks whether the byte array is a 32-byte long hex number.
func IsBytes32Hex(b []byte) (bool, error) {
	if b == nil {
		return false, nil
	}
	return regexp.Match("^0x[0-9a-fA-F]{64}$", []byte(hexutil.Encode(b)))
}




// ToHex returns the hex representation of b, prefixed with '0x'.
// For empty slices, the return value is "0x0".
//
// Deprecated: use hexutil.Encode instead.
func ToHex(b []byte) string {
	hex := Bytes2Hex(b)
	if len(hex) == 0 {
		hex = "0"
	}
	return "0x" + hex
}

// ToHexArray creates a array of hex-string based on []byte
func ToHexArray(b [][]byte) []string {
	r := make([]string, len(b))
	for i := range b {
		r[i] = ToHex(b[i])
	}
	return r
}

// FromHex returns the bytes represented by the hexadecimal string s.
// s may be prefixed with "0x".
func FromHex(s string) []byte {
	if len(s) > 1 {
		if s[0:2] == "0x" || s[0:2] == "0X" {
			s = s[2:]
		}
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return Hex2Bytes(s)
}

// CopyBytes returns an exact copy of the provided bytes.
func CopyBytes(b []byte) (copiedBytes []byte) {
	if b == nil {
		return nil
	}
	copiedBytes = make([]byte, len(b))
	copy(copiedBytes, b)

	return
}

// hasHexPrefix validates str begins with '0x' or '0X'.
func hasHexPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

// isHexCharacter returns bool of c being a valid hexadecimal.
func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
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

// Bytes2Hex returns the hexadecimal encoding of d.
func Bytes2Hex(d []byte) string {
	return hex.EncodeToString(d)
}

/*func Bytes2Hex0x(d []byte) string {
	return hexutil.Encode(d)
}*/

// Hex2Bytes returns the bytes represented by the hexadecimal string str.
func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}

/*func Hex0x2Bytes(str string) []byte {
	h, _ := hexutil.Decode(str)
	return h
}*/

/*// Hex2BytesFixed returns bytes of a specified fixed length flen.
func Hex2BytesFixed(str string, flen int) []byte {
	h, _ := hex.DecodeString(str)
	if len(h) == flen {
		return h
	}
	if len(h) > flen {
		return h[len(h)-flen:]
	}
	hh := make([]byte, flen)
	copy(hh[flen-len(h):flen], h)
	return hh
}*/

// RightPadBytes zero-pads slice to the right up to length l.
func RightPadBytes(slice []byte, l int) []byte {
	if l <= len(slice) {
		return slice
	}

	padded := make([]byte, l)
	copy(padded, slice)

	return padded
}

// LeftPadBytes zero-pads slice to the left up to length l.
func LeftPadBytes(slice []byte, l int) []byte {
	if l <= len(slice) {
		return slice
	}

	padded := make([]byte, l)
	copy(padded[l-len(slice):], slice)

	return padded
}

/*func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}*/

func Int32ToBytes(n int32) []byte {
	tmp := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes()
}

func BytesToInt32(b []byte) int32 {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int32
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return int32(tmp)
}

func Int64ToBytes(n int64) []byte {
	tmp := int64(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes()
}

func BytesToInt64(b []byte) int64 {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int64
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return int64(tmp)
}

func Float32ToBytes(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

func BytesToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

func Float64ToBytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func BytesToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}

/*func PaddingLeft(src []byte, bytes int) []byte {
	if len(src) >= bytes {
		return src
	}
	dst := make([]byte, bytes)
	copy(dst, src)
	return reverse(dst)
}*/

/*func reverse(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}*/

func BytesToUint64(b []byte) uint64 {
	b = append(make([]byte, 8-len(b)), b...)
	return binary.BigEndian.Uint64(b)
}
func Uint64ToBytes(val uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, val)
	return buf[:]
}

func BytesToUint32(b []byte) uint32 {
	b = append(make([]byte, 4-len(b)), b...)
	return binary.BigEndian.Uint32(b)
}

func Uint32ToBytes(val uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, val)
	return buf[:]
}

func BytesToUint16(b []byte) uint16 {
	b = append(make([]byte, 2-len(b)), b...)
	return binary.BigEndian.Uint16(b)
}

func Uint16ToBytes(val uint16) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, val)
	return buf[:]
}

