package task

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

func (svr *Server) TerminateTask(context.Context, *pb.TerminateTaskRequest) (*apipb.SimpleResponse, error) {
	return nil, nil
}

func (svr *Server) GetTaskDetailList(ctx context.Context, req *emptypb.Empty) (*pb.GetTaskDetailListResponse, error) {
	tasks, err := svr.B.GetTaskDetailList()
	if nil != err {
		log.WithError(err).Error("RPC-API:GetTaskDetailList failed")
		return nil, ErrGetNodeTaskList
	}
	arr := make([]*pb.GetTaskDetailResponse, len(tasks))
	for i, task := range tasks {
		t := &pb.GetTaskDetailResponse{
			Information: task,
			//TODO: 待确认
			//Role: task.Role,
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

func (svr *Server) GetTaskEventList(ctx context.Context, req *pb.GetTaskEventListRequest) (*pb.GetTaskEventListResponse, error) {

	events, err := svr.B.GetTaskEventList(req.TaskId)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:GetTaskEventList failed, taskId: {%s}", req.TaskId)
		return nil, ErrGetNodeTaskEventList
	}
	log.Debugf("RPC-API:GetTaskEventList succeed, taskId: {%s},  eventList len: {%d}", req.TaskId, len(events))
	return &pb.GetTaskEventListResponse{
		Status:        0,
		Msg:           backend.OK,
		TaskEventList: events,
	}, nil
}


func (svr *Server) GetTaskEventListByTaskIds (ctx context.Context, req *pb.GetTaskEventListByTaskIdsRequest) (*pb.GetTaskEventListResponse, error) {

	events, err := svr.B.GetTaskEventListByTaskIds(req.TaskIds)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:GetTaskEventListByTaskIds failed, taskId: {%v}", req.TaskIds)
		return nil, ErrGetNodeTaskEventList
	}
	log.Debugf("RPC-API:GetTaskEventListByTaskIds succeed, taskId: {%v},  eventList len: {%d}", req.TaskIds, len(events))
	return &pb.GetTaskEventListResponse{
		Status:        0,
		Msg:           backend.OK,
		TaskEventList: events,
	}, nil
}

func (svr *Server) PublishTaskDeclare(ctx context.Context, req *pb.PublishTaskDeclareRequest) (*pb.PublishTaskDeclareResponse, error) {
	if req.OperationCost == nil {
		return nil, errors.New("required operationCost")
	}
	if len(req.Receivers) == 0 {
		return nil, errors.New("required receivers")
	}
	if len( req.DataSupplier) == 0 {
		return nil, errors.New("required partners")
	}
	if "" == req.CalculateContractCode {
		return nil, errors.New("required CalculateContractCode")
	}

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishTaskDeclare failed, query local identity failed, can not publish task")
		return nil, ErrSendTaskMsg
	}

	taskMsg := types.NewTaskMessageFromRequest(req)

	// add  dataSuppliers
	dataSuppliers := make([]*libTypes.TaskDataSupplier, len(req.DataSupplier))
	for i, v := range req.DataSupplier {

		metaData, err := svr.B.GetMetaDataDetail(v.MemberInfo.IdentityId, v.MetaDataInfo.MetaDataId)
		if nil != err {
			log.WithError(err).Errorf("RPC-API:PublishTaskDeclare failed, query metadata of partner failed, identityId: {%s}, metadataId: {%s}",
				v.MemberInfo.IdentityId, v.MetaDataInfo.MetaDataId)
			return nil, fmt.Errorf("failed to query metadata of partner, identityId: {%s}, metadataId: {%s}",
				v.MemberInfo.IdentityId, v.MetaDataInfo.MetaDataId)
		}

		colTmp := make(map[uint32]*libTypes.MetadataColumn, len(metaData.Information.MetadataColumns))
		for _, col := range metaData.Information.MetadataColumns {
			colTmp[col.CIndex] = col
		}

		columnArr := make([]*libTypes.MetadataColumn, len(v.MetaDataInfo.ColumnIndexList))
		for j, colIndex := range v.MetaDataInfo.ColumnIndexList {
			if col, ok := colTmp[uint32(colIndex)]; ok {
				columnArr[j] = col
				/*columnArr[j] = &libTypes.MetadataColumn{
					CIndex:   col.CIndex,
					CType:    col.CType,
					CName:    col.CName,
					CSize:    col.CSize,
					CComment: col.CComment,
				}*/
			} else {
				return nil, fmt.Errorf("not found column of metadata, identityId: {%s}, metadataId: {%s}, columnIndex: {%d}",
					v.MemberInfo.IdentityId, v.MetaDataInfo.MetaDataId, colIndex)
			}
		}

		dataSuppliers[i] = &libTypes.TaskDataSupplier{
			MemberInfo: &apipb.TaskOrganization{
				PartyId:  v.MemberInfo.PartyId,
				NodeName: v.MemberInfo.NodeName,
				NodeId:   v.MemberInfo.NodeId,
				IdentityId: v.MemberInfo.IdentityId,
			},
			MetadataId:     v.MetaDataInfo.MetaDataId,
			MetadataName:   metaData.Information.MetaDataSummary.TableName,
			Columns: columnArr,
		}
	}

	//// TODO mock dataSuppliers
	//dataSuppliers := make([]*libTypes.TaskMetadataSupplierData, 0)

	taskMsg.Data.SetMetadataSupplierArr(dataSuppliers)

	// add receivers
	receivers := req.Receivers
	//receivers := make([]*libTypes.TaskResultReceiver, len(req.Receivers))
	//for i, v := range req.Receivers {
	//	providers := make([]*apipb.TaskOrganization, len(v.Providers))
	//	for j, val := range v.Providers {
	//		providers[j] = val
	//		/*providers[j] = &apipb.OrganizationData{
	//			PartyId:  val.PartyId,
	//			NodeName: val.Name,
	//			NodeId:   val.NodeId,
	//			Identity: val.IdentityId,
	//		}*/
	//	}
	//	receivers[i] = &libTypes.TaskResultReceiver{
	//		Receiver: &apipb.TaskOrganization{
	//			PartyId:  v.MemberInfo.PartyId,
	//			NodeName: v.MemberInfo.Name,
	//			NodeId:   v.MemberInfo.NodeId,
	//			Identity: v.MemberInfo.IdentityId,
	//		},
	//		Providers: providers,
	//	}
	//}*/
	taskMsg.Data.SetReceivers(receivers)

	// add empty powerSuppliers
	taskMsg.Data.GetTaskData().PowerSuppliers = make([]*libTypes.TaskPowerSupplier, 0)

	// add taskId
	taskId := taskMsg.SetTaskId()
	taskMsg.Data.GetTaskData().TaskId = taskId

	err = svr.B.SendMsg(taskMsg)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishTaskDeclare failed, query metadata of partner failed, taskId: {%s}",
			taskId)
		return nil, ErrSendTaskMsg
	}
	//log.Debugf("RPC-API:PublishTaskDeclare succeed, taskId: {%s}, taskMsg: %s", taskId, taskMsg.String())
	log.Debugf("RPC-API:PublishTaskDeclare succeed, taskId: {%s}", taskId)
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