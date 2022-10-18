package rawdb

var (
	workflowMsgKeyPrefix             = []byte("workflowMsgKeyPrefix:")
	sendToTaskManagerCacheKeyPrefix  = []byte("sendToTaskManagerCacheKeyPrefix:")
	workflowsCacheKeyPrefix          = []byte("workflowsCacheKeyPrefix:")
	workflowStatusCacheKeyPrefix     = []byte("workflowStatusCacheKeyPrefix:")
	workflowTaskStatusCacheKeyPrefix = []byte("workflowTaskStatusCacheKeyPrefix:")
)

func GetWorkflowMsgKey(workflowId string) []byte {
	return append(workflowMsgKeyPrefix, []byte(workflowId)...)
}

func GetSendToTaskManagerCacheKeyPrefix(taskId string) []byte {
	return append(sendToTaskManagerCacheKeyPrefix, []byte(taskId)...)
}

func GetWorkflowsCacheKeyPrefix(workflowId string) []byte {
	return append(workflowsCacheKeyPrefix, []byte(workflowId)...)
}

func GetWorkflowStatusCacheKeyPrefix(workflowId string) []byte {
	return append(workflowStatusCacheKeyPrefix, []byte(workflowId)...)
}

func GetWorkflowTaskStatusCacheKeyPrefix(workflowId, taskName string) []byte {
	return append(append(workflowTaskStatusCacheKeyPrefix, []byte(workflowId)...), taskName...)
}

func QueryWorkflowMsgKey() []byte {
	return workflowMsgKeyPrefix
}
