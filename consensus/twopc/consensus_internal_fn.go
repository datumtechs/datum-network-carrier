package twopc

import (
	"github.com/RosettaFlow/Carrier-Go/common"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	"github.com/RosettaFlow/Carrier-Go/types"
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
	t.state.DelProposalState(proposalId)
	t.state.RemovePrepareVoteState(proposalId)
	t.state.RemoveConfirmVoteState(proposalId)
	t.state.CleanPrepareVoteState(proposalId)
	t.state.CleanConfirmVoteState(proposalId)
}