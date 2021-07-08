package twopc

import (
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/types"
)


func makePrepareMsg() *pb.PrepareMsg {

	// TODO 需要将 ScheduleTask 转换成  PrepareMsg
	return nil
}

func makePrepareVote() *pb.PrepareVote {

	// TODO 组装  PrepareVote
	return nil
}

func makeConfirmMsg() *pb.ConfirmMsg {

	// TODO 组装  ConfirmMsg
	return nil
}

func makeConfirmVote() *pb.ConfirmVote {

	// TODO 组装  ConfirmVote
	return nil
}

func makeCommitMsg() *pb.CommitMsg {

	// TODO 组装  CommitMsg
	return nil
}


func (t *TwoPC) fetchPrepareMsg(prepareMsg *types.PrepareMsgWrap) (*types.ProposalTask, error) {
	// TODO 需要实现
	return nil, nil
}
func (t *TwoPC) fetchPrepareVote(prepareVote *types.PrepareVoteWrap) (*types.PrepareVote, error) {

	// TODO 需要实现
	return nil, nil
}