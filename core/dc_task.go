package core

import (
	"errors"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	"github.com/RosettaFlow/Carrier-Go/lib/center/api"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
)

// about task on local
// local task
func (dc *DataCenter) StoreLocalTask(task *types.Task) error {
	if task == nil {
		return errors.New("invalid params for task")
	}
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreLocalTask(dc.db, task)
}

func (dc *DataCenter) RemoveLocalTask(taskId string) error {
	if taskId == "" {
		return errors.New("invalid params for taskId to DelLocalTask")
	}
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveLocalTask(dc.db, taskId)
}

func (dc *DataCenter) QueryLocalTask(taskId string) (*types.Task, error) {
	if taskId == "" {
		return nil, errors.New("invalid params taskId for QueryLocalTask")
	}
	dc.mu.Lock()
	defer dc.mu.Unlock()
	//log.Debugf("QueryLocalTask, taskId: {%s}", taskId)
	return rawdb.QueryLocalTask(dc.db, taskId)
}

func (dc *DataCenter) QueryLocalTaskListByIds(taskIds []string) (types.TaskDataArray, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryLocalTaskByIds(dc.db, taskIds)
}

func (dc *DataCenter) QueryLocalTaskList() (types.TaskDataArray, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryAllLocalTasks(dc.db)
}

func (dc *DataCenter) QueryLocalTaskAndEvents(taskId string) (*types.Task, error) {
	if taskId == "" {
		return nil, errors.New("invalid params taskId for QueryLocalTask")
	}
	dc.mu.Lock()
	defer dc.mu.Unlock()
	task, err := rawdb.QueryLocalTask(dc.db, taskId)
	if nil != err {
		return nil, err
	}
	list, err := rawdb.QueryTaskEvent(dc.db, taskId)
	if nil != err {
		return nil, err
	}
	task.SetEventList(list)
	return task, nil
}

func (dc *DataCenter) QueryLocalTaskAndEventsListByIds(taskIds []string) (types.TaskDataArray, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	tasks, err := rawdb.QueryLocalTaskByIds(dc.db, taskIds)
	if nil != err {
		return nil, err
	}
	for i, task := range tasks {
		list, err := rawdb.QueryTaskEvent(dc.db, task.GetTaskId())
		if nil != err {
			return nil, err
		}
		task.SetEventList(list)
		tasks[i] = task
	}
	return tasks, nil
}

func (dc *DataCenter) QueryLocalTaskAndEventsList() (types.TaskDataArray, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()

	tasks, err := rawdb.QueryAllLocalTasks(dc.db)
	if nil != err {
		return nil, err
	}
	for i, task := range tasks {
		list, err := rawdb.QueryTaskEvent(dc.db, task.GetTaskId())
		if nil != err {
			return nil, err
		}
		task.SetEventList(list)
		tasks[i] = task
	}
	return tasks, nil
}

func (dc *DataCenter) StoreJobNodeRunningTaskId(jobNodeId, taskId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreResourceTaskId(dc.db, jobNodeId, taskId)
}

func (dc *DataCenter) RemoveJobNodeRunningTaskId(jobNodeId, taskId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveResourceTaskId(dc.db, jobNodeId, taskId)
}

func (dc *DataCenter) QueryRunningTaskCountOnJobNode(jobNodeId string) (uint32, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	taskIds, err := rawdb.QueryResourceTaskIds(dc.db, jobNodeId)
	if rawdb.IsNoDBNotFoundErr(err) {
		return 0, err
	} else if rawdb.IsDBNotFoundErr(err) {
		return 0, nil
	}
	return uint32(len(taskIds)), nil
}

func (dc *DataCenter) QueryJobNodeRunningTaskIdList(jobNodeId string) ([]string, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	ids, err := rawdb.QueryResourceTaskIds(dc.db, jobNodeId)
	if rawdb.IsNoDBNotFoundErr(err) {
		return nil, err
	} else if rawdb.IsDBNotFoundErr(err) {
		return nil, nil
	}
	return ids, nil
}

func (dc *DataCenter) IncreaseResourceTaskPartyIdCount(jobNodeId, taskId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.IncreaseResourceTaskPartyIdCount(dc.db, jobNodeId, taskId)
}

func (dc *DataCenter) DecreaseResourceTaskPartyIdCount(jobNodeId, taskId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.DecreaseResourceTaskPartyIdCount(dc.db, jobNodeId, taskId)
}

func (dc *DataCenter) QueryResourceTaskPartyIdCount(jobNodeId, taskId string) (uint32, error) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.QueryResourceTaskPartyIdCount(dc.db, jobNodeId, taskId)
}

func (dc *DataCenter) IncreaseResourceTaskTotalCount(jobNodeId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.IncreaseResourceTaskTotalCount(dc.db, jobNodeId)
}

func (dc *DataCenter) DecreaseResourceTaskTotalCount(jobNodeId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.DecreaseResourceTaskTotalCount(dc.db, jobNodeId)
}

func (dc *DataCenter) RemoveResourceTaskTotalCount(jobNodeId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveResourceTaskTotalCount(dc.db, jobNodeId)
}

func (dc *DataCenter) QueryResourceTaskTotalCount(jobNodeId string) (uint32, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryResourceTaskTotalCount(dc.db, jobNodeId)
}

func (dc *DataCenter) StoreTaskEvent(event *libtypes.TaskEvent) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreTaskEvent(dc.db, event)
}

func (dc *DataCenter) QueryTaskEventList(taskId string) ([]*libtypes.TaskEvent, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryTaskEvent(dc.db, taskId)
}

func (dc *DataCenter) QueryTaskEventListByPartyId (taskId, partyId string) ([]*libtypes.TaskEvent, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryTaskEventByPartyId (dc.db, taskId, partyId)
}

func (dc *DataCenter) RemoveTaskEventList(taskId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveTaskEvent(dc.db, taskId)
}

func (dc *DataCenter) RemoveTaskEventListByPartyId(taskId, partyId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveTaskEventByPartyId(dc.db, taskId, partyId)
}

// about task on datacenter
func (dc *DataCenter) InsertTask(task *types.Task) error {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	response, err := dc.client.SaveTask(dc.ctx, types.NewSaveTaskRequest(task))
	if err != nil {
		log.WithError(err).WithField("taskId", task.GetTaskId()).Errorf("InsertTask failed")
		return err
	}
	if response.Status != 0 {
		return fmt.Errorf("insert task, taskId: {%s},  error: %s", task.GetTaskId(), response.Msg)
	}
	return nil
}

func (dc *DataCenter) QueryTaskListByIdentityId(identityId string) (types.TaskDataArray, error) {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	taskListResponse, err := dc.client.ListTaskByIdentity(dc.ctx, &api.ListTaskByIdentityRequest{
		LastUpdated: timeutils.UnixMsecUint64(),
		IdentityId:  identityId,
	})
	return types.NewTaskArrayFromResponse(taskListResponse), err
}

func (dc *DataCenter) QueryRunningTaskCountOnOrg() uint32 {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	taskList, err := rawdb.QueryAllLocalTasks(dc.db)
	if nil != err {
		return 0
	}
	if taskList != nil {
		return uint32(taskList.Len())
	}
	return 0
}

func (dc *DataCenter) QueryTaskEventListByTaskId(taskId string) ([]*libtypes.TaskEvent, error) {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	taskEventResponse, err := dc.client.ListTaskEvent(dc.ctx, &api.ListTaskEventRequest{
		TaskId: taskId,
	})
	if nil != err {
		return nil, err
	}
	if nil == taskEventResponse {
		return nil, rawdb.ErrNotFound
	}
	return taskEventResponse.GetTaskEvents(), nil
}

func (dc *DataCenter) QueryTaskEventListByTaskIds(taskIds []string) ([]*libtypes.TaskEvent, error) {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()

	eventList := make([]*libtypes.TaskEvent, 0)
	for _, taskId := range taskIds {
		taskEventResponse, err := dc.client.ListTaskEvent(dc.ctx, &api.ListTaskEventRequest{
			TaskId: taskId,
		})
		if nil != err {
			//return nil, err
			continue
		}
		if nil == taskEventResponse {
			//return nil, rawdb.ErrNotFound
			continue
		}
		eventList = append(eventList, taskEventResponse.GetTaskEvents()...)
	}
	return eventList, nil
}


// about TaskPowerUsed
func (dc *DataCenter) StoreLocalTaskPowerUsed(taskPowerUsed *types.LocalTaskPowerUsed) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreLocalTaskPowerUsed(dc.db, taskPowerUsed)
}

func (dc *DataCenter) StoreLocalTaskPowerUseds(taskPowerUseds []*types.LocalTaskPowerUsed) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreLocalTaskPowerUseds(dc.db, taskPowerUseds)
}

func (dc *DataCenter) HasLocalTaskPowerUsed(taskId, partyId string) (bool, error) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.HasLocalTaskPowerUsed(dc.db, taskId, partyId)
}

func (dc *DataCenter) RemoveLocalTaskPowerUsed(taskId, partyId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveLocalTaskPowerUsed(dc.db, taskId, partyId)
}

func (dc *DataCenter) RemoveLocalTaskPowerUsedByTaskId(taskId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveLocalTaskPowerUsedByTaskId(dc.db, taskId)
}

func (dc *DataCenter) QueryLocalTaskPowerUsed(taskId, partyId string) (*types.LocalTaskPowerUsed, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryLocalTaskPowerUsed(dc.db, taskId, partyId)
}

func (dc *DataCenter) QueryLocalTaskPowerUsedsByTaskId(taskId string) ([]*types.LocalTaskPowerUsed, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	res, _ := rawdb.QueryLocalTaskPowerUsedsByTaskId(dc.db, taskId)
	return res, nil
}

func (dc *DataCenter) QueryLocalTaskPowerUseds() ([]*types.LocalTaskPowerUsed, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	res, _ := rawdb.QueryLocalTaskPowerUseds(dc.db)
	return res, nil
}

// about TaskResultFileMetadataId
func (dc *DataCenter) StoreTaskUpResultFile(turf *types.TaskUpResultFile)  error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreTaskUpResultFile(dc.db, turf)
}

func (dc *DataCenter) QueryTaskUpResultFile(taskId string)  (*types.TaskUpResultFile, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryTaskUpResultFile(dc.db, taskId)
}

func (dc *DataCenter) QueryTaskUpResultFileList () ([]*types.TaskUpResultFile, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryTaskUpResultFileList(dc.db)
}

func (dc *DataCenter) RemoveTaskUpResultFile(taskId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveTaskUpResultFile(dc.db, taskId)
}

func (dc *DataCenter) StoreTaskResuorceUsage(usage *types.TaskResuorceUsage) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreTaskResuorceUsage(dc.db, usage)
}

func (dc *DataCenter) QueryTaskResuorceUsage(taskId, partyId string) (*types.TaskResuorceUsage, error)  {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryTaskResuorceUsage(dc.db, taskId, partyId)
}

func (dc *DataCenter) RemoveTaskResuorceUsage(taskId, partyId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveTaskResuorceUsage(dc.db, taskId, partyId)
}

func (dc *DataCenter) RemoveTaskResuorceUsageByTaskId (taskId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveTaskResuorceUsageByTaskId(dc.db, taskId)
}

func (dc *DataCenter) StoreTaskPowerPartyIds(taskId string, powerPartyIds []string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreTaskPowerPartyIds(dc.db, taskId, powerPartyIds)
}

func (dc *DataCenter) QueryTaskPowerPartyIds(taskId string) ([]string, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryTaskPowerPartyIds(dc.db, taskId)
}

func (dc *DataCenter) RemoveTaskPowerPartyIds (taskId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveTaskPowerPartyIds(dc.db, taskId)
}

func (dc *DataCenter) StoreTaskPartnerPartyIds(taskId string, partyIds []string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreTaskPartnerPartyIds(dc.db, taskId, partyIds)
}

func (dc *DataCenter) HasTaskPartnerPartyIds(taskId string) (bool, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.HasTaskPartnerPartyIds(dc.db, taskId)
}

func (dc *DataCenter) QueryTaskPartnerPartyIds(taskId string) ([]string, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryTaskPartnerPartyIds(dc.db, taskId)
}

func (dc *DataCenter) RemoveTaskPartnerPartyId (taskId, partyId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	log.Debugf("Start remove partyId of local task's partner arr, taskId: {%s}, partyId: {%s}", taskId, partyId)
	return rawdb.RemoveTaskPartnerPartyId(dc.db, taskId, partyId)
}

func (dc *DataCenter) RemoveTaskPartnerPartyIds (taskId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveTaskPartnerPartyIds(dc.db, taskId)
}

// v 1.0 -> v 2.0 about task exec status (prefix + taskId + partyId -> "cons"|"exec")
func (dc *DataCenter) StoreLocalTaskExecuteStatusValConsByPartyId(taskId, partyId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreLocalTaskExecuteStatusValConsByPartyId(dc.db, taskId, partyId)
}

func (dc *DataCenter) StoreLocalTaskExecuteStatusValExecByPartyId(taskId, partyId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreLocalTaskExecuteStatusValExecByPartyId(dc.db, taskId, partyId)
}

func (dc *DataCenter) RemoveLocalTaskExecuteStatusByPartyId(taskId, partyId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveLocalTaskExecuteStatusByPartyId(dc.db, taskId, partyId)
}

func (dc *DataCenter) HasLocalTaskExecuteStatusParty(taskId string) (bool, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.HasLocalTaskExecuteStatusParty(dc.db, taskId)
}

func (dc *DataCenter) HasLocalTaskExecuteStatusByPartyId(taskId, partyId string) (bool, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.HasLocalTaskExecuteStatusByPartyId(dc.db, taskId, partyId)
}

func (dc *DataCenter) HasLocalTaskExecuteStatusValConsByPartyId(taskId, partyId string) (bool, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.HasLocalTaskExecuteStatusValConsByPartyId(dc.db, taskId, partyId)
}

func (dc *DataCenter) HasLocalTaskExecuteStatusValExecByPartyId(taskId, partyId string) (bool, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.HasLocalTaskExecuteStatusValExecByPartyId(dc.db, taskId, partyId)
}

// v 2.0 about local needExecuteTask
func (dc *DataCenter) StoreNeedExecuteTask(task *types.NeedExecuteTask) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreNeedExecuteTask(dc.db, task)
}

func (dc *DataCenter) RemoveNeedExecuteTaskByPartyId(taskId, partyId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveNeedExecuteTaskByPartyId(dc.db, taskId, partyId)
}

func (dc *DataCenter) RemoveNeedExecuteTask(taskId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveNeedExecuteTask(dc.db, taskId)
}

func (dc *DataCenter) ForEachNeedExecuteTaskWwithPrefix (prifix []byte, f func(key, value []byte) error) error {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.ForEachNeedExecuteTaskWwithPrefix(dc.db, prifix, f)
}

func (dc *DataCenter) ForEachNeedExecuteTask (f func(key, value []byte) error) error {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.ForEachNeedExecuteTask(dc.db, f)
}

// v 2.0 about taskbullet
func (dc *DataCenter) StoreTaskBullet(bullet *types.TaskBullet) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreTaskBullet(dc.db, bullet)
}

func (dc *DataCenter) RemoveTaskBullet(taskId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveTaskBullet(dc.db, taskId)
}

func (dc *DataCenter) ForEachTaskBullets(f func(key, value []byte) error) error {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.ForEachTaskBullets(dc.db, f)
}
