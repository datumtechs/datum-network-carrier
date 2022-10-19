package tk

import (
	"github.com/datumtechs/chainclient"
	"github.com/datumtechs/datum-network-carrier/carrierdb"
	"github.com/datumtechs/datum-network-carrier/core"
	"github.com/datumtechs/datum-network-carrier/db"
	"github.com/ethereum/go-ethereum/common"
	"testing"
)

func Test_Address(t *testing.T) {

	if addr := getAddress(); addr == (common.Address{}) {
		t.Log("zero length")
	} else {
		t.Log("not zero length")
	}
}

func getAddress() common.Address {
	return common.Address{}
}

func newMemoryCarrierDB() carrierdb.CarrierDB {
	db := db.NewMemoryDatabase()
	return core.NewDataCenter(nil, db)
}
func Test_InitWalletManager(t *testing.T) {

	InitWalletManager(newMemoryCarrierDB(), nil)

	addr, err := WalletManagerInstance().GenerateWallet()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("generate Addr:%s", addr.Hex())
	WalletManagerInstance().loadPrivateKey()

	t.Logf("load addr:%s", addr.Hex())
	var ethContext *chainclient.EthContext
	ethContext = chainclient.NewEthClientContext("http://8.219.126.197:6789", "lat", WalletManagerInstance())

	t.Logf("ethContext.GetAddress():%s", ethContext.GetAddress().Hex())
}
