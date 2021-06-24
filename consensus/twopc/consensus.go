package twopc

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/RosettaFlow/Carrier-Go/types"
	"strings"
)


type twoPC struct {
	config           *Config
	Errs 			 []error
	p2p 			 p2p.P2P
	currTask 		 *types.ScheduleTask
}

func New(conf *Config) *twoPC {

	t := &twoPC{
		config: conf,
		Errs: make([]error, 0),
	}

	return t
}

func (t *twoPC) OnPrepare(task *types.ScheduleTask) error {
	return nil
}
func (t *twoPC) OnStart(task *types.ScheduleTask, result chan<- *types.ScheduleResult) error {
	return nil
}
func (t *twoPC) OnError() error {
	if len(t.Errs) == 0 {
		return nil
	}
	errStrs := make([]string, len(t.Errs))
	for _, err := range t.Errs {
		errStrs = append(errStrs, err.Error())
	}
	// reset Errs
	t.Errs = make([]error, 0)
	return fmt.Errorf("%s", strings.Join(errStrs, "\n"))
}

func (t *twoPC) OnPrepareMsg (proposal *PrepareMsg) error {

	return nil
}