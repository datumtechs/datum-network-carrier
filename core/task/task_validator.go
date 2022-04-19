package task

import (
	"fmt"
	auth2 "github.com/RosettaFlow/Carrier-Go/ach/auth"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	"github.com/RosettaFlow/Carrier-Go/types"
)

type  TaskValidator struct {
	resourceMng     *resource.Manager
	authMng         *auth2.AuthorityManager
}

func newTaskValidator (resourceMng *resource.Manager, authMng *auth2.AuthorityManager) *TaskValidator {
	return &TaskValidator{
		resourceMng: resourceMng,
		authMng: authMng,
	}
}

func (tv *TaskValidator) validateTaskMsg (msgs types.TaskMsgArr) (types.BadTaskMsgArr, types.TaskMsgArr) {

	badMsgs := make(types.BadTaskMsgArr, 0)
	goodMsgs := make(types.TaskMsgArr, 0)

	identity, err := tv.resourceMng.GetDB().QueryIdentity()
	if nil != err {
		for _, msg := range msgs {
			badMsgs = append(badMsgs, types.NewBadTaskMsg(msg, fmt.Sprintf("query indeityt failed, %s", err)))
		}
		return badMsgs, goodMsgs
	}

	for _, msg := range msgs {
		if msg.GetSenderIdentityId() != identity.GetIdentityId() {
			log.Errorf("Failed to check sender's identity of task, is not current identity on TaskValidator.validateTaskMsg(), taskId: {%s}, sender's idntityId: {%s}, current identityId: {%s}",
				msg.GetTaskId(), msg.GetSenderIdentityId(), identity.GetIdentityId())
			badMsgs = append(badMsgs, types.NewBadTaskMsg(msg, fmt.Sprintf("the sender identity of task is not current identity")))
			continue
		}
		goodMsgs = append(goodMsgs, msg)
	}
	return badMsgs, goodMsgs
}