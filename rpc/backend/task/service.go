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
		return &pb.GetTaskDetailListResponse { Status: backend.ErrQueryNodeTaskList.ErrCode(), Msg: backend.ErrQueryNodeTaskList.Error()}, nil
	}

	arr := make([]*pb.GetTaskDetail, len(tasks))
	for i, task := range tasks {
		t := &pb.GetTaskDetail{
			Information: task,
		}
		arr[i] = t
	}
	log.Debugf("RPC-API:GetTaskDetailList succeed, taskList len: {%d}", len(arr))
	return &pb.GetTaskDetailListResponse{
		Status:   0,
		Msg:      backend.OK,
		Tasks: arr,
	}, nil
}

func (svr *Server) GetTaskDetailListByTaskIds(ctx context.Context, req *pb.GetTaskDetailListByTaskIdsRequest) (*pb.GetTaskDetailListResponse, error) {

	if len(req.GetTaskIds()) == 0 {
		return &pb.GetTaskDetailListResponse { Status: backend.ErrRequireParams.ErrCode(), Msg: "required taskIds"}, nil
	}

	tasks, err := svr.B.GetTaskDetailListByTaskIds(req.GetTaskIds())
	if nil != err {
		log.WithError(err).Error("RPC-API:GetTaskDetailListByTaskIds failed")
		return &pb.GetTaskDetailListResponse { Status: backend.ErrQueryNodeTaskList.ErrCode(), Msg: backend.ErrQueryNodeTaskList.Error()}, nil
	}

	arr := make([]*pb.GetTaskDetail, len(tasks))
	for i, task := range tasks {
		t := &pb.GetTaskDetail{
			Information: task,
		}
		arr[i] = t
	}
	log.Debugf("RPC-API:GetTaskDetailListByTaskIds succeed, taskIds len: {%d}, taskList len: {%d}", len(req.GetTaskIds()), len(arr))
	return &pb.GetTaskDetailListResponse{
		Status:   0,
		Msg:      backend.OK,
		Tasks: arr,
	}, nil
}

func (svr *Server) GetTaskEventList(ctx context.Context, req *pb.GetTaskEventListRequest) (*pb.GetTaskEventListResponse, error) {

	if "" == strings.Trim(req.GetTaskId(), "") {
		return &pb.GetTaskEventListResponse { Status: backend.ErrRequireParams.ErrCode(), Msg: "required taskId"}, nil
	}

	events, err := svr.B.GetTaskEventList(req.GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:QueryTaskEventList failed, taskId: {%s}", req.GetTaskId())
		errMsg := fmt.Sprintf("%s, taskId: {%s}", backend.ErrQueryNodeTaskEventList.Error(), req.GetTaskId())
		return &pb.GetTaskEventListResponse { Status: backend.ErrQueryNodeTaskEventList.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:QueryTaskEventList succeed, taskId: {%s},  eventList len: {%d}", req.GetTaskId(), len(events))
	return &pb.GetTaskEventListResponse{
		Status:        0,
		Msg:           backend.OK,
		TaskEvents: events,
	}, nil
}

func (svr *Server) GetTaskEventListByTaskIds(ctx context.Context, req *pb.GetTaskEventListByTaskIdsRequest) (*pb.GetTaskEventListResponse, error) {

	if len(req.GetTaskIds()) == 0 {
		return &pb.GetTaskEventListResponse { Status: backend.ErrRequireParams.ErrCode(), Msg: "required taskIdx"}, nil
	}

	events, err := svr.B.GetTaskEventListByTaskIds(req.GetTaskIds())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:QueryTaskEventListByTaskIds failed, taskId: {%v}", req.GetTaskIds())
		errMsg := fmt.Sprintf("%s, taskId: {%s}",backend.ErrQueryNodeTaskEventList.Error(), req.GetTaskIds())
		return &pb.GetTaskEventListResponse { Status: backend.ErrQueryNodeTaskEventList.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:QueryTaskEventListByTaskIds succeed, taskIds: %v,  eventList len: {%d}", req.GetTaskIds(), len(events))
	return &pb.GetTaskEventListResponse{
		Status:        0,
		Msg:           backend.OK,
		TaskEvents: events,
	}, nil
}

func (svr *Server) PublishTaskDeclare(ctx context.Context, req *pb.PublishTaskDeclareRequest) (*pb.PublishTaskDeclareResponse, error) {

	if req.GetUserType() == apicommonpb.UserType_User_Unknown {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check userType failed, wrong userType: {%s}", req.GetUserType().String())
		return &pb.PublishTaskDeclareResponse{ Status:  backend.ErrRequireParams.ErrCode(), Msg: "unknown userType"}, nil
	}
	if "" == req.GetUser() {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check user failed, user is empty")
		return &pb.PublishTaskDeclareResponse{ Status:  backend.ErrRequireParams.ErrCode(), Msg: "require user"}, nil
	}
	if len(req.GetSign()) == 0 {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check sign failed, sign is empty")
		return &pb.PublishTaskDeclareResponse{ Status:  backend.ErrRequireParams.ErrCode(), Msg: "require sign"}, nil
	}
	if req.GetOperationCost() == nil {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check OperationCost failed, OperationCost is empty")
		return &pb.PublishTaskDeclareResponse{ Status:  backend.ErrRequireParams.ErrCode(), Msg: "require operationCost"}, nil
	}
	if len(req.GetReceivers()) == 0 {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check receivers failed, receivers is empty")
		return &pb.PublishTaskDeclareResponse{ Status:  backend.ErrRequireParams.ErrCode(), Msg: "require receivers"}, nil
	}
	if len(req.GetDataSuppliers()) == 0 {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check dataSuppliers failed, dataSuppliers is empty")
		return &pb.PublishTaskDeclareResponse{ Status:  backend.ErrRequireParams.ErrCode(), Msg: "require dataSuppliers"}, nil
	}
	if "" == req.GetAlgorithmCode() && "" == req.GetMetaAlgorithmId() {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check AlgorithmCode AND MetaAlgorithmId failed, AlgorithmCode AND MetaAlgorithmId is empty")
		return &pb.PublishTaskDeclareResponse{ Status:  backend.ErrRequireParams.ErrCode(), Msg: "require algorithmCode OR metaAlgorithmId"}, nil
	}
	if 0 == req.GetDataPolicyType() {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check DataPolicyType failed, DataPolicyType is zero value")
		return &pb.PublishTaskDeclareResponse{ Status:  backend.ErrRequireParams.ErrCode(), Msg: "unknown dataPolicyType"}, nil
	}
	if "" == req.GetDataPolicyOption() {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check DataPolicyOption failed, DataPolicyOption is empty")
		return &pb.PublishTaskDeclareResponse{ Status:  backend.ErrRequireParams.ErrCode(), Msg: "require dataPolicyOption"}, nil
	}
	if 0 == req.GetPowerPolicyType() {

	}
	if "" == req.GetPowerPolicyOption() {

	}

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishTaskDeclare failed, query local identity failed")
		return &pb.PublishTaskDeclareResponse { Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}

	taskMsg := types.NewTaskMessageFromRequest(req)

	checkPartyIdCache := make(map[string]struct{}, 0)
	checkPartyIdCache[req.GetSender().GetPartyId()] = struct{}{}

	if _, ok := checkPartyIdCache[req.GetAlgoSupplier().GetPartyId()]; ok {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check partyId of algoSupplier failed, this partyId has alreay exist, partyId: {%s}",
			req.GetAlgoSupplier().GetPartyId())
		return &pb.PublishTaskDeclareResponse { Status: backend.ErrPublishTaskMsg.ErrCode(), Msg: fmt.Sprintf("The partyId of the task participants cannot be repeated on algoSupplier, partyId: {%s}", req.GetAlgoSupplier().GetPartyId())}, nil
	}
	checkPartyIdCache[req.GetAlgoSupplier().GetPartyId()] = struct{}{}

	for _, v := range req.GetDataSuppliers() {
		if _, ok := checkPartyIdCache[v.GetPartyId()]; ok {
			log.Errorf("RPC-API:PublishTaskDeclare failed, check partyId of dataSuppliers failed, this partyId has alreay exist, partyId: {%s}",
				v.GetPartyId())
			return &pb.PublishTaskDeclareResponse { Status: backend.ErrPublishTaskMsg.ErrCode(), Msg: fmt.Sprintf("The partyId of the task participants cannot be repeated on dataSuppliers, partyId: {%s}", v.GetPartyId())}, nil
		}
		checkPartyIdCache[v.GetPartyId()] = struct{}{}
	}

	// check partyId of powerSuppliers
	for _, v := range req.GetReceivers() {
		if _, ok := checkPartyIdCache[v.GetPartyId()]; ok {
			log.Errorf("RPC-API:PublishTaskDeclare failed, check partyId of receiver failed, this partyId has alreay exist, partyId: {%s}",
				v.GetPartyId())
			return &pb.PublishTaskDeclareResponse { Status: backend.ErrPublishTaskMsg.ErrCode(), Msg: fmt.Sprintf("The partyId of the task participants cannot be repeated on receivers, partyId: {%s}", v.GetPartyId())}, nil
		}
		checkPartyIdCache[v.GetPartyId()] = struct{}{}
	}

	// generate and store taskId
	taskId := taskMsg.GenTaskId()

	if err = svr.B.SendMsg(taskMsg); nil != err {
		log.WithError(err).Errorf("RPC-API:PublishTaskDeclare failed, send task msg failed, taskId: {%s}",
			taskId)
		errMsg := fmt.Sprintf("%s, taskId: {%s}", backend.ErrPublishTaskMsg.Error(), taskId)
		return &pb.PublishTaskDeclareResponse { Status: backend.ErrPublishTaskMsg.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:PublishTaskDeclare succeed, userType: {%s}, user: {%s}, publish taskId: {%s}", req.GetUserType(), req.GetUser(), taskId)
	return &pb.PublishTaskDeclareResponse{
		Status: 0,
		Msg:    backend.OK,
		TaskId: taskId,
	}, nil
}

func (svr *Server) TerminateTask(ctx context.Context, req *pb.TerminateTaskRequest) (*apicommonpb.SimpleResponse, error) {
	if req.GetUserType() == apicommonpb.UserType_User_Unknown {
		return &apicommonpb.SimpleResponse{ Status:  backend.ErrRequireParams.ErrCode(), Msg: "unknown userType"}, nil
	}
	if "" == req.GetUser() {
		return &apicommonpb.SimpleResponse{ Status:  backend.ErrRequireParams.ErrCode(), Msg: "require user"}, nil
	}
	if "" == req.GetTaskId() {
		return &apicommonpb.SimpleResponse{ Status:  backend.ErrRequireParams.ErrCode(), Msg: "require taskId"}, nil
	}
	if len(req.GetSign()) == 0 {
		return &apicommonpb.SimpleResponse{ Status:  backend.ErrRequireParams.ErrCode(), Msg: "require sign"}, nil
	}

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:TerminateTask failed, query local identity failed")
		return &apicommonpb.SimpleResponse{ Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}

	task, err := svr.B.GetLocalTask(req.GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:TerminateTask failed, query local task failed")
		return &apicommonpb.SimpleResponse{ Status: backend.ErrTerminateTaskMsg.ErrCode(), Msg: "query local task failed"}, nil
	}

	// check user
	if task.GetUser() != req.GetUser() ||
		task.GetUserType() != req.GetUserType() {
		log.WithError(err).Errorf("terminate task user and publish task user must be same, taskId: {%s}",
			task.GetTaskId())
		return &apicommonpb.SimpleResponse{ Status: backend.ErrTerminateTaskMsg.ErrCode(), Msg: fmt.Sprintf("terminate task user and publish task user must be same, taskId: {%s}",
			task.GetTaskId())}, nil
	}
	// todo verify user sign with terminate task


	taskTerminateMsg := types.NewTaskTerminateMsg(req.GetUserType(), req.GetUser(), req.GetTaskId(), req.GetSign())

	if err = svr.B.SendMsg(taskTerminateMsg); nil != err {
		log.WithError(err).Errorf("RPC-API:TerminateTask failed, send task terminate msg failed, taskId: {%s}",
			req.GetTaskId())

		errMsg := fmt.Sprintf("%s, taskId: {%s}", backend.ErrTerminateTaskMsg.Error(), req.GetTaskId())
		return &apicommonpb.SimpleResponse{ Status: backend.ErrTerminateTaskMsg.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:TerminateTask succeed, userType: {%s}, user: {%s}, taskId: {%s}", req.GetUserType(), req.GetUser(), req.GetTaskId())
	return &apicommonpb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}
