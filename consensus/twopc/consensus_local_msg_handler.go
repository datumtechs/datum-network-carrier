package twopc

import (
	twopcpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/consensus/twopc"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/libp2p/go-libp2p-core/peer"
)

func (t *Twopc) sendLocalPrepareMsg(pid peer.ID, req *twopcpb.PrepareMsg) error {
	return t.onPrepareMsg(pid, &types.PrepareMsgWrap{PrepareMsg: req}, LocalConsensusMsg)
}

func (t *Twopc) sendLocalPrepareVote(pid peer.ID, req *twopcpb.PrepareVote) error {
	return t.onPrepareVote(pid, &types.PrepareVoteWrap{PrepareVote: req}, LocalConsensusMsg)
}

func (t *Twopc) sendLocalConfirmMsg(pid peer.ID, req *twopcpb.ConfirmMsg) error {
	return t.onConfirmMsg(pid, &types.ConfirmMsgWrap{ConfirmMsg: req}, LocalConsensusMsg)
}

func (t *Twopc) sendLocalConfirmVote(pid peer.ID, req *twopcpb.ConfirmVote) error {
	return t.onConfirmVote(pid, &types.ConfirmVoteWrap{ConfirmVote: req}, LocalConsensusMsg)
}

func (t *Twopc) sendLocalCommitMsg(pid peer.ID, req *twopcpb.CommitMsg) error {
	return t.onCommitMsg(pid, &types.CommitMsgWrap{CommitMsg: req}, LocalConsensusMsg)
}