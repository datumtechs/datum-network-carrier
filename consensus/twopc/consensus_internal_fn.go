package twopc

import "github.com/RosettaFlow/Carrier-Go/types"

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
