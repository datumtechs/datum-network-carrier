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

	identity, err := tv.resourceMng.GetDB().GetIdentity()
	if nil != err {
		return nil, nil, err
	}

	badMsgs := make(types.TaskMsgArr, 0)
	goodMsgs := make(types.TaskMsgArr, 0)

	for _, msg := range msgs {
		if msg.GetSenderIdentityId() != identity.GetIdentityId() {
			badMsgs = append(badMsgs, msg)
			continue
		}
		for _, dataSupplier := range msg.GetTaskMetadataSupplierDatas() {
			if dataSupplier.GetOrganization().GetIdentityId() == identity.GetIdentityId() {
					if !tv.authMng.VerifyMetadataAuth(msg.GetUserType(), msg.GetUser(), dataSupplier.GetMetadataId()) {
						badMsgs = append(badMsgs, msg)
						continue
					}
			}
		}
		goodMsgs = append(goodMsgs, msg)
	}

	return badMsgs, goodMsgs, nil
}