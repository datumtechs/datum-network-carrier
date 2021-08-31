package p2p

import (
	"context"
	"crypto/ecdsa"
	pb "github.com/RosettaFlow/Carrier-Go/lib/p2p/v1"
	"github.com/RosettaFlow/Carrier-Go/p2p/encoder"
	"github.com/RosettaFlow/Carrier-Go/p2p/peers"
	"github.com/ethereum/go-ethereum/p2p/enr"
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/connmgr"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/multiformats/go-multiaddr"
)

// P2P represents the full p2p interface composed of all of the sub-interfaces.
type P2P interface {
	Broadcaster
	SetStreamHandler
	EncodingProvider
	PubSubProvider
	PubSubTopicUser
	PeerManager
	Sender
	ConnectionHandler
	PeersProvider
	MetadataProvider
}

// Broadcaster broadcasts messages to peers over the p2p pubsub protocol.
type Broadcaster interface {
	Broadcast(ctx context.Context, message proto.Message) error
	//BroadcastTask(ctx context.Context, task *libtypes.TaskData) error
}

// SetStreamHandler configures p2p to handle streams of a certain topic ID.
type SetStreamHandler interface {
	SetStreamHandler(topic string, handler network.StreamHandler)
}

// PubSubTopicUser providers way to join, use and leave PubSub topics.
type PubSubTopicUser interface {
	JoinTopic(topic string, opts ...pubsub.TopicOpt) (*pubsub.Topic, error)
	LeaveTopic(topic string) error
	PublishToTopic(ctx context.Context, topic string, data []byte, opts ...pubsub.PubOpt) error
	SubscribeToTopic(topic string, opts ...pubsub.SubOpt) (*pubsub.Subscription, error)
}

// ConnectionHandler configures p2p to handle connection with a peer.
type ConnectionHandler interface {
	AddConnectionHandler(f func(ctx context.Context, id peer.ID) error, j func(ctx context.Context, id peer.ID) error)
	AddDisconnectionHandler(f func(ctx context.Context, id peer.ID) error)
	connmgr.ConnectionGater
}

// EncodingProvider provides p2p network encoding.
type EncodingProvider interface {
	Encoding() encoder.NetworkEncoding
}

// PubSubProvider provides the p2p pubsub protocol.
type PubSubProvider interface {
	PubSub() *pubsub.PubSub
}

// PeerManager abstracts some peer management methods from libp2p
type PeerManager interface {
	AddPeer(string) error
	Disconnect(peer.ID) error
	PeerID() peer.ID
	NodeId() string
	PirKey() *ecdsa.PrivateKey
	Host() host.Host
	ENR() *enr.Record
	DiscoveryAddresses() ([]multiaddr.Multiaddr, error)
	RefreshENR()
	FindPeersWithSubnet(ctx context.Context, topic string, index, threshold uint64) (bool, error)
	AddPingMethod(reqFunc func(ctx context.Context, id peer.ID) error)
}

// Sender abstracts the sending functionality from libp2p.
type Sender interface {
	Send(context.Context, interface{}, string, peer.ID) (network.Stream, error)
}

//
type PeersProvider interface {
	Peers() *peers.Status
}

// MetadataProvider returns the metadata related information for the local peer.
type MetadataProvider interface {
	Metadata() *pb.MetaData
	MetadataSeq() uint64
}
