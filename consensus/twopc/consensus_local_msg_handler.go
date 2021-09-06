package twopc

import (
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/libp2p/go-libp2p-core/peer"
)

func (t *TwoPC) sendLocalPrepareMsg(pid peer.ID, req *pb.PrepareMsg) error {
	return t.onPrepareMsg(pid, &types.PrepareMsgWrap{PrepareMsg: req})
}

func (t *TwoPC) sendLocalPrepareVote(pid peer.ID, req *pb.PrepareVote) error {
	return t.onPrepareVote(pid, &types.PrepareVoteWrap{PrepareVote: req})
}

func (t *TwoPC) sendLocalConfirmMsg(pid peer.ID, req *pb.ConfirmMsg) error {
	return t.onConfirmMsg(pid, &types.ConfirmMsgWrap{ConfirmMsg: req})
}

func (t *TwoPC) sendLocalConfirmVote(pid peer.ID, req *pb.ConfirmVote) error {
	return t.onConfirmVote(pid, &types.ConfirmVoteWrap{ConfirmVote: req})
}

func (t *TwoPC) sendLocalCommitMsg(pid peer.ID, req *pb.CommitMsg) error {
	return t.onCommitMsg(pid, &types.CommitMsgWrap{CommitMsg: req})
}

//func (t *TwoPC) sendLocalTaskResultMsg(pid peer.ID, req *pb.TaskResultMsg) error {
//	return t.onTaskResultMsg(pid, &types.TaskResultMsgWrap{TaskResultMsg: req})
//}