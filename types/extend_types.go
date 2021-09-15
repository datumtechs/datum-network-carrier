package types

import (
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
)

func NewTaskDetailShowFromTaskData(input *Task, role apicommonpb.TaskRole) *pb.TaskDetailShow {
	taskData := input.GetTaskData()
	detailShow := &pb.TaskDetailShow{
		TaskId:   taskData.GetTaskId(),
		TaskName: taskData.GetTaskName(),
		//TODO: 需要确认部分
		//Role:     role,
		Sender: &apicommonpb.TaskOrganization{
			PartyId:    taskData.GetPartyId(),
			NodeName:   taskData.GetNodeName(),
			NodeId:     taskData.GetNodeId(),
			IdentityId: taskData.GetIdentityId(),
		},
		AlgoSupplier: &apicommonpb.TaskOrganization{
			PartyId:    taskData.GetPartyId(),
			NodeName:   taskData.GetNodeName(),
			NodeId:     taskData.GetNodeId(),
			IdentityId: taskData.GetIdentityId(),
		},
		DataSuppliers:  make([]*pb.TaskDataSupplierShow, 0, len(taskData.GetDataSuppliers())),
		PowerSuppliers: make([]*pb.TaskPowerSupplierShow, 0, len(taskData.GetPowerSuppliers())),
		Receivers:      taskData.GetReceivers(),
		CreateAt:       taskData.GetCreateAt(),
		StartAt:        taskData.GetStartAt(),
		EndAt:          taskData.GetEndAt(),
		State:          taskData.GetState(),
		OperationCost: &apicommonpb.TaskResourceCostDeclare{
			Processor: taskData.GetOperationCost().GetProcessor(),
			Memory:    taskData.GetOperationCost().GetMemory(),
			Bandwidth: taskData.GetOperationCost().GetBandwidth(),
			Duration:  taskData.GetOperationCost().GetDuration(),
		},
	}
	// DataSupplier
	for _, metadataSupplier := range taskData.GetDataSuppliers() {
		dataSupplier := &pb.TaskDataSupplierShow{
			Organization: &apicommonpb.TaskOrganization{
				PartyId:    metadataSupplier.GetOrganization().GetPartyId(),
				NodeName:   metadataSupplier.GetOrganization().GetNodeName(),
				NodeId:     metadataSupplier.GetOrganization().GetNodeId(),
				IdentityId: metadataSupplier.GetOrganization().GetIdentityId(),
			},
			MetadataId:   metadataSupplier.GetMetadataId(),
			MetadataName: metadataSupplier.GetMetadataName(),
		}
		detailShow.DataSuppliers = append(detailShow.DataSuppliers, dataSupplier)
	}
	// powerSupplier
	for _, data := range taskData.GetPowerSuppliers() {
		detailShow.PowerSuppliers = append(detailShow.PowerSuppliers, &pb.TaskPowerSupplierShow{
			Organization: &apicommonpb.TaskOrganization{
				PartyId:    data.GetOrganization().GetPartyId(),
				NodeName:   data.GetOrganization().GetNodeName(),
				NodeId:     data.GetOrganization().GetNodeId(),
				IdentityId: data.GetOrganization().GetIdentityId(),
			},
			PowerInfo: &libtypes.ResourceUsageOverview{
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

func NewTaskEventFromAPIEvent(input []*libtypes.TaskEvent) []*pb.TaskEventShow {
	result := make([]*pb.TaskEventShow, 0, len(input))
	for _, event := range input {
		result = append(result, &pb.TaskEventShow{
			TaskId:   event.GetTaskId(),
			Type:     event.GetType(),
			CreateAt: event.GetCreateAt(),
			Content:  event.GetContent(),
			Owner: &apicommonpb.Organization{
				IdentityId: event.GetIdentityId(),
			},
		})
	}
	return result
}

func NewOrgMetadataInfoFromMetadata(input *Metadata) *pb.GetMetadataDetailResponse {
	response := &pb.GetMetadataDetailResponse{
		Owner: &apicommonpb.Organization{
			NodeName:   input.data.GetNodeName(),
			NodeId:     input.data.GetNodeId(),
			IdentityId: input.data.GetIdentityId(),
		},
		Information: &libtypes.MetadataDetail{
			MetadataSummary: &libtypes.MetadataSummary{
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
		identityId: input.data.GetIdentityId(),
		total: &resource{
			mem:       input.data.GetTotalMem(),
			processor: input.data.GetTotalProcessor(),
			bandwidth: input.data.GetTotalBandwidth(),
		},
		used: &resource{
			mem:       input.data.GetUsedMem(),
			processor: input.data.GetUsedProcessor(),
			bandwidth: input.data.GetUsedBandwidth(),
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
