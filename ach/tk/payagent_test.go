package tk

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/datumtechs/chainclient"
	"github.com/datumtechs/datum-network-carrier/ach/tk/contracts"
	"github.com/datumtechs/datum-network-carrier/core"
	"github.com/datumtechs/datum-network-carrier/db"
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	"github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcore "github.com/ethereum/go-ethereum/core"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var (
	timespan = 4
)

var (
	walletKey, _     = ethcrypto.HexToECDSA("b24285967575de7d5563e35213a806c60d69094faa509025f2ab5437017d343a")
	walletAddress    = ethcrypto.PubkeyToAddress(walletKey.PublicKey) //0xEFb8aeE7c9BC8c8f1472976299855e7059b8Ecda
	walletLatBalance = big.NewInt(2e15)

	tkKey, _  = ethcrypto.HexToECDSA("42fe5edfa8327ceb7fe8ae059251f38528fd3bf8d65e26619bafcf60849790ec")
	tkAddress = ethcrypto.PubkeyToAddress(tkKey.PublicKey) //0xC00b0635a7660f1e7AA3Cfe789Eb04311c3A6E40
)

var payAgentContractProxy = ethcommon.HexToAddress("0x263B1D39843BF2e1DA27d827e749992fbD1f1577")

var payAgent *PayAgent

func setup() {
	database := db.NewMemoryDatabase()
	carrierDB := core.NewDataCenter(context.Background(), database)

	priKey := hex.EncodeToString(ethcrypto.FromECDSA(walletKey))
	carrierDB.SaveOrgPriKey(priKey)

	//wss://devnetopenapi2.platon.network/ws
	//chainid:2203181
	InitWalletManager(carrierDB, nil)

	WalletManagerInstance().GenerateWallet()

	ethContext := chainclient.NewEthClientContext("https://devnetopenapi2.platon.network/rpc", "lat", WalletManagerInstance())

	payAgent = NewPayAgent(ethContext, payAgentContractProxy)
}

func Test_getChainID(t *testing.T) {
	setup()

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
	setup()

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
func TestToken20Pay_DeployToken20Pay(t *testing.T) {
	setup()

	var genAlloc ethcore.GenesisAlloc
	var gasLimit uint64 = 8000029
	var sim *backends.SimulatedBackend
	genAlloc = make(ethcore.GenesisAlloc)
	genAlloc[walletAddress] = ethcore.GenesisAccount{Balance: big.NewInt(2e13)}
	genAlloc[tkAddress] = ethcore.GenesisAccount{Balance: big.NewInt(2e10)}
	sim = backends.NewSimulatedBackend(genAlloc, gasLimit)
	defer sim.Close()

	txOpts, _ := bind.NewKeyedTransactorWithChainID(walletKey, big.NewInt(1337))

	txOpts.Nonce = big.NewInt(int64(0))
	txOpts.Value = big.NewInt(0) // in wei
	txOpts.GasLimit = gasLimit   // in units
	txOpts.GasPrice = new(big.Int).SetUint64(1)

	// Deploy a tk contract on the simulated blockchain
	addr, _, tk, err := contracts.DeployDatumPay(txOpts, sim)
	if err != nil {
		t.Fatalf("Failed to deploy new tk contract: %v", err)
	}
	t.Logf("Token20Pay address: %s", addr.Hex())
	// Commit all pending transactions in the simulator and print the names again
	sim.Commit()
	taskID := new(big.Int).SetUint64(111)
	// Print the current (non existent) and pending name of the contract
	taskState, _ := tk.TaskState(nil, taskID)
	fmt.Println("Pre-mining taskState:", taskState)

	taskState, _ = tk.TaskState(&bind.CallOpts{Pending: true}, taskID)
	fmt.Println("Pre-mining pending taskState:", taskState)
	//-------------
}

//增加白名单, 需要在https://devnetopenapi2.platon.network/rpc开发链上， walletAddress上有LAT
func TestToken20Pay_AddWhiteList(t *testing.T) {
	setup()

	opts, err := payAgent.ethContext.BuildTxOpts(0, 500000)
	if err != nil {
		t.Fatalf("failed to build transact options: %v", err)
	}

	tx, err := payAgent.payAgentContractInstance.AddWhitelist(opts, walletAddress)
	if err != nil {
		t.Fatalf("failed to AddWhitelist : %v", err)
	}
	timeout := time.Duration(10) * time.Second
	ctx, cancelFn := context.WithTimeout(context.Background(), timeout)
	defer cancelFn()
	receipt := payAgent.GetReceipt(ctx, tx.Hash(), time.Duration(1000)*time.Millisecond)
	t.Logf("AddWhitelist receipt: %v", receipt)
}

//任务gas预估，需要在https://devnetopenapi2.platon.network/rpc开发链上，walletAddress上有LAT，并兑换有wLAT， tkAddress上有tk
func TestToken20Pay_EstimateTaskGas(t *testing.T) {
	setup()

	tkItem := new(carrierapipb.TkItem)
	tkItem.TkAddress = tkAddress.Hex()
	tkItem.TkType = constant.TkType_Tk20
	tkItem.Value = "1"

	dataTokenTransferItemList := []*carrierapipb.TkItem{tkItem}
	gasLimit, gasPrice, err := payAgent.EstimateTaskGas(walletAddress.Hex(), dataTokenTransferItemList)
	if err != nil {
		t.Fatalf("Failed to EstimateTaskGas : %v", err)
	}
	t.Logf("gasLimit: %d, gasPrice: %d", gasLimit, gasPrice)
}

func TestToken20Pay_Prepay(t *testing.T) {
	setup()

	database := db.NewMemoryDatabase()
	carrierDB := core.NewDataCenter(context.Background(), database)

	key, _ := ethcrypto.HexToECDSA("0481a0c35a0e22d25aeae127e948d02ebe7eb315620fb83421a5c8318260bb97")

	priKey := hex.EncodeToString(ethcrypto.FromECDSA(key))
	carrierDB.SaveOrgPriKey(priKey)

	ethContext := chainclient.NewEthClientContext("https://devnetopenapi2.platon.network/rpc", "lat", WalletManagerInstance())

	payAgent = NewPayAgent(ethContext, payAgentContractProxy)

	taskIdBytes, _ := hex.DecodeString("9977f8c9962d4eb67815022b7a079ba67382afd1bd3ed5d2df65d995d2ca6c41")

	taskID := new(big.Int).SetBytes(taskIdBytes)
	t.Logf("taskID:%s", hexutil.EncodeBig(taskID))

	taskSponsor := common.HexToAddress("0x6f852ba98639a001a315065ecaf2069c7479f4cc")

	tkItem1 := new(carrierapipb.TkItem)
	tkItem1.TkAddress = "0xe19Cfd8F9173155C26149818abd5dEcAA6F705F3"
	tkItem1.TkType = constant.TkType_Tk20
	tkItem1.Value = "1"

	tkItem2 := new(carrierapipb.TkItem)
	tkItem2.TkAddress = "0xE88695D3a3BA03ee6bB2130Ffd7869a8E368a0b4"
	tkItem2.TkType = constant.TkType_Tk20
	tkItem2.Value = "1"

	dataTokenTransferItemList := []*carrierapipb.TkItem{tkItem1, tkItem2}

	txHash, err := payAgent.PrepayTk20(taskID, taskSponsor, dataTokenTransferItemList)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("txHash:%s", txHash.Hex())

	timeout := time.Duration(60000) * time.Millisecond

	ctx, cancelFn := context.WithTimeout(context.Background(), timeout)
	//ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	receipt := payAgent.GetReceipt(ctx, txHash, time.Duration(500)*time.Millisecond)
	t.Logf("receipt.status: %d", receipt.Status)
}

func Test_bigintmul(t *testing.T) {
	//估算gas
	estimatedGas := uint64(430000)

	gasPrice := new(big.Int).SetUint64(1e11)
	//gas fee
	totalFeeUsed := new(big.Int).Mul(new(big.Int).SetUint64(estimatedGas), gasPrice)

	t.Logf("totalFeeUsed:%d", totalFeeUsed)
}

var (
	nonceLocker  sync.Mutex
	pendingNonce = uint64(0)
)

func getAndIncreaseNonce() uint64 {
	nonceLocker.Lock()
	defer nonceLocker.Unlock()

	current := pendingNonce
	atomic.AddUint64(&pendingNonce, 1)

	return current
}

func Test_getAndIncreaseNonce(t *testing.T) {
	for i := 0; i < 10; i++ {
		go func() {
			t.Logf("pendingNonce: %d", getAndIncreaseNonce())
		}()
	}
	time.Sleep(time.Duration(2) * time.Second)
}
