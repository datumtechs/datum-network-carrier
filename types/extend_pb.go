package types

import "github.com/RosettaFlow/Carrier-Go/lib/center/api"

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
		JobNodeId:            "", // 废弃
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
		Receivers:            make([]*api.Organization, len(task.data.GetReceivers())),
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
		request.Receivers = append(request.Receivers, &api.Organization{
			Name:                 v.GetNodeName(),
			NodeId:               v.GetNodeId(),
			IdentityId:           v.GetIdentity(),
		})
	}
	return request
}

func NewMetadataArrayFromResponse(response *api.MetaDataSummaryListResponse) MetadataArray {
	return nil
}

func NewResourceArrayFromResponse(response *api.PowerTotalSummaryListResponse) ResourceArray {
	return nil
}

func NewTaskArrayFromResponse(response *api.TaskListResponse) TaskDataArray {
	return nil
}

