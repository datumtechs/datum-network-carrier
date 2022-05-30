package handler

import (
	"bytes"
	carrierrpcdebugpbv1 "github.com/datumtechs/datum-network-carrier/pb/carrier/rpc/debug/v1"
	"github.com/datumtechs/datum-network-carrier/p2p"
	p2ptesting "github.com/datumtechs/datum-network-carrier/p2p/testing"
	"github.com/d4l3k/messagediff"
	"github.com/gogo/protobuf/proto"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	pb "github.com/libp2p/go-libp2p-pubsub/pb"
	"reflect"
	"testing"
)

func TestService_decodePubsubMessage(t *testing.T) {
	tests := []struct {
		name    string
		topic   string
		input   *pubsub.Message
		want    proto.Message
		wantErr error
	}{
		{
			name:    "Nil message",
			input:   nil,
			wantErr: errNilPubsubMessage,
		},
		{
			name: "nil topic",
			input: &pubsub.Message{
				Message: &pb.Message{
					Topic: nil,
				},
			},
			wantErr: errNilPubsubMessage,
		},
		{
			name:    "invalid topic format",
			topic:   "foo",
			wantErr: errInvalidTopic,
		},
		{
			name:    "topic not mapped to any message type",
			topic:   "/carrier/abcdef/foo",
			wantErr: p2p.ErrMessageNotMapped,
		},
		{
			name:  "valid message -- beacon block",
			topic: p2p.GossipTypeMapping[reflect.TypeOf(&carrierrpcdebugpbv1.GossipTestData{})],
			input: &pubsub.Message{
				Message: &pb.Message{
					Data: func() []byte {
						buf := new(bytes.Buffer)
						if _, err := p2ptesting.NewTestP2P(t).Encoding().EncodeGossip(buf, NewGossipTestData()); err != nil {
							t.Fatal(err)
						}
						return buf.Bytes()
					}(),
				},
			},
			wantErr: nil,
			want:    NewGossipTestData(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				cfg: &Config{P2P: p2ptesting.NewTestP2P(t)},
			}
			if tt.topic != "" {
				if tt.input == nil {
					tt.input = &pubsub.Message{Message: &pb.Message{}}
				} else if tt.input.Message == nil {
					tt.input.Message = &pb.Message{}
				}
				tt.input.Message.Topic = &tt.topic
			}
			got, err := s.decodePubsubMessage(tt.input)
			if err != tt.wantErr {
				t.Errorf("decodePubsubMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				diff, _ := messagediff.PrettyDiff(got, tt.want)
				t.Log(diff)
				t.Errorf("decodePubsubMessage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func NewGossipTestData() *carrierrpcdebugpbv1.GossipTestData {
	return  &carrierrpcdebugpbv1.GossipTestData{
		Data:                 []byte("data"),
		Count:                11,
		Step:                 23,
	}
}

