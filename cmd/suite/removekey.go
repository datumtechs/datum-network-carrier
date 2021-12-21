package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/core"
	"github.com/RosettaFlow/Carrier-Go/db"
)

func main() {

	var taskId string
	var dbpath string
	var datacenterHost string
	var datacenterPort uint64
	flag.StringVar(&taskId, "taskId", "", "remove taskId")
	flag.StringVar(&dbpath, "dbpath", "", "level db local path")
	flag.StringVar(&datacenterHost, "datacenterHost", "", "datacenter host")
	flag.Uint64Var(&datacenterPort, "datacenterPort", 0, "datacenter port")

	flag.Parse()

	fmt.Printf("Start remove task, taskId: {%s},\ndbpath: {%s},\ndatacenterHost: {%s},\ndatacenterPort: {%d}\n",
		taskId, dbpath, datacenterHost, datacenterPort)

	db, err := db.NewLDBDatabase(dbpath, 768, 8192)
	if err != nil {
		fmt.Printf("Failed to open local db, %s\n", err)
		return
	}

	carrierDB := core.NewDataCenter(context.TODO(), db)

	// Remove the only task that everyone refers to together
	if err := carrierDB.RemoveLocalTask(taskId); nil != err {
		fmt.Printf("Failed to remove local task, taskId: {%s}, %s\n,", taskId, err)
	}

	// Remove the only things in task that everyone refers to together
	if err := carrierDB.RemoveTaskPowerPartyIds(taskId); nil != err {
		fmt.Printf("Failed to remove power's partyIds of local task, taskId: {%s}, %s\n,", taskId, err)
	}

	// Remove the partyId list of current task participants saved by the task sender
	if err := carrierDB.RemoveTaskPartnerPartyIds(taskId); nil != err {
		fmt.Printf("Failed to remove handler partner's partyIds of local task, taskId: {%s}, %s\n,", taskId, err)
	}

	// Remove the task event of all partys
	if err := carrierDB.RemoveTaskEventList(taskId); nil != err {
		fmt.Printf("Failed to clean all event list of task, taskId: {%s}, %s\n,", taskId, err)
	}

}



