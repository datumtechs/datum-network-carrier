package peerdata

import (
	"context"
	"errors"
	pb "github.com/RosettaFlow/Carrier-Go/lib/p2p/v1"
	pbrpc "github.com/RosettaFlow/Carrier-Go/lib/rpc/v1"
	"github.com/ethereum/go-ethereum/p2p/enr"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1"
	"sync"
	"time"
)

var (
	// ErrPeerKnow is returned when there is an attempt to abtain data from a peer that is not know.
	ErrPeerUnknown = errors.New("peer unknown")
)

// PeerConnectionState is the state of the connection
type PeerConnectionState ethpb.ConnectionState

// StoreConfig holds peer store parameters.
type StoreConfig struct {
	MaxPeers int
}

// Store is a container for various peer related data (both protocol and app level).
// Container implements RWMutex, so data access can be restricted on the container level.
// This allows different components rely on the very same peer map container.
//
// Note: access to data is controlled by clients i.e. client code is responsible for locking/unlocking
// the mutex when accessing data.
type Store struct {
	sync.RWMutex
	ctx    context.Context
	config *StoreConfig
	peers  map[peer.ID]*PeerData
}

// PeerData aggregates protocol and application level info about a single peer.
type PeerData struct {
	Address       ma.Multiaddr
	Direction     network.Direction
	ConnState     PeerConnectionState
	Enr           *enr.Record
	NextValidTime time.Time
	// Chain related data
	MetaData                  *pb.MetaData
	ChainState                *pb.Status
	ChainStateLastUpdated     time.Time
	ChainStateValidationError error
	// Scores internal data.
	BadResponses         int
	ProcessedBlocks      uint64
	BlockProviderUpdated time.Time
	// Gossip Scoring data.
	TopicScores map[string]*pbrpc.TopicScoreSnapshot
	GossipScore      float64
	BehaviourPenalty float64
}

// NewStore creates new peer data store.
func NewStore(ctx context.Context, config *StoreConfig) *Store {
	return &Store{
		ctx:    ctx,
		config: config,
		peers:  make(map[peer.ID]*PeerData),
	}
}

// PeerData returns data associated with a given peer, if any.
// Notes: it is assumed that sotre mutex is locked when calling this method.
func (s *Store) PeerData(pid peer.ID) (*PeerData, bool) {
	peerData, ok := s.peers[pid]
	return peerData, ok
}

// PeerDataGetOrCreate returns data associated with a given peer.
// If no data has been associated yet, newly created and associated data object is returned.
// Notes: it is assumed that store mutex is locked when calling this method.
func (s *Store) PeerDataGetOrCreate(pid peer.ID) *PeerData {
	if peerData, ok := s.peers[pid]; ok {
		return peerData
	}
	s.peers[pid] = &PeerData{}
	return s.peers[pid]
}

// SetPeerData updates data associated with a given peer.
// Notes: it is assumed that score mutex is locked when calling this method.
func (s *Store) SetPeerData(pid peer.ID, data *PeerData) {
	s.peers[pid] = data
}

// DeletePeerData removes data associated with a given peer.
// Notes: it is assumed that store mutex is locked when calling this method.
func (s *Store) DeletePeerData(pid peer.ID) {
	delete(s.peers, pid)
}

// Peers returns map of peer data objects.
// Notes: it is assumed that store mutex is locked when calling this method.
func (s *Store) Peers() map[peer.ID]*PeerData {
	return s.peers
}

// Config exposes store configuration params.
func (s *Store) Config() *StoreConfig {
	return s.config
}
