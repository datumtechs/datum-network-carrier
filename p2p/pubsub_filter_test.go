package p2p

import (
	"fmt"
	"github.com/datumtechs/datum-network-carrier/p2p/encoder"
	"testing"
	"time"
)

func TestService_CanSubscribe(t *testing.T) {
	currentFork := [4]byte{0x01, 0x01, 0x01, 0x01}
	validProtocolSuffix := "/" + encoder.ProtocolSuffixSSZSnappy
	type test struct {
		name  string
		topic string
		want  bool
	}
	tests := []test{
		{
			name:  "topic on current fork",
			topic: fmt.Sprintf(GossipTestDataTopicFormat, currentFork) + validProtocolSuffix,
			want:  true,
		},
		{
			name:  "topic on unknown fork",
			topic: fmt.Sprintf(GossipTestDataTopicFormat, [4]byte{0xFF, 0xEE, 0x56, 0x21}) + validProtocolSuffix,
			want:  false,
		},
		{
			name:  "topic missing protocol suffix",
			topic: fmt.Sprintf(GossipTestDataTopicFormat, currentFork),
			want:  false,
		},
		{
			name:  "topic wrong protocol suffix",
			topic: fmt.Sprintf(GossipTestDataTopicFormat, currentFork) + "/foobar",
			want:  false,
		},
		{
			name:  "erroneous topic",
			topic: "hey, want to foobar?",
			want:  false,
		},
		{
			name:  "erroneous topic that has the correct amount of slashes",
			topic: "hey, want to foobar?////",
			want:  false,
		},
		{
			name:  "bad prefix",
			topic: fmt.Sprintf("/carrier/%x/foobar", currentFork) + validProtocolSuffix,
			want:  false,
		},
		{
			name:  "topic not in gossip mapping",
			topic: fmt.Sprintf("/carrier/%x/foobar", currentFork) + validProtocolSuffix,
			want:  false,
		},
	}

	// Ensure all gossip topic mappings pass validation.
	for topic := range GossipTopicMappings {
		formatting := []interface{}{currentFork}

		tt := test{
			name:  topic,
			topic: fmt.Sprintf(topic, formatting...) + validProtocolSuffix,
			want:  true,
		}
		tests = append(tests, tt)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				currentForkDigest:     currentFork,
				genesisValidatorsRoot: make([]byte, 32),
				genesisTime:           time.Now(),
			}
			if got := s.CanSubscribe(tt.topic); got != tt.want {
				t.Errorf("CanSubscribe(%s) = %v, want %v", tt.topic, got, tt.want)
			}
		})
	}
}