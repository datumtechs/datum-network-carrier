package rlputil

import (
	"github.com/Metisnetwork/Metis-Carrier/common"
	"github.com/Metisnetwork/Metis-Carrier/crypto/sha3"
	"github.com/ethereum/go-ethereum/rlp"
)

// hasherPool holds Keccak hashers.
//var hasherPool = sync.Pool{
//	New: func() interface{} {
//		return sha3.NewKeccak256()
//	},
//}


//func RlpHash(x interface{}) (h common.Hash) {
//	sha := hasherPool.Get().(crypto.KeccakState)
//	defer hasherPool.Put(sha)
//	sha.Reset()
//	rlp.Encode(sha, x)
//	sha.Read(h[:])
//	return h
//}


func RlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}


