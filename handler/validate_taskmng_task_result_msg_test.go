package handler

import (
	"bytes"
	"context"
	"github.com/datumtechs/datum-network-carrier/common/timeutils"
	"github.com/datumtechs/datum-network-carrier/common/traceutil"
	"github.com/datumtechs/datum-network-carrier/lib/netmsg/common"
	"github.com/datumtechs/datum-network-carrier/lib/netmsg/taskmng"
	"github.com/datumtechs/datum-network-carrier/p2p"
	p2ptest "github.com/datumtechs/datum-network-carrier/p2p/testing"
	lru "github.com/hashicorp/golang-lru"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	pubsubpb "github.com/libp2p/go-libp2p-pubsub/pb"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	"reflect"
	"testing"
)

func TestValidateTwopc_ValidTaskResult(t *testing.T) {
	p := p2ptest.NewTestP2P(t)
	ctx := context.Background()
	c, err := lru.New(10)
	r := &Service{
		cfg: &Config{
			P2P:         p,
			InitialSync: &p2ptest.Sync{IsSyncing: false},
			TaskManager: &p2ptest.MockTaskManager{},
		},
		seenTaskResultMsgCache: c,
	}

	buf := new(bytes.Buffer)
	pbmsg := &taskmng.TaskResultMsg{
		MsgOption: &common.MsgOption{
			ProposalId:      []byte("proposalId"),
			SenderRole:      0,
			SenderPartyId:   []byte("SenderPartyId"),
			ReceiverRole:    0,
			ReceiverPartyId: []byte("ReceiverPartyId"),
			MsgOwner:        nil,
		},
		CreateAt: timeutils.UnixMsecUint64(),
		Sign:     []byte("sign"),
	}
	_, err = p.Encoding().EncodeGossip(buf, pbmsg)
	require.NoError(t, err)

	topic := p2p.GossipTypeMapping[reflect.TypeOf(&taskmng.TaskResultMsg{})]
	msg := &pubsub.Message{
		Message: &pubsubpb.Message{
			Data:  buf.Bytes(),
			Topic: &topic,
		},
	}
	t.Log(traceutil.GenerateTraceID(pbmsg))
	t.Log(traceutil.GenerateTraceIDForPub(msg))
	valid := r.validateTaskResultMessagePubSub(ctx, "foobar", msg) == pubsub.ValidationIgnore
	//todo: Need to add validation against consensus modules at a later stage.
	assert.Equal(t, true, valid, "Failed Validation")
	//require.NotNil(t, msg.ValidatorData, "Decoded message was not set on the message validator data")
}