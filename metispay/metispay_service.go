package metispay

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"github.com/RosettaFlow/Carrier-Go/common/hexutil"
	"github.com/RosettaFlow/Carrier-Go/metispay/contract"
	"github.com/RosettaFlow/Carrier-Go/metispay/kms"
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
	"sync"
	"time"
)

const keystoreFile = ".keystore"

var (
	//Metispay *MetisPayService
	contractMetisPayAddress = common.HexToAddress("0x147B8eb97fD247D06C4006D269c90C1908Fb5D54")

	once     sync.Once
	metisPay *MetisPayService
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

type MetisPayService struct {
	Config                   *Config
	Kms                      kms.KmsService
	client                   *ethclient.Client
	chainID                  *big.Int
	abi                      abi.ABI
	contractMetisPayInstance *contract.MetisPay
}

func (metisPay *MetisPayService) loadKeystore() {
	keyBytes, err := ioutil.ReadFile(keystoreFile)
	if err != nil {
		log.Fatalf("read keystore error, %v", err)
	}
	if metisPay.Kms != nil {
		keyBytes, err = metisPay.Kms.Decrypt(keyBytes)
		if err != nil {
			log.Fatalf("KMS decrypt keystore error, %v", err)
		}
	}

	key, err := keystore.DecryptKey(keyBytes, metisPay.Config.WalletPwd)
	if err != nil {
		log.Fatalf("decrypt keystore error, %v", err)
	}
	metisPay.Config.privateKey = key.PrivateKey
	metisPay.Config.walletAddress = key.Address
}

func GetMetisPayService() *MetisPayService {
	return metisPay
}

func InitMetisPayService(config *Config, kmsConfig *kms.Config) {
	once.Do(func() {
		log.Info("Init MetisPay service ...")
		metisPay = new(MetisPayService)
		if config != nil {
			metisPay.Config = config

			client, err := ethclient.Dial(config.URL)
			if err != nil {
				log.Fatal(err)
			}

			metisPay.client = client

			chainID, err := metisPay.client.NetworkID(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			metisPay.chainID = chainID

			instance, err := contract.NewMetisPay(contractMetisPayAddress, client)
			if err != nil {
				log.Fatal(err)
			}
			metisPay.contractMetisPayInstance = instance

			if kmsConfig != nil {
				metisPay.Kms = &kms.AliKms{Config: kmsConfig}
			}
			metisPay.loadKeystore()
		}

		abiCode, err := abi.JSON(strings.NewReader(contract.MetisPayABI))
		if err != nil {
			log.Fatal(err)
		}
		metisPay.abi = abiCode
	})
}

func (metisPay *MetisPayService) buildTxOpts() *bind.TransactOpts {
	nonce, err := metisPay.client.PendingNonceAt(context.Background(), metisPay.Config.walletAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := metisPay.client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	//txOpts := bind.NewKeyedTransactor(metisPay.Config.privateKey)
	txOpts, err := bind.NewKeyedTransactorWithChainID(metisPay.Config.privateKey, metisPay.chainID)
	if err != nil {
		log.Fatal(err)
	}

	txOpts.Nonce = big.NewInt(int64(nonce))
	txOpts.Value = big.NewInt(0)     // in wei
	txOpts.GasLimit = uint64(300000) // in units
	txOpts.GasPrice = gasPrice
	return txOpts
}

func (metisPay *MetisPayService) estimateGas(method string, params ...interface{}) uint64 {
	input, err := metisPay.abi.Pack(method, params)
	if err != nil {
		log.Fatal(err)
	}

	msg := ethereum.CallMsg{From: metisPay.Config.walletAddress, To: &contractMetisPayAddress, Data: input}
	estimatedGas, err := metisPay.client.EstimateGas(context.Background(), msg)
	if err != nil {
		log.Fatal(err)
	}
	return uint64(float64(estimatedGas) * 1.30)
}

func (metisPay *MetisPayService) Start(taskID string, userAccount common.Address, dataTokenList []common.Address) error {
	response := metisPay.Prepay(taskID, userAccount, dataTokenList)
	if !response.success {
		return errors.New("prepayment failure")
	}
	if len(response.settleId) == 0 {
		return errors.New("settleID missing from task prepayment results")
	}

	if !metisPay.Settle(response.settleId, response.gasRefund) {
		return errors.New("settlement failure")
	}
	return nil
}

type PrepayResponse struct {
	txHash    common.Hash
	gasRefund int64
	settleId  []byte
	success   bool
}

func (metisPay *MetisPayService) Prepay(taskID string, userAccount common.Address, dataTokenList []common.Address) PrepayResponse {
	spliterCharIdx := 0
	if idx := strings.Index(taskID, ":"); idx >= 0 {
		spliterCharIdx = idx
	}
	taskIDBigInt := hexutil.MustDecodeBig(taskID[spliterCharIdx:])

	//估算gas, +30%
	gasLimit := metisPay.estimateGas("Prepay", taskIDBigInt, new(big.Int).SetUint64(1), dataTokenList)

	tx, err := metisPay.contractMetisPayInstance.Prepay(metisPay.buildTxOpts(), taskIDBigInt, userAccount, new(big.Int).SetUint64(gasLimit), dataTokenList)
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("prepay txHash:%v, taskID:%s ", tx.Hash().Hex(), taskID)

	ch := make(chan *prepayReceipt)
	ticker := time.Tick(1 * time.Second)
	var receipt *prepayReceipt
	select {
	case receipt = <-ch:
		break
	case <-ticker:
		go metisPay.getPrepayReceipt(tx.Hash(), ch)
	}

	response := PrepayResponse{}
	response.txHash = tx.Hash()

	if receipt != nil && !receipt.success {
		response.success = false
		//todo:交易没有成功，那么如果gasUsed大于gasLimit,是否需要再发个交易补上？
		response.gasRefund = int64(gasLimit) - int64(receipt.gasUsed)
	}

	if receipt != nil && receipt.success {
		response.success = true
		//需要退回的（有可能是负数需要清算时补上）
		response.gasRefund = int64(gasLimit) - int64(receipt.gasUsed)
		response.settleId = receipt.settleId
	}

	return response
}

type prepayReceipt struct {
	gasUsed  uint64
	settleId []byte
	success  bool
}

func (metisPay *MetisPayService) getPrepayReceipt(txHash common.Hash, ch chan *prepayReceipt) {
	receipt, err := metisPay.client.TransactionReceipt(context.Background(), txHash)
	/*if err == platon.NotFound {
		return nil
	}else*/
	if err != nil {
		return
	}
	prepayReceipt := &prepayReceipt{}

	if receipt.Status == 0 { //FAILURE
		prepayReceipt.success = false
		prepayReceipt.gasUsed = receipt.GasUsed
	} else { //SUCCESS
		prepayReceipt.success = true
		prepayReceipt.gasUsed = receipt.GasUsed
		log := receipt.Logs[0]
		prepayReceipt.settleId = log.Data
	}
	ch <- prepayReceipt
}

func (metisPay *MetisPayService) Settle(settleID []byte, gasRefundPrepayment int64) bool {
	//估算gas, +30%
	gasLimit := metisPay.estimateGas("Settle", settleID, new(big.Int).SetUint64(1))
	if int64(gasLimit) > gasRefundPrepayment {
		gasLimit = uint64(int64(gasLimit) - gasRefundPrepayment)
	}

	tx, err := metisPay.contractMetisPayInstance.Settle(metisPay.buildTxOpts(), new(big.Int).SetBytes(settleID), new(big.Int).SetUint64(gasLimit))
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("settle txHash:%v, settleID:%s ", tx.Hash().Hex(), hexutil.Encode(settleID))

	ch := make(chan bool)
	ticker := time.Tick(1 * time.Second)
	var result bool
	select {
	case result = <-ch:
		break
	case <-ticker:
		go metisPay.getSettleReceipt(tx.Hash(), ch)
	}
	return result
}

func (metisPay *MetisPayService) getSettleReceipt(txHash common.Hash, ch chan bool) {
	receipt, err := metisPay.client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return
	}
	ch <- receipt.Status == 1
}
