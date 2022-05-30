package messagehandler

import (
	"context"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"github.com/gogo/protobuf/proto"
	logTest "github.com/sirupsen/logrus/hooks/test"
	"testing"
)

func TestSafelyHandleMessage(t *testing.T) {
	SafelyHandleMessage(context.Background(), func(_ context.Context, _ proto.Message) error {
		panic("bad!")
		return nil
	}, &carriertypespb.BlockData{})
}

func TestSafelyHandleMessage_NoData(t *testing.T) {
	hook := logTest.NewGlobal()

	SafelyHandleMessage(context.Background(), func(_ context.Context, _ proto.Message) error {
		panic("bad!")
		return nil
	}, nil)

	entry := hook.LastEntry()
	if entry.Data["msg"] != "message contains no data" {
		t.Errorf("Message logged was not what was expected: %s", entry.Data["msg"])
	}
}
