package scheduler

import (
	"container/heap"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	"github.com/RosettaFlow/Carrier-Go/types"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	StarveTerm                  = uint32(3)
	defaultScheduleTaskInterval = 20 * time.Millisecond
	electionCondition           = 10000
	taskComputeOrgCount         = 3
)

var (
	ErrEnoughResourceOrgCountLessCalculateCount = fmt.Errorf("the enough resource org count is less calculate count")
)

type DataCenter interface {
	GetIdentityList() (types.IdentityArray, error)
	GetRegisterNodeList(typ types.RegisteredNodeType) ([]*types.RegisteredNodeInfo, error)
	GetIdentity() (*types.NodeAlias, error)
	GetResourceList() (types.ResourceArray, error)
}
type SchedulerStarveFIFO struct {
	resourceMng *resource.Manager
	// the local task into this queue, first
	queue *types.TaskBullets
	// the very very starve local task by priority
	starveQueue *types.TaskBullets
	// the cache with scheduled local task, will be send to `Consensus`
	scheduledQueue []*types.ScheduleTask
	// fetch local task from taskManager`
	localTaskCh chan types.TaskMsgs
	// send local task scheduled to `Consensus`
	schedTaskCh chan *types.ConsensusTaskWrap
	// receive remote task to replay from `Consensus`
	remoteTaskCh chan *types.ScheduleTaskWrap
	dataCenter   DataCenter
	err          error
}

func NewSchedulerStarveFIFO(
	localTaskCh chan types.TaskMsgs, schedTaskCh chan *types.ConsensusTaskWrap,
	remoteTaskCh chan *types.ScheduleTaskWrap, dataCenter DataCenter,
	mng *resource.Manager) *SchedulerStarveFIFO {

	return &SchedulerStarveFIFO{
		resourceMng:    mng,
		queue:          new(types.TaskBullets),
		starveQueue:    new(types.TaskBullets),
		scheduledQueue: make([]*types.ScheduleTask, 0),
		localTaskCh:    localTaskCh,
		schedTaskCh:    schedTaskCh,
		remoteTaskCh:   remoteTaskCh,
		dataCenter:     dataCenter,
	}
}
func (sche *SchedulerStarveFIFO) loop() {
	//taskTimer := time.NewTimer(defaultScheduleTaskInterval)
	for {
		select {
		case tasks := <-sche.localTaskCh:

			for _, task := range tasks {
				bullet := types.NewTaskBullet(task)
				sche.addTaskBullet(bullet)
				sche.trySchedule()
			}

			//case task := <-sche.remoteTaskCh:
			// todo 让自己的Scheduler 重演选举

		}

		// todo 这里还需要写上 定时调度 队列中的任务信息
	}
}

func (sche *SchedulerStarveFIFO) OnStart() error {
	err := sche.resourceMng.Start()
	if nil != err {
		return err
	}
	go sche.loop()
	return nil
}
func (sche *SchedulerStarveFIFO) OnStop() error  { return nil }
func (sche *SchedulerStarveFIFO) OnError() error { return sche.err }
func (sche *SchedulerStarveFIFO) Name() string   { return "SchedulerStarveFIFO" }
func (sche *SchedulerStarveFIFO) addTaskBullet(bullet *types.TaskBullet) {
	heap.Push(sche.queue, bullet) //
}
func (sche *SchedulerStarveFIFO) trySchedule() error {

	if sche.starveQueue.Len() != 0 {
		x := heap.Pop(sche.starveQueue)
		task := x.(*types.TaskBullet).TaskMsg

		powers, err := sche.electionConputeOrg(taskComputeOrgCount, &types.TaskOperationCost{Mem:
			task.OperationCost().Mem, Processor: task.OperationCost().Processor,
			Bandwidth:  task.OperationCost().Bandwidth})
		if nil != err {
			log.Errorf("Failed to election power org, err: %s", err)
			return err
		}
		scheduleTask := buildScheduleTask(task, powers)
		go func() {
			resCh := make(chan *types.TaskConsResult, 0)
			sche.schedTaskCh <- &types.ConsensusTaskWrap{
				Task: scheduleTask,
				ResultCh: resCh,
			}
			//res := <- resCh
			// TODO 处理 共识的 任务信息
		}()
	}

	return nil
}

func (sche *SchedulerStarveFIFO) replaySchedule() error {

	return nil
}

func (sche *SchedulerStarveFIFO) electionConputeOrg(calculateCount int, cost *types.TaskOperationCost) ([]*types.NodeAlias, error) {

	slot := &types.Slot{Mem: cost.Mem, Processor: cost.Processor, Bandwidth: cost.Bandwidth}
	orgs := make([]*types.NodeAlias, 0)
	identityIds := make([]string, 0)

	for _, r := range sche.resourceMng.GetRemoteResouceTables() {
		if r.IsEnough(slot) {
			identityIds = append(identityIds, r.GetIdentityId())
			//identityIdTmp[r.GetIdentityId()] = struct{}{}
		}
	}
	if calculateCount > len(identityIds) {
		return nil, ErrEnoughResourceOrgCountLessCalculateCount
	}
	// Election
	index := electionCondition%len(identityIds)
	identityIdTmp := make(map[string]struct{}, 0)
	for i := calculateCount; i > 0; i-- {

		identityIdTmp[identityIds[index]] = struct{}{}
		index ++

	}

	identityArr, err := sche.dataCenter.GetIdentityList()
	if nil != err {
		return nil, err
	}
	for _, iden := range identityArr {
		if _, ok := identityIdTmp[iden.IdentityId()]; ok {
			orgs = append(orgs, &types.NodeAlias{
				Name: iden.Name(),
				NodeId: iden.NodeId(),
				IdentityId: iden.IdentityId(),
			})
		}
	}
	return orgs, nil
}

func buildScheduleTask (task *types.TaskMsg, powers []*types.NodeAlias) *types.ScheduleTask {

	partners := make([]*types.ScheduleTaskDataSupplier, len(task.PartnerTaskSuppliers()))
	for i, p := range task.PartnerTaskSuppliers() {
		partner := &types.ScheduleTaskDataSupplier{
			NodeAlias: &types.NodeAlias{
				Name: p.Name,
				NodeId: p.NodeId,
				IdentityId: p.IdentityId,
			},
			MetaData:  p.MetaData,
		}
		partners[i] = partner
	}

	powerArr := make([]*types.ScheduleTaskPowerSupplier, len(powers))
	for i, p := range powers {
		power := &types.ScheduleTaskPowerSupplier{
			NodeAlias: p,
		}
		powerArr[i] = power
	}

	receivers := make([]*types.ScheduleTaskResultReceiver, len(task.ReceiverDetails()))
	for i, r := range task.ReceiverDetails() {
		receiver := &types.ScheduleTaskResultReceiver{
			NodeAlias: r.NodeAlias,
			Providers: r.Providers,
		}
		receivers[i] = receiver
	}
	return &types.ScheduleTask{
		TaskId: task.TaskId,
		TaskName: task.TaskName(),
		Owner:   &types.ScheduleTaskDataSupplier{
			NodeAlias: task.Onwer(),
			MetaData:  task.OwnerTaskSupplier().MetaData,
		},
		Partners: partners,
		PowerSuppliers: powerArr,
		Receivers: receivers,
		CalculateContractCode: task.CalculateContractCode(),
		DataSplitContractCode: task.DataSplitContractCode(),
		OperationCost: task.OperationCost(),
		CreateAt: task.CreateAt(),
	}
}