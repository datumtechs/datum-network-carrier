package types

import (
	"github.com/RosettaFlow/Carrier-Go/lib/center/api"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
)

func NewTaskDetailShowArrayFromTaskDataArray(input TaskDataArray) []*TaskDetailShow {
	taskDetailShowArray := make([]*TaskDetailShow, input.Len())
	for _, task := range input {
		taskData := task.TaskData()
		detailShow := &TaskDetailShow{
			TaskId:        taskData.GetTaskId(),
			TaskName:      taskData.GetTaskName(),
			Owner:         &NodeAlias{
				Name:       taskData.GetNodeName(),
				NodeId:     taskData.GetNodeId(),
				IdentityId: taskData.GetIdentity(),
			},
			AlgoSupplier:  &NodeAlias{
				Name:       taskData.GetNodeName(),
				NodeId:     taskData.GetNodeId(),
				IdentityId: taskData.GetIdentity(),
			},
			DataSupplier:  make([]*TaskDataSupplierShow, len(taskData.GetMetadataSupplier())),
			PowerSupplier: make([]*TaskPowerSupplierShow, len(taskData.GetResourceSupplier())),
			Receivers:     make([]*NodeAlias, len(taskData.GetReceivers())),
			CreateAt:      taskData.GetCreateAt(),
			EndAt:         taskData.GetEndAt(),
			State:         taskData.GetState(),
			OperationCost: &TaskOperationCost{
				Processor: uint64(taskData.GetTaskResource().GetCostProcessor()),
				Mem:       taskData.GetTaskResource().GetCostMem(),
				Bandwidth: taskData.GetTaskResource().GetCostBandwidth(),
				Duration:  taskData.GetTaskResource().GetDuration(),
			},
		}
		// DataSupplier
		for _, metadataSupplier := range taskData.GetMetadataSupplier() {
			dataSupplier := &TaskDataSupplierShow{
				MemberInfo: &NodeAlias{
					Name:       metadataSupplier.GetOrganization().GetNodeName(),
					NodeId:     metadataSupplier.GetOrganization().GetNodeId(),
					IdentityId: metadataSupplier.GetOrganization().GetIdentity(),
				},
				MetaDataId:     metadataSupplier.GetMetaId(),
				MetaDataName:   metadataSupplier.GetMetaName(),
			}
			detailShow.DataSupplier = append(detailShow.DataSupplier, dataSupplier)
		}
		// powerSupplier
		for _, data := range taskData.GetResourceSupplier() {
			detailShow.PowerSupplier = append(detailShow.PowerSupplier, &TaskPowerSupplierShow{
				MemberInfo:    &NodeAlias{
					Name:       data.GetOrganization().GetNodeName(),
					NodeId:     data.GetOrganization().GetNodeId(),
					IdentityId: data.GetOrganization().GetIdentity(),
				},
				ResourceUsage: &ResourceUsage{
					TotalMem:       data.GetResourceUsedOverview().GetTotalMem(),
					UsedMem:        data.GetResourceUsedOverview().GetUsedMem(),
					TotalProcessor: uint64(data.GetResourceUsedOverview().GetTotalProcessor()),
					UsedProcessor:  uint64(data.GetResourceUsedOverview().GetUsedProcessor()),
					TotalBandwidth: data.GetResourceUsedOverview().GetTotalBandwidth(),
					UsedBandwidth:  data.GetResourceUsedOverview().GetUsedBandwidth(),
				},
			})
		}
		// Receivers
		for _, receiver := range taskData.GetReceivers() {
			detailShow.Receivers = append(detailShow.Receivers, &NodeAlias{
				Name:       receiver.GetReceiver().GetNodeName(),
				NodeId:     receiver.GetReceiver().GetNodeId(),
				IdentityId: receiver.GetReceiver().GetIdentity(),
			})
		}
		taskDetailShowArray = append(taskDetailShowArray, detailShow)
	}
	return taskDetailShowArray
}

func NewTaskEventFromAPIEvent(input []*api.TaskEvent) []*TaskEvent  {
	result := make([]*TaskEvent, len(input))
	for _, event := range input {
		result = append(result, &TaskEvent{
			TaskId:   event.GetTaskId(),
			Type:     event.GetType(),
			CreateAt: event.GetCreateAt(),
			Content:  event.GetContent(),
			Owner:    &NodeAlias{
				Name:       event.GetOwner().GetName(),
				NodeId:     event.GetOwner().GetNodeId(),
				IdentityId: event.GetOwner().GetIdentityId(),
			},
		})
	}
	return result
}

func NewOrgMetaDataInfoFromMetadata(input *Metadata) *OrgMetaDataInfo  {
	orgMetaDataInfo := &OrgMetaDataInfo{
		Owner:    &NodeAlias{
			Name:       input.data.GetNodeName(),
			NodeId:     input.data.GetNodeId(),
			IdentityId: input.data.GetIdentity(),
		},
		MetaData: &MetaDataInfo{
			MetaDataSummary: &MetaDataSummary{
				OriginId:  input.data.GetOriginId(),
				TableName: input.data.GetTableName(),
				Desc:      input.data.GetDesc(),
				FilePath:  input.data.GetFilePath(),
				Rows:      uint32(input.data.GetRows()),
				Columns:   uint32(input.data.GetColumns()),
				Size:      uint32(input.data.GetSize_()),
				FileType:  input.data.GetFileType(),
				HasTitle:  input.data.GetHasTitleRow(),
				State:     input.data.GetState(),
			},
			ColumnMetas:     make([]*libtypes.ColumnMeta, len(input.data.GetColumnMetaList())),
		},
	}
	for _, columnMeta := range input.data.GetColumnMetaList() {
		orgMetaDataInfo.MetaData.ColumnMetas = append(orgMetaDataInfo.MetaData.ColumnMetas, &libtypes.ColumnMeta{
			Cindex:   columnMeta.GetCindex(),
			Cname:    columnMeta.GetCname(),
			Ctype:    columnMeta.GetCtype(),
			Csize:    columnMeta.GetCsize(),
			Ccomment: columnMeta.GetCcomment(),
		})
	}
	return orgMetaDataInfo
}

func NewOrgMetaDataInfoArrayFromMetadataArray(input MetadataArray) []*OrgMetaDataInfo  {
	result := make([]*OrgMetaDataInfo, input.Len())
	for _, metadata := range input {
		result = append(result, NewOrgMetaDataInfoFromMetadata(metadata))
	}
	return result
}