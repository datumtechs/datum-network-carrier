package signer

import "math/big"

type LatSigner struct {
	chainId, chainIdMul *big.Int
}


// Sender returns the sender address of the transaction.
func (s *LatSigner) Sender(sig []byte) (string, error)  {

}

// SignatureValues returns the raw R, S, V values corresponding to the
// given signature.
func (s *LatSigner) SignatureValues(sig []byte) (R, S, V *big.Int, err error) {

}

//// Hash returns the hash to be signed.
//func (s *LatSigner) Hash(tx *Transaction) common.Hash

// Equal returns true if the given signer is the same as the receiver.
func (s *LatSigner) Equal(s2 Signer) bool {
	eip155, ok := s2.(EIP155Signer)
	return ok && eip155.chainId.Cmp(s.chainId) == 0
}
