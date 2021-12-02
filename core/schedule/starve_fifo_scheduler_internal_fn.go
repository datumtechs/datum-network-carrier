package schedule

import (
	"container/heap"
	"errors"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"

	log "github.com/sirupsen/logrus"
	"strings"
)

func (sche *SchedulerStarveFIFO) pushTaskBullet(bullet *types.TaskBullet) error {
	sche.scheduleMutex.Lock()
	// The bullet is first into queue
	_, ok := sche.schedulings[bullet.TaskId]
	if !ok {
		heap.Push(sche.queue, bullet)
		sche.schedulings[bullet.TaskId] = bullet
		sche.resourceMng.GetDB().StoreTaskBullet(bullet)
	}
	sche.scheduleMutex.Unlock()
	log.Debugf("Succeed pushed local task into queue on scheduler, taskId: {%s}", bullet.TaskId)
	return nil
}

func (sche *SchedulerStarveFIFO) repushTaskBullet(bullet *types.TaskBullet) error {
	sche.scheduleMutex.Lock()

	if bullet.Starve {
		heap.Push(sche.starveQueue, bullet)

		log.Debugf("Succeed repushed task into starve queue on scheduler, taskId: {%s}, reschedCount: {%d}, max threshold: {%d}",
			bullet.TaskId, bullet.Resched, ReschedMaxCount)
	} else {
		heap.Push(sche.queue, bullet)

		log.Debugf("Succeed repushed task into queue on scheduler, taskId: {%s}, reschedCount: {%d}, max threshold: {%d}",
			bullet.TaskId, bullet.Resched, ReschedMaxCount)
	}
	sche.resourceMng.GetDB().StoreTaskBullet(bullet)  // cover old value with new value into db
	sche.scheduleMutex.Unlock()
	return nil
}

func (sche *SchedulerStarveFIFO) removeTaskBullet(taskId string) error {
	sche.scheduleMutex.Lock()
	defer sche.scheduleMutex.Unlock()

	_ ,ok := sche.schedulings[taskId]
	if !ok {
		return nil
	}

	log.Debugf("Succeed removed local bullet task on scheduler, taskId: {%s}", taskId)

	delete(sche.schedulings, taskId)
	sche.resourceMng.GetDB().RemoveTaskBullet(taskId)

	// traversal the queue to remove task bullet, first.
	i := 0
	for {
		if i == sche.queue.Len() {
			break
		}
		qbullet := (*(sche.queue))[i]
		// When found the bullet with taskId, removed it from queue.
		if qbullet.GetTaskId() == taskId {
			heap.Remove(sche.queue, i)
			return nil
		}
		(*(sche.queue))[i] = qbullet
		i++
	}

	// otherwise, traversal the starveQueue to remove task bullet, second.
	i = 0
	for {
		if i == sche.starveQueue.Len() {
			break
		}
		qbullet := (*(sche.starveQueue))[i]

		// When found the bullet with taskId, removed it from starveQueue.
		if qbullet.GetTaskId() == taskId {
			heap.Remove(sche.starveQueue, i)
			return nil
		}
		(*(sche.starveQueue))[i] = qbullet
		i++
	}
	return nil
}

func (sche *SchedulerStarveFIFO) popTaskBullet() *types.TaskBullet {
	sche.scheduleMutex.Lock()

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
	sche.scheduleMutex.Unlock()
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

func (sche *SchedulerStarveFIFO) electionJobNode(mem, bandwidth, disk uint64, processor uint32) (*pb.YarnRegisteredPeerDetail, error) {

	if nil == sche.internalNodeSet || 0 == sche.internalNodeSet.JobNodeClientSize() {
		return nil, errors.New("not found alive jobNode")
	}

	resourceNodeIdArr := make([]string, 0)

	tables, err := sche.resourceMng.QueryLocalResourceTables()
	if nil != err {
		return nil, err
	}
	log.Debugf("QueryLocalResourceTables on electionJobNode, localResources: %s", utilLocalResourceArrString(tables))
	for _, r := range tables {
		isEnough := r.IsEnough(mem, bandwidth, disk, processor)
		log.Debugf("Call electionJobNode, resource: %s, r.RemainMem(): %d, r.RemainBandwidth(): %d, r.RemainDisk(): %d, r.RemainProcessor(): %d, needMem: %d, needBandwidth: %d, needDisk: %d, needProcessor: %d, isEnough: %v",
			r.String(), r.RemainMem(), r.RemainBandwidth(), r.RemainDisk(), r.RemainProcessor(), mem, bandwidth, disk, processor, isEnough)
		if isEnough {
			jobNodeClient, find := sche.internalNodeSet.QueryJobNodeClient(r.GetNodeId())
			if find && jobNodeClient.IsConnected() {
				resourceNodeIdArr = append(resourceNodeIdArr, r.GetNodeId())
				log.Debugf("Call electionJobNode, append resourceId: %s", r.GetNodeId())
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
		return nil, fmt.Errorf("not found jobNode information")
	}
	return jobNode, nil
}

func (sche *SchedulerStarveFIFO) electionPowerOrg(
	powerPartyIds []string,
	skipIdentityIdCache map[string]struct{},
	mem, bandwidth, disk uint64, processor uint32,
) ([]*libtypes.TaskPowerSupplier, error) {

	calculateCount := len(powerPartyIds)

	// Find global identitys
	identityInfoArr, err := sche.resourceMng.GetDB().QueryIdentityList(timeutils.BeforeYearUnixMsecUint64(), backend.DefaultMaxPageSize)
	if nil != err {
		return nil, err
	}

	if len(identityInfoArr) < calculateCount {
		return nil, fmt.Errorf("query identityList count less calculate count")
	}

	log.Debugf("QueryIdentityList by dataCenter on electionPowerOrg, len: {%d}, identityList: %s", len(identityInfoArr), identityInfoArr.String())
	identityInfoTmp := make(map[string]*types.Identity, calculateCount)
	for _, identityInfo := range identityInfoArr {

		// Skip the mock identityId
		if sche.resourceMng.IsMockIdentityId(identityInfo.GetIdentityId()) {
			continue
		}

		identityInfoTmp[identityInfo.GetIdentityId()] = identityInfo
	}

	if len(identityInfoTmp) < calculateCount {
		return nil, fmt.Errorf("find valid identityIds count less calculate count")
	}

	// Find global power resources
	globalResources, err := sche.resourceMng.GetDB().QueryGlobalResourceSummaryList(timeutils.BeforeYearUnixMsecUint64(), backend.DefaultMaxPageSize)
	if nil != err {
		return nil, err
	}
	//log.Debugf("GetRemoteResouceTables on electionComputeOrg, globalResources: %s", utilRemoteResourceArrString(globalResources))
	log.Debugf("GetRemoteResouceTables on electionPowerOrg, len: {%d}, globalResources: %s", len(globalResources), globalResources.String())

	if len(globalResources) < calculateCount {
		return nil, fmt.Errorf("query org's power resource count less calculate count")
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
		if rMem < mem {
			continue
		}
		if rProcessor < processor {
			continue
		}
		if rBandwidth < bandwidth {
			continue
		}
		// ignore disk for power resource.

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
					TotalMem:       r.GetTotalMem(),   			// total resource value of org.
					UsedMem: 0,     							// used resource of this task (real time max used)
					TotalBandwidth: r.GetTotalBandWidth(),
					UsedBandwidth: 0, 							// used resource of this task (real time max used)
					TotalDisk:      r.GetTotalDisk(),
					UsedDisk:       0,
					TotalProcessor: r.GetTotalProcessor(),
					UsedProcessor: 0,							// used resource of this task (real time max used)
				},
			})
			i++
		}
	}
	if len(orgs) < calculateCount {
		return nil, ErrEnoughResourceOrgCountLessCalculateCount
	}
	return orgs, nil
}

func (sche *SchedulerStarveFIFO) verifyUserMetadataAuthOnTask(userType apicommonpb.UserType, user, metadataId string) error {
	return sche.authMng.VerifyMetadataAuth(userType, user, metadataId)
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
