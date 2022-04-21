package metispay

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/Metisnetwork/Metis-Carrier/ach/metispay/contracts"
	"github.com/Metisnetwork/Metis-Carrier/common"
	"github.com/Metisnetwork/Metis-Carrier/core"
	"github.com/Metisnetwork/Metis-Carrier/db"
	"github.com/Metisnetwork/Metis-Carrier/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	ethcore "github.com/ethereum/go-ethereum/core"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"testing"
	"time"
)

var (
	timespan = 4
)

func TestPrepayReceipt(t *testing.T) {
	ch := make(chan *Receipt)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var receipt *Receipt
LOOP:
	for {
		select {
		case receipt = <-ch:
			t.Log("get the receipt")
			break LOOP
		case <-ticker.C:
			t.Log("try to get the receipt")
			go getPrepayReceipt(common.HexToHash("0x00"), ch)
		}
	}
	t.Logf("receipt.success:%t", receipt.success)
	t.Logf("receipt.gasUsed:%d", receipt.gasUsed)
}

func getPrepayReceipt(txHash common.Hash, ch chan *Receipt) {
	if timespan > 0 {
		timespan--
		return
	}
	prepayReceipt := new(Receipt)
	prepayReceipt.success = true
	prepayReceipt.gasUsed = 100
	ch <- prepayReceipt
}

var (
	walletKey, _     = ethcrypto.HexToECDSA("b24285967575de7d5563e35213a806c60d69094faa509025f2ab5437017d343a")
	walletAddress    = ethcrypto.PubkeyToAddress(walletKey.PublicKey) //0xEFb8aeE7c9BC8c8f1472976299855e7059b8Ecda
	walletLatBalance = big.NewInt(2e15)

	tokenKey, _  = ethcrypto.HexToECDSA("42fe5edfa8327ceb7fe8ae059251f38528fd3bf8d65e26619bafcf60849790ec")
	tokenAddress = ethcrypto.PubkeyToAddress(tokenKey.PublicKey) //0xC00b0635a7660f1e7AA3Cfe789Eb04311c3A6E40
)

var metisManager *MetisPayManager

func init() {
	database := db.NewMemoryDatabase()
	carrierDB := core.NewDataCenter(context.Background(), database)

	wallet := &types.OrgWallet{Address: walletAddress, PriKey: hex.EncodeToString(ethcrypto.FromECDSA(walletKey))}
	carrierDB.StoreOrgWallet(wallet)

	//wss://devnetopenapi2.platon.network/ws
	//chainid:2203181

	config := &Config{
		URL:           "https://devnetopenapi2.platon.network/rpc",
		walletAddress: walletAddress,
		privateKey:    walletKey,
	}
	metisManager = NewMetisPayManager(carrierDB, config, nil)
}

func Test_getChainID(t *testing.T) {
	client, err := ethclient.Dial("https://devnetopenapi2.platon.network/rpc")
	if err != nil {
		t.Fatal(err)
	}
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("chainID:%s", chainID.String())
}

func Test_genKey(t *testing.T) {
	key, _ := ethcrypto.GenerateKey()
	keyHex := hex.EncodeToString(ethcrypto.FromECDSA(key))
	addr := ethcrypto.PubkeyToAddress(key.PublicKey)
	t.Log(keyHex)
	t.Log(addr)

	walletKey, _ = ethcrypto.HexToECDSA("b24285967575de7d5563e35213a806c60d69094faa509025f2ab5437017d343a")
	walletAddress = ethcrypto.PubkeyToAddress(walletKey.PublicKey)

	t.Log(walletAddress)
	add := common.Address{}
	t.Logf("add.hex:%s", add.Hex())
}

//在SimulatedBackend部署合约
func TestMetisPay_DeployMetisPay(t *testing.T) {
	var gasLimit uint64 = 8000029
	txOpts, _ := bind.NewKeyedTransactorWithChainID(walletKey, big.NewInt(1337))
	genAlloc := make(ethcore.GenesisAlloc)
	genAlloc[txOpts.From] = ethcore.GenesisAccount{Balance: big.NewInt(2e13)}

	sim := backends.NewSimulatedBackend(genAlloc, gasLimit)
	defer sim.Close()

	txOpts.Nonce = big.NewInt(int64(0))
	txOpts.Value = big.NewInt(0) // in wei
	txOpts.GasLimit = gasLimit   // in units
	txOpts.GasPrice = new(big.Int).SetUint64(1)

	// Deploy a token contract on the simulated blockchain
	addr, _, token, err := contracts.DeployMetisPay(txOpts, sim)
	if err != nil {
		t.Fatalf("Failed to deploy new token contract: %v", err)
	}
	t.Logf("MetisPay address: %s", addr.Hex())
	// Commit all pending transactions in the simulator and print the names again
	sim.Commit()
	taskID := new(big.Int).SetUint64(111)
	// Print the current (non existent) and pending name of the contract
	taskState, _ := token.TaskState(nil, taskID)
	fmt.Println("Pre-mining taskState:", taskState)

	taskState, _ = token.TaskState(&bind.CallOpts{Pending: true}, taskID)
	fmt.Println("Pre-mining pending taskState:", taskState)
	//-------------
}

//增加白名单, 需要walletAddress上有LAT
func TestMetisPay_AddWhiteList(t *testing.T) {
	opts, err := metisManager.buildTxOpts()
	if err != nil {
		t.Fatalf("failed to build transact options: %v", err)
	}

	tx, err := metisManager.contractMetisPayInstance.AddWhitelist(opts, walletAddress)
	if err != nil {
		t.Fatalf("failed to AddWhitelist : %v", err)
	}

	receipt := metisManager.GetReceipt(tx.Hash())
	t.Logf("AddWhitelist receipt: %v", receipt)
}

//任务gas预估，需要walletAddress上有LAT，并兑换有wLAT， tokenAddress上有token
func TestMetisPay_EstimateTaskGas(t *testing.T) {
	dataTokenTransferItemList := []string{tokenAddress.Hex()}
	gasLimit, gasPrice, err := metisManager.EstimateTaskGas(dataTokenTransferItemList)
	if err != nil {
		t.Fatalf("Failed to EstimateTaskGas : %v", err)
	}
	t.Logf("gasLimit: %d, gasPrice: %d", gasLimit, gasPrice)
}
