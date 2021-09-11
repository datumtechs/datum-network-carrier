package types

import (
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
)

func NewTaskDetailShowFromTaskData(input *Task, role apipb.TaskRole) *pb.TaskDetailShow {
	taskData := input.GetTaskData()
	detailShow := &pb.TaskDetailShow{
		TaskId:   taskData.GetTaskId(),
		TaskName: taskData.GetTaskName(),
		//TODO: 需要确认部分
		//Role:     role,
		Owner: &apipb.TaskOrganization{
			PartyId:    taskData.GetPartyId(),
			NodeName:       taskData.GetNodeName(),
			NodeId:     taskData.GetNodeId(),
			IdentityId: taskData.GetIdentityId(),
		},
		AlgoSupplier: &apipb.TaskOrganization{
			PartyId:    taskData.GetPartyId(),
			NodeName:       taskData.GetNodeName(),
			NodeId:     taskData.GetNodeId(),
			IdentityId: taskData.GetIdentityId(),
		},
		DataSupplier:  make([]*pb.TaskDataSupplierShow, 0, len(taskData.GetDataSuppliers())),
		PowerSupplier: make([]*pb.TaskPowerSupplierShow, 0, len(taskData.GetPowerSuppliers())),
		Receivers:     taskData.GetReceivers(),
		CreateAt:      taskData.GetCreateAt(),
		StartAt:       taskData.GetStartAt(),
		EndAt:         taskData.GetEndAt(),
		State:         taskData.GetState(),
		OperationCost: &apipb.TaskResourceCostDeclare{
			Processor: taskData.GetOperationCost().GetProcessor(),
			Memory:       taskData.GetOperationCost().GetMemory(),
			Bandwidth: taskData.GetOperationCost().GetBandwidth(),
			Duration:  taskData.GetOperationCost().GetDuration(),
		},
	}
	// DataSupplier
	for _, metadataSupplier := range taskData.GetDataSuppliers() {
		dataSupplier := &pb.TaskDataSupplierShow{
			Organization: &apipb.TaskOrganization{
				PartyId:    metadataSupplier.GetOrganization().GetPartyId(),
				NodeName:   metadataSupplier.GetOrganization().GetNodeName(),
				NodeId:     metadataSupplier.GetOrganization().GetNodeId(),
				IdentityId: metadataSupplier.GetOrganization().GetIdentityId(),
			},
			MetadataId:   metadataSupplier.GetMetadataId(),
			MetadataName: metadataSupplier.GetMetadataName(),
		}
		detailShow.DataSupplier = append(detailShow.DataSupplier, dataSupplier)
	}
	// powerSupplier
	for _, data := range taskData.GetPowerSuppliers() {
		detailShow.PowerSupplier = append(detailShow.PowerSupplier, &pb.TaskPowerSupplierShow{
			Organization: &apipb.TaskOrganization{
				PartyId:    data.GetOrganization().GetPartyId(),
				NodeName:   data.GetOrganization().GetNodeName(),
				NodeId:     data.GetOrganization().GetNodeId(),
				IdentityId: data.GetOrganization().GetIdentityId(),
			},
			PowerInfo: &libTypes.ResourceUsageOverview{
				TotalMem:       data.GetResourceUsedOverview().GetTotalMem(),
				UsedMem:        data.GetResourceUsedOverview().GetUsedMem(),
				TotalProcessor: data.GetResourceUsedOverview().GetTotalProcessor(),
				UsedProcessor:  data.GetResourceUsedOverview().GetUsedProcessor(),
				TotalBandwidth: data.GetResourceUsedOverview().GetTotalBandwidth(),
				UsedBandwidth:  data.GetResourceUsedOverview().GetUsedBandwidth(),
			},
		})
	}
	return detailShow
}

func NewTaskEventFromAPIEvent(input []*libTypes.TaskEvent) []*pb.TaskEventShow {
	result := make([]*pb.TaskEventShow, 0, len(input))
	for _, event := range input {
		result = append(result, &pb.TaskEventShow{
			TaskId:   event.GetTaskId(),
			Type:     event.GetType(),
			CreateAt: event.GetCreateAt(),
			Content:  event.GetContent(),
			Owner:    &apipb.Organization{
				IdentityId: event.GetIdentityId(),
			},
		})
	}
	return result
}

func NewOrgMetadataInfoFromMetadata(input *Metadata) *pb.GetMetadataDetailResponse {
	response := &pb.GetMetadataDetailResponse{
		Owner: &apipb.Organization{
			NodeName:   input.data.GetNodeName(),
			NodeId:     input.data.GetNodeId(),
			IdentityId: input.data.GetIdentityId(),
		},
		Information: &libTypes.MetadataDetail{
			MetadataSummary: &libTypes.MetadataSummary{
				MetadataId: input.data.GetDataId(),
				OriginId:   input.data.GetOriginId(),
				TableName:  input.data.GetTableName(),
				Desc:       input.data.GetDesc(),
				FilePath:   input.data.GetFilePath(),
				Rows:       input.data.GetRows(),
				Columns:    input.data.GetColumns(),
				Size_:      input.data.GetSize_(),
				FileType:   input.data.GetFileType(),
				HasTitle:   input.data.GetHasTitle(),
				State:      input.data.GetState(),
			},
			MetadataColumns: input.data.GetMetadataColumns(),
		},
	}
	return response
}

func NewOrgMetadataInfoArrayFromMetadataArray(input MetadataArray) []*pb.GetMetadataDetailResponse {
	result := make([]*pb.GetMetadataDetailResponse, 0, input.Len())
	for _, metadata := range input {
		if metadata == nil {
			continue
		}
		result = append(result, NewOrgMetadataInfoFromMetadata(metadata))
	}
	return result
}

func NewOrgResourceFromResource(input *Resource) *RemoteResourceTable {
	return &RemoteResourceTable{
		identityId: input.data.IdentityId,
		total: &resource{
			mem:       input.data.TotalMem,
			processor: input.data.TotalProcessor,
			bandwidth: input.data.TotalBandwidth,
		},
		used: &resource{
			mem:       input.data.UsedMem,
			processor: input.data.UsedProcessor,
			bandwidth: input.data.UsedBandwidth,
		},
	}
}

//func NewOrgResourceArrayFromResourceArray(input ResourceArray) []*RemoteResourceTable {
//	result := make([]*RemoteResourceTable, input.Len())
//	for i, resource := range input {
//		result[i] = NewOrgResourceFromResource(resource)
//	}
//	return result
//}
