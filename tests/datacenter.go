package tests

import (
	"context"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	"github.com/RosettaFlow/Carrier-Go/core"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
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
	identities := types.NewIdentity(&libtypes.IdentityPB{
		IdentityId: Identity,
		NodeId:     NodeId,
		NodeName:   NodeName,
		DataId:     DataId,
		DataStatus: apicommonpb.DataStatus_DataStatus_Normal,
		Status:     apicommonpb.CommonStatus_CommonStatus_Normal,
		Credential: "",
	})
	merr := dc.InsertIdentity(identities)

	if merr != nil {
		panic("from data center InsertIdentity fail," + merr.Error() + "\n")
	} else {
		fmt.Println("InsertIdentity successful")
		result, err := dc.QueryIdentityList()
		if nil != err {
			fmt.Println("err", err)
		}
		fmt.Println(result)
	}
}

func InsertMetadata() {
	dc := serverObj()
	metadata := types.NewMetadata(&libtypes.MetadataPB{
		IdentityId: Identity,
		NodeId:     NodeId,
		DataId:     DataId,
		DataStatus: apicommonpb.DataStatus_DataStatus_Normal,
		OriginId:   OriginId,
		TableName:  "test_table1",
		FilePath:   "/c/c",
		Desc:       "desc",
		Rows:       1,
		Columns:    2,
		Size_:      3,
		FileType:   apicommonpb.OriginFileType_FileType_CSV,
		State:      apicommonpb.MetadataState_MetadataState_Released,
		HasTitle:   false,
		MetadataColumns: []*libtypes.MetadataColumn{
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
		result, err := dc.QueryMetadataList()
		if nil != err {
			fmt.Println("err", err)
		}
		fmt.Println(result)
	}
}

func InsertResource() {
	dc := serverObj()
	resource := types.NewResource(&libtypes.ResourcePB{
		IdentityId:     Identity,
		NodeId:         NodeId,
		NodeName:       NodeName,
		DataId:         DataId,
		DataStatus:     apicommonpb.DataStatus_DataStatus_Normal,
		State:          apicommonpb.PowerState_PowerState_Released,
		TotalMem:       1,
		UsedMem:        2,
		TotalProcessor: 0,
		TotalBandwidth: 1,
	})

	terr := dc.InsertResource(resource)
	if terr != nil {
		panic("from data center InsertResource fail," + terr.Error() + "\n")
	} else {
		fmt.Print("InsertResource successful\n")
		//result, err := dc.QueryResourceList()
		//if nil != err {
		//	fmt.Println("err", err)
		//}
		//fmt.Println(result)
	}

}

func InsertTask() {
	dc := serverObj()
	taskdata := types.NewTask(&libtypes.TaskPB{
		IdentityId: Identity,
		NodeId:     NodeId,
		NodeName:   NodeName,
		DataId:     DataId,
		DataStatus: apicommonpb.DataStatus_DataStatus_Normal,
		TaskId:     TaskId,
		TaskName:   "task001",
		State:      apicommonpb.TaskState_TaskState_Succeed,
		Reason:     "",
		EventCount: 0,
		Desc:       "",
		CreateAt:   timeutils.UnixMsecUint64(),
		EndAt:      timeutils.UnixMsecUint64(),
		AlgoSupplier: &apicommonpb.TaskOrganization{
			PartyId:    "",
			NodeId:     Identity,
			NodeName:   NodeName,
			IdentityId: Identity,
		},
		OperationCost: &apicommonpb.TaskResourceCostDeclare{
			Memory:       12,
			Bandwidth: 120,
			Processor: 8,
		},
		DataSuppliers: []*libtypes.TaskDataSupplier{
			{
				Organization: &apicommonpb.TaskOrganization{
					PartyId:    "",
					NodeId:     Identity,
					NodeName:   NodeName,
					IdentityId: Identity,
				},
				MetadataId:   DataId,
				MetadataName: "meta1",
				KeyColumn: &libtypes.MetadataColumn{

						CIndex:   2,
						CName:    "cname",
						CType:    "ctype",
						CSize:    10,
						CComment: "this test",
					},

			},
		},
		PowerSuppliers: []*libtypes.TaskPowerSupplier{
			{
				Organization: &apicommonpb.TaskOrganization{
					PartyId:    "",
					NodeId:     Identity,
					NodeName:   NodeName,
					IdentityId: Identity,
				},
				ResourceUsedOverview: &libtypes.ResourceUsageOverview{
					TotalMem:       12,
					UsedMem:        8,
					TotalProcessor: 8,
					UsedProcessor:  4,
					TotalBandwidth: 120,
					UsedBandwidth:  25,
				},
			},
		},
		TaskEvents: []*libtypes.TaskEvent{
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
		result, err := dc.QueryLocalTaskList()
		if nil != err {
			fmt.Println("err", err)
		}
		fmt.Println(result)
	}
}

func RevokeIdentity() {
	dc := serverObj()
	identities := types.NewIdentity(&libtypes.IdentityPB{
		IdentityId: Identity,
		NodeId:     NodeId,
		NodeName:   NodeName,
		DataId:     DataId,
		DataStatus: apicommonpb.DataStatus_DataStatus_Normal,
		Status:     apicommonpb.CommonStatus_CommonStatus_Normal,
		Credential: "",
	})
	result := dc.RevokeIdentity(identities)
	if result == nil {
		fmt.Printf("RevokeIdentity successful")
	}
}
func GetData() {
	dc := serverObj()

	//region QueryIdentityList
	identity, err := dc.QueryIdentityList()
	if nil != err {
		fmt.Println("err", err)
	}
	fmt.Println("QueryIdentityList result is:")
	for _, value := range identity {
		fmt.Println(types.IdentityDataTojson(value))
	}
	//endregion

	// region QueryMetadataList
	metadata, err := dc.QueryMetadataList()
	if nil != err {
		fmt.Println("err", err)
	}
	fmt.Println("QueryMetadataList result is:")
	fmt.Println(metadata)
	//endregion

	//region QueryMetadataByDataId
	MetadataByDataId, err := dc.QueryMetadataByDataId(DataId)
	if nil != err {
		fmt.Println("err", err)
	}
	fmt.Println("QueryMetadataByDataId result is:")
	fmt.Println(types.MetadataToJson(MetadataByDataId))
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
	GetTaskList, err := dc.QueryLocalTaskList()
	if nil != err {
		fmt.Println("err", err)
	}
	fmt.Println("GetTaskList result is:")
	fmt.Println(GetTaskList)
	// endregion

	// region HasIdentity
	result, err := dc.HasIdentity(&apicommonpb.Organization{
		NodeName:   NodeName,
		NodeId:     NodeId,
		IdentityId: Identity})
	if nil != err {
		fmt.Println("err", err)
	}
	fmt.Println("HasIdentity result is:", result)
	fmt.Println(result)
	// endregion
}
