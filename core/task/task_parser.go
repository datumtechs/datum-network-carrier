package task

import (
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	"github.com/RosettaFlow/Carrier-Go/types"
)

type TaskParser struct {
	resourceMng     *resource.Manager
}

func newTaskParser (resourceMng *resource.Manager) *TaskParser {
	return &TaskParser{resourceMng: resourceMng}
}

func (tp *TaskParser) ParseTask(tasks types.TaskMsgArr) (types.TaskMsgArr, types.TaskMsgArr, error) {
	return nil, nil, nil
}