package twopc

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	"github.com/RosettaFlow/Carrier-Go/core/evengine"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	"github.com/RosettaFlow/Carrier-Go/handler"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/libp2p/go-libp2p-core/peer"
	"strings"
	"sync"
)




func (t *TwoPC) isConsensusTask(taskId string) bool {
	if _, ok := t.GetSendTaskWithOk(taskId); ok {
		return true
	}
	if _, ok := t.GetRecvTaskWithOk(taskId); ok {
		return true
	}
	return false
}

func (t *TwoPC) addSendTask(task *types.Task)  {
	t.sendTaskLock.Lock()
	t.sendTaskCache[task.TaskId()] = task
	t.sendTaskLock.Unlock()
}
func (t *TwoPC) addRecvTask(task *types.Task) {
	t.recvTaskLock.Lock()
	t.recvTaskCache[task.TaskId()] = task
	t.recvTaskLock.Unlock()
}
func (t *TwoPC) delTaskCache(taskId string) {
	t.delSendTask(taskId)
	t.delRecvTask(taskId)
}
func (t *TwoPC) delSendTask(taskId string) {
	t.sendTaskLock.Lock()
	delete(t.sendTaskCache, taskId)
	t.sendTaskLock.Unlock()
}
func (t *TwoPC) delRecvTask(taskId string) {
	t.recvTaskLock.Lock()
	delete(t.recvTaskCache, taskId)
	t.recvTaskLock.Unlock()
}
func (t *TwoPC) GetSendTaskWithOk(taskId string) (*types.Task, bool) {
	t.sendTaskLock.RLock()
	task, ok := t.sendTaskCache[taskId]
	t.sendTaskLock.RUnlock()
	return task, ok
}
func (t *TwoPC) GetSendTask(taskId string) *types.Task {
	task, _ := t.GetSendTaskWithOk(taskId)
	return task
}
func (t *TwoPC) GetRecvTaskWithOk(taskId string) (*types.Task, bool) {
	t.recvTaskLock.RLock()
	task, ok := t.recvTaskCache[taskId]
	t.recvTaskLock.RUnlock()
	return task, ok
}
func (t *TwoPC) GetRecvTask(taskId string) *types.Task {
	task, _ := t.GetRecvTaskWithOk(taskId)
	return task
}

func (t *TwoPC) addTaskResultCh(taskId string, resultCh chan<- *types.ConsensuResult) {
	t.taskResultLock.Lock()
	log.Debugf("AddTaskResultCh taskId: {%s}", taskId)
	t.taskResultChs[taskId] = resultCh
	t.taskResultLock.Unlock()
}
func (t *TwoPC) removeTaskResultCh(taskId string) {
	t.taskResultLock.Lock()
	log.Debugf("RemoveTaskResultCh taskId: {%s}", taskId)
	delete(t.taskResultChs, taskId)
	t.taskResultLock.Unlock()
}
func (t *TwoPC) collectTaskResultWillSendToSched(result *types.ConsensuResult) {
	t.taskResultCh <- result
}
func (t *TwoPC) sendConsensusTaskResultToSched (result *types.ConsensuResult) {
	t.taskResultLock.Lock()
	log.Debugf("Need SendTaskResultCh taskId: {%s}, result: {%s}", result.TaskId, result.String())
	if ch, ok := t.taskResultChs[result.TaskId]; ok {
		log.Debugf("Start SendTaskResultCh taskId: {%s}, result: {%s}", result.TaskId, result.String())
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
		t.delTaskCache(state.TaskId)
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
			log.Debugf("Started refresh proposalState loop, the proposalState direct be deadline, proposalId: {%s}, taskId: {%s}",
				id.String(), proposalState.TaskId)
			t.handleInvalidProposal(proposalState)
			continue
		}

		switch proposalState.CurrPeriodNum() {
		case ctypes.PeriodPrepare:
			if proposalState.IsPrepareTimeout() {
				log.Debugf("Started refresh proposalState loop, the proposalState was prepareTimeout, change to confirm epoch, proposalId: {%s}, taskId: {%s}",
					id.String(), proposalState.TaskId)
				proposalState.ChangeToConfirm(proposalState.PeriodStartTime + uint64(ctypes.PrepareMsgVotingTimeout.Milliseconds()))
				t.state.UpdateProposalState(proposalState)
			}
		case ctypes.PeriodConfirm:
			if proposalState.IsConfirmTimeout() {
				log.Debugf("Started refresh proposalState loop, the proposalState was confirmTimeout, change to commit epoch, proposalId: {%s}, taskId: {%s}",
					id.String(), proposalState.TaskId)
				proposalState.ChangeToCommit(proposalState.PeriodStartTime + uint64(ctypes.ConfirmMsgVotingTimeout.Milliseconds()))
				t.state.UpdateProposalState(proposalState)
			}
		case ctypes.PeriodCommit:
			if proposalState.IsCommitTimeout() {
				log.Debugf("Started refresh proposalState loop, the proposalState was commitTimeout, change to finished epoch, proposalId: {%s}, taskId: {%s}",
					id.String(), proposalState.TaskId)
				proposalState.ChangeToFinished(proposalState.PeriodStartTime + uint64(ctypes.CommitMsgEndingTimeout.Milliseconds()))
				//t.state.UpdateProposalState(proposalState)
				t.handleInvalidProposal(proposalState)
			}
		case ctypes.PeriodFinished:
			//
			if proposalState.IsDeadline() {
				log.Debugf("Started refresh proposalState loop, the proposalState was finished but coming deadline now, proposalId: {%s}, taskId: {%s}",
					id.String(), proposalState.TaskId)
				t.handleInvalidProposal(proposalState)
			}

		default:
			log.Errorf("Unknown the proposalState period, proposalId: {%s}, taskId: {%s}", id.String(), proposalState.TaskId)
			t.handleInvalidProposal(proposalState)
		}
	}
}

func (t *TwoPC) handleInvalidProposal(proposalState *ctypes.ProposalState) {

	log.Debugf("Call handleInvalidProposal(), handle and clean proposalState and task, proposalId: {%s}, taskId: {%s}, taskDir: {%s}", proposalState.ProposalId, proposalState.TaskId, proposalState.TaskDir.String())

	has, err := t.dataCenter.HasLocalTaskExecute(proposalState.TaskId)
	if nil != err {
		log.Errorf("Failed to query local task exec status with task on handleInvalidProposal(), taskId: {%s}, err: {%s}", proposalState.TaskId, err)
		// 最终 clean some data
		t.delProposalStateAndTask(proposalState.ProposalId)
		return
	}

	if has {
		log.Debugf("The local task have been executing, direct clean proposalStateAndTaskCache of consensus, taskId: {%s}", proposalState.TaskId)
		// 最终 clean some data
		t.delProposalStateAndTask(proposalState.ProposalId)
		return
	}


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

	} else {

		task, ok := t.GetRecvTaskWithOk(proposalState.TaskId)
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
			Content: fmt.Sprintf("%s for myself", evengine.TaskProposalStateDeadline.Msg),
			CreateTime: uint64(timeutils.UnixMsec()),
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
				CreateAt: uint64(timeutils.UnixMsec()),
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

		t.resourceMng.ReleaseLocalResourceWithTask("on consensus.handleInvalidProposal()", proposalState.TaskId, resource.SetAllReleaseResourceOption())

	}

	// 最终 clean some data
	t.delProposalStateAndTask(proposalState.ProposalId)
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
	selfIdentity *libTypes.OrganizationData,
	task *types.Task,
	) {

	log.Debugf("Start to call `driveTask`, proposalId: {%s}, taskId: {%s}, taskDir: {%s}, taskState: {%s}, taskRole: {%s}, myselfIdentityId: {%s}",
		proposalId.String(), task.TaskId(), taskDir.String(), taskState.String(), taskRole.String(), selfIdentity.Identity)

	selfVotePeerInfo := t.state.GetSelfPeerInfo(proposalId)
	if nil == selfVotePeerInfo {
		log.Errorf("Failed to find local cache about prepareVote myself internal resource, proposalId: {%s}, taskId: {%s}, taskDir: {%s}, taskState: {%s}, taskRole: {%s}, myselfIdentityId: {%s}",
			proposalId.String(), task.TaskId(), taskDir.String(), taskState.String(), taskRole.String(), selfIdentity.Identity)
		return
	}

	confirmTaskPeerInfo := t.state.GetConfirmTaskPeerInfo(proposalId)
	if nil == confirmTaskPeerInfo {
		log.Errorf("Failed to find local cache about prepareVote all peer resource {externalIP:externalPORT}, proposalId: {%s}, taskId: {%s}, taskDir: {%s}, taskState: {%s}, taskRole: {%s}, myselfIdentityId: {%s}",
			proposalId.String(), task.TaskId(), taskDir.String(), taskState.String(), taskRole.String(), selfIdentity.Identity)
		return
	}

	// Store task exec status
	if err := t.dataCenter.StoreLocalTaskExecuteStatus(task.TaskId()); nil != err {
		log.Errorf("Failed to store local task about exec status, proposalId: {%s}, taskId: {%s}, taskDir: {%s}, taskState: {%s}, taskRole: {%s}, myselfIdentityId: {%s}, err: {%s}",
			proposalId.String(), task.TaskId(), taskDir.String(), taskState.String(), taskRole.String(), selfIdentity.Identity, err)
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
			SelfVotePeerInfo: &types.PrepareVoteResource{
				Id: selfVotePeerInfo.Id,
				Ip: selfVotePeerInfo.Ip,
				Port: selfVotePeerInfo.Port,
				PartyId: selfVotePeerInfo.PartyId,
			},
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
				t.resourceMng.ReleaseLocalResourceWithTask("on consensus.driveTask()", task.TaskId(), resource.SetAllReleaseResourceOption())
				//// clean some data
				//t.delProposalStateAndTask(proposalId)
			}
		} else {
			<-taskWrap.ResultCh  // publish taskInfo to dataCenter done ..
			//// clean local proposalState and task cache
			//t.delProposalStateAndTask(proposalId)
		}
	}()
}

func (t *TwoPC) sendPrepareMsg(proposalId common.Hash, task *types.Task, startTime uint64) error {

	sendTaskFn := func(wg *sync.WaitGroup, proposalId common.Hash, taskRole types.TaskRole, partyId, identityId, nodeId, taskId string, errCh chan<- error) {

		defer wg.Done()

		var pid, err = p2p.HexPeerID(nodeId)
		if nil != err {
			errCh <- fmt.Errorf("failed to nodeId => peerId, proposalId: %s, taskId: %s, other peer taskRole: %s, other peer taskPartyId: %s, identityId: %s, pid: %s, err: %s",
				proposalId.String(), taskId, taskRole.String(), partyId, identityId, pid, err)
			return
		}

		prepareMsg, err := makePrepareMsgWithoutTaskRole(proposalId, task, startTime)

		if nil != err {
			errCh <- fmt.Errorf("failed to make prepareMsg, proposalId: %s, taskId: %s, other peer taskRole: %s, other peer taskPartyId: %s, identityId: %s, pid: %s, err: %s",
				proposalId.String(), taskId, taskRole.String(), partyId, identityId, pid, err)
			return
		}

		// set other peer's role and partyId
		prepareMsg.TaskRole = taskRole.Bytes()
		prepareMsg.TaskPartyId = []byte(partyId)

		if err = handler.SendTwoPcPrepareMsg(context.TODO(), t.p2p, pid, prepareMsg); nil != err {
			errCh <- fmt.Errorf("failed to call `SendTwoPcPrepareMsg` proposalId: %s, taskId: %s, other peer taskRole: %s, other peer taskPartyId: %s, identityId: %s, pid: %s, err: %s",
				proposalId.String(), taskId, taskRole.String(), partyId, identityId, pid, err)
			return
		}

		log.Debugf("Succceed to call `SendTwoPcPrepareMsg` proposalId: %s, taskId: %s, other peer taskRole: %s, other peer taskPartyId: %s, identityId: %s, pid: %s",
			proposalId.String(), taskId, taskRole.String(), partyId, identityId, pid)
	}

	size := (len(task.TaskData().MetadataSupplier) - 1) + len(task.TaskData().ResourceSupplier) + len(task.TaskData().Receivers)
	errCh := make(chan error, size)
	var wg sync.WaitGroup

	for i := 0; i < len(task.TaskData().MetadataSupplier); i++ {
		dataSupplier := task.TaskData().MetadataSupplier[i]
		// 排除掉 task 发起方 ...
		if task.TaskData().Identity != dataSupplier.Organization.Identity && task.TaskData().PartyId != dataSupplier.Organization.PartyId {
			wg.Add(1)
			go sendTaskFn(&wg, proposalId, types.DataSupplier, dataSupplier.Organization.PartyId,
				dataSupplier.Organization.Identity, dataSupplier.Organization.NodeId, task.TaskId(), errCh)
		}
	}
	for i := 0; i < len(task.TaskData().ResourceSupplier); i++ {
		powerSupplier := task.TaskData().ResourceSupplier[i]
		wg.Add(1)
		go sendTaskFn(&wg, proposalId, types.PowerSupplier, powerSupplier.Organization.PartyId,
			powerSupplier.Organization.Identity, powerSupplier.Organization.NodeId, task.TaskId(), errCh)
	}

	for i := 0; i < len(task.TaskData().Receivers); i++ {
		receiver := task.TaskData().Receivers[i]
		wg.Add(1)
		go sendTaskFn(&wg, proposalId, types.ResultSupplier, receiver.Receiver.PartyId,
			receiver.Receiver.Identity, receiver.Receiver.NodeId, task.TaskId(), errCh)
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

func (t *TwoPC) sendConfirmMsg(proposalId common.Hash, task *types.Task, startTime uint64) error {

	peerDesc := t.makeConfirmTaskPeerDesc(proposalId)
	// store the proposal about all partner peerInfo of task to local cache
	t.state.StoreConfirmTaskPeerInfo(proposalId, peerDesc)

	sendConfirmMsgFn := func(wg *sync.WaitGroup, proposalId common.Hash, taskRole types.TaskRole, taskPartyId, identityId, nodeId, taskId string, peerDesc *pb.ConfirmTaskPeerInfo, errCh chan<- error) {

		defer wg.Done()

		pid, err := p2p.HexPeerID(nodeId)
		if nil != err {
			errCh <- fmt.Errorf("failed to nodeId => peerId, proposalId: %s, taskId: %s, other peer's taskRole: %s, other peer's partyId: %s, other identityId: %s, pid: %s, err: %s",
				proposalId.String(), taskId, taskRole.String(), taskPartyId, identityId, pid, err)
			return
		}

		confirmMsg := makeConfirmMsg(proposalId, task, startTime)
		confirmMsg.PeerDesc = peerDesc

		// set other peer's role and partyId
		confirmMsg.TaskRole = taskRole.Bytes()
		confirmMsg.TaskPartyId = []byte(taskPartyId)

		// Send the ConfirmMsg to other peer
		if err := handler.SendTwoPcConfirmMsg(context.TODO(), t.p2p, pid, confirmMsg); nil != err {
			errCh <- fmt.Errorf("failed to call`SendTwoPcConfirmMsg` proposalId: %s, taskId: %s,other peer's taskRole: %s, other peer's partyId: %s, other identityId: %s, pid: %s, err: %s",
				proposalId.String(), taskId, taskRole.String(), taskPartyId, identityId, pid, err)
			errCh <- err
			return
		}

		log.Debugf("Succceed to call`SendTwoPcConfirmMsg` proposalId: %s, taskId: %s,other peer's taskRole: %s, other peer's partyId: %s, other identityId: %s, pid: %s",
			proposalId.String(), taskId, taskRole.String(), taskPartyId, identityId, pid)

	}

	size := (len(task.TaskData().MetadataSupplier) - 1) + len(task.TaskData().ResourceSupplier) + len(task.TaskData().Receivers)
	errCh := make(chan error, size)
	var wg sync.WaitGroup

	for i := 0; i < len(task.TaskData().MetadataSupplier); i++ {
		dataSupplier := task.TaskData().MetadataSupplier[i]
		// 排除掉 task 发起方 ...
		if task.TaskData().Identity != dataSupplier.Organization.Identity && task.TaskData().PartyId != dataSupplier.Organization.PartyId {
			wg.Add(1)
			go sendConfirmMsgFn(&wg, proposalId, types.DataSupplier, dataSupplier.Organization.PartyId,
				dataSupplier.Organization.Identity, dataSupplier.Organization.NodeId, task.TaskId(), peerDesc, errCh)
		}
	}
	for i := 0; i < len(task.TaskData().ResourceSupplier); i++ {
		powerSupplier := task.TaskData().ResourceSupplier[i]
		wg.Add(1)
		go sendConfirmMsgFn(&wg, proposalId, types.PowerSupplier, powerSupplier.Organization.PartyId,
			powerSupplier.Organization.Identity, powerSupplier.Organization.NodeId, task.TaskId(), peerDesc, errCh)
	}
	for i := 0; i < len(task.TaskData().Receivers); i++ {
		receiver := task.TaskData().Receivers[i]
		wg.Add(1)
		go sendConfirmMsgFn(&wg, proposalId, types.ResultSupplier, receiver.Receiver.PartyId,
			receiver.Receiver.Identity, receiver.Receiver.NodeId, task.TaskId(), peerDesc, errCh)
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

func (t *TwoPC) sendCommitMsg(proposalId common.Hash, task *types.Task, startTime uint64) error {

	sendCommitMsgFn := func(wg *sync.WaitGroup, proposalId common.Hash, taskRole types.TaskRole, taskPartyId, identityId, nodeId, taskId string, errCh chan<- error) {

		defer wg.Done()

		pid, err := p2p.HexPeerID(nodeId)
		if nil != err {
			errCh <- fmt.Errorf("failed to nodeId => peerId, proposalId: %s, taskId: %s, other peer's taskRole: %s, other peer's partyId: %s, identityId: %s, pid: %s, err: %s",
				proposalId.String(), taskId, taskRole.String(), taskPartyId, identityId, pid, err)
			return
		}
		commitMsg := makeCommitMsg(proposalId, task, startTime)
		// set other peer's role and partyId
		commitMsg.TaskRole = taskRole.Bytes()
		commitMsg.TaskPartyId = []byte(taskPartyId)

		// Send the ConfirmMsg to other peer
		if err := handler.SendTwoPcCommitMsg(context.TODO(), t.p2p, pid, commitMsg); nil != err {
			errCh <- fmt.Errorf("failed to call`SendTwoPcCommitMsg` proposalId: %s, taskId: %s,  other peer's taskRole: %s, other peer's partyId: %s, identityId: %s, pid: %s, err: %s",
				proposalId.String(), taskId, taskRole.String(), taskPartyId, identityId, pid, err)
			errCh <- err
			return
		}

		log.Debugf("Succceed to call`SendTwoPcCommitMsg` proposalId: %s, taskId: %s,  other peer's taskRole: %s, other peer's partyId: %s, identityId: %s, pid: %s",
			proposalId.String(), taskId, taskRole.String(), taskPartyId, identityId, pid)

	}


	size := (len(task.TaskData().MetadataSupplier) - 1) + len(task.TaskData().ResourceSupplier) + len(task.TaskData().Receivers)
	errCh := make(chan error, size)
	var wg sync.WaitGroup

	for i := 0; i < len(task.TaskData().MetadataSupplier); i++ {
		dataSupplier := task.TaskData().MetadataSupplier[i]
		// 排除掉 task 发起方 ...
		if task.TaskData().Identity != dataSupplier.Organization.Identity && task.TaskData().PartyId != dataSupplier.Organization.PartyId {
			wg.Add(1)
			go sendCommitMsgFn(&wg, proposalId, types.DataSupplier, dataSupplier.Organization.PartyId,
				dataSupplier.Organization.Identity, dataSupplier.Organization.NodeId, task.TaskId(), errCh)
		}
	}
	for i := 0; i < len(task.TaskData().ResourceSupplier); i++ {
		powerSupplier := task.TaskData().ResourceSupplier[i]
		wg.Add(1)
		go sendCommitMsgFn(&wg, proposalId, types.PowerSupplier, powerSupplier.Organization.PartyId,
			powerSupplier.Organization.Identity, powerSupplier.Organization.NodeId, task.TaskId(), errCh)
	}
	for i := 0; i < len(task.TaskData().Receivers); i++ {
		receiver := task.TaskData().Receivers[i]
		wg.Add(1)
		go sendCommitMsgFn(&wg, proposalId, types.ResultSupplier, receiver.Receiver.PartyId,
			receiver.Receiver.Identity, receiver.Receiver.NodeId, task.TaskId(), errCh)
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

