package schedule

import (
	"container/heap"
	"errors"
	twopctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"

	log "github.com/sirupsen/logrus"
	"strings"
)

func (sche *SchedulerStarveFIFO) pushTaskBullet(bullet *types.TaskBullet) error {
	sche.schedulingsMutex.Lock()
	// The bullet is first into queue
	_, ok := sche.schedulings[bullet.TaskId]
	if !ok {
		heap.Push(sche.queue, bullet)
		sche.schedulings[bullet.TaskId] = bullet
		sche.resourceMng.GetDB().StoreScheduling(bullet)
	}
	sche.schedulingsMutex.Unlock()
	return nil
}

func (sche *SchedulerStarveFIFO) repushTaskBullet(bullet *types.TaskBullet) error {
	sche.schedulingsMutex.Lock()

	if bullet.Starve {
		log.Debugf("repush task into starve queue, taskId: {%s}, reschedCount: {%d}, max threshold: {%d}",
			bullet.TaskId, bullet.Resched, ReschedMaxCount)
		heap.Push(sche.starveQueue, bullet)
	} else {
		log.Debugf("repush task into queue, taskId: {%s}, reschedCount: {%d}, max threshold: {%d}",
			bullet.TaskId, bullet.Resched, ReschedMaxCount)
		heap.Push(sche.queue, bullet)
	}
	sche.schedulingsMutex.Unlock()
	return nil
}

func (sche *SchedulerStarveFIFO) removeTaskBullet(taskId string) error {
	sche.schedulingsMutex.Lock()
	defer sche.schedulingsMutex.Unlock()

	// traversal the queue to remove task bullet, first.
	i := 0
	for {
		if i == sche.queue.Len() {
			break
		}
		bullet := (*(sche.queue))[i]

		// When found the bullet with taskId, removed it from queue.
		if bullet.GetTaskId() == taskId {
			heap.Remove(sche.queue, i)
			sche.resourceMng.GetDB().DeleteScheduling(sche.schedulings[taskId])
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
		if bullet.GetTaskId() == taskId {
			heap.Remove(sche.starveQueue, i)
			sche.resourceMng.GetDB().DeleteScheduling(sche.schedulings[taskId])
			delete(sche.schedulings, taskId)
			return nil // todo 这里需要做一次 持久化
		}
		(*(sche.starveQueue))[i] = bullet
		i++
	}
	return nil
}

func (sche *SchedulerStarveFIFO) popTaskBullet() *types.TaskBullet {
	sche.schedulingsMutex.Lock()

	var bullet *types.TaskBullet

	if sche.starveQueue.Len() != 0 {
		x := heap.Pop(sche.starveQueue)
		bullet = x.(*types.TaskBullet)
	} else {
		if sche.queue.Len() != 0 {
			x := heap.Pop(sche.queue)
			bullet = x.(*types.TaskBullet)
		}
	}
	sche.schedulingsMutex.Unlock()
	return bullet
}

func (sche *SchedulerStarveFIFO) increaseTotalTaskTerm() {
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

	resourceId := resourceNodeIdArr[len(resourceNodeIdArr)-1]
	jobNode, err := sche.resourceMng.GetDB().QueryRegisterNode(pb.PrefixTypeJobNode, resourceId)
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
	skipIdentityIdCache map[string]struct{},
	cost *twopctypes.TaskOperationCost,
) ([]*libtypes.TaskPowerSupplier, error) {

	calculateCount := len(powerPartyIds)

	// Find global identitys
	identityInfoArr, err := sche.resourceMng.GetDB().QueryIdentityList()
	if nil != err {
		return nil, err
	}

	if len(identityInfoArr) != calculateCount {
		return nil, ErrEnoughResourceOrgCountLessCalculateCount
	}

	log.Debugf("QueryIdentityList by dataCenter on electionConputeOrg, len: {%d}, identityList: %s", len(identityInfoArr), identityInfoArr.String())
	identityInfoTmp := make(map[string]*types.Identity, calculateCount)
	for _, identityInfo := range identityInfoArr {

		// Skip the mock identityId
		if sche.resourceMng.IsMockIdentityId(identityInfo.GetIdentityId()) {
			continue
		}

		identityInfoTmp[identityInfo.GetIdentityId()] = identityInfo
	}

	if len(identityInfoTmp) != calculateCount {
		return nil, ErrEnoughResourceOrgCountLessCalculateCount
	}

	// Find global power resources
	globalResources, err := sche.resourceMng.GetDB().QueryResourceList()
	if nil != err {
		return nil, err
	}
	//log.Debugf("GetRemoteResouceTables on electionConputeOrg, globalResources: %s", utilRemoteResourceArrString(globalResources))
	log.Debugf("GetRemoteResouceTables on electionConputeOrg, len: {%d}, globalResources: %s", len(globalResources), globalResources.String())

	if len(globalResources) != calculateCount {
		return nil, ErrEnoughResourceOrgCountLessCalculateCount
	}

	orgs := make([]*libtypes.TaskPowerSupplier, 0)
	i := 0
	for _, r := range globalResources {

		if i == calculateCount {
			break
		}

		// skip
		if len(skipIdentityIdCache) != 0 {
			if _, ok := skipIdentityIdCache[r.GetIdentityId()]; ok {
				continue
			}
		}

		// Find one, if have enough resource
		rMem, rBandwidth, rProcessor := r.GetTotalMem()-r.GetUsedMem(), r.GetTotalBandWidth()-r.GetUsedBandWidth(), r.GetTotalProcessor()-r.GetUsedProcessor()
		if rMem < cost.Mem {
			continue
		}
		if rProcessor < cost.Processor {
			continue
		}
		if rBandwidth < cost.Bandwidth {
			continue
		}

		// append one, if it enouph
		if info, ok := identityInfoTmp[r.GetIdentityId()]; ok {
			orgs = append(orgs, &libtypes.TaskPowerSupplier{
				Organization: &apicommonpb.TaskOrganization{
					PartyId:    powerPartyIds[i],
					NodeName:   info.GetName(),
					NodeId:     info.GetNodeId(),
					IdentityId: info.GetIdentityId(),
				},
				ResourceUsedOverview: &libtypes.ResourceUsageOverview{
					TotalMem:       r.GetTotalMem(),
					UsedMem:        cost.Mem,
					TotalProcessor: r.GetTotalProcessor(),
					UsedProcessor:  cost.Processor,
					TotalBandwidth: r.GetTotalBandWidth(),
					UsedBandwidth:  cost.Bandwidth,
				},
			})
			i++
		}
	}
	return orgs, nil
}

func (sche *SchedulerStarveFIFO) verifyUserMetadataAuthOnTask(userType apicommonpb.UserType, user, metadataId string) bool {
	if !sche.authMng.VerifyMetadataAuth(userType, user, metadataId) {
		return false
	}
	return true
}

func utilOrgPowerArrString(powers []*libtypes.TaskPowerSupplier) string {
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
