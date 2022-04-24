package types

import (
	"github.com/Metisnetwork/Metis-Carrier/lib/center/api"
	libtypes "github.com/Metisnetwork/Metis-Carrier/lib/types"
)

// NewMetadataSaveRequest converts Metadata object to SaveMetadataRequest object.
func NewMetadataSaveRequest(metadata *Metadata) *api.SaveMetadataRequest {
	request := &api.SaveMetadataRequest{
		Metadata: metadata.GetData(),
	}
	return request
}

func NewMetadataUpdateRequest(metadata *Metadata) *api.UpdateMetadataRequest {
	request := &api.UpdateMetadataRequest{
		Metadata: metadata.GetData(),
	}
	return request
}

func NewMetadataRevokeRequest(metadata *Metadata) *api.RevokeMetadataRequest {
	request := &api.RevokeMetadataRequest{
		Owner: &libtypes.Organization{
			IdentityId: metadata.GetData().GetOwner().GetIdentityId(),
			NodeId:     metadata.GetData().GetOwner().GetNodeId(),
			NodeName:   metadata.GetData().GetOwner().GetNodeName(),
		},
		MetadataId: metadata.GetData().GetDataId(),
	}
	return request
}

func NewPublishPowerRequest(resource *Resource) *api.PublishPowerRequest {
	request := &api.PublishPowerRequest{
		Power: resource.data,
	}
	return request
}

func RevokePowerRequest(resource *Resource) *api.RevokePowerRequest {
	request := &api.RevokePowerRequest{
		Owner: &libtypes.Organization{
			NodeName:   resource.GetNodeName(),
			NodeId:     resource.GetNodeId(),
			IdentityId: resource.GetIdentityId(),
		},
		PowerId: resource.GetDataId(),
	}
	return request
}

func NewSyncPowerRequest(resource *LocalResource) *api.SyncPowerRequest {
	return &api.SyncPowerRequest{
		Power: resource.GetData(),
	}
}

func NewSaveIdentityRequest(identity *Identity) *api.SaveIdentityRequest {
	request := &api.SaveIdentityRequest{
		Member: &libtypes.Organization{
			NodeName:   identity.GetName(),
			NodeId:     identity.GetNodeId(),
			IdentityId: identity.GetIdentityId(),
			ImageUrl:   identity.GetImageUrl(),
			Details:    identity.GetDetails(),
		},
		Credential: identity.GetCredential(),
	}
	return request
}

func NewSaveTaskRequest(task *Task) *api.SaveTaskRequest {
	request := &api.SaveTaskRequest{
		Task: task.data,
	}
	return request
}

func NewMetadataArrayFromDetailListResponse(response *api.ListMetadataResponse) MetadataArray {
	var metadataArray MetadataArray
	for _, v := range response.GetMetadata() {
		metadata := NewMetadata(v)
		metadataArray = append(metadataArray, metadata)
	}
	return metadataArray
}

func NewResourceArrayFromPowerTotalSummaryListResponse(response *api.ListPowerSummaryResponse) ResourceArray {
	resourceArray := make(ResourceArray, 0, len(response.GetPowers()))
	for _, v := range response.GetPowers() {
		resource := NewResource(&libtypes.ResourcePB{
			/**
			Owner                *Organization
			DataId               string
			DataStatus           DataStatus
			State                PowerState
			TotalMem             uint64
			UsedMem              uint64
			TotalProcessor       uint32
			UsedProcessor        uint32
			TotalBandwidth       uint64
			UsedBandwidth        uint64
			TotalDisk            uint64
			UsedDisk             uint64
			PublishAt            uint64
			UpdateAt             uint64
			Nonce                uint64
			*/
			Owner:          v.GetOwner(),
			DataId:         "", // todo: to be determined
			DataStatus:     libtypes.DataStatus_DataStatus_Valid,
			State:          v.GetPowerSummary().GetState(),
			TotalMem:       v.GetPowerSummary().GetInformation().GetTotalMem(),
			TotalProcessor: v.GetPowerSummary().GetInformation().GetTotalProcessor(),
			TotalBandwidth: v.GetPowerSummary().GetInformation().GetTotalBandwidth(),
			TotalDisk:      v.GetPowerSummary().GetInformation().GetTotalDisk(),
			UsedMem:        v.GetPowerSummary().GetInformation().GetUsedMem(),
			UsedProcessor:  v.GetPowerSummary().GetInformation().GetUsedProcessor(),
			UsedBandwidth:  v.GetPowerSummary().GetInformation().GetUsedBandwidth(),
			UsedDisk:       v.GetPowerSummary().GetInformation().GetUsedDisk(),
			// todo Summary is aggregate information and does not require paging, so there are no `publishat` and `updateat` and `nonce`
		})
		resourceArray = append(resourceArray, resource)
	}
	return resourceArray
}

func NewResourceArrayFromPowerDetailListResponse(response *api.ListPowerResponse) ResourceArray {
	resourceArray := make(ResourceArray, 0, len(response.GetPowers()))
	for _, v := range response.GetPowers() {
		resource := NewResource(v)
		resourceArray = append(resourceArray, resource)
	}
	return resourceArray
}

func NewResourceFromPowerSummaryResponse(response *api.PowerSummaryResponse) ResourceArray {
	resourceArray := make(ResourceArray, 0)
	resource := NewResource(&libtypes.ResourcePB{
		/**
		// todo summary 不需要加上 nonce 字段
		Owner                *Organization
		DataId               string
		DataStatus           DataStatus
		State                PowerState
		TotalMem             uint64
		UsedMem              uint64
		TotalProcessor       uint32
		UsedProcessor        uint32
		TotalBandwidth       uint64
		UsedBandwidth        uint64
		TotalDisk            uint64
		UsedDisk             uint64
		PublishAt            uint64
		UpdateAt             uint64
		Nonce                uint64
		*/
		Owner:          response.GetOwner(),
		DataId:         "",
		DataStatus:     libtypes.DataStatus_DataStatus_Valid,
		State:          response.GetPowerSummary().GetState(),
		TotalMem:       response.GetPowerSummary().GetInformation().GetTotalMem(),
		TotalProcessor: response.GetPowerSummary().GetInformation().GetTotalProcessor(),
		TotalBandwidth: response.GetPowerSummary().GetInformation().GetTotalBandwidth(),
		TotalDisk:      response.GetPowerSummary().GetInformation().GetTotalDisk(),
		UsedMem:        response.GetPowerSummary().GetInformation().GetUsedMem(),
		UsedProcessor:  response.GetPowerSummary().GetInformation().GetUsedProcessor(),
		UsedBandwidth:  response.GetPowerSummary().GetInformation().GetUsedBandwidth(),
		UsedDisk:       response.GetPowerSummary().GetInformation().GetUsedDisk(),
		//PublishAt:
		//UpdateAt:
		//Nonce:
	})
	resourceArray = append(resourceArray, resource)
	return resourceArray
}

func NewTaskArrayFromResponse(response *api.ListTaskResponse) TaskDataArray {
	taskArray := make(TaskDataArray, 0, len(response.GetTasks()))
	for _, v := range response.GetTasks() {
		taskArray = append(taskArray, NewTask(v))
	}
	return taskArray
}

func NewMetadataFromResponse(response *api.FindMetadataByIdResponse) *Metadata {
	if response == nil {
		return nil
	}
	return NewMetadata(response.GetMetadata())
}

func NewIdentityArrayFromIdentityListResponse(response *api.ListIdentityResponse) IdentityArray {
	if response == nil {
		return nil
	}
	var result IdentityArray
	for _, organization := range response.GetIdentities() {
		result = append(result, NewIdentity(&libtypes.IdentityPB{

			/**
			IdentityId           string
			NodeId               string
			NodeName             string
			DataId               string
			DataStatus           DataStatus
			Status               CommonStatus
			Credential           string
			UpdateAt             uint64
			ImageUrl             string
			Details              string
			Nonce                uint64
			*/
			IdentityId: organization.GetIdentityId(),
			NodeId:     organization.GetNodeId(),
			NodeName:   organization.GetNodeName(),
			DataId:     organization.GetIdentityId(),
			DataStatus: organization.GetDataStatus(),
			Status:     organization.GetStatus(),
			Credential: organization.GetCredential(),
			UpdateAt: organization.GetUpdateAt(),
			ImageUrl:   organization.GetImageUrl(),
			Details:    organization.GetDetails(),
			Nonce:      organization.GetNonce(),
		}))
	}
	return result
}

func NewMetadataAuthArrayFromResponse(responseList []*libtypes.MetadataAuthorityPB) MetadataAuthArray {
	if responseList == nil {
		return nil
	}
	var result MetadataAuthArray
	for _, auth := range responseList {
		result = append(result, NewMetadataAuthority(auth))
	}
	return result
}
