package task

import "github.com/RosettaFlow/Carrier-Go/types"

type  TaskValidator struct {

}

func newTaskValidator () *TaskValidator {
	return &TaskValidator{}
}

func (tv *TaskValidator) validateTaskMsg (tasks types.TaskMsgs) (types.TaskMsgs, error) {
	return nil, nil
}