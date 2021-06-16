package types

import "github.com/RosettaFlow/Carrier-Go/lib/center/api"

// NewMetaDataSaveRequest converts Metadata object to MetaDataSaveRequest object.
func NewMetaDataSaveRequest(metadata *Metadata) *api.MetaDataSaveRequest {
	return nil
}

func NewMetadataArrayFromResponse(response *api.MetaDataSummaryListResponse) MetadataArray {
	return nil
}

func NewPublishPowerRequest(resource *Resource) *api.PublishPowerRequest {
	return nil
}

func NewResourceArrayFromResponse(response *api.PowerTotalSummaryListResponse) ResourceArray {
	return nil
}

func NewSaveIdentityRequest(identity *Identity) *api.SaveIdentityRequest {
	return nil
}

func NewTaskDetail(task *Task) *api.TaskDetail {
	return nil
}

func NewTaskArrayFromResponse(response *api.TaskListResponse) TaskDataArray {
	return nil
}
