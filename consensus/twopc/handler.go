package twopc

import lru "github.com/hashicorp/golang-lru"

// EngineManager responsibles for processing the messages in the network.
type EngineManager struct {
	engine             *twoPC
	peers              *PeerSet
	//sendQueue          chan *types.MsgPackage
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
