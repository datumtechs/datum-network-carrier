package messagehandler

import (
	"context"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/gogo/protobuf/proto"
	logTest "github.com/sirupsen/logrus/hooks/test"
	"testing"
)

func TestSafelyHandleMessage(t *testing.T) {
	SafelyHandleMessage(context.Background(), func(_ context.Context, _ proto.Message) error {
		panic("bad!")
		return nil
	}, &libtypes.BlockData{})
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
