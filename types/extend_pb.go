package types

import (
	"github.com/RosettaFlow/Carrier-Go/lib/center/api"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
)

// NewMetaDataSaveRequest converts Metadata object to MetaDataSaveRequest object.
func NewMetaDataSaveRequest(metadata *Metadata) *api.MetaDataSaveRequest {
	request := &api.MetaDataSaveRequest{
		MetaSummary:          &api.MetaDataSummary{
			MetaDataId:           metadata.data.DataId,
			OriginId:             metadata.data.OriginId,
			TableName:            metadata.data.TableName,
			Desc:                 metadata.data.Desc,
			FilePath:             metadata.data.FilePath,
			Rows:                 uint32(metadata.data.Rows),
			Columns:              uint32(metadata.data.Columns),
			Size_:                metadata.data.Size_,
			FileType:             metadata.data.FileType,
			HasTitle:             metadata.data.HasTitleRow,
			State:                metadata.data.State,
		},
		ColumnMeta:         make([]*api.MetaDataColumnDetail, len(metadata.data.ColumnMetaList)),
	}
	for _, column := range metadata.data.ColumnMetaList {
		request.ColumnMeta = append(request.ColumnMeta, &api.MetaDataColumnDetail{
			Cindex:               column.GetCindex(),
			Cname:                column.GetCname(),
			Ctype:                column.GetCtype(),
			Csize:                column.GetCsize(),
			Ccomment:             column.GetCcomment(),
		})
	}
	return request
}

func NewPublishPowerRequest(resource *Resource) *api.PublishPowerRequest {
	request := &api.PublishPowerRequest{
		Owner:                &api.Organization{
			Name:                 resource.data.GetNodeName(),
			NodeId:               resource.data.GetNodeId(),
			IdentityId:           resource.data.GetIdentity(),
		},
		PowerId:          	  resource.data.DataId,
		Information:          &api.PurePower{
			Mem:                  resource.data.GetTotalMem(),
			Processor:            resource.data.GetTotalProcessor(),
			Bandwidth:            resource.data.GetTotalBandWidth(),
		},
	}
	return request
}

func NewSaveIdentityRequest(identity *Identity) *api.SaveIdentityRequest {
	request := &api.SaveIdentityRequest{
		Member:               &api.Organization{
			Name:                 identity.data.GetNodeName(),
			NodeId:               identity.data.GetNodeId(),
			IdentityId:           identity.data.GetIdentity(),
		},
		Credential:           identity.data.GetCredential(),
	}
	return request
}

func NewTaskDetail(task *Task) *api.TaskDetail {
	request := &api.TaskDetail{
		TaskId:               task.data.GetTaskId(),
		TaskName:             task.data.GetTaskName(),
		Owner:                &api.Organization{
			Name:                 task.data.GetNodeName(),
			NodeId:               task.data.GetNodeId(),
			IdentityId:           task.data.GetIdentity(),
		},
		AlgoSupplier:         &api.Organization{
			Name:                 task.data.AlgoSupplier.GetNodeName(),
			NodeId:               task.data.AlgoSupplier.GetNodeId(),
			IdentityId:           task.data.AlgoSupplier.GetIdentity(),
		},
		DataSupplier:         make([]*api.TaskDataSupplier, len(task.data.GetMetadataSupplier())),
		PowerSupplier:        make([]*api.TaskPowerSupplier, len(task.data.GetResourceSupplier())),
		Receivers:            make([]*api.TaskResultReceiver, len(task.data.GetReceivers())),
		CreateAt:             task.data.GetCreateAt(),
		EndAt:                task.data.GetEndAt(),
		State:                task.data.GetState(),
		OperationCost:        &api.TaskOperationCostDeclare{
			CostMem:              task.data.GetTaskResource().GetCostMem(),
			CostProcessor:        task.data.GetTaskResource().GetCostProcessor(),
			CostBandwidth:        task.data.GetTaskResource().GetCostBandwidth(),
			Duration:             task.data.GetTaskResource().GetDuration(),
		},
	}
	for _, v := range task.data.GetMetadataSupplier() {
		dataSupplier :=  &api.TaskDataSupplier{
			MemberInfo:           &api.Organization{
				Name:                 v.GetOrganization().GetNodeName(),
				NodeId:               v.GetOrganization().GetNodeId(),
				IdentityId:           v.GetOrganization().GetIdentity(),
			},
			MetaId:               v.GetMetaId(),
			MetaName:             v.GetMetaName(),
			ColumnMeta:			 make([]*api.MetaDataColumnDetail, len(v.GetColumnList())),
		}
		for _, column := range v.GetColumnList() {
			dataSupplier.ColumnMeta = append(dataSupplier.ColumnMeta, &api.MetaDataColumnDetail{
				Cindex:               column.GetCindex(),
				Cname:                column.GetCname(),
				Ctype:                column.GetCtype(),
				Csize:                column.GetCsize(),
				Ccomment:             column.GetCcomment(),
			})
		}
		request.DataSupplier = append(request.DataSupplier, dataSupplier)
	}
	for _, v := range task.data.GetResourceSupplier() {
		request.PowerSupplier = append(request.PowerSupplier, &api.TaskPowerSupplier{
			MemberInfo:           &api.Organization{
				Name:                 v.GetOrganization().GetNodeName(),
				NodeId:               v.GetOrganization().GetNodeId(),
				IdentityId:           v.GetOrganization().GetIdentity(),
			},
			PowerInfo:            &api.ResourceUsedDetail{
				TotalMem:             v.GetResourceUsedOverview().GetTotalMem(),
				UsedMem:              v.GetResourceUsedOverview().GetUsedMem(),
				TotalProcessor:       v.GetResourceUsedOverview().GetTotalProcessor(),
				UsedProcessor:        v.GetResourceUsedOverview().GetUsedProcessor(),
				TotalBandwidth:       v.GetResourceUsedOverview().GetTotalBandwidth(),
				UsedBandwidth:        v.GetResourceUsedOverview().GetUsedBandwidth(),
			},
		})
	}
	for _, v := range task.data.GetReceivers() {
		receive := v.GetReceiver()
		providers := v.GetProvider()
		taskResultReceiver := &api.TaskResultReceiver{
			MemberInfo:           &api.Organization{
				Name:                 receive.GetNodeName(),
				NodeId:               receive.GetNodeId(),
				IdentityId:           receive.GetIdentity(),
			},
			Provider:             make([]*api.Organization, len(providers)),
		}
		for _, provider := range providers {
			taskResultReceiver.Provider = append(taskResultReceiver.Provider, &api.Organization{
				Name:                 provider.GetNodeName(),
				NodeId:               provider.GetNodeId(),
				IdentityId:           provider.GetIdentity(),
			})
		}
		request.Receivers = append(request.Receivers, taskResultReceiver)
	}
	return request
}

func NewMetadataArrayFromResponse(response *api.MetaDataSummaryListResponse) MetadataArray {
	var metadataArray MetadataArray
	for _, v := range response.GetMetadataSummaryList() {
		metadata := NewMetadata(&libTypes.MetaData{
			Identity:             v.GetOwner().GetIdentityId(),
			NodeId:               v.GetOwner().GetNodeId(),
			NodeName:  			  v.GetOwner().GetName(),
			DataId:               v.GetInformation().GetMetaDataId(),
			DataStatus:           "N",	// todo: 待定
			OriginId:             v.GetInformation().GetOriginId(),
			TableName:            v.GetInformation().GetTableName(),
			FilePath:             v.GetInformation().GetFilePath(),
			Desc:                 v.GetInformation().GetDesc(),
			Rows:                 uint64(v.GetInformation().GetRows()),
			Columns:              uint64(v.GetInformation().GetColumns()),
			Size_:                v.GetInformation().GetSize_(),
			FileType:             v.GetInformation().GetFileType(),
			State:                v.GetInformation().GetState(),
			HasTitleRow:          v.GetInformation().GetHasTitle(),
			ColumnMetaList:       make([]*libTypes.ColumnMeta, 0),
		})
		metadataArray = append(metadataArray, metadata)
	}
	return metadataArray
}

func NewResourceArrayFromPowerListResponse(response *api.PowerListResponse) ResourceArray {
	return nil
}

func NewResourceArrayFromPowerTotalSummaryListResponse(response *api.PowerTotalSummaryListResponse) ResourceArray {
	resourceArray := make(ResourceArray, len(response.GetPowerList()))
	for _, v := range response.GetPowerList() {
		resource := NewResource(&libTypes.ResourceData{
			Identity:   v.GetOwner().GetIdentityId(),
			NodeId:     v.GetOwner().GetNodeId(),
			NodeName:   v.GetOwner().GetName(),
			DataId:     "", // todo: to be determined
			DataStatus: "", // todo: to be determined
			State:      v.GetPower().GetState(),
			TotalMem:   v.GetPower().GetInformation().GetTotalMem(),
			UsedMem:    v.GetPower().GetInformation().GetUsedMem(),
			TotalProcessor: v.GetPower().GetInformation().GetTotalProcessor(),
			TotalBandWidth:       v.GetPower().GetInformation().GetTotalBandwidth(),
		})
		resourceArray = append(resourceArray, resource)
	}
	return resourceArray
}

func NewResourceFromResponse(response *api.PowerTotalSummaryResponse) ResourceArray {
	resourceArray := make(ResourceArray, 1)
	resource := NewResource(&libTypes.ResourceData{
		Identity:   response.GetOwner().GetIdentityId(),
		NodeId:     response.GetOwner().GetNodeId(),
		NodeName:   response.GetOwner().GetName(),
		DataId:     "", // todo: to be determined
		DataStatus: "", // todo: to be determined
		State:      response.GetPower().GetState(),
		TotalMem:   response.GetPower().GetInformation().GetTotalMem(),
		UsedMem:    response.GetPower().GetInformation().GetUsedMem(),
		TotalProcessor: response.GetPower().GetInformation().GetTotalProcessor(),
		TotalBandWidth:       response.GetPower().GetInformation().GetTotalBandwidth(),
	})
	resourceArray = append(resourceArray, resource)
	return resourceArray
}

func NewTaskArrayFromResponse(response *api.TaskListResponse) TaskDataArray {
	taskArray := make(TaskDataArray, len(response.GetTaskList()))
	for _, v := range response.GetTaskList() {
		task := NewTask(&libTypes.TaskData{
			Identity:             v.GetOwner().GetIdentityId(),
			NodeId:               v.GetOwner().GetNodeId(),
			NodeName:             v.GetOwner().GetName(),
			DataId:               v.GetTaskId(), // todo: to be determined
			DataStatus:           "Y",           // todo: to be determined
			TaskId:               v.GetTaskId(),
			TaskName:             v.GetTaskName(),
			State:                v.GetState(),
			Desc:                 v.GetDesc(),
			CreateAt:             v.GetCreateAt(),
			EndAt:                v.GetEndAt(),
			AlgoSupplier:         &libTypes.OrganizationData{
				Alias:                v.GetAlgoSupplier().GetName(),
				Identity:             v.GetAlgoSupplier().GetIdentityId(),
				NodeId:               v.GetAlgoSupplier().GetNodeId(),
				NodeName:             v.GetAlgoSupplier().GetName(),
			},
			TaskResource:         &libTypes.TaskResourceData{
				CostMem:              v.GetOperationCost().GetCostMem(),
				CostProcessor:        v.GetOperationCost().GetCostProcessor(),
				CostBandwidth:        v.GetOperationCost().GetCostBandwidth(),
				Duration:             v.GetOperationCost().GetDuration(),
			},
			MetadataSupplier:     make([]*libTypes.TaskMetadataSupplierData, len(v.GetDataSupplier())),
			ResourceSupplier:     make([]*libTypes.TaskResourceSupplierData, len(v.GetPowerSupplier())),
			Receivers:            make([]*libTypes.TaskResultReceiverData, len(v.GetReceivers())),
			PartnerList:          make([]*libTypes.OrganizationData, len(v.GetDataSupplier())),
			EventDataList:        nil,
		})

		// MetadataSupplier filling
		for _, supplier := range v.GetDataSupplier() {
			// partner == dataSupplier
			task.data.PartnerList = append(task.data.PartnerList, &libTypes.OrganizationData{
				Identity:             supplier.GetMemberInfo().GetIdentityId(),
				NodeId:               supplier.GetMemberInfo().GetNodeId(),
				NodeName:             supplier.GetMemberInfo().GetName(),
			})
			supplierData := &libTypes.TaskMetadataSupplierData{
				Organization:         &libTypes.OrganizationData{
					Alias:                supplier.GetMemberInfo().GetName(),
					Identity:             supplier.GetMemberInfo().GetIdentityId(),
					NodeId:               supplier.GetMemberInfo().GetNodeId(),
					NodeName:             supplier.GetMemberInfo().GetName(),
				},
				MetaId:               supplier.GetMetaId(),
				MetaName:             supplier.GetMetaName(),
				ColumnList:           make([]*libTypes.ColumnMeta, len(supplier.GetColumnMeta())),
			}
			for _, columnMeta := range supplier.GetColumnMeta() {
				supplierData.ColumnList = append(supplierData.ColumnList, &libTypes.ColumnMeta{
					Cindex:               columnMeta.GetCindex(),
					Cname:                columnMeta.GetCname(),
					Ctype:                columnMeta.GetCtype(),
					Csize:                columnMeta.GetCsize(),
					Ccomment:             columnMeta.GetCcomment(),
				})
			}
			task.data.MetadataSupplier = append(task.data.MetadataSupplier, supplierData)
		}
		// ResourceSupplier
		for _, power := range v.GetPowerSupplier() {
			supplierData := &libTypes.TaskResourceSupplierData{
				Organization:         &libTypes.OrganizationData{
					Alias:                power.GetMemberInfo().GetName(),
					Identity:             power.GetMemberInfo().GetIdentityId(),
					NodeId:               power.GetMemberInfo().GetNodeId(),
					NodeName:             power.GetMemberInfo().GetName(),
				},
				ResourceUsedOverview: &libTypes.ResourceUsedOverview{
					TotalMem:             power.GetPowerInfo().GetTotalMem(),
					UsedMem:              power.GetPowerInfo().GetUsedMem(),
					TotalProcessor:       power.GetPowerInfo().GetTotalProcessor(),
					UsedProcessor:        power.GetPowerInfo().GetUsedProcessor(),
					TotalBandwidth:       power.GetPowerInfo().GetTotalBandwidth(),
					UsedBandwidth:        power.GetPowerInfo().GetUsedBandwidth(),
				},
			}
			task.data.ResourceSupplier = append(task.data.ResourceSupplier, supplierData)
		}
		// Receivers filling.
		for _, receiver := range v.GetReceivers() {
			receiverData := &libTypes.TaskResultReceiverData{
				Receiver:             &libTypes.OrganizationData{
					Alias:                receiver.GetMemberInfo().GetName(),
					Identity:             receiver.GetMemberInfo().GetIdentityId(),
					NodeId:               receiver.GetMemberInfo().GetNodeId(),
					NodeName:             receiver.GetMemberInfo().GetName(),
				},
				Provider:             make([]*libTypes.OrganizationData, len(receiver.GetProvider())),
			}
			for _, provider := range receiver.GetProvider() {
				receiverData.Provider = append(receiverData.Provider, &libTypes.OrganizationData{
					Alias:                provider.GetName(),
					Identity:             provider.GetIdentityId(),
					NodeId:               provider.GetNodeId(),
					NodeName:             provider.GetName(),
				})
			}
			task.data.Receivers = append(task.data.Receivers, receiverData)
		}
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
	metadata := &libTypes.MetaData{
		Identity:             response.GetMetadata().GetOwner().GetIdentityId(),
		NodeId:               response.GetMetadata().GetOwner().GetNodeId(),
		NodeName:             response.GetMetadata().GetOwner().GetName(),
		DataId:               metadataSummary.GetMetaDataId(),
		DataStatus:           "Y",
		OriginId:             metadataSummary.GetOriginId(),
		TableName:            metadataSummary.GetTableName(),
		FilePath:             metadataSummary.GetFilePath(),
		Desc:                 metadataSummary.GetDesc(),
		Rows:                 uint64(metadataSummary.GetRows()),
		Columns:              uint64(metadataSummary.GetColumns()),
		Size_:                metadataSummary.GetSize_(),
		FileType:             metadataSummary.GetFileType(),
		State:                metadataSummary.GetState(),
		HasTitleRow:          metadataSummary.GetHasTitle(),
		ColumnMetaList:       make([]*libTypes.ColumnMeta, len(response.GetMetadata().GetColumnMeta())),
	}
	for _, v := range response.GetMetadata().GetColumnMeta() {
		metadata.ColumnMetaList = append(metadata.ColumnMetaList, &libTypes.ColumnMeta{
			Cindex:               v.GetCindex(),
			Cname:                v.GetCname(),
			Ctype:                v.GetCtype(),
			Csize:                v.GetCsize(),
			Ccomment:             v.GetCcomment(),
		})
	}
	return NewMetadata(metadata)
}

func NewIdentityArrayFromIdentityListResponse(response *api.IdentityListResponse) IdentityArray {

	return nil
}