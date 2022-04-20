package policy

import (
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
)

func NewTaskDetailShowFromTaskData(input *types.Task) *libtypes.TaskDetail {
	return &libtypes.TaskDetail{
		Information: &libtypes.TaskDetailSummary{
			TaskId:                   input.GetTaskData().GetTaskId(),
			TaskName:                 input.GetTaskData().GetTaskName(),
			UserType:                 input.GetTaskData().GetUserType(),
			User:                     input.GetTaskData().GetUser(),
			Sender:                   input.GetTaskSender(),
			AlgoSupplier:             input.GetTaskData().GetAlgoSupplier(),
			DataSuppliers:            input.GetTaskData().GetDataSuppliers(),
			PowerSuppliers:           input.GetTaskData().GetPowerSuppliers(),
			Receivers:                input.GetTaskData().GetReceivers(),
			DataPolicyType:           input.GetTaskData().GetDataPolicyType(),
			DataPolicyOption:         input.GetTaskData().GetDataPolicyOption(),
			PowerPolicyType:          input.GetTaskData().GetPowerPolicyType(),
			PowerPolicyOption:        input.GetTaskData().GetPowerPolicyOption(),
			DataFlowPolicyType:       input.GetTaskData().GetDataFlowPolicyType(),
			DataFlowPolicyOption:     input.GetTaskData().GetDataFlowPolicyOption(),
			OperationCost:            input.GetTaskData().GetOperationCost(),
			AlgorithmCode:            input.GetTaskData().GetAlgorithmCode(),
			MetaAlgorithmId:          input.GetTaskData().GetMetaAlgorithmId(),
			AlgorithmCodeExtraParams: input.GetTaskData().GetAlgorithmCodeExtraParams(),
			PowerResourceOptions:     input.GetTaskData().GetPowerResourceOptions(),
			State:                    input.GetTaskData().GetState(),
			Reason:                   input.GetTaskData().GetReason(),
			Desc:                     input.GetTaskData().GetDesc(),
			CreateAt:                 input.GetTaskData().GetCreateAt(),
			StartAt:                  input.GetTaskData().GetStartAt(),
			EndAt:                    input.GetTaskData().GetEndAt(),
			Sign:                     input.GetTaskData().GetSign(),
			Nonce:                    input.GetTaskData().GetNonce(),
			UpdateAt:                 input.GetTaskData().GetEndAt(), // The endAt of the task is the updateAt in the data center database
		},
	}
}

func NewGlobalMetadataInfoFromMetadata(input *types.Metadata) *pb.GetGlobalMetadataDetail {
	response := &pb.GetGlobalMetadataDetail{
		Owner: input.GetData().GetOwner(),
		Information: &libtypes.MetadataDetail{
			MetadataSummary: &libtypes.MetadataSummary{
				MetadataId:     input.GetData().GetDataId(),
				MetadataName:   input.GetData().GetMetadataName(),
				MetadataType:   input.GetData().GetMetadataType(),
				DataHash:       input.GetData().GetDataHash(),
				Desc:           input.GetData().GetDesc(),
				DataType:       input.GetData().GetDataType(),
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

func NewLocalMetadataInfoFromMetadata(isInternal bool, input *types.Metadata) *pb.GetLocalMetadataDetail {
	response := &pb.GetLocalMetadataDetail{
		Owner: input.GetData().GetOwner(),
		Information: &libtypes.MetadataDetail{
			MetadataSummary: &libtypes.MetadataSummary{
				MetadataId:     input.GetData().GetDataId(),
				MetadataName:   input.GetData().GetMetadataName(),
				MetadataType:   input.GetData().GetMetadataType(),
				DataHash:       input.GetData().GetDataHash(),
				Desc:           input.GetData().GetDesc(),
				DataType:       input.GetData().GetDataType(),
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

func NewGlobalMetadataInfoArrayFromMetadataArray(input types.MetadataArray) []*pb.GetGlobalMetadataDetail {
	result := make([]*pb.GetGlobalMetadataDetail, 0, input.Len())
	for _, metadata := range input {
		if metadata == nil {
			continue
		}
		result = append(result, NewGlobalMetadataInfoFromMetadata(metadata))
	}
	return result
}

func NewLocalMetadataInfoArrayFromMetadataArray(internalArr, publishArr types.MetadataArray) []*pb.GetLocalMetadataDetail {
	result := make([]*pb.GetLocalMetadataDetail, 0, internalArr.Len()+publishArr.Len())

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
