package tk

import (
	"context"
	"crypto/ecdsa"
	"errors"
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	"github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"strconv"
	"sync"
	"time"

	"github.com/bglmmz/chainclient"
	"github.com/datumtechs/datum-network-carrier/ach/tk/contracts"
	"github.com/datumtechs/datum-network-carrier/common/hexutil"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethereumtypes "github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"math/big"
)

const keystoreFile = ".keystore"

//"task:0x81050ea1c64ab0ed96e50b151c36bcef180eea18d3b3245e7c4a42aa08638c58"
var (
	datumPayAddress               = common.HexToAddress("0x263B1D39843BF2e1DA27d827e749992fbD1f1577")
	defaultDataTkPrepaymentAmount = big.NewInt(1e18)

	mockTaskID, _ = hexutil.DecodeBig("0x81050ea1c64ab0ed96e50b151c36bcef180eea18d3b3245e7c4a42aa08638c58")
)

type Config struct {
	URL           string `json:"chain"`
	WalletPwd     string
	privateKey    *ecdsa.PrivateKey
	walletAddress common.Address
}

type PayAgent struct {
	ethContext            *chainclient.EthContext
	abi                   abi.ABI
	tkPayContractInstance *contracts.DatumPay
	txSyncLocker          sync.Mutex
}

func NewPayAgent(ethContext *chainclient.EthContext) *PayAgent {
	log.Info("Init payAgent agent...")
	m := new(PayAgent)
	m.ethContext = ethContext

	instance, err := contracts.NewDatumPay(datumPayAddress, m.ethContext.GetClient())
	if err != nil {
		log.Fatal(err)
	}
	m.tkPayContractInstance = instance

	return m
}

func groupingTkList(tkItemList []*carrierapipb.TkItem) ([]common.Address, []*big.Int, []*carrierapipb.TkItem) {
	tkErc20AddressList := make([]common.Address, 0)
	tkErc20AmountList := make([]*big.Int, 0)
	tkErc721ItemList := make([]*carrierapipb.TkItem, 0)

	for _, item := range tkItemList {
		if item.TkType == constant.TkType_Tk20 {
			tkErc20AddressList = append(tkErc20AddressList, common.HexToAddress(item.TkAddress))
			if len(item.Value) == 0 {
				tkErc20AmountList = append(tkErc20AmountList, defaultDataTkPrepaymentAmount)
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
	log.Debugf("call EstimateTaskGas, sponsorAddress: %s, tknAddressList:%v", taskSponsorAddress, tkItemList)

	tk20AddressList, tk20AmountList, _ := groupingTkList(tkItemList)

	// Estimating task gas happens before task starting, and the task ID has not been generated at this moment, so, apply a mock task ID.
	input := m.buildInput("prepay", mockTaskID, common.HexToAddress(taskSponsorAddress), big.NewInt(1), tk20AddressList, tk20AmountList)
	gasLimit, err := m.ethContext.EstimateGas(context.Background(), datumPayAddress, input)
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

// Prepay transfers more than enough gas from task sponsor to DatumPay, this gas will payAgent carrier to call Prepay()/Settle(), and remaining gas will refund to task sponsor.
// Prepay returns hx.Hash, and error.
// The complete procedure consists of two calls to DatumPay, the first is Prepay, the second is Settle.
func (m *PayAgent) Prepay(taskID *big.Int, taskSponsorAccount common.Address, tkItemList []*carrierapipb.TkItem) (common.Hash, error) {
	m.txSyncLocker.Lock()
	defer m.txSyncLocker.Unlock()

	taskIDHex := hexutil.EncodeBig(taskID)
	if WalletManagerInstance().GetPrivateKey() == nil {
		log.Errorf("cannot send Prepay transaction cause organization wallet missing")
		return common.Hash{}, errors.New("organization private key is missing")
	}

	tkErc20AddressList, tkErc20AmountList, tkErc721ItemList := groupingTkList(tkItemList)
	for _, amt := range tkErc20AmountList {
		if amt == nil || amt.Int64() == 0 {
			amt = defaultDataTkPrepaymentAmount
		}
	}

	for _, tkErc721 := range tkErc721ItemList {
		if err := m.inspectTkErc721ExtInfo(taskSponsorAccount, tkErc721); err != nil {
			return common.Hash{}, err
		}
	}

	input := m.buildInput("prepay", mockTaskID, taskSponsorAccount, big.NewInt(1), tkErc20AddressList, tkErc20AmountList)

	//估算gas
	gasEstimated, err := m.ethContext.EstimateGas(context.Background(), datumPayAddress, input)
	if err != nil {
		log.Errorf("failed to estimate gas for DatumPay.Prepay(), taskID: %s, error: %v", taskIDHex, err)
		return common.Hash{}, errors.New("failed to estimate gas for DatumPay.Prepay()")
	}

	//交易参数直接使用用户预付的总的gas，尽量放大，以防止交易执行gas不足
	gasEstimated = uint64(float64(gasEstimated) * 1.30)
	opts, err := m.ethContext.BuildTxOpts(0, gasEstimated)
	if err != nil {
		log.Errorf("failed to build transact options to call DatumPay.Prepay(): %v", err)
		return common.Hash{}, errors.New("failed to build transact options to call DatumPay.Prepay()")
	}

	//gas fee, 支付助手合约，需要记录用户预付的手续费
	feePrepaid := new(big.Int).Mul(new(big.Int).SetUint64(gasEstimated), opts.GasPrice)

	log.Debugf("Start call contract prepay(), taskID: %s, nonce: %d, gasEstimated: %d, gasLimit: %d, gasPrice: %d, feePrepaid: %d, taskSponsorAccount: %s, dataTokenAddressList: %v, dataTokenAmountList: %v",
		taskIDHex, opts.Nonce, gasEstimated, opts.GasLimit, opts.GasPrice, feePrepaid, taskSponsorAccount.String(), tkErc20AddressList, tkErc20AmountList)

	tx, err := m.tkPayContractInstance.Prepay(opts, taskID, taskSponsorAccount, feePrepaid, tkErc20AddressList, tkErc20AmountList)
	if err != nil {
		log.WithError(err).Errorf("failed to call DatumPay.Prepay(), taskID: %s", taskIDHex)
		return common.Hash{}, errors.New("failed to call DatumPay.Prepay()")
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

func (m *PayAgent) Settle(taskID *big.Int, gasUsedPrepay uint64) (common.Hash, error) {
	m.txSyncLocker.Lock()
	defer m.txSyncLocker.Unlock()

	taskIDHex := hexutil.EncodeBig(taskID)

	if WalletManagerInstance().GetPrivateKey() == nil {
		log.Errorf("cannot send Settle transaction cause organization wallet missing")
		return common.Hash{}, errors.New("organization private key is missing")
	}

	input := m.buildInput("settle", taskID, big.NewInt(1), new(big.Int).SetUint64(1))

	//估算gas
	gasEstimated, err := m.ethContext.EstimateGas(context.Background(), datumPayAddress, input)
	if err != nil {
		log.Errorf("failed to estimate gas for DatumPay.Settle(), taskID: %s, error: %v", hexutil.EncodeBig(taskID), err)
		return common.Hash{}, errors.New("failed to estimate gas for DatumPay.Settle()")
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
	tx, err := m.tkPayContractInstance.Settle(opts, taskID, totalFeeUsed)
	if err != nil {
		log.Errorf("failed to call DatumPay.Settle(), taskID: %s, error: %v", taskIDHex, err)
		return common.Hash{}, errors.New("failed to call DatumPay.Settle()")
	}
	log.Debugf("call DatumPay.Settle() txHash:%v, taskID:%s", tx.Hash().Hex(), taskIDHex)

	return tx.Hash(), nil
}

// GetReceipt returns the tx receipt. The caller goroutine will be blocked, and the caller could receive the receipt by channel.
func (m *PayAgent) GetReceipt(ctx context.Context, txHash common.Hash, interval time.Duration) *ethereumtypes.Receipt {
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
	if state, err := m.tkPayContractInstance.TaskState(&bind.CallOpts{}, taskId); err != nil {
		return -1, err
	} else {
		return int(state), nil
	}
}
