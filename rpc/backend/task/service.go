package task

import (
	"context"
	"fmt"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
	"strings"
)

func (svr *Server) GetTaskDetailList(ctx context.Context, req *pb.GetTaskDetailListRequest) (*pb.GetTaskDetailListResponse, error) {
	pageSize := req.GetPageSize()
	if pageSize == 0 {
		pageSize = backend.DefaultPageSize
	}
	tasks, err := svr.B.GetTaskDetailList(req.GetLastUpdated(), pageSize)
	if nil != err {
		log.WithError(err).Error("RPC-API:GetTaskDetailList failed")
		return nil, ErrGetNodeTaskList
	}

	arr := make([]*pb.GetTaskDetailResponse, len(tasks))
	for i, task := range tasks {
		t := &pb.GetTaskDetailResponse{
			Information: task,
		}
		arr[i] = t
	}
	log.Debugf("RPC-API:GetTaskDetailList succeed, taskList len: {%d}", len(arr))
	return &pb.GetTaskDetailListResponse{
		Status:   0,
		Msg:      backend.OK,
		TaskList: arr,
	}, nil
}

func (svr *Server) GetTaskDetailListByTaskIds(ctx context.Context, req *pb.GetTaskDetailListByTaskIdsRequest) (*pb.GetTaskDetailListResponse, error) {

	if len(req.GetTaskIds()) == 0 {
		return nil, fmt.Errorf("%s, required taskIds", ErrGetNodeTaskList.Msg)
	}

	tasks, err := svr.B.GetTaskDetailListByTaskIds(req.GetTaskIds())
	if nil != err {
		log.WithError(err).Error("RPC-API:GetTaskDetailListByTaskIds failed")
		return nil, ErrGetNodeTaskList
	}

	arr := make([]*pb.GetTaskDetailResponse, len(tasks))
	for i, task := range tasks {
		t := &pb.GetTaskDetailResponse{
			Information: task,
		}
		arr[i] = t
	}
	log.Debugf("RPC-API:GetTaskDetailListByTaskIds succeed, taskIds len: {%d}, taskList len: {%d}", len(req.GetTaskIds()), len(arr))
	return &pb.GetTaskDetailListResponse{
		Status:   0,
		Msg:      backend.OK,
		TaskList: arr,
	}, nil
}

func (svr *Server) GetTaskEventList(ctx context.Context, req *pb.GetTaskEventListRequest) (*pb.GetTaskEventListResponse, error) {

	if "" == strings.Trim(req.GetTaskId(), "") {
		return nil, backend.NewRpcBizErr(ErrGetNodeTaskEventList.Code, "require taskId")
	}

	events, err := svr.B.GetTaskEventList(req.GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:QueryTaskEventList failed, taskId: {%s}", req.GetTaskId())
		errMsg := fmt.Sprintf("%s, taskId: {%s}", ErrGetNodeTaskEventList.Msg, req.GetTaskId())
		return nil, backend.NewRpcBizErr(ErrGetNodeTaskEventList.Code, errMsg)
	}
	log.Debugf("RPC-API:QueryTaskEventList succeed, taskId: {%s},  eventList len: {%d}", req.GetTaskId(), len(events))
	return &pb.GetTaskEventListResponse{
		Status:        0,
		Msg:           backend.OK,
		TaskEventList: events,
	}, nil
}

func (svr *Server) GetTaskEventListByTaskIds(ctx context.Context, req *pb.GetTaskEventListByTaskIdsRequest) (*pb.GetTaskEventListResponse, error) {

	if len(req.GetTaskIds()) == 0 {
		return nil, backend.NewRpcBizErr(ErrGetNodeTaskEventList.Code, "require taskIds")
	}

	events, err := svr.B.GetTaskEventListByTaskIds(req.GetTaskIds())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:QueryTaskEventListByTaskIds failed, taskId: {%v}", req.GetTaskIds())
		errMsg := fmt.Sprintf("%s, taskId: {%s}", ErrGetNodeTaskEventList.Msg, req.GetTaskIds())
		return nil, backend.NewRpcBizErr(ErrGetNodeTaskEventList.Code, errMsg)
	}
	log.Debugf("RPC-API:QueryTaskEventListByTaskIds succeed, taskIds: %v,  eventList len: {%d}", req.GetTaskIds(), len(events))
	return &pb.GetTaskEventListResponse{
		Status:        0,
		Msg:           backend.OK,
		TaskEventList: events,
	}, nil
}

func (svr *Server) PublishTaskDeclare(ctx context.Context, req *pb.PublishTaskDeclareRequest) (*pb.PublishTaskDeclareResponse, error) {

	if req.GetUserType() == apicommonpb.UserType_User_Unknown {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check userType failed, wrong userType: {%s}", req.GetUserType().String())
		return nil, ErrReqUserTypePublishTask
	}
	if "" == req.GetUser() {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check user failed, user is empty")
		return nil, ErrReqUserPublishTask
	}
	if len(req.GetSign()) == 0 {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check sign failed, sign is empty")
		return nil, ErrReqUserSignPublishTask
	}
	if req.GetOperationCost() == nil {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check OperationCost failed, OperationCost is empty")
		return nil, ErrReqOperationCostForPublishTask
	}
	if len(req.GetReceivers()) == 0 {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check receivers failed, receivers is empty")
		return nil, ErrReqReceiversForPublishTask
	}
	if len(req.GetDataSuppliers()) == 0 {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check dataSuppliers failed, dataSuppliers is empty")
		return nil, ErrReqDataSuppliersForPublishTask
	}
	if "" == req.GetAlgorithmCode() && "" == req.GetMetaAlgorithmId() {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check AlgorithmCode AND MetaAlgorithmId failed, AlgorithmCode AND MetaAlgorithmId is empty")
		return nil, ErrReqCalculateContractCodeForPublishTask
	}

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishTaskDeclare failed, query local identity failed")
		return nil, fmt.Errorf("query local identity failed")
	}

	taskMsg := types.NewTaskMessageFromRequest(req)

	checkPartyIdCache := make(map[string]struct{}, 0)
	checkPartyIdCache[req.GetSender().GetPartyId()] = struct{}{}

	if _, ok := checkPartyIdCache[req.GetAlgoSupplier().GetPartyId()]; ok {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check partyId of algoSupplier failed, this partyId has alreay exist, partyId: {%s}",
			req.GetAlgoSupplier().GetPartyId())
		return nil, fmt.Errorf("The partyId of the task participants cannot be repeated on algoSupplier, partyId: {%s}", req.GetAlgoSupplier().GetPartyId())
	}
	checkPartyIdCache[req.GetAlgoSupplier().GetPartyId()] = struct{}{}

	for _, v := range req.GetReceivers() {
		if _, ok := checkPartyIdCache[v.GetPartyId()]; ok {
			log.Errorf("RPC-API:PublishTaskDeclare failed, check partyId of receiver failed, this partyId has alreay exist, partyId: {%s}",
				v.GetPartyId())
			return nil, fmt.Errorf("The partyId of the task participants cannot be repeated on receivers, partyId: {%s}", v.GetPartyId())
		}
		checkPartyIdCache[v.GetPartyId()] = struct{}{}
	}

	// add taskId
	taskId := taskMsg.GenTaskId()

	if err = svr.B.SendMsg(taskMsg); nil != err {
		log.WithError(err).Errorf("RPC-API:PublishTaskDeclare failed, send task msg failed, taskId: {%s}",
			taskId)
		errMsg := fmt.Sprintf("%s, taskId: {%s}", ErrSendTaskMsgByTaskId.Msg, taskId)
		return nil, backend.NewRpcBizErr(ErrSendTaskMsgByTaskId.Code, errMsg)
	}
	//log.Debugf("RPC-API:PublishTaskDeclare succeed, taskId: {%s}, taskMsg: %s", taskId, taskMsg.String())
	log.Debugf("RPC-API:PublishTaskDeclare succeed, userType: {%s}, user: {%s}, publish taskId: {%s}", req.GetUserType(), req.GetUser(), taskId)
	return &pb.PublishTaskDeclareResponse{
		Status: 0,
		Msg:    backend.OK,
		TaskId: taskId,
	}, nil
}

func (svr *Server) TerminateTask(ctx context.Context, req *pb.TerminateTaskRequest) (*apicommonpb.SimpleResponse, error) {
	if req.GetUserType() == apicommonpb.UserType_User_Unknown {
		return nil, ErrReqUserTypeTerminateTask
	}
	if "" == req.GetUser() {
		return nil, ErrReqUserTerminateTask
	}
	if "" == req.GetTaskId() {
		return nil, ErrReqTaskIdTerminateTask
	}
	if len(req.GetSign()) == 0 {
		return nil, ErrReqUserSignTerminateTask
	}

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:TerminateTask failed, query local identity failed")
		return nil, fmt.Errorf("query local identity failed")
	}

	task, err := svr.B.GetLocalTask(req.GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:TerminateTask failed, query local task failed")
		return nil, ErrTerminateTask
	}

	// check user
	if task.GetUser() != req.GetUser() ||
		task.GetUserType() != req.GetUserType() {
		log.WithError(err).Errorf("terminate task user and publish task user must be same, taskId: {%s}",
			task.GetTaskId())
		return nil, ErrTerminateTask
	}
	// todo verify user sign with terminate task


	taskTerminateMsg := types.NewTaskTerminateMsg(req.GetUserType(), req.GetUser(), req.GetTaskId(), req.GetSign())

	if err = svr.B.SendMsg(taskTerminateMsg); nil != err {
		log.WithError(err).Errorf("RPC-API:TerminateTask failed, send task terminate msg failed, taskId: {%s}",
			req.GetTaskId())

		errMsg := fmt.Sprintf("%s, taskId: {%s}", ErrTerminateTaskMsg.Msg, req.GetTaskId())
		return nil, backend.NewRpcBizErr(ErrTerminateTaskMsg.Code, errMsg)
	}
	log.Debugf("RPC-API:TerminateTask succeed, userType: {%s}, user: {%s}, taskId: {%s}", req.GetUserType(), req.GetUser(), req.GetTaskId())
	return &apicommonpb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
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
