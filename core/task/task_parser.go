package task

import "github.com/RosettaFlow/Carrier-Go/types"

type TaskParser struct {

}

func newTaskParser () *TaskParser {
	return &TaskParser{}
}

func (tp *TaskParser) ParseTask(tasks types.TaskMsgs) (types.TaskMsgs, error) {
	return nil, nil
}