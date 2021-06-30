package types

import (
	"sync"
)


type peer struct {
	conn *conn
}


func (p *peer) SendMsg() error {return nil}
func (p *peer) RecvMsg() error {return nil}


// PeerSet represents the collection of active peers currently participating
// in the TwoPC protocol.
type PeerSet struct {
	peers  map[string]*peer
	lock   sync.RWMutex
	closed bool
}




