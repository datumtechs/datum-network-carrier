package utils

import (
	"bytes"
	"github.com/datumtechs/datum-network-carrier/common"
	"github.com/datumtechs/datum-network-carrier/crypto/sha3"
	"io"
)

// BuildHash converts byte array to hash. Use sha256 to
// generate a unique message hash.
func BuildHash(msgType byte, bytes []byte) common.Hash {
	bytes[0] = msgType
	hashBytes := sha3.Sum256(bytes)
	result := common.Hash{}
	result.SetBytes(hashBytes[:])
	return result
}

// MergeBytes merges multiple bytes of data and
// returns the merged byte array.
func MergeBytes(bts ...[]byte) []byte {
	buffer := bytes.NewBuffer(make([]byte, 0, 128))
	for _, v := range bts {
		io.Copy(buffer, bytes.NewReader(v))
	}
	temp := buffer.Bytes()
	length := len(temp)
	var response []byte
	if cap(temp) > (length + length/10) {
		response = make([]byte, length)
		copy(response, temp)
	} else {
		response = temp
	}
	return response
}