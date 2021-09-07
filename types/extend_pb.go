package types

import (
	"github.com/RosettaFlow/Carrier-Go/lib/center/api"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
)

// NewMetaDataSaveRequest converts Metadata object to MetaDataSaveRequest object.
func NewMetaDataSaveRequest(metadata *Metadata) *api.MetaDataSaveRequest {
	request := &api.MetaDataSaveRequest{
		MetaSummary: &libTypes.MetaDataSummary{
			MetaDataId: metadata.data.DataId,
			OriginId:   metadata.data.OriginId,
			TableName:  metadata.data.TableName,
			Desc:       metadata.data.Desc,
			FilePath:   metadata.data.FilePath,
			Rows:       uint32(metadata.data.Rows),
			Columns:    uint32(metadata.data.Columns),
			Size_:      metadata.data.Size_,
			FileType:   metadata.data.FileType,
			HasTitle:   metadata.data.HasTitle,
			State:      metadata.data.State,
		},
		ColumnMeta: make([]*libTypes.MetadataColumn, 0),
		Owner: &apipb.Organization{
			NodeName:   metadata.data.GetNodeName(),
			NodeId:     metadata.data.GetNodeId(),
			IdentityId: metadata.data.GetIdentityId(),
		},
	}
	for _, column := range metadata.data.MetadataColumns {
		request.ColumnMeta = append(request.ColumnMeta, column)
	}
	return request
}

func NewMetaDataRevokeRequest(metadata *Metadata) *api.RevokeMetaDataRequest {
	request := &api.RevokeMetaDataRequest{
		Owner: &apipb.Organization{
			IdentityId: metadata.MetadataData().IdentityId,
			NodeId:     metadata.MetadataData().NodeId,
			NodeName:   metadata.MetadataData().NodeName,
		},
		MetaDataId: metadata.MetadataData().DataId,
	}
	return request
}

func NewPublishPowerRequest(resource *Resource) *api.PublishPowerRequest {
	request := &api.PublishPowerRequest{
		Owner: &apipb.Organization{
			NodeName:   resource.data.GetNodeName(),
			NodeId:     resource.data.GetNodeId(),
			IdentityId: resource.data.GetIdentityId(),
		},
		PowerId: resource.data.DataId,
		Information: &api.PurePower{
			Mem:       resource.data.GetTotalMem(),
			Processor: uint32(resource.data.GetTotalProcessor()),
			Bandwidth: resource.data.GetTotalBandwidth(),
		},
	}
	return request
}

func RevokePowerRequest(resource *Resource) *api.RevokePowerRequest {
	request := &api.RevokePowerRequest{
		Owner: &apipb.Organization{
			NodeName:   resource.data.GetNodeName(),
			NodeId:     resource.data.GetNodeId(),
			IdentityId: resource.data.GetIdentityId(),
		},
		PowerId: resource.data.DataId,
	}
	return request
}

func NewSyncPowerRequest(resource *LocalResource) *api.SyncPowerRequest {
	return &api.SyncPowerRequest{
		Power: &libTypes.Power{
			JobNodeId: resource.data.JobNodeId,
			PowerId:   resource.data.DataId,
			Information: &libTypes.ResourceUsageOverview{
				TotalMem:       resource.data.TotalMem,
				TotalProcessor: uint32(resource.data.TotalProcessor),
				TotalBandwidth: resource.data.TotalBandwidth,
				UsedMem:        resource.data.UsedMem,
				UsedProcessor:  uint32(resource.data.UsedProcessor),
				UsedBandwidth:  resource.data.UsedBandwidth,
			},
			State: resource.data.State,
		},
	}
}

func NewSaveIdentityRequest(identity *Identity) *api.SaveIdentityRequest {
	request := &api.SaveIdentityRequest{
		Member: &apipb.Organization{
			NodeName:   identity.data.GetNodeName(),
			NodeId:     identity.data.GetNodeId(),
			IdentityId: identity.data.GetIdentityId(),
		},
		Credential: identity.data.GetCredential(),
	}
	return request
}

func NewTaskDetail(task *Task) *libTypes.TaskDetail {
	request := &libTypes.TaskDetail{
		TaskId:        task.data.GetTaskId(),
		TaskName:      task.data.GetTaskName(),
		AlgoSupplier:  task.data.GetAlgoSupplier(),
		DataSuppliers:  task.data.GetDataSuppliers(),
		PowerSuppliers: task.data.GetPowerSuppliers(),
		Receivers:     task.data.GetReceivers(),
		CreateAt:      task.data.GetCreateAt(),
		StartAt:       task.data.GetStartAt(),
		EndAt:         task.data.GetEndAt(),
		State:         task.data.GetState(),
		OperationCost: task.data.GetOperationCost(),
		TaskEvents: task.data.GetTaskEvents(),
	}
	return request
}

func NewMetadataArrayFromResponse(response *api.MetaDataSummaryListResponse) MetadataArray {
	var metadataArray MetadataArray
	for _, v := range response.GetMetadataSummaries() {
		metadata := NewMetadata(&libTypes.MetadataPB{
			IdentityId: v.GetOwner().GetIdentityId(),
			NodeId:     v.GetOwner().GetNodeId(),
			NodeName:   v.GetOwner().GetNodeName(),
			DataId:     v.GetInformation().GetMetaDataId(),
			DataStatus: apipb.DataStatus_DataStatus_Normal,
			OriginId:   v.GetInformation().GetOriginId(),
			TableName:  v.GetInformation().GetTableName(),
			FilePath:   v.GetInformation().GetFilePath(),
			Desc:       v.GetInformation().GetDesc(),
			Rows:       v.GetInformation().GetRows(),
			Columns:    v.GetInformation().GetColumns(),
			Size_:      uint64(v.GetInformation().GetSize_()),
			FileType:   v.GetInformation().GetFileType(),
			State:      v.GetInformation().GetState(),
			HasTitle:   v.GetInformation().GetHasTitle(),
			MetadataColumns: make([]*libTypes.MetadataColumn, 0),
		})
		metadataArray = append(metadataArray, metadata)
	}
	return metadataArray
}

func NewMetadataArrayFromDetailListResponse(response *api.MetadataListResponse) MetadataArray {
	var metadataArray MetadataArray
	for _, v := range response.GetMetadatas() {
		data := &libTypes.MetadataPB{
			IdentityId: v.GetOwner().GetIdentityId(),
			NodeId:     v.GetOwner().GetNodeId(),
			NodeName:   v.GetOwner().GetNodeName(),
			DataId:     v.GetMetaSummary().GetMetaDataId(),
			DataStatus: apipb.DataStatus_DataStatus_Normal,
			OriginId:   v.GetMetaSummary().GetOriginId(),
			TableName:  v.GetMetaSummary().GetTableName(),
			FilePath:   v.GetMetaSummary().GetFilePath(),
			Desc:       v.GetMetaSummary().GetDesc(),
			Rows:       v.GetMetaSummary().GetRows(),
			Columns:    v.GetMetaSummary().GetColumns(),
			Size_:      uint64(v.GetMetaSummary().GetSize_()),
			FileType:   v.GetMetaSummary().GetFileType(),
			State:      v.GetMetaSummary().GetState(),
			HasTitle:   v.GetMetaSummary().GetHasTitle(),
			MetadataColumns: v.GetMetadataColumns(),
		}
		metadata := NewMetadata(data)
		metadataArray = append(metadataArray, metadata)
	}
	return metadataArray
}

//func NewResourceArrayFromPowerListResponse(response *api.PowerTotalSummaryListResponse) ResourceArray {
//	return nil
//}

func NewResourceArrayFromPowerTotalSummaryListResponse(response *api.PowerTotalSummaryListResponse) ResourceArray {
	resourceArray := make(ResourceArray, 0, len(response.GetPowers()))
	for _, v := range response.GetPowers() {
		resource := NewResource(&libTypes.ResourcePB{
			IdentityId:     v.GetOwner().GetIdentityId(),
			NodeId:         v.GetOwner().GetNodeId(),
			NodeName:       v.GetOwner().GetNodeName(),
			DataId:         "", // todo: to be determined
			DataStatus:     apipb.DataStatus_DataStatus_Normal,
			State:          v.GetPower().GetState(),
			TotalMem:       v.GetPower().GetInformation().GetTotalMem(),
			TotalProcessor: v.GetPower().GetInformation().GetTotalProcessor(),
			TotalBandwidth: v.GetPower().GetInformation().GetTotalBandwidth(),
			UsedMem:        v.GetPower().GetInformation().GetUsedMem(),
			UsedProcessor:  v.GetPower().GetInformation().GetUsedProcessor(),
			UsedBandwidth:  v.GetPower().GetInformation().GetUsedBandwidth(),
		})
		resourceArray = append(resourceArray, resource)
	}
	return resourceArray
}

func NewResourceFromResponse(response *api.PowerTotalSummaryResponse) ResourceArray {
	resourceArray := make(ResourceArray, 0)
	resource := NewResource(&libTypes.ResourcePB{
		IdentityId:     response.GetOwner().GetIdentityId(),
		NodeId:         response.GetOwner().GetNodeId(),
		NodeName:       response.GetOwner().GetNodeName(),
		DataId:         "", // todo: to be determined
		DataStatus:     apipb.DataStatus_DataStatus_Normal,
		State:          response.GetPower().GetState(),
		TotalMem:       response.GetPower().GetInformation().GetTotalMem(),
		TotalProcessor: response.GetPower().GetInformation().GetTotalProcessor(),
		TotalBandwidth: response.GetPower().GetInformation().GetTotalBandwidth(),
		UsedMem:        response.GetPower().GetInformation().GetUsedMem(),
		UsedProcessor:  response.GetPower().GetInformation().GetUsedProcessor(),
		UsedBandwidth:  response.GetPower().GetInformation().GetUsedBandwidth(),
	})
	resourceArray = append(resourceArray, resource)
	return resourceArray
}

func NewTaskArrayFromResponse(response *api.TaskListResponse) TaskDataArray {
	taskArray := make(TaskDataArray, 0, len(response.GetTaskDetails()))
	for _, v := range response.GetTaskDetails() {
		task := NewTask(&libTypes.TaskPB{
			// TODO: 任务的所有者标识明确
			IdentityId:    v.GetSender().GetIdentityId(),
			NodeId:        v.GetSender().GetNodeId(),
			NodeName:      v.GetSender().GetNodeName(),
			DataId:        v.GetTaskId(),
			DataStatus:    apipb.DataStatus_DataStatus_Normal,
			TaskId:        v.GetTaskId(),
			TaskName:      v.GetTaskName(),
			State:         v.GetState(),
			Desc:          v.GetDesc(),
			CreateAt:      v.GetCreateAt(),
			StartAt:       v.GetStartAt(),
			EndAt:         v.GetEndAt(),
			AlgoSupplier:  v.GetAlgoSupplier(),
			OperationCost: v.GetOperationCost(),
			DataSuppliers:  v.GetDataSuppliers(),
			PowerSuppliers: v.GetPowerSuppliers(),
			Receivers:     v.GetReceivers(),
			TaskEvents: nil,
		})
		taskArray = append(taskArray, task)
	}
	return taskArray
}

func NewMetadataFromResponse(response *api.MetadataByIdResponse) *Metadata {
	if response == nil {
		return nil
	}
	metadataSummary := response.GetMetadata().GetMetaSummary()
	if metadataSummary == nil {
		return nil
	}
	metadata := &libTypes.MetadataPB{
		IdentityId:         response.GetMetadata().GetOwner().GetIdentityId(),
		NodeId:             response.GetMetadata().GetOwner().GetNodeId(),
		NodeName:           response.GetMetadata().GetOwner().GetNodeName(),
		DataId:             metadataSummary.GetMetaDataId(),
		DataStatus:         apipb.DataStatus_DataStatus_Normal,
		OriginId:           metadataSummary.GetOriginId(),
		TableName:          metadataSummary.GetTableName(),
		FilePath:           metadataSummary.GetFilePath(),
		Desc:               metadataSummary.GetDesc(),
		Rows:               metadataSummary.GetRows(),
		Columns:            metadataSummary.GetColumns(),
		Size_:              uint64(metadataSummary.GetSize_()),
		FileType:           metadataSummary.GetFileType(),
		State:              metadataSummary.GetState(),
		HasTitle:           metadataSummary.GetHasTitle(),
		MetadataColumns: response.GetMetadata().GetMetadataColumns(),
	}
	return NewMetadata(metadata)
}

func NewIdentityArrayFromIdentityListResponse(response *api.IdentityListResponse) IdentityArray {
	if response == nil {
		return nil
	}
	var result IdentityArray
	for _, organization := range response.GetIdentities() {
		result = append(result, NewIdentity(&libTypes.IdentityPB{
			IdentityId: organization.GetIdentityId(),
			NodeId:     organization.GetNodeId(),
			NodeName:   organization.GetNodeName(),
			DataId:     organization.GetIdentityId(),
			DataStatus: apipb.DataStatus_DataStatus_Normal,
		}))
	}
	// todo: need more fields
	return result
}

func NewMetadataAuthArrayFromResponse(responseList []*api.MetaDataAuthorityDetail) MetadataAuthArray {
	if responseList == nil {
		return nil
	}
	var result MetadataAuthArray
	for _, auth := range responseList {
		result = append(result, NewMedataAuth(&libTypes.AuthRecordPB{
			AuthRecordId:         auth.MetaDataAuthId,
			User:                 auth.User,
			UserType:             auth.UserType,
			DataRecord:           auth.Auth,
			AuditResult:          auth.Audit,
			AuditSuggest:         "",
		}))
	}
	return result
}