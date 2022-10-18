package workflow

import (
	"encoding/json"
	"fmt"
	"github.com/datumtechs/datum-network-carrier/carrierdb"
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/datumtechs/datum-network-carrier/rpc/backend"
	"github.com/datumtechs/datum-network-carrier/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gogo/protobuf/proto"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	timeout                                               int64 = 72 * 3600
	defaultRemoveWorkflowExecuteResultSaveTimeoutInterval       = 30 * time.Second
	workflowTaskStatusCacheKeyPrefix                            = []byte("workflowTaskStatusCacheKeyPrefix:")
	workflowStatusCacheKeyPrefix                                = []byte("workflowStatusCacheKeyPrefix:")
)

type Manager struct {
	taskExecuteResultCh       chan *carrierapipb.WorkFlowTaskStatus // Trigger the channel when saveTask(InsertTask)
	TaskMsgToMessageManagerCh chan *types.TaskMsg                   // This channel is triggered when a task in the workflow executes successfully
	dataCenter                carrierdb.CarrierDB
	sendToTaskManagerCache    map[string]string                                      // {"task1":"workFlowId1","task2":"workFlowId1","task3":"workFlowId2"}
	workflowsCache            map[string]*types.Workflow                             // {"workflowId1":[],"workflowId2":[]}
	workflowStatusCache       map[string]*types.WorkflowStatus                       // {"workflowId1":{},"workflowId2":{}}
	workflowTaskStatusCache   map[string]map[string]*carrierapipb.WorkFlowTaskStatus // {"workflowId1":{"taskName":{}}}
	sendToTaskManagerLock     sync.RWMutex
	workflowsLock             sync.RWMutex
	workflowStatusLock        sync.RWMutex
	workflowTaskStatusLock    sync.RWMutex
	quit                      chan struct{}
}

func NewWorkflowService(
	db carrierdb.CarrierDB,
	taskExecuteResultCh chan *carrierapipb.WorkFlowTaskStatus,
	TaskMsgToMessageManagerCh chan *types.TaskMsg,
) *Manager {
	return &Manager{
		sendToTaskManagerCache:    make(map[string]string, 0),
		workflowsCache:            make(map[string]*types.Workflow, 0),
		workflowStatusCache:       make(map[string]*types.WorkflowStatus, 0),
		workflowTaskStatusCache:   make(map[string]map[string]*carrierapipb.WorkFlowTaskStatus, 0),
		taskExecuteResultCh:       taskExecuteResultCh,
		TaskMsgToMessageManagerCh: TaskMsgToMessageManagerCh,
		quit:                      make(chan struct{}),
		dataCenter:                db,
	}
}

func (m *Manager) AddWorkflow(workflow *types.Workflow) error {
	log.Debugf("AddWorkflow workflowId is:{%s}", workflow.WorkflowId)
	m.workflowsLock.RLock()
	defer m.workflowsLock.RUnlock()
	if workflow.GetWorkflowId() == "" {
		return fmt.Errorf("workflow name is %s,it's workflow id %s", workflow.GetWorkflowId(), workflow.GetWorkflowName())
	}
	err := m.taskMsgSendToMessageManager(workflow)
	if err != nil {
		return err
	}
	log.Debugf("AddWorkflow successful is:{%s}", workflow.WorkflowId)
	return nil
}

func (m *Manager) GetWorkflowStatus(workflowIds []string) (*carrierapipb.QueryWorkStatusResponse, error) {
	log.Debugf("GetWorkflowStatus workflowIds %v", workflowIds)
	workflowStatusList := make([]*carrierapipb.WorkFlowStatus, 0)
	m.workflowStatusLock.RLock()
	m.workflowTaskStatusLock.RLock()
	defer func() {
		m.workflowStatusLock.RUnlock()
		m.workflowTaskStatusLock.RUnlock()
	}()
	for _, workflowId := range workflowIds {
		if status, ok := m.workflowStatusCache[workflowId]; !ok {
			log.Errorf("no status information for {%s} was found in workflowStatusCache", workflowId)
		} else {
			taskStatusList := make([]*carrierapipb.WorkFlowTaskStatus, 0)
			if taskStatusMap, ok := m.workflowTaskStatusCache[workflowId]; !ok {
				log.Errorf("no status information for {%s} was found in workflowTaskStatusCache", workflowId)
			} else {
				for _, status := range taskStatusMap {
					taskStatusList = append(taskStatusList, &carrierapipb.WorkFlowTaskStatus{
						TaskId:   status.GetTaskId(),
						Status:   status.GetStatus(),
						TaskName: status.GetTaskName(),
					})
				}
			}
			workflowStatus := &carrierapipb.WorkFlowStatus{
				Status:   status.Status,
				TaskList: taskStatusList,
			}
			workflowStatusList = append(workflowStatusList, workflowStatus)
		}
	}
	log.Debugf("GetWorkflowStatus over workflowIds %v", workflowIds)
	return &carrierapipb.QueryWorkStatusResponse{
		Status:             0,
		Msg:                backend.OK,
		WorkflowStatusList: workflowStatusList,
	}, nil
}

func (m *Manager) loop() {
	checkExecuteResultTicker := time.NewTicker(defaultRemoveWorkflowExecuteResultSaveTimeoutInterval)
	for {
		select {
		case result := <-m.taskExecuteResultCh:
			log.Debugf("taskExecuteResultCh result is %v", result)
			m.sendToTaskManagerLock.RLock()
			workflowId := m.sendToTaskManagerCache[result.GetTaskId()]
			m.sendToTaskManagerLock.RUnlock()
			m.updateWorkflowTaskStatus(workflowId, result)
			switch result.GetStatus() {
			case commonconstantpb.TaskState_TaskState_Succeed:
				m.workflowsLock.RLock()
				workflow := m.workflowsCache[workflowId]
				if len(workflow.Tasks) == 0 {
					m.updateWorkflowStatus(workflowId, commonconstantpb.WorkFlowState_WorkFlowState_Succeed)
					m.DeleteWorkflowCache(workflowId)
				} else {
					if err := m.taskMsgSendToMessageManager(workflow); err != nil {
						log.Warnf("taskMsgSendToMessageManager fail,%s", err.Error())
					}
				}
				m.workflowsLock.RUnlock()
			case commonconstantpb.TaskState_TaskState_Failed:
				m.updateWorkflowStatus(workflowId, commonconstantpb.WorkFlowState_WorkFlowState_Failed)
				m.DeleteWorkflowCache(workflowId)
			}
		case <-checkExecuteResultTicker.C:
			m.removeWorkflowExecuteResultSaveTimeout()
		case <-m.quit:
			log.Info("Stopped workflowManager ...")
			return
		}
	}
}

func (m *Manager) Start() error {
	log.Info("Started workflowManager ...")
	m.removeWorkflowExecuteResultSaveTimeout()
	m.recoveryCache()
	log.Info("recoveryCache over ...")
	go m.loop()
	return nil
}

func (m *Manager) Stop() error {
	close(m.quit)
	return nil
}

func (m *Manager) sendTaskMsg(tm *types.TaskMsg, workflowId string) {
	log.Debugf("sendTaskMsg taskId:{%s},workflowId:{%s}", tm.GenTaskId(), workflowId)
	taskId := tm.GetTaskId()
	m.sendToTaskManagerLock.Lock()
	defer m.sendToTaskManagerLock.Unlock()
	if _, ok := m.sendToTaskManagerCache[taskId]; !ok {
		m.sendToTaskManagerCache[taskId] = workflowId
		if err := m.dataCenter.SaveSendToTaskManager(taskId, workflowId); err != nil {
			log.WithError(err).Errorf("sendTaskMsg SaveSendToTaskManager fail.")
		}
		m.sendToTaskMsg(tm)
	} else {
		log.Warnf("initWorkflowFirstTask taskId %s alerady exits in alreadySendTaskManager.", taskId)
	}
}

func (m *Manager) sendToTaskMsg(tm *types.TaskMsg) {
	m.TaskMsgToMessageManagerCh <- tm
}

func (m *Manager) updateWorkflowTaskStatus(workflowId string, status *carrierapipb.WorkFlowTaskStatus) {
	m.workflowTaskStatusLock.Lock()
	defer m.workflowTaskStatusLock.Unlock()
	if taskStatus, ok := m.workflowTaskStatusCache[workflowId]; !ok {
		log.Errorf("workflow update task status fail,workflowTaskStatusCache not exits workflowId %s", workflowId)
	} else {
		taskStatus[status.GetTaskName()].Status = status.GetStatus()
		m.workflowTaskStatusCache[workflowId] = taskStatus
		if err := m.dataCenter.SaveWorkflowTaskStatusCache(workflowId, status); err != nil {
			log.WithError(err).Errorf("updateWorkflowTaskStatus SaveWorkflowTaskStatusCache fail.")
		}
	}
}

func (m *Manager) updateWorkflowStatus(workflowId string, status commonconstantpb.WorkFlowState) {
	m.workflowStatusLock.Lock()
	m.sendToTaskManagerLock.RLock()
	defer func() {
		m.workflowStatusLock.Unlock()
		m.sendToTaskManagerLock.RUnlock()
	}()
	workflowStatus := &types.WorkflowStatus{
		Status:   status,
		UpdateAt: time.Now().Unix(),
	}
	m.workflowStatusCache[workflowId] = workflowStatus
	workflow, ok := m.workflowsCache[workflowId]
	if !ok {
		return
	}
	for _, v := range workflow.Tasks {
		delete(m.sendToTaskManagerCache, v.GetTaskId())
		if err := m.dataCenter.RemoveSendToTaskManager(v.GetTaskId()); err != nil {
			log.WithError(err).Errorf("updateWorkflowStatus RemoveSendToTaskManager fail.")
		}
	}
	if err := m.dataCenter.SaveWorkflowStatusCache(workflowId, workflowStatus); err != nil {
		log.WithError(err).Errorf("updateWorkflowStatus SaveWorkflowStatusCache fail.")
	}
}

func (m *Manager) taskMsgSendToMessageManager(workflow *types.Workflow) error {
	workflowId := workflow.GetWorkflowId()

	m.workflowStatusLock.Lock()
	m.workflowTaskStatusLock.Lock()
	defer func() {
		m.workflowStatusLock.Unlock()
		m.workflowTaskStatusLock.Unlock()
	}()
	if _, ok := m.workflowsCache[workflowId]; ok {
		return fmt.Errorf("workflowId %s alerady exits workflowsCache", workflowId)
	}
	// equivalent to pop and check defer to task execute result
	task, err := m.assemblyTaskParameters(workflow)
	if err != nil {
		return err
	}
	m.sendTaskMsg(task, workflowId)
	workflowStatus := &types.WorkflowStatus{
		Status:   commonconstantpb.WorkFlowState_WorkFlowState_Running,
		UpdateAt: time.Now().Unix(),
	}
	m.workflowStatusCache[workflowId] = workflowStatus
	if err := m.dataCenter.SaveWorkflowStatusCache(workflowId, workflowStatus); err != nil {
		log.WithError(err).Errorf("taskMsgSendToMessageManager SaveWorkflowStatusCache fail.")
	}
	//update workflowTaskStatus
	workflowTaskStatus := make(map[string]*carrierapipb.WorkFlowTaskStatus, 0)
	statusWorkflowTask := &carrierapipb.WorkFlowTaskStatus{
		Status:   task.GetState(),
		TaskId:   task.GetTaskId(),
		TaskName: task.GetTaskName(),
	}
	workflowTaskStatus[task.GetTaskName()] = statusWorkflowTask
	m.workflowTaskStatusCache[workflowId] = workflowTaskStatus
	if err := m.dataCenter.SaveWorkflowTaskStatusCache(workflowId, statusWorkflowTask); err != nil {
		log.WithError(err).Errorf("updateWorkflowTaskStatus SaveWorkflowTaskStatusCache fail.")
	}
	// update workflow Tasks
	workflow.Tasks = workflow.Tasks[1:]
	m.workflowsCache[workflow.GetWorkflowId()] = workflow
	taskList := make([]*carriertypespb.TaskMsg, 0)
	for _, v := range workflow.Tasks {
		taskList = append(taskList, &carriertypespb.TaskMsg{
			Data: v.GetTaskData(),
		})
	}
	if err := m.dataCenter.SaveWorkflowCache(&carriertypespb.Workflow{
		WorkflowId:   workflow.GetWorkflowId(),
		Desc:         workflow.Desc,
		WorkflowName: workflow.GetWorkflowName(),
		PolicyType:   workflow.PolicyType,
		Policy:       workflow.Policy,
		User:         workflow.User,
		UserType:     workflow.UserType,
		Sign:         workflow.Sign,
		Tasks:        taskList,
		CreateAt:     workflow.CreateAt,
	}); err != nil {
		log.WithError(err).Errorf("taskMsgSendToMessageManager SaveWorkflowCache fail.")
	}
	return nil
}

func (m *Manager) assemblyTaskParameters(workflow *types.Workflow) (*types.TaskMsg, error) {
	task := workflow.Tasks[0]
	task.GenTaskId()
	switch workflow.PolicyType {
	case commonconstantpb.WorkFlowPolicyType_Ordinary_Policy:
		// Check if the current task has dependencies
		var wp *types.WorkflowPolicy
		if err := json.Unmarshal([]byte(workflow.Policy), &wp); err != nil {
			return nil, err
		}
		for _, v := range *wp {
			if v.Origin == task.GetTaskName() {
				if len(v.Reference) == 0 {
					return task, nil
				}
				for _, ref := range v.Reference {
					if len(ref.DependParams) == 0 {
						return task, nil
					}
					referToTaskId := m.getWorkflowTaskStatusCacheTaskId(workflow.WorkflowId, ref.Target)
					if referToTaskId == "" {
						log.Warnf("getWorkflowTaskStatusCacheTaskId get result is empty,workflowId {%s},taskName {%s}", workflow.WorkflowId, ref.Target)
						continue
					}
					for index, params := range ref.DependParams {
						var (
							dataPolicyOption string
							dataPolicyType   uint32
						)
						dependParamsType := ref.DependParamsType[index]
						switch dependParamsType {
						case types.TASK_DATA_POLICY_CSV_WITH_TASKRESULTDATA:
							var p *types.PSIParams
							if err := json.Unmarshal([]byte(params), p); err != nil {
								log.WithError(err).Errorf("json Unmarshal fail,%s, dependParamsType:%d", params, dependParamsType)
								continue
							}
							taskResultParams := &types.TaskMetadataPolicyCSVWithTaskResultData{
								PartyId:             fmt.Sprintf("data%s", strconv.FormatInt(time.Now().Unix(), 10)),
								TaskId:              referToTaskId,
								InputType:           p.InputType,
								KeyColumnName:       p.KeyColumnName,
								SelectedColumnNames: p.SelectedColumnNames,
							}
							result, _ := json.Marshal(taskResultParams)
							dataPolicyOption = string(result)
							dataPolicyType = types.TASK_DATA_POLICY_CSV_WITH_TASKRESULTDATA
						case types.TASK_DATA_POLICY_DIR:
							var p *types.MODELParams
							if err := json.Unmarshal([]byte(params), p); err != nil {
								log.WithError(err).Errorf("json Unmarshal fail,%s, dependParamsType:%d", params, dependParamsType)
								continue
							}
							resultDataSummary, err := m.dataCenter.QueryTaskResultDataSummary(referToTaskId)
							if err != nil {
								log.WithError(err).Errorf("query task {%s} resultDataSummary fail!", referToTaskId)
								continue
							}
							taskResultParams := &types.TaskMetadataPolicyDIR{
								PartyId:      fmt.Sprintf("data%s", strconv.FormatInt(time.Now().Unix(), 10)),
								MetadataId:   resultDataSummary.GetMetadataId(),
								MetadataName: resultDataSummary.GetMetadataName(),
								InputType:    p.InputType,
							}
							result, _ := json.Marshal(taskResultParams)
							dataPolicyOption = string(result)
							dataPolicyType = types.TASK_DATA_POLICY_DIR
						default:
							log.Errorf("assemblyTaskParameters unknown params type")
						}
						if dataPolicyOption != "" {
							task.GetTaskData().DataPolicyTypes = append(task.GetTaskData().DataPolicyTypes, dataPolicyType)
							task.GetTaskData().DataPolicyOptions = append(task.GetTaskData().DataPolicyOptions, dataPolicyOption)
						}
					}
				}
			}
		}
	default:
		return nil, fmt.Errorf("unknown workflow policy type %s", workflow.PolicyType.String())
	}
	return task, nil
}

func (m *Manager) getWorkflowTaskStatusCacheTaskId(workflowId, taskName string) string {
	m.workflowTaskStatusLock.RLock()
	defer m.workflowTaskStatusLock.RUnlock()
	if workflowTaskStatus, ok := m.workflowTaskStatusCache[workflowId]; !ok {
		return ""
	} else {
		if status, ok := workflowTaskStatus[taskName]; !ok {
			return ""
		} else {
			return status.GetTaskId()
		}
	}
}

func (m *Manager) DeleteWorkflowCache(workflowId string) {
	m.workflowsLock.RLock()
	delete(m.workflowsCache, workflowId)
	m.workflowsLock.RUnlock()
	if err := m.dataCenter.RemoveWorkflowCache(workflowId); err != nil {
		log.WithError(err).Errorf("DeleteWorkflowCache dataCenter.RemoveWorkflowCache fail")
	}
}
func (m *Manager) recoveryCache() {
	errCh := make(chan error, 4)
	var wg sync.WaitGroup
	wg.Add(4)

	// recovery sendToTaskManagerCache
	go func(wg *sync.WaitGroup, errCh chan<- error) {
		defer wg.Done()
		sendToTaskManagerCacheKeyPrefix := []byte("sendToTaskManagerCacheKeyPrefix:")
		prefixLength := len(sendToTaskManagerCacheKeyPrefix)
		if err := m.dataCenter.ForEachKVWithPrefix(sendToTaskManagerCacheKeyPrefix, func(key, value []byte) error {
			taskId := string(key[prefixLength:])
			log.Debugf("recovery sendToTaskManagerCache taskId {%s}", taskId)
			if len(key) != 0 && len(value) != 0 {
				var workflowId string
				if err := rlp.DecodeBytes(value, &workflowId); err != nil {
					return err
				} else {
					m.sendToTaskManagerCache[taskId] = workflowId
				}
			}
			return nil
		}); err != nil {
			errCh <- err
			return
		}
	}(&wg, errCh)
	// recovery workflowsCache
	go func(wg *sync.WaitGroup, errCh chan<- error) {
		defer wg.Done()
		workflowsCacheKeyPrefix := []byte("workflowsCacheKeyPrefix:")
		prefixLength := len(workflowsCacheKeyPrefix)
		if err := m.dataCenter.ForEachKVWithPrefix(workflowsCacheKeyPrefix, func(key, value []byte) error {
			workflowId := string(key[prefixLength:])
			log.Debugf("reovery workflowsCache workflowId {%s}", workflowId)
			var workflow carriertypespb.Workflow
			if len(key) != 0 && len(value) != 0 {
				if err := proto.Unmarshal(value, &workflow); err != nil {
					return err
				} else {
					taskList := make([]*types.TaskMsg, 0)
					for _, v := range workflow.Tasks {
						taskList = append(taskList, &types.TaskMsg{Data: types.NewTask(v.GetData())})
					}

					m.workflowsCache[workflowId] = &types.Workflow{
						WorkflowId:   workflow.GetWorkflowId(),
						Desc:         workflow.GetDesc(),
						WorkflowName: workflow.GetWorkflowName(),
						PolicyType:   workflow.GetPolicyType(),
						Policy:       workflow.GetPolicy(),
						User:         workflow.GetUser(),
						UserType:     workflow.GetUserType(),
						Sign:         workflow.GetSign(),
						Tasks:        taskList,
						CreateAt:     workflow.GetCreateAt(),
					}
				}
			}
			return nil
		}); err != nil {
			errCh <- err
			return
		}
	}(&wg, errCh)
	// recovery workflowStatusCache
	go func(wg *sync.WaitGroup, errCh chan<- error) {
		defer wg.Done()
		prefixLength := len(workflowStatusCacheKeyPrefix)
		if err := m.dataCenter.ForEachKVWithPrefix(workflowStatusCacheKeyPrefix, func(key, value []byte) error {
			workflowId := string(key[prefixLength:])
			log.Debugf("recovery workflowStatusCache workflowId {%s}", workflowId)
			var workflowState *types.WorkflowStatus
			if len(key) != 0 && len(value) != 0 {
				if err := json.Unmarshal(value, &workflowState); err != nil {
					return err
				} else {
					m.workflowStatusCache[workflowId] = workflowState
				}
			}
			return nil
		}); err != nil {
			errCh <- err
			return
		}
	}(&wg, errCh)
	// recovery workflowTaskStatusCache
	go func(wg *sync.WaitGroup, errCh chan<- error) {
		defer wg.Done()
		prefixLength := len(workflowTaskStatusCacheKeyPrefix)
		if err := m.dataCenter.ForEachKVWithPrefix(workflowTaskStatusCacheKeyPrefix, func(key, value []byte) error {
			workflowIdTaskName := string(key[prefixLength:])
			workflowId, taskName := workflowIdTaskName[:75], workflowIdTaskName[75:]
			log.Debugf("recovery workflowTaskStatusCache,workflowId{%s},taskName{%s}", workflowId, taskName)
			var taskState carrierapipb.WorkFlowTaskStatus
			if len(key) != 0 && len(value) != 0 {
				if err := proto.Unmarshal(value, &taskState); err != nil {
					return err
				} else {
					result, ok := m.workflowTaskStatusCache[workflowId]
					if !ok {
						result = make(map[string]*carrierapipb.WorkFlowTaskStatus, 0)
					}
					result[taskName] = &taskState
					m.workflowTaskStatusCache[workflowId] = result
				}
			}
			return nil
		}); err != nil {
			errCh <- err
			return
		}
	}(&wg, errCh)

	wg.Wait()
	close(errCh)

	errStrs := make([]string, 0)

	for err := range errCh {
		if nil != err {
			errStrs = append(errStrs, err.Error())
		}
	}
	if len(errStrs) != 0 {
		log.Fatalf("recover workflow state failed: \n%s", strings.Join(errStrs, "\n"))
	}
}

func (m *Manager) removeWorkflowExecuteResultSaveTimeout() {
	prefixLength := len(workflowStatusCacheKeyPrefix)
	saveRemoveWorkflowIds := make(map[string]struct{}, 0)
	// remove workflowStatusCache
	if err := m.dataCenter.ForEachKVWithPrefix(workflowStatusCacheKeyPrefix, func(key, value []byte) error {
		workflowId := string(key[prefixLength:])
		log.Debugf("recovery workflowStatusCache workflowId {%s}", workflowId)
		var workflowState *types.WorkflowStatus
		if len(key) != 0 && len(value) != 0 {
			if err := json.Unmarshal(value, &workflowState); err != nil {
				log.WithError(err).Errorf("removeWorkflowExecuteResultSaveTimeout json.Unmarshal workflowState fail")
			} else {
				if (time.Now().Unix()-workflowState.UpdateAt) > timeout && (workflowState.Status == commonconstantpb.WorkFlowState_WorkFlowState_Running || workflowState.Status == commonconstantpb.WorkFlowState_WorkFlowState_Failed) {
					saveRemoveWorkflowIds[workflowId] = struct{}{}
					return m.dataCenter.RemoveWorkflowStatusCache(workflowId)
				}
			}
		}
		return nil
	}); err != nil {
		log.WithError(err).Errorf("removeWorkflowExecuteResultSaveTimeout workflowStatusCache fail")
	}
	//remove workflowTaskStatusCache
	prefixLength = len(workflowTaskStatusCacheKeyPrefix)
	if err := m.dataCenter.ForEachKVWithPrefix(workflowTaskStatusCacheKeyPrefix, func(key, value []byte) error {
		workflowIdTaskName := string(key[prefixLength:])
		workflowId, taskName := workflowIdTaskName[:75], workflowIdTaskName[75:]
		log.Debugf("recovery workflowTaskStatusCache,workflowId{%s},taskName{%s}", workflowId, taskName)
		var taskState carrierapipb.WorkFlowTaskStatus
		if len(key) != 0 && len(value) != 0 {
			if err := proto.Unmarshal(value, &taskState); err != nil {
				log.WithError(err).Errorf("removeWorkflowExecuteResultSaveTimeout json.Unmarshal workflowState fail")
			} else {
				workflowIdTaskName := string(key[prefixLength:])
				workflowId := workflowIdTaskName[:75]
				if _, ok := saveRemoveWorkflowIds[workflowId]; ok {
					if err := m.dataCenter.RemoveWorkflowTaskStatusCache(workflowIdTaskName); err != nil {
						log.Warnf("removeWorkflowExecuteResultSaveTimeout RemoveWorkflowTaskStatusCache fail,workflowIdTaskName %s", workflowIdTaskName)
					}
				}
			}
		}
		return nil
	}); err != nil {
		log.WithError(err).Errorf("removeWorkflowExecuteResultSaveTimeout workflowTaskStatusCache fail")
	}
}
