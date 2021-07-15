package chaincons

import (
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/libp2p/go-libp2p-core/peer"
)

type Chaincons struct {

}

func New() *Chaincons {return &Chaincons{}}
func (c *Chaincons)Start() error {return nil}
func (c *Chaincons) Close() error {return nil}
func (c *Chaincons)OnPrepare(task *types.Task) error {return nil}
func (c *Chaincons)OnHandle(task *types.Task,selfPeerResource *types.PrepareVoteResource, result chan<- *types.ConsensuResult) error  {return nil}
func (c *Chaincons) ValidateConsensusMsg(pid peer.ID, msg types.ConsensusMsg) error {return nil}
func (c *Chaincons) OnConsensusMsg(pid peer.ID, msg types.ConsensusMsg) error {return nil}
func (c *Chaincons)OnError() error  {return nil}