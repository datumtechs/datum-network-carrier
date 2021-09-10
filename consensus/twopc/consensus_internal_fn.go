package twopc

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	"github.com/RosettaFlow/Carrier-Go/core/evengine"
	"github.com/RosettaFlow/Carrier-Go/handler"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/libp2p/go-libp2p-core/peer"
	"strings"
	"sync"
)

func (t *TwoPC) isProposalTask(taskId string) bool {
	t.proposalTaskLock.RLock()
	_, ok := t.proposalTaskCache[taskId]
	t.proposalTaskLock.RUnlock()
	if !ok {
		return true
	}
	return false
}
func (t *TwoPC) isNotProposalTask(taskId string) bool { return !t.isProposalTask(taskId) }

func (t *TwoPC) hasOrgProposal(proposalId common.Hash, partyId string) bool {
	return t.state.HasOrgProposal(proposalId, partyId)
}
func (t *TwoPC) hasNotOrgProposal(proposalId common.Hash, partyId string) bool {
	return t.state.HasNotOrgProposal(proposalId, partyId)
}

func (t *TwoPC) addProposalTask(task *types.ProposalTask) {
	t.proposalTaskLock.Lock()
	_, ok := t.proposalTaskCache[task.GetTaskId()]
	if !ok {
		t.proposalTaskCache[task.GetTaskId()] = task
	}
	t.proposalTaskLock.Unlock()
}
func (t *TwoPC) removeProposalTask(taskId string) {
	t.proposalTaskLock.Lock()
	delete(t.proposalTaskCache, taskId)
	t.proposalTaskLock.Unlock()
}
func (t *TwoPC) hasProposalTask(taskId string) bool {
	t.proposalTaskLock.RLock()
	_, ok := t.proposalTaskCache[taskId]
	t.proposalTaskLock.RUnlock()
	if ok {
		return true
	}
	return false
}
func (t *TwoPC) hasNotProposalTask(taskId string) bool { return !t.hasProposalTask(taskId) }

func (t *TwoPC) getProposalTask(taskId string) (*types.ProposalTask, bool) {
	t.proposalTaskLock.RLock()
	task, ok := t.proposalTaskCache[taskId]
	t.proposalTaskLock.RUnlock()
	return task, ok
}

func (t *TwoPC) mustGetProposalTask(taskId string) *types.ProposalTask {
	t.proposalTaskLock.RLock()
	task, _ := t.getProposalTask(taskId)
	t.proposalTaskLock.RUnlock()
	return task
}

func (t *TwoPC) hasPrepareVoting(proposalId common.Hash, org *apipb.TaskOrganization) bool {
	return t.state.HasPrepareVoting(proposalId, org)
}

func (t *TwoPC) storePrepareVote(vote *types.PrepareVote) {
	t.state.StorePrepareVote(vote)
}

func (t *TwoPC) removePrepareVote(proposalId common.Hash, partyId string, role apipb.TaskRole) {
	t.state.RemovePrepareVote(proposalId, partyId, role)
}

func (t *TwoPC) hasConfirmVoting(proposalId common.Hash, org *apipb.TaskOrganization) bool {
	return t.state.HasConfirmVoting(proposalId, org)
}

func (t *TwoPC) storeConfirmVote(vote *types.ConfirmVote) {
	t.state.StoreConfirmVote(vote)
}

func (t *TwoPC) removeConfirmVote(proposalId common.Hash, partyId string, role apipb.TaskRole) {
	t.state.RemoveConfirmVote(proposalId, partyId, role)
}

func (t *TwoPC) getTaskPrepareYesVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskPrepareYesVoteCount(proposalId)
}
func (t *TwoPC) getTaskDataSupplierPrepareYesVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskDataSupplierPrepareYesVoteCount(proposalId)
}
func (t *TwoPC) getTaskPowerSupplierPrepareYesVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskPowerSupplierPrepareYesVoteCount(proposalId)
}
func (t *TwoPC) getTaskReceiverPrepareYesVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskReceiverPrepareYesVoteCount(proposalId)
}
func (t *TwoPC) getTaskPrepareTotalVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskPrepareTotalVoteCount(proposalId)
}
func (t *TwoPC) getTaskDataSupplierPrepareTotalVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskDataSupplierPrepareTotalVoteCount(proposalId)
}
func (t *TwoPC) getTaskPowerSupplierPrepareTotalVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskPowerSupplierPrepareTotalVoteCount(proposalId)
}
func (t *TwoPC) getTaskReceiverPrepareTotalVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskReceiverPrepareTotalVoteCount(proposalId)
}

func (t *TwoPC) getTaskConfirmYesVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskConfirmYesVoteCount(proposalId)
}
func (t *TwoPC) getTaskDataSupplierConfirmYesVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskDataSupplierConfirmYesVoteCount(proposalId)
}
func (t *TwoPC) getTaskPowerSupplierConfirmYesVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskPowerSupplierConfirmYesVoteCount(proposalId)
}
func (t *TwoPC) getTaskReceiverConfirmYesVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskReceiverConfirmYesVoteCount(proposalId)
}
func (t *TwoPC) getTaskConfirmTotalVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskConfirmTotalVoteCount(proposalId)
}
func (t *TwoPC) getTaskDataSupplierConfirmTotalVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskDataSupplierConfirmTotalVoteCount(proposalId)
}
func (t *TwoPC) getTaskPowerSupplierConfirmTotalVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskPowerSupplierConfirmTotalVoteCount(proposalId)
}
func (t *TwoPC) getTaskReceiverConfirmTotalVoteCount(proposalId common.Hash) uint32 {
	return t.state.GetTaskReceiverConfirmTotalVoteCount(proposalId)
}

func (t *TwoPC) changeToConfirm(proposalId common.Hash, partyId string, startTime uint64) {
	t.state.ChangeToConfirm(proposalId, partyId, startTime)
}
func (t *TwoPC) changeToCommit(proposalId common.Hash, partyId string, startTime uint64) {
	t.state.ChangeToCommit(proposalId, partyId, startTime)
}

func (t *TwoPC) addTaskResultCh(taskId string, resultCh chan<- *types.TaskConsResult) {
	t.taskResultLock.Lock()
	log.Debugf("AddTaskResultCh taskId: {%s}", taskId)
	t.taskResultChSet[taskId] = resultCh
	t.taskResultLock.Unlock()
}
func (t *TwoPC) removeTaskResultCh(taskId string) {
	t.taskResultLock.Lock()
	log.Debugf("RemoveTaskResultCh taskId: {%s}", taskId)
	delete(t.taskResultChSet, taskId)
	t.taskResultLock.Unlock()
}
func (t *TwoPC) replyTaskConsensusResult(result *types.TaskConsResult) {
	t.taskResultBusCh <- result
}
func (t *TwoPC) handleTaskConsensusResult(result *types.TaskConsResult) {
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

func (t *TwoPC) sendNeedReplayScheduleTask(task *types.NeedReplayScheduleTask) {
	t.needReplayScheduleTaskCh <- task
}

func (t *TwoPC) sendNeedExecuteTask(task *types.NeedExecuteTask) {
	t.needExecuteTaskCh <- task
}

func (t *TwoPC) storeProposalState(proposalState *ctypes.ProposalState) {
	t.state.StoreProposalState(proposalState)
}
func (t *TwoPC) removeProposalState(proposalId common.Hash) {
	t.state.CleanProposalState(proposalId)
}
func (t *TwoPC) removeProposalStateAndTask(proposalId common.Hash) {
	if state := t.state.GetProposalState(proposalId); state.IsNotEmpty() {
		log.Infof("Start remove proposalState and task cache on Consensus, proposalId {%s}, taskId {%s}", proposalId, state.GetTaskId())
		t.removeProposalTask(state.GetTaskId())
		t.state.CleanProposalState(proposalId)

	}
}
func (t *TwoPC) storeOrgProposalState (proposalId common.Hash, taskId string, sender *apipb.TaskOrganization, orgState *ctypes.OrgProposalState) {
	pstate := t.state.GetProposalState(proposalId)

	var first bool
	if pstate.IsEmpty() {
		pstate = ctypes.NewProposalState(proposalId, taskId, sender)
		first = true
	}
	pstate.StoreOrgProposalState(orgState)
	if first {
		t.state.StoreProposalState(pstate)
	} else {
		t.state.UpdateProposalState(pstate)
	}
}

func (t *TwoPC) removeOrgProposalState (proposalId common.Hash, partyId string) {
	pstate := t.state.GetProposalState(proposalId)

	if pstate.IsEmpty() {
		return
	}
	pstate.RemoveOrgProposalState(partyId)
}

func (t *TwoPC) getOrgProposalState (proposalId common.Hash, partyId string) (*ctypes.OrgProposalState, bool) {
	pstate := t.state.GetProposalState(proposalId)
	if pstate.IsEmpty() {
		return nil, false
	}
	return pstate.GetOrgProposalState(partyId)
}

func (t *TwoPC) mustGetOrgProposalState (proposalId common.Hash, partyId string) *ctypes.OrgProposalState {
	pstate := t.state.GetProposalState(proposalId)
	if pstate.IsEmpty() {
		return nil
	}
	return pstate.MustGetOrgProposalState(partyId)
}

func (t *TwoPC) makeConfirmTaskPeerDesc (proposalId common.Hash) *pb.ConfirmTaskPeerInfo {

	var sender *pb.TaskPeerInfo
	dataSuppliers, powerSuppliers, receivers := make([]*pb.TaskPeerInfo, 0), make([]*pb.TaskPeerInfo, 0), make([]*pb.TaskPeerInfo, 0)

	for _, vote := range t.state.GetPrepareVoteArr(proposalId) {

		if vote.MsgOption.SenderRole == apipb.TaskRole_TaskRole_Sender && nil != vote.PeerInfo {
			sender = types.ConvertTaskPeerInfo(vote.PeerInfo)
		}

		if vote.MsgOption.SenderRole == apipb.TaskRole_TaskRole_DataSupplier && nil != vote.PeerInfo {
			dataSuppliers = append(dataSuppliers, types.ConvertTaskPeerInfo(vote.PeerInfo))
		}
		if vote.MsgOption.SenderRole == apipb.TaskRole_TaskRole_PowerSupplier && nil != vote.PeerInfo {
			powerSuppliers = append(powerSuppliers, types.ConvertTaskPeerInfo(vote.PeerInfo))
		}
		if vote.MsgOption.SenderRole == apipb.TaskRole_TaskRole_Receiver && nil != vote.PeerInfo {
			receivers = append(receivers, types.ConvertTaskPeerInfo(vote.PeerInfo))
		}
	}
	return &pb.ConfirmTaskPeerInfo{
		OwnerPeerInfo:              sender,
		DataSupplierPeerInfoList:   dataSuppliers,
		PowerSupplierPeerInfoList:  powerSuppliers,
		ResultReceiverPeerInfoList: receivers,
	}
}

func (t *TwoPC) refreshProposalState () {

	t.state.proposalsLock.Lock()

	for proposalId, proposalState := range t.state.proposalSet {

		proposalState.RefreshProposalState()

		if proposalState.IsEmpty() {

			log.Debugf("Start clean proposalState and task on consensus.refreshProposalState(), proposalId: {%s}, taskId: {%s}",
				proposalId.String(), proposalState.GetTaskId())

			has, err := t.resourceMng.GetDB().HasLocalTaskExecute(proposalState.GetTaskId())
			if nil != err {
				log.Errorf("Failed to query local task exec status with task on consensus.refreshProposalState(),  proposalId: {%s}, taskId: {%s}, err: {%s}",
					proposalId.String(), proposalState.GetTaskId(), err)

				t.removeProposalStateAndTask(proposalState.GetProposalId())
				continue
			}

			if has {
				log.Debugf("The local task have been executing, direct clean proposalStateAndTaskCache of consensus on consensus.refreshProposalState(), proposalId: {%s}, taskId: {%s}",
					proposalId.String(), proposalState.GetTaskId())

				t.removeProposalStateAndTask(proposalState.GetProposalId())
				continue
			}


			identity, err := t.resourceMng.GetDB().GetIdentity()
			if nil != err {
				log.Errorf("Failed to query local identity on consensus.refreshProposalState(), proposalId: {%s}, taskId: {%s}, err: {%s}",
					proposalId.String(), proposalState.GetTaskId(), err)

				t.removeProposalStateAndTask(proposalState.GetProposalId())
				continue
			}

			// Send task result msg to remote peer, if current org identityId is not task sender identityId
			if identity.GetIdentityId() == proposalState.GetTaskSender().GetIdentityId() {
				t.replyTaskConsensusResult(types.NewTaskConsResult(proposalState.GetTaskId(), types.TaskConsensusInterrupt, fmt.Errorf("the task proposalState coming deadline")))
			} else {

				if t.hasNotProposalTask(proposalState.GetTaskId()) {
					t.removeProposalStateAndTask(proposalState.GetProposalId())
					continue
				}

				pid, err := p2p.HexPeerID(identity.GetNodeId())
				if nil == err {

					t.resourceMng.GetDB().StoreTaskEvent(&libTypes.TaskEvent{
						Type:       evengine.TaskProposalStateDeadline.Type,
						IdentityId: identity.GetIdentityId(),
						TaskId:     proposalState.GetTaskId(),
						Content:    fmt.Sprintf("%s for myself", evengine.TaskProposalStateDeadline.Msg),
						CreateAt:   uint64(timeutils.UnixMsec()),
					})

					t.sendNeedExecuteTask(types.NewNeedExecuteTask(
						pid,
						proposalId,
						apipb.TaskRole_TaskRole_Unknown,
						&apipb.TaskOrganization{
							PartyId: "",
							NodeName: identity.GetNodeName(),
							NodeId: identity.GetNodeId(),
							IdentityId: identity.GetIdentityId(),
						},
						apipb.TaskRole_TaskRole_Sender,
						proposalState.GetTaskSender(),
						t.mustGetProposalTask(proposalState.GetTaskId()).Task,
						types.TaskConsensusInterrupt,
						nil,
						nil,
					))
				}
			}
			// finally clean all things about proposal and task
			t.removeProposalStateAndTask(proposalState.GetProposalId())
		}
	}

	t.state.proposalsLock.Unlock()
}

func (t *TwoPC) handleInvalidProposal(proposalState *ctypes.ProposalState) {


}


func (t *TwoPC) driveTask(
	pid peer.ID,
	proposalId common.Hash,
	localTaskRole apipb.TaskRole,
	localTaskOrganization *apipb.TaskOrganization,
	remoteTaskRole apipb.TaskRole,
	remoteTaskOrganization *apipb.TaskOrganization,
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

	t.removeProposalStateAndTask(proposalId)

}

func (t *TwoPC) sendPrepareMsg(proposalId common.Hash, task *types.Task, startTime uint64) error {

	sender := task.GetTaskSender()

	sendTaskFn := func(wg *sync.WaitGroup, sender, receiver *apipb.TaskOrganization, senderRole, receiverRole apipb.TaskRole, errCh chan<- error) {

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
		go sendTaskFn(&wg, sender, receiver, apipb.TaskRole_TaskRole_Sender, apipb.TaskRole_TaskRole_DataSupplier, errCh)

	}
	for i := 0; i < len(task.GetTaskData().GetPowerSuppliers()); i++ {

		wg.Add(1)
		powerSupplier := task.GetTaskData().GetPowerSuppliers()[i]
		receiver := powerSupplier.GetOrganization()
		go sendTaskFn(&wg, sender, receiver, apipb.TaskRole_TaskRole_Sender, apipb.TaskRole_TaskRole_PowerSupplier, errCh)

	}

	for i := 0; i < len(task.GetTaskData().GetReceivers()); i++ {

		wg.Add(1)
		receiver := task.GetTaskData().GetReceivers()[i]
		go sendTaskFn(&wg, sender, receiver, apipb.TaskRole_TaskRole_Sender, apipb.TaskRole_TaskRole_Receiver, errCh)
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

func (t *TwoPC) sendPrepareVote(pid peer.ID, sender, receiver *apipb.TaskOrganization, req *pb.PrepareVote) error {
	if types.IsNotSameTaskOrg(sender, receiver) {
		return handler.SendTwoPcPrepareVote(context.TODO(), t.p2p, pid, req)
	} else {
		return t.sendLocalPrepareVote(pid, req)
	}
}

func (t *TwoPC) sendConfirmMsg(proposalId common.Hash, task *types.Task, peers *pb.ConfirmTaskPeerInfo, startTime uint64) error {

	sender := task.GetTaskSender()

	sendConfirmMsgFn := func(wg *sync.WaitGroup, sender, receiver *apipb.TaskOrganization, senderRole, receiverRole apipb.TaskRole, errCh chan<- error) {

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
		go sendConfirmMsgFn(&wg, sender, receiver, apipb.TaskRole_TaskRole_Sender, apipb.TaskRole_TaskRole_DataSupplier, errCh)

	}
	for i := 0; i < len(task.GetTaskData().GetPowerSuppliers()); i++ {

		wg.Add(1)
		powerSupplier := task.GetTaskData().GetPowerSuppliers()[i]
		receiver := powerSupplier.GetOrganization()
		go sendConfirmMsgFn(&wg, sender, receiver, apipb.TaskRole_TaskRole_Sender, apipb.TaskRole_TaskRole_PowerSupplier, errCh)

	}

	for i := 0; i < len(task.GetTaskData().GetReceivers()); i++ {

		wg.Add(1)
		receiver := task.GetTaskData().GetReceivers()[i]
		go sendConfirmMsgFn(&wg, sender, receiver, apipb.TaskRole_TaskRole_Sender, apipb.TaskRole_TaskRole_Receiver, errCh)
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

func (t *TwoPC) sendConfirmVote(pid peer.ID, sender, receiver *apipb.TaskOrganization, req *pb.ConfirmVote) error {
	if types.IsNotSameTaskOrg(sender, receiver) {
		return handler.SendTwoPcConfirmVote(context.TODO(), t.p2p, pid, req)
	} else {
		return t.sendLocalConfirmVote(pid, req)
	}
}

func (t *TwoPC) sendCommitMsg(proposalId common.Hash, task *types.Task, startTime uint64) error {

	sender := task.GetTaskSender()

	sendCommitMsgFn := func(wg *sync.WaitGroup, sender, receiver *apipb.TaskOrganization, senderRole, receiverRole apipb.TaskRole, errCh chan<- error) {

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
		go sendCommitMsgFn(&wg, sender, receiver, apipb.TaskRole_TaskRole_Sender, apipb.TaskRole_TaskRole_DataSupplier, errCh)

	}
	for i := 0; i < len(task.GetTaskData().GetPowerSuppliers()); i++ {

		wg.Add(1)
		powerSupplier := task.GetTaskData().GetPowerSuppliers()[i]
		receiver := powerSupplier.GetOrganization()
		go sendCommitMsgFn(&wg, sender, receiver, apipb.TaskRole_TaskRole_Sender, apipb.TaskRole_TaskRole_PowerSupplier, errCh)

	}

	for i := 0; i < len(task.GetTaskData().GetReceivers()); i++ {

		wg.Add(1)
		receiver := task.GetTaskData().GetReceivers()[i]
		go sendCommitMsgFn(&wg, sender, receiver, apipb.TaskRole_TaskRole_Sender, apipb.TaskRole_TaskRole_Receiver, errCh)
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

func verifyPartyRole(partyId string, role apipb.TaskRole, task *types.Task) bool {

	switch role {
	case apipb.TaskRole_TaskRole_Sender:
		if partyId == task.GetTaskSender().GetPartyId() {
			return true
		}
	case apipb.TaskRole_TaskRole_DataSupplier:
		for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
			if partyId == dataSupplier.GetOrganization().GetPartyId() {
				return true
			}
		}
	case apipb.TaskRole_TaskRole_PowerSupplier:
		for _, powerSupplier := range task.GetTaskData().GetPowerSuppliers() {
			if partyId == powerSupplier.GetOrganization().GetPartyId() {
				return true
			}
		}
	case apipb.TaskRole_TaskRole_Receiver:
		for _, receiver := range task.GetTaskData().GetReceivers() {
			if partyId == receiver.GetPartyId() {
				return true
			}
		}
	}
	return false
}

func fetchOrgByPartyRole(partyId string, role apipb.TaskRole, task *types.Task) *apipb.TaskOrganization {

	switch role {
	case apipb.TaskRole_TaskRole_Sender:
		if partyId == task.GetTaskSender().GetPartyId() {
			return task.GetTaskSender()
		}
	case apipb.TaskRole_TaskRole_DataSupplier:
		for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
			if partyId == dataSupplier.GetOrganization().GetPartyId() {
				return dataSupplier.GetOrganization()
			}
		}
	case apipb.TaskRole_TaskRole_PowerSupplier:
		for _, powerSupplier := range task.GetTaskData().GetPowerSuppliers() {
			if partyId == powerSupplier.GetOrganization().GetPartyId() {
				return powerSupplier.GetOrganization()
			}
		}
	case apipb.TaskRole_TaskRole_Receiver:
		for _, receiver := range task.GetTaskData().GetReceivers() {
			if partyId == receiver.GetPartyId() {
				return receiver
			}
		}
	}
	return nil
}

func (t *TwoPC) isIncludeSenderParty(partyId, identityId string, role apipb.TaskRole, task *types.Task) bool {
	if t.getIncludeSenderPartyCount(partyId, identityId, role, task) != 0 {
		return true
	}
	return false
}

func (t *TwoPC) getIncludeSenderPartyCount(partyId, identityId string, role apipb.TaskRole, task *types.Task) uint32 {

	var count int
	switch role {
	case apipb.TaskRole_TaskRole_DataSupplier:
		for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
			if partyId == dataSupplier.GetOrganization().GetPartyId() && identityId == dataSupplier.GetOrganization().GetIdentityId() {
				count++
			}
		}
	case apipb.TaskRole_TaskRole_PowerSupplier:
		for _, powerSupplier := range task.GetTaskData().GetPowerSuppliers() {
			if partyId == powerSupplier.GetOrganization().GetPartyId() && identityId == powerSupplier.GetOrganization().GetIdentityId() {
				count++
			}
		}
	case apipb.TaskRole_TaskRole_Receiver:
		for _, receiver := range task.GetTaskData().GetReceivers() {
			if partyId == receiver.GetPartyId() && identityId == receiver.GetIdentityId() {
				count++
			}
		}
	}
	return uint32(count)
}
func (t *TwoPC) getNeedVotingCount(role apipb.TaskRole, task *types.Task) uint32 {
	sender := task.GetTaskSender()
	includeCount := t.getIncludeSenderPartyCount(sender.GetPartyId(), sender.GetIdentityId(), role, task)
	switch role {
	case apipb.TaskRole_TaskRole_DataSupplier:
		return uint32(len(task.GetTaskData().GetDataSuppliers())) - includeCount
	case apipb.TaskRole_TaskRole_PowerSupplier:
		return uint32(len(task.GetTaskData().GetPowerSuppliers())) - includeCount
	case apipb.TaskRole_TaskRole_Receiver:
		return uint32(len(task.GetTaskData().GetReceivers())) - includeCount
	}
	return 0
}

func (t *TwoPC) verifyPrepareVoteRole(proposalId common.Hash, partyId, identityId string, role apipb.TaskRole, task *types.Task) (bool, error) {
	var identityValid bool
	switch role {
	case apipb.TaskRole_TaskRole_DataSupplier:

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
	case apipb.TaskRole_TaskRole_PowerSupplier:

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
	case apipb.TaskRole_TaskRole_Receiver:

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
func (t *TwoPC) verifyConfirmVoteRole(proposalId common.Hash, partyId, identityId string, role apipb.TaskRole, task *types.Task) (bool, error) {
	var identityValid bool
	switch role {
	case apipb.TaskRole_TaskRole_DataSupplier:

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
	case apipb.TaskRole_TaskRole_PowerSupplier:

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
	case apipb.TaskRole_TaskRole_Receiver:

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

func (t *TwoPC) getPrepareVote(proposalId common.Hash, partyId string) *types.PrepareVote {
	return t.state.GetPrepareVote(proposalId, partyId)
}

func (t *TwoPC) getConfirmVote(proposalId common.Hash, partyId string) *types.ConfirmVote {
	return t.state.GetConfirmVote(proposalId, partyId)
}

func (t *TwoPC) storeConfirmTaskPeerInfo(proposalId common.Hash, peers *pb.ConfirmTaskPeerInfo) {
	t.state.StoreConfirmTaskPeerInfo(proposalId, peers)
}

func (t *TwoPC) hasConfirmTaskPeerInfo(proposalId common.Hash) bool {
	return t.state.HasConfirmTaskPeerInfo(proposalId)
}

func (t *TwoPC) getConfirmTaskPeerInfo(proposalId common.Hash) (*pb.ConfirmTaskPeerInfo, bool) {
	return t.state.GetConfirmTaskPeerInfo(proposalId)
}

func (t *TwoPC) mustGetConfirmTaskPeerInfo(proposalId common.Hash) *pb.ConfirmTaskPeerInfo {
	return t.state.MustGetConfirmTaskPeerInfo(proposalId)
}

func (t *TwoPC) removeConfirmTaskPeerInfo(proposalId common.Hash) {
	t.state.RemoveConfirmTaskPeerInfo(proposalId)
}
