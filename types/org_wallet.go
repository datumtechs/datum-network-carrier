package types

import "github.com/ethereum/go-ethereum/common"

type OrgWallet struct {
	address common.Address // wallet address
	priKey  string         // private key ciphertext
}

func (w *OrgWallet) SetAddress(addr common.Address) {
	w.address = addr
}

func (w *OrgWallet) GetAddress() common.Address {
	return w.address
}

func (w *OrgWallet) SetPriKey(privateKeyHex string) {
	w.priKey = privateKeyHex
}

func (w *OrgWallet) GetPriKey() string {
	return w.priKey
}
