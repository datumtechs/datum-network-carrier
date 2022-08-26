package tk

import (
	"context"
	"crypto/ecdsa"
	"errors"
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	"github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bglmmz/chainclient"
	"github.com/datumtechs/datum-network-carrier/ach/tk/contracts"
	"github.com/datumtechs/datum-network-carrier/common/hexutil"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"math/big"
)

const keystoreFile = ".keystore"

//"task:0x81050ea1c64ab0ed96e50b151c36bcef180eea18d3b3245e7c4a42aa08638c58"
var (
	//payAgentAddress           = ethcommon.HexToAddress("0x263B1D39843BF2e1DA27d827e749992fbD1f1577")
	defaultTkPrepaymentAmount = big.NewInt(1e18)

	mockTaskID, _ = hexutil.DecodeBig("0x81050ea1c64ab0ed96e50b151c36bcef180eea18d3b3245e7c4a42aa08638c58")
)

type Config struct {
	URL           string `json:"chain"`
	WalletPwd     string
	privateKey    *ecdsa.PrivateKey
	walletAddress ethcommon.Address
}

type PayAgent struct {
	ethContext               *chainclient.EthContext
	abi                      abi.ABI
	payAgentContractInstance *contracts.DatumPay
	payAgentContractProxy    ethcommon.Address
	txSyncLocker             sync.Mutex
}

func NewPayAgent(ethContext *chainclient.EthContext, payAgentContractProxy ethcommon.Address) *PayAgent {
	log.Infof("Init pay agent, wallet address: %s...", ethContext.GetAddress())
	m := new(PayAgent)
	m.ethContext = ethContext
	m.payAgentContractProxy = payAgentContractProxy

	instance, err := contracts.NewDatumPay(m.payAgentContractProxy, m.ethContext.GetClient())
	if err != nil {
		log.Fatal(err)
	}
	m.payAgentContractInstance = instance

	abiCode, err := abi.JSON(strings.NewReader(contracts.DatumPayMetaData.ABI))
	if err != nil {
		log.Fatal(err)
	}
	m.abi = abiCode

	return m
}

func groupingTkList(tkItemList []*carrierapipb.TkItem) ([]ethcommon.Address, []*big.Int, []*carrierapipb.TkItem) {
	tkErc20AddressList := make([]ethcommon.Address, 0)
	tkErc20AmountList := make([]*big.Int, 0)
	tkErc721ItemList := make([]*carrierapipb.TkItem, 0)

	for _, item := range tkItemList {
		if item.TkType == constant.TkType_Tk20 {
			tkErc20AddressList = append(tkErc20AddressList, ethcommon.HexToAddress(item.TkAddress))
			if len(item.Value) == 0 {
				tkErc20AmountList = append(tkErc20AmountList, defaultTkPrepaymentAmount)
			} else {
				vInt, err := strconv.ParseUint(item.Value, 10, 64)
				if err != nil {
					log.WithError(err).Errorf("token value is not a valid integer, it will be reset to the default value of 1. value:=%s", item.Value)
					vInt = 1
				}
				tkErc20AmountList = append(tkErc20AmountList, new(big.Int).SetUint64(vInt))
			}

		} else if item.TkType == constant.TkType_Tk721 {
			tkErc721ItemList = append(tkErc721ItemList, item)
		}
	}
	return tkErc20AddressList, tkErc20AmountList, tkErc721ItemList
}

// EstimateTaskGas estimates gas fee for a task's sponsor.
// EstimateTaskGas returns estimated gas and suggested gas price.
func (m *PayAgent) EstimateTaskGas(taskSponsorAddress string, tkItemList []*carrierapipb.TkItem) (uint64, *big.Int, error) {
	log.Debugf("call EstimateTaskGas, sender: %s, sponsorAddress: %s, tknAddressList:%v", m.ethContext.GetAddress().Hex(), taskSponsorAddress, tkItemList)

	tk20AddressList, tk20AmountList, _ := groupingTkList(tkItemList)

	// Estimating task gas happens before task starting, and the task ID has not been generated at this moment, so, apply a mock task ID.
	input := m.buildInput("prepay", mockTaskID, ethcommon.HexToAddress(taskSponsorAddress), big.NewInt(1), tk20AddressList, tk20AmountList)
	gasLimit, err := m.ethContext.EstimateGas(context.Background(), m.payAgentContractProxy, input)
	if err != nil {
		log.Errorf("call EstimateTaskGas error: %v", err)
		return 0, nil, err
	}

	gasPrice, err := m.ethContext.SuggestGasPrice(context.Background())
	if err != nil {
		log.Errorf("call SuggestGasPrice error: %v", err)
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
	// todo: assure gas limit of settle = gas limit of prepay
	return gasLimit + gasLimit, gasPrice, nil
}

// estimateGas call EstimateGas() by RPC
func (m *PayAgent) buildInput(method string, params ...interface{}) []byte {
	input, err := m.abi.Pack(method, params...)
	if err != nil {
		log.WithError(err).Errorf("failed to build input for method: %s", method)
		return nil
	}
	return input
}

// VerifyNFT verify each NTF
func (m *PayAgent) VerifyTk721(taskSponsorAccount ethcommon.Address, tk721ItemList []*carrierapipb.TkItem) error {
	m.txSyncLocker.Lock()
	defer m.txSyncLocker.Unlock()

	for _, tkErc721 := range tk721ItemList {
		if err := m.inspectTk721ExtInfo(taskSponsorAccount, tkErc721); err != nil {
			return err
		}
	}
	return nil
}

// PrepayTk20 transfers more than enough gas from task sponsor to DatumPay, this gas will payAgent carrier to call PrepayTk20()/Settle(), and remaining gas will refund to task sponsor.
// PrepayTk20 returns hx.Hash, and error.
// The complete procedure consists of two calls to DatumPay, the first is PrepayTk20, the second is Settle.
func (m *PayAgent) PrepayTk20(taskID *big.Int, taskSponsorAccount ethcommon.Address, tk20ItemList []*carrierapipb.TkItem) (ethcommon.Hash, error) {
	m.txSyncLocker.Lock()
	defer m.txSyncLocker.Unlock()

	taskIDHex := hexutil.EncodeBig(taskID)
	if WalletManagerInstance().GetPrivateKey() == nil {
		log.Errorf("cannot send Prepay transaction cause organization wallet missing")
		return ethcommon.Hash{}, errors.New("organization private key is missing")
	}

	//tk20AddressList, tk20AmountList, tk721ItemList := groupingTkList(tk20ItemList)
	tk20AddressList, tk20AmountList, _ := groupingTkList(tk20ItemList)
	for _, amt := range tk20AmountList {
		if amt == nil || amt.Int64() == 0 {
			amt = defaultTkPrepaymentAmount
		}
	}

	/*for _, tkErc721 := range tk721ItemList {
		if err := m.inspectTk721ExtInfo(taskSponsorAccount, tkErc721); err != nil {
			return common.Hash{}, err
		}
	}*/

	input := m.buildInput("prepay", mockTaskID, taskSponsorAccount, big.NewInt(1), tk20AddressList, tk20AmountList)

	//估算gas
	gasEstimated, err := m.ethContext.EstimateGas(context.Background(), m.payAgentContractProxy, input)
	if err != nil {
		log.Errorf("failed to estimate gas for DatumPay.Prepay(), taskID: %s, error: %v", taskIDHex, err)
		return ethcommon.Hash{}, errors.New("failed to estimate gas for DatumPay.Prepay()")
	}

	//交易参数直接使用用户预付的总的gas，尽量放大，以防止交易执行gas不足
	gasEstimated = uint64(float64(gasEstimated) * 1.30)
	opts, err := m.ethContext.BuildTxOpts(0, gasEstimated)
	if err != nil {
		log.Errorf("failed to build transact options to call DatumPay.Prepay(): %v", err)
		return ethcommon.Hash{}, errors.New("failed to build transact options to call DatumPay.Prepay()")
	}

	//gas fee, 支付助手合约，需要记录用户预付的手续费
	feePrepaid := new(big.Int).Mul(new(big.Int).SetUint64(gasEstimated), opts.GasPrice)

	log.Debugf("Start call contract prepay(), taskID: %s, nonce: %d, gasEstimated: %d, gasLimit: %d, gasPrice: %d, feePrepaid: %d, taskSponsorAccount: %s, dataTokenAddressList: %v, dataTokenAmountList: %v",
		taskIDHex, opts.Nonce, gasEstimated, opts.GasLimit, opts.GasPrice, feePrepaid, taskSponsorAccount.String(), tk20AddressList, tk20AmountList)

	tx, err := m.payAgentContractInstance.Prepay(opts, taskID, taskSponsorAccount, feePrepaid, tk20AddressList, tk20AmountList)
	if err != nil {
		log.WithError(err).Errorf("failed to call DatumPay.Prepay(), taskID: %s", taskIDHex)
		return ethcommon.Hash{}, errors.New("failed to call DatumPay.Prepay()")
	}
	log.Debugf("call DatumPay.Prepay() txHash:%v, taskID:%s", tx.Hash().Hex(), taskIDHex)

	return tx.Hash(), nil
}

// Settle get funds clearing,
// 1. transfers tks from DatumPay to DataToken owner,
// 2. transfers gas used to Carrier,
// 3. refunds gas to task sponsor.
// Settle returns hx.Hash, and error.
// gasUsedPrepay - carrier used gas for Prepay()

func (m *PayAgent) Settle(taskID *big.Int, gasUsedPrepay uint64) (ethcommon.Hash, error) {
	m.txSyncLocker.Lock()
	defer m.txSyncLocker.Unlock()

	taskIDHex := hexutil.EncodeBig(taskID)

	if WalletManagerInstance().GetPrivateKey() == nil {
		log.Errorf("cannot send Settle transaction cause organization wallet missing")
		return ethcommon.Hash{}, errors.New("organization private key is missing")
	}

	input := m.buildInput("settle", taskID, big.NewInt(1), new(big.Int).SetUint64(1))

	//估算gas
	gasEstimated, err := m.ethContext.EstimateGas(context.Background(), m.payAgentContractProxy, input)
	if err != nil {
		log.Errorf("failed to estimate gas for DatumPay.Settle(), taskID: %s, error: %v", hexutil.EncodeBig(taskID), err)
		return ethcommon.Hash{}, errors.New("failed to estimate gas for DatumPay.Settle()")
	}

	//交易参数的gasLimit可以放大，以防止交易执行gas不足；实际并不会真的消耗这么多
	gasEstimated = uint64(float64(gasEstimated) * 1.30)
	opts, err := m.ethContext.BuildTxOpts(0, gasEstimated)
	if err != nil {
		log.Errorf("failed to build transact options: %v", err)
	}

	//carrier付出的总的gas，gasUsedPrepay是准确的，estimatedGas是估计的
	totalGasUsed := gasUsedPrepay + gasEstimated

	//gas fee
	totalFeeUsed := new(big.Int).Mul(new(big.Int).SetUint64(totalGasUsed), opts.GasPrice)

	log.Debugf("call contract settle(),  taskID: %s, nonce: %d, gasUsedPrepay: %d, gasEstimated: %d, gasLimit: %d, gasPrice: %d, totalFeeUsed: %d", taskIDHex, opts.Nonce, gasUsedPrepay, gasEstimated, opts.GasLimit, opts.GasPrice, totalFeeUsed)

	//合约
	tx, err := m.payAgentContractInstance.Settle(opts, taskID, totalFeeUsed)
	if err != nil {
		log.Errorf("failed to call DatumPay.Settle(), taskID: %s, error: %v", taskIDHex, err)
		return ethcommon.Hash{}, errors.New("failed to call DatumPay.Settle()")
	}
	log.Debugf("call DatumPay.Settle() txHash:%v, taskID:%s", tx.Hash().Hex(), taskIDHex)

	return tx.Hash(), nil
}

// GetReceipt returns the tx receipt. The caller goroutine will be blocked, and the caller could receive the receipt by channel.
func (m *PayAgent) GetReceipt(ctx context.Context, txHash ethcommon.Hash, interval time.Duration) *ethtypes.Receipt {
	return m.ethContext.WaitReceipt(ctx, txHash, interval)
}

// GetTaskState returns the task payment state.
// task state in contract
// constant int8 private NOTEXIST = -1;
// constant int8 private BEGIN = 0;
// constant int8 private PREPAY = 1;
// constant int8 private SETTLE = 2;
// constant int8 private END = 3;
func (m *PayAgent) GetTaskState(taskId *big.Int) (int, error) {
	log.Debugf("Start call contract taskState(), params{taskID: %s}", hexutil.EncodeBig(taskId))
	if state, err := m.payAgentContractInstance.TaskState(&bind.CallOpts{}, taskId); err != nil {
		return -1, err
	} else {
		return int(state), nil
	}
}
