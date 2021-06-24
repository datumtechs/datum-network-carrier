package types

import "github.com/RosettaFlow/Carrier-Go/consensus/twopc"

const (
	NoneMode = iota // none consensus node
	PartMode        // partial node
	FullMode        // all node
)

// MsgPackage represents a specific message package.
// It contains the node ID, the message body, and
// the forwarding mode from the sender.
type MsgPackage struct {
	peerID string        // from the sender of the message
	msg    twopc.Message // message body
	mode   uint64        // forwarding mode.
}

// Create a new MsgPackage based on params.
func NewMsgPackage(pid string, msg twopc.Message, mode uint64) *MsgPackage {
	return &MsgPackage{
		peerID: pid,
		msg:    msg,
		mode:   mode,
	}
}

func (m *MsgPackage) Message() twopc.Message {
	return m.msg
}

func (m *MsgPackage) PeerID() string {
	return m.peerID
}

func (m *MsgPackage) Mode() uint64 {
	return m.mode
}
