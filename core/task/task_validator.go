package task

import "github.com/RosettaFlow/Carrier-Go/types"

type  TaskValidator struct {

}

func newTaskValidator () *TaskValidator {
	return &TaskValidator{}
}

func (tv *TaskValidator) validateTaskMsg (tasks types.TaskMsgArr) (types.TaskMsgArr, types.TaskMsgArr, error) {
	return nil, nil, nil
}