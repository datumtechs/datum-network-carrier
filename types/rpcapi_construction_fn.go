package types

import (
	"github.com/RosettaFlow/Carrier-Go/lib/center/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
)

// NewMetadataSaveRequest converts Metadata object to SaveMetadataRequest object.
func NewMetadataSaveRequest(metadata *Metadata) *api.SaveMetadataRequest {
	request := &api.SaveMetadataRequest{
		Metadata: metadata.GetData(),
	}
	return request
}

func NewMetadataRevokeRequest(metadata *Metadata) *api.RevokeMetadataRequest {
	request := &api.RevokeMetadataRequest{
		Owner: &apicommonpb.Organization{
			IdentityId: metadata.GetData().IdentityId,
			NodeId:     metadata.GetData().NodeId,
			NodeName:   metadata.GetData().NodeName,
		},
		MetadataId: metadata.GetData().DataId,
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
		Owner: &apicommonpb.Organization{
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
		Power: resource.GetData(),
	}
}

func NewSaveIdentityRequest(identity *Identity) *api.SaveIdentityRequest {
	request := &api.SaveIdentityRequest{
		Member: &apicommonpb.Organization{
			NodeName:   identity.data.GetNodeName(),
			NodeId:     identity.data.GetNodeId(),
			IdentityId: identity.data.GetIdentityId(),
		},
		Credential: identity.data.GetCredential(),
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
			IdentityId:     v.GetOwner().GetIdentityId(),
			NodeId:         v.GetOwner().GetNodeId(),
			NodeName:       v.GetOwner().GetNodeName(),
			DataId:         "", // todo: to be determined
			DataStatus:     apicommonpb.DataStatus_DataStatus_Normal,
			State:          v.GetPowerSummary().GetState(),
			TotalMem:       v.GetPowerSummary().GetInformation().GetTotalMem(),
			TotalProcessor: v.GetPowerSummary().GetInformation().GetTotalProcessor(),
			TotalBandwidth: v.GetPowerSummary().GetInformation().GetTotalBandwidth(),
			UsedMem:        v.GetPowerSummary().GetInformation().GetUsedMem(),
			UsedProcessor:  v.GetPowerSummary().GetInformation().GetUsedProcessor(),
			UsedBandwidth:  v.GetPowerSummary().GetInformation().GetUsedBandwidth(),
		})
		resourceArray = append(resourceArray, resource)
	}
	return resourceArray
}

func NewResourceFromResponse(response *api.PowerSummaryResponse) ResourceArray {
	resourceArray := make(ResourceArray, 0)
	resource := NewResource(&libtypes.ResourcePB{
		IdentityId:     response.GetOwner().GetIdentityId(),
		NodeId:         response.GetOwner().GetNodeId(),
		NodeName:       response.GetOwner().GetNodeName(),
		DataId:         "", // todo: to be determined
		DataStatus:     apicommonpb.DataStatus_DataStatus_Normal,
		State:          response.GetPowerSummary().GetState(),
		TotalMem:       response.GetPowerSummary().GetInformation().GetTotalMem(),
		TotalProcessor: response.GetPowerSummary().GetInformation().GetTotalProcessor(),
		TotalBandwidth: response.GetPowerSummary().GetInformation().GetTotalBandwidth(),
		UsedMem:        response.GetPowerSummary().GetInformation().GetUsedMem(),
		UsedProcessor:  response.GetPowerSummary().GetInformation().GetUsedProcessor(),
		UsedBandwidth:  response.GetPowerSummary().GetInformation().GetUsedBandwidth(),
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
	return NewMetadata(response.Metadata)
}

func NewIdentityArrayFromIdentityListResponse(response *api.ListIdentityResponse) IdentityArray {
	if response == nil {
		return nil
	}
	var result IdentityArray
	for _, organization := range response.GetIdentities() {
		result = append(result, NewIdentity(&libtypes.IdentityPB{
			IdentityId: organization.GetIdentityId(),
			NodeId:     organization.GetNodeId(),
			NodeName:   organization.GetNodeName(),
			DataId:     organization.GetIdentityId(),
			DataStatus: apicommonpb.DataStatus_DataStatus_Normal,
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