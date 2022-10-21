package workflow

import (
	"encoding/json"
	"fmt"
	"github.com/datumtechs/datum-network-carrier/carrierdb"
	"github.com/datumtechs/datum-network-carrier/carrierdb/rawdb"
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/datumtechs/datum-network-carrier/rpc/backend"
	"github.com/datumtechs/datum-network-carrier/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gogo/protobuf/proto"
	"strings"
	"sync"
	"time"
)

var (
	timeout                                               int64 = 72 * 3600 // second
	defaultRemoveWorkflowExecuteResultSaveTimeoutInterval       = 60 * time.Second
)

type Manager struct {
	taskExecuteResultCh       chan *carrierapipb.WorkFlowTaskStatus // Trigger the channel when saveTask(InsertTask)
	TaskMsgToMessageManagerCh chan *types.TaskMsg                   // This channel is triggered when a task in the workflow executes successfully
	dataCenter                carrierdb.CarrierDB
	//key=>taskId,value=>workflowId
	//eg:{"taskId1":"workFlowId1","taskId2":"workFlowId1","taskId3":"workFlowId2"}
	sendToTaskManagerCache  map[string]string
	workflowsCache          map[string]*types.Workflow                             // {"workflowId1":[],"workflowId2":[]}
	workflowStatusCache     map[string]*types.WorkflowStatus                       // {"workflowId1":{},"workflowId2":{}}
	workflowTaskStatusCache map[string]map[string]*carrierapipb.WorkFlowTaskStatus // {"workflowId1":{"taskName":{}}}
	sendToTaskManagerLock   sync.RWMutex
	workflowsLock           sync.RWMutex
	workflowStatusLock      sync.RWMutex
	workflowTaskStatusLock  sync.RWMutex
	quit                    chan struct{}
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
	if _, ok := m.workflowsCache[workflow.WorkflowId]; ok {
		return fmt.Errorf("AddWorkflow WorkflowId {%s} alerady exits", workflow.WorkflowId)
	}
	m.workflowsLock.RUnlock()
	if workflow.GetWorkflowId() == "" {
		return fmt.Errorf("workflow name is %s,it's workflow id %s", workflow.GetWorkflowId(), workflow.GetWorkflowName())
	}
	return m.taskMsgSendToMessageManager(workflow, true)
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
						StartAt:  status.GetStartAt(),
						EndAt:    status.GetEndAt(),
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

func (m *Manager) GetAllWorkflowDetails() (*carrierapipb.QueryWorkflowDetailsResponse, error) {
	workflowList := make([]*carriertypespb.Workflow, 0)
	workflowsBackupCacheKeyPrefix := rawdb.QueryWorkflowsBackupCacheKeyPrefix()
	prefixLength := len(workflowsBackupCacheKeyPrefix)
	if err := m.dataCenter.ForEachKVWithPrefix(workflowsBackupCacheKeyPrefix, func(key, value []byte) error {
		workflowId := string(key[prefixLength:])
		log.Debugf("reovery workflowsBackupCache workflowId {%s}", workflowId)
		var workflow carriertypespb.Workflow
		if len(key) != 0 && len(value) != 0 {
			if err := proto.Unmarshal(value, &workflow); err != nil {
				return err
			} else {
				taskList := make([]*carriertypespb.TaskMsg, 0)
				for _, v := range workflow.Tasks {
					taskList = append(taskList, &carriertypespb.TaskMsg{Data: v.GetData()})
				}

				workflowList = append(workflowList, &carriertypespb.Workflow{
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
				})
			}
		}
		return nil
	}); err != nil {
		log.WithError(err).Errorf("GetAllWorkflowDetails ForEachKVWithPrefix fail.")
	}
	return &carrierapipb.QueryWorkflowDetailsResponse{
		WorkflowDetailsList: workflowList,
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

			m.workflowsLock.RLock()
			workflow := m.workflowsCache[workflowId]
			m.workflowsLock.RUnlock()
			switch result.GetStatus() {
			case commonconstantpb.TaskState_TaskState_Succeed:
				if len(workflow.Tasks) == 0 {
					m.updateWorkflowStatus(workflowId, commonconstantpb.WorkFlowState_WorkFlowState_Succeed)
					m.DeleteWorkflowCache(workflowId)
				} else {
					if err := m.taskMsgSendToMessageManager(workflow, false); err != nil {
						log.Warnf("taskMsgSendToMessageManager fail,%s", err.Error())
					}
				}
			case commonconstantpb.TaskState_TaskState_Failed:
				// a=>b, c=>d
				// tasks = {b,d,a,c} or {b,a,c,d}
				if len(workflow.Tasks) != 0 {
					log.Debugf("m.taskExecuteResultCh status is fail,workflowId %s", workflowId)
					m.adjustTasksList(workflowId, result.TaskName)
				}
				if len(workflow.Tasks) != 0 {
					if err := m.taskMsgSendToMessageManager(workflow, false); err != nil {
						log.Warnf("taskMsgSendToMessageManager fail,%s", err.Error())
					}
				} else {
					m.updateWorkflowStatus(workflowId, commonconstantpb.WorkFlowState_WorkFlowState_Failed)
					m.DeleteWorkflowCache(workflowId)
				}
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
	if taskStatus, ok := m.workflowTaskStatusCache[workflowId]; ok {
		taskStatus[status.GetTaskName()] = status
		m.workflowTaskStatusCache[workflowId] = taskStatus
	} else {
		taskStatus = make(map[string]*carrierapipb.WorkFlowTaskStatus, 0)
		taskStatus[status.GetTaskName()] = status
		m.workflowTaskStatusCache[workflowId] = taskStatus
	}
	if err := m.dataCenter.SaveWorkflowTaskStatusCache(workflowId, status); err != nil {
		log.WithError(err).Errorf("updateWorkflowTaskStatus SaveWorkflowTaskStatusCache fail.")
	}
}

func (m *Manager) updateWorkflowStatus(workflowId string, status commonconstantpb.WorkFlowState) {
	m.workflowStatusLock.Lock()
	m.sendToTaskManagerLock.RLock()
	m.workflowsLock.RLock()
	defer func() {
		m.workflowStatusLock.Unlock()
		m.sendToTaskManagerLock.RUnlock()
		m.workflowsLock.RUnlock()
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

func (m *Manager) taskMsgSendToMessageManager(workflow *types.Workflow, isFirst bool) error {
	workflowId := workflow.GetWorkflowId()
	// equivalent to pop and check defer to task execute result
	task, err := m.assemblyTaskParameters(workflow)
	if err != nil {
		return err
	}
	m.workflowStatusLock.Lock()
	m.workflowTaskStatusLock.Lock()
	m.workflowsLock.Lock()
	defer func() {
		m.workflowStatusLock.Unlock()
		m.workflowTaskStatusLock.Unlock()
		m.workflowsLock.Unlock()
	}()

	m.sendTaskMsg(task, workflowId)
	workflowStatus := &types.WorkflowStatus{
		Status:   commonconstantpb.WorkFlowState_WorkFlowState_Running,
		UpdateAt: time.Now().Unix(),
	}
	m.workflowStatusCache[workflowId] = workflowStatus
	if err := m.dataCenter.SaveWorkflowStatusCache(workflowId, workflowStatus); err != nil {
		log.WithError(err).Errorf("taskMsgSendToMessageManager SaveWorkflowStatusCache fail.")
	}

	statusWorkflowTask := &carrierapipb.WorkFlowTaskStatus{
		Status:   task.GetState(),
		TaskId:   task.GetTaskId(),
		TaskName: task.GetTaskName(),
		StartAt:  task.GetStartAt(),
		EndAt:    task.GetEndAt(),
	}
	if statusDetails, ok := m.workflowTaskStatusCache[workflowId]; ok {
		statusDetails[task.GetTaskName()] = statusWorkflowTask
		m.workflowTaskStatusCache[workflowId] = statusDetails
	} else {
		statusDetails = make(map[string]*carrierapipb.WorkFlowTaskStatus, 0)
		statusDetails[task.GetTaskName()] = statusWorkflowTask
		m.workflowTaskStatusCache[workflowId] = statusDetails
	}
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
	workflowPB := &carriertypespb.Workflow{
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
	}
	if err := m.dataCenter.SaveWorkflowCache(workflowPB); err != nil {
		log.WithError(err).Errorf("taskMsgSendToMessageManager SaveWorkflowCache fail.")
	}
	if isFirst {
		if err := m.dataCenter.SaveWorkflowCacheBackup(workflowPB); err != nil {
			log.WithError(err).Errorf("taskMsgSendToMessageManager SaveWorkflowCacheBackup fail.")
		}
	}
	return nil
}

func (m *Manager) assemblyTaskParameters(workflow *types.Workflow) (*types.TaskMsg, error) {
	task := workflow.Tasks[0]
	task.GenTaskId()
	log.Infof("assemblyTaskParameters begin taskId {%s},DataPolicyTypes {%v},DataPolicyOptions %s", task.GetTaskId(), task.GetTaskData().GetDataPolicyTypes(), task.GetTaskData().GetDataPolicyOptions())
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
							partyId          string
						)
						dependParamsType := ref.DependParamsType[index]
						switch dependParamsType {
						case types.TASK_DATA_POLICY_CSV_WITH_TASKRESULTDATA:
							var taskResultParams *types.TaskMetadataPolicyCSVWithTaskResultData
							if err := json.Unmarshal([]byte(params), &taskResultParams); err != nil {
								log.WithError(err).Errorf("json Unmarshal fail,%s, dependParamsType:%d", params, dependParamsType)
								return nil, err
							}
							taskResultParams.TaskId = referToTaskId
							result, _ := json.Marshal(taskResultParams)
							dataPolicyOption = string(result)
							dataPolicyType = types.TASK_DATA_POLICY_CSV_WITH_TASKRESULTDATA
							partyId = taskResultParams.GetPartyId()
						case types.TASK_DATA_POLICY_DIR_WITH_TASKRESULTDATA:
							var taskResultParams *types.TaskMetadataPolicyDIRWithTaskResultData
							if err := json.Unmarshal([]byte(params), &taskResultParams); err != nil {
								log.WithError(err).Errorf("json Unmarshal fail,%s, dependParamsType:%d", params, dependParamsType)
								return nil, err
							}
							taskResultParams.TaskId = referToTaskId
							result, _ := json.Marshal(taskResultParams)
							dataPolicyOption = string(result)
							dataPolicyType = types.TASK_DATA_POLICY_DIR_WITH_TASKRESULTDATA
							partyId = taskResultParams.GetPartyId()
						default:
							return nil, fmt.Errorf("assemblyTaskParameters unknown params type,it is %d", dependParamsType)
						}
						if partyId == "" {
							return nil, fmt.Errorf("DependParams include empty partyId, %s", params)
						}
						if dataPolicyOption != "" {
							// Check if partyId exists in DataSuppliers
							if ok := m.checkDependParamsPartyId(partyId, task.GetTaskData().DataSuppliers); !ok {
								return nil, fmt.Errorf("partyId %s does not exist in DataSuppliers,params %s", partyId, params)
							}
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
	if len(task.GetTaskData().GetDataPolicyTypes()) != len(task.GetTaskData().GetDataPolicyOptions()) || len(task.GetTaskData().GetDataPolicyTypes()) != len(task.GetTaskData().GetDataSuppliers()) {
		return nil, fmt.Errorf("assemblyTaskParameters fail")
	}
	log.Infof("assemblyTaskParameters over taskId {%s},DataPolicyTypes {%v},DataPolicyOptions %s", task.GetTaskId(), task.GetTaskData().GetDataPolicyTypes(), task.GetTaskData().GetDataPolicyOptions())
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
	m.workflowsLock.Lock()
	delete(m.workflowsCache, workflowId)
	m.workflowsLock.Unlock()
	if err := m.dataCenter.RemoveWorkflowCache(workflowId); err != nil {
		log.WithError(err).Errorf("DeleteWorkflowCache dataCenter.RemoveWorkflowCache fail")
	}
	if err := m.dataCenter.RemoveWorkflowCacheBackup(workflowId); err != nil {
		log.WithError(err).Errorf("DeleteWorkflowCache dataCenter.RemoveWorkflowCacheBackup fail")
	}
}
func (m *Manager) recoveryCache() {
	errCh := make(chan error, 4)
	var wg sync.WaitGroup
	wg.Add(4)

	// recovery sendToTaskManagerCache
	go func(wg *sync.WaitGroup, errCh chan<- error) {
		defer wg.Done()
		sendToTaskManagerCacheKeyPrefix := rawdb.QuerySendToTaskManagerCacheKeyPrefix()
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
		workflowsCacheKeyPrefix := rawdb.QueryWorkflowsCacheKeyPrefix()
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
		workflowStatusCacheKeyPrefix := rawdb.QueryWorkflowStatusCacheKeyPrefix()
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
		workflowTaskStatusCacheKeyPrefix := rawdb.QueryWorkflowTaskStatusCacheKeyPrefix()
		prefixLength := len(workflowTaskStatusCacheKeyPrefix)
		if err := m.dataCenter.ForEachKVWithPrefix(workflowTaskStatusCacheKeyPrefix, func(key, value []byte) error {
			workflowIdTaskName := string(key[prefixLength:])
			// workflow:${workflowId hex} == len(7 + 2 + 64) == "workflowId:" + "0x" + "e33...fe4"
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
	workflowStatusCacheKeyPrefix := rawdb.QueryWorkflowStatusCacheKeyPrefix()
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
				now := time.Now().Unix()
				if (now-workflowState.UpdateAt) > timeout && (workflowState.Status == commonconstantpb.WorkFlowState_WorkFlowState_Succeed || workflowState.Status == commonconstantpb.WorkFlowState_WorkFlowState_Failed) {
					log.Warnf("workflowId save status time,%s,now is {%d},workflowState.UpdateAt {%d}", workflowId, now, workflowState.UpdateAt)
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
	workflowTaskStatusCacheKeyPrefix := rawdb.QueryWorkflowTaskStatusCacheKeyPrefix()
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

func (m *Manager) adjustTasksList(workflowId, taskName string) {
	m.workflowsLock.RLock()
	defer m.workflowsLock.RUnlock()
	workflow := m.workflowsCache[workflowId]
	var wp *types.WorkflowPolicy
	_ = json.Unmarshal([]byte(workflow.Policy), &wp)

	notExecuteTask := make([]string, 0)
	for _, v := range *wp {
		for _, vv := range v.Reference {
			if taskName == vv.Target {
				notExecuteTask = append(notExecuteTask, v.Origin)
			}
		}
	}

	for _, taskName := range notExecuteTask {
		j := 0
		for _, task := range workflow.Tasks {
			if taskName != task.GetTaskName() {
				workflow.Tasks[j] = task
				j++
			}
		}
		workflow.Tasks = workflow.Tasks[:j]
	}

}

func (m *Manager) checkDependParamsPartyId(partyId string, dataSuppliers []*carriertypespb.TaskOrganization) bool {
	for _, value := range dataSuppliers {
		if partyId == value.PartyId {
			return true
		}
	}
	return false
}
