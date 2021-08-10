package p2p

import (
	"bytes"
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/traceutil"
	"github.com/RosettaFlow/Carrier-Go/params"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"reflect"
	"time"
)

// ErrMessageNotMapped occurs on a Broadcast attempt when a message has not been defined in the GossipTypeMapping.
var ErrMessageNotMapped = errors.New("message type is not mapped to a PubSub topic")

// Broadcast a message to the p2p network.
func (s *Service) Broadcast(ctx context.Context, msg proto.Message) error {
	ctx, span := trace.StartSpan(ctx, "p2p.Broadcast")
	defer span.End()

	twoSlots := time.Duration(2*params.CarrierChainConfig().SecondsPerSlot) * time.Second
	ctx, cancel := context.WithTimeout(ctx, twoSlots)
	defer cancel()

	forkDigest, err := s.forkDigest()
	if err != nil {
		err := errors.Wrap(err, "could not retrieve fork digest")
		traceutil.AnnotateError(span, err)
		return err
	}

	topic, ok := GossipTypeMapping[reflect.TypeOf(msg)]
	if !ok {
		traceutil.AnnotateError(span, ErrMessageNotMapped)
		return ErrMessageNotMapped
	}
	return s.broadcastObject(ctx, msg, fmt.Sprintf(topic, forkDigest))
}

// To broadcast message to other peers in our gossip mesh.
func (s *Service) broadcastObject(ctx context.Context, obj interface{}, topic string) error {
	buf := new(bytes.Buffer)
	if _, err := s.Encoding().EncodeGossip(buf, obj); err != nil {
		err := errors.Wrap(err, "could not encode message")
		return err
	}
	if err := s.PublishToTopic(ctx, topic+s.Encoding().ProtocolSuffix(), buf.Bytes()); err != nil {
		err := errors.Wrap(err, "could not publish message")
		return err
	}
	return nil
}

func attestationToTopic(subnet uint64, forkDigest [4]byte) string {
	return fmt.Sprintf(AttestationSubnetTopicFormat, forkDigest, subnet)
}

