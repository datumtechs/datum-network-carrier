package task

import (
	"context"
	"encoding/json"
	"fmt"
	carriercommon "github.com/datumtechs/datum-network-carrier/common"
	"github.com/datumtechs/datum-network-carrier/common/hexutil"
	"github.com/datumtechs/datum-network-carrier/common/runutil"
	"github.com/datumtechs/datum-network-carrier/common/signutil"
	"github.com/datumtechs/datum-network-carrier/common/timeutils"
	"github.com/datumtechs/datum-network-carrier/common/traceutil"
	ctypes "github.com/datumtechs/datum-network-carrier/consensus/twopc/types"
	ev "github.com/datumtechs/datum-network-carrier/core/evengine"
	"github.com/datumtechs/datum-network-carrier/core/rawdb"
	"github.com/datumtechs/datum-network-carrier/core/resource"
	"github.com/datumtechs/datum-network-carrier/core/schedule"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	ethereumcommon "github.com/ethereum/go-ethereum/common"
	"github.com/prysmaticlabs/prysm/shared/hashutil"
	"math/big"

	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	carriernetmsgcommonpb "github.com/datumtechs/datum-network-carrier/pb/carrier/netmsg/common"
	carriertwopcpb "github.com/datumtechs/datum-network-carrier/pb/carrier/netmsg/consensus/twopc"
	carriernetmsgtaskmngpb "github.com/datumtechs/datum-network-carrier/pb/carrier/netmsg/taskmng"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	fighterpbtypes "github.com/datumtechs/datum-network-carrier/pb/fighter/types"
	"github.com/datumtechs/datum-network-carrier/policy"
	"github.com/datumtechs/datum-network-carrier/types"
	"github.com/libp2p/go-libp2p-core/peer"
	"strconv"
	"strings"
	"time"
)

func (m *Manager) tryScheduleTask() error {

	nonConsTask, taskId, err := m.scheduler.TrySchedule()
	if nil == err && nil == nonConsTask {
		return nil
	} else if nil != err && err == schedule.ErrAbandonTaskWithNotFoundTask {
		m.scheduler.RemoveTask(taskId)
		log.WithError(err).Errorf("Failed to call scheduler.TrySchedule(), then remove bullet task, taskId: {%s}", taskId)
	} else if nil != err && err == schedule.ErrAbandonTaskWithNotFoundPowerPartyIds {
		m.scheduler.RemoveTask(taskId)
		log.WithError(err).Errorf("Failed to call scheduler.TrySchedule(), then remove bullet task, taskId: {%s}", taskId)
		m.sendNeedExecuteTaskByAction(types.NewNeedExecuteTask(
			"",
			commonconstantpb.TaskRole_TaskRole_Sender,
			commonconstantpb.TaskRole_TaskRole_Sender,
			nonConsTask.GetTask().GetTaskSender(),
			nonConsTask.GetTask().GetTaskSender(),
			nonConsTask.GetTask().GetTaskId(),
			types.TaskScheduleFailed,
			&types.PrepareVoteResource{},          // zero value
			&carriertwopcpb.ConfirmTaskPeerInfo{}, // zero value
			fmt.Errorf("schedule failed: "+schedule.ErrAbandonTaskWithNotFoundPowerPartyIds.Error()),
		))
		return err
	} else if nil != err {
		if nil != nonConsTask {

			m.resourceMng.GetDB().StoreTaskEvent(m.eventEngine.GenerateEvent(ev.TaskScheduleFailed.GetType(),
				nonConsTask.GetTask().GetTaskId(), nonConsTask.GetTask().GetTaskSender().GetIdentityId(),
				nonConsTask.GetTask().GetTaskSender().GetPartyId(), "schedule failed: "+err.Error()))

			if er := m.scheduler.RepushTask(nonConsTask.GetTask()); er == schedule.ErrRescheduleLargeThreshold {
				log.WithError(err).Errorf("Failed to repush local task into queue/starve queue after trySchedule failed %s on `taskManager.tryScheduleTask()`, taskId: {%s}",
					err, nonConsTask.GetTask().GetTaskId())

				m.scheduler.RemoveTask(nonConsTask.GetTask().GetTaskId())
				m.sendNeedExecuteTaskByAction(types.NewNeedExecuteTask(
					"",
					commonconstantpb.TaskRole_TaskRole_Sender,
					commonconstantpb.TaskRole_TaskRole_Sender,
					nonConsTask.GetTask().GetTaskSender(),
					nonConsTask.GetTask().GetTaskSender(),
					nonConsTask.GetTask().GetTaskId(),
					types.TaskScheduleFailed,
					&types.PrepareVoteResource{},          // zero value
					&carriertwopcpb.ConfirmTaskPeerInfo{}, // zero value
					fmt.Errorf("schedule failed: "+err.Error()+" and "+schedule.ErrRescheduleLargeThreshold.Error()),
				))
			}
		}
		return err
	}

	go func(nonConsTask *types.NeedConsensusTask) {

		log.Debugf("Start send `NEED-CONSENSUS` task to [consensus engine] on `taskManager.tryScheduleTask()`, taskId: {%s}", nonConsTask.GetTask().GetTaskId())

		if err := m.consensusEngine.OnPrepare(nonConsTask); nil != err {
			log.WithError(err).Errorf("Failed to call `OnPrepare()` of [consensus engine] on `taskManager.tryScheduleTask()`, taskId: {%s}", nonConsTask.GetTask().GetTaskId())
			// re push task into queue ,if anything else
			if err := m.scheduler.RepushTask(nonConsTask.GetTask()); err == schedule.ErrRescheduleLargeThreshold {
				log.WithError(err).Errorf("Failed to repush local task into queue/starve queue after call `OnPrepare()` of [consensus engine] on `taskManager.tryScheduleTask()`, taskId: {%s}",
					nonConsTask.GetTask().GetTaskId())

				m.scheduler.RemoveTask(nonConsTask.GetTask().GetTaskId())
				m.sendNeedExecuteTaskByAction(types.NewNeedExecuteTask(
					"",
					commonconstantpb.TaskRole_TaskRole_Sender,
					commonconstantpb.TaskRole_TaskRole_Sender,
					nonConsTask.GetTask().GetTaskSender(),
					nonConsTask.GetTask().GetTaskSender(),
					nonConsTask.GetTask().GetTaskId(),
					types.TaskScheduleFailed,
					&types.PrepareVoteResource{},          // zero value
					&carriertwopcpb.ConfirmTaskPeerInfo{}, // zero value
					fmt.Errorf("consensus onPrepare failed: "+err.Error()+" and "+schedule.ErrRescheduleLargeThreshold.Error()),
				))
			} else {
				log.Debugf("Succeed to repush local task into queue/starve queue after call `OnPrepare()` of [consensus engine] on `taskManager.tryScheduleTask()`, taskId: {%s}",
					nonConsTask.GetTask().GetTaskId())
			}
			return
		}
		if err := m.consensusEngine.OnHandle(nonConsTask); nil != err {
			log.WithError(err).Errorf("Failed to call `OnHandle()` of [consensus engine] on `taskManager.tryScheduleTask()`, taskId: {%s}", nonConsTask.GetTask().GetTaskId())
		}
	}(nonConsTask)
	return nil
}

func (m *Manager) beginConsumeMetadataOrPower(task *types.NeedExecuteTask, localTask *types.Task) error {

	switch m.config.MetadataConsumeOption {
	case 1: // use metadataAuth
		return m.beginConsumeByMetadataAuth(task, localTask)
	case 2: // use datatoken
		return m.beginConsumeByDataToken(task, localTask)
	default: // use nothing
		return nil
	}
}
func (m *Manager) beginConsumeByMetadataAuth(task *types.NeedExecuteTask, localTask *types.Task) error {

	partyId := task.GetLocalTaskOrganization().GetPartyId()

	switch task.GetLocalTaskRole() {
	case commonconstantpb.TaskRole_TaskRole_Sender:
		return nil // do nothing ...
	case commonconstantpb.TaskRole_TaskRole_DataSupplier:
		for _, dataSupplier := range localTask.GetTaskData().GetDataSuppliers() {
			if partyId == dataSupplier.GetPartyId() {

				userType := localTask.GetTaskData().GetUserType()
				user := localTask.GetTaskData().GetUser()

				metadataId, err := policy.FetchMetedataIdByPartyIdFromDataPolicy(partyId, localTask.GetTaskData().GetDataPolicyTypes(), localTask.GetTaskData().GetDataPolicyOptions())
				if nil != err {
					return fmt.Errorf("not fetch metadataId from task dataPolicy when call beginConsumeByMetadataAuth(), %s, taskId: {%s}, partyId: {%s}",
						err, localTask.GetTaskId(), partyId)
				}
				// verify metadataAuth first
				if err := m.authMng.VerifyMetadataAuth(userType, user, metadataId); nil != err {
					return fmt.Errorf("verify user metadataAuth failed when call beginConsumeByMetadataAuth(), %s, userType: {%s}, user: {%s}, taskId: {%s}, partyId: {%s}, metadataId: {%s}",
						err, userType, user, localTask.GetTaskId(), partyId, metadataId)
				}

				internalMetadataFlag, err := m.resourceMng.GetDB().IsInternalMetadataById(metadataId)
				if nil != err {
					return fmt.Errorf("check metadata whether internal metadata failed %s when call beginConsumeByMetadataAuth(), taskId: {%s}, partyId: {%s}, metadataId: {%s}",
						err, localTask.GetTaskId(), partyId, metadataId)
				}

				// only consume metadata auth when metadata is not internal metadata.
				if !internalMetadataFlag {
					// query metadataAuthId by metadataId
					metadataAuthId, err := m.authMng.QueryMetadataAuthIdByMetadataId(userType, user, metadataId)
					if nil != err {
						return fmt.Errorf("query metadataAuthId failed %s when call beginConsumeByMetadataAuth(), metadataId: {%s}", err, metadataId)
					}
					// ConsumeMetadataAuthority
					if err = m.authMng.ConsumeMetadataAuthority(metadataAuthId); nil != err {
						return fmt.Errorf("consume metadataAuth failed %s when call beginConsumeByMetadataAuth(), metadataAuthId: {%s}", err, metadataAuthId)
					} else {
						log.Debugf("Succeed consume metadataAuth when call beginConsumeByMetadataAuth(), taskId: {%s}, metadataAuthId: {%s}", task.GetTaskId(), metadataAuthId)
					}
				}
				break
			}
		}
		return nil
	case commonconstantpb.TaskRole_TaskRole_PowerSupplier:
		return nil // do nothing ...
	case commonconstantpb.TaskRole_TaskRole_Receiver:
		return nil // do nothing ...
	default:
		return fmt.Errorf("unknown task role on beginConsumeByMetadataAuth()")

	}
}
func (m *Manager) beginConsumeByDataToken(task *types.NeedExecuteTask, localTask *types.Task) error {

	partyId := task.GetLocalTaskOrganization().GetPartyId()

	switch task.GetLocalTaskRole() {
	case commonconstantpb.TaskRole_TaskRole_Sender:

		if partyId != localTask.GetTaskSender().GetPartyId() {
			return fmt.Errorf("this partyId is not task sender on beginConsumeByDataToken()")
		}

		taskIdBigInt, err := hexutil.DecodeBig("0x" + strings.TrimLeft(strings.Trim(task.GetTaskId(), types.PREFIX_TASK_ID+"0x"), "\x00"))
		if nil != err {
			return fmt.Errorf("cannot decode taskId to big.Int on beginConsumeByDataToken(), %s", err)
		}

		// verify user
		/**
		  User_1 = 1;    // PlatON
		  User_2 = 2;    // Alaya
		  User_3 = 3;    // Ethereum
		*/

		// fetch all datatoken contract adresses of metadata of task
		metadataIds, err := policy.FetchAllMetedataIdsFromDataPolicy(localTask.GetTaskData().GetDataPolicyTypes(), localTask.GetTaskData().GetDataPolicyOptions())
		if nil != err {
			return fmt.Errorf("cannot fetch all metadataIds of dataPolicyOption on beginConsumeByDataToken(), %s", err)
		}
		metadataList, err := m.resourceMng.GetDB().QueryMetadataByIds(metadataIds)
		if nil != err {
			return fmt.Errorf("call QueryMetadataByIds() failed on beginConsumeByDataToken(), %s", err)
		}

		// for debug log...
		addrs := make([]string, len(metadataList))

		dataTokenAaddresses := make([]ethereumcommon.Address, len(metadataList))
		for i, metadata := range metadataList {
			addr := ethereumcommon.HexToAddress(metadata.GetData().GetTokenAddress())
			dataTokenAaddresses[i] = addr
			addrs[i] = addr.String()
		}

		user := ethereumcommon.HexToAddress(localTask.GetTaskData().GetUser())

		log.Debugf("Start call token20PayManager.prepay(), taskId: {%s}, partyId: {%s}, call params{taskIdBigInt: %d, taskSponsorAccount: %s, dataTokenAaddresses: %s}",
			task.GetTaskId(), partyId, taskIdBigInt, user.String(), "["+strings.Join(addrs, ",")+"]")

		// start prepay dataToken
		txHash, err := m.token20PayMng.Prepay(taskIdBigInt, user, dataTokenAaddresses)
		if nil != err {
			return fmt.Errorf("cannot call token20Pay to prepay datatoken on beginConsumeByDataToken(), %s", err)
		}

		log.Debugf("Succeed send `contract prepay()` tx to blockchain on beginConsumeByDataToken(), taskId: {%s}, partyId: {%s}, taskIdBigInt: {%d}, txHash: {%s}",
			task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), taskIdBigInt, txHash.String())

		// make sure the `prepay` tx into blockchain
		timeout := time.Duration(localTask.GetTaskData().GetOperationCost().GetDuration()) * time.Millisecond
		ctx, cancelFn := context.WithTimeout(context.Background(), timeout)
		//ctx, cancelFn := context.WithCancel(context.Background())
		defer cancelFn()

		receipt := m.token20PayMng.GetReceipt(ctx, txHash, time.Duration(500)*time.Millisecond) // period 500 ms
		if nil == receipt {
			return fmt.Errorf("prepay dataToken failed, the transaction had not receipt on beginConsumeByDataToken(), txHash: {%s}", txHash.String())
		}
		// contract tx execute failed.
		if receipt.Status == 0 {
			return fmt.Errorf("prepay dataToken failed, the transaction receipt status is %d on beginConsumeByDataToken(), txHash: {%s}", receipt.Status, txHash.String())
		}

		// query task state
		state, err := m.token20PayMng.GetTaskState(taskIdBigInt)
		if nil != err {
			//including NotFound
			return fmt.Errorf("query task state of token20Pay failed, %s on beginConsumeByDataToken()", err)
		}
		// task state in contract
		// constant int8 private NOTEXIST = -1;
		// constant int8 private BEGIN = 0;
		// constant int8 private PREPAY = 1;
		// constant int8 private SETTLE = 2;
		// constant int8 private END = 3;
		if state == -1 { //  We need to know if the task status value is 1.
			return fmt.Errorf("task state is not existing in token20Pay contract on beginConsumeByDataToken()")
		}

		log.Debugf("Succeed execute `contract prepay()` tx on blockchain on beginConsumeByDataToken(), taskId: {%s}, partyId: {%s}, taskIdBigInt: {%d}, txHash: {%s}, task.state: {%d}",
			task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), taskIdBigInt, txHash.String(), state)

		// update consumeSpec into needExecuteTask
		if "" == strings.Trim(task.GetConsumeSpec(), "") {
			return fmt.Errorf("consumeSpec about task is empty on beginConsumeByDataToken(), consumeSpec: %s", task.GetConsumeSpec())
		}
		var consumeSpec *types.DatatokenPaySpec
		if err := json.Unmarshal([]byte(task.GetConsumeSpec()), &consumeSpec); nil != err {
			return fmt.Errorf("cannot json unmarshal consumeSpec on beginConsumeByDataToken(), consumeSpec: %s, %s", task.GetConsumeSpec(), err)
		}

		// task state in contract
		// constant int8 private NOTEXIST = -1;
		// constant int8 private BEGIN = 0;
		// constant int8 private PREPAY = 1;
		// constant int8 private SETTLE = 2;
		// constant int8 private END = 3;
		consumeSpec.Consumed = int32(state)
		//consumeSpec.GasEstimated = gasLimit
		consumeSpec.GasUsed = receipt.GasUsed

		b, err := json.Marshal(consumeSpec)
		if nil != err {
			return fmt.Errorf("connot json marshal task consumeSpec on beginConsumeByDataToken(), consumeSpec: %v, %s", consumeSpec, err)
		}
		task.SetConsumeSpec(string(b))
		// update needExecuteTask into cache
		m.updateNeedExecuteTaskCache(task)

		return nil
	case commonconstantpb.TaskRole_TaskRole_DataSupplier, commonconstantpb.TaskRole_TaskRole_PowerSupplier, commonconstantpb.TaskRole_TaskRole_Receiver:

		taskIdBigInt, err := hexutil.DecodeBig("0x" + strings.TrimLeft(strings.Trim(task.GetTaskId(), types.PREFIX_TASK_ID+"0x"), "\x00"))
		if nil != err {
			return fmt.Errorf("cannot decode taskId to big.Int on beginConsumeByDataToken(), %s", err)
		}

		// make sure the `prepay` tx of task sender into blockchain
		timeout := time.Duration(localTask.GetTaskData().GetOperationCost().GetDuration()) * time.Millisecond
		ctx, cancelFn := context.WithTimeout(context.Background(), timeout)
		//ctx, cancelFn := context.WithCancel(context.Background())
		defer cancelFn()

		queryTaskState := func(ctx context.Context, taskIdBigInt *big.Int, period time.Duration) (int, error) {

			start := timeutils.UnixMsec()

			ticker := time.NewTicker(period)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():

					end := timeutils.UnixMsec()
					// time.Unix(end/1000, 0)
					log.Warnf("Warning query task state of token20Pay time out on blockchain on beginConsumeByDataToken(), taskId: {%s}, partyId: {%s}, taskIdBigInt: {%d}, startTime: {%d <==> %s}, endTime: {%d <==> %s}, duration: %d ms",
						task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), taskIdBigInt, start, time.Unix(start/1000, 0).Format("2006-01-02 15:04:05"), end, time.Unix(end/1000, 0).Format("2006-01-02 15:04:05"), end-start)

					return 0, fmt.Errorf("query task state of token20Pay time out")
				case <-ticker.C:
					state, err := m.token20PayMng.GetTaskState(taskIdBigInt)
					if nil != err {
						//including NotFound
						log.WithError(err).Warnf("Warning cannot query task state of token20Pay on beginConsumeByDataToken(), taskId: {%s}, partyId: {%s}, taskIdBigInt: {%d}",
							task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), taskIdBigInt)
						continue
					}
					// task state in contract
					// constant int8 private NOTEXIST = -1;
					// constant int8 private BEGIN = 0;
					// constant int8 private PREPAY = 1;
					// constant int8 private SETTLE = 2;
					// constant int8 private END = 3;
					if state == -1 { //  We need to know if the task status value is 1.
						//log.Warnf("query task state value is equal `NOTEXIST` on beginConsumeByDataToken(), taskId: {%s}, partyId: {%s}, taskId.bigInt: {%d}, task.state: {%d}",
						//	task.GetTaskId(), partyId, taskId.Uint64(), state)
						continue
					}
					log.Debugf("Succeed query task.state value is not equal `NOTEXIST` on blockchain on beginConsumeByDataToken(), taskId: {%s}, partyId: {%s}, taskIdBigInt: {%d}, task.state: {%d}",
						task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), taskIdBigInt, state)
					return state, nil
				}
			}
		}

		state, err := queryTaskState(ctx, taskIdBigInt, time.Duration(500)*time.Millisecond) // period 500 ms
		if nil != err {
			return err
		}
		if state != 1 {
			return fmt.Errorf("check task prepay state failed, task state is not in `prepay` on beginConsumeByDataToken()")
		}

		return nil
	//case carriertypespb.TaskRole_TaskRole_PowerSupplier:
	//	return nil // do nothing ...
	//case carriertypespb.TaskRole_TaskRole_Receiver:
	//	return nil // do nothing ...
	default:
		return fmt.Errorf("unknown task role on beginConsumeByDataToken()")
	}
}

func (m *Manager) endConsumeMetadataOrPower(task *types.NeedExecuteTask, localTask *types.Task) error {
	switch m.config.MetadataConsumeOption {
	case 1: // use metadataAuth
		return m.endConsumeByMetadataAuth(task, localTask)
	case 2: // use datatoken
		return m.endConsumeByDataToken(task, localTask)
	default: // use nothing
		return nil
	}
}
func (m *Manager) endConsumeByMetadataAuth(task *types.NeedExecuteTask, localTask *types.Task) error {
	return nil // do nothing.
}
func (m *Manager) endConsumeByDataToken(task *types.NeedExecuteTask, localTask *types.Task) error {

	partyId := task.GetLocalTaskOrganization().GetPartyId()

	switch task.GetLocalTaskRole() {
	case commonconstantpb.TaskRole_TaskRole_Sender:

		// query consumeSpec of task
		var consumeSpec *types.DatatokenPaySpec

		if "" == strings.Trim(task.GetConsumeSpec(), "") {
			return fmt.Errorf("consumeSpec about task is empty on endConsumeByDataToken(), consumeSpec: %s", task.GetConsumeSpec())
		}

		if err := json.Unmarshal([]byte(task.GetConsumeSpec()), &consumeSpec); nil != err {
			return fmt.Errorf("cannot json unmarshal consumeSpec on endConsumeByDataToken(), consumeSpec: %s, %s", task.GetConsumeSpec(), err)
		}

		if partyId != localTask.GetTaskSender().GetPartyId() {
			return fmt.Errorf("this partyId is not task sender on endConsumeByDataToken(), partyId: %s, sender partyId: %s",
				partyId, localTask.GetTaskSender().GetPartyId())
		}

		// check task state in contract
		//
		// If the state of the contract is not "prepay",
		// the settlement action will not be performed
		if consumeSpec.GetConsumed() != 1 {
			log.Warnf("Warning the task state value of contract is not `prepay`, the settlement action will not be performed, taskId: {%s}, partyId: {%s}, state: {%d}",
				task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), consumeSpec.GetConsumed())
			return nil
		}

		taskIdBigInt, err := hexutil.DecodeBig("0x" + strings.TrimLeft(strings.Trim(task.GetTaskId(), types.PREFIX_TASK_ID+"0x"), "\x00"))
		if nil != err {
			return fmt.Errorf("cannot decode taskId to big.Int on endConsumeByDataToken(), %s", err)
		}

		log.Debugf("Start call token20Pay.settle(), taskId: {%s}, partyId: {%s}, call params{taskIdBigInt: %d, gasUsedPrepay: %d}",
			task.GetTaskId(), partyId, taskIdBigInt, consumeSpec.GetGasUsed())

		// start prepay dataToken
		txHash, err := m.token20PayMng.Settle(taskIdBigInt, consumeSpec.GetGasUsed())
		if nil != err {
			return fmt.Errorf("cannot call token20Pay to settle datatoken on endConsumeByDataToken(), %s", err)
		}

		log.Debugf("Succeed send `contract settle()` tx to blockchain on endConsumeByDataToken(), taskId: {%s}, partyId: {%s}, taskIdBigInt: {%d}, txHash: {%s}",
			task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), taskIdBigInt, txHash.String())

		// make sure the `prepay` tx into blockchain
		timeout := time.Duration(localTask.GetTaskData().GetOperationCost().GetDuration()) * time.Millisecond
		ctx, cancelFn := context.WithTimeout(context.Background(), timeout)
		//ctx, cancelFn := context.WithCancel(context.Background())
		defer cancelFn()

		receipt := m.token20PayMng.GetReceipt(ctx, txHash, time.Duration(500)*time.Millisecond) // period 500 ms
		if nil == receipt {
			return fmt.Errorf("settle dataToken failed, the transaction had not receipt on endConsumeByDataToken(), txHash: {%s}", txHash.String())
		}

		// contract tx execute failed.
		if receipt.Status == 0 {
			return fmt.Errorf("settle dataToken failed, the transaction receipt status is %d on endConsumeByDataToken(), txHash: {%s}", receipt.Status, txHash.String())
		}

		log.Debugf("Succeed execute `contract settle()` tx on blockchain on endConsumeByDataToken(), taskId: {%s}, partyId: {%s}, taskIdBigInt: {%d}, txHash: {%s}",
			task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), taskIdBigInt, txHash.String())

		return nil
	case commonconstantpb.TaskRole_TaskRole_DataSupplier:
		return nil // do nothing ...
	case commonconstantpb.TaskRole_TaskRole_PowerSupplier:
		return nil // do nothing ...
	case commonconstantpb.TaskRole_TaskRole_Receiver:
		return nil // do nothing ...
	default:
		return fmt.Errorf("unknown task role on endConsumeByDataToken()")
	}
}

// To execute task
func (m *Manager) driveTaskForExecute(task *types.NeedExecuteTask, localTask *types.Task) error {

	// 1、 consume the resource of task
	// TODO 打开这里 ...
	//if err := m.beginConsumeMetadataOrPower(task, localTask); nil != err {
	//	return err
	//}

	// 2、 update needExecuteTask to disk
	if err := m.resourceMng.GetDB().StoreNeedExecuteTask(task); nil != err {
		log.WithError(err).Errorf("store needExecuteTask failed, taskId: {%s}, partyId: {%s}", task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId())
	}
	// 3、 let's task execute
	switch task.GetLocalTaskRole() {
	case commonconstantpb.TaskRole_TaskRole_DataSupplier, commonconstantpb.TaskRole_TaskRole_Receiver:
		return m.executeTaskOnDataNode(task, localTask)
	case commonconstantpb.TaskRole_TaskRole_PowerSupplier:
		return m.executeTaskOnJobNode(task, localTask)
	default:
		return nil
	}

}
func (m *Manager) executeTaskOnDataNode(task *types.NeedExecuteTask, localTask *types.Task) error {

	// find dataNodeId with self vote
	var dataNodeId string
	dataNodes, err := m.resourceMng.GetDB().QueryRegisterNodeList(carrierapipb.PrefixTypeDataNode)
	if nil != err {
		log.Errorf("Failed to query internal dataNode arr on `taskManager.executeTaskOnDataNode()`, taskId: {%s}, role: {%s}, partyId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId())
		return fmt.Errorf("query internal dataNodes failed when ready to execute task")
	}
	for _, dataNode := range dataNodes {
		if dataNode.GetExternalIp() == task.GetLocalResource().GetIp() && dataNode.GetExternalPort() == task.GetLocalResource().GetPort() {
			dataNodeId = dataNode.GetId()
			break
		}
	}

	if "" == strings.Trim(dataNodeId, "") {
		log.Errorf("Failed to find dataNodeId of self vote resource on `taskManager.executeTaskOnDataNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId)
		return fmt.Errorf("not find dataNodeId of self vote resource when ready to execute task")
	}

	// clinet *grpclient.DataNodeClient,
	client, has := m.resourceMng.QueryDataNodeClient(dataNodeId)
	if !has {
		log.Errorf("Failed to query internal data node on `taskManager.executeTaskOnDataNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId)
		return fmt.Errorf("dataNode client not found when ready to execute task")
	}
	if client.IsNotConnected() {
		if err := client.Reconnect(); nil != err {
			log.WithError(err).Errorf("Failed to connect internal data node on `taskManager.executeTaskOnDataNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
				task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId)
			return fmt.Errorf("the dataNode is not connected when ready to execute task")
		}
	}

	req, err := m.makeTaskReadyGoReq(task, localTask)
	if nil != err {
		log.WithError(err).Errorf("Falied to make TaskReadyGoReq on `taskManager.executeTaskOnDataNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId)
		return fmt.Errorf("make dataNode rpc req failed when ready to execute task")
	}

	resp, err := client.HandleTaskReadyGo(req)
	if nil != err {
		log.WithError(err).Errorf("Falied to call publish schedTask to `data-Fighter` node to executing, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId)
		return fmt.Errorf("call dataNode rpc api failed when ready to execute task")
	}
	if resp.GetStatus() != 0 {
		log.Errorf("Falied to executing task from `data-Fighter` node to executing the resp code is not `ok`, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}, resp: %s",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId, resp.String())
		return fmt.Errorf("dataNode rpc api wrong resp %s when ready to execute task", resp.GetMsg())
	}

	log.Infof("Success to publish schedTask to `data-Fighter` node to executing, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
		task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId)
	return nil
}
func (m *Manager) executeTaskOnJobNode(task *types.NeedExecuteTask, localTask *types.Task) error {

	// find jobNodeId with self vote
	var jobNodeId string
	jobNodes, err := m.resourceMng.GetDB().QueryRegisterNodeList(carrierapipb.PrefixTypeJobNode)
	if nil != err {
		log.Errorf("Failed to query internal jobNode arr on `taskManager.executeTaskOnJobNode()`, taskId: {%s}, role: {%s}, partyId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId())
		return fmt.Errorf("query internal jobNodes failed when ready to execute task")
	}
	for _, jobNode := range jobNodes {
		if jobNode.GetExternalIp() == task.GetLocalResource().GetIp() && jobNode.GetExternalPort() == task.GetLocalResource().GetPort() {
			jobNodeId = jobNode.GetId()
			break
		}
	}

	if "" == strings.Trim(jobNodeId, "") {
		log.Errorf("Failed to find jobNodeId of self vote resource on `taskManager.executeTaskOnJobNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId)
		return fmt.Errorf("not find jobNodeId of self vote resource when ready to execute task")
	}

	// clinet *grpclient.JobNodeClient,
	client, has := m.resourceMng.QueryJobNodeClient(jobNodeId)
	if !has {
		log.Errorf("Failed to query internal job node on `taskManager.executeTaskOnJobNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId)
		return fmt.Errorf("jobNode client not found when ready to execute task")
	}
	if client.IsNotConnected() {
		if err := client.Reconnect(); nil != err {
			log.WithError(err).Errorf("Failed to connect internal job node on `taskManager.executeTaskOnJobNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
				task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId)
			return fmt.Errorf("the jobNode is not connected when ready to execute task")
		}
	}

	req, err := m.makeTaskReadyGoReq(task, localTask)
	if nil != err {
		log.WithError(err).Errorf("Falied to make TaskReadyGoReq on `taskManager.executeTaskOnJobNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId)
		return fmt.Errorf("make jobNode rpc req failed when ready to execute task")
	}

	resp, err := client.HandleTaskReadyGo(req)
	if nil != err {
		log.WithError(err).Errorf("Falied to publish schedTask to `job-Fighter` node to executing, taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId)
		return fmt.Errorf("call jobNode rpc api failed when ready to execute task")
	}
	if resp.GetStatus() != 0 {
		log.Errorf("Falied to publish schedTask to `job-Fighter` node to executing the resp code is not `ok`, taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}, resp: %s",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId, resp.String())
		return fmt.Errorf("jobNode rpc api wrong resp %s when ready to execute task", resp.GetMsg())
	}

	log.Infof("Success to publish schedTask to `job-Fighter` node to executing, taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
		task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId)
	return nil
}

// To terminate task
func (m *Manager) driveTaskForTerminate(task *types.NeedExecuteTask) error {

	// annonce fighter processor to terminate this task
	switch task.GetLocalTaskRole() {
	case commonconstantpb.TaskRole_TaskRole_DataSupplier, commonconstantpb.TaskRole_TaskRole_Receiver:
		return m.terminateTaskOnDataNode(task)
	case commonconstantpb.TaskRole_TaskRole_PowerSupplier:
		return m.terminateTaskOnJobNode(task)
	default:
		log.Errorf("Faided to driveTaskForTerminate(), Unknown task role, taskId: {%s}, taskRole: {%s}", task.GetTaskId(), task.GetLocalTaskRole().String())
		return fmt.Errorf("unknown resource node type when ready to terminate task")
	}
}
func (m *Manager) terminateTaskOnDataNode(task *types.NeedExecuteTask) error {

	// find dataNodeId with self vote
	var dataNodeId string
	dataNodes, err := m.resourceMng.GetDB().QueryRegisterNodeList(carrierapipb.PrefixTypeDataNode)
	if nil != err {
		log.Errorf("Failed to query internal dataNode arr on `taskManager.terminateTaskOnDataNode()`, taskId: {%s}, role: {%s}, partyId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId())
		return fmt.Errorf("query internal dataNodes failed when ready to terminate task")
	}
	for _, dataNode := range dataNodes {
		if dataNode.GetExternalIp() == task.GetLocalResource().GetIp() && dataNode.GetExternalPort() == task.GetLocalResource().GetPort() {
			dataNodeId = dataNode.GetId()
			break
		}
	}

	if "" == strings.Trim(dataNodeId, "") {
		log.Errorf("Failed to find dataNodeId of self vote resource on `taskManager.terminateTaskOnDataNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId)
		return fmt.Errorf("not find dataNodeId of self vote resource when ready to terminate task")
	}

	// clinet *grpclient.DataNodeClient,
	client, has := m.resourceMng.QueryDataNodeClient(dataNodeId)
	if !has {
		log.Errorf("Failed to query internal data node on `taskManager.terminateTaskOnDataNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId)
		return fmt.Errorf("dataNode client not found when ready to terminate task")
	}
	if client.IsNotConnected() {
		if err := client.Reconnect(); nil != err {
			log.WithError(err).Errorf("Failed to connect internal data node on `taskManager.terminateTaskOnDataNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
				task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId)
			return fmt.Errorf("the dataNode is not connected when ready to terminate task")
		}
	}

	req, err := m.makeTerminateTaskReq(task)
	if nil != err {
		log.WithError(err).Errorf("Falied to make TaskCancelReq on `taskManager.terminateTaskOnDataNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId)
		return fmt.Errorf("make dataNode rpc req failed when ready to terminate task")
	}

	resp, err := client.HandleCancelTask(req)
	if nil != err {
		log.WithError(err).Errorf("Falied to call publish schedTask to `data-Fighter` node to terminating, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId)
		return fmt.Errorf("call dataNode rpc api failed when ready to terminate task")
	}
	if resp.GetStatus() != 0 {
		log.Errorf("Falied to executing task from `data-Fighter` node to terminating the resp code is not `ok`, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}, resp: %s",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId, resp.String())
		return fmt.Errorf("dataNode rpc api wrong resp %s when ready to terminate task", resp.GetMsg())
	}

	log.Infof("Success to publish schedTask to `data-Fighter` node to terminating, taskId: {%s}, role: {%s}, partyId: {%s}, dataNodeId: {%s}",
		task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), dataNodeId)
	return nil
}
func (m *Manager) terminateTaskOnJobNode(task *types.NeedExecuteTask) error {

	// find jobNodeId with self vote
	var jobNodeId string
	jobNodes, err := m.resourceMng.GetDB().QueryRegisterNodeList(carrierapipb.PrefixTypeJobNode)
	if nil != err {
		log.Errorf("Failed to query internal jobNode arr on `taskManager.terminateTaskOnJobNode()`, taskId: {%s}, role: {%s}, partyId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId())
		return fmt.Errorf("query internal jobNodes failed when ready to terminate task")
	}
	for _, jobNode := range jobNodes {
		if jobNode.GetExternalIp() == task.GetLocalResource().GetIp() && jobNode.GetExternalPort() == task.GetLocalResource().GetPort() {
			jobNodeId = jobNode.GetId()
			break
		}
	}

	if "" == strings.Trim(jobNodeId, "") {
		log.Errorf("Failed to find jobNodeId of self vote resource on `taskManager.terminateTaskOnJobNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId)
		return fmt.Errorf("not find jobNodeId of self vote resource when ready to terminate task")
	}

	// clinet *grpclient.JobNodeClient,
	client, has := m.resourceMng.QueryJobNodeClient(jobNodeId)
	if !has {
		log.Errorf("Failed to query internal job node on `taskManager.terminateTaskOnJobNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId)
		return fmt.Errorf("jobNode client not found when ready to terminate task")
	}
	if client.IsNotConnected() {
		if err := client.Reconnect(); nil != err {
			log.WithError(err).Errorf("Failed to connect internal job node on `taskManager.terminateTaskOnJobNode()`, taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
				task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId)
			return fmt.Errorf("the jobNode is not connected when ready to terminate task")
		}
	}

	req, err := m.makeTerminateTaskReq(task)
	if nil != err {
		log.WithError(err).Errorf("Falied to make TaskCancelReq on `taskManager.terminateTaskOnJobNode()`,taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId)
		return fmt.Errorf("make jobNode rpc req failed when ready to terminate task")
	}

	resp, err := client.HandleCancelTask(req)
	if nil != err {
		log.WithError(err).Errorf("Falied to publish schedTask to `job-Fighter` node to terminating, taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId)
		return fmt.Errorf("call jobNode rpc api failed when ready to terminate task")
	}
	if resp.GetStatus() != 0 {
		log.Errorf("Falied to publish schedTask to `job-Fighter` node to terminating the resp code is not `ok`,taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}, resp: %s",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId, resp.String())
		return fmt.Errorf("jobNode rpc api wrong resp %s when ready to terminate task", resp.GetMsg())
	}

	log.Infof("Success to publish schedTask to `job-Fighter` node to terminating, taskId: {%s}, role: {%s}, partyId: {%s}, jobNodeId: {%s}",
		task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId(), jobNodeId)
	return nil
}

func (m *Manager) publishFinishedTaskToDataCenter(task *types.NeedExecuteTask, localTask *types.Task, delay bool) {

	/**
	++++++++++++++++++++++++++++
	NOTE: the needExecuteTask must be sender's here. (task.GetLocalTaskOrganization is sender's identity)
	++++++++++++++++++++++++++++
	*/

	handleFn := func() {

		eventList, err := m.resourceMng.GetDB().QueryTaskEventList(task.GetTaskId())
		if nil != err {
			log.WithError(err).Errorf("Failed to Query all task event list for sending datacenter on publishFinishedTaskToDataCenter, taskId: {%s}, partyId: {%s}",
				task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId())
			return
		}

		// check all events of this task, and change task state finally.
		var isFailed bool
		for _, event := range eventList {
			if event.GetType() == ev.TaskFailed.GetType() {
				isFailed = true
				break
			}
		}
		var taskState commonconstantpb.TaskState
		if isFailed {
			taskState = commonconstantpb.TaskState_TaskState_Failed
		} else {
			taskState = commonconstantpb.TaskState_TaskState_Succeed
		}

		// 1、settle metadata or power usage.
		// TODO 打开这里 ...
		//if err := m.endConsumeMetadataOrPower(task, localTask); nil != err {
		//	log.WithError(err).Errorf("Failed to settle consume metadata or power on publishFinishedTaskToDataCenter, taskId: {%s}, partyId: {%s}, taskState: {%s}",
		//		task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), taskState.String())
		//}

		log.Debugf("Start publishFinishedTaskToDataCenter, taskId: {%s}, partyId: {%s}, taskState: {%s}",
			task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), taskState.String())

		// 2、fill task by events AND push task to datacenter
		finalTask := m.fillTaskEventAndFinishedState(localTask, eventList, taskState)
		if err := m.resourceMng.GetDB().InsertTask(finalTask); nil != err {
			log.WithError(err).Errorf("Failed to save task to datacenter on publishFinishedTaskToDataCenter, taskId: {%s}, partyId: {%s}",
				task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId())
		}

		// 3、clear up all cache of task
		if err := m.RemoveExecuteTaskStateAfterExecuteTask("on taskManager.publishFinishedTaskToDataCenter()", task.GetTaskId(),
			task.GetLocalTaskOrganization().GetPartyId(), resource.SetAllReleaseResourceOption(), true); nil != err {
			log.WithError(err).Errorf("Failed to call RemoveExecuteTaskStateAfterExecuteTask() on publishFinishedTaskToDataCenter(), taskId: {%s},  partyId: {%s}, remote pid: {%s}",
				task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), task.GetRemotePID())
			return
		}

		log.Debugf("Finished publishFinishedTaskToDataCenter, taskId: {%s}, partyId: {%s}, taskState: {%s}",
			task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), taskState)
	}

	if delay {
		// delays handling some logic. (default delay 10s)
		runutil.RunOnce(context.TODO(), senderExecuteTaskExpire, handleFn)
	} else {
		handleFn()
	}
}
func (m *Manager) sendTaskResultMsgToTaskSender(task *types.NeedExecuteTask, localTask *types.Task) {

	// 1、settle metadata or power usage.
	// TODO 打开这里 ...
	//if err := m.endConsumeMetadataOrPower(task, localTask); nil != err {
	//	log.WithError(err).Errorf("Failed to settle consume metadata or power on sendTaskResultMsgToTaskSender, taskId: {%s},  partyId: {%s}",
	//		task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId())
	//}

	// 2、push all events of task to task sender.
	log.Debugf("Start sendTaskResultMsgToTaskSender, taskId: {%s}, partyId: {%s}, remote pid: {%s}",
		task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), task.GetRemotePID())

	var (
		option        resource.ReleaseResourceOption
		taskResultMsg *carriernetmsgtaskmngpb.TaskResultMsg
	)

	// when other task partner and task sender is same identity,
	// we don't need to removed local task and local eventList
	if task.GetLocalTaskOrganization().GetIdentityId() == task.GetRemoteTaskOrganization().GetIdentityId() {
		option = resource.SetUnlockLocalResorce() // unlock local resource of partyId, but don't remove local task and events of partyId
	} else {
		option = resource.SetAllReleaseResourceOption() // unlock local resource and remove local task and events
		// broadcast `task result msg` to reply remote peer
		taskResultMsg = m.makeTaskResultMsgWithEventList(task)
	}

	// 3、clear up all cache of task
	if err := m.RemoveExecuteTaskStateAfterExecuteTask("on sendTaskResultMsgToTaskSender()", task.GetTaskId(),
		task.GetLocalTaskOrganization().GetPartyId(), option, false); nil != err {
		log.WithError(err).Errorf("Failed to call RemoveExecuteTaskStateAfterExecuteTask() on sendTaskResultMsgToTaskSender(), taskId: {%s}, partyId: {%s}, remote pid: {%s}",
			task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), task.GetRemotePID())
		return
	}

	if nil != taskResultMsg {

		msg := types.FetchTaskResultMsg(taskResultMsg)
		// signature the msg and fill sign field of taskResourceUsageMsg
		sign, err := signutil.SignMsg(msg.Hash().Bytes(), m.nodePriKey)
		if nil != err {
			log.WithError(err).Errorf("failed to sign taskResultMsg, taskId: {%s}, partyId: {%s}, remote pid: {%s}, err: %s",
				task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), task.GetRemotePID(), err)
			return
		}
		taskResultMsg.Sign = sign

		if err := m.p2p.Broadcast(context.TODO(), taskResultMsg); nil != err {
			log.WithError(err).Errorf("failed to call `SendTaskResultMsg` on sendTaskResultMsgToTaskSender(), taskId: {%s}, partyId: {%s}, remote pid: {%s}",
				task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), task.GetRemotePID())
		} else {
			log.WithField("traceId", traceutil.GenerateTraceID(taskResultMsg)).Debugf("Succeed broadcast taskResultMsg to taskSender on sendTaskResultMsgToTaskSender(), taskId: {%s}, partyId: {%s}, remote pid: {%s}",
				task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), task.GetRemotePID())
		}
	}

	log.Debugf("Finished sendTaskResultMsgToTaskSender, taskId: {%s}, partyId: {%s}, remote pid: {%s}",
		task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), task.GetRemotePID())
}
func (m *Manager) broadcastTaskTerminateMsg(task *types.Task) error {

	sender := task.GetTaskSender()

	needSendLocalMsgFn := func() bool {
		for i := 0; i < len(task.GetTaskData().GetDataSuppliers()); i++ {
			if sender.GetIdentityId() == task.GetTaskData().GetDataSuppliers()[i].GetIdentityId() {
				return true
			}
		}

		for i := 0; i < len(task.GetTaskData().GetPowerSuppliers()); i++ {
			if sender.GetIdentityId() == task.GetTaskData().GetPowerSuppliers()[i].GetIdentityId() {
				return true
			}
		}

		for i := 0; i < len(task.GetTaskData().GetReceivers()); i++ {
			if sender.GetIdentityId() == task.GetTaskData().GetReceivers()[i].GetIdentityId() {
				return true
			}
		}
		return false
	}

	terminateMsg := &carriernetmsgtaskmngpb.TaskTerminateMsg{
		MsgOption: &carriernetmsgcommonpb.MsgOption{
			ProposalId:      carriercommon.Hash{}.Bytes(),
			SenderRole:      uint64(commonconstantpb.TaskRole_TaskRole_Sender),
			SenderPartyId:   []byte(sender.GetPartyId()),
			ReceiverRole:    uint64(commonconstantpb.TaskRole_TaskRole_Unknown),
			ReceiverPartyId: []byte{},
			MsgOwner: &carriernetmsgcommonpb.TaskOrganizationIdentityInfo{
				Name:       []byte(sender.GetNodeName()),
				NodeId:     []byte(sender.GetNodeId()),
				IdentityId: []byte(sender.GetIdentityId()),
				PartyId:    []byte(sender.GetPartyId()),
			},
		},
		TaskId:   []byte(task.GetTaskId()),
		CreateAt: timeutils.UnixMsecUint64(),
		Sign:     nil,
	}
	msg := &types.TaskTerminateTaskMngMsg{
		MsgOption: types.FetchMsgOption(terminateMsg.GetMsgOption()),
		TaskId:    string(terminateMsg.GetTaskId()),
		CreateAt:  terminateMsg.GetCreateAt(),
	}

	// signature the msg and fill sign field of taskTerminateMsg
	sign, err := signutil.SignMsg(msg.Hash().Bytes(), m.nodePriKey)
	if nil != err {
		return fmt.Errorf("sign taskTerminateMsg, %s", err)
	}
	terminateMsg.Sign = sign

	errs := make([]string, 0)
	if needSendLocalMsgFn() {
		if err := m.onTaskTerminateMsg("", terminateMsg, types.LocalNetworkMsg); nil != err {
			errs = append(errs, fmt.Sprintf("broadcast taskTerminateMsg to local peer, %s", err))
		} else {
			log.WithField("traceId", traceutil.GenerateTraceID(terminateMsg)).Debugf("Succeed to call local `onTaskTerminateMsg` taskId: %s",
				task.GetTaskId())
		}
	}

	if err := m.p2p.Broadcast(context.TODO(), terminateMsg); nil != err {
		errs = append(errs, fmt.Sprintf("broadcast taskTerminateMsg to remote peer, %s", err))
	} else {
		log.WithField("traceId", traceutil.GenerateTraceID(terminateMsg)).Debugf("Succeed to broadcast `TaskTerminateMsg` taskId: %s",
			task.GetTaskId())
	}

	if len(errs) != 0 {
		return fmt.Errorf(
			"\n######################################################## \n%s\n########################################################\n",
			strings.Join(errs, "\n"))
	}

	return nil
}
func (m *Manager) sendTaskEvent(event *carriertypespb.TaskEvent) {
	go func(event *carriertypespb.TaskEvent) {
		m.eventCh <- event
		log.Debugf("Succeed send a task event to manager.loop() on taskManager.sendTaskEvent(), event: %s", event.String())
	}(event)
}
func (m *Manager) sendNeedExecuteTaskByAction(task *types.NeedExecuteTask) {
	go func(task *types.NeedExecuteTask) { // asynchronous transmission to reduce Chan blocking
		m.needExecuteTaskCh <- task
	}(task)
}

func (m *Manager) StoreExecuteTaskStateBeforeExecuteTask(logdesc, taskId, partyId string) error {
	// Store task exec status
	if err := m.resourceMng.GetDB().StoreLocalTaskExecuteStatusValRunningByPartyId(taskId, partyId); nil != err {
		log.WithError(err).Errorf("Failed to store local task about `running` status %s, taskId: {%s}, partyId: {%s}",
			logdesc, taskId, partyId)
		return err
	}
	log.Debugf("Succeed store local task about `running` status %s, taskId: {%s}, partyId: {%s}",
		logdesc, taskId, partyId)
	// do anythings else?
	return nil
}
func (m *Manager) RemoveExecuteTaskStateAfterExecuteTask(logdesc, taskId, partyId string, option resource.ReleaseResourceOption, isSender bool) error {
	if err := m.resourceMng.GetDB().RemoveLocalTaskExecuteStatusByPartyId(taskId, partyId); nil != err {
		log.WithError(err).Errorf("Failed to remove task executing status %s, taskId: {%s} partyId: {%s}, isSender: {%v}",
			logdesc, taskId, partyId, isSender)
		return err
	}

	log.Debugf("Succeed remove task executing status %s, taskId: {%s} partyId: {%s}, isSender: {%v}",
		logdesc, taskId, partyId, isSender)

	// clean local task cache
	m.resourceMng.ReleaseLocalResourceWithTask(logdesc, taskId, partyId, option, isSender)
	return nil
}

func (m *Manager) publishBadTaskToDataCenter(task *types.Task, events []*carriertypespb.TaskEvent, reason string) error {
	task.GetTaskData().TaskEvents = events
	task.GetTaskData().State = commonconstantpb.TaskState_TaskState_Failed
	task.GetTaskData().Reason = reason
	task.GetTaskData().EndAt = timeutils.UnixMsecUint64()

	m.resourceMng.GetDB().RemoveLocalTask(task.GetTaskId())
	m.resourceMng.GetDB().RemoveTaskEventList(task.GetTaskId())

	return m.resourceMng.GetDB().InsertTask(task)
}

func (m *Manager) fillTaskEventAndFinishedState(task *types.Task, eventList []*carriertypespb.TaskEvent, state commonconstantpb.TaskState) *types.Task {
	task.GetTaskData().TaskEvents = eventList
	task.GetTaskData().EndAt = timeutils.UnixMsecUint64()
	task.GetTaskData().State = state
	return task
}

func (m *Manager) makeTaskReadyGoReq(task *types.NeedExecuteTask, localTask *types.Task) (*fighterpbtypes.TaskReadyGoReq, error) {

	var dataPartyArr []string
	var powerPartyArr []string
	var receiverPartyArr []string

	peerList := make([]*fighterpbtypes.Party, 0)

	for _, dataSupplier := range task.GetResources().GetDataSupplierPeerInfos() {
		portStr := string(dataSupplier.GetPort())
		port, err := strconv.Atoi(portStr)
		if nil != err {
			return nil, err
		}
		peerList = append(peerList, &fighterpbtypes.Party{
			Ip:      string(dataSupplier.GetIp()),
			Port:    int32(port),
			PartyId: string(dataSupplier.GetPartyId()),
		})
		dataPartyArr = append(dataPartyArr, string(dataSupplier.GetPartyId()))
	}

	for _, powerSupplier := range task.GetResources().GetPowerSupplierPeerInfos() {
		portStr := string(powerSupplier.GetPort())
		port, err := strconv.Atoi(portStr)
		if nil != err {
			return nil, err
		}
		peerList = append(peerList, &fighterpbtypes.Party{
			Ip:      string(powerSupplier.GetIp()),
			Port:    int32(port),
			PartyId: string(powerSupplier.GetPartyId()),
		})

		powerPartyArr = append(powerPartyArr, string(powerSupplier.GetPartyId()))
	}

	for _, receiver := range task.GetResources().GetResultReceiverPeerInfos() {
		portStr := string(receiver.GetPort())
		port, err := strconv.Atoi(portStr)
		if nil != err {
			return nil, err
		}
		peerList = append(peerList, &fighterpbtypes.Party{
			Ip:      string(receiver.GetIp()),
			Port:    int32(port),
			PartyId: string(receiver.GetPartyId()),
		})

		receiverPartyArr = append(receiverPartyArr, string(receiver.GetPartyId()))
	}

	selfCfgParams, err := m.makeReqCfgParams(task, localTask)
	if nil != err {
		return nil, fmt.Errorf("make contractParams failed, %s", err)
	}

	connectPolicyFormat, connectPolicy, err := m.makeConnectPolicy(task, localTask)
	if nil != err {
		return nil, fmt.Errorf("make makeConnectPolicy failed, %s", err)
	}
	log.Debugf("Succeed make selfCfgParams field of req, taskId:{%s}, partyId: {%s}, selfCfgParams: %s, connectPolicyType: %s, connectPolicy: %s",
		task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), selfCfgParams, connectPolicyFormat.String(), connectPolicy)

	req := &fighterpbtypes.TaskReadyGoReq{

		/**
		TaskId                 string
		PartyId                string
		EnvId                  string
		Parties                []*Party
		AlgorithmCode          string
		SelfCfgParams          string
		AlgorithmDynamicParams string
		DataPartyIds           []string
		ComputationPartyIds    []string
		ResultPartyIds         []string
		Duration               uint64
		Memory                 uint64
		Processor              uint32
		Bandwidth              uint64
		ConnectPolicyFormat    ConnectPolicyFormat
		ConnectPolicy          string
		*/
		TaskId:  task.GetTaskId(),
		PartyId: task.GetLocalTaskOrganization().GetPartyId(),
		//EnvId: "",
		Parties:                peerList,
		AlgorithmCode:          localTask.GetTaskData().GetAlgorithmCode(),
		SelfCfgParams:          selfCfgParams,
		AlgorithmDynamicParams: localTask.GetTaskData().GetAlgorithmCodeExtraParams(),
		DataPartyIds:           dataPartyArr,
		ComputationPartyIds:    powerPartyArr,
		ResultPartyIds:         receiverPartyArr,
		Duration:               localTask.GetTaskData().GetOperationCost().GetDuration(),
		Memory:                 localTask.GetTaskData().GetOperationCost().GetMemory(),
		Processor:              localTask.GetTaskData().GetOperationCost().GetProcessor(),
		Bandwidth:              localTask.GetTaskData().GetOperationCost().GetBandwidth(),
		ConnectPolicyFormat:    connectPolicyFormat,
		ConnectPolicy:          connectPolicy,
	}

	return req, nil
}

func (m *Manager) makeReqCfgParams(task *types.NeedExecuteTask, localTask *types.Task) (string, error) {

	/**
	# FOR DATANODE:

	{
		"part_id": "p0",
	    "input_data": [
	          {
	              "input_type": 3,  # 输入数据的类型，(算法用标识数据使用方式). 0:unknown, 1:origin_data, 2:psi_output, 3:model
	              "access_type": 1, # 访问数据的方式，(fighter用决定是否预先加载数据). 0:unknown, 1:local, 2:url
	              "data_type": 0,   # 数据的格式，(算法用标识数据格式). 0:unknown, 1:csv, 2:dir, 3:binary, 4:xls, 5:xlsx, 6:txt, 7:json
	              "data_path": "/task_result/task:0xdeefff3434..556/"  # 数据所在的本地路径
	          }
	    ]
	}

	# FOR JOBNODE:

	{
		"part_id": "y0",
	    "input_data": []
	}

	*/

	partyId := task.GetLocalTaskOrganization().GetPartyId()
	inputDataArr := make([]interface{}, 0)

	switch task.GetLocalTaskRole() {
	case commonconstantpb.TaskRole_TaskRole_DataSupplier:

		for i, policyType := range localTask.GetTaskData().GetDataPolicyTypes() {

			switch policyType {
			case uint32(commonconstantpb.OrigindataType_OrigindataType_CSV):

				var dataPolicy *types.TaskMetadataPolicyCSV
				if err := json.Unmarshal([]byte(localTask.GetTaskData().GetDataPolicyOptions()[i]), &dataPolicy); nil != err {
					return "", fmt.Errorf("can not unmarshal dataPolicyOption, %s, taskId: {%s}, partyId: {%s}", err, localTask.GetTaskId(), partyId)
				}

				if dataPolicy.GetPartyId() == partyId {
					inputData, err := m.metadataInputCSV(task, localTask, dataPolicy)
					if nil != err {
						return "", fmt.Errorf("can not unmarshal metadataInputCSV, %s, taskId: {%s}, partyId: {%s}, metadataId: {%s}, metadataName: {%s}",
							err, localTask.GetTaskId(), partyId, dataPolicy.GetMetadataId(), dataPolicy.GetMetadataName())
					}
					inputDataArr = append(inputDataArr, inputData)
				}
			case uint32(commonconstantpb.OrigindataType_OrigindataType_DIR):
				var dataPolicy *types.TaskMetadataPolicyDIR
				if err := json.Unmarshal([]byte(localTask.GetTaskData().GetDataPolicyOptions()[i]), &dataPolicy); nil != err {
					return "", fmt.Errorf("can not unmarshal dataPolicyOption, %s, taskId: {%s}, partyId: {%s}", err, localTask.GetTaskId(), partyId)
				}

				if dataPolicy.GetPartyId() == partyId {
					inputData, err := m.metadataInputDIR(task, localTask, dataPolicy)
					if nil != err {
						return "", fmt.Errorf("can not unmarshal metadataInputDIR, %s, taskId: {%s}, partyId: {%s}, metadataId: {%s}, metadataName: {%s}",
							err, localTask.GetTaskId(), partyId, dataPolicy.GetMetadataId(), dataPolicy.GetMetadataName())
					}
					inputDataArr = append(inputDataArr, inputData)
				}

			case uint32(commonconstantpb.OrigindataType_OrigindataType_BINARY):
				var dataPolicy *types.TaskMetadataPolicyBINARY
				if err := json.Unmarshal([]byte(localTask.GetTaskData().GetDataPolicyOptions()[i]), &dataPolicy); nil != err {
					return "", fmt.Errorf("can not unmarshal dataPolicyOption, %s, taskId: {%s}, partyId: {%s}", err, localTask.GetTaskId(), partyId)
				}

				if dataPolicy.GetPartyId() == partyId {
					inputData, err := m.metadataInputBINARY(task, localTask, dataPolicy)
					if nil != err {
						return "", fmt.Errorf("can not unmarshal metadataInputBINARY, %s, taskId: {%s}, partyId: {%s}, metadataId: {%s}, metadataName: {%s}",
							err, localTask.GetTaskId(), partyId, dataPolicy.GetMetadataId(), dataPolicy.GetMetadataName())
					}
					inputDataArr = append(inputDataArr, inputData)
				}
			default:
				return "", fmt.Errorf("unknown dataPolicy type, taskId: {%s}, partyId: {%s}, dataPolicyType: {%d}", task.GetTaskId(), partyId, policyType)
			}
		}
	case commonconstantpb.TaskRole_TaskRole_PowerSupplier:
		// do nothing...
	}

	scps := &types.SelfCfgParams{
		PartyId:   partyId,
		InputData: inputDataArr,
	}

	b, err := json.Marshal(scps)
	if nil != err {
		return "", fmt.Errorf("cannot json marshal selfCfgParams, %s, taskId: {%s}, partyId: {%s}", err, task.GetTaskId(), partyId)
	}

	return string(b), err
}

func (m *Manager) metadataInputCSV(task *types.NeedExecuteTask, localTask *types.Task, dataPolicy *types.TaskMetadataPolicyCSV) (*types.InputDataCSV, error) {

	var (
		dataPath        string
		keyColumn       string
		selectedColumns []string
	)

	metadataId := dataPolicy.GetMetadataId()
	internalMetadataFlag, err := m.resourceMng.GetDB().IsInternalMetadataById(metadataId)
	if nil != err {
		return nil, fmt.Errorf("cannot check metadata whether internal metadata, %s", err)
	}

	var metadata *types.Metadata

	// whether the metadata is internal metadata ?
	if internalMetadataFlag {
		// query internal metadata
		metadata, err = m.resourceMng.GetDB().QueryInternalMetadataById(metadataId)
		if nil != err {
			return nil, fmt.Errorf("cannot query internale metadata, %s", err)
		}
	} else {
		// query published metadata
		metadata, err = m.resourceMng.GetDB().QueryMetadataById(metadataId)
		if nil != err {
			return nil, fmt.Errorf("cannot query publish metadata, %s", err)
		}
	}

	if types.IsNotCSVdata(metadata.GetData().GetDataType()) {
		return nil, fmt.Errorf("dataType of metadata is not `CSV`, dataType: %s", metadata.GetData().GetDataType().String())
	}

	var metadataOption *types.MetadataOptionCSV
	if err := json.Unmarshal([]byte(metadata.GetData().GetMetadataOption()), &metadataOption); nil != err {
		return nil, fmt.Errorf("can not unmarshal `CSV` metadataOption, %s", err)
	}

	// collection all the column name cache
	columnNameCache := make(map[uint32]string, 0)
	for _, mdop := range metadataOption.GetMetadataColumns() {
		columnNameCache[mdop.GetIndex()] = mdop.GetName()
	}
	// find key column name
	if kname, ok := columnNameCache[dataPolicy.QueryKeyColumn()]; ok {
		keyColumn = kname
	} else {
		return nil, fmt.Errorf("not found the keyColumn of task dataPolicy on `CSV` metadataOption, columnIndex: {%d}", dataPolicy.QueryKeyColumn())
	}

	// find all select column names
	selectedColumns = make([]string, len(dataPolicy.QuerySelectedColumns()))
	for i, selectedColumnIndex := range dataPolicy.QuerySelectedColumns() {

		if sname, ok := columnNameCache[selectedColumnIndex]; ok {
			selectedColumns[i] = sname
		} else {
			return nil, fmt.Errorf("not found the selectColumn of task dataPolicy on `CSV` metadataOption, columnIndex: {%d}", selectedColumnIndex)
		}
	}

	dataPath = metadataOption.GetDataPath()
	if strings.Trim(dataPath, "") == "" {
		return nil, fmt.Errorf("dataPath is empty")
	}

	return &types.InputDataCSV{
		InputType:       dataPolicy.QueryInputType(),
		AccessType:      uint32(metadata.GetData().GetLocationType()),
		DataType:        uint32(metadata.GetData().GetDataType()),
		DataPath:        dataPath,
		KeyColumn:       keyColumn,
		SelectedColumns: selectedColumns,
	}, nil
}
func (m *Manager) metadataInputDIR(task *types.NeedExecuteTask, localTask *types.Task, dataPolicy *types.TaskMetadataPolicyDIR) (*types.InputDataDIR, error) {

	metadataId := dataPolicy.GetMetadataId()
	internalMetadataFlag, err := m.resourceMng.GetDB().IsInternalMetadataById(metadataId)
	if nil != err {
		return nil, fmt.Errorf("cannot check metadata whether internal metadata, %s", err)
	}

	var metadata *types.Metadata

	// whether the metadata is internal metadata ?
	if internalMetadataFlag {
		// query internal metadata
		metadata, err = m.resourceMng.GetDB().QueryInternalMetadataById(metadataId)
		if nil != err {
			return nil, fmt.Errorf("cannot query internale metadata, %s", err)
		}
	} else {
		// query published metadata
		metadata, err = m.resourceMng.GetDB().QueryMetadataById(metadataId)
		if nil != err {
			return nil, fmt.Errorf("cannot query publish metadata, %s", err)
		}
	}

	if types.IsNotDIRdata(metadata.GetData().GetDataType()) {
		return nil, fmt.Errorf("dataType of metadata is not `DIR`, dataType: %s", metadata.GetData().GetDataType().String())
	}

	var metadataOption *types.MetadataOptionDIR
	if err := json.Unmarshal([]byte(metadata.GetData().GetMetadataOption()), &metadataOption); nil != err {
		return nil, fmt.Errorf("can not unmarshal `DIR` metadataOption, %s", err)
	}

	dirPath := metadataOption.GetDirPath()
	if strings.Trim(dirPath, "") == "" {
		return nil, fmt.Errorf("dirPath is empty")
	}

	return &types.InputDataDIR{
		InputType:  dataPolicy.QueryInputType(),
		AccessType: uint32(metadata.GetData().GetLocationType()),
		DataType:   uint32(metadata.GetData().GetDataType()),
		DataPath:   dirPath,
	}, nil
}
func (m *Manager) metadataInputBINARY(task *types.NeedExecuteTask, localTask *types.Task, dataPolicy *types.TaskMetadataPolicyBINARY) (*types.InputDataBINARY, error) {

	metadataId := dataPolicy.GetMetadataId()
	internalMetadataFlag, err := m.resourceMng.GetDB().IsInternalMetadataById(metadataId)
	if nil != err {
		return nil, fmt.Errorf("cannot check metadata whether internal metadata, %s", err)
	}

	var metadata *types.Metadata

	// whether the metadata is internal metadata ?
	if internalMetadataFlag {
		// query internal metadata
		metadata, err = m.resourceMng.GetDB().QueryInternalMetadataById(metadataId)
		if nil != err {
			return nil, fmt.Errorf("cannot query internale metadata, %s", err)
		}
	} else {
		// query published metadata
		metadata, err = m.resourceMng.GetDB().QueryMetadataById(metadataId)
		if nil != err {
			return nil, fmt.Errorf("cannot query publish metadata, %s", err)
		}
	}

	if types.IsNotBINARYdata(metadata.GetData().GetDataType()) {
		return nil, fmt.Errorf("dataType of metadata is not `BINARY`, dataType: %s", metadata.GetData().GetDataType().String())
	}

	var metadataOption *types.MetadataOptionBINARY
	if err := json.Unmarshal([]byte(metadata.GetData().GetMetadataOption()), &metadataOption); nil != err {
		return nil, fmt.Errorf("can not unmarshal `BINARY` metadataOption, %s", err)
	}

	dataPath := metadataOption.GetDataPath()
	if strings.Trim(dataPath, "") == "" {
		return nil, fmt.Errorf("dataPath is empty")
	}

	return &types.InputDataBINARY{
		InputType:  dataPolicy.QueryInputType(),
		AccessType: uint32(metadata.GetData().GetLocationType()),
		DataType:   uint32(metadata.GetData().GetDataType()),
		DataPath:   dataPath,
	}, nil
}

func (m *Manager) makeConnectPolicy(task *types.NeedExecuteTask, localTask *types.Task) (commonconstantpb.ConnectPolicyFormat, string, error) {

	//partyId := task.GetLocalTaskOrganization().GetPartyId()

	var (
		format        commonconstantpb.ConnectPolicyFormat
		connectPolicy string
		err           error
	)

	for i, policyType := range localTask.GetTaskData().GetDataFlowPolicyTypes() {

		policy := localTask.GetTaskData().GetDataFlowPolicyOptions()[i]

		switch policyType {
		case types.TASK_DATAFLOW_POLICY_GENERAL_FULL_CONNECT, types.TASK_DATAFLOW_POLICY_GENERAL_DIRECTIONAL_CONNECT:
			format = commonconstantpb.ConnectPolicyFormat_ConnectPolicyFormat_Json
			connectPolicy = policy
		default:
			format, connectPolicy, err =
				commonconstantpb.ConnectPolicyFormat_ConnectPolicyFormat_Unknown, "", fmt.Errorf("unknown dataFlowPolicyType")
		}
		break
	}
	return format, connectPolicy, err
}

// make terminate rpc req
func (m *Manager) makeTerminateTaskReq(task *types.NeedExecuteTask) (*fighterpbtypes.TaskCancelReq, error) {
	return &fighterpbtypes.TaskCancelReq{
		TaskId:  task.GetTaskId(),
		PartyId: task.GetLocalTaskOrganization().GetPartyId(),
	}, nil
}

func (m *Manager) initConsumeSpecByConsumeOption(task *types.NeedExecuteTask) {
	// add consumeSpec
	switch m.config.MetadataConsumeOption {
	case 1: // use metadataAuth
		// pass
	case 2: // use datatoken
		taskIdBigInt, err := hexutil.DecodeBig("0x" + strings.TrimLeft(strings.Trim(task.GetTaskId(), types.PREFIX_TASK_ID+"0x"), "\x00"))
		if nil != err {
			log.WithError(err).Errorf("cannot decode taskId to big.Int on initConsumeSpecByConsumeOption()")
			return
		}
		// store consumeSpec into needExecuteTask
		consumeSpec := &types.DatatokenPaySpec{
			// task state in contract
			// constant int8 private NOTEXIST = -1;
			// constant int8 private BEGIN = 0;
			// constant int8 private PREPAY = 1;
			// constant int8 private SETTLE = 2;
			// constant int8 private END = 3;
			Consumed: int32(-1),
			//GasEstimated: 0,
			GasUsed: 0,
		}

		b, err := json.Marshal(consumeSpec)
		if nil != err {
			log.WithError(err).Errorf("cannot json marshal task consumeSpec on initConsumeSpecByConsumeOption(), consumeSpec: %v, %s", consumeSpec, err)
			return
		}
		task.SetConsumeQueryId(taskIdBigInt.String())
		task.SetConsumeSpec(string(b))

		//default: // use nothing
		//	// pass
	}
}

func (m *Manager) addNeedExecuteTaskCacheAndminotor(task *types.NeedExecuteTask, when int64) {

	m.addNeedExecuteTaskCache(task)

	// v0.3.0 add NeedExecuteTask Expire Monitor
	m.addmonitor(task, when)

}

func (m *Manager) addmonitor(task *types.NeedExecuteTask, when int64) {

	taskId, partyId := task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId()

	m.syncExecuteTaskMonitors.AddMonitor(types.NewExecuteTaskMonitor(taskId, partyId, when, func() {
		m.runningTaskCacheLock.Lock()
		defer m.runningTaskCacheLock.Unlock()

		cache, ok := m.runningTaskCache[taskId]
		if !ok {
			return
		}

		// 1、check local task from taskId
		localTask, err := m.resourceMng.GetDB().QueryLocalTask(taskId)
		if nil != err {
			for pid, _ := range cache {
				log.WithError(err).Warnf("Can not query local task info, clean current party task cache short circuit AND skip it, on `taskManager.expireTaskMonitor()`, taskId: {%s}, partyId: {%s}",
					taskId, pid)
				// clean current party task cache short circuit.
				delete(cache, pid)
				go m.resourceMng.GetDB().RemoveNeedExecuteTaskByPartyId(taskId, pid)
				if len(cache) == 0 {
					delete(m.runningTaskCache, taskId)
				} else {
					m.runningTaskCache[taskId] = cache
				}
				log.Debugf("Call expireTaskMonitor remove NeedExecuteTask as query local task info failed when task was expired, taskId: {%s}, partyId: {%s}", taskId, partyId)
				continue
			}
			return
		}

		// 2、 check partyId from cache and query needExecuteTask
		nt, ok := cache[partyId]
		if !ok {
			return
		}

		// 3、 handle ExpireTask
		if localTask.GetTaskData().GetState() == commonconstantpb.TaskState_TaskState_Running && localTask.GetTaskData().GetStartAt() != 0 {
			var duration uint64

			duration = timeutils.UnixMsecUint64() - localTask.GetTaskData().GetStartAt()

			log.Infof("Has task running expire, taskId: {%s}, partyId: {%s}, current running duration: {%d ms}, need running duration: {%d ms}",
				taskId, partyId, duration, localTask.GetTaskData().GetOperationCost().GetDuration())

			// 1、 store task expired (failed) event with current party
			m.resourceMng.GetDB().StoreTaskEvent(m.eventEngine.GenerateEvent(ev.TaskFailed.GetType(), taskId,
				nt.GetLocalTaskOrganization().GetIdentityId(), partyId,
				fmt.Sprintf("task running expire")))

			switch nt.GetLocalTaskRole() {
			case commonconstantpb.TaskRole_TaskRole_Sender:
				m.publishFinishedTaskToDataCenter(nt, localTask, true)
			default:
				// 2、terminate fighter processor for this task with current party
				m.driveTaskForTerminate(nt)
				m.sendTaskResultMsgToTaskSender(nt, localTask)
			}

			// clean current party task cache
			delete(cache, partyId)
			go m.resourceMng.GetDB().RemoveNeedExecuteTaskByPartyId(taskId, partyId)
			if len(cache) == 0 {
				delete(m.runningTaskCache, taskId)
			} else {
				m.runningTaskCache[taskId] = cache
			}
			log.Debugf("Call expireTaskMonitor remove NeedExecuteTask when task was expired, taskId: {%s}, partyId: {%s}", taskId, partyId)
		}
	}))
}

func (m *Manager) addNeedExecuteTaskCache(task *types.NeedExecuteTask) {
	m.runningTaskCacheLock.Lock()

	taskId, partyId := task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId()

	cache, ok := m.runningTaskCache[taskId]
	if !ok {
		cache = make(map[string]*types.NeedExecuteTask, 0)
	}
	cache[partyId] = task
	m.runningTaskCache[taskId] = cache
	if err := m.resourceMng.GetDB().StoreNeedExecuteTask(task); nil != err {
		log.WithError(err).Errorf("cannot store needExecuteTask on addNeedExecuteTaskCache(), taskId: {%s}, partyId: {%s}", taskId, partyId)
	}

	log.Debugf("Succeed call addNeedExecuteTaskCache, taskId: {%s}, partyId: {%s}", taskId, partyId)
	m.runningTaskCacheLock.Unlock()
}

func (m *Manager) updateNeedExecuteTaskCache(task *types.NeedExecuteTask) {
	m.runningTaskCacheLock.Lock()

	taskId, partyId := task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId()

	cache, ok := m.runningTaskCache[taskId]
	if !ok {
		return
	}

	_, ok = cache[partyId]
	if !ok {
		return
	}

	cache[partyId] = task
	m.runningTaskCache[taskId] = cache
	if err := m.resourceMng.GetDB().StoreNeedExecuteTask(task); nil != err {
		log.WithError(err).Errorf("cannot store needExecuteTask on updateNeedExecuteTaskCache(), taskId: {%s}, partyId: {%s}", taskId, partyId)
	}

	log.Debugf("Succeed call updateNeedExecuteTaskCache, taskId: {%s}, partyId: {%s}", taskId, partyId)
	m.runningTaskCacheLock.Unlock()
}

func (m *Manager) removeNeedExecuteTaskCache(taskId, partyId string) {
	m.runningTaskCacheLock.Lock()
	defer m.runningTaskCacheLock.Unlock()
	cache, ok := m.runningTaskCache[taskId]
	if !ok {
		return
	}
	delete(cache, partyId)
	go m.resourceMng.GetDB().RemoveNeedExecuteTaskByPartyId(taskId, partyId)
	if len(cache) == 0 {
		delete(m.runningTaskCache, taskId) // delete empty map[partyId]task
	} else {
		m.runningTaskCache[taskId] = cache // restore map[partyId]task if it is not empty
	}
	// v3.0  remove executeTask monitor
	m.syncExecuteTaskMonitors.DelMonitor(taskId, partyId)

	log.Debugf("Call removeNeedExecuteTaskCache, taskId: {%s}, partyId: {%s}",
		taskId, partyId)
}

func (m *Manager) hasNeedExecuteTaskCache(taskId, partyId string) bool {
	m.runningTaskCacheLock.RLock()
	defer m.runningTaskCacheLock.RUnlock()
	cache, ok := m.runningTaskCache[taskId]
	if !ok {
		log.Debugf("Call queryNeedExecuteTaskCache, taskId: {%s}, partyId: {%s}, has: {%v}",
			taskId, partyId, ok)
		return false
	}
	_, ok = cache[partyId]
	log.Debugf("Call queryNeedExecuteTaskCache, taskId: {%s}, partyId: {%s}, has: {%v}",
		taskId, partyId, ok)
	return ok
}

func (m *Manager) queryNeedExecuteTaskCache(taskId, partyId string) (*types.NeedExecuteTask, bool) {
	m.runningTaskCacheLock.RLock()
	defer m.runningTaskCacheLock.RUnlock()
	cache, ok := m.runningTaskCache[taskId]
	if !ok {
		//log.Debugf("Call queryNeedExecuteTaskCache, taskId: {%s}, partyId: {%s}, has: {%v}",
		//	taskId, partyId, ok)
		return nil, false
	}
	task, ok := cache[partyId]
	//log.Debugf("Call queryNeedExecuteTaskCache, taskId: {%s}, partyId: {%s}, has: {%v}",
	//	taskId, partyId, ok)
	return task, ok
}

func (m *Manager) mustQueryNeedExecuteTaskCache(taskId, partyId string) *types.NeedExecuteTask {
	task, _ := m.queryNeedExecuteTaskCache(taskId, partyId)
	return task
}

func (m *Manager) makeTaskResultMsgWithEventList(task *types.NeedExecuteTask) *carriernetmsgtaskmngpb.TaskResultMsg {

	if task.GetLocalTaskRole() == commonconstantpb.TaskRole_TaskRole_Sender {
		log.Errorf("the task sender can not make taskResultMsg")
		return nil
	}

	eventList, err := m.resourceMng.GetDB().QueryTaskEventListByPartyId(task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId())
	if nil != err {
		log.WithError(err).Errorf("Failed to make taskResultMsg with query task eventList, taskId {%s}, partyId {%s}",
			task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId())
		return nil
	}

	if len(eventList) == 0 {
		log.Errorf("Failed to make taskResultMsg with query task eventList is empty, taskId {%s}, partyId {%s}",
			task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId())
		return nil
	}

	return &carriernetmsgtaskmngpb.TaskResultMsg{
		MsgOption: &carriernetmsgcommonpb.MsgOption{
			ProposalId:      carriercommon.Hash{}.Bytes(),
			SenderRole:      uint64(task.GetLocalTaskRole()),
			SenderPartyId:   []byte(task.GetLocalTaskOrganization().GetPartyId()),
			ReceiverRole:    uint64(task.GetRemoteTaskRole()),
			ReceiverPartyId: []byte(task.GetRemoteTaskOrganization().GetPartyId()),
			MsgOwner: &carriernetmsgcommonpb.TaskOrganizationIdentityInfo{
				Name:       []byte(task.GetLocalTaskOrganization().GetNodeName()),
				NodeId:     []byte(task.GetLocalTaskOrganization().GetNodeId()),
				IdentityId: []byte(task.GetLocalTaskOrganization().GetIdentityId()),
				PartyId:    []byte(task.GetLocalTaskOrganization().GetPartyId()),
			},
		},
		TaskEventList: types.ConvertTaskEventArr(eventList),
		CreateAt:      timeutils.UnixMsecUint64(),
		Sign:          nil,
	}
}

func (m *Manager) handleTaskEventWithCurrentOranization(event *carriertypespb.TaskEvent) error {
	if len(event.GetType()) != ev.EventTypeCharLen {
		return ev.IncEventType
	}

	identityId, err := m.resourceMng.GetDB().QueryIdentityId()
	if nil != err {
		log.WithError(err).Errorf("Failed to query self identityId on taskManager.SendTaskEvent()")
		return fmt.Errorf("query local identityId failed, %s", err)
	}
	event.IdentityId = identityId

	if task, ok := m.queryNeedExecuteTaskCache(event.GetTaskId(), event.GetPartyId()); ok {

		// need to validate the task that have been processing ? Maybe~
		// While task is consensus or executing, can terminate.
		has, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusByPartyId(event.GetTaskId(), event.GetPartyId())
		if nil != err {
			log.WithError(err).Errorf("Failed to check local task execute status whether exist on `taskManager.handleTaskEventWithCurrentOranization()`, taskId: {%s}, partyId: {%s}",
				event.GetTaskId(), event.GetPartyId())
			// remove wrong task cache
			m.removeNeedExecuteTaskCache(event.GetTaskId(), event.GetPartyId())
			return err
		}

		if !has {
			log.Warnf("Warn ignore event, `event is the end` but not find party task executeStatus on `taskManager.handleTaskEventWithCurrentOranization()`, event: %s",
				event.String())
			// remove wrong task cache
			m.removeNeedExecuteTaskCache(event.GetTaskId(), event.GetPartyId())
			return nil
		}

		localTask, err := m.resourceMng.GetDB().QueryLocalTask(event.GetTaskId())
		if nil != err {
			log.WithError(err).Errorf("Failed to query local task info on `taskManager.handleTaskEventWithCurrentOranization()`, taskId: {%s}, partyId: {%s}",
				event.GetTaskId(), event.GetPartyId())
			// remove wrong task cache
			m.removeNeedExecuteTaskCache(event.GetTaskId(), event.GetPartyId())
			return err
		}

		return m.executeTaskEvent("on `taskManager.handleTaskEventWithCurrentOranization()`", types.LocalNetworkMsg, event, task, localTask)
	}
	return nil // ignore event while task is not exist.
}

func (m *Manager) handleNeedReplayScheduleTask(task *types.NeedReplayScheduleTask) {

	// Do duplication check ...
	has, err := m.resourceMng.GetDB().HasLocalTask(task.GetTask().GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("Failed to query local task when received remote task, taskId: {%s}", task.GetTask().GetTaskId())
		return
	}

	// There is no need to store local tasks repeatedly
	if !has {

		log.Infof("Start to store local task on taskManager.loop() when received needReplayScheduleTask, taskId: {%s}", task.GetTask().GetTaskId())

		// store metadata used taskId
		if err := m.storeMetadataUsedTaskId(task.GetTask()); nil != err {
			log.WithError(err).Errorf("Failed to store metadata used taskId when received remote task, taskId: {%s}", task.GetTask().GetTaskId())
		}
		if err := m.resourceMng.GetDB().StoreLocalTask(task.GetTask()); nil != err {
			log.WithError(err).Errorf("Failed to call StoreLocalTask when replay schedule remote task, taskId: {%s}", task.GetTask().GetTaskId())
			task.SendFailedResult(task.GetTask().GetTaskId(), err)
			return
		}
		log.Infof("Finished to store local task on taskManager.loop() when received needReplayScheduleTask, taskId: {%s}", task.GetTask().GetTaskId())

	}

	// Start replay schedule remote task ...
	result := m.scheduler.ReplaySchedule(task.GetLocalPartyId(), task.GetLocalTaskRole(), task)
	task.SendResult(result)
}
func (m *Manager) handleNormalNeedExecuteTask(task *types.NeedExecuteTask, localTask *types.Task) {

	log.Debugf("Start handle needExecuteTask on handleNormalNeedExecuteTask(), taskId: {%s}, role: {%s}, partyId: {%s}",
		task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId())

	switch task.GetConsStatus() {
	case types.TaskConsensusFinished:
		// store task consensus result (failed or succeed) event with sender party
		m.resourceMng.GetDB().StoreTaskEvent(&carriertypespb.TaskEvent{
			Type:       ev.TaskSucceedConsensus.GetType(),
			TaskId:     task.GetTaskId(),
			IdentityId: task.GetLocalTaskOrganization().GetIdentityId(),
			PartyId:    task.GetLocalTaskOrganization().GetPartyId(),
			Content:    "succeed consensus",
			CreateAt:   timeutils.UnixMsecUint64(),
		})

		// the task that are 'consensusFinished' status should be removed from the queue/starvequeue as task sender.
		if task.GetLocalTaskRole() == commonconstantpb.TaskRole_TaskRole_Sender {

			if err := m.scheduler.RemoveTask(task.GetTaskId()); nil != err {
				log.WithError(err).Errorf("Failed to remove local task from queue/starve queue when received `consensusFinished` needExecuteTask of `NEED-CONSENSUS` task consensus result from [consensus engine], taskId: {%s}, role: {%s}, partyId: {%s}",
					task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId())
			}
		}
	}

	// Store task exec status
	if err := m.StoreExecuteTaskStateBeforeExecuteTask("on handleNormalNeedExecuteTask()", task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId()); nil != err {
		log.WithError(err).Errorf("Failed to call StoreExecuteTaskStateBeforeExecuteTask() on handleNormalNeedExecuteTask(), taskId: {%s}, role: {%s}, partyId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId())
		return
	}

	if localTask.GetTaskData().GetState() == commonconstantpb.TaskState_TaskState_Pending &&
		localTask.GetTaskData().GetStartAt() == 0 {

		localTask.GetTaskData().State = commonconstantpb.TaskState_TaskState_Running
		localTask.GetTaskData().StartAt = timeutils.UnixMsecUint64()
		if err := m.resourceMng.GetDB().StoreLocalTask(localTask); nil != err {
			log.WithError(err).Errorf("Failed to update local task state before executing task on handleNormalNeedExecuteTask(), taskId: {%s}, need update state: {%s}",
				task.GetTaskId(), commonconstantpb.TaskState_TaskState_Running.String())
		}
	}

	// store local cache
	m.addNeedExecuteTaskCacheAndminotor(task, int64(localTask.GetTaskData().GetStartAt()+localTask.GetTaskData().GetOperationCost().GetDuration()))

	// The task sender will not execute the task
	// driving task to executing
	if err := m.driveTaskForExecute(task, localTask); nil != err {
		log.WithError(err).Errorf("Failed to execute task on internal node on handleNormalNeedExecuteTask(), taskId: {%s}, role: {%s}, partyId: {%s}",
			task.GetTaskId(), task.GetLocalTaskRole().String(), task.GetLocalTaskOrganization().GetPartyId())

		m.SendTaskEvent(m.eventEngine.GenerateEvent(ev.TaskExecuteFailedEOF.GetType(), task.GetTaskId(), task.GetLocalTaskOrganization().GetIdentityId(),
			task.GetLocalTaskOrganization().GetPartyId(), fmt.Sprintf("%s, %s with %s", ev.TaskExecuteFailedEOF.GetMsg(), err,
				task.GetLocalTaskOrganization().GetPartyId())))
	}
}
func (m *Manager) handleAbnormalNeedExecuteTask(task *types.NeedExecuteTask, localTask *types.Task) {
	// store a bad event into local db before handle bad task.
	var (
		eventTyp string
		reason   string
		done bool
	)

	switch task.GetConsStatus() {

	// The consensus cycle is interrupted in the process of consensus.
	case types.TaskConsensusInterrupt:

		// check if repush task into queue/starvequeue by task sender.
		if task.GetLocalTaskRole() == commonconstantpb.TaskRole_TaskRole_Sender {

			log.Debugf("Received `consensusInterrupt` needExecuteTask of `NEED-CONSENSUS` task consensus result from [consensus engine], taskId: {%s}, needExecuteTask: {%s}", task.GetTaskId(), task.String())

			// clean old powerSuppliers and update local task
			localTask.RemovePowerSuppliers()
			localTask.RemovePowerResources()
			// restore task by power
			if err := m.resourceMng.GetDB().StoreLocalTask(localTask); nil != err {
				log.WithError(err).Errorf("Failed to update local task whit clean powers after consensus interrupted when received needExecuteTask of `NEED-CONSENSUS` task consensus result from [consensus engine], taskId: {%s}", task.GetTaskId())
			}

			// repush task into queue, if anything else
			if err := m.scheduler.RepushTask(localTask); err == schedule.ErrRescheduleLargeThreshold {
				log.WithError(err).Errorf("Failed to repush local task into queue/starve queue when received needExecuteTask of `NEED-CONSENSUS` task consensus result from [consensus engine], taskId: {%s}, consensus status: {%s}",
					task.GetTaskId(), task.GetConsStatus().String())

				m.scheduler.RemoveTask(task.GetTaskId())

				eventTyp = ev.TaskFailed.GetType()
				reason = fmt.Sprintf("consensus interrupted: "+task.GetErr().Error()+" and "+schedule.ErrRescheduleLargeThreshold.Error())
				done = true

			} else {
				log.Debugf("Succeed to repush local task into queue/starve queue when received needExecuteTask of `NEED-CONSENSUS` task consensus result from [consensus engine], taskId: {%s}, consensus status: {%s}",
					task.GetTaskId(), task.GetConsStatus().String())

				// When the repush task is successfully sent to the queue,
				// we need to clear the task sender's neexexecutetask
				// so that the task sender can schedule the task later.
				m.removeNeedExecuteTaskCache(task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId())

				eventTyp = ev.TaskFailedConsensus.GetType()
				reason = task.GetErr().Error()
			}
		} else {
			eventTyp = ev.TaskFailedConsensus.GetType()
			reason = task.GetErr().Error()
			done = true
		}

	// May have failed before the consensus started
	case types.TaskScheduleFailed:

		if task.GetLocalTaskRole() == commonconstantpb.TaskRole_TaskRole_Sender {
			log.Debugf("Received `taskScheduleFailed` needExecuteTask of `NEED-CONSENSUS` task consensus result from [consensus engine], taskId: {%s}, needExecuteTask: {%s}",
				task.GetTaskId(), task.String())
		}

		eventTyp = ev.TaskFailed.GetType()
		reason = task.GetErr().Error()
		done = true

	case types.TaskTerminate:

		// the task that are 'terminated' status should be removed from the queue/starvequeue as task sender.
		if task.GetLocalTaskRole() == commonconstantpb.TaskRole_TaskRole_Sender {

			log.Debugf("Received `taskTerminate` needExecuteTask of `NEED-CONSENSUS` task consensus result from [consensus engine], taskId: {%s}, needExecuteTask: {%s}",
				task.GetTaskId(), task.String())

			if err := m.scheduler.RemoveTask(task.GetTaskId()); nil != err {
				log.WithError(err).Errorf("Failed to remove local task from queue/starve queue when received needExecuteTask of `NEED-CONSENSUS` task consensus result from [consensus engine],taskId: {%s}, partyId: {%s}, status: {%s}",
					task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), task.GetConsStatus().String())
			}
		}

		// Both task sender and task partners will end tasks in 'terminate' status.
		eventTyp = ev.TaskFailed.GetType()
		reason = task.GetErr().Error()
		done = true

	}

	// store event.
	m.resourceMng.GetDB().StoreTaskEvent(m.eventEngine.GenerateEvent(eventTyp, task.GetTaskId(), task.GetLocalTaskOrganization().GetIdentityId(),
		task.GetLocalTaskOrganization().GetPartyId(), reason))

	// Finally handle the fate of needexecutetask
	if done {
		switch task.GetLocalTaskRole() {
		case commonconstantpb.TaskRole_TaskRole_Sender:
			m.publishFinishedTaskToDataCenter(task, localTask, true)
		default:
			m.sendTaskResultMsgToTaskSender(task, localTask)
		}
		// finally, remove needExecuteTask after call `publishFinishedTaskToDataCenter` or `sendTaskResultMsgToTaskSender`
		m.removeNeedExecuteTaskCache(task.GetTaskId(), task.GetLocalTaskOrganization().GetIdentityId())
	}
}
func (m *Manager) handleNeedExecuteTask(task *types.NeedExecuteTask) {

	localTask, err := m.resourceMng.GetDB().QueryLocalTask(task.GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("Failed to query local task info on taskManager.loop() when received needExecuteTask, taskId: {%s}, partyId: {%s}, status: {%s}",
			task.GetTaskId(), task.GetLocalTaskOrganization().GetPartyId(), task.GetConsStatus().String())
		//continue
		return
	}

	// init consumeSpec of NeedExecuteTask first (by v0.4.0)
	m.initConsumeSpecByConsumeOption(task)

	switch task.GetConsStatus() {
	// sender and partner to handle needExecuteTask when consensus succeed.
	// sender need to store some cache, partner need to execute task.
	case types.TaskConsensusFinished:

		// to execute the task
		m.handleNormalNeedExecuteTask(task, localTask)

	// sender and partner to clear local task things after received status:
	// #### `scheduleFailed` ####
	// #### `interrupt` ####
	// #### `terminate` ####.
	//
	// sender need to publish local task and event to datacenter,
	// partner need to send task's event to remote task's sender.
	//
	// #### NOTE: ####
	//
	//	sender never received status `interrupt`,
	//	and partners never received status `scheduleFailed`.
	default:
		m.handleAbnormalNeedExecuteTask(task, localTask)
	}
}

func (m *Manager) executeTaskEvent(logkeyword string, symbol types.NetworkMsgLocationSymbol, event *carriertypespb.TaskEvent, localNeedtask *types.NeedExecuteTask, localTask *types.Task) error {

	if err := m.resourceMng.GetDB().StoreTaskEvent(event); nil != err {
		log.WithError(err).Errorf("Failed to store %s taskEvent %s, event: %s", symbol.String(), logkeyword, event.String())
	} else {
		log.Infof("Started store %s taskEvent %s, event: %s", symbol.String(), logkeyword, event.String())
	}

	switch event.GetType() {
	case ev.TaskExecuteSucceedEOF.GetType():

		log.Infof("Started handle taskEvent with currentIdentity, `event is the task final succeed EOF finished` %s, event: %s", logkeyword, event.String())

		// #### NOTE: ####
		//
		// The type '0100005' has been generated when the fighter reports to the carrier,
		// but it is carried with it when it is sent to the task sender
		// (on the contrary, the remote message received by the task sender has already carried the type '0100005')
		if symbol == types.LocalNetworkMsg {
			m.resourceMng.GetDB().StoreTaskEvent(m.eventEngine.GenerateEvent(ev.TaskSucceed.GetType(), event.GetTaskId(),
				event.GetIdentityId(), event.GetPartyId(), "task execute succeed"))
		}

		// ### NOTE ###
		//
		// remove current partyId by local with taskId
		// then check tail partyIds count if count value is zero
		// (When the task sender is not same one with the identityid of the current organization,
		// this function will return `false`)
		publish, err := m.checkTaskSenderPublishOpportunity(localTask, event)
		if nil != err {
			log.WithError(err).Errorf("Failed to check task sender publish opportunity %s, event: %s",
				logkeyword, event.GetPartyId())
			return err
		}

		// 1、 handle last party
		//   send this task result to remote target peer,
		//   but they belong to same organization, call local msg.
		if symbol == types.LocalNetworkMsg {
			m.sendTaskResultMsgToTaskSender(localNeedtask, localTask)
			m.removeNeedExecuteTaskCache(event.GetTaskId(), event.GetPartyId())
		}

		if publish {

			// 2、 handle sender party
			senderNeedTask, ok := m.queryNeedExecuteTaskCache(event.GetTaskId(), localTask.GetTaskSender().GetPartyId())
			if ok {
				log.Debugf("Need to call `publishFinishedTaskToDataCenter` %s, taskId: {%s}, sender partyId: {%s}",
					logkeyword, event.GetTaskId(), localTask.GetTaskSender().GetPartyId())
				// handle this task result with current peer
				m.publishFinishedTaskToDataCenter(senderNeedTask, localTask, true)
				m.removeNeedExecuteTaskCache(event.GetTaskId(), localTask.GetTaskSender().GetPartyId())
			}
		}
	case ev.TaskExecuteFailedEOF.GetType():

		log.Infof("Started handle taskEvent with currentIdentity, `event is the task final [FAILED] EOF finished` will terminate task %s, event: %s", logkeyword, event.String())

		if symbol == types.LocalNetworkMsg {

			// #### NOTE: ####
			//
			// The type '0100004' has been generated when the fighter reports to the carrier,
			// but it is carried with it when it is sent to the task sender
			// (on the contrary, the remote message received by the task sender has already carried the type '0100004')
			if symbol == types.LocalNetworkMsg {
				m.resourceMng.GetDB().StoreTaskEvent(m.eventEngine.GenerateEvent(ev.TaskFailed.GetType(), event.GetTaskId(),
					event.GetIdentityId(), event.GetPartyId(), "task execute failed"))
			}

			// ### NOTE ###
			//
			// When processing local events, regardless of the partyid of the current needexecutetask,
			// as long as' 0008001 'appears,
			// you should check whether the current partyid belongs to the same organization as the sender.
			// If it belongs to the same organization,
			// 		we need the sender to terminate the task and broadcast the' terminate 'message;
			// If it does not belong to the same organization,
			// 		we only need to terminate the needexecutetask<terminate> message of the current partyid.
			//
			//
			// First, verify whether the current event partyId and task sender partyId belong to the same organization.
			// (that is, whether there is a `needexecutetask` for task sender).
			if m.hasNeedExecuteTaskCache(localTask.GetTaskId(), localTask.GetTaskSender().GetPartyId()) {

				// remove all partyIds by local with taskId
				if err := m.resourceMng.GetDB().RemoveTaskPartnerPartyIds(event.GetTaskId()); nil != err {
					log.WithError(err).Errorf("Failed to remove all partyId of local task's partner arr %s, taskId: {%s}",
						logkeyword, event.GetTaskId())
				}

				if err := m.terminateExecuteTaskBySender(localTask); nil != err {
					log.WithError(err).Errorf("Failed to call `terminateExecuteTaskBySender()` %s, taskId: {%s}", logkeyword, localTask.GetTaskId())
				}
			} else {
				if needExecuteTask, ok := m.queryNeedExecuteTaskCache(localTask.GetTaskId(), event.GetPartyId()); ok {
					if err := m.terminateExecuteTaskByPartner(localTask, needExecuteTask); nil != err {
						log.WithError(err).Errorf("Failed to call `terminateExecuteTaskByPartner()` %s, taskId: {%s}, partyId: {%s}", logkeyword, localTask.GetTaskId(), event.GetPartyId())
					} else {
						log.Debugf("Finished [terminate task] when received report local event, taskId: {%s}, partyId: {%s}", localTask.GetTaskId(), event.GetPartyId())
					}
				}
			}
		}
	}
	return nil
}

func (m *Manager) storeMetadataUsedTaskId(task *types.Task) error {
	identityId, err := m.resourceMng.GetDB().QueryIdentityId()
	if nil != err {
		return err
	}
	for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
		if dataSupplier.GetIdentityId() == identityId {

			metadataId, err := policy.FetchMetedataIdByPartyIdFromDataPolicy(dataSupplier.GetPartyId(), task.GetTaskData().GetDataPolicyTypes(), task.GetTaskData().GetDataPolicyOptions())
			if nil != err {
				return fmt.Errorf("not fetch metadataId from task dataPolicy, %s, partyId: {%s}", err, dataSupplier.GetPartyId())
			}
			if err := m.resourceMng.GetDB().StoreMetadataHistoryTaskId(metadataId, task.GetTaskId()); nil != err {
				return err
			}
		}
	}
	return nil
}

func (m *Manager) checkTaskSenderPublishOpportunity(task *types.Task, event *carriertypespb.TaskEvent) (bool, error) {

	if "" == strings.Trim(event.GetPartyId(), "") {
		log.Errorf("Failed to check partyId of event, partyId is empty on taskManager.checkTaskSenderPublishOpportunity(), event: %s", event.String())
		return false, fmt.Errorf("empty partyId of event")
	}

	identityId, err := m.resourceMng.GetDB().QueryIdentityId()
	if nil != err {
		log.WithError(err).Errorf("Failed to query self identityId on taskManager.checkTaskSenderPublishOpportunity()")
		return false, fmt.Errorf("query local identityId failed, %s", err)
	}

	// Collect the task result things of other organization, When the current organization is the task sender
	if task.GetTaskSender().GetIdentityId() != identityId {
		return false, nil
	}

	// Remove the currently processed partyId from the partyIds array of the task partner to be processed
	if err := m.resourceMng.GetDB().RemoveTaskPartnerPartyId(event.GetTaskId(), event.GetPartyId()); nil != err {
		log.WithError(err).Errorf("Failed to remove partyId of local task's partner arr on `taskManager.checkTaskSenderPublishOpportunity()`, taskId: {%s}, partyId: {%s}",
			event.GetTaskId(), event.GetPartyId())
	}
	// has, err := m.resourceMng.GetDB().HasTaskPartnerPartyIds(event.GetTaskId())
	partyIds, err := m.resourceMng.GetDB().QueryTaskPartnerPartyIds(event.GetTaskId())
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("Failed to query task partner partyIds on `taskManager.checkTaskSenderPublishOpportunity()`, taskId: {%s}",
			event.GetTaskId())
		return false, err
	}

	log.Debugf("Query partyIds on `taskManager.checkTaskSenderPublishOpportunity()`, current partyId: {%s}, partyIds: %s",
		event.GetPartyId(), "["+strings.Join(partyIds, ",")+"]")

	if rawdb.IsDBNotFoundErr(err) {
		return true, nil
	}

	// if has {
	if len(partyIds) != 0 {
		return false, nil
	}
	return true, nil
}

func (m *Manager) handleResourceUsage(keyword, usageIdentityId string, usage *types.TaskResuorceUsage, localTask *types.Task, nmls types.NetworkMsgLocationSymbol) (bool, error) {

	var (
		partyId string
	)

	if nmls == types.LocalNetworkMsg {
		partyId = usage.GetPartyId()
	} else {
		partyId = localTask.GetTaskSender().GetPartyId()
	}

	// ## 1、 check whether task status is terminate (with party self | with task sender) ?
	terminating, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusTerminateByPartyId(usage.GetTaskId(), partyId)
	if nil != err {
		log.WithError(err).Errorf("Failed to call HasLocalTaskExecuteStatusTerminateByPartyId() on taskManager.handleResourceUsage() %s, taskId: {%s}, partyId: {%s}, remote partyId: {%s}",
			keyword, usage.GetTaskId(), partyId, usage.GetPartyId())
		return false, fmt.Errorf("check current party has `terminate` status needExecuteTask failed, %s", err)
	}
	if terminating {
		log.Warnf("The localTask execute status has `terminate` on taskManager.handleResourceUsage() %s, taskId: {%s}, partyId: {%s}, remote partyId: {%s}",
			keyword, usage.GetTaskId(), partyId, usage.GetPartyId())
		return false, fmt.Errorf("task was terminated")
	}

	// ## 2、 check whether task status is running (with party self | with task sender) ?
	running, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusRunningByPartyId(usage.GetTaskId(), partyId)
	if nil != err {
		log.WithError(err).Errorf("Failed to call HasLocalTaskExecuteStatusRunningByPartyId() on taskManager.handleResourceUsage() %s, taskId: {%s}, partyId: {%s}, remote partyId: {%s}",
			keyword, usage.GetTaskId(), partyId, usage.GetPartyId())
		return false, fmt.Errorf("check current party has `running` status needExecuteTask failed, %s", err)
	}
	if !running {
		log.Warnf("Not found localTask execute status `running` on taskManager.handleResourceUsage() %s, taskId: {%s}, partyId: {%s}, remote partyId: {%s}",
			keyword, usage.GetTaskId(), partyId, usage.GetPartyId())
		return false, fmt.Errorf("task is not executed")
	}

	var needUpdate bool

	// 1、collected the powerSuppliers into cache
	powerCache := make(map[string]*carriertypespb.TaskOrganization)
	for _, power := range localTask.GetTaskData().GetPowerSuppliers() {
		powerCache[power.GetPartyId()] = power
	}

	// 2、update the power resource with partyId
	for i, powerOption := range localTask.GetTaskData().GetPowerResourceOptions() {

		power, ok := powerCache[powerOption.GetPartyId()]
		if ok {
			// find power supplier info by identity and partyId with msg from reomte peer
			// (find the target power supplier, it maybe local power supplier or remote power supplier)
			// and update its' resource usage info.
			if usage.GetPartyId() == powerOption.GetPartyId() &&
				usageIdentityId == power.GetIdentityId() {

				resourceUsage := localTask.GetTaskData().GetPowerResourceOptions()[i].GetResourceUsedOverview()
				// update ...
				if usage.GetUsedMem() > resourceUsage.GetUsedMem() {
					if usage.GetUsedMem() > localTask.GetTaskData().GetOperationCost().GetMemory() {
						resourceUsage.UsedMem = localTask.GetTaskData().GetOperationCost().GetMemory()
					} else {
						resourceUsage.UsedMem = usage.GetUsedMem()
					}
					needUpdate = true
				}
				if usage.GetUsedProcessor() > resourceUsage.GetUsedProcessor() {
					if usage.GetUsedProcessor() > localTask.GetTaskData().GetOperationCost().GetProcessor() {
						resourceUsage.UsedProcessor = localTask.GetTaskData().GetOperationCost().GetProcessor()
					} else {
						resourceUsage.UsedProcessor = usage.GetUsedProcessor()
					}
					needUpdate = true
				}
				if usage.GetUsedBandwidth() > resourceUsage.GetUsedBandwidth() {
					if usage.GetUsedBandwidth() > localTask.GetTaskData().GetOperationCost().GetBandwidth() {
						resourceUsage.UsedBandwidth = localTask.GetTaskData().GetOperationCost().GetBandwidth()
					} else {
						resourceUsage.UsedBandwidth = usage.GetUsedBandwidth()
					}
					needUpdate = true
				}
				if usage.GetUsedDisk() > resourceUsage.GetUsedDisk() {
					resourceUsage.UsedDisk = usage.GetUsedDisk()
					needUpdate = true
				}
				// update ...
				localTask.GetTaskData().GetPowerResourceOptions()[i].ResourceUsedOverview = resourceUsage
			}
		}
	}

	if needUpdate {
		log.Debugf("Need to update local task on taskManager.handleResourceUsage() %s, usage: %s", keyword, usage.String())

		// Updata task when resourceUsed change.
		if err := m.resourceMng.GetDB().StoreLocalTask(localTask); nil != err {
			log.WithError(err).Errorf("Failed to call StoreLocalTask() on taskManager.handleResourceUsage() %s, taskId: {%s}, partyId: {%s}, remote partyId: {%s}",
				keyword, usage.GetTaskId(), partyId, usage.GetPartyId())
			return false, fmt.Errorf("update local task by usage change failed, %s", err)
		}
	}

	return needUpdate, nil
}

func (m *Manager) ValidateTaskResultMsg(pid peer.ID, taskResultMsg *carriernetmsgtaskmngpb.TaskResultMsg) error {
	msg := types.FetchTaskResultMsg(taskResultMsg) // fetchTaskResultMsg(taskResultMsg)

	if len(msg.GetTaskEventList()) == 0 {
		return nil
	}

	taskId := msg.GetTaskEventList()[0].GetTaskId()

	for _, event := range msg.GetTaskEventList() {
		if taskId != event.GetTaskId() {
			return fmt.Errorf("received event failed, has invalid taskId: {%s}, right taskId: {%s}", event.GetTaskId(), taskId)
		}
	}

	return nil
}

func (m *Manager) OnTaskResultMsg(pid peer.ID, taskResultMsg *carriernetmsgtaskmngpb.TaskResultMsg) error {

	if err := m.ContainsOrAddMsg(taskResultMsg); nil != err {
		return err
	}

	msg := types.FetchTaskResultMsg(taskResultMsg)

	// Verify the signature
	if _, err := signutil.VerifyMsgSign(msg.GetMsgOption().GetOwner().GetNodeId(), msg.Hash().Bytes(), msg.GetSign()); err != nil {
		log.WithError(err).Errorf("Failed to call `VerifyMsgSign()` when received taskResultMsg on `taskManager.OnTaskResultMsg()`, proposalId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetMsgOption().GetProposalId().String(), msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("verify taskResultMsg sign %s", err)
	}

	if len(msg.GetTaskEventList()) == 0 {
		return nil
	}

	taskId := msg.GetTaskEventList()[0].GetTaskId()

	if _, ok := m.queryNeedExecuteTaskCache(taskId, msg.GetMsgOption().GetReceiverPartyId()); !ok {
		log.Errorf("Failed to query needExecuteTask when received taskResultMsg on `taskManager.OnTaskResultMsg()`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetMsgOption().GetProposalId().String(), taskId, msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("cannot query needExecuteTask by msg receiver partyId")
	}

	// While task is consensus or executing, handle task resultMsg.  (In this case, the receiver must be task sender)
	has, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusByPartyId(taskId, msg.GetMsgOption().GetReceiverPartyId())
	if nil != err {
		log.WithError(err).Errorf("Failed to query local task executing status when received taskResultMsg on `taskManager.OnTaskResultMsg()`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetMsgOption().GetProposalId().String(), taskId, msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("query local task executing status failed")
	}

	if !has {
		log.Warnf("Not found local task executing status when received taskResultMsg on `taskManager.OnTaskResultMsg()`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetMsgOption().GetProposalId().String(), taskId, msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		//// remove wrong task cache
		//m.removeNeedExecuteTaskCache(taskId, msg.GetMsgOption().GetReceiverPartyId())
		return nil
	}

	localTask, err := m.resourceMng.GetDB().QueryLocalTask(taskId)
	if nil != err {
		log.WithError(err).Errorf("Failed to call `QueryLocalTask()` when received taskResultMsg on `taskManager.OnTaskResultMsg()`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetMsgOption().GetProposalId().String(), taskId, msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("query local task failed, %s", err)
	}

	receiver := fetchOrgByPartyRole(msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetReceiverRole(), localTask)
	identity, err := m.resourceMng.GetDB().QueryIdentity()
	if nil != err {
		log.WithError(err).Errorf("Failed to call `QueryIdentity()` when received taskResultMsg on `taskManager.OnTaskResultMsg()`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetMsgOption().GetProposalId().String(), taskId, msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("query local identity failed, %s", err)
	}
	// verify the receiver is myself ?
	if identity.GetIdentityId() != receiver.GetIdentityId() {
		log.Warnf("Warning verify receiver identityId of taskResultMsg, receiver is not me when received taskResultMsg on `taskManager.OnTaskResultMsg()`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, remote role: {%s}, remote partyId: {%s}",
			msg.GetMsgOption().GetProposalId().String(), taskId, msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderRole().String(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("receiver is not me of taskResultMsg")
	}

	/**
	+++++++++++++++++++++++++++++++++++++++++++
	+++++++++++++++++++++++++++++++++++++++++++
	NOTE: receiverPartyId must be task sender partyId.
	+++++++++++++++++++++++++++++++++++++++++++
	+++++++++++++++++++++++++++++++++++++++++++
	*/
	if localTask.GetTaskSender().GetIdentityId() != receiver.GetIdentityId() ||
		localTask.GetTaskSender().GetPartyId() != receiver.GetPartyId() {
		log.Warnf("Warning the receiver of taskResultMsg and sender of task is not the same when received taskResultMsg on `taskManager.OnTaskResultMsg()`, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, msg receiver: %s, task sender: %s",
			msg.GetMsgOption().GetProposalId().String(), localTask.GetTaskId(), msg.GetMsgOption().GetReceiverRole().String(), msg.GetMsgOption().GetReceiverPartyId(), receiver.String(), localTask.GetTaskSender().String())
		return fmt.Errorf("receiver is not task sender of taskResultMsg")
	}

	log.WithField("traceId", traceutil.GenerateTraceID(taskResultMsg)).Debugf("Received remote taskResultMsg on `taskManager.OnTaskResultMsg()`, remote pid: {%s}, taskId: {%s}, taskResultMsg: %s", pid, taskId, msg.String())

	for _, event := range msg.GetTaskEventList() {

		if "" == strings.Trim(event.GetPartyId(), "") || msg.GetMsgOption().GetSenderPartyId() != strings.Trim(event.GetPartyId(), "") {
			continue
		}

		if err := m.executeTaskEvent("on `taskManager.OnTaskResultMsg()`", types.RemoteNetworkMsg, event, nil, localTask); nil != err {
			return err
		}
	}
	return nil
}

func (m *Manager) ValidateTaskResourceUsageMsg(pid peer.ID, taskResourceUsageMsg *carriernetmsgtaskmngpb.TaskResourceUsageMsg) error {
	return nil
}

func (m *Manager) OnTaskResourceUsageMsg(pid peer.ID, usageMsg *carriernetmsgtaskmngpb.TaskResourceUsageMsg) error {
	return m.onTaskResourceUsageMsg(pid, usageMsg)
}

func (m *Manager) onTaskResourceUsageMsg(pid peer.ID, usageMsg *carriernetmsgtaskmngpb.TaskResourceUsageMsg) error {

	if err := m.ContainsOrAddMsg(usageMsg); nil != err {
		return err
	}

	msg := types.FetchTaskResourceUsageMsg(usageMsg)

	if msg.GetMsgOption().GetSenderPartyId() != msg.GetUsage().GetPartyId() {
		log.Errorf("msg sender partyId AND usageMsg partyId is not same when received taskResourceUsageMsg, taskId: {%s}, partyId: {%s}, remote partyId: {%s}, usagemsgPartyId: {%s}",
			msg.GetUsage().GetTaskId(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderPartyId(), msg.GetUsage().GetPartyId())
		return fmt.Errorf("invalid partyId of usageMsg")
	}

	// Verify the signature
	if _, err := signutil.VerifyMsgSign(msg.GetMsgOption().GetOwner().GetNodeId(), msg.Hash().Bytes(), msg.GetSign()); err != nil {
		log.WithError(err).Errorf("Failed to call `VerifyMsgSign()` when received taskResourceUsageMsg, taskId: {%s}, partyId: {%s}, remote partyId: {%s}, usagemsgPartyId: {%s}",
			msg.GetUsage().GetTaskId(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderPartyId(), msg.GetUsage().GetPartyId())
		return fmt.Errorf("verify usageMsg sign %s", err)
	}

	// Note: the needexecutetask obtained here is generally the needexecutetask of the task sender,
	//       so the remote organization is also the task sender.
	_, ok := m.queryNeedExecuteTaskCache(msg.GetUsage().GetTaskId(), msg.GetMsgOption().GetReceiverPartyId())
	if !ok {
		log.Warnf("Not found needExecuteTask when received taskResourceUsageMsg, taskId: {%s}, partyId: {%s}, remote partyId: {%s}",
			msg.GetUsage().GetTaskId(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("can not find `need execute task` cache")
	}

	// Update task resourceUsed of powerSuppliers of local task
	localTask, err := m.resourceMng.GetDB().QueryLocalTask(msg.GetUsage().GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("Failed to call `QueryLocalTask()` when received taskResourceUsageMsg, taskId: {%s}, partyId: {%s}, remote partyId: {%s}",
			msg.GetUsage().GetTaskId(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("query local task failed, %s", err)
	}

	receiver := fetchOrgByPartyRole(msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetReceiverRole(), localTask)
	identity, err := m.resourceMng.GetDB().QueryIdentity()
	if nil != err {
		log.WithError(err).Errorf("Failed to call `QueryIdentity()` when received taskResourceUsageMsg, taskId: {%s}, partyId: {%s}, remote partyId: {%s}",
			msg.GetUsage().GetTaskId(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("query local identity failed, %s", err)
	}
	// verify the receiver is myself ?
	if identity.GetIdentityId() != receiver.GetIdentityId() {
		log.Warnf("Warning verify receiver identityId of taskResourceUsageMsg, receiver is not me when received taskResourceUsageMsg, taskId: {%s}, partyId: {%s}, remote partyId: {%s}",
			msg.GetUsage().GetTaskId(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("receiver is not me of taskResourceUsageMsg")
	}

	/**
	+++++++++++++++++++++++++++++++++++++++++++
	+++++++++++++++++++++++++++++++++++++++++++
	NOTE: receiverPartyId must be task sender partyId.
	+++++++++++++++++++++++++++++++++++++++++++
	+++++++++++++++++++++++++++++++++++++++++++
	*/
	if localTask.GetTaskSender().GetIdentityId() != receiver.GetIdentityId() ||
		localTask.GetTaskSender().GetPartyId() != receiver.GetPartyId() {
		log.Warnf("Warning the receiver of the usageMsg and sender of task is not the same when received taskResourceUsageMsg, taskId: {%s}, partyId: {%s}, remote partyId: {%s}",
			msg.GetUsage().GetTaskId(), msg.GetMsgOption().GetReceiverPartyId(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("receiver is not task sender of taskResourceUsageMsg")
	}

	log.WithField("traceId", traceutil.GenerateTraceID(usageMsg)).Debugf("Received remote taskResourceUsageMsg when received taskResourceUsageMsg, remote pid: {%s}, taskResourceUsageMsg: %s", pid, msg.String())

	needUpdate, err := m.handleResourceUsage("when received remote resourceUsage", msg.GetMsgOption().GetOwner().GetIdentityId(), msg.GetUsage(), localTask, types.RemoteNetworkMsg)
	if nil != err {
		return err
	}

	if needUpdate {
		log.Debugf("Succeed handle remote resourceUsage when received taskResourceUsageMsg, remote pid: {%s}, taskResourceUsageMsg: %s", pid, msg.String())
	}
	return nil
}

func (m *Manager) ValidateTaskTerminateMsg(pid peer.ID, terminateMsg *carriernetmsgtaskmngpb.TaskTerminateMsg) error {
	return nil
}

func (m *Manager) OnTaskTerminateMsg(pid peer.ID, terminateMsg *carriernetmsgtaskmngpb.TaskTerminateMsg) error {
	return m.onTaskTerminateMsg(pid, terminateMsg, types.RemoteNetworkMsg)
}

func (m *Manager) onTaskTerminateMsg(pid peer.ID, terminateMsg *carriernetmsgtaskmngpb.TaskTerminateMsg, nmls types.NetworkMsgLocationSymbol) error {

	if err := m.ContainsOrAddMsg(terminateMsg); nil != err {
		return err
	}

	msg := types.FetchTaskTerminateTaskMngMsg(terminateMsg)

	// Verify the signature of remote msg
	if nmls == types.RemoteNetworkMsg {
		if _, err := signutil.VerifyMsgSign(msg.GetMsgOption().GetOwner().GetNodeId(), msg.Hash().Bytes(), msg.GetSign()); err != nil {
			log.WithError(err).Errorf("Failed to call `VerifyMsgSign()` when received terminateMsg on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, sender partyId: {%s}",
				msg.GetTaskId(), msg.GetMsgOption().GetSenderPartyId())
			return fmt.Errorf("verify remote terminateMsg sign %s", err)
		}
	}

	identity, err := m.resourceMng.GetDB().QueryIdentity()
	if nil != err {
		log.WithError(err).Errorf("Failed to call `QueryIdentity()` when received terminateMsg on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, sender partyId: {%s}",
			msg.GetTaskId(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("query local identity failed when received terminateMsg, %s", err)
	}

	task, err := m.resourceMng.GetDB().QueryLocalTask(msg.GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("Failed to query local task when received terminateMsg on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, sender partyId: {%s}",
			msg.GetTaskId(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("query local task failed, %s", err)
	}
	if nil == task {
		log.Errorf("Not found local task when received terminateMsg on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, sender partyId: {%s}",
			msg.GetTaskId(), msg.GetMsgOption().GetSenderPartyId())
		return err
	}

	sender := fetchOrgByPartyRole(msg.GetMsgOption().GetSenderPartyId(), msg.GetMsgOption().GetSenderRole(), task)
	if nil == sender {
		log.Errorf("Failed to check sender of msg when received terminateMsg on `taskManager.OnTaskTerminateMsg()`, it is empty, taskId: {%s}, sender partyId: {%s}",
			msg.GetTaskId(), msg.GetMsgOption().GetSenderPartyId())
		return fmt.Errorf("%s when received terminateMsg", ctypes.ErrConsensusMsgInvalid)
	}

	// Check whether the sender of the message is the same organization as the sender of the task.
	// If not, this message is illegal.
	if task.GetTaskSender().GetIdentityId() != sender.GetIdentityId() ||
		task.GetTaskSender().GetPartyId() != sender.GetPartyId() {
		log.Warnf("Warning the sender of the msg is not the same organization as the sender of the task when received terminateMsg on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, msg sender: %s, task sender: %s",
			task.GetTaskId(), sender.String(), task.GetTaskSender().String())
		return fmt.Errorf("the sender of the msg and sender of local task is not same when received terminateMsg")
	}

	log.WithField("traceId", traceutil.GenerateTraceID(terminateMsg)).Debugf("Received taskTerminateMsg, consensusSymbol: {%s}, remote pid: {%s}, taskTerminateMsg: %s", nmls.String(), pid, msg.String())

	terminateFn := func(party *carriertypespb.TaskOrganization, role commonconstantpb.TaskRole) (uint32, error) {
		// ## 1、 check whether the task has been terminated

		log.Debugf("Prepare [terminate task] when received taskTerminateMsg on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, partyId: {%s}, consensusSymbol: {%s}", task.GetTaskId(), party.GetPartyId(), nmls.String())

		var consensusFlag uint32

		terminating, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusTerminateByPartyId(task.GetTaskId(), party.GetPartyId())
		if nil != err {
			log.WithError(err).Errorf("Failed to query local task execute `termining` status when received taskTerminateMsg on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, partyId: {%s}",
				task.GetTaskId(), party.GetPartyId())
			return consensusFlag, err
		}
		// If so, we will directly short circuit
		if terminating {
			log.Warnf("Warning query local task execute status has `termining` when received taskTerminateMsg on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, partyId: {%s}",
				task.GetTaskId(), party.GetPartyId())
			return consensusFlag, nil
		}

		// ## 2、 check whether the task is running.
		running, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusRunningByPartyId(task.GetTaskId(), party.GetPartyId())
		if nil != err {
			log.WithError(err).Errorf("Failed to query local task execute `running` status when received taskTerminateMsg on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, partyId: {%s}",
				task.GetTaskId(), party.GetPartyId())
			return consensusFlag, err
		}
		// If it is, we will terminate the task
		if running {
			// ## 2、 check whether the task is running.
			if needExecuteTask, ok := m.queryNeedExecuteTaskCache(task.GetTaskId(), party.GetPartyId()); ok {
				err = m.startTerminateWithNeedExecuteTask(needExecuteTask)
				if nil == err {
					log.Debugf("Finished [terminate task] that is `running` status when received taskTerminateMsg on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, partyId: {%s}, consensusSymbol: {%s}",
						task.GetTaskId(), party.GetPartyId(), nmls.String())
				}
			}
			return consensusFlag, err
		}

		// ## 3、 check whether the task is in consensus

		// While task is consensus or executing, can terminate.
		consensusing, err := m.resourceMng.GetDB().HasLocalTaskExecuteStatusConsensusByPartyId(task.GetTaskId(), party.GetPartyId())
		if nil != err {
			log.WithError(err).Errorf("Failed to query local task execute `cons` status on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, partyId: {%s}",
				task.GetTaskId(), party.GetPartyId())
			return 0, err
		}

		if consensusing {
			return 1, nil
		}

		return 0, nil
	}

	failedPartyIds := make([]string, 0)
	var consensusing uint32

	for _, data := range task.GetTaskData().GetDataSuppliers() {
		if identity.GetIdentityId() == data.GetIdentityId() {
			flag, err := terminateFn(data, commonconstantpb.TaskRole_TaskRole_DataSupplier)
			if nil != err {
				failedPartyIds = append(failedPartyIds, data.GetPartyId())
			}
			consensusing |= flag
		}
	}
	for _, data := range task.GetTaskData().GetPowerSuppliers() {
		if identity.GetIdentityId() == data.GetIdentityId() {
			flag, err := terminateFn(data, commonconstantpb.TaskRole_TaskRole_PowerSupplier)
			if nil != err {
				failedPartyIds = append(failedPartyIds, data.GetPartyId())
			}
			consensusing |= flag
		}
	}

	for _, data := range task.GetTaskData().GetReceivers() {
		if identity.GetIdentityId() == data.GetIdentityId() {
			flag, err := terminateFn(data, commonconstantpb.TaskRole_TaskRole_Receiver)
			if nil != err {
				failedPartyIds = append(failedPartyIds, data.GetPartyId())
			}
			consensusing |= flag
		}
	}
	if len(failedPartyIds) != 0 {
		return fmt.Errorf("terminate task failed by [%s], proposaId: {%s}, taskId: {%s}",
			strings.Join(failedPartyIds, ","), msg.GetMsgOption().GetProposalId(), task.GetTaskId())
	}

	// ## 3、 maybe the task is in consensus

	// interrupt consensus with sender AND send terminateMsg to remote partners
	// OR terminate executing task AND send terminateMsg to remote partners
	if consensusing == 1 {
		if err = m.consensusEngine.OnConsensusMsg(pid, types.NewInterruptMsgWrap(task.GetTaskId(), terminateMsg.GetMsgOption(), terminateMsg.GetCreateAt(), terminateMsg.GetSign())); nil != err {
			log.WithError(err).Errorf("Failed to call `OnConsensusMsg()` on `taskManager.OnTaskTerminateMsg()`, taskId: {%s}, sender partyId: {%s}",
				task.GetTaskId(), task.GetTaskSender().GetPartyId())
			return err
		}
	}

	return nil
}

func (m *Manager) startTerminateWithNeedExecuteTask(needExecuteTask *types.NeedExecuteTask) error {

	// 1、terminate fighter processor for this task with current party
	if err := m.driveTaskForTerminate(needExecuteTask); nil != err {
		log.WithError(err).Errorf("Failed to call driveTaskForTerminate() on `taskManager.startTerminateWithNeedExecuteTask()`, taskId: {%s}, role: {%s}, partyId: {%s}",
			needExecuteTask.GetTaskId(), needExecuteTask.GetLocalTaskRole().String(), needExecuteTask.GetLocalTaskOrganization().GetPartyId())
		return err
	}
	// 2、 remove needExecuteTask cache with current party  #### Need remove after call sendTaskResultMsgToTaskSender
	//m.removeNeedExecuteTaskCache(needExecuteTask.GetTaskId(), needExecuteTask.GetLocalTaskOrganization().GetPartyId())
	// 3、Set the execution status of the task to being terminated`
	if err := m.resourceMng.GetDB().StoreLocalTaskExecuteStatusValTerminateByPartyId(needExecuteTask.GetTaskId(), needExecuteTask.GetLocalTaskOrganization().GetPartyId()); nil != err {
		log.WithError(err).Errorf("Failed to store needExecute task status to `terminate` on `taskManager.startTerminateWithNeedExecuteTask()`, taskId: {%s}, partyId: {%s}",
			needExecuteTask.GetTaskId(), needExecuteTask.GetLocalTaskOrganization().GetPartyId())
	}
	// 5、 send a new needExecuteTask(status: types.TaskTerminate) for terminate with current party
	m.sendNeedExecuteTaskByAction(types.NewNeedExecuteTask(
		"",
		needExecuteTask.GetLocalTaskRole(),
		needExecuteTask.GetRemoteTaskRole(),
		needExecuteTask.GetLocalTaskOrganization(),
		needExecuteTask.GetRemoteTaskOrganization(),
		needExecuteTask.GetTaskId(),
		types.TaskTerminate,
		&types.PrepareVoteResource{},          // zero value
		&carriertwopcpb.ConfirmTaskPeerInfo{}, // zero value
		fmt.Errorf("task was terminated"),
	))
	return nil
}

func (m *Manager) checkNeedExecuteTaskMonitors(now int64, syncCall bool) int64 {
	return m.syncExecuteTaskMonitors.CheckMonitors(now, syncCall)
}

func (m *Manager) needExecuteTaskMonitorTimer() *time.Timer {
	return m.syncExecuteTaskMonitors.Timer()
}

func (m *Manager) AddMsg(msg interface{}) bool {

	var key string
	switch msg.(type) {
	case *carriernetmsgtaskmngpb.TaskResourceUsageMsg:
		pure := msg.(*carriernetmsgtaskmngpb.TaskResourceUsageMsg)
		v := hashutil.Hash(append(taskResourceUsageMsgCacheKeyPrefix, []byte(pure.String())...))
		key = carriercommon.Bytes2Hex(v[0:])

	case *carriernetmsgtaskmngpb.TaskResultMsg:
		pure := msg.(*carriernetmsgtaskmngpb.TaskResultMsg)
		v := hashutil.Hash(append(taskResultMsgCacheKeyPrefix, []byte(pure.String())...))
		key = carriercommon.Bytes2Hex(v[0:])

	case *carriernetmsgtaskmngpb.TaskTerminateMsg:
		pure := msg.(*carriernetmsgtaskmngpb.TaskTerminateMsg)
		v := hashutil.Hash(append(taskTerminateMsgCacheKeyPrefix, []byte(pure.String())...))
		key = carriercommon.Bytes2Hex(v[0:])

		//default:
		//	return false
	}
	if "" != key {
		m.msgCache.Add(key, struct{}{})
		return true
	}
	return false
}

func (m *Manager) ContainMsg(msg interface{}) bool {
	switch msg.(type) {
	case *carriernetmsgtaskmngpb.TaskResourceUsageMsg:
		pure := msg.(*carriernetmsgtaskmngpb.TaskResourceUsageMsg)
		v := hashutil.Hash(append(taskResourceUsageMsgCacheKeyPrefix, []byte(pure.String())...))
		return m.msgCache.Contains(carriercommon.Bytes2Hex(v[0:]))
	case *carriernetmsgtaskmngpb.TaskResultMsg:
		pure := msg.(*carriernetmsgtaskmngpb.TaskResultMsg)
		v := hashutil.Hash(append(taskResultMsgCacheKeyPrefix, []byte(pure.String())...))
		return m.msgCache.Contains(carriercommon.Bytes2Hex(v[0:]))
	case *carriernetmsgtaskmngpb.TaskTerminateMsg:
		pure := msg.(*carriernetmsgtaskmngpb.TaskTerminateMsg)
		v := hashutil.Hash(append(taskTerminateMsgCacheKeyPrefix, []byte(pure.String())...))
		return m.msgCache.Contains(carriercommon.Bytes2Hex(v[0:]))
	default:
		return false
	}
}

// return: ok, evict
func (m *Manager) ContainsOrAddMsg(msg interface{}) error {

	var key string
	switch msg.(type) {
	case *carriernetmsgtaskmngpb.TaskResourceUsageMsg:
		pure := msg.(*carriernetmsgtaskmngpb.TaskResourceUsageMsg)
		v := hashutil.Hash(append(taskResourceUsageMsgCacheKeyPrefix, []byte(pure.String())...))
		key = carriercommon.Bytes2Hex(v[0:])

	case *carriernetmsgtaskmngpb.TaskResultMsg:
		pure := msg.(*carriernetmsgtaskmngpb.TaskResultMsg)
		v := hashutil.Hash(append(taskResultMsgCacheKeyPrefix, []byte(pure.String())...))
		key = carriercommon.Bytes2Hex(v[0:])

	case *carriernetmsgtaskmngpb.TaskTerminateMsg:
		pure := msg.(*carriernetmsgtaskmngpb.TaskTerminateMsg)
		v := hashutil.Hash(append(taskTerminateMsgCacheKeyPrefix, []byte(pure.String())...))
		key = carriercommon.Bytes2Hex(v[0:])

		//default:
		//	has, evict = false, false
	}

	if "" == key {
		return fmt.Errorf("not match msg type")
	}
	if has, _ := m.msgCache.ContainsOrAdd(key, struct{}{}); has {
		return fmt.Errorf("key value already exists in lru cache")
	}
	return nil
}

func fetchOrgByPartyRole(partyId string, role commonconstantpb.TaskRole, task *types.Task) *carriertypespb.TaskOrganization {

	switch role {
	case commonconstantpb.TaskRole_TaskRole_Sender:
		if partyId == task.GetTaskSender().GetPartyId() {
			return task.GetTaskSender()
		}
	case commonconstantpb.TaskRole_TaskRole_DataSupplier:
		for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
			if partyId == dataSupplier.GetPartyId() {
				return dataSupplier
			}
		}
	case commonconstantpb.TaskRole_TaskRole_PowerSupplier:
		for _, powerSupplier := range task.GetTaskData().GetPowerSuppliers() {
			if partyId == powerSupplier.GetPartyId() {
				return powerSupplier
			}
		}
	case commonconstantpb.TaskRole_TaskRole_Receiver:
		for _, receiver := range task.GetTaskData().GetReceivers() {
			if partyId == receiver.GetPartyId() {
				return receiver
			}
		}
	}
	return nil
}
