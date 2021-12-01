package twopc

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	"github.com/RosettaFlow/Carrier-Go/core/evengine"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	twopcpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/consensus/twopc"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/peer"
	"strings"
	"sync"
)

func (t *Twopc) removeOrgProposalStateAndTask(proposalId common.Hash, partyId string) {
	if state := t.state.GetProposalState(proposalId); state.IsNotEmpty() {
		log.Infof("Start remove org proposalState and task cache on Consensus, proposalId {%s}, taskId {%s}, partyId: {%s}", proposalId, state.GetTaskId(), partyId)
		t.state.RemoveOrgProposalStateAnyCache(proposalId, state.GetTaskId(), partyId) // remove task/proposal state/prepare vote/ confirm vote and wal with partyId
	} else {
		log.Infof("Start remove confirm taskPeerInfo when proposalState is empty, proposalId {%s}, final remove partyId: {%s}", proposalId, partyId)
		t.state.RemoveConfirmTaskPeerInfo(proposalId) // remove confirm peers and wal
	}
}

func (t *Twopc) storeOrgProposalState(proposalId common.Hash, taskId string, sender *apicommonpb.TaskOrganization, orgState *ctypes.OrgProposalState) {
	pstate := t.state.GetProposalState(proposalId)

	var first bool
	if pstate.IsEmpty() {
		pstate = ctypes.NewProposalState(proposalId, taskId, sender)
		first = true
	}
	pstate.StoreOrgProposalState(orgState)
	t.wal.StoreOrgProposalState(proposalId, pstate.GetTaskSender(), orgState)
	if first {
		t.state.StoreProposalState(pstate)
	} else {
		t.state.UpdateProposalState(pstate)
	}
}

func (t *Twopc) removeOrgProposalState(proposalId common.Hash, partyId string) {
	pstate := t.state.GetProposalState(proposalId)

	if pstate.IsEmpty() {
		return
	}
	pstate.RemoveOrgProposalState(partyId)
	t.wal.DeleteState(t.wal.GetProposalSetKey(proposalId, partyId))
}

func (t *Twopc) getOrgProposalState(proposalId common.Hash, partyId string) (*ctypes.OrgProposalState, bool) {
	pstate := t.state.GetProposalState(proposalId)
	if pstate.IsEmpty() {
		return nil, false
	}
	return pstate.GetOrgProposalState(partyId)
}

func (t *Twopc) mustGetOrgProposalState(proposalId common.Hash, partyId string) *ctypes.OrgProposalState {
	pstate := t.state.GetProposalState(proposalId)
	if pstate.IsEmpty() {
		return nil
	}
	return pstate.MustGetOrgProposalState(partyId)
}

func (t *Twopc) makeEmptyConfirmTaskPeerDesc() *twopcpb.ConfirmTaskPeerInfo {
	return &twopcpb.ConfirmTaskPeerInfo{
		DataSupplierPeerInfos:   make([]*twopcpb.TaskPeerInfo, 0),
		PowerSupplierPeerInfos:  make([]*twopcpb.TaskPeerInfo, 0),
		ResultReceiverPeerInfos: make([]*twopcpb.TaskPeerInfo, 0),
	}
}

func (t *Twopc) makeConfirmTaskPeerDesc(proposalId common.Hash) *twopcpb.ConfirmTaskPeerInfo {

	dataSuppliers, powerSuppliers, receivers := make([]*twopcpb.TaskPeerInfo, 0), make([]*twopcpb.TaskPeerInfo, 0), make([]*twopcpb.TaskPeerInfo, 0)

	for _, vote := range t.state.GetPrepareVoteArr(proposalId) {

		if vote.MsgOption.SenderRole == apicommonpb.TaskRole_TaskRole_DataSupplier && nil != vote.PeerInfo {
			dataSuppliers = append(dataSuppliers, types.ConvertTaskPeerInfo(vote.PeerInfo))
		}
		if vote.MsgOption.SenderRole == apicommonpb.TaskRole_TaskRole_PowerSupplier && nil != vote.PeerInfo {
			powerSuppliers = append(powerSuppliers, types.ConvertTaskPeerInfo(vote.PeerInfo))
		}
		if vote.MsgOption.SenderRole == apicommonpb.TaskRole_TaskRole_Receiver && nil != vote.PeerInfo {
			receivers = append(receivers, types.ConvertTaskPeerInfo(vote.PeerInfo))
		}
	}
	return &twopcpb.ConfirmTaskPeerInfo{
		DataSupplierPeerInfos:   dataSuppliers,
		PowerSupplierPeerInfos:  powerSuppliers,
		ResultReceiverPeerInfos: receivers,
	}
}

func (t *Twopc) refreshProposalState() {

	identity, err := t.resourceMng.GetDB().QueryIdentity()
	if nil != err {
		//log.WithError(err).Errorf("Failed to query local identity on consensus.refreshProposalState()")
		return
	}

	t.state.proposalsLock.Lock()
	defer t.state.proposalsLock.Unlock()

	for proposalId, pstate := range t.state.proposalSet {

		//pstate.MustLock()
		for partyId, orgState := range pstate.GetStateCache() {

			if orgState.IsDeadline() {
				log.Debugf("Started refresh proposalState loop, the proposalState direct be deadline, remove org proposalState and task cache on Consensus, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
					pstate.GetProposalId().String(), pstate.GetTaskId(), partyId)

				_, ok := t.state.GetProposalTaskWithPartyId(pstate.GetTaskId(), partyId)

				t.state.RemoveProposalTaskWithPartyId(pstate.GetTaskId(), partyId) // remove proposal task with partyId
				pstate.RemoveOrgProposalStateUnSafe(partyId)                       // remove state with partyId
				t.state.RemoveOrgPrepareVoteState(proposalId, partyId)             // remove prepare vote with partyId
				t.state.RemoveOrgConfirmVoteState(proposalId, partyId)             // remove confirm vote with partyId

				if pstate.IsEmpty() {
					delete(t.state.proposalSet, proposalId)
					t.removeConfirmTaskPeerInfo(proposalId) // remove confirm peers and wal
				} else {
					t.state.proposalSet[proposalId] = pstate
				}
				t.wal.DeleteState(t.wal.GetProposalTaskCacheKey(pstate.GetTaskId(), partyId))
				t.wal.DeleteState(t.wal.GetPrepareVotesKey(proposalId, partyId))
				t.wal.DeleteState(t.wal.GetConfirmVotesKey(proposalId, partyId))
				t.wal.DeleteState(t.wal.GetProposalSetKey(proposalId, partyId))

				if !ok {
					log.Errorf("Not found the proposalTask on `twopc.refreshProposalState()`, skip this proposal state, taskId: {%s}, partyId: {%s}",
						pstate.GetTaskId(), partyId)
					continue
				}

				// release local resource and clean some data  (on task partenr)
				t.resourceMng.GetDB().StoreTaskEvent(&libtypes.TaskEvent{
					Type:       evengine.TaskProposalStateDeadline.Type,
					IdentityId: identity.GetIdentityId(),
					PartyId:    partyId,
					TaskId:     pstate.GetTaskId(),
					Content:    fmt.Sprintf("%s for myself", evengine.TaskProposalStateDeadline.Msg),
					CreateAt:   timeutils.UnixMsecUint64(),
				})

				t.stopTaskConsensus("the proposalState had been deadline",
					proposalId,
					pstate.GetTaskId(),
					orgState.GetTaskRole(),
					apicommonpb.TaskRole_TaskRole_Sender,
					&apicommonpb.TaskOrganization{
						PartyId:    partyId,
						NodeName:   identity.GetNodeName(),
						NodeId:     identity.GetNodeId(),
						IdentityId: identity.GetIdentityId(),
					},
					pstate.GetTaskSender(), types.TaskConsensusInterrupt)
				continue
			}

			switch orgState.CurrPeriodNum() {

			case ctypes.PeriodPrepare:

				if orgState.IsPrepareTimeout() {
					log.Debugf("Started refresh org proposalState, the org proposalState was prepareTimeout, change to confirm epoch, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
						pstate.GetProposalId().String(), pstate.GetTaskId(), partyId)

					orgState.ChangeToConfirm(orgState.PeriodStartTime + uint64(ctypes.PrepareMsgVotingDuration.Milliseconds()))
					pstate.StoreOrgProposalStateUnSafe(orgState)

					t.wal.StoreOrgProposalState(pstate.GetProposalId(), pstate.GetTaskSender(), orgState)
				}

			case ctypes.PeriodConfirm:

				if orgState.IsConfirmTimeout() {
					log.Debugf("Started refresh org proposalState, the org proposalState was confirmTimeout, change to commit epoch, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
						pstate.GetProposalId().String(), pstate.GetTaskId(), partyId)

					orgState.ChangeToCommit(orgState.PeriodStartTime + uint64(ctypes.ConfirmMsgVotingDuration.Milliseconds()))
					pstate.StoreOrgProposalStateUnSafe(orgState)

					t.wal.StoreOrgProposalState(pstate.GetProposalId(), pstate.GetTaskSender(), orgState)
				}

			case ctypes.PeriodCommit:

				if orgState.IsCommitTimeout() {
					log.Debugf("Started refresh org proposalState, the org proposalState was commitTimeout, change to finished epoch, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
						pstate.GetProposalId().String(), pstate.GetTaskId(), partyId)

					orgState.ChangeToFinished(orgState.PeriodStartTime + uint64(ctypes.CommitMsgEndingDuration.Milliseconds()))
					pstate.StoreOrgProposalStateUnSafe(orgState)

					t.wal.StoreOrgProposalState(pstate.GetProposalId(), pstate.GetTaskSender(), orgState)
				}

			case ctypes.PeriodFinished:

				if orgState.IsDeadline() {
					log.Debugf("Started refresh org proposalState, the org proposalState was finished, but coming deadline now, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
						pstate.GetProposalId().String(), pstate.GetTaskId(), partyId)

					_, ok := t.state.GetProposalTaskWithPartyId(pstate.GetTaskId(), partyId)

					t.state.RemoveProposalTaskWithPartyId(pstate.GetTaskId(), partyId) // remove proposal task with partyId
					pstate.RemoveOrgProposalStateUnSafe(partyId)                       // remove state with partyId
					t.state.RemoveOrgPrepareVoteState(proposalId, partyId)             // remove prepare vote with partyId
					t.state.RemoveOrgConfirmVoteState(proposalId, partyId)             // remove confirm vote with partyId

					if pstate.IsEmpty() {
						delete(t.state.proposalSet, proposalId)
						t.removeConfirmTaskPeerInfo(proposalId) // remove confirm peers and wal
					} else {
						t.state.proposalSet[proposalId] = pstate
					}
					t.wal.DeleteState(t.wal.GetProposalTaskCacheKey(pstate.GetTaskId(), partyId))
					t.wal.DeleteState(t.wal.GetPrepareVotesKey(proposalId, partyId))
					t.wal.DeleteState(t.wal.GetConfirmVotesKey(proposalId, partyId))
					t.wal.DeleteState(t.wal.GetProposalSetKey(proposalId, partyId))

					if !ok {
						log.Errorf("Not found the proposalTask on `twopc.refreshProposalState()`, skip this proposal state, taskId: {%s}, partyId: {%s}",
							pstate.GetTaskId(), partyId)
						continue
					}

					// release local resource and clean some data  (on task partenr)
					t.resourceMng.GetDB().StoreTaskEvent(&libtypes.TaskEvent{
						Type:       evengine.TaskProposalStateDeadline.Type,
						IdentityId: identity.GetIdentityId(),
						PartyId:    partyId,
						TaskId:     pstate.GetTaskId(),
						Content:    fmt.Sprintf("%s for myself", evengine.TaskProposalStateDeadline.Msg),
						CreateAt:   timeutils.UnixMsecUint64(),
					})

					t.stopTaskConsensus("the proposalState had been deadline",
						proposalId,
						pstate.GetTaskId(),
						orgState.GetTaskRole(),
						apicommonpb.TaskRole_TaskRole_Sender,
						&apicommonpb.TaskOrganization{
							PartyId:    partyId,
							NodeName:   identity.GetNodeName(),
							NodeId:     identity.GetNodeId(),
							IdentityId: identity.GetIdentityId(),
						},
						pstate.GetTaskSender(), types.TaskConsensusInterrupt)
				}

			default:
				log.Errorf("Unknown the proposalState period,  proposalId: {%s}, taskId: {%s}, partyId: {%s}", pstate.GetProposalId().String(), pstate.GetTaskId(), partyId)
			}
		}

		//pstate.MustUnLock()
	}

}

func (t *Twopc) stopTaskConsensus(
	reason string,
	proposalId common.Hash,
	taskId string,
	senderRole, receiverRole apicommonpb.TaskRole,
	sender, receiver *apicommonpb.TaskOrganization,
	taskActionStatus types.TaskActionStatus,
) {

	log.Debugf("Call stopTaskConsensus() to interrupt consensus msg %s with role is %s, proposalId: {%s}, taskId: {%s}, partyId: {%s}, remote partyId: {%s}, taskActionStatus: {%s}",
		reason, senderRole, proposalId.String(), taskId, sender.GetPartyId(), receiver.GetPartyId(), taskActionStatus.String())
	// Send consensus result to interrupt consensus epoch and clean some data (on task sender)
	if senderRole == apicommonpb.TaskRole_TaskRole_Sender {
		t.replyTaskConsensusResult(types.NewTaskConsResult(taskId, taskActionStatus, fmt.Errorf(reason)))
	} else {

		var remotePID peer.ID
		// remote task, so need convert remote pid by nodeId of task sender identity
		if receiver.GetIdentityId() != sender.GetIdentityId() {
			receiverPID, err := p2p.HexPeerID(receiver.GetNodeId())
			if nil != err {
				log.WithError(err).Errorf("Failed to convert nodeId to remote pid receiver identity when call stopTaskConsensus() to interrupt consensus msg %s with role is %s, proposalId: {%s}, taskId: {%s}, partyId: {%s}, remote partyId: {%s}",
					reason, senderRole, proposalId.String(), taskId, sender.GetPartyId(), receiver.GetPartyId())
			}
			remotePID = receiverPID
		}

		t.sendNeedExecuteTask(types.NewNeedExecuteTask(
			remotePID,
			senderRole,
			receiverRole,
			sender,
			receiver,
			taskId,
			taskActionStatus,
			&types.PrepareVoteResource{},   // zero value
			&twopcpb.ConfirmTaskPeerInfo{}, // zero value
			fmt.Errorf(reason),
		))
		// Finally, release local task cache that task manager will do it. (to call `resourceMng.ReleaseLocalResourceWithTask()` by taskManager)
	}
}

func (t *Twopc) driveTask(
	remotePid peer.ID,
	proposalId common.Hash,
	localTaskRole apicommonpb.TaskRole,
	localTaskOrganization *apicommonpb.TaskOrganization,
	remoteTaskRole apicommonpb.TaskRole,
	remoteTaskOrganization *apicommonpb.TaskOrganization,
	taskId string,
) {

	log.Debugf("Start to call `2pc.driveTask()`, proposalId: {%s}, taskId: {%s}, partyId: {%s}, identityId: {%s}, nodeName: {%s}",
		proposalId.String(), taskId, localTaskOrganization.GetPartyId(), localTaskOrganization.GetIdentityId(), localTaskOrganization.GetNodeName())

	selfvote := t.getPrepareVote(proposalId, localTaskOrganization.GetPartyId())
	if nil == selfvote {
		log.Errorf("Failed to find local cache about prepareVote myself internal resource on `2pc.driveTask()`, proposalId: {%s}, taskId: {%s}, partyId: {%s}, identityId: {%s}, nodeName: {%s}",
			proposalId.String(), taskId, localTaskOrganization.GetPartyId(), localTaskOrganization.GetIdentityId(), localTaskOrganization.GetNodeName())
		return
	}

	peers, ok := t.getConfirmTaskPeerInfo(proposalId)
	if !ok {
		log.Errorf("Failed to find local cache about prepareVote all peer resource {externalIP:externalPORT} on `2pc.driveTask()`, proposalId: {%s}, taskId: {%s}, partyId: {%s}, identityId: {%s}, nodeName: {%s}",
			proposalId.String(), taskId, localTaskOrganization.GetPartyId(), localTaskOrganization.GetIdentityId(), localTaskOrganization.GetNodeName())

		return
	}

	log.Debugf("Find vote resources on `2pc.driveTask()` proposalId: {%s}, taskId: {%s}, partyId: {%s}, identityId: {%s}, nodeName: {%s}, self vote: %s, peers: %s",
		proposalId.String(), taskId, localTaskOrganization.GetPartyId(), localTaskOrganization.GetIdentityId(), localTaskOrganization.GetNodeName(),
		selfvote.String(), peers.String())

	// Send task to TaskManager to execute
	t.sendNeedExecuteTask(types.NewNeedExecuteTask(
		remotePid,
		localTaskRole,
		remoteTaskRole,
		localTaskOrganization,
		remoteTaskOrganization,
		taskId,
		types.TaskNeedExecute,
		selfvote.GetPeerInfo(),
		peers,
		nil,
	))
}

func (t *Twopc) sendPrepareMsg(proposalId common.Hash, nonConsTask *types.NeedConsensusTask, startTime uint64) error {


	task := nonConsTask.GetTask()
	sender := task.GetTaskSender()

	sendTaskFn := func(wg *sync.WaitGroup, sender, receiver *apicommonpb.TaskOrganization, senderRole, receiverRole apicommonpb.TaskRole, errCh chan<- error) {

		defer wg.Done()

		var pid, err = p2p.HexPeerID(receiver.NodeId)
		if nil != err {
			errCh <- fmt.Errorf("failed to nodeId => peerId, proposalId: %s, taskId: %s, other peer taskRole: %s, other peer taskPartyId: %s, identityId: %s, pid: %s, err: %s",
				proposalId.String(), task.GetTaskId(), receiverRole.String(), receiver.GetPartyId(), receiver.GetIdentityId(), pid, err)
			return
		}

		prepareMsg, err := makePrepareMsg(proposalId, senderRole, receiverRole, sender.GetPartyId(), receiver.GetPartyId(), nonConsTask, startTime)

		if nil != err {
			errCh <- fmt.Errorf("failed to make prepareMsg, proposalId: %s, taskId: %s, other peer taskRole: %s, other peer taskPartyId: %s, identityId: %s, pid: %s, err: %s",
				proposalId.String(), task.GetTaskId(), receiverRole.String(), receiver.GetPartyId(), receiver.GetIdentityId(), pid, err)
			return
		}

		var sendErr error
		var logdesc string
		if types.IsSameTaskOrg(sender, receiver) {
			sendErr = t.sendLocalPrepareMsg(pid, prepareMsg)
			logdesc = "sendLocalPrepareMsg()"
		} else {
			//sendErr = handler.SendTwoPcPrepareMsg(context.TODO(), t.p2p, pid, prepareMsg)
			sendErr = t.p2p.Broadcast(context.TODO(), prepareMsg)
			logdesc = "Broadcast()"
		}

		if nil != sendErr {
			errCh <- fmt.Errorf("failed to call `sendPrepareMsg.%s` proposalId: %s, taskId: %s, other peer taskRole: %s, other peer taskPartyId: %s, identityId: %s, pid: %s, err: %s",
				logdesc, proposalId.String(), task.GetTaskId(), receiverRole.String(), receiver.GetPartyId(), receiver.GetIdentityId(), pid, sendErr)
			return
		}

		log.Debugf("Succceed to call `sendPrepareMsg.%s` proposalId: %s, taskId: %s, other peer taskRole: %s, other peer taskPartyId: %s, identityId: %s, pid: %s",
			logdesc, proposalId.String(), task.GetTaskId(), receiverRole.String(), receiver.GetPartyId(), receiver.GetIdentityId(), pid)
	}

	size := (len(task.GetTaskData().GetDataSuppliers())) + len(task.GetTaskData().GetPowerSuppliers()) + len(task.GetTaskData().GetReceivers())
	errCh := make(chan error, size)
	var wg sync.WaitGroup

	for i := 0; i < len(task.GetTaskData().GetDataSuppliers()); i++ {

		wg.Add(1)
		dataSupplier := task.GetTaskData().GetDataSuppliers()[i]
		receiver := dataSupplier.GetOrganization()
		go sendTaskFn(&wg, sender, receiver, apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_DataSupplier, errCh)

	}
	for i := 0; i < len(task.GetTaskData().GetPowerSuppliers()); i++ {

		wg.Add(1)
		powerSupplier := task.GetTaskData().GetPowerSuppliers()[i]
		receiver := powerSupplier.GetOrganization()
		go sendTaskFn(&wg, sender, receiver, apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_PowerSupplier, errCh)

	}

	for i := 0; i < len(task.GetTaskData().GetReceivers()); i++ {

		wg.Add(1)
		receiver := task.GetTaskData().GetReceivers()[i]
		go sendTaskFn(&wg, sender, receiver, apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_Receiver, errCh)
	}

	wg.Wait()
	close(errCh)

	errStrs := make([]string, 0)

	for err := range errCh {
		if nil != err {
			errStrs = append(errStrs, err.Error())
		}
	}
	if len(errStrs) != 0 {
		return fmt.Errorf(
			"\n######################################################## \n%s\n########################################################\n",
			strings.Join(errStrs, "\n"))
	}
	return nil
}

func (t *Twopc) addTaskResultCh(taskId string, resultCh chan<- *types.TaskConsResult) {
	t.taskResultLock.Lock()
	log.Debugf("AddTaskResultCh taskId: {%s}", taskId)
	t.taskResultChSet[taskId] = resultCh
	t.taskResultLock.Unlock()
}
func (t *Twopc) removeTaskResultCh(taskId string) {
	t.taskResultLock.Lock()
	log.Debugf("RemoveTaskResultCh taskId: {%s}", taskId)
	delete(t.taskResultChSet, taskId)
	t.taskResultLock.Unlock()
}
func (t *Twopc) replyTaskConsensusResult(result *types.TaskConsResult) {
	t.taskResultBusCh <- result
}
func (t *Twopc) handleTaskConsensusResult(result *types.TaskConsResult) {
	t.taskResultLock.Lock()
	log.Debugf("Need SendTaskResultCh taskId: {%s}, result: {%s}", result.GetTaskId(), result.String())
	if ch, ok := t.taskResultChSet[result.GetTaskId()]; ok {
		log.Debugf("Start SendTaskResultCh taskId: {%s}, result: {%s}", result.GetTaskId(), result.String())
		ch <- result
		close(ch)
		delete(t.taskResultChSet, result.GetTaskId())
	}
	t.taskResultLock.Unlock()
}

func (t *Twopc) sendNeedReplayScheduleTask(task *types.NeedReplayScheduleTask) {
	t.needReplayScheduleTaskCh <- task
}

func (t *Twopc) sendNeedExecuteTask(task *types.NeedExecuteTask) {
	t.needExecuteTaskCh <- task
}

func (t *Twopc) sendPrepareVote(pid peer.ID, sender, receiver *apicommonpb.TaskOrganization, req *twopcpb.PrepareVote) error {
	if types.IsNotSameTaskOrg(sender, receiver) {
		//return handler.SendTwoPcPrepareVote(context.TODO(), t.p2p, pid, req)
		return t.p2p.Broadcast(context.TODO(), req)
	} else {
		return t.sendLocalPrepareVote(pid, req)
	}
}

func (t *Twopc) sendConfirmMsg(proposalId common.Hash, task *types.Task, peers *twopcpb.ConfirmTaskPeerInfo, option types.TwopcMsgOption, startTime uint64) error {

	sender := task.GetTaskSender()

	sendConfirmMsgFn := func(wg *sync.WaitGroup, sender, receiver *apicommonpb.TaskOrganization, senderRole, receiverRole apicommonpb.TaskRole, errCh chan<- error) {

		defer wg.Done()

		pid, err := p2p.HexPeerID(receiver.NodeId)
		if nil != err {
			errCh <- fmt.Errorf("failed to nodeId => peerId, proposalId: %s, taskId: %s, other peer's taskRole: %s, other peer's partyId: %s, other identityId: %s, pid: %s, err: %s",
				proposalId.String(), task.GetTaskId(), receiverRole.String(), receiver.GetPartyId(), receiver.GetIdentityId(), pid, err)
			return
		}

		confirmMsg := makeConfirmMsg(proposalId, senderRole, receiverRole, sender.GetPartyId(), receiver.GetPartyId(), sender, peers, option, startTime)

		var sendErr error
		var logdesc string
		if types.IsSameTaskOrg(sender, receiver) {
			sendErr = t.sendLocalConfirmMsg(pid, confirmMsg)
			logdesc = "sendLocalConfirmMsg()"
		} else {
			//sendErr = handler.SendTwoPcConfirmMsg(context.TODO(), t.p2p, pid, confirmMsg)
			sendErr = t.p2p.Broadcast(context.TODO(), confirmMsg)
			logdesc = "Broadcast()"
		}

		// Send the ConfirmMsg to other peer
		if nil != sendErr {
			errCh <- fmt.Errorf("failed to call`sendConfirmMsg.%s` proposalId: %s, taskId: %s,other peer's taskRole: %s, other peer's partyId: %s, other identityId: %s, pid: %s, err: %s",
				logdesc, proposalId.String(), task.GetTaskId(), receiverRole.String(), receiver.GetPartyId(), receiver.GetIdentityId(), pid, sendErr)
			errCh <- err
			return
		}

		log.Debugf("Succceed to call`sendConfirmMsg.%s` proposalId: %s, taskId: %s,other peer's taskRole: %s, other peer's partyId: %s, other identityId: %s, pid: %s",
			logdesc, proposalId.String(), task.GetTaskId(), receiverRole.String(), receiver.GetPartyId(), receiver.GetIdentityId(), pid)

	}

	size := (len(task.GetTaskData().GetDataSuppliers())) + len(task.GetTaskData().GetPowerSuppliers()) + len(task.GetTaskData().GetReceivers())
	errCh := make(chan error, size)
	var wg sync.WaitGroup

	for i := 0; i < len(task.GetTaskData().GetDataSuppliers()); i++ {

		wg.Add(1)
		dataSupplier := task.GetTaskData().GetDataSuppliers()[i]
		receiver := dataSupplier.GetOrganization()
		go sendConfirmMsgFn(&wg, sender, receiver, apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_DataSupplier, errCh)

	}
	for i := 0; i < len(task.GetTaskData().GetPowerSuppliers()); i++ {

		wg.Add(1)
		powerSupplier := task.GetTaskData().GetPowerSuppliers()[i]
		receiver := powerSupplier.GetOrganization()
		go sendConfirmMsgFn(&wg, sender, receiver, apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_PowerSupplier, errCh)

	}

	for i := 0; i < len(task.GetTaskData().GetReceivers()); i++ {

		wg.Add(1)
		receiver := task.GetTaskData().GetReceivers()[i]
		go sendConfirmMsgFn(&wg, sender, receiver, apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_Receiver, errCh)
	}

	wg.Wait()
	close(errCh)

	errStrs := make([]string, 0)

	for err := range errCh {
		if nil != err {
			errStrs = append(errStrs, err.Error())
		}
	}
	if len(errStrs) != 0 {
		return fmt.Errorf(
			"\n######################################################## \n%s\n########################################################\n",
			strings.Join(errStrs, "\n"))
	}
	return nil
}

func (t *Twopc) sendConfirmVote(pid peer.ID, sender, receiver *apicommonpb.TaskOrganization, req *twopcpb.ConfirmVote) error {
	if types.IsNotSameTaskOrg(sender, receiver) {
		//return handler.SendTwoPcConfirmVote(context.TODO(), t.p2p, pid, req)
		return t.p2p.Broadcast(context.TODO(), req)
	} else {
		return t.sendLocalConfirmVote(pid, req)
	}
}

func (t *Twopc) sendCommitMsg(proposalId common.Hash, task *types.Task, option types.TwopcMsgOption, startTime uint64) error {

	sender := task.GetTaskSender()

	sendCommitMsgFn := func(wg *sync.WaitGroup, sender, receiver *apicommonpb.TaskOrganization, senderRole, receiverRole apicommonpb.TaskRole, errCh chan<- error) {

		defer wg.Done()

		pid, err := p2p.HexPeerID(receiver.NodeId)
		if nil != err {
			errCh <- fmt.Errorf("failed to nodeId => peerId, proposalId: %s, taskId: %s, other peer's taskRole: %s, other peer's partyId: %s, identityId: %s, pid: %s, err: %s",
				proposalId.String(), task.GetTaskId(), receiverRole.String(), receiver.GetPartyId(), receiver.GetIdentityId(), pid, err)
			return
		}

		commitMsg := makeCommitMsg(proposalId, senderRole, receiverRole, sender.GetPartyId(), receiver.GetPartyId(), sender, option, startTime)

		var sendErr error
		var logdesc string
		if types.IsSameTaskOrg(sender, receiver) {
			sendErr = t.sendLocalCommitMsg(pid, commitMsg)
			logdesc = "sendLocalCommitMsg()"
		} else {
			//sendErr = handler.SendTwoPcCommitMsg(context.TODO(), t.p2p, pid, commitMsg)
			sendErr = t.p2p.Broadcast(context.TODO(), commitMsg)
			logdesc = "Broadcast()"
		}

		// Send the ConfirmMsg to other peer
		if nil != sendErr {
			errCh <- fmt.Errorf("failed to call`sendCommitMsg.%s` proposalId: %s, taskId: %s,  other peer's taskRole: %s, other peer's partyId: %s, identityId: %s, pid: %s, err: %s",
				logdesc, proposalId.String(), task.GetTaskId(), receiverRole.String(), receiver.GetPartyId(), receiver.GetIdentityId(), pid, sendErr)
			errCh <- err
			return
		}

		log.Debugf("Succceed to call`sendCommitMsg.%s` proposalId: %s, taskId: %s,  other peer's taskRole: %s, other peer's partyId: %s, identityId: %s, pid: %s",
			logdesc, proposalId.String(), task.GetTaskId(), receiverRole.String(), receiver.GetPartyId(), receiver.GetIdentityId(), pid)

	}

	size := (len(task.GetTaskData().GetDataSuppliers())) + len(task.GetTaskData().GetPowerSuppliers()) + len(task.GetTaskData().GetReceivers())
	errCh := make(chan error, size)
	var wg sync.WaitGroup

	for i := 0; i < len(task.GetTaskData().GetDataSuppliers()); i++ {

		wg.Add(1)
		dataSupplier := task.GetTaskData().GetDataSuppliers()[i]
		receiver := dataSupplier.GetOrganization()
		go sendCommitMsgFn(&wg, sender, receiver, apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_DataSupplier, errCh)

	}
	for i := 0; i < len(task.GetTaskData().GetPowerSuppliers()); i++ {

		wg.Add(1)
		powerSupplier := task.GetTaskData().GetPowerSuppliers()[i]
		receiver := powerSupplier.GetOrganization()
		go sendCommitMsgFn(&wg, sender, receiver, apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_PowerSupplier, errCh)

	}

	for i := 0; i < len(task.GetTaskData().GetReceivers()); i++ {

		wg.Add(1)
		receiver := task.GetTaskData().GetReceivers()[i]
		go sendCommitMsgFn(&wg, sender, receiver, apicommonpb.TaskRole_TaskRole_Sender, apicommonpb.TaskRole_TaskRole_Receiver, errCh)
	}

	wg.Wait()
	close(errCh)

	errStrs := make([]string, 0)

	for err := range errCh {
		if nil != err {
			errStrs = append(errStrs, err.Error())
		}
	}
	if len(errStrs) != 0 {
		return fmt.Errorf(
			"\n######################################################## \n%s\n########################################################\n",
			strings.Join(errStrs, "\n"))
	}
	return nil
}

func verifyPartyRole(partyId string, role apicommonpb.TaskRole, task *types.Task) bool {

	switch role {
	case apicommonpb.TaskRole_TaskRole_Sender:
		if partyId == task.GetTaskSender().GetPartyId() {
			return true
		}
	case apicommonpb.TaskRole_TaskRole_DataSupplier:
		for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
			if partyId == dataSupplier.GetOrganization().GetPartyId() {
				return true
			}
		}
	case apicommonpb.TaskRole_TaskRole_PowerSupplier:
		for _, powerSupplier := range task.GetTaskData().GetPowerSuppliers() {
			if partyId == powerSupplier.GetOrganization().GetPartyId() {
				return true
			}
		}
	case apicommonpb.TaskRole_TaskRole_Receiver:
		for _, receiver := range task.GetTaskData().GetReceivers() {
			if partyId == receiver.GetPartyId() {
				return true
			}
		}
	}
	return false
}

func fetchOrgByPartyRole(partyId string, role apicommonpb.TaskRole, task *types.Task) *apicommonpb.TaskOrganization {

	switch role {
	case apicommonpb.TaskRole_TaskRole_Sender:
		if partyId == task.GetTaskSender().GetPartyId() {
			return task.GetTaskSender()
		}
	case apicommonpb.TaskRole_TaskRole_DataSupplier:
		for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
			if partyId == dataSupplier.GetOrganization().GetPartyId() {
				return dataSupplier.GetOrganization()
			}
		}
	case apicommonpb.TaskRole_TaskRole_PowerSupplier:
		for _, powerSupplier := range task.GetTaskData().GetPowerSuppliers() {
			if partyId == powerSupplier.GetOrganization().GetPartyId() {
				return powerSupplier.GetOrganization()
			}
		}
	case apicommonpb.TaskRole_TaskRole_Receiver:
		for _, receiver := range task.GetTaskData().GetReceivers() {
			if partyId == receiver.GetPartyId() {
				return receiver
			}
		}
	}
	return nil
}

func (t *Twopc) verifyPrepareVoteRoleIsTaskPartner(identityId, partyId string, role apicommonpb.TaskRole, task *types.Task) (bool, error) {
	var find bool
	switch role {
	case apicommonpb.TaskRole_TaskRole_DataSupplier:
		for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
			// identity + partyId
			if dataSupplier.GetOrganization().GetIdentityId() == identityId && dataSupplier.GetOrganization().GetPartyId() == partyId {
				find = true
				break
			}
		}
	case apicommonpb.TaskRole_TaskRole_PowerSupplier:
		for _, powerSupplier := range task.GetTaskData().GetPowerSuppliers() {
			// identity + partyId
			if powerSupplier.GetOrganization().GetIdentityId() == identityId && powerSupplier.GetOrganization().GetPartyId() == partyId {
				find = true
				break
			}
		}
	case apicommonpb.TaskRole_TaskRole_Receiver:
		for _, receiver := range task.GetTaskData().GetReceivers() {
			// identity + partyId
			if receiver.GetIdentityId() == identityId && receiver.GetPartyId() == partyId {
				find = true
				break
			}
		}
	default:
		return false, fmt.Errorf("Invalid role, the role is not task partners role on the prepare vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
			task.GetTaskData().GetTaskId(), role.String(), identityId, partyId)

	}
	return find, nil
}
func (t *Twopc) verifyConfirmVoteRoleIsTaskPartner(identityId, partyId string, role apicommonpb.TaskRole, task *types.Task) (bool, error) {
	var identityValid bool
	switch role {
	case apicommonpb.TaskRole_TaskRole_DataSupplier:

		for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
			// identity + partyId
			if dataSupplier.GetOrganization().GetIdentityId() == identityId && dataSupplier.GetOrganization().GetPartyId() == partyId {
				identityValid = true
				break
			}
		}
	case apicommonpb.TaskRole_TaskRole_PowerSupplier:

		for _, powerSupplier := range task.GetTaskData().GetPowerSuppliers() {
			// identity + partyId
			if powerSupplier.GetOrganization().GetIdentityId() == identityId && powerSupplier.GetOrganization().GetPartyId() == partyId {
				identityValid = true
				break
			}
		}
	case apicommonpb.TaskRole_TaskRole_Receiver:

		for _, receiver := range task.GetTaskData().GetReceivers() {
			// identity + partyId
			if receiver.GetIdentityId() == identityId && receiver.GetPartyId() == partyId {
				identityValid = true
				break
			}
		}
	default:
		return false, fmt.Errorf("Invalid role, the role is not task partners role on the confirm vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
			task.GetTaskData().GetTaskId(), role.String(), identityId, partyId)

	}
	return identityValid, nil
}

func (t *Twopc) getPrepareVote(proposalId common.Hash, partyId string) *types.PrepareVote {
	return t.state.GetPrepareVote(proposalId, partyId)
}

func (t *Twopc) getConfirmVote(proposalId common.Hash, partyId string) *types.ConfirmVote {
	return t.state.GetConfirmVote(proposalId, partyId)
}

func (t *Twopc) storeConfirmTaskPeerInfo(proposalId common.Hash, peers *twopcpb.ConfirmTaskPeerInfo) {
	t.state.StoreConfirmTaskPeerInfo(proposalId, peers)
}

func (t *Twopc) hasConfirmTaskPeerInfo(proposalId common.Hash) bool {
	return t.state.HasConfirmTaskPeerInfo(proposalId)
}

func (t *Twopc) getConfirmTaskPeerInfo(proposalId common.Hash) (*twopcpb.ConfirmTaskPeerInfo, bool) {
	return t.state.GetConfirmTaskPeerInfo(proposalId)
}

func (t *Twopc) mustGetConfirmTaskPeerInfo(proposalId common.Hash) *twopcpb.ConfirmTaskPeerInfo {
	return t.state.MustGetConfirmTaskPeerInfo(proposalId)
}

func (t *Twopc) removeConfirmTaskPeerInfo(proposalId common.Hash) {
	t.state.RemoveConfirmTaskPeerInfo(proposalId)
}

func (t *Twopc) recoverCache() {

	errCh := make(chan error, 5)
	var wg sync.WaitGroup
	wg.Add(5)
	// recovery proposalTaskCache  (taskId -> partyId -> proposalTask)
	go func(wg *sync.WaitGroup, errCh chan<- error) {

		defer wg.Done()

		prefixLength := len(proposalTaskCachePrefix)
		if err := t.wal.ForEachKVWithPrefix(proposalTaskCachePrefix, func(key, value []byte) error {

			if len(key) != 0 && len(value) != 0 {
				taskId, partyId := string(key[prefixLength:prefixLength+71]), string(key[prefixLength+71:])
				proposalTaskPB := &libtypes.ProposalTask{}
				if err := proto.Unmarshal(value, proposalTaskPB); err != nil {
					t.wal.DeleteState(key)
					return fmt.Errorf("unmarshal proposalTask failed, %s", err)
				}

				//task, err := t.resourceMng.GetDB().QueryLocalTask(proposalTaskPB.GetTaskId())
				//if nil != err {
				//	return fmt.Errorf("query local task failed on recover proposalTask from wal, %s, taskId: {%s}", err, proposalTaskPB.GetTaskId())
				//}

				cache, ok := t.state.proposalTaskCache[taskId]
				if !ok {
					cache = make(map[string]*types.ProposalTask, 0)
				}
				cache[partyId] = types.NewProposalTask(common.HexToHash(proposalTaskPB.GetProposalId()), taskId, proposalTaskPB.GetCreateAt())
				t.state.proposalTaskCache[taskId] = cache
			}
			return nil
		}); nil != err {
			errCh <- err
			return
		}
	}(&wg, errCh)

	// recovery proposalSet (proposalId -> partyId -> orgState)
	go func(wg *sync.WaitGroup, errCh chan<- error) {

		defer wg.Done()

		prefixLength := len(proposalSetPrefix)
		if err := t.wal.ForEachKVWithPrefix(proposalSetPrefix, func(key, value []byte) error {

			if len(key) != 0 && len(value) != 0 {

				// prefix + len(common.Hash) == prefix + len(proposalId)
				proposalId := common.BytesToHash(key[prefixLength : prefixLength+32])

				libOrgProposalState := &libtypes.OrgProposalState{}
				if err := proto.Unmarshal(value, libOrgProposalState); err != nil {
					return fmt.Errorf("unmarshal org proposalState failed, %s", err)
				}
				proposalState, ok := t.state.proposalSet[proposalId]
				if !ok {
					proposalState = ctypes.NewProposalState(proposalId, libOrgProposalState.GetTaskId(), libOrgProposalState.GetTaskSender())
				}
				proposalState.StoreOrgProposalStateUnSafe(&ctypes.OrgProposalState{
					PrePeriodStartTime: libOrgProposalState.GetPrePeriodStartTime(),
					PeriodStartTime:    libOrgProposalState.GetPeriodStartTime(),
					DeadlineDuration:   libOrgProposalState.GetDeadlineDuration(),
					CreateAt:           libOrgProposalState.GetCreateAt(),
					TaskId:             libOrgProposalState.GetTaskId(),
					TaskRole:           libOrgProposalState.GetTaskRole(),
					TaskOrg:            libOrgProposalState.GetTaskOrg(),
					PeriodNum:          ctypes.ProposalStatePeriod(libOrgProposalState.GetPeriodNum()),
				})
				t.state.proposalSet[proposalId] = proposalState
			}
			return nil
		}); nil != err {
			errCh <- err
			return
		}
	}(&wg, errCh)

	// recovery prepareVotes (proposalId -> partyId -> prepareVote)
	go func(wg *sync.WaitGroup, errCh chan<- error) {

		defer wg.Done()

		prefixLength := len(prepareVotesPrefix)
		if err := t.wal.ForEachKVWithPrefix(prepareVotesPrefix, func(key, value []byte) error {

			if len(key) != 0 && len(value) != 0 {
				proposalId, partyId := common.BytesToHash(key[prefixLength:prefixLength+33]), string(key[prefixLength+34:])

				vote := &libtypes.PrepareVote{}
				if err := proto.Unmarshal(value, vote); err != nil {
					return fmt.Errorf("unmarshal org prepareVote failed, %s", err)
				}

				prepareVoteState, ok := t.state.prepareVotes[proposalId]
				if !ok {
					prepareVoteState = newPrepareVoteState()
				}
				prepareVoteState.addVote(&types.PrepareVote{
					MsgOption: &types.MsgOption{
						ProposalId:      proposalId,
						SenderRole:      vote.MsgOption.SenderRole,
						SenderPartyId:   partyId,
						ReceiverRole:    vote.MsgOption.ReceiverRole,
						ReceiverPartyId: vote.MsgOption.ReceiverPartyId,
						Owner:           vote.MsgOption.Owner,
					},
					VoteOption: types.VoteOption(vote.VoteOption),
					PeerInfo: &types.PrepareVoteResource{
						Id:      vote.PeerInfo.Id,
						Ip:      vote.PeerInfo.Ip,
						Port:    vote.PeerInfo.Port,
						PartyId: vote.PeerInfo.PartyId,
					},
					CreateAt: vote.CreateAt,
					Sign:     vote.Sign,
				})
				t.state.prepareVotes[proposalId] = prepareVoteState
			}
			return nil
		}); nil != err {
			errCh <- err
			return
		}
	}(&wg, errCh)

	// recovery confirmVotes (proposalId -> partyId -> confirmVote)
	go func(wg *sync.WaitGroup, errCh chan<- error) {

		defer wg.Done()

		prefixLength := len(confirmVotesPrefix)
		if err := t.wal.ForEachKVWithPrefix(confirmVotesPrefix, func(key, value []byte) error {

			if len(key) != 0 && len(value) != 0 {
				proposalId, partyId := common.BytesToHash(key[prefixLength:prefixLength+32]), string(key[prefixLength+32:])

				vote := &libtypes.ConfirmVote{}
				if err := proto.Unmarshal(value, vote); err != nil {
					return fmt.Errorf("unmarshal org confirmVote failed, %s", err)
				}
				confirmVoteState, ok := t.state.confirmVotes[proposalId]
				if !ok {
					confirmVoteState = newConfirmVoteState()
				}
				confirmVoteState.addVote(&types.ConfirmVote{
					MsgOption: &types.MsgOption{
						ProposalId:      proposalId,
						SenderRole:      vote.MsgOption.SenderRole,
						SenderPartyId:   partyId,
						ReceiverRole:    vote.MsgOption.ReceiverRole,
						ReceiverPartyId: vote.MsgOption.ReceiverPartyId,
						Owner:           vote.MsgOption.Owner,
					},
					VoteOption: types.VoteOption(vote.VoteOption),
					CreateAt:   vote.CreateAt,
					Sign:       vote.Sign,
				})
				t.state.confirmVotes[proposalId] = confirmVoteState
			}
			return nil
		}); nil != err {
			errCh <- err
			return
		}
	}(&wg, errCh)

	// recovery proposalPeerInfoCache (proposalId -> ConfirmTaskPeerInfo)
	go func(wg *sync.WaitGroup, errCh chan<- error) {

		defer wg.Done()

		prefixLength := len(proposalPeerInfoCachePrefix)
		if err := t.wal.ForEachKVWithPrefix(proposalPeerInfoCachePrefix, func(key, value []byte) error {

			if len(key) != 0 && len(value) != 0 {
				proposalId := common.BytesToHash(key[prefixLength:])
				confirmTaskPeerInfo := &twopcpb.ConfirmTaskPeerInfo{}
				if err := proto.Unmarshal(value, confirmTaskPeerInfo); err != nil {
					return fmt.Errorf("unmarshal confirmTaskPeerInfo failed, %s", err)
				}
				t.state.proposalPeerInfoCache[proposalId] = confirmTaskPeerInfo
			}
			return nil
		}); nil != err {
			errCh <- err
			return
		}
	}(&wg, errCh)

	wg.Wait()
	close(errCh)

	errStrs := make([]string, 0)

	for err := range errCh {
		if nil != err {
			errStrs = append(errStrs, err.Error())
		}
	}
	if len(errStrs) != 0 {
		log.Fatalf(
			"recover consensus state failed: \n%s",
			strings.Join(errStrs, "\n"))
	}
}
