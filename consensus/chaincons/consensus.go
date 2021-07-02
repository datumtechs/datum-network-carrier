package chaincons

import "github.com/RosettaFlow/Carrier-Go/types"

type Chaincons struct {

}

func New() *Chaincons {return &Chaincons{}}

func (c *Chaincons)OnPrepare(task *types.ScheduleTask) error {return nil}
func (c *Chaincons)OnStart(task *types.ScheduleTask, result chan<- *types.TaskConsResult) error  {return nil}
func (c *Chaincons) ValidateConsensusMsg(msg types.ConsensusMsg) error {return nil}
func (c *Chaincons) OnConsensusMsg(msg types.ConsensusMsg) error {return nil}
func (c *Chaincons)OnError() error  {return nil}