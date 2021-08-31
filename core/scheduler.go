package core

import "github.com/RosettaFlow/Carrier-Go/types"

type Scheduler interface {
	Start() error
	Stop() error
	Error () error
	Name() string
	AddTask(task *types.Task) error
	FetchTask() error

}