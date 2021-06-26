package types

import pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"

type PrepareMsgWrap struct {
	*pb.PrepareMsg
}
func (msg *PrepareMsgWrap) String() string {return ""}

type PrepareVoteWrap struct {
	*pb.PrepareVote
}
func (msg *PrepareVoteWrap) String() string {return ""}

type ConfirmMsgWrap struct {
	*pb.ConfirmMsg
}
func (msg *ConfirmMsgWrap) String() string {return ""}

type CommitMsgWrap struct {
	*pb.CommitMsg
}
func (msg *CommitMsgWrap) String() string {return ""}