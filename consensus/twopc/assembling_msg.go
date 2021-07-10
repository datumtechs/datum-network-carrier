package twopc

import (
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/types"
)


func makePrepareMsg(startTime uint64) *pb.PrepareMsg {

	// TODO 需要将 ScheduleTask 转换成  PrepareMsg
	return nil
}

func makePrepareVote() *pb.PrepareVote {

	// TODO 组装  PrepareVote
	return nil
}

func makeConfirmMsg(epoch, startTime uint64) *pb.ConfirmMsg {

	// TODO 组装  ConfirmMsg
	return nil
}

func makeConfirmVote() *pb.ConfirmVote {

	// TODO 组装  ConfirmVote
	return nil
}

func makeCommitMsg(startTime uint64) *pb.CommitMsg {

	// TODO 组装  CommitMsg
	return nil
}

func makeTaskResultMsg(startTime uint64) *pb.TaskResultMsg {

	// TODO 组装  TaskResultMsg
	return nil
}

func fetchPrepareMsg(prepareMsg *types.PrepareMsgWrap) (*types.ProposalTask, error) {
	// TODO 需要实现
	return nil, nil
}
func fetchPrepareVote(prepareVote *types.PrepareVoteWrap) (*types.PrepareVote, error) {

	// TODO 需要实现
	return nil, nil
}

func fetchConfirmMsg (confirmMsg *types.ConfirmMsgWrap) (*types.ConfirmMsg, error) {


	// TODO 需要实现
	return nil, nil
}

func fetchConfirmVote (confirmVote *types.ConfirmVoteWrap) (*types.ConfirmVote, error) {

	// TODO 需要实现
	return nil, nil
}

func fetchCommitMsg (commitMsg *types.CommitMsgWrap) (*types.CommitMsg, error) {


	// TODO 需要实现
	return nil, nil
}

func fetchTaskResultMsg (commitMsg *types.TaskResultMsgWrap) (*types.TaskResultMsg, error) {


	// TODO 需要实现
	return nil, nil
}