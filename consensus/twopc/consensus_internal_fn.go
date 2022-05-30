package twopc

import (
	"context"
	"fmt"
	"github.com/datumtechs/datum-network-carrier/common"
	"github.com/datumtechs/datum-network-carrier/common/timeutils"
	"github.com/datumtechs/datum-network-carrier/common/traceutil"
	ctypes "github.com/datumtechs/datum-network-carrier/consensus/twopc/types"
	"github.com/datumtechs/datum-network-carrier/core/evengine"
	"github.com/datumtechs/datum-network-carrier/core/resource"
	twopcpb "github.com/datumtechs/datum-network-carrier/pb/carrier/netmsg/consensus/twopc"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"github.com/datumtechs/datum-network-carrier/p2p"
	"github.com/datumtechs/datum-network-carrier/types"
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/peer"
	"strings"
	"sync"
	"time"
)

func (t *Twopc) removeOrgProposalStateAndTask(proposalId common.Hash, partyId string) {

	if hasCache := t.state.HasOrgProposalWithProposalId(proposalId); hasCache {
		orgState, hasOrgState := t.state.QueryOrgProposalStateWithProposalIdAndPartyId(proposalId, partyId)
		if hasOrgState {
			log.Infof("Start remove org proposalState and task cache on Consensus, proposalId {%s}, taskId {%s}, partyId: {%s}",
				proposalId, orgState.GetTaskId(), partyId)
			t.state.RemoveOrgProposalStateAnyCache(proposalId, orgState.GetTaskId(), partyId) // remove task/proposal state/prepare vote/ confirm vote and wal with partyId
		}
	} else {
		log.Infof("Start remove confirm taskPeerInfo when proposalState is empty, proposalId {%s}, final remove partyId: {%s}", proposalId, partyId)
		t.state.RemoveConfirmTaskPeerInfo(proposalId) // remove confirm peers and wal
	}
}

func (t *Twopc) storeOrgProposalState(orgState *ctypes.OrgProposalState) {

	first := t.state.HasNotOrgProposalWithProposalIdAndPartyId(orgState.GetProposalId(), orgState.GetTaskOrg().GetPartyId())

	t.state.StoreOrgProposalState(orgState)
	t.wal.StoreOrgProposalState(orgState)
	if first {
		// v 0.3.0 add proposal state monitor
		t.addmonitor(orgState)
	}
}

func (t *Twopc) addmonitor(orgState *ctypes.OrgProposalState) {

	var (
		when int64
		next int64
	)

	switch orgState.GetPeriodNum() {

	case ctypes.PeriodPrepare:
		when = orgState.GetPrepareExpireTime()
		next = orgState.GetConfirmExpireTime()
	case ctypes.PeriodConfirm:
		when = orgState.GetConfirmExpireTime()
		next = orgState.GetCommitExpireTime()
	case ctypes.PeriodCommit, ctypes.PeriodFinished:
		if !orgState.IsDeadline() {
			if orgState.IsCommitTimeout() { // finished period => when: deadline, next: 0
				when = orgState.GetDeadlineExpireTime()
				next = 0
			} else { // commit period => when: commit, next: deadline
				when = orgState.GetCommitExpireTime()
				next = orgState.GetDeadlineExpireTime()
			}
		}
	default:
		log.Warnf("It is a unknown period org state on twopc.refreshProposalStateMonitor")
		return
	}

	monitor := ctypes.NewProposalStateMonitor(orgState, when, next, nil)

	monitor.SetCallBackFn(func(orgState *ctypes.OrgProposalState) {

		t.state.proposalsLock.Lock()
		defer t.state.proposalsLock.Unlock()

		switch orgState.GetPeriodNum() {

		case ctypes.PeriodPrepare:

			if orgState.IsPrepareTimeout() {
				log.Debugf("Started refresh org proposalState, the org proposalState was prepareTimeout, change to confirm epoch on twopc.refreshProposalStateMonitor, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
					orgState.GetProposalId().String(), orgState.GetTaskId(), orgState.GetTaskOrg().GetPartyId())

				orgState.ChangeToConfirm()
				monitor.SetNext(orgState.GetCommitExpireTime())
				t.state.StoreOrgProposalStateUnsafe(orgState)
				t.wal.StoreOrgProposalState(orgState)

			} else {
				// There is only one possibility that the `prepare epoch` expire time is set incorrectly,
				// and the monitor arrived early. At this time, you need to ensure the `prepare epoch` expire time again (reset).
				monitor.SetWhen(orgState.GetPrepareExpireTime())
				monitor.SetNext(orgState.GetConfirmExpireTime())
			}

		case ctypes.PeriodConfirm:

			if orgState.IsConfirmTimeout() {
				log.Debugf("Started refresh org proposalState, the org proposalState was confirmTimeout, change to commit epoch on twopc.refreshProposalStateMonitor, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
					orgState.GetProposalId().String(), orgState.GetTaskId(), orgState.GetTaskOrg().GetPartyId())

				orgState.ChangeToCommit()
				monitor.SetNext(orgState.GetDeadlineExpireTime()) // Ensure that the next expire time is `dealine time`.
				t.state.StoreOrgProposalStateUnsafe(orgState)
				t.wal.StoreOrgProposalState(orgState)

			} else {
				// There is only one possibility that the `confirm epoch` expire time is set incorrectly,
				// and the monitor arrived early. At this time, you need to ensure the `confirm epoch` expire time again (reset).
				monitor.SetWhen(orgState.GetConfirmExpireTime())
				monitor.SetNext(orgState.GetCommitExpireTime())
			}

		case ctypes.PeriodCommit, ctypes.PeriodFinished:

			// when epoch is commit or finished AND time is dealine.
			if orgState.IsDeadline() {
				identity, err := t.resourceMng.GetDB().QueryIdentity()
				if nil != err {
					log.WithError(err).Errorf("Failed to query local identity on proposate monitor expire deadline on twopc.refreshProposalStateMonitor, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
						orgState.GetProposalId().String(), orgState.GetTaskId(), orgState.GetTaskOrg().GetPartyId())
					return
				}

				log.Debugf("Started refresh org proposalState, the org proposalState was finished, but coming deadline now on twopc.refreshProposalStateMonitor, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
					orgState.GetProposalId().String(), orgState.GetTaskId(), orgState.GetTaskOrg().GetPartyId())

				_, ok := t.state.QueryProposalTaskWithTaskIdAndPartyId(orgState.GetTaskId(), orgState.GetTaskOrg().GetPartyId())

				t.state.RemoveProposalTaskWithTaskIdAndPartyId(orgState.GetTaskId(), orgState.GetTaskOrg().GetPartyId())                   // remove proposal task with partyId
				t.state.RemoveOrgProposalStateWithProposalIdAndPartyIdUnsafe(orgState.GetProposalId(), orgState.GetTaskOrg().GetPartyId()) // remove state with partyId
				t.state.RemoveOrgPrepareVoteState(orgState.GetProposalId(), orgState.GetTaskOrg().GetPartyId())                            // remove prepare vote with partyId
				t.state.RemoveOrgConfirmVoteState(orgState.GetProposalId(), orgState.GetTaskOrg().GetPartyId())                            // remove confirm vote with partyId

				if _, has := t.state.proposalSet[orgState.GetProposalId()]; !has {
					t.removeConfirmTaskPeerInfo(orgState.GetProposalId()) // remove confirm peers and wal
				}
				t.wal.DeleteState(t.wal.GetProposalTaskCacheKey(orgState.GetTaskId(), orgState.GetTaskOrg().GetPartyId()))
				t.wal.DeleteState(t.wal.GetPrepareVotesKey(orgState.GetProposalId(), orgState.GetTaskOrg().GetPartyId()))
				t.wal.DeleteState(t.wal.GetConfirmVotesKey(orgState.GetProposalId(), orgState.GetTaskOrg().GetPartyId()))
				t.wal.DeleteState(t.wal.GetProposalSetKey(orgState.GetProposalId(), orgState.GetTaskOrg().GetPartyId()))

				if !ok {
					log.Errorf("Not found the proposalTask, skip this proposal state on twopc.refreshProposalStateMonitor, taskId: {%s}, partyId: {%s}",
						orgState.GetTaskId(), orgState.GetTaskOrg().GetPartyId())
					return
				}

				_, err = t.resourceMng.GetDB().QueryLocalTask(orgState.GetTaskId())
				if nil == err {
					// release local resource and clean some data  (on task partenr)
					t.resourceMng.GetDB().StoreTaskEvent(&carriertypespb.TaskEvent{
						Type:       evengine.TaskProposalStateDeadline.GetType(),
						IdentityId: identity.GetIdentityId(),
						PartyId:    orgState.GetTaskOrg().GetPartyId(),
						TaskId:     orgState.GetTaskId(),
						Content:    fmt.Sprintf("proposalId: %s, %s for myself", orgState.GetProposalId().TerminalString(), evengine.TaskProposalStateDeadline.GetMsg()),
						CreateAt:   timeutils.UnixMsecUint64(),
					})

					t.stopTaskConsensus("the proposalState had been deadline",
						orgState.GetProposalId(),
						orgState.GetTaskId(),
						orgState.GetTaskRole(),
						carriertypespb.TaskRole_TaskRole_Sender,
						&carriertypespb.TaskOrganization{
							PartyId:    orgState.GetTaskOrg().GetPartyId(),
							NodeName:   identity.GetNodeName(),
							NodeId:     identity.GetNodeId(),
							IdentityId: identity.GetIdentityId(),
						},
						orgState.GetTaskSender(), types.TaskConsensusInterrupt)
				}
			} else {
				// OR epoch is commit but time just is commit time out
				if orgState.IsCommitTimeout() {
					log.Debugf("Started refresh org proposalState, the org proposalState was commitTimeout, change to finished epoch on twopc.refreshProposalStateMonitor, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
						orgState.GetProposalId().String(), orgState.GetTaskId(), orgState.GetTaskOrg().GetPartyId())

					orgState.ChangeToFinished()
					monitor.SetNext(0) // clean next target timestamp, let monitor would been able to clean from monitor queue.
					t.state.StoreOrgProposalStateUnsafe(orgState)
					t.wal.StoreOrgProposalState(orgState)

				} else {
					// There is only one possibility that the `commit epoch` expire time is set incorrectly,
					// and the monitor arrived early. At this time, you need to ensure the `commit epoch` expire time again (reset).
					monitor.SetWhen(orgState.GetCommitExpireTime())
					monitor.SetNext(orgState.GetDeadlineExpireTime())
				}
			}

		default:
			log.Errorf("Unknown the proposalState period on twopc.refreshProposalStateMonitor,  proposalId: {%s}, taskId: {%s}, partyId: {%s}",
				orgState.GetProposalId().String(), orgState.GetTaskId(), orgState.GetTaskOrg().GetPartyId())
		}
	})

	// add monitor to queue
	t.state.AddMonitor(monitor)
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

		if vote.GetMsgOption().GetSenderRole() == carriertypespb.TaskRole_TaskRole_DataSupplier && nil != vote.GetPeerInfo() {
			dataSuppliers = append(dataSuppliers, types.ConvertTaskPeerInfo(vote.GetPeerInfo()))
		}
		if vote.GetMsgOption().GetSenderRole() == carriertypespb.TaskRole_TaskRole_PowerSupplier && nil != vote.GetPeerInfo() {
			powerSuppliers = append(powerSuppliers, types.ConvertTaskPeerInfo(vote.GetPeerInfo()))
		}
		if vote.GetMsgOption().GetSenderRole() == carriertypespb.TaskRole_TaskRole_Receiver && nil != vote.GetPeerInfo() {
			receivers = append(receivers, types.ConvertTaskPeerInfo(vote.GetPeerInfo()))
		}
	}
	return &twopcpb.ConfirmTaskPeerInfo{
		DataSupplierPeerInfos:   dataSuppliers,
		PowerSupplierPeerInfos:  powerSuppliers,
		ResultReceiverPeerInfos: receivers,
	}
}

func (t *Twopc) checkProposalStateMonitors(now int64, syncCall bool) int64 {
	return t.state.CheckProposalStateMonitors(now, syncCall)
}

func (t *Twopc) proposalStateMonitorTimer() *time.Timer {
	return t.state.Timer()
}

func (t *Twopc) proposalStateMonitorsLen() int {
	return t.state.MonitorsLen()
}

func (t *Twopc) stopTaskConsensus(
	reason string,
	proposalId common.Hash,
	taskId string,
	senderRole, receiverRole carriertypespb.TaskRole,
	sender, receiver *carriertypespb.TaskOrganization,
	taskActionStatus types.TaskActionStatus,
) {

	log.Debugf("Call stopTaskConsensus() to interrupt consensus msg %s with role is %s, proposalId: {%s}, taskId: {%s}, partyId: {%s}, remote partyId: {%s}, taskActionStatus: {%s}",
		reason, senderRole, proposalId.String(), taskId, sender.GetPartyId(), receiver.GetPartyId(), taskActionStatus.String())
	// Send consensus result to interrupt consensus epoch and clean some data (on task sender)
	if senderRole == carriertypespb.TaskRole_TaskRole_Sender {
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
	localTaskRole carriertypespb.TaskRole,
	localTaskOrganization *carriertypespb.TaskOrganization,
	remoteTaskRole carriertypespb.TaskRole,
	remoteTaskOrganization *carriertypespb.TaskOrganization,
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

	sendTaskFn := func(wg *sync.WaitGroup, sender, receiver *carriertypespb.TaskOrganization, senderRole, receiverRole carriertypespb.TaskRole, errCh chan<- error) {

		defer wg.Done()

		var pid, err = p2p.HexPeerID(receiver.GetNodeId())
		if nil != err {
			errCh <- fmt.Errorf("failed to nodeId => peerId when send prepareMsg, proposalId: %s, taskId: %s, other peer taskRole: %s, other peer taskPartyId: %s, identityId: %s, pid: %s, err: %s",
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

		log.WithField("traceId", traceutil.GenerateTraceID(prepareMsg)).Debugf("Succeed to call `sendPrepareMsg.%s` proposalId: %s, taskId: %s, other peer taskRole: %s, other peer taskPartyId: %s, identityId: %s, pid: %s",
			logdesc, proposalId.String(), task.GetTaskId(), receiverRole.String(), receiver.GetPartyId(), receiver.GetIdentityId(), pid)
	}

	size := (len(task.GetTaskData().GetDataSuppliers())) + len(task.GetTaskData().GetPowerSuppliers()) + len(task.GetTaskData().GetReceivers())
	errCh := make(chan error, size)
	var wg sync.WaitGroup

	for i := 0; i < len(task.GetTaskData().GetDataSuppliers()); i++ {

		wg.Add(1)
		dataSupplier := task.GetTaskData().GetDataSuppliers()[i]
		receiver := dataSupplier
		go sendTaskFn(&wg, sender, receiver, carriertypespb.TaskRole_TaskRole_Sender, carriertypespb.TaskRole_TaskRole_DataSupplier, errCh)

	}
	for i := 0; i < len(task.GetTaskData().GetPowerSuppliers()); i++ {

		wg.Add(1)
		powerSupplier := task.GetTaskData().GetPowerSuppliers()[i]
		receiver := powerSupplier
		go sendTaskFn(&wg, sender, receiver, carriertypespb.TaskRole_TaskRole_Sender, carriertypespb.TaskRole_TaskRole_PowerSupplier, errCh)

	}

	for i := 0; i < len(task.GetTaskData().GetReceivers()); i++ {

		wg.Add(1)
		receiver := task.GetTaskData().GetReceivers()[i]
		go sendTaskFn(&wg, sender, receiver, carriertypespb.TaskRole_TaskRole_Sender, carriertypespb.TaskRole_TaskRole_Receiver, errCh)
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

func (t *Twopc) replyTaskConsensusResult(result *types.TaskConsResult) {
	go func(result *types.TaskConsResult) { // asynchronous transmission to reduce Chan blocking
		t.taskConsResultCh <- result
	}(result)
}

func (t *Twopc) sendNeedReplayScheduleTask(task *types.NeedReplayScheduleTask) {
	go func(task *types.NeedReplayScheduleTask) { // asynchronous transmission to reduce Chan blocking
		t.needReplayScheduleTaskCh <- task
	}(task)
}

func (t *Twopc) sendNeedExecuteTask(task *types.NeedExecuteTask) {
	go func(task *types.NeedExecuteTask) { // asynchronous transmission to reduce Chan blocking
		t.needExecuteTaskCh <- task
	}(task)
}

func (t *Twopc) sendPrepareVote(pid peer.ID, sender, receiver *carriertypespb.TaskOrganization, req *twopcpb.PrepareVote) error {
	if types.IsNotSameTaskOrg(sender, receiver) {
		//return handler.SendTwoPcPrepareVote(context.TODO(), t.p2p, pid, req)
		return t.p2p.Broadcast(context.TODO(), req)
	} else {
		return t.sendLocalPrepareVote(pid, req)
	}
}

func (t *Twopc) sendConfirmMsg(proposalId common.Hash, task *types.Task, peers *twopcpb.ConfirmTaskPeerInfo, option types.TwopcMsgOption, startTime uint64) error {

	sender := task.GetTaskSender()

	sendConfirmMsgFn := func(wg *sync.WaitGroup, sender, receiver *carriertypespb.TaskOrganization, senderRole, receiverRole carriertypespb.TaskRole, errCh chan<- error) {

		defer wg.Done()

		pid, err := p2p.HexPeerID(receiver.GetNodeId())
		if nil != err {
			errCh <- fmt.Errorf("failed to nodeId => peerId when send confirmMsg, proposalId: %s, taskId: %s, other peer's taskRole: %s, other peer's partyId: %s, other identityId: %s, pid: %s, err: %s",
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

		log.WithField("traceId", traceutil.GenerateTraceID(confirmMsg)).Debugf("Succeed to call`sendConfirmMsg.%s` proposalId: %s, taskId: %s,other peer's taskRole: %s, other peer's partyId: %s, other identityId: %s, pid: %s",
			logdesc, proposalId.String(), task.GetTaskId(), receiverRole.String(), receiver.GetPartyId(), receiver.GetIdentityId(), pid)

	}

	size := (len(task.GetTaskData().GetDataSuppliers())) + len(task.GetTaskData().GetPowerSuppliers()) + len(task.GetTaskData().GetReceivers())
	errCh := make(chan error, size)
	var wg sync.WaitGroup

	for i := 0; i < len(task.GetTaskData().GetDataSuppliers()); i++ {

		wg.Add(1)
		dataSupplier := task.GetTaskData().GetDataSuppliers()[i]
		receiver := dataSupplier
		go sendConfirmMsgFn(&wg, sender, receiver, carriertypespb.TaskRole_TaskRole_Sender, carriertypespb.TaskRole_TaskRole_DataSupplier, errCh)

	}
	for i := 0; i < len(task.GetTaskData().GetPowerSuppliers()); i++ {

		wg.Add(1)
		powerSupplier := task.GetTaskData().GetPowerSuppliers()[i]
		receiver := powerSupplier
		go sendConfirmMsgFn(&wg, sender, receiver, carriertypespb.TaskRole_TaskRole_Sender, carriertypespb.TaskRole_TaskRole_PowerSupplier, errCh)

	}

	for i := 0; i < len(task.GetTaskData().GetReceivers()); i++ {

		wg.Add(1)
		receiver := task.GetTaskData().GetReceivers()[i]
		go sendConfirmMsgFn(&wg, sender, receiver, carriertypespb.TaskRole_TaskRole_Sender, carriertypespb.TaskRole_TaskRole_Receiver, errCh)
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

func (t *Twopc) sendConfirmVote(pid peer.ID, sender, receiver *carriertypespb.TaskOrganization, req *twopcpb.ConfirmVote) error {
	if types.IsNotSameTaskOrg(sender, receiver) {
		//return handler.SendTwoPcConfirmVote(context.TODO(), t.p2p, pid, req)
		return t.p2p.Broadcast(context.TODO(), req)
	} else {
		return t.sendLocalConfirmVote(pid, req)
	}
}

func (t *Twopc) sendCommitMsg(proposalId common.Hash, task *types.Task, option types.TwopcMsgOption, startTime uint64) error {

	sender := task.GetTaskSender()

	sendCommitMsgFn := func(wg *sync.WaitGroup, sender, receiver *carriertypespb.TaskOrganization, senderRole, receiverRole carriertypespb.TaskRole, errCh chan<- error) {

		defer wg.Done()

		pid, err := p2p.HexPeerID(receiver.GetNodeId())
		if nil != err {
			errCh <- fmt.Errorf("failed to nodeId => peerId when send commitMsg, proposalId: %s, taskId: %s, other peer's taskRole: %s, other peer's partyId: %s, identityId: %s, pid: %s, err: %s",
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

		log.WithField("traceId", traceutil.GenerateTraceID(commitMsg)).Debugf("Succeed to call`sendCommitMsg.%s` proposalId: %s, taskId: %s,  other peer's taskRole: %s, other peer's partyId: %s, identityId: %s, pid: %s",
			logdesc, proposalId.String(), task.GetTaskId(), receiverRole.String(), receiver.GetPartyId(), receiver.GetIdentityId(), pid)

	}

	size := (len(task.GetTaskData().GetDataSuppliers())) + len(task.GetTaskData().GetPowerSuppliers()) + len(task.GetTaskData().GetReceivers())
	errCh := make(chan error, size)
	var wg sync.WaitGroup

	for i := 0; i < len(task.GetTaskData().GetDataSuppliers()); i++ {

		wg.Add(1)
		dataSupplier := task.GetTaskData().GetDataSuppliers()[i]
		receiver := dataSupplier
		go sendCommitMsgFn(&wg, sender, receiver, carriertypespb.TaskRole_TaskRole_Sender, carriertypespb.TaskRole_TaskRole_DataSupplier, errCh)

	}
	for i := 0; i < len(task.GetTaskData().GetPowerSuppliers()); i++ {

		wg.Add(1)
		powerSupplier := task.GetTaskData().GetPowerSuppliers()[i]
		receiver := powerSupplier
		go sendCommitMsgFn(&wg, sender, receiver, carriertypespb.TaskRole_TaskRole_Sender, carriertypespb.TaskRole_TaskRole_PowerSupplier, errCh)

	}

	for i := 0; i < len(task.GetTaskData().GetReceivers()); i++ {

		wg.Add(1)
		receiver := task.GetTaskData().GetReceivers()[i]
		go sendCommitMsgFn(&wg, sender, receiver, carriertypespb.TaskRole_TaskRole_Sender, carriertypespb.TaskRole_TaskRole_Receiver, errCh)
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

func verifyPartyRole(partyId string, role carriertypespb.TaskRole, task *types.Task) bool {

	switch role {
	case carriertypespb.TaskRole_TaskRole_Sender:
		if partyId == task.GetTaskSender().GetPartyId() {
			return true
		}
	case carriertypespb.TaskRole_TaskRole_DataSupplier:
		for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
			if partyId == dataSupplier.GetPartyId() {
				return true
			}
		}
	case carriertypespb.TaskRole_TaskRole_PowerSupplier:
		for _, powerSupplier := range task.GetTaskData().GetPowerSuppliers() {
			if partyId == powerSupplier.GetPartyId() {
				return true
			}
		}
	case carriertypespb.TaskRole_TaskRole_Receiver:
		for _, receiver := range task.GetTaskData().GetReceivers() {
			if partyId == receiver.GetPartyId() {
				return true
			}
		}
	}
	return false
}

func fetchLocalIdentityAndOriganizationFromMsg(resourceMng *resource.Manager, option *types.MsgOption, task *types.Task) (*carriertypespb.Organization, *carriertypespb.TaskOrganization, *carriertypespb.TaskOrganization, error) {
	identity, err := resourceMng.GetDB().QueryIdentity()
	if nil != err {
		//log.WithError(err).Errorf("Failed to call `QueryIdentity()` %s, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
		//	logdesc, option.GetProposalId().String(), task.GetTaskId(), option.GetReceiverRole().String(), option.GetReceiverPartyId())
		return nil, nil, nil, fmt.Errorf("query local identity failed, %s", err)
	}

	sender := fetchOrgByPartyRole(option.GetSenderPartyId(), option.GetSenderRole(), task)
	receiver := fetchOrgByPartyRole(option.GetReceiverPartyId(), option.GetReceiverRole(), task)
	if nil == sender || nil == receiver {
		//log.Errorf("Failed to check msg.MsgOption sender and receiver %s, someone is empty, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}",
		//	logdesc, option.GetProposalId().String(), task.GetTaskId(), option.GetReceiverRole().String(), option.GetReceiverPartyId())
		return nil, nil, nil, fmt.Errorf("%s, sender and receiver of msg that someone is empty", ctypes.ErrConsensusMsgInvalid)
	}

	// verify the receiver is myself ?
	if identity.GetIdentityId() != receiver.GetIdentityId() {
		//log.Warnf("Warning verify receiver identityId of msg %s, receiver is not me, proposalId: {%s}, taskId: {%s}, role: {%s}, partyId: {%s}, my identityId: {%s}, receiver identityId: {%s}",
		//	logdesc, option.GetProposalId().String(), task.GetTaskId(), option.GetReceiverRole().String(), option.GetReceiverPartyId(), identity.GetIdentityId(), receiver.GetIdentityId())
		return nil, nil, nil, fmt.Errorf("%s, receiver is not current identity", ctypes.ErrConsensusMsgInvalid)
	}
	return identity, sender, receiver, nil
}

func fetchOrgByPartyRole(partyId string, role carriertypespb.TaskRole, task *types.Task) *carriertypespb.TaskOrganization {

	switch role {
	case carriertypespb.TaskRole_TaskRole_Sender:
		if partyId == task.GetTaskSender().GetPartyId() {
			return task.GetTaskSender()
		}
	case carriertypespb.TaskRole_TaskRole_DataSupplier:
		for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
			if partyId == dataSupplier.GetPartyId() {
				return dataSupplier
			}
		}
	case carriertypespb.TaskRole_TaskRole_PowerSupplier:
		for _, powerSupplier := range task.GetTaskData().GetPowerSuppliers() {
			if partyId == powerSupplier.GetPartyId() {
				return powerSupplier
			}
		}
	case carriertypespb.TaskRole_TaskRole_Receiver:
		for _, receiver := range task.GetTaskData().GetReceivers() {
			if partyId == receiver.GetPartyId() {
				return receiver
			}
		}
	}
	return nil
}

func (t *Twopc) verifyPrepareVoteRoleIsTaskPartner(identityId, partyId string, role carriertypespb.TaskRole, task *types.Task) (bool, error) {
	var find bool
	switch role {
	case carriertypespb.TaskRole_TaskRole_DataSupplier:
		for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
			// identity + partyId
			if dataSupplier.GetIdentityId() == identityId && dataSupplier.GetPartyId() == partyId {
				find = true
				break
			}
		}
	case carriertypespb.TaskRole_TaskRole_PowerSupplier:
		for _, powerSupplier := range task.GetTaskData().GetPowerSuppliers() {
			// identity + partyId
			if powerSupplier.GetIdentityId() == identityId && powerSupplier.GetPartyId() == partyId {
				find = true
				break
			}
		}
	case carriertypespb.TaskRole_TaskRole_Receiver:
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
func (t *Twopc) verifyConfirmVoteRoleIsTaskPartner(identityId, partyId string, role carriertypespb.TaskRole, task *types.Task) (bool, error) {
	var identityValid bool
	switch role {
	case carriertypespb.TaskRole_TaskRole_DataSupplier:

		for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
			// identity + partyId
			if dataSupplier.GetIdentityId() == identityId && dataSupplier.GetPartyId() == partyId {
				identityValid = true
				break
			}
		}
	case carriertypespb.TaskRole_TaskRole_PowerSupplier:

		for _, powerSupplier := range task.GetTaskData().GetPowerSuppliers() {
			// identity + partyId
			if powerSupplier.GetIdentityId() == identityId && powerSupplier.GetPartyId() == partyId {
				identityValid = true
				break
			}
		}
	case carriertypespb.TaskRole_TaskRole_Receiver:

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
	return t.state.QueryConfirmTaskPeerInfo(proposalId)
}

func (t *Twopc) mustGetConfirmTaskPeerInfo(proposalId common.Hash) *twopcpb.ConfirmTaskPeerInfo {
	return t.state.MustQueryConfirmTaskPeerInfo(proposalId)
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
				proposalTaskPB := &carriertypespb.ProposalTask{}
				if err := proto.Unmarshal(value, proposalTaskPB); err != nil {
					t.wal.DeleteState(key)
					return fmt.Errorf("unmarshal proposalTask failed, %s", err)
				}

				cache, ok := t.state.proposalTaskCache[taskId]
				if !ok {
					cache = make(map[string]*ctypes.ProposalTask, 0)
				}
				cache[partyId] = ctypes.NewProposalTask(common.HexToHash(proposalTaskPB.GetProposalId()), taskId, proposalTaskPB.GetCreateAt())
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

				libOrgProposalState := &carriertypespb.OrgProposalState{}
				if err := proto.Unmarshal(value, libOrgProposalState); err != nil {
					return fmt.Errorf("unmarshal org proposalState failed, %s", err)
				}

				orgState := ctypes.NewOrgProposalStateWithFields(
					proposalId,
					libOrgProposalState.GetTaskId(),
					libOrgProposalState.GetTaskRole(),
					libOrgProposalState.GetTaskSender(),
					libOrgProposalState.GetTaskOrg(),
					ctypes.ProposalStatePeriod(libOrgProposalState.GetPeriodNum()),
					libOrgProposalState.GetDeadlineDuration(),
					libOrgProposalState.GetCreateAt(),
					libOrgProposalState.GetStartAt())

				t.state.StoreOrgProposalStateUnsafe(orgState)
				// v 0.3.0 add proposal state monitor
				t.addmonitor(orgState)
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

				vote := &carriertypespb.PrepareVote{}
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

				vote := &carriertypespb.ConfirmVote{}
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
