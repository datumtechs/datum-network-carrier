package twopc

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	"github.com/RosettaFlow/Carrier-Go/core/evengine"
	"github.com/RosettaFlow/Carrier-Go/handler"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/RosettaFlow/Carrier-Go/types"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/libp2p/go-libp2p-core/peer"
	"strings"
	"time"
)

func (t *TwoPC) isProcessingTask(taskId string) bool {
	if _, ok := t.sendTasks[taskId]; ok {
		return true
	}
	if _, ok := t.recvTasks[taskId]; ok {
		return true
	}
	return false
}

//func (t *TwoPC) cleanExpireProposal() {
//	expireProposalSendTaskIds, expireProposalRecvTaskIds := t.state.CleanExpireProposal()
//	for _, sendTaskId := range expireProposalSendTaskIds {
//		delete(t.sendTasks, sendTaskId)
//	}
//	for _, recvTaskId := range expireProposalRecvTaskIds {
//		delete(t.recvTasks, recvTaskId)
//	}
//}

func (t *TwoPC) addSendTask(task *types.Task)  {
	t.sendTasks[task.TaskId()] = task
}
func (t *TwoPC) addRecvTask(task *types.Task) {
	t.recvTasks[task.TaskId()] = task
}
func (t *TwoPC) delTask(taskId string) {
	t.delSendTask(taskId)
	t.delRecvTask(taskId)
}
func (t *TwoPC) delSendTask(taskId string) {
	delete(t.sendTasks, taskId)
}
func (t *TwoPC) delRecvTask(taskId string) {
	delete(t.recvTasks, taskId)
}
func (t *TwoPC) addTaskResultCh(taskId string, resultCh chan<- *types.ConsensuResult) {
	t.taskResultLock.Lock()
	t.taskResultChs[taskId] = resultCh
	t.taskResultLock.Unlock()
}
func (t *TwoPC) removeTaskResultCh(taskId string) {
	t.taskResultLock.Lock()
	delete(t.taskResultChs, taskId)
	t.taskResultLock.Unlock()
}
func (t *TwoPC) collectTaskResultWillSendToSched(result *types.ConsensuResult) {
	t.taskResultCh <- result
}
func (t *TwoPC) sendConsensusTaskResultToSched (result *types.ConsensuResult) {
	t.taskResultLock.Lock()
	if ch, ok := t.taskResultChs[result.TaskId]; ok {
		ch <- result
		close(ch)
		delete(t.taskResultChs, result.TaskId)
	}
	t.taskResultLock.Unlock()
}

func (t *TwoPC) sendReplaySchedTaskToScheduler(replaySchedTask *types.ReplayScheduleTaskWrap) {
	t.replayTaskCh <- replaySchedTask
}

func (t *TwoPC) addProposalState(proposalState *ctypes.ProposalState) {
	t.state.AddProposalState(proposalState)
}
func (t *TwoPC) delProposalState(proposalId common.Hash) {
	t.state.CleanProposalState(proposalId)
}
func (t *TwoPC) delProposalStateAndTask(proposalId common.Hash) {
	if state := t.state.GetProposalState(proposalId); t.state.EmptyInfo() != state {
		log.Infof("Start remove proposalState and task cache on Consensus, proposalId {%s}, taskId {%s}", proposalId, state.TaskId)
		t.state.CleanProposalState(proposalId)
		t.delTask(state.TaskId)
	}
}

func (t *TwoPC) sendTaskToTaskManagerForExecute(task *types.DoneScheduleTaskChWrap) {
	t.doneScheduleTaskCh <- task
}

func (t *TwoPC) makeConfirmTaskPeerDesc(proposalId common.Hash) *pb.ConfirmTaskPeerInfo {

	dataSuppliers, powerSuppliers, receivers := make([]*pb.TaskPeerInfo, 0), make([]*pb.TaskPeerInfo, 0), make([]*pb.TaskPeerInfo, 0)

	for _, vote := range t.state.GetPrepareVoteArr(proposalId) {
		if vote.TaskRole == types.DataSupplier && nil != vote.PeerInfo {
			dataSuppliers = append(dataSuppliers, types.ConvertTaskPeerInfo(vote.PeerInfo))
		}
		if vote.TaskRole == types.PowerSupplier && nil != vote.PeerInfo {
			powerSuppliers = append(powerSuppliers, types.ConvertTaskPeerInfo(vote.PeerInfo))
		}
		if vote.TaskRole == types.ResultSupplier && nil != vote.PeerInfo {
			receivers = append(receivers, types.ConvertTaskPeerInfo(vote.PeerInfo))
		}
	}
	owner := t.state.GetSelfPeerInfo(proposalId)
	if nil == owner {
		return nil
	}
	return &pb.ConfirmTaskPeerInfo{
		OwnerPeerInfo:              types.ConvertTaskPeerInfo(owner),
		DataSupplierPeerInfoList:   dataSuppliers,
		PowerSupplierPeerInfoList:  powerSuppliers,
		ResultReceiverPeerInfoList: receivers,
	}
}

func (t *TwoPC) refreshProposalState() {


	for id, proposalState := range t.state.GetProposalStates() {

		if proposalState.IsDeadline() {
			t.handleInvalidProposal(proposalState)
			continue
		}

		switch proposalState.GetPeriod() {
		case ctypes.PeriodPrepare:
			if proposalState.IsPrepareTimeout() {
				proposalState.ChangeToConfirm(proposalState.PeriodStartTime + uint64(ctypes.PrepareMsgVotingTimeout.Nanoseconds()))
				t.state.UpdateProposalState(proposalState)
			}
		case ctypes.PeriodConfirm:
			if proposalState.IsConfirmTimeout() {
				proposalState.ChangeToCommit(proposalState.PeriodStartTime + uint64(ctypes.ConfirmMsgVotingTimeout.Nanoseconds()))
				t.state.UpdateProposalState(proposalState)
			}
		case ctypes.PeriodCommit:
			if proposalState.IsCommitTimeout() {
				proposalState.ChangeToFinished(proposalState.PeriodStartTime + uint64(ctypes.CommitMsgEndingTimeout.Nanoseconds()))
				t.state.UpdateProposalState(proposalState)
			}
		case ctypes.PeriodFinished:
			//
			if proposalState.IsDeadline() {
				t.handleInvalidProposal(proposalState)
			}

		default:
			log.Error("Unknown the proposalState period", "proposalId", id.String())
			t.handleInvalidProposal(proposalState)
		}
	}
}

func (t *TwoPC) handleInvalidProposal(proposalState *ctypes.ProposalState) {

	log.Debug("Call handleInvalidProposal(), handle and clean proposalState and task", "proposalId", proposalState.ProposalId, "taskId", proposalState.TaskId, "taskDir", proposalState.TaskDir.String())

	if proposalState.TaskDir == types.SendTaskDir {
		// Send consensus result to Scheduler
		t.collectTaskResultWillSendToSched(&types.ConsensuResult{
			TaskConsResult: &types.TaskConsResult{
				TaskId: proposalState.TaskId,
				Status: types.TaskConsensusInterrupt,
				Done:   false,
				Err:    fmt.Errorf("the task proposalState coming deadline"),
			},
		})
		// clean some invalid data
		t.delProposalStateAndTask(proposalState.ProposalId)
	} else {

		task, ok := t.recvTasks[proposalState.TaskId]
		if !ok {
			log.Errorf("Failed to query recvTaskInfo on consensus.handleInvalidProposal(), taskId: {%s}", proposalState.TaskId)
			return
		}
		eventList, err := t.dataCenter.GetTaskEventList(proposalState.TaskId)
		if nil != err {
			log.Errorf("Failed to GetTaskEventList() on consensus.handleInvalidProposal(), taskId: {%s}, err: {%s}", proposalState.TaskId, err)
			eventList = make([]*types.TaskEventInfo, 0)
		}
		eventList = append(eventList, &types.TaskEventInfo{
			Type: evengine.TaskProposalStateDeadline.Type,
			Identity: proposalState.SelfIdentity.IdentityId,
			TaskId: proposalState.TaskId,
			Content: evengine.TaskProposalStateDeadline.Msg,
			CreateTime: uint64(time.Now().UnixNano()),
		})
		taskResultWrap := &types.TaskResultMsgWrap{
			TaskResultMsg: &pb.TaskResultMsg{
				ProposalId: proposalState.ProposalId.Bytes(),
				TaskRole: proposalState.TaskRole.Bytes(),
				TaskId: []byte(proposalState.TaskId),
				Owner: &pb.TaskOrganizationIdentityInfo{
					PartyId: []byte(proposalState.SelfIdentity.PartyId),
					Name: []byte(proposalState.SelfIdentity.Name),
					NodeId: []byte(proposalState.SelfIdentity.NodeId),
					IdentityId: []byte(proposalState.SelfIdentity.IdentityId),
				},
				TaskEventList: types.ConvertTaskEventArr(eventList),
				CreateAt: uint64(time.Now().UnixNano()),
				Sign: nil,
			},
		}

		// Send taskResultMsg to task Owner
		pid, err := p2p.HexPeerID(task.TaskData().NodeId)
		if nil == err {
			if err := t.sendTaskResultMsg(pid, taskResultWrap); nil != err {
				log.Error(err)
			}
		}

		t.resourceMng.UnLockLocalResourceWithTask(proposalState.TaskId)
		// 因为在 scheduler 那边已经对 task 做了 StoreLocalTask
		t.dataCenter.RemoveLocalTask(proposalState.TaskId)
		t.dataCenter.CleanTaskEventList(proposalState.TaskId)
		// clean some data
		t.delProposalStateAndTask(proposalState.ProposalId)
	}
}
//
//func (t *TwoPC) pulishFinishedTaskToDataCenter(taskId string) {
//	eventList, err := t.dataCenter.GetTaskEventList(taskId)
//	if nil != err {
//		log.Errorf("Failed to Query local task event list for sending datacenter, taskId {%s}, err {%s}", taskId, err)
//		return
//	}
//	task, err := t.dataCenter.GetLocalTask(taskId)
//	if nil != err {
//		log.Errorf("Failed to Query local task info for sending datacenter, taskId {%s}, err {%s}", taskId, err)
//		return
//	}
//	task.SetEventList(eventList)
//	task.TaskData().EndAt = uint64(time.Now().UnixNano())
//	//task.TaskData().State =
//	if err := t.dataCenter.InsertTask(task); nil != err {
//		log.Error("Failed to pulish task and eventlist to datacenter, taskId {%s}, err {%s}", taskId, err)
//		return
//	}
//}

func (t *TwoPC) storeTaskEvent(pid peer.ID, taskId string, events []*types.TaskEventInfo) error {
	for _, event := range events {
		if err := t.dataCenter.StoreTaskEvent(event); nil != err {
			log.Error("Failed to store local task event from remote peer", "remote peerId", pid, "taskId", taskId)
		}
	}
	return nil
}

func (t *TwoPC) driveTask(
	pid peer.ID,
	proposalId common.Hash,
	taskDir types.ProposalTaskDir,
	taskState types.TaskState,
	taskRole  types.TaskRole,
	selfIdentity *libTypes.OrganizationData,
	task *types.Task,
	) {



	confirmTaskPeerInfo := t.state.GetConfirmTaskPeerInfo(proposalId)
	if nil == confirmTaskPeerInfo {
		log.Error("Failed to find local cache about all peer resource {externalIP:externalPORT}")
		return
	}
	// Send task to TaskManager to execute
	taskWrap := &types.DoneScheduleTaskChWrap{
		ProposalId: proposalId,
		SelfTaskRole: taskRole,
		SelfIdentity: selfIdentity,
		Task: &types.ConsensusScheduleTask{
			TaskDir:   taskDir,
			TaskState: taskState,
			SchedTask: task,
			Resources: confirmTaskPeerInfo,
		},
		ResultCh: make(chan *types.TaskResultMsgWrap, 0),
	}
	// 发给 taskManager 去执行 task
	t.sendTaskToTaskManagerForExecute(taskWrap)
	go func() {
		if taskDir == types.RecvTaskDir {
			if taskResultWrap, ok := <-taskWrap.ResultCh; ok {
				if err := t.sendTaskResultMsg(pid, taskResultWrap); nil != err {
					log.Error(err)
				}

				t.resourceMng.UnLockLocalResourceWithTask(task.TaskId())
				// 因为在 scheduler 那边已经对 task 做了 StoreLocalTask
				t.dataCenter.RemoveLocalTask(task.TaskId())
				t.dataCenter.CleanTaskEventList(task.TaskId())
				// clean some data
				t.delProposalStateAndTask(proposalId)
			}
		} else {
			<-taskWrap.ResultCh  // publish taskInfo to dataCenter done ..
			// clean local proposalState and task cache
			t.delProposalStateAndTask(proposalId)
		}
	}()
}

func (t *TwoPC) sendPrepareMsg(proposalId common.Hash, task *types.Task, startTime uint64) error {

	prepareMsg := makePrepareMsgWithoutTaskRole(proposalId, task, startTime)

	sendTaskFn := func(proposalId common.Hash, taskRole types.TaskRole, partyId, identityId, nodeId, taskId string, prepareMsg *pb.PrepareMsg, errCh chan<- error) {
		var pid, err = p2p.HexPeerID(nodeId)
		if nil != err {
			errCh <- fmt.Errorf("failed to nodeId => peerId, proposalId: %s, taskId: %s, other peer taskRole: %s, other peer taskPartyId: %s, identityId: %s, nodeId: %s, err: %s",
				proposalId.String(), taskId, taskRole.String(), partyId, identityId, nodeId, err)
			return
		}
		// set other peer's role and partyId
		prepareMsg.TaskRole = taskRole.Bytes()
		prepareMsg.TaskPartyId = []byte(partyId)

		if err = handler.SendTwoPcPrepareMsg(context.TODO(), t.p2p, pid, prepareMsg); nil != err {
			errCh <- fmt.Errorf("failed to call `SendTwoPcPrepareMsg` proposalId: %s, taskId: %s, other peer taskRole: %s, other peer taskPartyId: %s, identityId: %s, nodeId: %s, err: %s",
				proposalId.String(), taskId, taskRole.String(), partyId, identityId, nodeId, err)
			return
		}
		errCh <- nil
	}

	errCh := make(chan error, (len(task.TaskData().MetadataSupplier) - 1) + len(task.TaskData().ResourceSupplier)+len(task.TaskData().Receivers))

	go func() {
		for _, dataSupplier := range task.TaskData().MetadataSupplier {
			// 排除掉 task 发起方 ...
			if task.TaskData().Identity != dataSupplier.Organization.Identity && task.TaskData().PartyId != dataSupplier.Organization.PartyId {
				go sendTaskFn(proposalId, types.DataSupplier, dataSupplier.Organization.PartyId,
					dataSupplier.Organization.Identity, dataSupplier.Organization.NodeId, task.TaskId(), prepareMsg, errCh)
			}
		}
	}()
	go func() {
		for _, powerSupplier := range task.TaskData().ResourceSupplier {
			go sendTaskFn(proposalId, types.PowerSupplier, powerSupplier.Organization.PartyId,
				powerSupplier.Organization.Identity, powerSupplier.Organization.NodeId, task.TaskId(), prepareMsg, errCh)
		}
	}()

	go func() {
		for _, receiver := range task.TaskData().Receivers {
			go sendTaskFn(proposalId, types.ResultSupplier, receiver.Receiver.PartyId,
				receiver.Receiver.Identity, receiver.Receiver.NodeId, task.TaskId(), prepareMsg, errCh)
		}
	}()
	errStrs := make([]string, 0)
	for err := range errCh {
		if nil != err {
			errStrs = append(errStrs, err.Error())
		}
	}
	if len(errStrs) != 0 {
		return fmt.Errorf(
			`failed to Send PrepareMsg for task:
%s`, strings.Join(errStrs, "\n"))
	}

	return nil
}

func (t *TwoPC) sendConfirmMsg(proposalId common.Hash, task *types.Task, startTime uint64) error {

	confirmMsg := makeConfirmMsg(proposalId, task, startTime)
	confirmMsg.PeerDesc = t.makeConfirmTaskPeerDesc(proposalId)

	// store the proposal about all partner peerInfo of task to local cache
	t.state.StoreConfirmTaskPeerInfo(proposalId, confirmMsg.PeerDesc)

	sendConfirmMsgFn := func(proposalId common.Hash, taskRole types.TaskRole, taskPartyId, identityId, nodeId, taskId string, confirmMsg *pb.ConfirmMsg, errCh chan<- error) {
		pid, err := p2p.HexPeerID(nodeId)
		if nil != err {
			errCh <- fmt.Errorf("failed to nodeId => peerId, proposalId: %s, taskId: %s, other peer's taskRole: %s, other peer's partyId: %s, identityId: %s, nodeId: %s, err: %s",
				proposalId.String(), taskId, taskRole.String(), taskPartyId, identityId, nodeId, err)
			return
		}

		// set other peer's role and partyId
		confirmMsg.TaskRole = taskRole.Bytes()
		confirmMsg.TaskPartyId = []byte(taskPartyId)

		// Send the ConfirmMsg to other peer
		if err := handler.SendTwoPcConfirmMsg(context.TODO(), t.p2p, pid, confirmMsg); nil != err {
			errCh <- fmt.Errorf("failed to call`SendTwoPcConfirmMsg` proposalId: %s, taskId: %s,other peer's taskRole: %s, other peer's partyId: %s, identityId: %s, nodeId: %s, err: %s",
				proposalId.String(), taskId, taskRole.String(), taskPartyId, identityId, nodeId, err)
			errCh <- err
			return
		}

		errCh <- nil
	}



	errCh := make(chan error, (len(task.TaskData().MetadataSupplier) - 1) + len(task.TaskData().ResourceSupplier)+len(task.TaskData().Receivers))

	go func() {
		for _, dataSupplier := range task.TaskData().MetadataSupplier {
			// 排除掉 task 发起方 ...  identityId + partyId
			if task.TaskData().Identity != dataSupplier.Organization.Identity && task.TaskData().PartyId != dataSupplier.Organization.PartyId {
				go sendConfirmMsgFn(proposalId, types.DataSupplier, dataSupplier.Organization.PartyId,
					dataSupplier.Organization.Identity, dataSupplier.Organization.NodeId, task.TaskId(), confirmMsg, errCh)
			}
		}
	}()
	go func() {
		for _, powerSupplier := range task.TaskData().ResourceSupplier {
			go sendConfirmMsgFn(proposalId, types.PowerSupplier, powerSupplier.Organization.PartyId,
				powerSupplier.Organization.Identity, powerSupplier.Organization.NodeId, task.TaskId(), confirmMsg, errCh)
		}
	}()

	go func() {
		for _, receiver := range task.TaskData().Receivers {
			go sendConfirmMsgFn(proposalId, types.ResultSupplier, receiver.Receiver.PartyId,
				receiver.Receiver.Identity, receiver.Receiver.NodeId, task.TaskId(), confirmMsg, errCh)
		}
	}()

	errStrs := make([]string, 0)
	for err := range errCh {
		if nil != err {
			errStrs = append(errStrs, err.Error())
		}
	}
	if len(errStrs) != 0 {
		return fmt.Errorf(
			`failed to Send ConfirmMsg for task:
%s`, strings.Join(errStrs, "\n"))
	}

	return nil
}

func (t *TwoPC) sendCommitMsg(proposalId common.Hash, task *types.Task, startTime uint64) error {

	commitMsg := makeCommitMsg(proposalId, task, startTime)

	sendCommitMsgFn := func(proposalId common.Hash, taskRole types.TaskRole, taskPartyId, identityId, nodeId, taskId string, commitMsg *pb.CommitMsg, errCh chan<- error) {
		pid, err := p2p.HexPeerID(nodeId)
		if nil != err {
			errCh <- fmt.Errorf("failed to nodeId => peerId, proposalId: %s, taskId: %s, other peer's taskRole: %s, other peer's partyId: %s, identityId: %s, nodeId: %s, err: %s",
				proposalId.String(), taskId, taskRole.String(), taskPartyId, identityId, nodeId, err)
			return
		}

		// set other peer's role and partyId
		commitMsg.TaskRole = taskRole.Bytes()
		commitMsg.TaskPartyId = []byte(taskPartyId)

		// Send the ConfirmMsg to other peer
		if err := handler.SendTwoPcCommitMsg(context.TODO(), t.p2p, pid, commitMsg); nil != err {
			errCh <- fmt.Errorf("failed to call`SendTwoPcCommitMsg` proposalId: %s, taskId: %s,  other peer's taskRole: %s, other peer's partyId: %s, identityId: %s, nodeId: %s, err: %s",
				proposalId.String(), taskId, taskRole.String(), taskPartyId, identityId, nodeId, err)
			errCh <- err
			return
		}

		errCh <- nil
	}


	errCh := make(chan error, (len(task.TaskData().MetadataSupplier) - 1) + len(task.TaskData().ResourceSupplier)+len(task.TaskData().Receivers))

	go func() {
		for _, dataSupplier := range task.TaskData().MetadataSupplier {
			// 排除掉 task 发起方 ...  identityId + partyId
			if task.TaskData().Identity != dataSupplier.Organization.Identity && task.TaskData().PartyId != dataSupplier.Organization.PartyId {
				go sendCommitMsgFn(proposalId, types.DataSupplier, dataSupplier.Organization.PartyId,
					dataSupplier.Organization.Identity, dataSupplier.Organization.NodeId, task.TaskId(), commitMsg, errCh)
			}
		}
	}()
	go func() {
		for _, powerSupplier := range task.TaskData().ResourceSupplier {
			go sendCommitMsgFn(proposalId, types.PowerSupplier, powerSupplier.Organization.PartyId,
				powerSupplier.Organization.Identity, powerSupplier.Organization.NodeId, task.TaskId(), commitMsg, errCh)
		}
	}()

	go func() {
		for _, receiver := range task.TaskData().Receivers {
			go sendCommitMsgFn(proposalId, types.ResultSupplier, receiver.Receiver.PartyId,
				receiver.Receiver.Identity, receiver.Receiver.NodeId, task.TaskId(), commitMsg, errCh)
		}
	}()

	errStrs := make([]string, 0)
	for err := range errCh {
		if nil != err {
			errStrs = append(errStrs, err.Error())
		}
	}
	if len(errStrs) != 0 {
		return fmt.Errorf(
			`failed to Send CommitMsg for task:
%s`, strings.Join(errStrs, "\n"))
	}

	return nil
}
