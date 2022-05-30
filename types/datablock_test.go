package types
//
//import (
//	"github.com/datumtechs/datum-network-carrier/common"
//	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
//	"reflect"
//	"testing"
//)
//
//func newBlock() *Block {
//	metadatas := func() MetadataArray {
//		return MetadataArray{
//			&Metadata{
//				data: &carriertypespb.MetadataPB{
//					IdentityId: "identity",
//					NodeId:     "NodeId",
//					DataId:     "dataId",
//					DataStatus: carriertypespb.DataStatus_DataStatus_Deleted,
//					DataPath:   "/a/b/c",
//					Desc:       "desc",
//					Rows:       10,
//					Columns:    100,
//					FileType:   carriertypespb.OriginFileType_FileType_CSV,
//					State:      carriertypespb.MetadataState_MetadataState_Unknown,
//					HasTitle:   true,
//					MetadataColumns: []*carriertypespb.MetadataColumn{
//						{
//							CIndex: 0,
//							CName:  "cname",
//							CType:  "ctype",
//							CSize:  0,
//						},
//					},
//				},
//			},
//		}
//	}
//	resources := func() ResourceArray {
//		return ResourceArray{
//			{
//				data: &carriertypespb.ResourcePB{
//					IdentityId:     "resource-identity",
//					NodeId:         "resource-nodeId",
//					NodeName:       "resource-nodeName",
//					DataId:         "resource-dataId",
//					DataStatus:     carriertypespb.DataStatus_DataStatus_Deleted,
//					State:          carriertypespb.PowerState_PowerState_Unknown,
//					TotalMem:       100 * 1024,
//					UsedMem:        100 * 2014,
//					TotalProcessor: 0,
//					UsedProcessor: 0,
//					TotalBandwidth: 100 * 1024,
//					UsedBandwidth: 0,
//				},
//			},
//		}
//	}
//	identities := func() IdentityArray {
//		return IdentityArray{
//			{
//				data: &carriertypespb.IdentityPB{
//					IdentityId: "identity-identity",
//					NodeId:     "identity-nodeId",
//					NodeName:   "identity-nodeName",
//					DataId:     "identity-nodeId",
//					DataStatus: carriertypespb.DataStatus_DataStatus_Deleted,
//					Status:     carriertypespb.CommonStatus_CommonStatus_NonNormal,
//					Credential: "{\"a\":\"b\"}",
//				},
//			},
//		}
//	}
//	taskdatas := func() TaskDataArray {
//		return TaskDataArray{
//			{
//				data: &carriertypespb.TaskPB{
//					IdentityId: "task-identity",
//					NodeId:     "task-nodeId",
//					NodeName:   "task-nodeName",
//					DataId:     "task-dataId",
//					DataStatus: carriertypespb.DataStatus_DataStatus_Deleted,
//					TaskId:     "task-taskId",
//					State:      carriertypespb.TaskState_TaskState_Unknown,
//					Reason:     "task-reason",
//					EventCount: 1,
//					Desc:       "task-desc",
//					TaskEvents: []*carriertypespb.TaskEvent{
//						{
//							TaskId: "1-evengine-taskId",
//						},
//					},
//				},
//			},
//		}
//	}
//	return NewBlock(newHeader(), metadatas(), resources(), identities(), taskdatas())
//}
//
//func newHeader() *Header {
//	return &Header{
//		ParentHash: []byte("parentHash"),
//		Version:    uint64(1),
//		Timestamp:  uint64(1000000),
//		Extra:      []byte("extraData"),
//	}
//}
//
//func TestBlockEncoding(t *testing.T) {
//	block := newBlock()
//	data, err := block.EncodePb()
//	if err != nil {
//		t.Fatal("EncodePB Failed")
//	}
//	t.Log(common.Bytes2Hex(data))
//
//	// test decode.
//	var dBlock Block
//	if err := dBlock.DecodePb(data); err != nil {
//		t.Fatal("decode error: ", err)
//	}
//
//	// check every element.
//	check := func(f string, got, want interface{}) {
//		if !reflect.DeepEqual(got, want) {
//			t.Errorf("%s mismatch: got %v, want %v", f, got, want)
//		}
//	}
//	check("metadatas.size", dBlock.metadatas.Len(), block.metadatas.Len())
//	check("resources.size", dBlock.resources.Len(), block.resources.Len())
//	check("identities.size", dBlock.identities.Len(), block.identities.Len())
//	check("taskdatas.size", dBlock.taskDatas.Len(), block.taskDatas.Len())
//
//	check("header.hash", dBlock.Hash(), block.Hash())
//	check("metadatas", dBlock.metadatas, block.metadatas)
//	check("resources", dBlock.resources, block.resources)
//	check("taskDatas", dBlock.taskDatas, block.taskDatas)
//
//	check("metadata.hash", dBlock.metadatas[0].Hash(), block.metadatas[0].Hash())
//	check("resources.hash", dBlock.resources[0].Hash(), block.resources[0].Hash())
//	check("identities.hash", dBlock.identities[0].Hash(), block.identities[0].Hash())
//	check("taskDatas.hash", dBlock.taskDatas[0].Hash(), block.taskDatas[0].Hash())
//}
