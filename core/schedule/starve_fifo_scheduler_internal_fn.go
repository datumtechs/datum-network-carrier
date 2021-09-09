package schedule

import (
	"container/heap"
	"errors"
	twopctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"

	log "github.com/sirupsen/logrus"
	"strings"
)

func (sche *SchedulerStarveFIFO) pushTaskBullet(bullet *types.TaskBullet) error {
	sche.queueMutex.Lock()
	defer sche.queueMutex.Unlock()
	// The bullet is first into queue
	old, ok := sche.schedulings[bullet.TaskId]
	if !ok {
		heap.Push(sche.queue, bullet)
		return nil
	}

	bullet = old
	delete(sche.schedulings, bullet.TaskId)

	if bullet.Resched >= ReschedMaxCount {
		return ErrRescheduleLargeThreshold
	} else {
		if bullet.Starve {
			log.Debugf("GetTask repush  into starve queue, taskId: {%s}, reschedCount: {%d}, max threshold: {%d}",
				bullet.TaskId, bullet.Resched, ReschedMaxCount)
			heap.Push(sche.starveQueue, bullet)
		} else {
			log.Debugf("GetTask repush  into queue, taskId: {%s}, reschedCount: {%d}, max threshold: {%d}",
				bullet.TaskId, bullet.Resched, ReschedMaxCount)
			heap.Push(sche.queue, bullet)
		}
	}
	return nil
}

func (sche *SchedulerStarveFIFO) removeTaskBullet(taskId string) error {
	sche.queueMutex.Lock()
	defer sche.queueMutex.Unlock()

	// traversal the queue to remove task bullet, first.
	i := 0
	for {
		if i == sche.queue.Len() {
			break
		}
		bullet := (*(sche.queue))[i]

		// When found the bullet with taskId, removed it from queue.
		if bullet.TaskId == taskId {
			heap.Remove(sche.queue, i)
			delete(sche.schedulings, taskId)
			return nil // todo 这里需要做一次 持久化
		}
		(*(sche.queue))[i] = bullet
		i++
	}

	// otherwise, traversal the starveQueue to remove task bullet, second.
	i = 0
	for {
		if i == sche.starveQueue.Len() {
			break
		}
		bullet := (*(sche.starveQueue))[i]

		// When found the bullet with taskId, removed it from starveQueue.
		if bullet.TaskId == taskId {
			heap.Remove(sche.starveQueue, i)
			delete(sche.schedulings, taskId)
			return nil // todo 这里需要做一次 持久化
		}
		(*(sche.starveQueue))[i] = bullet
		i++
	}
	return nil
}

func (sche *SchedulerStarveFIFO) popTaskBullet() (*types.TaskBullet, error) {
	sche.queueMutex.Lock()
	defer sche.queueMutex.Unlock()

	var bullet *types.TaskBullet

	if sche.starveQueue.Len() != 0 {
		x := heap.Pop(sche.starveQueue)
		bullet = x.(*types.TaskBullet)
	} else {
		if sche.queue.Len() != 0 {
			x := heap.Pop(sche.queue)
			bullet = x.(*types.TaskBullet)
		} else {
			return nil, nil
		}
	}
	bullet.IncreaseResched()
	sche.schedulings[bullet.TaskId] = bullet

	return bullet, nil
}



func (sche *SchedulerStarveFIFO) increaseTaskTerm() {
	// handle starve queue
	sche.starveQueue.IncreaseTerm()

	// handle queue
	i := 0
	for {
		if i == sche.queue.Len() {
			return
		}
		bullet := (*(sche.queue))[i]
		bullet.IncreaseTerm()

		// When the task in the queue meets hunger, it will be transferred to starveQueue
		if bullet.Term >= StarveTerm {
			bullet.Starve = true
			heap.Push(sche.starveQueue, bullet)
			heap.Remove(sche.queue, i)
			i = 0
			continue
		}
		(*(sche.queue))[i] = bullet
		i++
	}
}

func (sche *SchedulerStarveFIFO) electionComputeNode(needSlotCount uint64) (*pb.YarnRegisteredPeerDetail, error) {

	if nil == sche.internalNodeSet || 0 == sche.internalNodeSet.JobNodeClientSize() {
		return nil, errors.New("not found alive jobNode")
	}

	resourceNodeIdArr := make([]string, 0)

	tables, err := sche.resourceMng.GetLocalResourceTables()
	if nil != err {
		return nil, err
	}
	log.Debugf("GetLocalResourceTables on electionConputeNode, localResources: %s", utilLocalResourceArrString(tables))
	for _, r := range tables {
		if r.IsEnough(uint32(needSlotCount)) {

			jobNodeClient, find := sche.internalNodeSet.QueryJobNodeClient(r.GetNodeId())
			if find && jobNodeClient.IsConnected() {
				resourceNodeIdArr = append(resourceNodeIdArr, r.GetNodeId())
			}
		}
	}

	if len(resourceNodeIdArr) == 0 {
		return nil, ErrEnoughInternalResourceCount
	}

	resourceId := resourceNodeIdArr[len(resourceNodeIdArr)%electionLocalSeed]
	jobNode, err := sche.resourceMng.GetDB().GetRegisterNode(pb.PrefixTypeJobNode, resourceId)
	if nil != err {
		return nil, err
	}
	if nil == jobNode {
		return nil, errors.New("not found jobNode information")
	}
	return jobNode, nil
}

func (sche *SchedulerStarveFIFO) electionConputeOrg(
	powerPartyIds []string,
	dataIdentityIdCache map[string]struct{},
	cost *twopctypes.TaskOperationCost,
) ([]*libTypes.TaskPowerSupplier, error) {

	calculateCount := len(powerPartyIds)
	identityIds := make([]string, 0)

	remoteResources := sche.resourceMng.GetRemoteResourceTables()
	log.Debugf("GetRemoteResouceTables on electionConputeOrg, remoteResources: %s", utilRemoteResourceArrString(remoteResources))
	for _, r := range remoteResources {

		// Skip the mock identityId
		if sche.resourceMng.IsMockIdentityId(r.GetIdentityId()) {
			continue
		}

		// 计算方不可以是任务发起方 和 数据参与方 和 接收方
		if _, ok := dataIdentityIdCache[r.GetIdentityId()]; ok {
			continue
		}
		// 还需要有足够的 资源
		if r.IsEnough(cost.Mem, cost.Bandwidth, cost.Processor) {
			identityIds = append(identityIds, r.GetIdentityId())
		}
	}

	if calculateCount > len(identityIds) {
		return nil, ErrEnoughResourceOrgCountLessCalculateCount
	}

	// Election
	index := electionOrgCondition % len(identityIds)
	identityIdTmp := make(map[string]struct{}, calculateCount)
	for i := calculateCount; i > 0; i-- {
		identityIdTmp[identityIds[index]] = struct{}{}
		index++

	}

	if len(identityIdTmp) != calculateCount {
		return nil, ErrEnoughResourceOrgCountLessCalculateCount
	}

	identityInfoArr, err := sche.resourceMng.GetDB().GetIdentityList()
	if nil != err {
		return nil, err
	}

	log.Debugf("GetIdentityList by dataCenter on electionConputeOrg, identityList: %s", identityInfoArr.String())
	identityInfoTmp := make(map[string]*types.Identity, calculateCount)
	for _, identityInfo := range identityInfoArr {

		// Skip the mock identityId
		if sche.resourceMng.IsMockIdentityId(identityInfo.IdentityId()) {
			continue
		}

		if _, ok := identityIdTmp[identityInfo.IdentityId()]; ok {
			identityInfoTmp[identityInfo.IdentityId()] = identityInfo
		}
	}

	if len(identityInfoTmp) != calculateCount {
		return nil, ErrEnoughResourceOrgCountLessCalculateCount
	}

	resourceArr, err := sche.resourceMng.GetDB().GetResourceList()
	if nil != err {
		return nil, err
	}

	log.Debugf("GetResourceList by dataCenter on electionConputeOrg, resources: %s", resourceArr.String())

	orgs := make([]*libTypes.TaskPowerSupplier, calculateCount)
	i := 0
	for _, iden := range resourceArr {

		if i == calculateCount {
			break
		}

		if info, ok := identityInfoTmp[iden.GetIdentityId()]; ok {
			orgs[i] = &libTypes.TaskPowerSupplier{
				Organization: &apipb.TaskOrganization{
					PartyId:    powerPartyIds[i],
					NodeName:   info.Name(),
					NodeId:     info.NodeId(),
					IdentityId: info.IdentityId(),
				},
				// TODO 这里的 task 资源消耗是事先加上的 先在这里直接加上 写死的(任务定义的)
				ResourceUsedOverview: &libTypes.ResourceUsageOverview{
					TotalMem:       iden.GetTotalMem(),
					UsedMem:        cost.Mem,
					TotalProcessor: uint32(iden.GetTotalProcessor()),
					UsedProcessor:  uint32(cost.Processor),
					TotalBandwidth: iden.GetTotalBandWidth(),
					UsedBandwidth:  cost.Bandwidth,
				},
			}
			i++
			delete(identityInfoTmp, iden.GetIdentityId())
		}
	}

	return orgs, nil
}

func utilOrgPowerArrString(powers []*libTypes.TaskPowerSupplier) string {
	arr := make([]string, len(powers))
	for i, power := range powers {
		arr[i] = power.String()
	}
	if len(arr) != 0 {
		return "[" + strings.Join(arr, ",") + "]"
	}
	return "[]"
}
func utilLocalResourceArrString(resources []*types.LocalResourceTable) string {
	arr := make([]string, len(resources))
	for i, r := range resources {
		arr[i] = r.String()
	}
	if len(arr) != 0 {
		return "[" + strings.Join(arr, ",") + "]"
	}
	return "[]"
}
func utilRemoteResourceArrString(resources []*types.RemoteResourceTable) string {
	arr := make([]string, len(resources))
	for i, r := range resources {
		arr[i] = r.String()
	}
	if len(arr) != 0 {
		return "[" + strings.Join(arr, ",") + "]"
	}
	return "[]"
}
func utilDataResourceArrString(resources []*types.DataResourceTable) string {
	arr := make([]string, len(resources))
	for i, r := range resources {
		arr[i] = r.String()
	}
	if len(arr) != 0 {
		return "[" + strings.Join(arr, ",") + "]"
	}
	return "[]"
}
