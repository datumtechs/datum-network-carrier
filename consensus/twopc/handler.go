package twopc

import (
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	lru "github.com/hashicorp/golang-lru"
)

// EngineManager responsibles for processing the messages in the network.
type EngineManager struct {
	engine             *TwoPC
	peers              *PeerSet
	sendQueue          chan *types.MsgPackage
	//quitSend           chan struct{}
	//sendQueueHook      func(*types.MsgPackage)
	historyMessageHash *lru.ARCCache // Consensus message record that has been processed successfully.
	blacklist          *lru.Cache    // Save node blacklist.
}

func (e *EngineManager) handleMsg(p *peer) error {

	switch {

	// TODO 处理 Msg.Code
	
	}
	
	return nil
}

// Broadcast imports messages into the send queue and send it according to broadcast.
//
// Note: The broadcast of this method defaults to FULL mode.
func (h *EngineManager) Broadcast(msg Message) {
	msgPkg := types.NewMsgPackage("", msg, types.PartMode)
	select {
	case h.sendQueue <- msgPkg:
		log.Trace("Broadcast message to sendQueue", "msgHash", msg.MsgHash(), "msg", msg.String())
	default:
		log.Error("Broadcast message failed, message queue blocking", "msgHash", msg.MsgHash())
	}
}


// Message interface, all message structures must
// implement this interface.
type Message interface {
	String() string
	MsgHash() common.Hash
}