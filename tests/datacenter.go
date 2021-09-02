package tests

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	"github.com/RosettaFlow/Carrier-Go/core"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/params"
	"github.com/RosettaFlow/Carrier-Go/types"
)

var (
	NodeName   = "orgName_000103"
	NodeId     = "NodeId_000103"
	Identity   = "identityId_000103"
	DataId     = "DataId_000103"
	OriginId   = "OriginId_000103"
	TaskId     = "123456"
	DataStatus = "Y"
	Status     = "Y"
)

func serverObj() *core.DataCenter {

	server := params.DataCenterConfig{"192.168.112.32", 9099}
	ctx := context.Background()
	dc, err := core.NewDataCenter(ctx, nil, &server)
	if err != nil {
		panic("init data center fail," + err.Error())
	} else {
		fmt.Printf("GrupUrl:" + server.GrpcUrl + "\n")
	}
	return dc
}

func InsertData() {
	dc := serverObj()
	identities := types.NewIdentity(&libTypes.IdentityPB{
		IdentityId: Identity,
		NodeId:     NodeId,
		NodeName:   NodeName,
		DataId:     DataId,
		DataStatus: DataStatus,
		Status:     Status,
		Credential: "",
	})
	merr := dc.InsertIdentity(identities)

	if merr != nil {
		panic("from data center InsertIdentity fail," + merr.Error() + "\n")
	} else {
		fmt.Println("InsertIdentity successful")
		result, err := dc.GetIdentityList()
		if nil != err {
			fmt.Println("err", err)
		}
		fmt.Println(result)
	}
}

func InsertMetaData() {
	dc := serverObj()
	metadata := types.NewMetadata(&libTypes.MetadataPB{
		IdentityId: Identity,
		NodeId:     NodeId,
		DataId:     DataId,
		DataStatus: DataStatus,
		OriginId:   OriginId,
		TableName:  "test_table1",
		FilePath:   "/c/c",
		Desc:       "desc",
		Rows:       1,
		Columns:    2,
		Size_:      3,
		FileType:   "csv",
		State:      Status,
		HasTitle:   false,
		MetadataColumnList: []*libTypes.MetadataColumn{
			{
				CIndex: 2,
				CName:  "cname",
				CType:  "ctype",
				CSize:  10,
			},
		},
	})
	merr := dc.InsertMetadata(metadata)
	if merr != nil {
		panic("from data center InsertMetadata fail," + merr.Error() + "\n")
	} else {
		fmt.Println("InsertMetadata successful")
		result, err := dc.GetMetadataList()
		if nil != err {
			fmt.Println("err", err)
		}
		fmt.Println(result)
	}
}

func InsertResource() {
	dc := serverObj()
	resource := types.NewResource(&libTypes.ResourcePB{
		IdentityId:     Identity,
		NodeId:         NodeId,
		NodeName:       NodeName,
		DataId:         DataId,
		DataStatus:     DataStatus,
		State:          Status,
		TotalMem:       1,
		UsedMem:        2,
		TotalProcessor: 0,
		TotalBandWidth: 1,
	})

	terr := dc.InsertResource(resource)
	if terr != nil {
		panic("from data center InsertResource fail," + terr.Error() + "\n")
	} else {
		fmt.Print("InsertResource successful\n")
		//result, err := dc.GetResourceList()
		//if nil != err {
		//	fmt.Println("err", err)
		//}
		//fmt.Println(result)
	}

}

func InsertTask() {
	dc := serverObj()
	taskdata := types.NewTask(&libTypes.TaskPB{
		IdentityId: Identity,
		NodeId:     NodeId,
		NodeName:   NodeName,
		DataId:     DataId,
		DataStatus: DataStatus,
		TaskId:     TaskId,
		TaskName:   "task001",
		State:      Status,
		Reason:     "",
		EventCount: 0,
		Desc:       "",
		CreateAt:   uint64(timeutils.UnixMsec()),
		EndAt:      uint64(timeutils.UnixMsec()),
		AlgoSupplier: &apipb.TaskOrganization{
			PartyId:  "",
			NodeId:   Identity,
			NodeName: NodeName,
			IdentityId: Identity,
		},
		OperationCost: &apipb.TaskResourceCostDeclare{
			CostMem:       12,
			CostBandwidth: 120,
			CostProcessor: 8,
		},
		DataSupplier: []*libTypes.TaskDataSupplier{
			{
				MemberInfo: &apipb.TaskOrganization{
					PartyId:    "",
					NodeId:     Identity,
					NodeName:   NodeName,
					IdentityId: Identity,
				},
				MetadataId:   DataId,
				MetadataName: "meta1",
				ColumnList: []*libTypes.MetadataColumn{
					{
						CIndex:   2,
						CName:    "cname",
						CType:    "ctype",
						CSize:    10,
						CComment: "this test",
					},
				},
			},
		},
		PowerSupplier: []*libTypes.TaskPowerSupplier{
			{
				Organization: &apipb.TaskOrganization{
					PartyId:    "",
					NodeId:     Identity,
					NodeName:   NodeName,
					IdentityId: Identity,
				},
				ResourceUsedOverview: &libTypes.ResourceUsageOverview{
					TotalMem:       12,
					UsedMem:        8,
					TotalProcessor: 8,
					UsedProcessor:  4,
					TotalBandwidth: 120,
					UsedBandwidth:  25,
				},
			},
		},
		Receivers: []*libTypes.TaskResultReceiver{
			{
				Receiver: &apipb.TaskOrganization{
					PartyId:    "",
					NodeId:     Identity,
					NodeName:   NodeName,
					IdentityId: Identity,
				},
				Providers: []*apipb.TaskOrganization{
					{
						PartyId:    "",
						NodeId:     Identity,
						NodeName:   NodeName,
						IdentityId: Identity,
					},
				},
			},
		},
		PartnerList: []*apipb.TaskOrganization{
			{
				PartyId:    "",
				NodeId:     Identity,
				NodeName:   NodeName,
				IdentityId: Identity,
			},
		},
		TaskEventList: []*libTypes.TaskEvent{
			{
				TaskId:     "123456",
				Type:       "1-evengine-eventType",
				CreateAt:   0,
				Content:    "1-evengine-eventContent",
				IdentityId: Identity,
			},
		},
	})
	terr := dc.InsertTask(taskdata)
	if terr != nil {
		panic("from data center InsertTask fail," + terr.Error() + "\n")
	} else {
		fmt.Print("InsertTask successful\n")
		result, err := dc.GetLocalTaskList()
		if nil != err {
			fmt.Println("err", err)
		}
		fmt.Println(result)
	}
}

func RevokeIdentity() {
	dc := serverObj()
	identities := types.NewIdentity(&libTypes.IdentityPB{
		IdentityId: Identity,
		NodeId:     NodeId,
		NodeName:   NodeName,
		DataId:     DataId,
		DataStatus: DataStatus,
		Status:     Status,
		Credential: "",
	})
	result := dc.RevokeIdentity(identities)
	if result == nil {
		fmt.Printf("RevokeIdentity successful")
	}
}
func GetData() {
	dc := serverObj()

	//region GetIdentityList
	identity, err := dc.GetIdentityList()
	if nil != err {
		fmt.Println("err", err)
	}
	fmt.Println("GetIdentityList result is:")
	for _, value := range identity {
		fmt.Println(types.IdentityDataTojson(value))
	}
	//endregion

	// region GetMetadataList
	metadata, err := dc.GetMetadataList()
	if nil != err {
		fmt.Println("err", err)
	}
	fmt.Println("GetMetadataList result is:")
	fmt.Println(metadata)
	//endregion

	//region GetMetadataByDataId
	MetaDataByDataId, err := dc.GetMetadataByDataId(DataId)
	if nil != err {
		fmt.Println("err", err)
	}
	fmt.Println("GetMetadataByDataId result is:")
	fmt.Println(types.MetaDataToJson(MetaDataByDataId))
	// endregion

	// region GetResourceListByNodeId
	ResourceByNodeId, err := dc.GetResourceListByIdentityId(Identity)
	if nil != err {
		fmt.Println("err", err)
	}
	fmt.Printf("GetResourceListByNodeId result is:")
	fmt.Printf("%+v\n", ResourceByNodeId)
	// endregion

	// region GetTaskList
	GetTaskList, err := dc.GetLocalTaskList()
	if nil != err {
		fmt.Println("err", err)
	}
	fmt.Println("GetTaskList result is:")
	fmt.Println(GetTaskList)
	// endregion

	// region HasIdentity
	result, err := dc.HasIdentity(&apipb.Organization{
		Name:       NodeName,
		NodeId:     NodeId,
		IdentityId: Identity})
	if nil != err {
		fmt.Println("err", err)
	}
	fmt.Println("HasIdentity result is:", result)
	fmt.Println(result)
	// endregion
}
