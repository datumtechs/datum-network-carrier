package policy

import (
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"github.com/datumtechs/datum-network-carrier/types"
)

func NewTaskDetailShowFromTaskData(input *types.Task) *carriertypespb.TaskDetail {
	return &carriertypespb.TaskDetail{
		Information: &carriertypespb.TaskDetailSummary{

			TaskId:                   input.GetTaskData().GetTaskId(),
			TaskName:                 input.GetTaskData().GetTaskName(),
			UserType:                 input.GetTaskData().GetUserType(),
			User:                     input.GetTaskData().GetUser(),
			Sender:                   input.GetTaskSender(),
			AlgoSupplier:             input.GetTaskData().GetAlgoSupplier(),
			DataSuppliers:            input.GetTaskData().GetDataSuppliers(),
			PowerSuppliers:           input.GetTaskData().GetPowerSuppliers(),
			Receivers:                input.GetTaskData().GetReceivers(),
			DataPolicyTypes:          input.GetTaskData().GetDataPolicyTypes(),
			DataPolicyOptions:        input.GetTaskData().GetDataPolicyOptions(),
			PowerPolicyTypes:         input.GetTaskData().GetPowerPolicyTypes(),
			PowerPolicyOptions:       input.GetTaskData().GetPowerPolicyOptions(),
			ReceiverPolicyTypes:      input.GetTaskData().GetReceiverPolicyTypes(),
			ReceiverPolicyOptions:    input.GetTaskData().GetReceiverPolicyOptions(),
			DataFlowPolicyTypes:      input.GetTaskData().GetDataFlowPolicyTypes(),
			DataFlowPolicyOptions:    input.GetTaskData().GetDataFlowPolicyOptions(),
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

func NewGlobalMetadataInfoFromMetadata(input *types.Metadata) *carrierapipb.GetGlobalMetadataDetail {
	response := &carrierapipb.GetGlobalMetadataDetail{
		Owner: input.GetData().GetOwner(),
		Information: &carriertypespb.MetadataDetail{
			MetadataSummary: &carriertypespb.MetadataSummary{
				/**
				MetadataId           string
				MetadataName         string
				MetadataType         MetadataType
				DataHash             string
				Desc                 string
				LocationType         DataLocationType
				DataType             OrigindataType
				Industry             string
				State                MetadataState
				PublishAt            uint64
				UpdateAt             uint64
				Nonce                uint64
				MetadataOption       string
				AllowExpose          bool
				TokenAddress         string
				*/
				MetadataId:     input.GetData().GetDataId(),
				MetadataName:   input.GetData().GetMetadataName(),
				MetadataType:   input.GetData().GetMetadataType(),
				DataHash:       input.GetData().GetDataHash(),
				Desc:           input.GetData().GetDesc(),
				LocationType:   input.GetData().GetLocationType(),
				DataType:       input.GetData().GetDataType(),
				Industry:       input.GetData().GetIndustry(),
				State:          input.GetData().GetState(),
				PublishAt:      input.GetData().GetPublishAt(),
				UpdateAt:       input.GetData().GetUpdateAt(),
				Nonce:          input.GetData().GetNonce(),
				MetadataOption: input.GetData().GetMetadataOption(),
				AllowExpose:    input.GetData().GetAllowExpose(),
				TokenAddress:   input.GetData().GetTokenAddress(),
			},
		},
	}
	return response
}

func NewLocalMetadataInfoFromMetadata(isInternal bool, input *types.Metadata) *carrierapipb.GetLocalMetadataDetail {
	response := &carrierapipb.GetLocalMetadataDetail{
		Owner: input.GetData().GetOwner(),
		Information: &carriertypespb.MetadataDetail{
			MetadataSummary: &carriertypespb.MetadataSummary{
				/**
				MetadataId           string
				MetadataName         string
				MetadataType         MetadataType
				DataHash             string
				Desc                 string
				LocationType         DataLocationType
				DataType             OrigindataType
				Industry             string
				State                MetadataState
				PublishAt            uint64
				UpdateAt             uint64
				Nonce                uint64
				MetadataOption       string
				AllowExpose          bool
				TokenAddress         string
				*/
				MetadataId:     input.GetData().GetDataId(),
				MetadataName:   input.GetData().GetMetadataName(),
				MetadataType:   input.GetData().GetMetadataType(),
				DataHash:       input.GetData().GetDataHash(),
				Desc:           input.GetData().GetDesc(),
				LocationType:   input.GetData().GetLocationType(),
				DataType:       input.GetData().GetDataType(),
				Industry:       input.GetData().GetIndustry(),
				State:          input.GetData().GetState(),
				PublishAt:      input.GetData().GetPublishAt(),
				UpdateAt:       input.GetData().GetUpdateAt(),
				Nonce:          input.GetData().GetNonce(),
				MetadataOption: input.GetData().GetMetadataOption(),
				AllowExpose:    input.GetData().GetAllowExpose(),
				TokenAddress:   input.GetData().GetTokenAddress(),
			},
			TotalTaskCount: 0,
		},
		IsInternal: isInternal,
	}
	return response
}

func NewGlobalMetadataInfoArrayFromMetadataArray(input types.MetadataArray) []*carrierapipb.GetGlobalMetadataDetail {
	result := make([]*carrierapipb.GetGlobalMetadataDetail, 0, input.Len())
	for _, metadata := range input {
		if metadata == nil {
			continue
		}
		result = append(result, NewGlobalMetadataInfoFromMetadata(metadata))
	}
	return result
}

func NewLocalMetadataInfoArrayFromMetadataArray(internalArr, publishArr types.MetadataArray) []*carrierapipb.GetLocalMetadataDetail {
	result := make([]*carrierapipb.GetLocalMetadataDetail, 0, internalArr.Len()+publishArr.Len())

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