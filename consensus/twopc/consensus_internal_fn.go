package twopc

import (
	"context"
	"fmt"
	"github.com/Metisnetwork/Metis-Carrier/common"
	"github.com/Metisnetwork/Metis-Carrier/common/timeutils"
	"github.com/Metisnetwork/Metis-Carrier/common/traceutil"
	ctypes "github.com/Metisnetwork/Metis-Carrier/consensus/twopc/types"
	"github.com/Metisnetwork/Metis-Carrier/core/evengine"
	twopcpb "github.com/Metisnetwork/Metis-Carrier/lib/netmsg/consensus/twopc"
	libtypes "github.com/Metisnetwork/Metis-Carrier/lib/types"
	"github.com/Metisnetwork/Metis-Carrier/p2p"
	"github.com/Metisnetwork/Metis-Carrier/types"
	"github.com/ethereum/go-ethereum/crypto"
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
					t.resourceMng.GetDB().StoreTaskEvent(&libtypes.TaskEvent{
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
						libtypes.TaskRole_TaskRole_Sender,
						&libtypes.TaskOrganization{
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

		if vote.GetMsgOption().GetSenderRole() == libtypes.TaskRole_TaskRole_DataSupplier && nil != vote.GetPeerInfo() {
			dataSuppliers = append(dataSuppliers, types.ConvertTaskPeerInfo(vote.GetPeerInfo()))
		}
		if vote.GetMsgOption().GetSenderRole() == libtypes.TaskRole_TaskRole_PowerSupplier && nil != vote.GetPeerInfo() {
			powerSuppliers = append(powerSuppliers, types.ConvertTaskPeerInfo(vote.GetPeerInfo()))
		}
		if vote.GetMsgOption().GetSenderRole() == libtypes.TaskRole_TaskRole_Receiver && nil != vote.GetPeerInfo() {
			receivers = append(receivers, types.ConvertTaskPeerInfo(vote.GetPeerInfo()))
		}
	}
	return &twopcpb.ConfirmTaskPeerInfo{
		DataSupplierPeerInfos:   dataSuppliers,
		PowerSupplierPeerInfos:  powerSuppliers,
		ResultReceiverPeerInfos: receivers,
	}
}

func (t *Twopc) checkProposalStateMonitors(now int64) int64 {
	return t.state.CheckProposalStateMonitors(now)
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
	senderRole, receiverRole libtypes.TaskRole,
	sender, receiver *libtypes.TaskOrganization,
	taskActionStatus types.TaskActionStatus,
) {

	log.Debugf("Call stopTaskConsensus() to interrupt consensus msg %s with role is %s, proposalId: {%s}, taskId: {%s}, partyId: {%s}, remote partyId: {%s}, taskActionStatus: {%s}",
		reason, senderRole, proposalId.String(), taskId, sender.GetPartyId(), receiver.GetPartyId(), taskActionStatus.String())
	// Send consensus result to interrupt consensus epoch and clean some data (on task sender)
	if senderRole == libtypes.TaskRole_TaskRole_Sender {
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
	localTaskRole libtypes.TaskRole,
	localTaskOrganization *libtypes.TaskOrganization,
	remoteTaskRole libtypes.TaskRole,
	remoteTaskOrganization *libtypes.TaskOrganization,
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

func (t *Twopc) sendPrepareMsg(proposalId common.Hash, nonConsTask *types.NeedConsensusTask, startTime uint64) error {

	task := nonConsTask.GetTask()
	sender := task.GetTaskSender()

	needSendLocalMsgFn := func() bool {
		for i := 0; i < len(task.GetTaskData().GetDataSuppliers()); i++ {
			if sender.GetIdentityId() == task.GetTaskData().GetDataSuppliers()[i].GetIdentityId() {
				return true
			}
		}

		for i := 0; i < len(task.GetTaskData().GetPowerSuppliers()); i++ {
			if sender.GetIdentityId() == task.GetTaskData().GetPowerSuppliers()[i].GetIdentityId() {
				return true
			}
		}

		for i := 0; i < len(task.GetTaskData().GetReceivers()); i++ {
			if sender.GetIdentityId() == task.GetTaskData().GetReceivers()[i].GetIdentityId() {
				return true
			}
		}
		return false
	}

	prepareMsg, err := makePrepareMsg(proposalId, libtypes.TaskRole_TaskRole_Sender, libtypes.TaskRole_TaskRole_Unknown,
		sender.GetPartyId(), "", nonConsTask, startTime)
	if nil != err {
		return fmt.Errorf("failed to make prepareMsg, proposalId: %s, taskId: %s, err: %s",
			proposalId.String(), task.GetTaskId(), err)
	}

	msg := &types.PrepareMsg{
		MsgOption: types.FetchMsgOption(prepareMsg.GetMsgOption()),
		TaskInfo:  task,
		Evidence:  string(prepareMsg.GetEvidence()),
		CreateAt:  prepareMsg.GetCreateAt(),
	}

	// signature the msg and fill sign field of prepareMsg
	sign, err := crypto.Sign(msg.Hash().Bytes(), t.config.Option.NodePriKey)
	if nil != err {
		return fmt.Errorf("failed to sign prepareMsg, proposalId: %s, err: %s",
			msg.GetMsgOption().GetProposalId().String(), err)
	}
	prepareMsg.Sign = sign

	errs := make([]string, 0)
	if needSendLocalMsgFn() {
		if err := t.sendLocalPrepareMsg("", prepareMsg); nil != err {
			errs = append(errs, fmt.Sprintf("send prepareMsg to local peer, %s", err))
		}
	}

	if err := t.p2p.Broadcast(context.TODO(), prepareMsg); nil != err {
		errs = append(errs, fmt.Sprintf("send prepareMsg to remote peer, %s", err))
	}

	log.WithField("traceId", traceutil.GenerateTraceID(prepareMsg)).Debugf("Succeed to call `sendPrepareMsg` proposalId: %s, taskId: %s, startTime: %d ms",
		proposalId.String(), task.GetTaskId(), startTime)

	if len(errs) != 0 {
		return fmt.Errorf(
			"\n######################################################## \n%s\n########################################################\n",
			strings.Join(errs, "\n"))
	}
	return nil
}


func (t *Twopc) sendPrepareVote(pid peer.ID, sender, receiver *libtypes.TaskOrganization, req *twopcpb.PrepareVote) error {

	vote := &types.PrepareVote{
		MsgOption:  types.FetchMsgOption(req.GetMsgOption()),
		VoteOption: types.VoteOptionFromBytes(req.GetVoteOption()),
		PeerInfo:   types.FetchTaskPeerInfo(req.GetPeerInfo()),
		CreateAt:   req.GetCreateAt(),
	}

	// signature the msg and fill sign field of prepareVote
	sign, err := crypto.Sign(vote.Hash().Bytes(), t.config.Option.NodePriKey)
	if nil != err {
		return fmt.Errorf("failed to sign prepareVote, proposalId: %s, err: %s",
			vote.GetMsgOption().GetProposalId().String(), err)
	}
	req.Sign = sign

	if types.IsNotSameTaskOrg(sender, receiver) {
		//return handler.SendTwoPcPrepareVote(context.TODO(), t.p2p, pid, req)
		return t.p2p.Broadcast(context.TODO(), req)
	} else {
		return t.sendLocalPrepareVote(pid, req)
	}
}

func (t *Twopc) sendConfirmMsg(proposalId common.Hash, task *types.Task, peers *twopcpb.ConfirmTaskPeerInfo, option types.TwopcMsgOption, startTime uint64) error {

	sender := task.GetTaskSender()

	needSendLocalMsgFn := func() bool {
		for i := 0; i < len(task.GetTaskData().GetDataSuppliers()); i++ {
			if sender.GetIdentityId() == task.GetTaskData().GetDataSuppliers()[i].GetIdentityId() {
				return true
			}
		}

		for i := 0; i < len(task.GetTaskData().GetPowerSuppliers()); i++ {
			if sender.GetIdentityId() == task.GetTaskData().GetPowerSuppliers()[i].GetIdentityId() {
				return true
			}
		}

		for i := 0; i < len(task.GetTaskData().GetReceivers()); i++ {
			if sender.GetIdentityId() == task.GetTaskData().GetReceivers()[i].GetIdentityId() {
				return true
			}
		}
		return false
	}

	confirmMsg := makeConfirmMsg(proposalId, libtypes.TaskRole_TaskRole_Sender, libtypes.TaskRole_TaskRole_Unknown,
		sender.GetPartyId(), "", sender, peers, option, startTime)

	msg := &types.ConfirmMsg{
		MsgOption:     types.FetchMsgOption(confirmMsg.GetMsgOption()),
		ConfirmOption: types.TwopcMsgOptionFromBytes(confirmMsg.GetConfirmOption()),
		Peers:         confirmMsg.GetPeers(),
		CreateAt:      confirmMsg.GetCreateAt(),
	}

	// signature the msg and fill sign field of confirmMsg
	sign, err := crypto.Sign(msg.Hash().Bytes(), t.config.Option.NodePriKey)
	if nil != err {
		return fmt.Errorf("failed to sign confirmMsg, proposalId: %s, err: %s",
			msg.GetMsgOption().GetProposalId().String(), err)
	}
	confirmMsg.Sign = sign

	errs := make([]string, 0)
	if needSendLocalMsgFn() {
		if err := t.sendLocalConfirmMsg("", confirmMsg); nil != err {
			errs = append(errs, fmt.Sprintf("send confirmMsg to local peer, %s", err))
		}
	}

	if err := t.p2p.Broadcast(context.TODO(), confirmMsg); nil != err {
		errs = append(errs, fmt.Sprintf("send confirmMsg to remote peer, %s", err))
	}

	log.WithField("traceId", traceutil.GenerateTraceID(confirmMsg)).Debugf("Succeed to call`sendConfirmMsg.%s` proposalId: %s, taskId: %s",
		proposalId.String(), task.GetTaskId())


	if len(errs) != 0 {
		return fmt.Errorf(
			"\n######################################################## \n%s\n########################################################\n",
			strings.Join(errs, "\n"))
	}
	return nil
}

func (t *Twopc) sendConfirmVote(pid peer.ID, sender, receiver *libtypes.TaskOrganization, req *twopcpb.ConfirmVote) error {

	vote := &types.ConfirmVote{
		MsgOption:  types.FetchMsgOption(req.GetMsgOption()),
		VoteOption: types.VoteOptionFromBytes(req.GetVoteOption()),
		CreateAt:   req.GetCreateAt(),
	}

	// signature the msg and fill sign field of confirmVote
	sign, err := crypto.Sign(vote.Hash().Bytes(), t.config.Option.NodePriKey)
	if nil != err {
		return fmt.Errorf("failed to sign confirmVote, proposalId: %s, err: %s",
			vote.GetMsgOption().GetProposalId().String(), err)
	}
	req.Sign = sign

	if types.IsNotSameTaskOrg(sender, receiver) {
		//return handler.SendTwoPcConfirmVote(context.TODO(), t.p2p, pid, req)
		return t.p2p.Broadcast(context.TODO(), req)
	} else {
		return t.sendLocalConfirmVote(pid, req)
	}
}

func (t *Twopc) sendCommitMsg(proposalId common.Hash, task *types.Task, option types.TwopcMsgOption, startTime uint64) error {

	sender := task.GetTaskSender()


	needSendLocalMsgFn := func() bool {
		for i := 0; i < len(task.GetTaskData().GetDataSuppliers()); i++ {
			if sender.GetIdentityId() == task.GetTaskData().GetDataSuppliers()[i].GetIdentityId() {
				return true
			}
		}

		for i := 0; i < len(task.GetTaskData().GetPowerSuppliers()); i++ {
			if sender.GetIdentityId() == task.GetTaskData().GetPowerSuppliers()[i].GetIdentityId() {
				return true
			}
		}

		for i := 0; i < len(task.GetTaskData().GetReceivers()); i++ {
			if sender.GetIdentityId() == task.GetTaskData().GetReceivers()[i].GetIdentityId() {
				return true
			}
		}
		return false
	}

	commitMsg := makeCommitMsg(proposalId, libtypes.TaskRole_TaskRole_Sender, libtypes.TaskRole_TaskRole_Unknown,
		sender.GetPartyId(), "", sender, option, startTime)

	msg := &types.CommitMsg{
		MsgOption:     types.FetchMsgOption(commitMsg.GetMsgOption()),
		CommitOption: types.TwopcMsgOptionFromBytes(commitMsg.GetCommitOption()),
		CreateAt:      commitMsg.GetCreateAt(),
	}

	// signature the msg and fill sign field of commitMsg
	sign, err := crypto.Sign(msg.Hash().Bytes(), t.config.Option.NodePriKey)
	if nil != err {
		return fmt.Errorf("failed to sign commitMsg, proposalId: %s, err: %s",
			msg.GetMsgOption().GetProposalId().String(), err)
	}
	commitMsg.Sign = sign


	errs := make([]string, 0)
	if needSendLocalMsgFn() {
		if err := t.sendLocalCommitMsg("", commitMsg); nil != err {
			errs = append(errs, fmt.Sprintf("send commitMsg to local peer, %s", err))
		}
	}

	if err := t.p2p.Broadcast(context.TODO(), commitMsg); nil != err {
		errs = append(errs, fmt.Sprintf("send commitMsg to remote peer, %s", err))
	}

	log.WithField("traceId", traceutil.GenerateTraceID(commitMsg)).Debugf("Succeed to call`sendCommitMsg.%s` proposalId: %s, taskId: %s",
		proposalId.String(), task.GetTaskId())


	if len(errs) != 0 {
		return fmt.Errorf(
			"\n######################################################## \n%s\n########################################################\n",
			strings.Join(errs, "\n"))
	}

	return nil
}

func verifyPartyRole(partyId string, role libtypes.TaskRole, task *types.Task) bool {

	switch role {
	case libtypes.TaskRole_TaskRole_Sender:
		if partyId == task.GetTaskSender().GetPartyId() {
			return true
		}
	case libtypes.TaskRole_TaskRole_DataSupplier:
		for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
			if partyId == dataSupplier.GetPartyId() {
				return true
			}
		}
	case libtypes.TaskRole_TaskRole_PowerSupplier:
		for _, powerSupplier := range task.GetTaskData().GetPowerSuppliers() {
			if partyId == powerSupplier.GetPartyId() {
				return true
			}
		}
	case libtypes.TaskRole_TaskRole_Receiver:
		for _, receiver := range task.GetTaskData().GetReceivers() {
			if partyId == receiver.GetPartyId() {
				return true
			}
		}
	}
	return false
}

func fetchOrgByPartyRole(partyId string, role libtypes.TaskRole, task *types.Task) *libtypes.TaskOrganization {

	switch role {
	case libtypes.TaskRole_TaskRole_Sender:
		if partyId == task.GetTaskSender().GetPartyId() {
			return task.GetTaskSender()
		}
	case libtypes.TaskRole_TaskRole_DataSupplier:
		for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
			if partyId == dataSupplier.GetPartyId() {
				return dataSupplier
			}
		}
	case libtypes.TaskRole_TaskRole_PowerSupplier:
		for _, powerSupplier := range task.GetTaskData().GetPowerSuppliers() {
			if partyId == powerSupplier.GetPartyId() {
				return powerSupplier
			}
		}
	case libtypes.TaskRole_TaskRole_Receiver:
		for _, receiver := range task.GetTaskData().GetReceivers() {
			if partyId == receiver.GetPartyId() {
				return receiver
			}
		}
	}
	return nil
}

func (t *Twopc) verifyPartyAndTaskPartner(role libtypes.TaskRole, party *libtypes.TaskOrganization, task *types.Task) (bool, error) {
	var identityValid bool
	switch role {
	case libtypes.TaskRole_TaskRole_DataSupplier:

		for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
			// identity + partyId
			if dataSupplier.GetIdentityId() == party.GetIdentityId() && dataSupplier.GetPartyId() == party.GetPartyId() {
				identityValid = true
				break
			}
		}
	case libtypes.TaskRole_TaskRole_PowerSupplier:

		for _, powerSupplier := range task.GetTaskData().GetPowerSuppliers() {
			// identity + partyId
			if powerSupplier.GetIdentityId() == party.GetIdentityId() && powerSupplier.GetPartyId() == party.GetPartyId() {
				identityValid = true
				break
			}
		}
	case libtypes.TaskRole_TaskRole_Receiver:

		for _, receiver := range task.GetTaskData().GetReceivers() {
			// identity + partyId
			if receiver.GetIdentityId() == party.GetIdentityId() && receiver.GetPartyId() == party.GetPartyId() {
				identityValid = true
				break
			}
		}
	default:
		return false, fmt.Errorf("invalid party, the party is not task partners [taskId: %s, role: %s, identity: %s, partyId: %s]",
			task.GetTaskData().GetTaskId(), role.String(), party.GetIdentityId(), party.GetPartyId())

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
				proposalTaskPB := &libtypes.ProposalTask{}
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

				libOrgProposalState := &libtypes.OrgProposalState{}
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
