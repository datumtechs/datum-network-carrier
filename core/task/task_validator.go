package task

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/auth"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	"github.com/RosettaFlow/Carrier-Go/types"
)

type  TaskValidator struct {
	resourceMng     *resource.Manager
	authMng         *auth.AuthorityManager
}

func newTaskValidator (resourceMng *resource.Manager, authMng *auth.AuthorityManager) *TaskValidator {
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

	loop:
	for _, msg := range msgs {
		if msg.GetSenderIdentityId() != identity.GetIdentityId() {
			log.Errorf("Failed to check sender's identity of task, is not current identity on TaskValidator.validateTaskMsg(), taskId: {%s}, sender's idntityId: {%s}, current identityId: {%s}",
				msg.GetTaskId(), msg.GetSenderIdentityId(), identity.GetIdentityId())
			badMsgs = append(badMsgs, types.NewBadTaskMsg(msg, fmt.Sprintf("The identity of task is not current identity")))
			continue
		}
		for _, dataSupplier := range msg.GetTaskMetadataSupplierDatas() {
			if dataSupplier.GetOrganization().GetIdentityId() == identity.GetIdentityId() {
					if err := tv.authMng.VerifyMetadataAuth(msg.GetUserType(), msg.GetUser(), dataSupplier.GetMetadataId()); nil != err {
						log.WithError(err).Errorf("Failed to verify metadataAuth of task on TaskValidator.validateTaskMsg(), taskId: {%s}, partyId: {%s}, userType: {%s}, user: {%s}, metadataId: {%s}",
							msg.GetTaskId(), dataSupplier.GetOrganization().GetPartyId(), msg.GetUserType(), msg.GetUser(), dataSupplier.GetMetadataId())
						badMsgs = append(badMsgs, types.NewBadTaskMsg(msg, fmt.Sprintf("verify metadataAuth failed, %s", err)))
						continue loop // goto continue next msg...
					}
			}
		}
		goodMsgs = append(goodMsgs, msg)
	}
	return badMsgs, goodMsgs
}