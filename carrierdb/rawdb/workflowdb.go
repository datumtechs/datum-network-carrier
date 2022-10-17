package rawdb

import (
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/datumtechs/datum-network-carrier/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gogo/protobuf/proto"
)

func RemoveWorkflowMsg(db KeyValueStore, workflowId string) error {
	key := GetWorkflowMsgKey(workflowId)
	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

func SaveSendToTaskManager(db KeyValueStore, taskId, workflowId string) error {
	key := GetSendToTaskManagerCacheKeyPrefix(taskId)
	val, err := rlp.EncodeToBytes(workflowId)
	if nil != err {
		return err
	}
	return db.Put(key, val)
}

func RemoveSendToTaskManager(db KeyValueStore, taskId string) error {
	key := GetSendToTaskManagerCacheKeyPrefix(taskId)
	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

func SaveWorkflowCache(db KeyValueStore, workflow *carriertypespb.Workflow) error {
	//carriertypespb.Workflow
	key := GetWorkflowsCacheKeyPrefix(workflow.GetWorkflowId())
	val, err := proto.Marshal(workflow)
	if nil != err {
		return err
	}
	return db.Put(key, val)
}

func RemoveWorkflowCache(db KeyValueStore, workflowId string) error {
	key := GetWorkflowsCacheKeyPrefix(workflowId)
	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

func SaveWorkflowStatusCache(db KeyValueStore, workflowId string, status commonconstantpb.WorkFlowState) error {
	key := GetWorkflowStatusCacheKeyPrefix(workflowId)
	val, err := rlp.EncodeToBytes(uint32(status))
	if nil != err {
		return err
	}
	return db.Put(key, val)
}

func RemoveWorkflowStatusCache(db KeyValueStore, workflowId string) error {
	key := GetWorkflowStatusCacheKeyPrefix(workflowId)
	has, err := db.Has(key)
	switch {
	case IsNoDBNotFoundErr(err):
		return err
	case IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return db.Delete(key)
}

func SaveWorkflowTaskStatusCache(db KeyValueStore, workflowId string, taskState *carrierapipb.WorkFlowTaskStatus) error {
	key := GetWorkflowTaskStatusCacheKeyPrefix(workflowId, taskState.GetTaskName())
	val, err := proto.Marshal(taskState)
	if nil != err {
		return err
	}
	return db.Put(key, val)
}

func ForEachKVWithPrefix(db KeyValueStore, prefix []byte, f func(key, value []byte) error) error {
	it := db.NewIteratorWithPrefixAndStart(prefix, nil)
	defer it.Release()
	for it.Next() {
		if err := f(it.Key(), it.Value()); nil != err {
			return err
		}
	}
	return nil
}

func QueryWorkflowMsgArr(db KeyValueStore) (types.WorkflowMsgArr, error) {
	it := db.NewIteratorWithPrefixAndStart(QueryWorkflowMsgKey(), nil)
	defer it.Release()

	arr := make(types.WorkflowMsgArr, 0)

	for it.Next() {
		if val := it.Value(); len(val) != 0 {
			var res carriertypespb.WorkflowMsg
			if err := proto.Unmarshal(val, &res); nil != err {
				continue
			}
			taskList := make([]*types.TaskMsg, 0)
			for _, v := range res.Data.Tasks {
				taskList = append(taskList, &types.TaskMsg{
					Data: types.NewTask(v.GetData()),
				})
			}
			arr = append(arr, &types.WorkflowMsg{
				Data: &types.Workflow{
					WorkflowId:   res.Data.GetWorkflowId(),
					Desc:         res.Data.Desc,
					WorkflowName: res.Data.GetWorkflowName(),
					PolicyType:   res.Data.PolicyType,
					Policy:       res.Data.Policy,
					User:         res.Data.User,
					UserType:     res.Data.UserType,
					Sign:         res.Data.Sign,
					Tasks:        taskList,
					CreateAt:     res.Data.CreateAt,
				},
			})
		}
	}
	if len(arr) == 0 {
		return nil, ErrNotFound
	}
	return arr, nil
}
