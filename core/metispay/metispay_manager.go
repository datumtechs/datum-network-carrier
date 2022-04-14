package metispay

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"github.com/RosettaFlow/Carrier-Go/common/hexutil"
	"github.com/RosettaFlow/Carrier-Go/core"
	"github.com/RosettaFlow/Carrier-Go/core/metispay/contracts"
	"github.com/RosettaFlow/Carrier-Go/core/metispay/kms"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	"github.com/RosettaFlow/Carrier-Go/types"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/big"
	"strings"
	"time"
)

const keystoreFile = ".keystore"
const mockTaskID = "task:0x81050ea1c64ab0ed96e50b151c36bcef180eea18d3b3245e7c4a42aa08638c58"

var (
	contractMetisPayAddress = common.HexToAddress("0xa23A5ab6BA4095c9Eb3affd99c3F9949f3e6170a")
)

func LoadCarrierKey() (*ecdsa.PrivateKey, common.Address) {
	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	if err != nil {
		log.Fatal("cannot load private key")
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return privateKey, fromAddress
}

type Config struct {
	URL           string `json:"chain"`
	WalletPwd     string
	privateKey    *ecdsa.PrivateKey
	walletAddress common.Address
}

type MetisPayManager struct {
	dataCenter               core.CarrierDB
	Config                   *Config
	Kms                      kms.KmsService
	client                   *ethclient.Client
	chainID                  *big.Int
	abi                      abi.ABI
	contractMetisPayInstance *contracts.MetisPay
}

func (metisPay *MetisPayManager) loadPrivateKey() {
	if wallet, err := metisPay.dataCenter.QueryOrgWallet(); err != nil {
		log.Errorf("load organization wallet error: %v", err)
		return
	} else if wallet == nil {
		log.Warnf("organization wallet not found, please generate it.")
		return
	} else {
		metisPay.Config.walletAddress = wallet.Address
		if metisPay.Kms != nil {
			if key, err := metisPay.Kms.Decrypt(wallet.PriKey); err != nil {
				log.Errorf("decrypt organization wallet private key error: %v", err)
				return
			} else {
				priKey, err := crypto.ToECDSA([]byte(key))
				if err != nil {
					log.Errorf("convert organization wallet private key to ECDSA error: %v", err)
					return
				}
				metisPay.Config.privateKey = priKey

			}
		} else {
			key, err := crypto.HexToECDSA(wallet.PriKey)
			if err != nil {
				log.Errorf("not a valid private key hex: %v", err)
				return
			} else {
				metisPay.Config.privateKey = key
			}
		}
	}
}

func (metisPay *MetisPayManager) getPrivateKey() *ecdsa.PrivateKey {
	if metisPay.Config.privateKey == nil {
		metisPay.loadPrivateKey()
	}
	return metisPay.Config.privateKey
}

func (metisPay *MetisPayManager) loadKeystore() {
	var content string
	keyBytes, err := ioutil.ReadFile(keystoreFile)
	if err != nil {
		log.Fatalf("read keystore error, %v", err)
	}
	content = string(keyBytes)

	if metisPay.Kms != nil {
		content, err = metisPay.Kms.Decrypt(content)
		if err != nil {
			log.Fatalf("KMS decrypt keystore error, %v", err)
		}
	}

	key, err := keystore.DecryptKey([]byte(content), metisPay.Config.WalletPwd)
	if err != nil {
		log.Fatalf("decrypt keystore error, %v", err)
	}
	metisPay.Config.privateKey = key.PrivateKey
	metisPay.Config.walletAddress = key.Address
}

func NewMetisPayManager(db core.CarrierDB, config *Config, kmsConfig *kms.Config) *MetisPayManager {
	log.Info("Init MetisPay manager ...")
	metisPay := new(MetisPayManager)
	metisPay.dataCenter = db

	if config != nil && len(config.URL) > 0 {
		metisPay.Config = config

		client, err := ethclient.Dial(config.URL)
		if err != nil {
			log.Fatal(err)
		}
		chainID, err := client.NetworkID(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		metisPay.client = client
		metisPay.chainID = chainID

		instance, err := contracts.NewMetisPay(contractMetisPayAddress, client)
		if err != nil {
			log.Fatal(err)
		}
		metisPay.contractMetisPayInstance = instance

		if kmsConfig != nil {
			metisPay.Kms = &kms.AliKms{Config: kmsConfig}
		}
		metisPay.loadPrivateKey()
	}

	abiCode, err := abi.JSON(strings.NewReader(contracts.MetisPayABI))
	if err != nil {
		log.Fatal(err)
	}
	metisPay.abi = abiCode
	return metisPay
}

func (metisPay *MetisPayManager) buildTxOpts() *bind.TransactOpts {
	nonce, err := metisPay.client.PendingNonceAt(context.Background(), metisPay.Config.walletAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := metisPay.client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	//txOpts := bind.NewKeyedTransactor(metisPay.Config.privateKey)
	txOpts, err := bind.NewKeyedTransactorWithChainID(metisPay.getPrivateKey(), metisPay.chainID)
	if err != nil {
		log.Fatal(err)
	}

	txOpts.Nonce = big.NewInt(int64(nonce))
	txOpts.Value = big.NewInt(0)     // in wei
	txOpts.GasLimit = uint64(300000) // in units
	txOpts.GasPrice = gasPrice
	return txOpts
}

func convert(dataTokenTransferItemList []*pb.DataTokenTransferItem) ([]common.Address, []*big.Int) {
	tokenAddressList := make([]common.Address, len(dataTokenTransferItemList))
	amountList := make([]*big.Int, len(dataTokenTransferItemList))
	for idx, item := range dataTokenTransferItemList {
		tokenAddressList[idx] = common.HexToAddress(item.Address)
		amountList[idx] = new(big.Int).SetUint64(item.Amount)
	}
	return tokenAddressList, amountList
}
func (metisPay *MetisPayManager) EstimateTaskGas(dataTokenTransferItemList []*pb.DataTokenTransferItem) (uint64, error) {

	tokenAddressList, amountList := convert(dataTokenTransferItemList)

	gasLimit1, err := metisPay.estimateGas("Prepay", mockTaskID, new(big.Int).SetUint64(1), tokenAddressList, amountList)

	if err != nil {
		log.Errorf("call EstimateTaskGas error: %v", err)
		return 0, err
	}
	gasLimit2, err := metisPay.estimateGas("Settle", mockTaskID, new(big.Int).SetUint64(1))
	if err != nil {
		log.Errorf("call EstimateTaskGas error: %v", err)
		return 0, err
	}
	return gasLimit1 + gasLimit2, nil
}

func (metisPay *MetisPayManager) estimateGas(method string, params ...interface{}) (uint64, error) {
	input, err := metisPay.abi.Pack(method, params)
	if err != nil {
		return 0, err
	}

	msg := ethereum.CallMsg{From: metisPay.Config.walletAddress, To: &contractMetisPayAddress, Data: input}
	estimatedGas, err := metisPay.client.EstimateGas(context.Background(), msg)
	if err != nil {
		return 0, err
	}
	return estimatedGas, nil
}

// ReadyToStart get funds clearing to start a task.
// the caller should use a channel to receive the result.
func (metisPay *MetisPayManager) ReadyToStart(taskID string, userAccount common.Address, dataTokenTransferItemList []*pb.DataTokenTransferItem) bool {
	if metisPay.getPrivateKey() == nil {
		log.Errorf("cannot get funds clearing to start a task cause organization wallet missing")
		return false
	}

	spliterIdx := 0
	if idx := strings.Index(taskID, ":"); idx >= 0 {
		spliterIdx = idx
	}
	taskIDBigInt, err := hexutil.DecodeBig(taskID[spliterIdx:])
	if err != nil {
		log.Errorf("cannot decode taskID to big.Int, %v", err)
		return false
	}

	response := metisPay.Prepay(taskIDBigInt, userAccount, dataTokenTransferItemList)
	if !response.success {
		return false
	}

	if !metisPay.Settle(taskIDBigInt, response.gasRefund) {
		return false
	}
	return true
}

type PrepayResponse struct {
	txHash    common.Hash
	gasRefund int64
	success   bool
}

func (metisPay *MetisPayManager) QueryOrgWallet() (string, error) {
	wallet, err := metisPay.dataCenter.QueryOrgWallet()

	if nil != err {
		log.WithError(err).Error("failed to query organization wallet. ", err)
		return "", errors.New("failed to query organization wallet")
	}

	if wallet != nil {
		return wallet.Address.Hex(), nil
	}

	return "", nil
}

func (metisPay *MetisPayManager) GenerateOrgWallet() (string, error) {
	walletAddr, err := metisPay.QueryOrgWallet()
	if err != nil {
		return "", err
	}
	if len(walletAddr) > 0 {
		log.Warnf("organization wallet exists, just retuens current wallet: %s", walletAddr)
		return walletAddr, nil
	}

	key, _ := crypto.GenerateKey()
	keyHex := hex.EncodeToString(crypto.FromECDSA(key))
	addr := crypto.PubkeyToAddress(key.PublicKey)

	if metisPay.Kms != nil {
		if cipher, err := metisPay.Kms.Encrypt(keyHex); err != nil {
			return "", errors.New("cannot encrypt organization wallet private key")
		} else {
			keyHex = cipher
		}
	}

	wallet := new(types.OrgWallet)
	wallet.PriKey = keyHex
	wallet.Address = addr
	if err := metisPay.dataCenter.StoreOrgWallet(wallet); err != nil {
		log.WithError(err).Error("Failed to store organization wallet")
		return "", errors.New("failed to store organization wallet")
	} else {
		return addr.Hex(), nil
	}
}

func (metisPay *MetisPayManager) Prepay(taskID *big.Int, userAccount common.Address, dataTokenTransferItemList []*pb.DataTokenTransferItem) *PrepayResponse {
	response := new(PrepayResponse)
	if metisPay.getPrivateKey() == nil {
		log.Errorf("cannot send Prepay transaction cause organization wallet missing")
		response.success = false
		return response
	}

	tokenAddressList, amountList := convert(dataTokenTransferItemList)

	//估算gas, +30%
	gasLimit, err := metisPay.estimateGas("Prepay", taskID, new(big.Int).SetUint64(1), tokenAddressList, amountList)
	if err != nil {
		log.Errorf("call constract's estimateGas error: %v", err)
		response.success = false
		return response
	}
	gasLimit = uint64(float64(gasLimit) * 1.30)

	tx, err := metisPay.contractMetisPayInstance.Prepay(metisPay.buildTxOpts(), taskID, userAccount, new(big.Int).SetUint64(gasLimit), tokenAddressList, amountList)
	if err != nil {
		log.Errorf("call constract's Prepay error: %v", err)
		response.success = false
		return response
	}
	log.Debugf("call constract's Prepay txHash:%v, taskID:%s ", tx.Hash().Hex(), hexutil.EncodeBig(taskID))
	response.txHash = tx.Hash()

	ch := make(chan *prepayReceipt)
	ticker := time.NewTicker(1 * time.Second)
	var receipt *prepayReceipt

LOOP:
	for {
		select {
		case receipt = <-ch:
			break LOOP
		case <-ticker.C:
			go metisPay.getPrepayReceipt(tx.Hash(), ch)
		}
	}

	if receipt != nil && !receipt.success {
		response.success = false
		//todo:交易没有成功，那么如果gasUsed大于gasLimit,是否需要再发个交易补上？
		response.gasRefund = int64(gasLimit) - int64(receipt.gasUsed)
	}

	if receipt != nil && receipt.success {
		response.success = true
		//需要退回的（有可能是负数需要清算时补上）
		response.gasRefund = int64(gasLimit) - int64(receipt.gasUsed)
	}
	return response
}

type prepayReceipt struct {
	gasUsed uint64
	success bool
}

func (metisPay *MetisPayManager) getPrepayReceipt(txHash common.Hash, ch chan *prepayReceipt) {
	receipt, err := metisPay.client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		//including NotFound
		log.Errorf("get Prepay transaction receipt error: %v", err)
		return
	}
	log.Debugf("get Prepay transaction receipt status: %d", receipt.Status)

	prepayReceipt := new(prepayReceipt)
	if receipt.Status == 0 { //FAILURE
		prepayReceipt.success = false
		prepayReceipt.gasUsed = receipt.GasUsed
	} else { //SUCCESS
		prepayReceipt.success = true
		prepayReceipt.gasUsed = receipt.GasUsed
	}
	ch <- prepayReceipt
}

func (metisPay *MetisPayManager) Settle(taskID *big.Int, gasRefundPrepayment int64) bool {
	if metisPay.getPrivateKey() == nil {
		log.Errorf("cannot send Settle transaction cause organization wallet missing")
		return false
	}

	//估算gas, +30%
	gasLimit, err := metisPay.estimateGas("Settle", taskID, new(big.Int).SetUint64(1))
	if err != nil {
		log.Errorf("call constract's estimateGas error: %v", err)
		return false
	}
	gasLimit = uint64(float64(gasLimit) * 1.30)

	if int64(gasLimit) > gasRefundPrepayment {
		gasLimit = uint64(int64(gasLimit) - gasRefundPrepayment)
	}

	tx, err := metisPay.contractMetisPayInstance.Settle(metisPay.buildTxOpts(), taskID, new(big.Int).SetUint64(gasLimit))
	if err != nil {
		log.Errorf("call constract's Settle error: %v", err)
		return false
	}
	log.Debugf("call constract's Settle txHash:%v, taskID:%s ", tx.Hash().Hex(), hexutil.EncodeBig(taskID))

	ch := make(chan bool)
	ticker := time.NewTicker(1 * time.Second)
	var result bool

LOOP:
	select {
	case result = <-ch:
		break LOOP
	case <-ticker.C:
		go metisPay.getSettleReceipt(tx.Hash(), ch)
	}
	return result
}

func (metisPay *MetisPayManager) getSettleReceipt(txHash common.Hash, ch chan bool) {
	receipt, err := metisPay.client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		log.Errorf("get Settle transaction receipt error: %v", err)
		return
	}
	log.Debugf("get Settle transaction receipt status: %d", receipt.Status)
	ch <- receipt.Status == 1
}
