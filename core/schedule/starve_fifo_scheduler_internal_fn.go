package schedule

import (
	"container/heap"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"

	log "github.com/sirupsen/logrus"
)

func (sche *SchedulerStarveFIFO) pushTaskBullet(bullet *types.TaskBullet) error {
	sche.scheduleMutex.Lock()
	defer sche.scheduleMutex.Unlock()
	// The bullet is first into queue
	_, ok := sche.schedulings[bullet.GetTaskId()]
	if !ok {

		bullet.InQueueFlag = true
		heap.Push(sche.queue, bullet)
		sche.schedulings[bullet.GetTaskId()] = bullet
		sche.resourceMng.GetDB().StoreTaskBullet(bullet)
	}
	log.Debugf("Succeed pushed local task into queue on scheduler, taskId: {%s}", bullet.GetTaskId())
	return nil
}

func (sche *SchedulerStarveFIFO) repushTaskBullet(bullet *types.TaskBullet) error {

	if bullet.GetInQueueFlag() {
		log.Warnf("Warning repush local bullet task on scheduler, task was not exist in schedulings map, taskId: {%s}, inQueueFlag: {%v}",
			bullet.GetTaskId(), bullet.GetInQueueFlag())
		return nil
	}

	bullet.InQueueFlag = true
	if bullet.IsStarve() {
		heap.Push(sche.starveQueue, bullet)

		log.Debugf("Succeed repushed task into starve queue on scheduler, taskId: {%s}, reschedCount: {%d}, max threshold: {%d}",
			bullet.GetTaskId(), bullet.GetResched(), ReschedMaxCount)
	} else {
		heap.Push(sche.queue, bullet)

		log.Debugf("Succeed repushed task into queue on scheduler, taskId: {%s}, reschedCount: {%d}, max threshold: {%d}",
			bullet.GetTaskId(), bullet.GetResched(), ReschedMaxCount)
	}
	sche.resourceMng.GetDB().StoreTaskBullet(bullet)  // update bullet into wal

	return nil
}

func (sche *SchedulerStarveFIFO) removeTaskBullet(taskId string) error {
	sche.scheduleMutex.Lock()
	defer sche.scheduleMutex.Unlock()

	_ ,ok := sche.schedulings[taskId]
	if !ok {
		log.Warnf("Warning removed local bullet task on scheduler, task was not exist in schedulings map, taskId: {%s}", taskId)
		return nil
	}

	delete(sche.schedulings, taskId)
	sche.resourceMng.GetDB().RemoveTaskBullet(taskId)

	log.Debugf("Succeed removed local bullet task on scheduler, taskId: {%s}", taskId)

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
			return nil
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
			return nil
		}
		(*(sche.starveQueue))[i] = bullet
		i++
	}
	return nil
}

func (sche *SchedulerStarveFIFO) popTaskBullet() *types.TaskBullet {

	var bullet *types.TaskBullet
	if sche.starveQueue.Len() != 0 {
		x := heap.Pop(sche.starveQueue)
		bullet = x.(*types.TaskBullet)
		bullet.InQueueFlag = false
	} else {
		if sche.queue.Len() != 0 {
			x := heap.Pop(sche.queue)
			bullet = x.(*types.TaskBullet)
			bullet.InQueueFlag = false
		}
	}
	if nil != bullet {
		sche.resourceMng.GetDB().StoreTaskBullet(bullet)  // update bullet into wal
	}

	return bullet
}

func (sche *SchedulerStarveFIFO) increaseTotalTaskTerm() {
	// handle starve queue
	sche.starveQueue.IncreaseTermByCallbackFn(func(bullet *types.TaskBullet) {
		sche.resourceMng.GetDB().StoreTaskBullet(bullet)  		// update bullet into wal
	})

	// handle queue
	i := 0
	for {
		if i == sche.queue.Len() {
			return
		}
		bullet := (*(sche.queue))[i]
		bullet.IncreaseTerm()

		// When the task in the queue meets hunger, it will be transferred to starveQueue
		if bullet.GetTerm() >= StarveTerm {
			bullet.Starve = true
			heap.Push(sche.starveQueue, bullet)   				// push it into `starve queue`
			heap.Remove(sche.queue, i)							// remove it from `queue`
			sche.resourceMng.GetDB().StoreTaskBullet(bullet)  	// update bullet into wal
			i = 0
			continue
		}
		(*(sche.queue))[i] = bullet
		sche.resourceMng.GetDB().StoreTaskBullet(bullet)  		// update bullet into wal
		i++
	}
}

func (sche *SchedulerStarveFIFO) verifyUserMetadataAuthOnTask(userType libtypes.UserType, user, metadataId string) error {
	return sche.authMng.VerifyMetadataAuth(userType, user, metadataId)
}
