package task

import "github.com/RosettaFlow/Carrier-Go/types"

type TaskParser struct {

}

func newTaskParser () *TaskParser {
	return &TaskParser{}
}

func (tp *TaskParser) ParseTask(tasks types.TaskMsgArr) (types.TaskMsgArr, types.TaskMsgArr, error) {
	return nil, nil, nil
}