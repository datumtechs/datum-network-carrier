package signsuite

import (
	"errors"
	"fmt"
	"github.com/datumtechs/datum-network-carrier/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

var (
	ErrInvalidSig = errors.New("invalid msg v, r, s values")
)

type LatSigner struct {
	chainId, chainIdMul *big.Int
}
func NewLatSigner(chainId *big.Int) *LatSigner {
	if chainId == nil {
		chainId = new(big.Int)
	}
	return &LatSigner{
		chainId:    chainId,
		chainIdMul: new(big.Int).Mul(chainId, big.NewInt(2)),
	}
}

// Sender returns the sender address of the transaction.
func (s *LatSigner) Sender(hash common.Hash, sig []byte) (string, error)  {
	rb, sb, vb, err := s.SignatureValues(sig)
	if nil != err {

	}
	return recoverPlain(hash, rb, sb, vb)
}

// SignatureValues returns the raw R, S, V values corresponding to the
// given signature.
func (s *LatSigner) SignatureValues(sig []byte) (R, S, V *big.Int, err error) {
	if len(sig) != 65 {
		panic(fmt.Sprintf("wrong size for signature: got %d, want 65", len(sig)))
	}
	R = new(big.Int).SetBytes(sig[:32])
	S = new(big.Int).SetBytes(sig[32:64])
	V = new(big.Int).SetBytes([]byte{sig[64] + 35})
	V.Add(V, s.chainIdMul)

	return R, S, V, nil
}

//// Hash returns the hash to be signed.
//func (s *LatSigner) Hash(tx *Transaction) common.Hash

//// Equal returns true if the given signsuite is the same as the receiver.
//func (s *LatSigner) Equal(s2 Signer) bool {
//	eip155, ok := s2.(EIP155Signer)
//	return ok && eip155.chainId.Cmp(s.chainId) == 0
//}

func recoverPlain(sighash common.Hash, R, S, Vb *big.Int) (addr string, err error) {
	if Vb.BitLen() > 8 {
		return "", ErrInvalidSig
	}
	V := byte(Vb.Uint64() - 27)
	if !crypto.ValidateSignatureValues(V, R, S, true) {
		return "", ErrInvalidSig
	}
	// encode the signature in uncompressed format
	r, s := R.Bytes(), S.Bytes()
	sig := make([]byte, 65)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = V
	// recover the public key from the signature
	pub, err := crypto.Ecrecover(sighash[:], sig)
	if err != nil {
		return "", err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return "", errors.New("invalid public key")
	}
	var ad common.Address
	copy(ad[:], crypto.Keccak256(pub[1:])[12:])

	return ad.String(), nil
}
