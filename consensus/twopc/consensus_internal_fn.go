package twopc

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	"github.com/RosettaFlow/Carrier-Go/core/evengine"
	"github.com/RosettaFlow/Carrier-Go/handler"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	twopcpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/consensus/twopc"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/libp2p/go-libp2p-core/peer"
	"strings"
	"sync"
)

func (t *Twopc) isProposalTask(taskId, partyId string) bool {
	t.proposalTaskLock.RLock()
	defer t.proposalTaskLock.RUnlock()
	partyCache, ok := t.proposalTaskCache[taskId]
	if !ok {
		return false
	}
	_, ok = partyCache[partyId]
	if !ok {
		return false
	}
	return true
}
func (t *Twopc) isNotProposalTask(taskId, partyId string) bool {
	return !t.isProposalTask(taskId, partyId)
}

func (t *Twopc) hasOrgProposal(proposalId common.Hash, partyId string) bool {
	return t.state.HasOrgProposal(proposalId, partyId)
}
func (t *Twopc) hasNotOrgProposal(proposalId common.Hash, partyId string) bool {
	return t.state.HasNotOrgProposal(proposalId, partyId)
}

func (t *Twopc) addProposalTask(partyId string, task *types.ProposalTask) {
	t.proposalTaskLock.Lock()
	partyCache, ok := t.proposalTaskCache[task.GetTaskId()]
	if !ok {
		partyCache = make(map[string]*types.ProposalTask, 0)
	}
	partyCache[partyId] = task
	t.proposalTaskCache[task.GetTaskId()] = partyCache

	t.proposalTaskLock.Unlock()
}
func (t *Twopc) removeProposalTask(taskId, partyId string) {
	t.proposalTaskLock.Lock()
	partyCache, ok := t.proposalTaskCache[taskId]
	if ok {
		delete(partyCache, partyId)
		if len(partyCache) == 0 {
			delete(t.proposalTaskCache, taskId)
		}
	}
	t.proposalTaskLock.Unlock()
}

func (t *Twopc) hasProposalTask(taskId, partyId string) bool {
	t.proposalTaskLock.RLock()
	defer t.proposalTaskLock.RUnlock()
	partyCache, ok := t.proposalTaskCache[taskId]
	if !ok {
		return false
	}
	_, ok = partyCache[partyId]
	if !ok {
		return false
	}
	return true
}
func (t *Twopc) hasNotProposalTask(taskId, partyId string) bool {
	return !t.hasProposalTask(taskId, partyId)
}

func (t *Twopc) getProposalTask(taskId, partyId string) (*types.ProposalTask, bool) {
	t.proposalTaskLock.RLock()
	defer t.proposalTaskLock.RUnlock()
	partyCache, ok := t.proposalTaskCache[taskId]
	if !ok {
		return nil, false
	}
	task, ok := partyCache[partyId]
	return task, ok
}

func (t *Twopc) mustGetProposalTask(taskId, partyId string) *types.ProposalTask {
	t.proposalTaskLock.RLock()
	task, _ := t.getProposalTask(taskId, partyId)
	t.proposalTaskLock.RUnlock()
	return task
}

func (t *Twopc) hasPrepareVoting(proposalId common.Hash, org *apicommonpb.TaskOrganization) bool {
	return t.state.HasPrepareVoting(proposalId, org)
}

func (t *Twopc) storePrepareVote(vote *types.PrepareVote) {
	t.state.StorePrepareVote(vote)
}

//func (t *Twopc) removePrepareVote(proposalId common.Hash, partyId string, role apicommonpb.TaskRole) {
//	t.state.RemovePrepareVote(proposalId, partyId, role)
//}

func (t *Twopc) hasConfirmVoting(proposalId common.Hash, org *apicommonpb.TaskOrganization) bool {
	return t.state.HasConfirmVoting(proposalId, org)
}

func (t *Twopc) storeConfirmVote(vote *types.ConfirmVote) {
	t.state.StoreConfirmVote(vote)
}

//func (t *Twopc) removeConfirmVote(proposalId common.Hash, partyId string, role apicommonpb.TaskRole) {
//	t.state.RemoveConfirmVote(proposalId, partyId, role)
//}

func (t *Twopc) getTaskPrepareYesVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskPrepareYesVoteCount(proposalId)
}
func (t *Twopc) getTaskDataSupplierPrepareYesVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskDataSupplierPrepareYesVoteCount(proposalId)
}
func (t *Twopc) getTaskPowerSupplierPrepareYesVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskPowerSupplierPrepareYesVoteCount(proposalId)
}
func (t *Twopc) getTaskReceiverPrepareYesVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskReceiverPrepareYesVoteCount(proposalId)
}
func (t *Twopc) getTaskPrepareTotalVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskPrepareTotalVoteCount(proposalId)
}
func (t *Twopc) getTaskDataSupplierPrepareTotalVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskDataSupplierPrepareTotalVoteCount(proposalId)
}
func (t *Twopc) getTaskPowerSupplierPrepareTotalVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskPowerSupplierPrepareTotalVoteCount(proposalId)
}
func (t *Twopc) getTaskReceiverPrepareTotalVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskReceiverPrepareTotalVoteCount(proposalId)
}

func (t *Twopc) getTaskConfirmYesVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskConfirmYesVoteCount(proposalId)
}
func (t *Twopc) getTaskDataSupplierConfirmYesVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskDataSupplierConfirmYesVoteCount(proposalId)
}
func (t *Twopc) getTaskPowerSupplierConfirmYesVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskPowerSupplierConfirmYesVoteCount(proposalId)
}
func (t *Twopc) getTaskReceiverConfirmYesVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskReceiverConfirmYesVoteCount(proposalId)
}
func (t *Twopc) getTaskConfirmTotalVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskConfirmTotalVoteCount(proposalId)
}
func (t *Twopc) getTaskDataSupplierConfirmTotalVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskDataSupplierConfirmTotalVoteCount(proposalId)
}
func (t *Twopc) getTaskPowerSupplierConfirmTotalVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskPowerSupplierConfirmTotalVoteCount(proposalId)
}
func (t *Twopc) getTaskReceiverConfirmTotalVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskReceiverConfirmTotalVoteCount(proposalId)
}

func (t *Twopc) changeToConfirm(proposalId common.Hash, partyId string, startTime uint64) {
	t.state.ChangeToConfirm(proposalId, partyId, startTime)
}
func (t *Twopc) changeToCommit(proposalId common.Hash, partyId string, startTime uint64) {
	t.state.ChangeToCommit(proposalId, partyId, startTime)
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
	log.Debugf("Need SendTaskResultCh taskId: {%s}, result: {%s}", result.TaskId, result.String())
	if ch, ok := t.taskResultChSet[result.TaskId]; ok {
		log.Debugf("Start SendTaskResultCh taskId: {%s}, result: {%s}", result.TaskId, result.String())
		ch <- result
		close(ch)
		delete(t.taskResultChSet, result.TaskId)
	}
	t.taskResultLock.Unlock()
}

func (t *Twopc) sendNeedReplayScheduleTask(task *types.NeedReplayScheduleTask) {
	t.needReplayScheduleTaskCh <- task
}

func (t *Twopc) sendNeedExecuteTask(task *types.NeedExecuteTask) {
	t.needExecuteTaskCh <- task
}

func (t *Twopc) removeOrgProposalStateAndTask(proposalId common.Hash, partyId string) {
	if state := t.state.GetProposalState(proposalId); state.IsNotEmpty() {
		log.Infof("Start remove org proposalState and task cache on Consensus, proposalId {%s}, taskId {%s}", proposalId, state.GetTaskId())
		t.removeProposalTask(state.GetTaskId(), partyId)
		t.state.CleanOrgProposalState(proposalId, partyId)
		go func() {
			t.db.DeleteState(t.db.GetPrepareVotesKey(proposalId, partyId))
			t.db.DeleteState(t.db.GetConfirmVotesKey(proposalId, partyId))
			t.db.DeleteState(t.db.GetProposalSetKey(proposalId, partyId))
		}()
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
	t.db.UpdateOrgProposalState(proposalId, pstate.GetTaskSender(), orgState)
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
	t.db.DeleteState(t.db.GetProposalSetKey(proposalId, partyId))
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

func (t *Twopc) makeConfirmTaskPeerDesc(proposalId common.Hash) *twopcpb.ConfirmTaskPeerInfo {

	var sender *twopcpb.TaskPeerInfo
	dataSuppliers, powerSuppliers, receivers := make([]*twopcpb.TaskPeerInfo, 0), make([]*twopcpb.TaskPeerInfo, 0), make([]*twopcpb.TaskPeerInfo, 0)

	for _, vote := range t.state.GetPrepareVoteArr(proposalId) {

		if vote.MsgOption.SenderRole == apicommonpb.TaskRole_TaskRole_Sender && nil != vote.PeerInfo {
			sender = types.ConvertTaskPeerInfo(vote.PeerInfo)
		}

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
		OwnerPeerInfo:              sender,
		DataSupplierPeerInfoList:   dataSuppliers,
		PowerSupplierPeerInfoList:  powerSuppliers,
		ResultReceiverPeerInfoList: receivers,
	}
}

func (t *Twopc) refreshProposalState() {

	identity, err := t.resourceMng.GetDB().GetIdentity()
	if nil != err {
		log.Errorf("Failed to query local identity on consensus.refreshProposalState(), err: {%s}", err)
		return
	}

	pid, err := p2p.HexPeerID(identity.GetNodeId())
	if nil != err {
		log.Errorf("Failed to convert on consensus.refreshProposalState(), err: {%s}", err)
		return
	}

	t.state.proposalsLock.Lock()

	for proposalId, pstate := range t.state.proposalSet {

		pstate.MustLock()
		for partyId, orgState := range pstate.GetStateCache() {

			if orgState.IsDeadline() {
				log.Debugf("Started refresh proposalState loop, the proposalState direct be deadline, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
					pstate.GetProposalId().String(), pstate.GetTaskId(), partyId)

				pstate.RemoveOrgProposalStateUnSafe(partyId)
				t.removeOrgProposalStateAndTask(pstate.GetProposalId(), partyId)
				t.TaskConsensusInterrupt(proposalId, pid, pstate.GetTaskId(), partyId, identity, pstate.GetTaskSender())

				t.db.DeleteState(t.db.GetProposalSetKey(proposalId, partyId))
				continue
			}

			switch orgState.CurrPeriodNum() {

			case ctypes.PeriodPrepare:

				if orgState.IsPrepareTimeout() {
					log.Debugf("Started refresh org proposalState, the org proposalState was prepareTimeout, change to confirm epoch, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
						pstate.GetProposalId().String(), pstate.GetTaskId(), partyId)

					orgState.ChangeToConfirm(orgState.PeriodStartTime + uint64(ctypes.PrepareMsgVotingTimeout.Milliseconds()))
					pstate.StoreOrgProposalStateUnSafe(orgState)

					t.db.UpdateOrgProposalState(pstate.GetProposalId(), pstate.GetTaskSender(), orgState)
				}

			case ctypes.PeriodConfirm:

				if orgState.IsConfirmTimeout() {
					log.Debugf("Started refresh org proposalState, the org proposalState was confirmTimeout, change to commit epoch, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
						pstate.GetProposalId().String(), pstate.GetTaskId(), partyId)

					orgState.ChangeToCommit(orgState.PeriodStartTime + uint64(ctypes.ConfirmMsgVotingTimeout.Milliseconds()))
					pstate.StoreOrgProposalStateUnSafe(orgState)

					t.db.UpdateOrgProposalState(pstate.GetProposalId(), pstate.GetTaskSender(), orgState)
				}

			case ctypes.PeriodCommit:

				if orgState.IsCommitTimeout() {
					log.Debugf("Started refresh org proposalState, the org proposalState was commitTimeout, change to finished epoch, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
						pstate.GetProposalId().String(), pstate.GetTaskId(), partyId)

					orgState.ChangeToFinished(orgState.PeriodStartTime + uint64(ctypes.CommitMsgEndingTimeout.Milliseconds()))
					pstate.StoreOrgProposalStateUnSafe(orgState)

					t.db.UpdateOrgProposalState(pstate.GetProposalId(), pstate.GetTaskSender(), orgState)
				}

			case ctypes.PeriodFinished:

				if orgState.IsDeadline() {
					log.Debugf("Started refresh org proposalState, the org proposalState was finished, but coming deadline now, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
						pstate.GetProposalId().String(), pstate.GetTaskId(), partyId)

					pstate.RemoveOrgProposalStateUnSafe(partyId)
					t.removeOrgProposalStateAndTask(pstate.GetProposalId(), partyId)
					t.TaskConsensusInterrupt(proposalId, pid, pstate.GetTaskId(), partyId, identity, pstate.GetTaskSender())

					t.db.DeleteState(t.db.GetProposalSetKey(proposalId, partyId))
				}

			default:
				log.Errorf("Unknown the proposalState period,  proposalId: {%s}, taskId: {%s}, partyId: {%s}", pstate.GetProposalId().String(), pstate.GetTaskId(), partyId)
			}
		}

		pstate.MustUnLock()
	}

	t.state.proposalsLock.Unlock()
}

func (t *Twopc) TaskConsensusInterrupt(
	proposalId common.Hash,
	localPid peer.ID,
	taskId, partyId string,
	identity *apicommonpb.Organization,
	taskSender *apicommonpb.TaskOrganization,
) {
	// Send task result msg to remote peer, if current org identityId is not task sender identityId
	if identity.GetIdentityId() == taskSender.GetIdentityId() {
		t.replyTaskConsensusResult(types.NewTaskConsResult(taskId, types.TaskConsensusInterrupt, fmt.Errorf("the task proposalState coming deadline")))
	} else {

		if t.hasNotProposalTask(taskId, partyId) {
			t.removeOrgProposalStateAndTask(proposalId, partyId)
			return
		}

		t.resourceMng.GetDB().StoreTaskEvent(&libtypes.TaskEvent{
			Type:       evengine.TaskProposalStateDeadline.Type,
			IdentityId: identity.GetIdentityId(),
			TaskId:     taskId,
			Content:    fmt.Sprintf("%s for myself", evengine.TaskProposalStateDeadline.Msg),
			CreateAt:   timeutils.UnixMsecUint64(),
		})

		t.sendNeedExecuteTask(types.NewNeedExecuteTask(
			localPid,
			proposalId,
			apicommonpb.TaskRole_TaskRole_Unknown,
			&apicommonpb.TaskOrganization{
				PartyId:    partyId,
				NodeName:   identity.GetNodeName(),
				NodeId:     identity.GetNodeId(),
				IdentityId: identity.GetIdentityId(),
			},
			apicommonpb.TaskRole_TaskRole_Sender,
			taskSender,
			t.mustGetProposalTask(taskId, partyId).Task,
			types.TaskConsensusInterrupt,
			nil,
			nil,
		))
	}
}

func (t *Twopc) driveTask(
	pid peer.ID,
	proposalId common.Hash,
	localTaskRole apicommonpb.TaskRole,
	localTaskOrganization *apicommonpb.TaskOrganization,
	remoteTaskRole apicommonpb.TaskRole,
	remoteTaskOrganization *apicommonpb.TaskOrganization,
	task *types.Task,
) {

	log.Debugf("Start to call `driveTask`, proposalId: {%s}, taskId: {%s}, localTaskRole: {%s}, partyId: {%s}, identityId: {%s}, nodeName: {%s}",
		proposalId.String(), task.GetTaskId(), localTaskRole.String(), localTaskOrganization.GetPartyId(), localTaskOrganization.GetIdentityId(), localTaskOrganization.GetNodeName())

	selfVote := t.getPrepareVote(proposalId, localTaskOrganization.GetPartyId())
	if nil == selfVote {
		log.Errorf("Failed to find local cache about prepareVote myself internal resource, proposalId: {%s}, taskId: {%s}, localTaskRole: {%s}, partyId: {%s}, identityId: {%s}, nodeName: {%s}",
			proposalId.String(), task.GetTaskId(), localTaskRole.String(), localTaskOrganization.GetPartyId(), localTaskOrganization.GetIdentityId(), localTaskOrganization.GetNodeName())

		return
	}

	peers, ok := t.getConfirmTaskPeerInfo(proposalId)
	if !ok {
		log.Errorf("Failed to find local cache about prepareVote all peer resource {externalIP:externalPORT}, proposalId: {%s}, taskId: {%s}, localTaskRole: {%s}, partyId: {%s}, identityId: {%s}, nodeName: {%s}",
			proposalId.String(), task.GetTaskId(), localTaskRole.String(), localTaskOrganization.GetPartyId(), localTaskOrganization.GetIdentityId(), localTaskOrganization.GetNodeName())

		return
	}

	// Send task to TaskManager to execute
	t.sendNeedExecuteTask(types.NewNeedExecuteTask(
		pid,
		proposalId,
		localTaskRole,
		localTaskOrganization,
		remoteTaskRole,
		remoteTaskOrganization,
		task,
		types.TaskNeedExecute,
		selfVote.PeerInfo,
		peers,
	))

	t.removeOrgProposalStateAndTask(proposalId, localTaskOrganization.GetPartyId())

}

func (t *Twopc) sendPrepareMsg(proposalId common.Hash, task *types.Task, startTime uint64) error {

	sender := task.GetTaskSender()

	sendTaskFn := func(wg *sync.WaitGroup, sender, receiver *apicommonpb.TaskOrganization, senderRole, receiverRole apicommonpb.TaskRole, errCh chan<- error) {

		defer wg.Done()

		var pid, err = p2p.HexPeerID(receiver.NodeId)
		if nil != err {
			errCh <- fmt.Errorf("failed to nodeId => peerId, proposalId: %s, taskId: %s, other peer taskRole: %s, other peer taskPartyId: %s, identityId: %s, pid: %s, err: %s",
				proposalId.String(), task.GetTaskId(), receiverRole.String(), receiver.GetPartyId(), receiver.GetIdentityId(), pid, err)
			return
		}

		prepareMsg, err := makePrepareMsg(proposalId, senderRole, receiverRole, sender.GetPartyId(), receiver.GetPartyId(), task, startTime)

		if nil != err {
			errCh <- fmt.Errorf("failed to make prepareMsg, proposalId: %s, taskId: %s, other peer taskRole: %s, other peer taskPartyId: %s, identityId: %s, pid: %s, err: %s",
				proposalId.String(), task.GetTaskId(), receiverRole.String(), receiver.GetPartyId(), receiver.GetIdentityId(), pid, err)
			return
		}

		var sendErr error
		if types.IsSameTaskOrg(sender, receiver) {
			sendErr = t.sendLocalPrepareMsg(pid, prepareMsg)
		} else {
			sendErr = handler.SendTwoPcPrepareMsg(context.TODO(), t.p2p, pid, prepareMsg)
		}

		if nil != sendErr {
			errCh <- fmt.Errorf("failed to call `SendTwoPcPrepareMsg` proposalId: %s, taskId: %s, other peer taskRole: %s, other peer taskPartyId: %s, identityId: %s, pid: %s, err: %s",
				proposalId.String(), task.GetTaskId(), receiverRole.String(), receiver.GetPartyId(), receiver.GetIdentityId(), pid, sendErr)
			return
		}

		log.Debugf("Succceed to call `SendTwoPcPrepareMsg` proposalId: %s, taskId: %s, other peer taskRole: %s, other peer taskPartyId: %s, identityId: %s, pid: %s",
			proposalId.String(), task.GetTaskId(), receiverRole.String(), receiver.GetPartyId(), receiver.GetIdentityId(), pid)
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

func (t *Twopc) sendPrepareVote(pid peer.ID, sender, receiver *apicommonpb.TaskOrganization, req *twopcpb.PrepareVote) error {
	if types.IsNotSameTaskOrg(sender, receiver) {
		return handler.SendTwoPcPrepareVote(context.TODO(), t.p2p, pid, req)
	} else {
		return t.sendLocalPrepareVote(pid, req)
	}
}

func (t *Twopc) sendConfirmMsg(proposalId common.Hash, task *types.Task, peers *twopcpb.ConfirmTaskPeerInfo, startTime uint64) error {

	sender := task.GetTaskSender()

	sendConfirmMsgFn := func(wg *sync.WaitGroup, sender, receiver *apicommonpb.TaskOrganization, senderRole, receiverRole apicommonpb.TaskRole, errCh chan<- error) {

		defer wg.Done()

		pid, err := p2p.HexPeerID(receiver.NodeId)
		if nil != err {
			errCh <- fmt.Errorf("failed to nodeId => peerId, proposalId: %s, taskId: %s, other peer's taskRole: %s, other peer's partyId: %s, other identityId: %s, pid: %s, err: %s",
				proposalId.String(), task.GetTaskId(), receiverRole.String(), receiver.GetPartyId(), receiver.GetIdentityId(), pid, err)
			return
		}

		confirmMsg := makeConfirmMsg(proposalId, senderRole, receiverRole, sender.GetPartyId(), receiver.GetPartyId(), task, peers, startTime)

		var sendErr error
		if types.IsSameTaskOrg(sender, receiver) {
			sendErr = t.sendLocalConfirmMsg(pid, confirmMsg)
		} else {
			sendErr = handler.SendTwoPcConfirmMsg(context.TODO(), t.p2p, pid, confirmMsg)
		}

		// Send the ConfirmMsg to other peer
		if nil != sendErr {
			errCh <- fmt.Errorf("failed to call`SendTwoPcConfirmMsg` proposalId: %s, taskId: %s,other peer's taskRole: %s, other peer's partyId: %s, other identityId: %s, pid: %s, err: %s",
				proposalId.String(), task.GetTaskId(), receiverRole.String(), receiver.GetPartyId(), receiver.GetIdentityId(), pid, err)
			errCh <- err
			return
		}

		log.Debugf("Succceed to call`SendTwoPcConfirmMsg` proposalId: %s, taskId: %s,other peer's taskRole: %s, other peer's partyId: %s, other identityId: %s, pid: %s",
			proposalId.String(), task.GetTaskId(), receiverRole.String(), receiver.GetPartyId(), receiver.GetIdentityId(), pid)

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
		return handler.SendTwoPcConfirmVote(context.TODO(), t.p2p, pid, req)
	} else {
		return t.sendLocalConfirmVote(pid, req)
	}
}

func (t *Twopc) sendCommitMsg(proposalId common.Hash, task *types.Task, startTime uint64) error {

	sender := task.GetTaskSender()

	sendCommitMsgFn := func(wg *sync.WaitGroup, sender, receiver *apicommonpb.TaskOrganization, senderRole, receiverRole apicommonpb.TaskRole, errCh chan<- error) {

		defer wg.Done()

		pid, err := p2p.HexPeerID(receiver.NodeId)
		if nil != err {
			errCh <- fmt.Errorf("failed to nodeId => peerId, proposalId: %s, taskId: %s, other peer's taskRole: %s, other peer's partyId: %s, identityId: %s, pid: %s, err: %s",
				proposalId.String(), task.GetTaskId(), receiverRole.String(), receiver.GetPartyId(), receiver.GetIdentityId(), pid, err)
			return
		}

		commitMsg := makeCommitMsg(proposalId, senderRole, receiverRole, sender.GetPartyId(), receiver.GetPartyId(), task, startTime)

		var sendErr error
		if types.IsSameTaskOrg(sender, receiver) {
			sendErr = t.sendLocalCommitMsg(pid, commitMsg)
		} else {
			sendErr = handler.SendTwoPcCommitMsg(context.TODO(), t.p2p, pid, commitMsg)
		}

		// Send the ConfirmMsg to other peer
		if nil != sendErr {
			errCh <- fmt.Errorf("failed to call`SendTwoPcCommitMsg` proposalId: %s, taskId: %s,  other peer's taskRole: %s, other peer's partyId: %s, identityId: %s, pid: %s, err: %s",
				proposalId.String(), task.GetTaskId(), receiverRole.String(), receiver.GetPartyId(), receiver.GetIdentityId(), pid, err)
			errCh <- err
			return
		}

		log.Debugf("Succceed to call`SendTwoPcCommitMsg` proposalId: %s, taskId: %s,  other peer's taskRole: %s, other peer's partyId: %s, identityId: %s, pid: %s",
			proposalId.String(), task.GetTaskId(), receiverRole.String(), receiver.GetPartyId(), receiver.GetIdentityId(), pid)

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

func (t *Twopc) isIncludeSenderParty(partyId, identityId string, role apicommonpb.TaskRole, task *types.Task) bool {
	if t.getIncludeSenderPartyCount(partyId, identityId, role, task) != 0 {
		return true
	}
	return false
}

func (t *Twopc) getIncludeSenderPartyCount(partyId, identityId string, role apicommonpb.TaskRole, task *types.Task) uint32 {

	var count int
	switch role {
	case apicommonpb.TaskRole_TaskRole_DataSupplier:
		for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
			if partyId == dataSupplier.GetOrganization().GetPartyId() && identityId == dataSupplier.GetOrganization().GetIdentityId() {
				count++
			}
		}
	case apicommonpb.TaskRole_TaskRole_PowerSupplier:
		for _, powerSupplier := range task.GetTaskData().GetPowerSuppliers() {
			if partyId == powerSupplier.GetOrganization().GetPartyId() && identityId == powerSupplier.GetOrganization().GetIdentityId() {
				count++
			}
		}
	case apicommonpb.TaskRole_TaskRole_Receiver:
		for _, receiver := range task.GetTaskData().GetReceivers() {
			if partyId == receiver.GetPartyId() && identityId == receiver.GetIdentityId() {
				count++
			}
		}
	}
	return uint32(count)
}
func (t *Twopc) getNeedVotingCount(role apicommonpb.TaskRole, task *types.Task) uint32 {
	sender := task.GetTaskSender()
	includeCount := t.getIncludeSenderPartyCount(sender.GetPartyId(), sender.GetIdentityId(), role, task)
	switch role {
	case apicommonpb.TaskRole_TaskRole_DataSupplier:
		return uint32(len(task.GetTaskData().GetDataSuppliers())) - includeCount
	case apicommonpb.TaskRole_TaskRole_PowerSupplier:
		return uint32(len(task.GetTaskData().GetPowerSuppliers())) - includeCount
	case apicommonpb.TaskRole_TaskRole_Receiver:
		return uint32(len(task.GetTaskData().GetReceivers())) - includeCount
	}
	return 0
}

func (t *Twopc) verifyPrepareVoteRole(proposalId common.Hash, partyId, identityId string, role apicommonpb.TaskRole, task *types.Task) (bool, error) {
	var identityValid bool
	switch role {
	case apicommonpb.TaskRole_TaskRole_DataSupplier:

		dataSupplierCount := len(task.GetTaskData().GetDataSuppliers())
		includeCount := t.getIncludeSenderPartyCount(partyId, identityId, role, task)

		if (uint32(dataSupplierCount) - includeCount) == t.getTaskDataSupplierPrepareTotalVoteCount(proposalId) {
			return false, fmt.Errorf("%s, on the prepare vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
				ctypes.ErrVoteCountOverflow, task.GetTaskData().GetTaskId(), role.String(),
				identityId, partyId)
		}

		for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
			// identity + partyId
			if dataSupplier.GetOrganization().GetIdentityId() == identityId && dataSupplier.GetOrganization().GetPartyId() == partyId {
				identityValid = true
				break
			}
		}
	case apicommonpb.TaskRole_TaskRole_PowerSupplier:

		powerSupplierCount := len(task.GetTaskData().GetPowerSuppliers())
		includeCount := t.getIncludeSenderPartyCount(partyId, identityId, role, task)

		if (uint32(powerSupplierCount) - includeCount) == t.getTaskPowerSupplierPrepareTotalVoteCount(proposalId) {
			return false, fmt.Errorf("%s, on the prepare vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
				ctypes.ErrVoteCountOverflow, task.GetTaskData().GetTaskId(), role.String(),
				identityId, partyId)
		}

		for _, powerSupplier := range task.GetTaskData().GetPowerSuppliers() {
			// identity + partyId
			if powerSupplier.GetOrganization().GetIdentityId() == identityId && powerSupplier.GetOrganization().GetPartyId() == partyId {
				identityValid = true
				break
			}
		}
	case apicommonpb.TaskRole_TaskRole_Receiver:

		receiverCount := len(task.GetTaskData().GetReceivers())
		includeCount := t.getIncludeSenderPartyCount(partyId, identityId, role, task)

		if (uint32(receiverCount) - includeCount) == t.getTaskReceiverPrepareTotalVoteCount(proposalId) {
			return false, fmt.Errorf("%s, on the prepare vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
				ctypes.ErrVoteCountOverflow, task.GetTaskData().GetTaskId(), role.String(),
				identityId, partyId)
		}

		for _, receiver := range task.GetTaskData().GetReceivers() {
			// identity + partyId
			if receiver.GetIdentityId() == identityId && receiver.GetPartyId() == partyId {
				identityValid = true
				break
			}
		}
	default:
		return false, fmt.Errorf("%s, on the prepare vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
			ctypes.ErrMsgOwnerNodeIdInvalid, task.GetTaskData().GetTaskId(), role.String(),
			identityId, partyId)

	}
	return identityValid, nil
}
func (t *Twopc) verifyConfirmVoteRole(proposalId common.Hash, partyId, identityId string, role apicommonpb.TaskRole, task *types.Task) (bool, error) {
	var identityValid bool
	switch role {
	case apicommonpb.TaskRole_TaskRole_DataSupplier:

		dataSupplierCount := len(task.GetTaskData().GetDataSuppliers())
		includeCount := t.getIncludeSenderPartyCount(partyId, identityId, role, task)

		if (uint32(dataSupplierCount) - includeCount) == t.getTaskDataSupplierConfirmTotalVoteCount(proposalId) {
			return false, fmt.Errorf("%s, on the confirm vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
				ctypes.ErrVoteCountOverflow, task.GetTaskData().GetTaskId(), role.String(), identityId, partyId)
		}

		for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
			// identity + partyId
			if dataSupplier.GetOrganization().GetIdentityId() == identityId && dataSupplier.GetOrganization().GetPartyId() == partyId {
				identityValid = true
				break
			}
		}
	case apicommonpb.TaskRole_TaskRole_PowerSupplier:

		powerSupplierCount := len(task.GetTaskData().GetPowerSuppliers())
		includeCount := t.getIncludeSenderPartyCount(partyId, identityId, role, task)

		if (uint32(powerSupplierCount) - includeCount) == t.getTaskPowerSupplierConfirmTotalVoteCount(proposalId) {
			return false, fmt.Errorf("%s, on the confirm vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
				ctypes.ErrVoteCountOverflow, task.GetTaskData().GetTaskId(), role.String(), identityId, partyId)
		}

		for _, powerSupplier := range task.GetTaskData().GetPowerSuppliers() {
			// identity + partyId
			if powerSupplier.GetOrganization().GetIdentityId() == identityId && powerSupplier.GetOrganization().GetPartyId() == partyId {
				identityValid = true
				break
			}
		}
	case apicommonpb.TaskRole_TaskRole_Receiver:

		receiverCount := len(task.GetTaskData().GetReceivers())
		includeCount := t.getIncludeSenderPartyCount(partyId, identityId, role, task)

		if (uint32(receiverCount) - includeCount) == t.getTaskReceiverConfirmTotalVoteCount(proposalId) {
			return false, fmt.Errorf("%s, on the confirm vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
				ctypes.ErrVoteCountOverflow, task.GetTaskData().GetTaskId(), role.String(), identityId, partyId)
		}

		for _, receiver := range task.GetTaskData().GetReceivers() {
			// identity + partyId
			if receiver.GetIdentityId() == identityId && receiver.GetPartyId() == partyId {
				identityValid = true
				break
			}
		}
	default:
		return false, fmt.Errorf("%s, on the confirm vote [taskId: %s, taskRole: %s, identity: %s, partyId: %s]",
			ctypes.ErrMsgOwnerNodeIdInvalid, task.GetTaskData().GetTaskId(), role.String(), identityId, partyId)

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
