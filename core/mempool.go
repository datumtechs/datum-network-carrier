package core

import (
	"github.com/RosettaFlow/Carrier-Go/core/message"
	"github.com/RosettaFlow/Carrier-Go/types"
)

type MempoolConfig struct {
}

type Mempool struct {
	cfg MempoolConfig

	mateDataMsgQueue *message.MetaDataMsgList
	powerMsgQueue 	 *message.PowerMsgList
	taskMsgQueue     *message.TaskMsgList

}

type Msg interface {
	Marshal() ([]byte, error)
	Unmarshal(b []byte) error
	String() string
	MsgType () string
}

func (pool *Mempool) Add(msg Msg) error {


	switch msg.(type) {
	case *types.PowerMsg:
	case *types.MetaDataMsg:
	case *types.TaskMsg:

	default:
		log.Fatalf("Failed to add msg, can not match the msg type")
	}

	return nil
}