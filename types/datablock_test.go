package types

import (
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/lib/types"
	"reflect"
	"testing"
)

func newBlock() *Block {
	metadatas := func() MetadataArray {
		return MetadataArray{
			&Metadata{
				data: &types.MetaData{
					Identity:    "identity",
					NodeId:      "NodeId",
					DataId:      "dataId",
					DataStatus:  "D",
					FilePath:    "/a/b/c",
					Desc:        "desc",
					Rows:        10,
					Columns:     100,
					FileType:    "csv",
					State:       "statue",
					HasTitleRow: true,
					ColumnMetaList: []*types.ColumnMeta{
						{
							Cindex:               0,
							Cname:                "cname",
							Ctype:                "ctype",
							Csize:                0,
						},
					},
				},
			},
		}
	}
	resources := func() ResourceArray {
		return ResourceArray{
			{
				data: &types.ResourceData{
					Identity:             "resource-identity",
					NodeId:               "resource-nodeId",
					NodeName:             "resource-nodeName",
					DataId:               "resource-dataId",
					DataStatus:           "resource-dataStatus",
					State:                "resource-state",
					TotalMem:             "1000MB",
					UsedMem:              "200MB",
					TotalProcessor:       0,
					TotalBandWidth:       "1000",
				},
			},
		}
	}
	identities := func() IdentityArray {
		return IdentityArray{
			{
				data: &types.IdentityData{
					Identity:             "identity-identity",
					NodeId:               "identity-nodeId",
					NodeName:             "identity-nodeName",
					DataId:               "identity-nodeId",
					DataStatus:           "identity-dataStatus",
					Status:               "identity-status",
					Credential:           "{\"a\":\"b\"}",
				},
			},
		}
	}
	taskdatas := func() TaskDataArray {
		return TaskDataArray{
			{
				data: &types.TaskData{
					Identity:             "task-identity",
					NodeId:               "task-nodeId",
					NodeName:             "task-nodeName",
					DataId:               "task-dataId",
					DataStatus:           "task-dataStatus",
					TaskId:               "task-taskId",
					State:                "task-state",
					Reason:               "task-reason",
					EventCount:           1,
					Desc:                 "task-desc",
					PartnerList:          []*types.Partner{
						{
							Alias:                "1-partner-alias",
							Identity:             "1-partner-identity",
							NodeId:               "1-partner-nodeId",
							NodeName:             "1-partner-nodeName",
						},
					},
					EventDataList:        []*types.EventData{
						{
							TaskId:               "1-evengine-taskId",
							EventType:            "1-evengine-eventType",
							EventAt:              0,
							EventContent:         "1-evengine-eventContent",
							Identity:             "1-evengine-identity",
						},
					},
				},
			},
		}
	}
	return NewBlock(newHeader(), metadatas(), resources(), identities(), taskdatas())
}

func newHeader() *Header {
	return &Header{
		ParentHash: []byte("parentHash"),
		Version:    uint64(1),
		Timestamp: uint64(1000000),
		Extra: []byte("extraData"),
	}
}

func TestBlockEncoding(t *testing.T) {
	block := newBlock()
	data, err := block.EncodePb()
	if err != nil {
		t.Fatal("EncodePB Failed")
	}
	t.Log(common.Bytes2Hex(data))

	// test decode.
	var dBlock Block
	if err := dBlock.DecodePb(data); err != nil {
		t.Fatal("decode error: ", err)
	}

	// check every element.
	check := func(f string, got, want interface{}) {
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s mismatch: got %v, want %v", f, got, want)
		}
	}
	check("metadatas.size", dBlock.metadatas.Len(), block.metadatas.Len())
	check("resources.size", dBlock.resources.Len(), block.resources.Len())
	check("identities.size", dBlock.identities.Len(), block.identities.Len())
	check("taskdatas.size", dBlock.taskDatas.Len(), block.taskDatas.Len())

	check("header.hash", dBlock.Hash(), block.Hash())
	check("metadatas", dBlock.metadatas, block.metadatas)
	check("resources", dBlock.resources, block.resources)
	check("taskDatas", dBlock.taskDatas, block.taskDatas)

	check("metadata.hash", dBlock.metadatas[0].Hash(), block.metadatas[0].Hash())
	check("resources.hash", dBlock.resources[0].Hash(), block.resources[0].Hash())
	check("identities.hash", dBlock.identities[0].Hash(), block.identities[0].Hash())
	check("taskDatas.hash", dBlock.taskDatas[0].Hash(), block.taskDatas[0].Hash())
}