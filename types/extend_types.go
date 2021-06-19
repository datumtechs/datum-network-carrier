package types

import "github.com/RosettaFlow/Carrier-Go/lib/center/api"

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
				MetaId:     metadataSupplier.GetMetaId(),
				MetaName:   metadataSupplier.GetMetaName(),
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

		taskDetailShowArray = append(taskDetailShowArray, detailShow)
	}
	return taskDetailShowArray
}

func NewTaskEventFromAPIEvent(input []*api.TaskEvent) []*TaskEvent  {
	return nil
}

func NewOrgMetaDataInfoFromMetadata(input *Metadata) *OrgMetaDataInfo  {
	return nil
}

func NewOrgMetaDataInfoArrayFromMetadataArray(input MetadataArray) []*OrgMetaDataInfo  {
	return nil
}