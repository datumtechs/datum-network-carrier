package schedule

import (
	"container/heap"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/types"

	log "github.com/sirupsen/logrus"
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

func (sche *SchedulerStarveFIFO) verifyUserMetadataAuthOnTask(userType apicommonpb.UserType, user, metadataId string) error {
	return sche.authMng.VerifyMetadataAuth(userType, user, metadataId)
}
