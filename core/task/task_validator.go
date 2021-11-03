package task

import (
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

func (tv *TaskValidator) validateTaskMsg (msgs types.TaskMsgArr) (types.TaskMsgArr, types.TaskMsgArr, error) {

	identity, err := tv.resourceMng.GetDB().QueryIdentity()
	if nil != err {
		return nil, nil, err
	}

	badMsgs := make(types.TaskMsgArr, 0)
	goodMsgs := make(types.TaskMsgArr, 0)

	loop:
	for _, msg := range msgs {
		if msg.GetSenderIdentityId() != identity.GetIdentityId() {
			log.Errorf("Failed to verify sender's identity of task on TaskValidator.validateTaskMsg(), taskId: {%s}, sender's idntityId: {%s}, current identityId: {%s}",
				msg.GetTaskId(), msg.GetSenderIdentityId(), identity.GetIdentityId())
			badMsgs = append(badMsgs, msg)
			continue
		}
		for _, dataSupplier := range msg.GetTaskMetadataSupplierDatas() {
			if dataSupplier.GetOrganization().GetIdentityId() == identity.GetIdentityId() {
					if !tv.authMng.VerifyMetadataAuth(msg.GetUserType(), msg.GetUser(), dataSupplier.GetMetadataId()) {
						log.Errorf("Failed to verify metadataAuth of task on TaskValidator.validateTaskMsg(), taskId: {%s}, partyId: {%s}, userType: {%s}, user: {%s}, metadataId: {%s}",
							msg.GetTaskId(), dataSupplier.GetOrganization().GetPartyId(), msg.GetUserType(), msg.GetUser(), dataSupplier.GetMetadataId())
						badMsgs = append(badMsgs, msg)
						continue loop // goto continue next msg...
					}
			}
		}
		goodMsgs = append(goodMsgs, msg)
	}
	return badMsgs, goodMsgs, nil
}