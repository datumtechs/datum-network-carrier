package task

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
	"strings"
)

func (svr *TaskServiceServer) GetTaskDetailList(ctx context.Context, req *pb.EmptyGetParams) (*pb.GetTaskDetailListResponse, error) {
	tasks, err := svr.B.GetTaskDetailList()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetTaskDetailList failed")
		return nil, ErrGetNodeTaskList
	}
	arr := make([]*pb.GetTaskDetailResponse, len(tasks))
	for i, task := range tasks {
		t := &pb.GetTaskDetailResponse{
			Information: types.ConvertTaskDetailShowToPB(task),
			Role: task.Role,
		}
		arr[i] = t
	}
	log.Debugf("RPC-API:GetTaskDetailList succeed, taskList len: {%d}, json: %s", len(arr), utilTaskDetailResponseArrString(arr))
	return &pb.GetTaskDetailListResponse{
		Status:   0,
		Msg:      backend.OK,
		TaskList: arr,
	}, nil
}

func (svr *TaskServiceServer) GetTaskEventList(ctx context.Context, req *pb.GetTaskEventListRequest) (*pb.GetTaskEventListResponse, error) {

	events, err := svr.B.GetTaskEventList(req.TaskId)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:GetTaskEventList failed, taskId: {%s}", req.TaskId)
		return nil, ErrGetNodeTaskEventList
	}
	log.Debugf("RPC-API:GetTaskEventList succeed, taskId: {%s},  eventList len: {%d}", req.TaskId, len(events))
	return &pb.GetTaskEventListResponse{
		Status:        0,
		Msg:           backend.OK,
		TaskEventList: types.ConvertTaskEventArrToPB(events),
	}, nil
}


func (svr *TaskServiceServer) GetTaskEventListByTaskIds (ctx context.Context, req *pb.GetTaskEventListByTaskIdsRequest) (*pb.GetTaskEventListResponse, error) {

	events, err := svr.B.GetTaskEventListByTaskIds(req.TaskIds)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:GetTaskEventListByTaskIds failed, taskId: {%v}", req.TaskIds)
		return nil, ErrGetNodeTaskEventList
	}
	log.Debugf("RPC-API:GetTaskEventListByTaskIds succeed, taskId: {%v},  eventList len: {%d}", req.TaskIds, len(events))
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
	if len(req.Receivers) == 0 {
		return nil, errors.New("required receivers")
	}
	if len( req.DataSupplier) == 0 {
		return nil, errors.New("required partners")
	}
	if "" == req.CalculateContractcode {
		return nil, errors.New("required CalculateContractCode")
	}

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishTaskDeclare failed, query local identity failed, can not publish task")
		return nil, ErrSendTaskMsg
	}

	taskMsg := types.NewTaskMessageFromRequest(req)

	// add  dataSuppliers
	dataSuppliers := make([]*libTypes.TaskMetadataSupplierData, len(req.DataSupplier))
	for i, v := range req.DataSupplier {

		metaData, err := svr.B.GetMetaDataDetail(v.MemberInfo.IdentityId, v.MetaDataInfo.MetaDataId)
		if nil != err {
			log.WithError(err).Errorf("RPC-API:PublishTaskDeclare failed, query metadata of partner failed, identityId: {%s}, metadataId: {%s}",
				v.MemberInfo.IdentityId, v.MetaDataInfo.MetaDataId)
			return nil, fmt.Errorf("failed to query metadata of partner, identityId: {%s}, metadataId: {%s}",
				v.MemberInfo.IdentityId, v.MetaDataInfo.MetaDataId)
		}

		colTmp := make(map[uint32]*libTypes.ColumnMeta, len(metaData.MetaData.ColumnMetas))
		for _, col := range metaData.MetaData.ColumnMetas {
			colTmp[col.Cindex] = col
		}

		columnArr := make([]*libTypes.ColumnMeta, len(v.MetaDataInfo.ColumnIndexList))
		for j, colIndex := range v.MetaDataInfo.ColumnIndexList {
			if col, ok := colTmp[uint32(colIndex)]; ok {
				columnArr[j] = &libTypes.ColumnMeta{
					Cindex:   col.Cindex,
					Ctype:    col.Ctype,
					Cname:    col.Cname,
					Csize:    col.Csize,
					Ccomment: col.Ccomment,
				}
			} else {
				return nil, fmt.Errorf("not found column of metadata, identityId: {%s}, metadataId: {%s}, columnIndex: {%d}",
					v.MemberInfo.IdentityId, v.MetaDataInfo.MetaDataId, colIndex)
			}
		}

		dataSuppliers[i] = &libTypes.TaskMetadataSupplierData{
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

	//// TODO mock dataSuppliers
	//dataSuppliers := make([]*libTypes.TaskMetadataSupplierData, 0)

	taskMsg.Data.SetMetadataSupplierArr(dataSuppliers)

	// add receivers
	receivers := make([]*libTypes.TaskResultReceiverData, len(req.Receivers))
	for i, v := range req.Receivers {

		providers := make([]*libTypes.OrganizationData, len(v.Providers))
		for j, val := range v.Providers {
			providers[j] = &libTypes.OrganizationData{
				PartyId:  val.PartyId,
				NodeName: val.Name,
				NodeId:   val.NodeId,
				Identity: val.IdentityId,
			}
		}

		receivers[i] = &libTypes.TaskResultReceiverData{
			Receiver: &libTypes.OrganizationData{
				PartyId:  v.MemberInfo.PartyId,
				NodeName: v.MemberInfo.Name,
				NodeId:   v.MemberInfo.NodeId,
				Identity: v.MemberInfo.IdentityId,
			},
			Provider: providers,
		}
	}
	taskMsg.Data.SetReceivers(receivers)

	// add empty powerSuppliers
	taskMsg.Data.TaskData().ResourceSupplier = make([]*libTypes.TaskResourceSupplierData, 0)

	// add taskId
	taskId := taskMsg.SetTaskId()
	taskMsg.Data.TaskData().TaskId = taskId

	err = svr.B.SendMsg(taskMsg)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishTaskDeclare failed, query metadata of partner failed, taskId: {%s}",
			taskId)
		return nil, ErrSendTaskMsg
	}
	log.Debugf("RPC-API:PublishTaskDeclare succeed, taskId: {%s}, taskMsg: %s", taskId, taskMsg.String())
	return &pb.PublishTaskDeclareResponse{
		Status: 0,
		Msg:    backend.OK,
		TaskId: taskId,
	}, nil
}


func utilTaskDetailResponseArrString(tasks []*pb.GetTaskDetailResponse) string {
	arr := make([]string, len(tasks))
	for i, t := range tasks {
		arr[i] = t.String()
	}
	if len(arr) != 0 {
		return "[" + strings.Join(arr, ",") + "]"
	}
	return "[]"
}