package types

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
)

type OrgWallet struct {
	Address   common.Address    // wallet address
	PriKeyHex string            // private key ciphertext
	PriKey    *ecdsa.PrivateKey // private key
}
