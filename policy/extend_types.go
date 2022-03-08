package policy

import (
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
)

func NewTaskDetailShowFromTaskData(input *types.Task) *pb.TaskDetailShow {
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
			MetadataId:   metadataId,
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

func NewGlobalMetadataInfoFromMetadata(input *types.Metadata) *pb.GetGlobalMetadataDetailResponse {
	response := &pb.GetGlobalMetadataDetailResponse{
		Owner: input.GetData().GetOwner(),
		Information: &libtypes.MetadataDetail{
			MetadataSummary: &libtypes.MetadataSummary{
				MetadataId:     input.GetData().GetDataId(),
				MetadataName:   input.GetData().GetMetadataName(),
				MetadataType:   input.GetData().GetMetadataType(),
				FileHash:       input.GetData().GetFileHash(),
				Desc:           input.GetData().GetDesc(),
				FileType:       input.GetData().GetFileType(),
				Industry:       input.GetData().GetIndustry(),
				State:          input.GetData().GetState(),
				PublishAt:      input.GetData().GetPublishAt(),
				UpdateAt:       input.GetData().GetUpdateAt(),
				Nonce:          input.GetData().GetNonce(),
				MetadataOption: input.GetData().GetMetadataOption(),
			},
		},
	}
	return response
}

func NewLocalMetadataInfoFromMetadata(isInternal bool, input *types.Metadata) *pb.GetLocalMetadataDetailResponse {
	response := &pb.GetLocalMetadataDetailResponse{
		Owner: input.GetData().GetOwner(),
		Information: &libtypes.MetadataDetail{
			MetadataSummary: &libtypes.MetadataSummary{
				MetadataId:     input.GetData().GetDataId(),
				MetadataName:   input.GetData().GetMetadataName(),
				MetadataType:   input.GetData().GetMetadataType(),
				FileHash:       input.GetData().GetFileHash(),
				Desc:           input.GetData().GetDesc(),
				FileType:       input.GetData().GetFileType(),
				Industry:       input.GetData().GetIndustry(),
				State:          input.GetData().GetState(),
				PublishAt:      input.GetData().GetPublishAt(),
				UpdateAt:       input.GetData().GetUpdateAt(),
				Nonce:          input.GetData().GetNonce(),
				MetadataOption: input.GetData().GetMetadataOption(),
			},
			//TotalTaskCount: ,
		},
		IsInternal: isInternal,
	}
	return response
}

func NewGlobalMetadataInfoArrayFromMetadataArray(input types.MetadataArray) []*pb.GetGlobalMetadataDetailResponse {
	result := make([]*pb.GetGlobalMetadataDetailResponse, 0, input.Len())
	for _, metadata := range input {
		if metadata == nil {
			continue
		}
		result = append(result, NewGlobalMetadataInfoFromMetadata(metadata))
	}
	return result
}

func NewLocalMetadataInfoArrayFromMetadataArray(internalArr, publishArr types.MetadataArray) []*pb.GetLocalMetadataDetailResponse {
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
