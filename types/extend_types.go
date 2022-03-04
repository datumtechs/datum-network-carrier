package types

import (
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
)

func NewTaskDetailShowFromTaskData(input *Task) *pb.TaskDetailShow {
	taskData := input.GetTaskData()
	detailShow := &pb.TaskDetailShow{
		TaskId:   taskData.GetTaskId(),
		TaskName: taskData.GetTaskName(),
		UserType: taskData.GetUserType(),
		User:     taskData.GetUser(),
		Sender: &apicommonpb.TaskOrganization{
			PartyId:    input.GetTaskSender().GetPartyId(),
			NodeName:   input.GetTaskSender().GetNodeName(),
			NodeId:     input.GetTaskSender().GetNodeId(),
			IdentityId: input.GetTaskSender().GetIdentityId(),
		},
		AlgoSupplier: &apicommonpb.TaskOrganization{
			PartyId:    input.GetTaskData().GetAlgoSupplier().GetPartyId(),
			NodeName:   input.GetTaskData().GetAlgoSupplier().GetNodeName(),
			NodeId:     input.GetTaskData().GetAlgoSupplier().GetNodeId(),
			IdentityId: input.GetTaskData().GetAlgoSupplier().GetIdentityId(),
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
		UpdateAt: taskData.GetEndAt(), // The endAt of the task is the updateAt in the data center database
	}

	// DataSupplier
	for _, dataSupplier := range taskData.GetDataSuppliers() {
		metadataId, err := FetchMetedataIdByPartyId(dataSupplier.GetPartyId(), taskData.GetDataPolicyType(), taskData.GetDataPolicyOption())
		if nil != err {
			log.WithError(err).Errorf("failed to fetch metadataId by partyId from DataPolicyOption, taskId: {%s}, partyId: {%s}",
				taskData.GetTaskId(), dataSupplier.GetPartyId())
		}
		metadataName, err := FetchMetedataNameByPartyId(dataSupplier.GetPartyId(), taskData.GetDataPolicyType(), taskData.GetDataPolicyOption())
		if nil != err {
			log.WithError(err).Errorf("failed to fetch metadataName by partyId from DataPolicyOption, taskId: {%s}, partyId: {%s}",
				taskData.GetTaskId(), dataSupplier.GetPartyId())
		}
		supplier := &pb.TaskDataSupplierShow{
			Organization: &apicommonpb.TaskOrganization{
				PartyId:    dataSupplier.GetPartyId(),
				NodeName:   dataSupplier.GetNodeName(),
				NodeId:     dataSupplier.GetNodeId(),
				IdentityId: dataSupplier.GetIdentityId(),
			},
			MetadataId:  metadataId,
			MetadataName: metadataName,
		}
		detailShow.DataSuppliers = append(detailShow.DataSuppliers, supplier)
	}
	// powerSupplier
	for _, data := range taskData.GetPowerSuppliers() {

		var option *libtypes.TaskPowerResourceOption
		for _, op := range taskData.GetPowerResourceOptions() {
			if data.GetPartyId() == op.GetPartyId() {
				option = op
				break
			}
		}
		supplier := &pb.TaskPowerSupplierShow{
			Organization: &apicommonpb.TaskOrganization{
				PartyId:    data.GetPartyId(),
				NodeName:   data.GetNodeName(),
				NodeId:     data.GetNodeId(),
				IdentityId: data.GetIdentityId(),
			},
			PowerInfo: &libtypes.ResourceUsageOverview{
				TotalMem:       option.GetResourceUsedOverview().GetTotalMem(),
				UsedMem:        option.GetResourceUsedOverview().GetUsedMem(),
				TotalProcessor: option.GetResourceUsedOverview().GetTotalProcessor(),
				UsedProcessor:  option.GetResourceUsedOverview().GetUsedProcessor(),
				TotalBandwidth: option.GetResourceUsedOverview().GetTotalBandwidth(),
				UsedBandwidth:  option.GetResourceUsedOverview().GetUsedBandwidth(),
				TotalDisk:      option.GetResourceUsedOverview().GetTotalDisk(),
				UsedDisk:       option.GetResourceUsedOverview().GetUsedDisk(),
			},
		}
		detailShow.PowerSuppliers = append(detailShow.PowerSuppliers, supplier)
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
			PartyId: event.GetPartyId(),
		})
	}
	return result
}

func NewGlobalMetadataInfoFromMetadata(input *Metadata) *pb.GetGlobalMetadataDetailResponse {
	response := &pb.GetGlobalMetadataDetailResponse{
		Owner: input.GetData().GetOwner(),
		Information: &libtypes.MetadataDetail{
			MetadataSummary: &libtypes.MetadataSummary{
				/**
				MetadataId   string
				MetadataName string
				MetadataType uint32
				FileHash     string
				Desc         string
				FileType     common.OriginFileType
				Industry     string
				State        common.MetadataState
				// v 2.0
				PublishAt            uint64   `pro
				UpdateAt             uint64   `pro
				Nonce                uint64   `pro
				MetadataOption       string   `pro
				 */
				MetadataId: input.GetData().GetMetadataId(),
				MetadataName: input.GetData().GetMetadataName(),
				Desc:       input.GetData().GetDesc(),
				FileType:   input.GetData().GetFileType(),
				Industry:   input.GetData().GetIndustry(),
				State:      input.GetData().GetState(),
				PublishAt:  input.GetData().GetPublishAt(),
				UpdateAt:   input.GetData().GetUpdateAt(),
			},
		},
	}
	return response
}

func NewLocalMetadataInfoFromMetadata(isInternal bool, input *Metadata) *pb.GetLocalMetadataDetailResponse {
	response := &pb.GetLocalMetadataDetailResponse{
		Owner: &apicommonpb.Organization{
			NodeName:   input.data.GetNodeName(),
			NodeId:     input.data.GetNodeId(),
			IdentityId: input.data.GetIdentityId(),
		},
		Information: &libtypes.MetadataDetail{
			MetadataSummary: &libtypes.MetadataSummary{
				MetadataId: input.GetData().GetDataId(),
				OriginId:   input.GetData().GetOriginId(),
				TableName:  input.GetData().GetTableName(),
				Desc:       input.GetData().GetDesc(),
				FilePath:   input.GetData().GetFilePath(),
				Rows:       input.GetData().GetRows(),
				Columns:    input.GetData().GetColumns(),
				Size_:      input.GetData().GetSize_(),
				FileType:   input.GetData().GetFileType(),
				HasTitle:   input.GetData().GetHasTitle(),
				Industry:   input.GetData().GetIndustry(),
				State:      input.GetData().GetState(),
				PublishAt:  input.GetData().GetPublishAt(),
				UpdateAt:   input.GetData().GetUpdateAt(),
			},
			MetadataColumns: input.GetData().GetMetadataColumns(),
		},
		IsInternal: isInternal,
	}
	return response
}

func NewGlobalMetadataInfoArrayFromMetadataArray(input MetadataArray) []*pb.GetGlobalMetadataDetailResponse {
	result := make([]*pb.GetGlobalMetadataDetailResponse, 0, input.Len())
	for _, metadata := range input {
		if metadata == nil {
			continue
		}
		result = append(result, NewGlobalMetadataInfoFromMetadata(metadata))
	}
	return result
}

func NewLocalMetadataInfoArrayFromMetadataArray(internalArr, publishArr MetadataArray) []*pb.GetLocalMetadataDetailResponse {
	result := make([]*pb.GetLocalMetadataDetailResponse, 0, internalArr.Len()+publishArr.Len())

	for _, metadata := range internalArr {
		if metadata == nil {
			continue
		}
		result = append(result, NewLocalMetadataInfoFromMetadata(true, metadata))
	}

	for _, metadata := range publishArr {
		if metadata == nil {
			continue
		}
		result = append(result, NewLocalMetadataInfoFromMetadata(false, metadata))
	}

	return result
}
