package types

import (
	"github.com/datumtechs/datum-network-carrier/p2p"
	"github.com/libp2p/go-libp2p-core/peer"
	"sync"
)


type Peer struct {
	conn *conn
	nodeId p2p.NodeID
	pid peer.ID
}


func (p *Peer) SendMsg() error {return nil}
func (p *Peer) RecvMsg() error {return nil}
func (p *Peer) GetPid() peer.ID {return p.pid}
func (p *Peer) GetNodeId() p2p.NodeID {return p.nodeId}


// PeerSet represents the collection of active peers currently participating
// in the TwoPC protocol.
type PeerSet struct {
	peers  map[p2p.NodeID]*Peer
	lock   sync.RWMutex
	closed bool
}

func  NewPeerSet(size int) *PeerSet {
	return &PeerSet{
		peers: make(map[p2p.NodeID]*Peer, size),
	}
}

func (ps *PeerSet) GetPeer(nodeId p2p.NodeID) *Peer {
	return ps.peers[nodeId]
}
func (ps *PeerSet)SetPeer(peer *Peer) {
	ps.peers[peer.nodeId] = peer
}


