package types

import "github.com/ethereum/go-ethereum/common"

type OrgWallet struct {
	Address common.Address // wallet address
	PriKey  string         // private key ciphertext
}
