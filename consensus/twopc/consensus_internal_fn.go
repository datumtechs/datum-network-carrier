package twopc

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	"github.com/RosettaFlow/Carrier-Go/handler"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/libp2p/go-libp2p-core/peer"
	"strings"
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

func (t *TwoPC) addSendTask(task *types.Task) error {
	t.sendTasks[task.TaskId()] = task
	if err := t.dataCenter.StoreLocalTask(task); nil != err {
		delete(t.sendTasks, task.TaskId())
		return err
	}
	return nil
}
func (t *TwoPC) addRecvTask(task *types.Task) error {
	t.recvTasks[task.TaskId()] = task
	if err := t.dataCenter.StoreLocalTask(task); nil != err {
		delete(t.recvTasks, task.TaskId())
		return err
	}
	return nil
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
func (t *TwoPC) collectTaskResult(result *types.ConsensuResult) {
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

func (t *TwoPC) sendReplaySchedTask(replaySchedTask *types.ReplayScheduleTaskWrap) {
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
	if proposalState.TaskDir == types.SendTaskDir {
		// 发布 task 和  event 给 dataCenter
		t.pulishFinishedTaskToDataCenter(proposalState.TaskId)
	} else {
		// 给 task  owner 发出 taskResultMsg  TODO 先不做处理 ...

	}
	// 清空本地 资源占用 和 各种缓存...
	t.delProposalStateAndTask(proposalState.ProposalId)
	t.resourceMng.UnLockLocalResourceWithTask(proposalState.TaskId)
}

func (t *TwoPC) pulishFinishedTaskToDataCenter(taskId string) {
	eventList, err := t.dataCenter.GetTaskEventList(taskId)
	if nil != err {
		log.Errorf("Failed to Query local task event list for sending datacenter, taskId {%s}, err {%s}", taskId, err)
		return
	}
	task, err := t.dataCenter.GetLocalTask(taskId)
	if nil != err {
		log.Errorf("Failed to Query local task info for sending datacenter, taskId {%s}, err {%s}", taskId, err)
		return
	}
	task.SetEventList(eventList)
	if err := t.dataCenter.InsertTask(task); nil != err {
		log.Error("Failed to pulish task and eventlist to datacenter, taskId {%s}, err {%s}", taskId, err)
		return
	}
}

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
	task *types.Task,
	) {


	//var recvTaskResultCh chan *types.TaskResultMsgWrap

	recvTaskResultCh := make(chan *types.TaskResultMsgWrap, 0)

	confirmTaskPeerInfo := t.state.GetConfirmTaskPeerInfo(proposalId)
	if nil == confirmTaskPeerInfo {
		log.Error("Failed to find local cache about all peer resource {externalIP:externalPORT}")
		return
	}
	// Send task to TaskManager to execute
	taskWrap := &types.DoneScheduleTaskChWrap{
		ProposalId: proposalId,
		SelfTaskRole: taskRole,
		// SelfPeerInfo: // TODO
		Task: &types.ConsensusScheduleTask{
			TaskDir:   taskDir,
			TaskState: taskState,
			SchedTask: task,
			Resources: confirmTaskPeerInfo,
		},
		ResultCh: recvTaskResultCh,
	}
	// 发给 taskManager 去执行 task
	t.sendTaskToTaskManagerForExecute(taskWrap)
	go func() {
		if taskDir == types.RecvTaskDir {
			if taskResultWrap, ok := <-taskWrap.ResultCh; ok {
				if err := t.sendTaskResultMsg(pid, taskResultWrap); nil != err {
					log.Error(err)
				}
				// clean local proposalState and task cache
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

	sendTaskFn := func(proposalId common.Hash, taskRole types.TaskRole, identityId, nodeId, taskId string, prepareMsg *pb.PrepareMsg, errCh chan<- error) {
		var pid, err = p2p.HexPeerID(nodeId)
		if nil != err {
			errCh <- fmt.Errorf("failed to nodeId => peerId, proposalId: %s, taskId: %s, taskRole: %s, identityId: %s, nodeId: %s, err: %s",
				proposalId.String(), taskId, taskRole.String(), identityId, nodeId, err)
			return
		}
		prepareMsg.TaskRole = taskRole.Bytes()
		if err = handler.SendTwoPcPrepareMsg(context.TODO(), t.p2p, pid, prepareMsg); nil != err {
			errCh <- fmt.Errorf("failed to call `SendTwoPcPrepareMsg` proposalId: %s, taskId: %s, taskRole: %s, identityId: %s, nodeId: %s, err: %s",
				proposalId.String(), taskId, taskRole.String(), identityId, nodeId, err)
			return
		}
		errCh <- nil
	}

	errCh := make(chan error, (len(task.TaskData().MetadataSupplier) - 1) + len(task.TaskData().ResourceSupplier)+len(task.TaskData().Receivers))

	go func() {
		for _, partner := range task.TaskData().MetadataSupplier {
			// 排除掉 task 发起方 ...
			if task.TaskData().Identity != partner.Organization.Identity {
				go sendTaskFn(proposalId, types.DataSupplier, partner.Organization.Identity, partner.Organization.NodeId, task.TaskId(), prepareMsg, errCh)
			}
		}
	}()
	go func() {
		for _, powerSupplier := range task.TaskData().ResourceSupplier {
			go sendTaskFn(proposalId, types.PowerSupplier, powerSupplier.Organization.Identity, powerSupplier.Organization.NodeId, task.TaskId(), prepareMsg, errCh)
		}
	}()

	go func() {
		for _, receiver := range task.TaskData().Receivers {
			go sendTaskFn(proposalId, types.ResultSupplier, receiver.Receiver.Identity, receiver.Receiver.NodeId, task.TaskId(), prepareMsg, errCh)
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

	sendConfirmMsgFn := func(proposalId common.Hash, taskRole types.TaskRole, identityId, nodeId, taskId string, confirmMsg *pb.ConfirmMsg, errCh chan<- error) {
		pid, err := p2p.HexPeerID(nodeId)
		if nil != err {
			errCh <- fmt.Errorf("failed to nodeId => peerId, proposalId: %s, taskId: %s, taskRole: %s, identityId: %s, nodeId: %s, err: %s",
				proposalId.String(), taskId, taskRole.String(), identityId, nodeId, err)
			return
		}

		confirmMsg.TaskRole = taskRole.Bytes()
		// Send the ConfirmMsg to other peer
		if err := handler.SendTwoPcConfirmMsg(context.TODO(), t.p2p, pid, confirmMsg); nil != err {
			errCh <- fmt.Errorf("failed to call`SendTwoPcConfirmMsg` proposalId: %s, taskId: %s, taskRole: %s, identityId: %s, nodeId: %s, err: %s",
				proposalId.String(), taskId, taskRole.String(), identityId, nodeId, err)
			errCh <- err
			return
		}

		errCh <- nil
	}



	errCh := make(chan error, (len(task.TaskData().MetadataSupplier) - 1) + len(task.TaskData().ResourceSupplier)+len(task.TaskData().Receivers))

	go func() {
		for _, partner := range task.TaskData().MetadataSupplier {
			// 排除掉 task 发起方 ...
			if task.TaskData().Identity != partner.Organization.Identity {
				go sendConfirmMsgFn(proposalId, types.DataSupplier, partner.Organization.Identity, partner.Organization.NodeId, task.TaskId(), confirmMsg, errCh)
			}
		}
	}()
	go func() {
		for _, powerSupplier := range task.TaskData().ResourceSupplier {
			go sendConfirmMsgFn(proposalId, types.PowerSupplier, powerSupplier.Organization.Identity, powerSupplier.Organization.NodeId, task.TaskId(), confirmMsg, errCh)
		}
	}()

	go func() {
		for _, receiver := range task.TaskData().Receivers {
			go sendConfirmMsgFn(proposalId, types.ResultSupplier, receiver.Receiver.Identity, receiver.Receiver.NodeId, task.TaskId(), confirmMsg, errCh)
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

	sendCommitMsgFn := func(proposalId common.Hash, taskRole types.TaskRole, identityId, nodeId, taskId string, commitMsg *pb.CommitMsg, errCh chan<- error) {
		pid, err := p2p.HexPeerID(nodeId)
		if nil != err {
			errCh <- fmt.Errorf("failed to nodeId => peerId, proposalId: %s, taskId: %s, taskRole: %s, identityId: %s, nodeId: %s, err: %s",
				proposalId.String(), taskId, taskRole.String(), identityId, nodeId, err)
			return
		}

		commitMsg.TaskRole = taskRole.Bytes()
		// Send the ConfirmMsg to other peer
		if err := handler.SendTwoPcCommitMsg(context.TODO(), t.p2p, pid, commitMsg); nil != err {
			errCh <- fmt.Errorf("failed to call`SendTwoPcCommitMsg` proposalId: %s, taskId: %s, taskRole: %s, identityId: %s, nodeId: %s, err: %s",
				proposalId.String(), taskId, taskRole.String(), identityId, nodeId, err)
			errCh <- err
			return
		}

		errCh <- nil
	}


	errCh := make(chan error, (len(task.TaskData().MetadataSupplier) - 1) + len(task.TaskData().ResourceSupplier)+len(task.TaskData().Receivers))

	go func() {
		for _, partner := range task.TaskData().MetadataSupplier {
			// 排除掉 task 发起方 ...
			if task.TaskData().Identity != partner.Organization.Identity {
				go sendCommitMsgFn(proposalId, types.DataSupplier, partner.Organization.Identity, partner.Organization.NodeId, task.TaskId(), commitMsg, errCh)
			}
		}
	}()
	go func() {
		for _, powerSupplier := range task.TaskData().ResourceSupplier {
			go sendCommitMsgFn(proposalId, types.PowerSupplier, powerSupplier.Organization.Identity, powerSupplier.Organization.NodeId, task.TaskId(), commitMsg, errCh)
		}
	}()

	go func() {
		for _, receiver := range task.TaskData().Receivers {
			go sendCommitMsgFn(proposalId, types.ResultSupplier, receiver.Receiver.Identity, receiver.Receiver.NodeId, task.TaskId(), commitMsg, errCh)
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
