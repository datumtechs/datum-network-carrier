package task

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
)

//func (svr *TaskServiceServer) GetTaskSummaryList(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetTaskSummaryListResponse, error) {
//	return nil, nil
//}
//func (svr *TaskServiceServer) GetTaskJoinSummaryList(ctx context.Context, req *pb.GetTaskJoinSummaryListRequest) (*pb.GetTaskJoinSummaryListResponse, error) {
//	return nil, nil
//}
//func (svr *TaskServiceServer) GetTaskDetail(ctx context.Context, req *pb.GetTaskDetailRequest) (*pb.GetTaskDetailResponse, error) {
//	return nil, nil
//}
func (svr *TaskServiceServer) GetTaskDetailList(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetTaskDetailListResponse, error) {
	tasks, err := svr.B.GetTaskDetailList()
	if nil != err {
		return nil, backend.NewRpcBizErr(ErrGetNodeTaskListStr)
	}
	arr := make([]*pb.GetTaskDetailResponse, len(tasks))
	for i, task := range tasks {
		t := &pb.GetTaskDetailResponse{
			Information: types.ConvertTaskDetailShowToPB(task),
		}
		arr[i] = t
	}
	return &pb.GetTaskDetailListResponse{
		Status:   0,
		Msg:      backend.OK,
		TaskList: arr,
	}, nil
}

func (svr *TaskServiceServer) GetTaskEventList(ctx context.Context, req *pb.GetTaskEventListRequest) (*pb.GetTaskEventListResponse, error) {

	events, err := svr.B.GetTaskEventList(req.TaskId)
	if nil != err {
		return nil, backend.NewRpcBizErr(ErrGetNodeTaskEventListStr)
	}

	return &pb.GetTaskEventListResponse{
		Status:        0,
		Msg:           backend.OK,
		TaskEventList: types.ConvertTaskEventArrToPB(events),
	}, nil
}

func (svr *TaskServiceServer) PublishTaskDeclare(ctx context.Context, req *pb.PublishTaskDeclareRequest) (*pb.PublishTaskDeclareResponse, error) {
	if req == nil || req.Owner == nil {
		return nil, errors.New("required owner")
	}
	if req.OperationCost == nil {
		return nil, errors.New("required operationCost")
	}
	if req.Receivers == nil {
		return nil, errors.New("required receivers")
	}
	if req.Partners == nil {
		return nil, errors.New("required partners")
	}
	taskMsg := types.NewTaskMessageFromRequest(req)

	partners := make([]*libTypes.TaskMetadataSupplierData, len(req.Partners))
	for i, v := range req.Partners {

		metaData, err := svr.B.GetMetaDataDetail(v.MemberInfo.IdentityId, v.MetaDataInfo.MetaDataId)
		if nil != err {
			return nil, fmt.Errorf("failed to query metadata of partner, identityId: {%s}, metadataId: {%s}", v.MemberInfo.IdentityId, v.MetaDataInfo.MetaDataId)
		}

		columnArr := make([]*libTypes.ColumnMeta, len(v.MetaDataInfo.ColumnIndexList))
		for j, index := range v.MetaDataInfo.ColumnIndexList {
			col := metaData.MetaData.ColumnMetas[index]
			columnArr[j] = &libTypes.ColumnMeta{
				Cindex:   col.Cindex,
				Ctype:    col.Ctype,
				Cname:    col.Cname,
				Csize:    col.Csize,
				Ccomment: col.Ccomment,
			}
		}

		partners[i] = &libTypes.TaskMetadataSupplierData{
			Organization: &libTypes.OrganizationData{
				PartyId:  v.MemberInfo.PartyId,
				NodeName: v.MemberInfo.Name,
				NodeId:   v.MemberInfo.NodeId,
				Identity: v.MemberInfo.IdentityId,
			},
			MetaId:     v.MetaDataInfo.MetaDataId,
			MetaName:   metaData.MetaData.MetaDataSummary.TableName,
			ColumnList: columnArr,
		}
	}

	taskMsg.Data.SetMetadataSupplierArr(partners)

	receivers := make([]*libTypes.TaskResultReceiverData, len(req.Receivers))
	for i, v := range req.Receivers {

		providers := make([]*libTypes.OrganizationData, len(v.Providers))
		for j, val := range v.Providers {
			provider := &libTypes.OrganizationData{
				PartyId:  val.PartyId,
				NodeName: val.Name,
				NodeId:   val.NodeId,
				Identity: val.IdentityId,
			}
			providers[j] = provider
		}

		receiver := &libTypes.TaskResultReceiverData{

			Receiver: &libTypes.OrganizationData{
				PartyId:  v.MemberInfo.PartyId,
				NodeName: v.MemberInfo.Name,
				NodeId:   v.MemberInfo.NodeId,
				Identity: v.MemberInfo.IdentityId,
			},
			Provider: providers,
		}
		receivers[i] = receiver
	}
	taskMsg.Data.SetReceivers(receivers)
	taskId := taskMsg.GetTaskId()
	taskMsg.Data.TaskData().TaskId = taskId

	err := svr.B.SendMsg(taskMsg)
	if nil != err {
		return nil, backend.NewRpcBizErr(ErrSendTaskMsgStr)
	}
	return &pb.PublishTaskDeclareResponse{
		Status: 0,
		Msg:    backend.OK,
		TaskId: taskId,
	}, nil
}
