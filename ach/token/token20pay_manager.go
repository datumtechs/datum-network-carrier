package token

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/datumtechs/datum-network-carrier/ach/token/contracts"
	"github.com/datumtechs/datum-network-carrier/ach/token/kms"
	"github.com/datumtechs/datum-network-carrier/common/hexutil"
	"github.com/datumtechs/datum-network-carrier/core"
	"github.com/datumtechs/datum-network-carrier/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	ethereumtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/big"
	"strings"
)

const keystoreFile = ".keystore"

//"task:0x81050ea1c64ab0ed96e50b151c36bcef180eea18d3b3245e7c4a42aa08638c58"
var (
	contractToken20PayAddress        = common.HexToAddress("0x263B1D39843BF2e1DA27d827e749992fbD1f1577")
	defaultDataTokenPrepaymentAmount = big.NewInt(1e18)

	mockTaskID, _ = hexutil.DecodeBig("0x81050ea1c64ab0ed96e50b151c36bcef180eea18d3b3245e7c4a42aa08638c58")
)

type Config struct {
	URL           string `json:"chain"`
	WalletPwd     string
	privateKey    *ecdsa.PrivateKey
	walletAddress common.Address
}

type Token20PayManager struct {
	dataCenter                 core.CarrierDB
	Config                     *Config
	Kms                        kms.KmsService
	client                     *ethclient.Client
	chainID                    *big.Int
	abi                        abi.ABI
	contractToken20PayInstance *contracts.Token20Pay
	nonceLocker                sync.Mutex
	pendingNonce               uint64
}

func (m *Token20PayManager) loadPrivateKey() {
	if wallet, err := m.dataCenter.QueryOrgWallet(); err != nil {
		log.Errorf("load organization wallet error: %v", err)
		return
	} else if wallet == nil {
		log.Warnf("organization wallet not found, please generate it.")
		return
	} else {
		m.Config.walletAddress = wallet.Address
		if m.Kms != nil {
			if key, err := m.Kms.Decrypt(wallet.PriKey); err != nil {
				log.Errorf("decrypt organization wallet private key error: %v", err)
				return
			} else {
				priKey, err := crypto.HexToECDSA(key)
				if err != nil {
					log.Errorf("convert organization wallet private key to ECDSA error: %v", err)
					return
				}
				m.Config.privateKey = priKey
			}
		} else {
			key, err := crypto.HexToECDSA(wallet.PriKey)
			if err != nil {
				log.Errorf("not a valid private key hex: %v", err)
				return
			} else {
				m.Config.privateKey = key
			}
		}
	}
	log.Debugf("load organization wallet successful. address: %s, privateKey: %s", m.Config.walletAddress, hex.EncodeToString(crypto.FromECDSA(m.Config.privateKey)))
}

func (m *Token20PayManager) getPrivateKey() *ecdsa.PrivateKey {
	if m.Config.privateKey == nil {
		m.loadPrivateKey()
	}
	return m.Config.privateKey
}

func (m *Token20PayManager) loadKeystore() {
	var content string
	keyBytes, err := ioutil.ReadFile(keystoreFile)
	if err != nil {
		log.Fatalf("read keystore error, %v", err)
	}
	content = string(keyBytes)

	if m.Kms != nil {
		content, err = m.Kms.Decrypt(content)
		if err != nil {
			log.Fatalf("KMS decrypt keystore error, %v", err)
		}
	}

	key, err := keystore.DecryptKey([]byte(content), m.Config.WalletPwd)
	if err != nil {
		log.Fatalf("decrypt keystore error, %v", err)
	}
	m.Config.privateKey = key.PrivateKey
	m.Config.walletAddress = key.Address
}

func NewToken20PayManager(db core.CarrierDB, config *Config, kmsConfig *kms.Config) *Token20PayManager {
	log.Info("Init Token20Pay manager ...")
	m := new(Token20PayManager)
	m.dataCenter = db

	if config != nil && len(config.URL) > 0 {
		m.Config = config

		client, err := ethclient.Dial(config.URL)
		if err != nil {
			log.Fatal(err)
		}
		chainID, err := client.ChainID(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		m.client = client
		m.chainID = chainID

		instance, err := contracts.NewToken20Pay(contractToken20PayAddress, client)
		if err != nil {
			log.Fatal(err)
		}
		m.contractToken20PayInstance = instance

		if kmsConfig != nil {
			m.Kms = &kms.AliKms{Config: kmsConfig}
		}
		m.loadPrivateKey()

		m.initPendingNonce()
	}

	abiCode, err := abi.JSON(strings.NewReader(contracts.Token20PayABI))
	if err != nil {
		log.Fatal(err)
	}
	m.abi = abiCode
	return m
}

func (m *Token20PayManager) initPendingNonce() {
	if pendingNonce, err := m.client.PendingNonceAt(context.Background(), m.Config.walletAddress); err != nil {
		log.Fatalf("cannot init pending nonce: %v", err)
	} else {
		m.pendingNonce = pendingNonce
	}
}

func (m *Token20PayManager) getAndIncreasePendingNonce() uint64 {
	m.nonceLocker.Lock()
	defer m.nonceLocker.Unlock()

	current := m.pendingNonce
	atomic.AddUint64(&m.pendingNonce, 1)

	return current
}

func (m *Token20PayManager) buildTxOpts(gasLimit uint64) (*bind.TransactOpts, error) {

	gasPrice, err := m.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	//txOpts := bind.NewKeyedTransactor(m.Config.privateKey)
	txOpts, err := bind.NewKeyedTransactorWithChainID(m.getPrivateKey(), m.chainID)
	if err != nil {
		return nil, err
	}

	txOpts.Nonce = new(big.Int).SetUint64(m.getAndIncreasePendingNonce())
	txOpts.Value = big.NewInt(0) // in wei
	txOpts.GasLimit = gasLimit   // in units
	txOpts.GasPrice = gasPrice.Mul(gasPrice, big.NewInt(2))
	return txOpts, nil
}

func convert(dataTokenAddressList []string) ([]common.Address, []*big.Int) {
	tokenAddressList := make([]common.Address, len(dataTokenAddressList))
	dataTokenAmountList := make([]*big.Int, len(dataTokenAddressList))
	for idx, addrHex := range dataTokenAddressList {
		tokenAddressList[idx] = common.HexToAddress(addrHex)
		dataTokenAmountList[idx] = defaultDataTokenPrepaymentAmount
	}
	return tokenAddressList, dataTokenAmountList
}

// EstimateTaskGas estimates gas fee for a task's sponsor.
// EstimateTaskGas returns estimated gas and suggested gas price.
func (m *Token20PayManager) EstimateTaskGas(taskSponsorAddress string, dataTokenAddressList []string) (uint64, *big.Int, error) {
	log.Debugf("call EstimateTaskGas, sponsorAddress: %s, tokenAddressList:%v", taskSponsorAddress, dataTokenAddressList)

	tokenAddressList, tokenAmountList := convert(dataTokenAddressList)

	// Estimating task gas happens before task starting, and the task ID has not been generated at this moment, so, apply a mock task ID.
	gasLimit1, err := m.estimateGas("prepay", mockTaskID, common.HexToAddress(taskSponsorAddress), big.NewInt(1), tokenAddressList, tokenAmountList)
	if err != nil {
		log.Errorf("call EstimateTaskGas error: %v", err)
		return 0, nil, err
	}

	//cannot settle a mock task.
	// number of transactions in settlement is same as in prepayment, so settlement gas limit is considered same as prepayment gas limit.
	/*
		gasLimit2, err := m.estimateGas("settle", mockTaskID, new(big.Int).SetUint64(1))
		if err != nil {
			log.Errorf("call EstimateTaskGas error: %v", err)
			return 0, nil, err
		}
	*/

	gasPrice, err := m.client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Errorf("call SuggestGasPrice error: %v", err)
		return 0, nil, err
	}
	return gasLimit1 + gasLimit1, gasPrice, nil
}

// estimateGas call EstimateGas() by RPC
func (m *Token20PayManager) estimateGas(method string, params ...interface{}) (uint64, error) {
	input, err := m.abi.Pack(method, params...)
	if err != nil {
		return 0, err
	}

	msg := ethereum.CallMsg{From: m.Config.walletAddress, To: &contractToken20PayAddress, Data: input, Gas: 0, GasPrice: big.NewInt(0)}
	estimatedGas, err := m.client.EstimateGas(context.Background(), msg)
	if err != nil {
		return 0, err
	}
	return estimatedGas, nil
}

// QueryOrgWallet returns the organization wallet address
func (m *Token20PayManager) QueryOrgWallet() (common.Address, error) {
	wallet, err := m.dataCenter.QueryOrgWallet()

	if nil != err {
		log.WithError(err).Error("failed to query organization wallet. ", err)
		return common.Address{}, errors.New("failed to query organization wallet")
	}

	if wallet != nil {
		return wallet.Address, nil
	}

	return common.Address{}, nil
}

// GenerateOrgWallet generates a new wallet if there's no wallet for current organization, if there is an organization wallet already, just reuse it.
func (m *Token20PayManager) GenerateOrgWallet() (common.Address, error) {
	walletAddr, err := m.QueryOrgWallet()
	if err != nil {
		return common.Address{}, err
	}

	if walletAddr != (common.Address{}) {
		log.Warnf("organization wallet exists, just retuens current wallet: %s", walletAddr)
		return walletAddr, nil
	}

	key, _ := crypto.GenerateKey()
	keyHex := hex.EncodeToString(crypto.FromECDSA(key))
	addr := crypto.PubkeyToAddress(key.PublicKey)

	if m.Kms != nil {
		if cipher, err := m.Kms.Encrypt(keyHex); err != nil {
			return common.Address{}, errors.New("cannot encrypt organization wallet private key")
		} else {
			keyHex = cipher
		}
	}

	wallet := new(types.OrgWallet)
	wallet.PriKey = keyHex
	wallet.Address = addr
	if err := m.dataCenter.StoreOrgWallet(wallet); err != nil {
		log.WithError(err).Error("Failed to store organization wallet")
		return common.Address{}, errors.New("failed to store organization wallet")
	} else {
		m.Config.privateKey = key
		m.Config.walletAddress = addr
		return addr, nil
	}
}

// Prepay transfers more than enough gas from task sponsor to Token20Pay, this gas will pay carrier to call Prepay()/Settle(), and remaining gas will refund to task sponsor.
// Prepay returns hx.Hash, and error.
// The complete procedure consists of two calls to Token20Pay, the first is Prepay, the second is Settle.
func (m *Token20PayManager) Prepay(taskID *big.Int, taskSponsorAccount common.Address, dataTokenAddressList []common.Address) (common.Hash, error) {
	taskIDHex := hexutil.EncodeBig(taskID)
	if m.getPrivateKey() == nil {
		log.Errorf("cannot send Prepay transaction cause organization wallet missing")
		return common.Hash{}, errors.New("organization private key is missing")
	}

	// for debug log...
	addrs := make([]string, len(dataTokenAddressList))
	amounts := make([]string, len(dataTokenAddressList))

	dataTokenAmountList := make([]*big.Int, len(dataTokenAddressList))
	for idx, addr := range dataTokenAddressList {
		dataTokenAmountList[idx] = defaultDataTokenPrepaymentAmount

		addrs[idx] = addr.String()
		amounts[idx] = defaultDataTokenPrepaymentAmount.String()
	}
	//估算gas
	gasEstimated, err := m.estimateGas("prepay", taskID, taskSponsorAccount, new(big.Int).SetUint64(1), dataTokenAddressList, dataTokenAmountList)
	if err != nil {
		log.Errorf("failed to estimate gas for Token20Pay.Prepay(), taskID: %s, error: %v", taskIDHex, err)
		return common.Hash{}, errors.New("failed to estimate gas for Token20Pay.Prepay()")
	}

	//交易参数直接使用用户预付的总的gas，尽量放大，以防止交易执行gas不足
	opts, err := m.buildTxOpts(gasEstimated * 2)
	if err != nil {
		log.Errorf("failed to build transact options to call Token20Pay.Prepay(): %v", err)
		return common.Hash{}, errors.New("failed to build transact options to call Token20Pay.Prepay()")
	}

	//gas fee, 支付助手合约，需要记录用户预付的手续费
	feePrepaid := new(big.Int).Mul(new(big.Int).SetUint64(gasEstimated), opts.GasPrice)

	log.Debugf("Start call contract prepay(), taskID: %s, gasEstimated: %d, gasLimit: %d, gasPrice: %d, feePrepaid: %d, taskSponsorAccount: %s, dataTokenAddressList: %s, dataTokenAmountList: %s",
		taskIDHex, gasEstimated, opts.GasLimit, opts.GasPrice, feePrepaid, taskSponsorAccount.String(), "["+strings.Join(addrs, ",")+"]", "["+strings.Join(amounts, ",")+"]")

	tx, err := m.contractToken20PayInstance.Prepay(opts, taskID, taskSponsorAccount, feePrepaid, dataTokenAddressList, dataTokenAmountList)
	if err != nil {
		log.WithError(err).Errorf("failed to call Token20Pay.Prepay(), taskID: %s", taskIDHex)
		return common.Hash{}, errors.New("failed to call Token20Pay.Prepay()")
	}
	log.Debugf("call Token20Pay.Prepay() txHash:%v, taskID:%s", tx.Hash().Hex(), taskIDHex)

	return tx.Hash(), nil
}

// Settle get funds clearing,
// 1. transfers tokens from Token20Pay to DataToken owner,
// 2. transfers gas used to Carrier,
// 3. refunds gas to task sponsor.
// Settle returns hx.Hash, and error.
// gasUsedPrepay - carrier used gas for Prepay()

func (m *Token20PayManager) Settle(taskID *big.Int, gasUsedPrepay uint64) (common.Hash, error) {
	taskIDHex := hexutil.EncodeBig(taskID)

	if m.getPrivateKey() == nil {
		log.Errorf("cannot send Settle transaction cause organization wallet missing")
		return common.Hash{}, errors.New("organization private key is missing")
	}

	//估算gas
	gasEstimated, err := m.estimateGas("settle", taskID, new(big.Int).SetUint64(1))
	if err != nil {
		log.Errorf("failed to estimate gas for Token20Pay.Settle(), taskID: %s, error: %v", hexutil.EncodeBig(taskID), err)
		return common.Hash{}, errors.New("failed to estimate gas for Token20Pay.Settle()")
	}

	//交易参数的gasLimit可以放大，以防止交易执行gas不足；实际并不会真的消耗这么多
	opts, err := m.buildTxOpts(gasEstimated * 2)
	if err != nil {
		log.Errorf("failed to build transact options: %v", err)
	}

	//carrier付出的总的gas，gasUsedPrepay是准确的，estimatedGas是估计的
	totalGasUsed := gasUsedPrepay + gasEstimated

	//gas fee
	totalFeeUsed := new(big.Int).Mul(new(big.Int).SetUint64(totalGasUsed), opts.GasPrice)

	log.Debugf("call contract settle(),  taskID: %s, gasUsedPrepay: %d, gasEstimated: %d, gasLimit: %d, gasPrice: %d, totalFeeUsed: %d", taskIDHex, gasUsedPrepay, gasEstimated, opts.GasLimit, opts.GasPrice, totalFeeUsed)

	//合约
	tx, err := m.contractToken20PayInstance.Settle(opts, taskID, totalFeeUsed)
	if err != nil {
		log.Errorf("failed to call Token20Pay.Settle(), taskID: %s, error: %v", taskIDHex, err)
		return common.Hash{}, errors.New("failed to call Token20Pay.Settle()")
	}
	log.Debugf("call Token20Pay.Settle() txHash:%v, taskID:%s", tx.Hash().Hex(), taskIDHex)

	return tx.Hash(), nil
}

// GetReceipt returns the tx receipt. The caller goroutine will be blocked, and the caller could receive the receipt by channel.
func (m *Token20PayManager) GetReceipt(ctx context.Context, txHash common.Hash, period time.Duration) *ethereumtypes.Receipt {

	fetchReceipt := func(txHash common.Hash) (*ethereumtypes.Receipt, error) {
		receipt, err := m.client.TransactionReceipt(context.Background(), txHash)
		if nil != err {
			//including NotFound
			log.WithError(err).Warnf("Warning cannot query prepay transaction receipt, txHash: %s", txHash.Hex())
			return nil, err
		} else {
			log.Debugf("txHash:%s, receipt: %#v", txHash.Hex(), receipt)
			return receipt, nil
		}
	}

	if period < 0 { // do once only

		receipt, err := fetchReceipt(txHash)
		if nil != err {
			return nil
		}
		return receipt

	} else {
		ticker := time.NewTicker(period)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Warnf("query prepay transaction receipt timeout, txHash: %s", txHash.Hex())
				return nil
			case <-ticker.C:

				if receipt, err := fetchReceipt(txHash); nil == err {
					return receipt
				}
			}
		}
	}
}

// GetTaskState returns the task payment state.
// task state in contract
// constant int8 private NOTEXIST = -1;
// constant int8 private BEGIN = 0;
// constant int8 private PREPAY = 1;
// constant int8 private SETTLE = 2;
// constant int8 private END = 3;
func (m *Token20PayManager) GetTaskState(taskId *big.Int) (int, error) {
	log.Debugf("Start call contract taskState(), params{taskID: %s}", hexutil.EncodeBig(taskId))
	if state, err := m.contractToken20PayInstance.TaskState(&bind.CallOpts{}, taskId); err != nil {
		return -1, err
	} else {
		return int(state), nil
	}
}
