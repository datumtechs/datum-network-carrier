package handler

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/messagehandler"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/RosettaFlow/Carrier-Go/params"
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"go.opencensus.io/trace"
	"runtime/debug"
	"strings"
	"time"
)

const pubsubMessageTimeout = 30 * time.Second

// subHandler represents handler for a given subscription.
type subHandler func(context.Context, proto.Message) error

// noopValidator is a no-op that only decodes the message, but does not check its contents.
func (s *Service) noopValidator(_ context.Context, _ peer.ID, msg *pubsub.Message) pubsub.ValidationResult {
	m, err := s.decodePubsubMessage(msg)
	if err != nil {
		log.WithError(err).Debug("Could not decode message")
		return pubsub.ValidationReject
	}
	msg.ValidatorData = m
	return pubsub.ValidationAccept
}

// Register PubSub subscribers
func (s *Service) registerSubscribers() {
	s.subscribe(
		p2p.GossipTestDataTopicFormat,
		s.validateGossipTestData,
		s.gossipTestDataSubscriber)
	//TODO: more subscribe to be register...
}

// subscribe to a given topic with a given validator and subscription handler.
// The base protobuf message is used to initialize new messages for decoding.
func (s *Service) subscribe(topic string, validator pubsub.ValidatorEx, handle subHandler) *pubsub.Subscription {
	base := p2p.GossipTopicMappings[topic]
	if base == nil {
		panic(fmt.Sprintf("%s is not mapped to any message in GossipTopicMappings", topic))
	}
	return s.subscribeWithBase(s.addDigestToTopic(topic), validator, handle)
}

func (s *Service) subscribeWithBase(topic string, validator pubsub.ValidatorEx, handle subHandler) *pubsub.Subscription {
	topic += s.cfg.P2P.Encoding().ProtocolSuffix()
	log := log.WithField("topic", topic)

	if err := s.cfg.P2P.PubSub().RegisterTopicValidator(s.wrapAndReportValidation(topic, validator)); err != nil {
		log.WithError(err).Error("Could not register validator for topic")
		return nil
	}

	sub, err := s.cfg.P2P.SubscribeToTopic(topic)
	if err != nil {
		// Any error subscribing to a PubSub topic would be the result of a misconfiguration of
		// libp2p PubSub library or a subscription request to a topic that fails to match the topic
		// subscription filter.
		log.WithError(err).Error("Could not subscribe topic")
		return nil
	}

	// Pipeline decodes the incoming subscription data, runs the validation, and handles the message.
	pipeline := func(msg *pubsub.Message) {
		ctx, cancel := context.WithTimeout(s.ctx, pubsubMessageTimeout)
		defer cancel()
		ctx, span := trace.StartSpan(ctx, "sync.pubsub")
		defer span.End()

		defer func() {
			if r := recover(); r != nil {
				log.WithField("error", r).Error("Panic occurred")
				debug.PrintStack()
			}
		}()

		span.AddAttributes(trace.StringAttribute("topic", topic))

		if msg.ValidatorData == nil {
			log.Debug("Received nil message on pubsub")
			messageFailedProcessingCounter.WithLabelValues(topic).Inc()
			return
		}

		if err := handle(ctx, msg.ValidatorData.(proto.Message)); err != nil {
			log.WithError(err).Debug("Could not handle p2p pubsub")
			messageFailedProcessingCounter.WithLabelValues(topic).Inc()
			return
		}
	}

	// The main message loop for receiving incoming messages from this subscription.
	messageLoop := func() {
		for {
			msg, err := sub.Next(s.ctx)
			if err != nil {
				// This should only happen when the context is cancelled or subscription is cancelled.
				if err != pubsub.ErrSubscriptionCancelled { // Only log a warning on unexpected errors.
					log.WithError(err).Warn("Subscription next failed")
				}
				// Cancel subscription in the evengine of an error, as we are
				// now exiting topic evengine loop.
				sub.Cancel()
				return
			}

			if msg.ReceivedFrom == s.cfg.P2P.PeerID() {
				continue
			}

			go pipeline(msg)
		}
	}

	go messageLoop()
	return sub
}

// Wrap the pubsub validator with a metric monitoring function. This function increments the
// appropriate counter if the particular message fails to validate.
func (s *Service) wrapAndReportValidation(topic string, v pubsub.ValidatorEx) (string, pubsub.ValidatorEx) {
	return topic, func(ctx context.Context, pid peer.ID, msg *pubsub.Message) (res pubsub.ValidationResult) {
		defer messagehandler.HandlePanic(ctx, msg)
		res = pubsub.ValidationIgnore // Default: ignore any message that panics.
		ctx, cancel := context.WithTimeout(ctx, pubsubMessageTimeout)
		defer cancel()
		messageReceivedCounter.WithLabelValues(topic).Inc()
		if msg.Topic == nil {
			messageFailedValidationCounter.WithLabelValues(topic).Inc()
			return pubsub.ValidationReject
		}
		b := v(ctx, pid, msg)
		if b == pubsub.ValidationReject {
			messageFailedValidationCounter.WithLabelValues(topic).Inc()
		}
		return b
	}
}

// revalidate that our currently connected subnets are valid.
func (s *Service) reValidateSubscriptions(subscriptions map[uint64]*pubsub.Subscription,
	wantedSubs []uint64, topicFormat string, digest [4]byte) {
	for k, v := range subscriptions {
		var wanted bool
		for _, idx := range wantedSubs {
			if k == idx {
				wanted = true
				break
			}
		}
		if !wanted && v != nil {
			v.Cancel()
			fullTopic := fmt.Sprintf(topicFormat, digest, k) + s.cfg.P2P.Encoding().ProtocolSuffix()
			if err := s.cfg.P2P.PubSub().UnregisterTopicValidator(fullTopic); err != nil {
				log.WithError(err).Error("Could not unregister topic validator")
			}
			delete(subscriptions, k)
		}
	}
}

// find if we have peers who are subscribed to the same subnet
func (s *Service) validPeersExist(subnetTopic string) bool {
	numOfPeers := s.cfg.P2P.PubSub().ListPeers(subnetTopic + s.cfg.P2P.Encoding().ProtocolSuffix())
	return uint64(len(numOfPeers)) >= params.CarrierNetworkConfig().MinimumPeersInSubnet
}

// Add fork digest to topic.
func (s *Service) addDigestToTopic(topic string) string {
	if !strings.Contains(topic, "%x") {
		log.Fatal("Topic does not have appropriate formatter for digest")
	}
	digest, err := s.forkDigest()
	if err != nil {
		log.WithError(err).Fatal("Could not compute fork digest")
	}
	return fmt.Sprintf(topic, digest)
}

// forkDigest returns the current fork digest of the node.
func (s *Service) forkDigest() ([4]byte, error) {
	return [4]byte{0x1, 0x1, 0x1, 0x1,},nil
}
