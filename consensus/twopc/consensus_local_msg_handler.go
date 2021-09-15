package twopc

import (
	twopcpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/libp2p/go-libp2p-core/peer"
)

func (t *TwoPC) sendLocalPrepareMsg(pid peer.ID, req *twopcpb.PrepareMsg) error {
	return t.onPrepareMsg(pid, &types.PrepareMsgWrap{PrepareMsg: req})
}

func (t *TwoPC) sendLocalPrepareVote(pid peer.ID, req *twopcpb.PrepareVote) error {
	return t.onPrepareVote(pid, &types.PrepareVoteWrap{PrepareVote: req})
}

func (t *TwoPC) sendLocalConfirmMsg(pid peer.ID, req *twopcpb.ConfirmMsg) error {
	return t.onConfirmMsg(pid, &types.ConfirmMsgWrap{ConfirmMsg: req})
}

func (t *TwoPC) sendLocalConfirmVote(pid peer.ID, req *twopcpb.ConfirmVote) error {
	return t.onConfirmVote(pid, &types.ConfirmVoteWrap{ConfirmVote: req})
}

func (t *TwoPC) sendLocalCommitMsg(pid peer.ID, req *twopcpb.CommitMsg) error {
	return t.onCommitMsg(pid, &types.CommitMsgWrap{CommitMsg: req})
}

//func (t *Twopc) sendLocalTaskResultMsg(pid peer.ID, req *twopcpb.TaskResultMsg) error {
//	return t.onTaskResultMsg(pid, &types.TaskResultMsgWrap{TaskResultMsg: req})
//}