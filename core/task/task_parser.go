package task

import (
	"github.com/datumtechs/datum-network-carrier/core/resource"
	"github.com/datumtechs/datum-network-carrier/types"
)

type TaskParser struct {
	resourceMng     *resource.Manager
}

func newTaskParser (resourceMng *resource.Manager) *TaskParser {
	return &TaskParser{resourceMng: resourceMng}
}

func (tp *TaskParser) ParseTask(tasks types.TaskMsgArr) (types.BadTaskMsgArr, types.TaskMsgArr) {
	// nonParsedMsgArr, parsedMsgArr, error
	return nil, tasks
}