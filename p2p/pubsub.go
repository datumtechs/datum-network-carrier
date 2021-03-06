package p2p

import (
	"context"
	"github.com/RosettaFlow/Carrier-Go/common/hashutil"
	pbrpc "github.com/RosettaFlow/Carrier-Go/lib/rpc/v1"
	"github.com/RosettaFlow/Carrier-Go/p2p/encoder"
	"github.com/RosettaFlow/Carrier-Go/params"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	pubsub_pb "github.com/libp2p/go-libp2p-pubsub/pb"
	"time"
)

const (
	// overlay parameters
	gossipSubD   = 8  // topic stable mesh target count
	gossipSubDlo = 6  // topic stable mesh low watermark
	gossipSubDhi = 12 // topic stable mesh high watermark

	// gossip parameters
	gossipSubMcacheLen    = 6   // number of windows to retain full message in  cache for `IWANT` responses
	gossipSubMcacheGossip = 3   // number of windows to gossip about
	gossipSubSeenTTL      = 550 // number of heartbeat intervals to retain message IDs

	// fanout ttl: TTL for fanout maps for topics
	// we are not subscribed to but have published to,
	// in nano seconds
	gossipSubFanoutTTL = 60000000000

	// heartbeat interval
	gossipSubHeartbeatInterval = 700 * time.Millisecond // frequency of heartbeat, milliseconds

	// misc
	randomSubD = 6 // random gossip target
)

// JoinTopic will join PubSub topic, if not already joined
func (s *Service) JoinTopic(topic string, opts ...pubsub.TopicOpt) (*pubsub.Topic, error) {
	s.joinedTopicsLock.Lock()
	defer s.joinedTopicsLock.Unlock()

	if _, ok := s.joinedTopics[topic]; !ok {
		topicHandle, err := s.pubsub.Join(topic, opts...)
		if err != nil {
			return nil, err
		}
		s.joinedTopics[topic] = topicHandle
	}
	return s.joinedTopics[topic], nil
}

// LeaveTopic closes topic and removes corresponding handler from list of joined topics.
// This method will return error if there are outstanding evengine handlers or subscriptions.
func (s *Service) LeaveTopic(topic string) error {
	s.joinedTopicsLock.Lock()
	defer s.joinedTopicsLock.Unlock()

	if t, ok := s.joinedTopics[topic]; ok{
		if err := t.Close(); err != nil {
			return err
		}
		delete(s.joinedTopics, topic)
	}
	return nil
}

// PublishToTopic joins (if necessary) and publishes a message to a PubSub topic.
func (s *Service) PublishToTopic(ctx context.Context, topic string, data []byte, opts ...pubsub.PubOpt) error {
	topicHandle, err := s.JoinTopic(topic)
	if err != nil {
		return err
	}
	// Wait for at least 1 peer to be available to receive the published message.
	for {
		if len(topicHandle.ListPeers()) > 0 {
			return topicHandle.Publish(ctx, data, opts...)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// SubscribeToTopic joins (if necessary) and subscribes to PubSub topic.
func (s *Service) SubscribeToTopic(topic string, opts ...pubsub.SubOpt) (*pubsub.Subscription, error) {
	s.awaitStateInitialized()

	topicHandle, err := s.JoinTopic(topic)
	if err != nil {
		return nil, err
	}
	//TODO: add new topic score params...
	scoringParams, err := s.topicScoreParams(topic)
	if err != nil {
		return nil, err
	}
	if scoringParams != nil {
		if err := topicHandle.SetScoreParams(scoringParams); err != nil {
			return nil, err
		}
		logGossipParameters(topic, scoringParams)
	}
	return topicHandle.Subscribe(opts...)
}

// peerInspector will scrape all the relevant scoring data
// and add it to our peer handler.
func (s *Service) peerInspector(peerMap map[peer.ID]*pubsub.PeerScoreSnapshot) {
	// Iterate through all the connected peers and through any of their
	// relevant topics.
	for pid, snap := range peerMap {
		s.peers.Scorers().GossipScorer().SetGossipData(pid, snap.Score,
			snap.BehaviourPenalty, convertTopicScores(snap.Topics))
	}
}

// Content addressable ID function.
//
// ETH2 spec defines the message ID as:
//    The `message-id` of a gossipsub message MUST be the following 20 byte value computed from the message data:
//    If `message.data` has a valid snappy decompression, set `message-id` to the first 20 bytes of the `SHA256` hash of
//    the concatenation of `MESSAGE_DOMAIN_VALID_SNAPPY` with the snappy decompressed message data,
//    i.e. `SHA256(MESSAGE_DOMAIN_VALID_SNAPPY + snappy_decompress(message.data))[:20]`.
//
//    Otherwise, set `message-id` to the first 20 bytes of the `SHA256` hash of
//    the concatenation of `MESSAGE_DOMAIN_INVALID_SNAPPY` with the raw message data,
//    i.e. `SHA256(MESSAGE_DOMAIN_INVALID_SNAPPY + message.data)[:20]`.
func msgIDFunction(pmsg *pubsub_pb.Message) string {
	decodedData, err := encoder.DecodeSnappy(pmsg.Data, params.CarrierNetworkConfig().GossipMaxSize)
	if err != nil {
		combinedData := append(params.CarrierNetworkConfig().MessageDomainInvalidSnappy[:], pmsg.Data...)
		h := hashutil.Hash(combinedData)
		return string(h[:20])
	}
	combinedData := append(params.CarrierNetworkConfig().MessageDomainValidSnappy[:], decodedData...)
	h := hashutil.Hash(combinedData)
	return string(h[:20])
}

func setPubSubParameters() {
	pubsub.GossipSubDlo = gossipSubDlo
	pubsub.GossipSubD = gossipSubD
	pubsub.GossipSubHeartbeatInterval = gossipSubHeartbeatInterval
	pubsub.GossipSubHistoryLength = gossipSubMcacheLen
	pubsub.GossipSubHistoryGossip = gossipSubMcacheGossip
	pubsub.TimeCacheDuration = 550 * gossipSubHeartbeatInterval

	// Set a larger gossip history to ensure that slower
	// messages have a longer time to be propagated. This
	// comes with the tradeoff of larger memory usage and
	// size of the seen message cache.
	//TODO: add config....
	//if featureconfig.Get().EnableLargerGossipHistory {
	//	pubsub.GossipSubHistoryLength = 12
	//	pubsub.GossipSubHistoryLength = 5
	//}
}

// convert from libp2p's internal schema to a compatible protobuf format.
func convertTopicScores(topicMap map[string]*pubsub.TopicScoreSnapshot) map[string]*pbrpc.TopicScoreSnapshot {
	newMap := make(map[string]*pbrpc.TopicScoreSnapshot, len(topicMap))
	for t, s := range topicMap {
		newMap[t] = &pbrpc.TopicScoreSnapshot{
			TimeInMesh:               uint64(s.TimeInMesh.Milliseconds()),
			FirstMessageDeliveries:   float32(s.FirstMessageDeliveries),
			MeshMessageDeliveries:    float32(s.MeshMessageDeliveries),
			InvalidMessageDeliveries: float32(s.InvalidMessageDeliveries),
		}
	}
	return newMap
}
