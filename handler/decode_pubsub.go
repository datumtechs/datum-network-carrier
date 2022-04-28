package handler

import (
	"github.com/Metisnetwork/Metis-Carrier/p2p"
	"github.com/gogo/protobuf/proto"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/pkg/errors"
	"strings"
)

var errNilPubsubMessage = errors.New("nil pubsub message")
var errInvalidTopic = errors.New("invalid topic format")

func (s *Service) decodePubsubMessage(msg *pubsub.Message) (proto.Message, error) {
	if msg == nil || msg.Topic == nil || *msg.Topic == "" {
		return nil, errNilPubsubMessage
	}
	topic := *msg.Topic
	topic = strings.TrimSuffix(topic, s.cfg.P2P.Encoding().ProtocolSuffix())
	topic, err := s.replaceForkDigest(topic)
	if err != nil {
		return nil, err
	}
	base, ok := p2p.GossipTopicMappings[topic]
	if !ok {
		return nil, p2p.ErrMessageNotMapped
	}
	m := proto.Clone(base)
	if err := s.cfg.P2P.Encoding().DecodeGossip(msg.Data, m); err != nil {
		return nil, err
	}
	return m, nil
}

// Replaces our fork digest with the formatter.
func (s *Service) replaceForkDigest(topic string) (string, error) {
	subStrings := strings.Split(topic, "/")
	if len(subStrings) != 4 {
		return "", errInvalidTopic
	}
	subStrings[2] = "%x"
	return strings.Join(subStrings, "/"), nil
}