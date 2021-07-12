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


func (t *TwoPC) cleanExpireProposal() {
	expireProposalSendTaskIds, expireProposalRecvTaskIds := t.state.CleanExpireProposal()
	for _, sendTaskId := range expireProposalSendTaskIds {
		delete(t.sendTasks, sendTaskId)
	}
	for _, recvTaskId := range expireProposalRecvTaskIds {
		delete(t.recvTasks, recvTaskId)
	}
}


func (t *TwoPC) addSendTask(task *types.ScheduleTask) {
	t.sendTasks[task.TaskId] = task
}
func (t *TwoPC) addRecvTask(task *types.ScheduleTask) {
	t.recvTasks[task.TaskId] = task
}
func (t *TwoPC) delSendTask(taskId string) {
	delete(t.sendTasks, taskId)
}
func (t *TwoPC) delRecvTask(taskId string) {
	delete(t.recvTasks, taskId)
}
func (t *TwoPC) addTaskResultCh (taskId string, resultCh chan<- *types.ConsensuResult) {
	t.taskResultLock.Lock()
	t.taskResultChs[taskId] = resultCh
	t.taskResultLock.Unlock()
}
func (t *TwoPC) removeTaskResultCh (taskId string) {
	t.taskResultLock.Lock()
	delete(t.taskResultChs, taskId)
	t.taskResultLock.Unlock()
}
func (t *TwoPC) collectTaskResult (result *types.ConsensuResult) {
	t.taskResultCh <- result
}
func (t *TwoPC) sendTaskResult (result *types.ConsensuResult) {
	t.taskResultLock.Lock()
	if ch, ok := t.taskResultChs[result.TaskId]; ok {
		ch <- result
		close(ch)
		delete(t.taskResultChs, result.TaskId)
	}
	t.taskResultLock.Unlock()
}

func (t *TwoPC) sendReplaySchedTask (replaySchedTask *types.ScheduleTaskWrap) {
	t.replayTaskCh <- replaySchedTask
}

func (t *TwoPC) addProposalState(proposalState *ctypes.ProposalState) {
	t.state.AddProposalState(proposalState)
}
func (t *TwoPC) delProposalState(proposalId common.Hash) {
	t.state.CleanProposalState(proposalId)
}

func (t *TwoPC) sendTaskForStart (task *types.ConsensusScheduleTask) {
	t.recvSchedTaskCh <- task
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
		OwnerPeerInfo: types.ConvertTaskPeerInfo(owner),
		DataSupplierPeerInfoList: dataSuppliers,
		PowerSupplierPeerInfoList: powerSuppliers,
		ResultReceiverPeerInfoList: receivers,
	}
}

// TODO 写到这里
func (t *TwoPC) refreshProposalState() {
	//for _, proposalState := range t.state.GetProposalStates() {
	//
	//}
}

func (t *TwoPC) storeTaskEvent(pid peer.ID, taskId string, events []*types.TaskEventInfo) error {
	for _, event := range events {
		if err := t.dataCenter.StoreTaskEvent(event); nil != err {
			log.Error("Failed to store local task event from remote peer", "remote peerId", pid, "taskId", taskId)
		}
	}
	return nil
}

func (t *TwoPC) driveTask(taskDir ctypes.ProposalTaskDir, proposalId common.Hash, taskId string) error {

	var task *types.ScheduleTask
	if taskDir == ctypes.SendTaskDir {
		task = t.sendTasks[taskId]
	} else {
		task = t.sendTasks[taskId]
	}
	if nil == task {
		return ctypes.ErrProposalTaskNotFound
	}

	dataSuppliers, powerSuppliers, receivers :=
		make([]*types.PrepareVoteResource, 0), make([]*types.PrepareVoteResource, 0), make([]*types.PrepareVoteResource, 0)

	for _, vote := range t.state.GetPrepareVoteArr(proposalId) {
		if vote.TaskRole == types.DataSupplier && nil != vote.PeerInfo {
			dataSuppliers = append(dataSuppliers, vote.PeerInfo)
		}
		if vote.TaskRole == types.PowerSupplier && nil != vote.PeerInfo {
			powerSuppliers = append(powerSuppliers, vote.PeerInfo)
		}
		if vote.TaskRole == types.ResultSupplier && nil != vote.PeerInfo {
			receivers = append(receivers, vote.PeerInfo)
		}
	}
	//
	t.sendTaskForStart(&types.ConsensusScheduleTask{
		SchedTask: task,
		OwnerResource: t.state.GetSelfPeerInfo(proposalId),
		PartnersResource: dataSuppliers,
		PowerSuppliersResource: powerSuppliers,
		ReceiversResource: receivers,
	})

	return nil
}

func (t *TwoPC) sendPrepareMsg(proposalId common.Hash, task *types.ScheduleTask, startTime uint64) error {

	prepareMsg := makePrepareMsg(startTime)

	sendTaskFn := func(proposalId common.Hash, taskRole types.TaskRole, identityId, nodeId, taskId string, prepareMsg *pb.PrepareMsg, errCh chan<- error) {
		var pid, err = p2p.HexPeerID(nodeId)
		if nil != err {
			errCh <- fmt.Errorf("failed to nodeId => peerId, proposalId: %s, taskId: %s, taskRole: %s, identityId: %s, nodeId: %s, err: %s",
				proposalId.String(), taskId, taskRole.String(), identityId, nodeId, err)
			return
		}

		prepareMsg.TaskOption.TaskRole = taskRole.Bytes()
		if err = handler.SendTwoPcPrepareMsg(context.TODO(), t.p2p, pid, prepareMsg); nil != err {
			errCh <- fmt.Errorf("failed to call `SendTwoPcPrepareMsg` proposalId: %s, taskId: %s, taskRole: %s, identityId: %s, nodeId: %s, err: %s",
				proposalId.String(), taskId, taskRole.String(), identityId, nodeId, err)
			return
		}
		errCh <- nil
	}

	errCh := make(chan error, len(task.Partners)+len(task.PowerSuppliers)+len(task.Receivers))

	go func() {
		for _, partner := range task.Partners {
			go sendTaskFn(proposalId, types.DataSupplier, partner.IdentityId, partner.NodeId, task.TaskId, prepareMsg, errCh)
		}
	}()
	go func() {
		for _, powerSupplier := range task.PowerSuppliers {
			go sendTaskFn(proposalId, types.PowerSupplier, powerSupplier.IdentityId, powerSupplier.NodeId, task.TaskId, prepareMsg, errCh)
		}
	}()

	go func() {
		for _, receiver := range task.Receivers {
			go sendTaskFn(proposalId, types.ResultSupplier, receiver.IdentityId, receiver.NodeId, task.TaskId, prepareMsg, errCh)
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

func (t *TwoPC) sendConfirmMsg(proposalId common.Hash, task *types.ScheduleTask, epoch uint64, startTime uint64) error {

	confirmMsg := makeConfirmMsg(epoch, startTime)
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

	errCh := make(chan error, len(task.Partners)+len(task.PowerSuppliers)+len(task.Receivers))

	go func() {
		for _, partner := range task.Partners {
			go sendConfirmMsgFn(proposalId, types.DataSupplier, partner.IdentityId, partner.NodeId, task.TaskId, confirmMsg, errCh)
		}
	}()
	go func() {
		for _, powerSupplier := range task.PowerSuppliers {
			go sendConfirmMsgFn(proposalId, types.PowerSupplier, powerSupplier.IdentityId, powerSupplier.NodeId, task.TaskId, confirmMsg, errCh)
		}
	}()

	go func() {
		for _, receiver := range task.Receivers {
			go sendConfirmMsgFn(proposalId, types.ResultSupplier, receiver.IdentityId, receiver.NodeId, task.TaskId, confirmMsg, errCh)
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

func (t *TwoPC) sendCommitMsg(proposalId common.Hash, task *types.ScheduleTask, startTime uint64) error {

	commitMsg := makeCommitMsg(startTime)

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

	errCh := make(chan error, len(task.Partners)+len(task.PowerSuppliers)+len(task.Receivers))

	go func() {
		for _, partner := range task.Partners {
			go sendCommitMsgFn(proposalId, types.DataSupplier, partner.IdentityId, partner.NodeId, task.TaskId, commitMsg, errCh)
		}
	}()
	go func() {
		for _, powerSupplier := range task.PowerSuppliers {
			go sendCommitMsgFn(proposalId, types.PowerSupplier, powerSupplier.IdentityId, powerSupplier.NodeId, task.TaskId, commitMsg, errCh)
		}
	}()

	go func() {
		for _, receiver := range task.Receivers {
			go sendCommitMsgFn(proposalId, types.ResultSupplier, receiver.IdentityId, receiver.NodeId, task.TaskId, commitMsg, errCh)
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



//func (t *TwoPC) makeConfirmTaskPeerDesc(proposalId common.Hash) *pb.ConfirmTaskPeerInfo {
//
//	dataSuppliers, powerSuppliers, receivers := make([]*pb.TaskPeerInfo, 0), make([]*pb.TaskPeerInfo, 0), make([]*pb.TaskPeerInfo, 0)
//
//	for _, vote := range t.state.GetPrepareVoteArr(proposalId) {
//		if vote.TaskRole == types.DataSupplier && nil != vote.PeerInfo {
//			dataSuppliers = append(dataSuppliers, types.ConvertTaskPeerInfo(vote.PeerInfo))
//		}
//		if vote.TaskRole == types.PowerSupplier && nil != vote.PeerInfo {
//			powerSuppliers = append(powerSuppliers, types.ConvertTaskPeerInfo(vote.PeerInfo))
//		}
//		if vote.TaskRole == types.ResultSupplier && nil != vote.PeerInfo {
//			receivers = append(receivers, types.ConvertTaskPeerInfo(vote.PeerInfo))
//		}
//	}
//	owner := t.state.GetSelfPeerInfo(proposalId)
//	if nil == owner {
//		return nil
//	}
//	return &pb.ConfirmTaskPeerInfo{
//		OwnerPeerInfo: types.ConvertTaskPeerInfo(owner),
//		DataSupplierPeerInfoList: dataSuppliers,
//		PowerSupplierPeerInfoList: powerSuppliers,
//		ResultReceiverPeerInfoList: receivers,
//	}
//}
