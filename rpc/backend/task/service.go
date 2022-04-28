package task

import (
	"context"
	"fmt"
	pb "github.com/Metisnetwork/Metis-Carrier/lib/api"
	libtypes "github.com/Metisnetwork/Metis-Carrier/lib/types"
	"github.com/Metisnetwork/Metis-Carrier/rpc/backend"
	"github.com/Metisnetwork/Metis-Carrier/signsuite"
	"github.com/Metisnetwork/Metis-Carrier/types"

	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

func (svr *Server) GetLocalTaskDetailList(ctx context.Context, req *pb.GetTaskDetailListRequest) (*pb.GetTaskDetailListResponse, error) {
	pageSize := req.GetPageSize()
	if pageSize == 0 {
		pageSize = backend.DefaultPageSize
	}
	tasks, err := svr.B.GetLocalTaskDetailList(req.GetLastUpdated(), pageSize)
	if nil != err {
		log.WithError(err).Error("RPC-API:GetLocalTaskDetailList failed")
		return &pb.GetTaskDetailListResponse{Status: backend.ErrQueryLocalTaskList.ErrCode(), Msg: backend.ErrQueryLocalTaskList.Error()}, nil
	}

	log.Debugf("RPC-API:GetLocalTaskDetailList succeed, taskList len: {%d}", len(tasks))
	return &pb.GetTaskDetailListResponse{
		Status: 0,
		Msg:    backend.OK,
		Tasks:  tasks,
	}, nil
}

func (svr *Server) GetGlobalTaskDetailList(ctx context.Context, req *pb.GetTaskDetailListRequest) (*pb.GetTaskDetailListResponse, error) {
	pageSize := req.GetPageSize()
	if pageSize == 0 {
		pageSize = backend.DefaultPageSize
	}
	tasks, err := svr.B.GetGlobalTaskDetailList(req.GetLastUpdated(), pageSize)
	if nil != err {
		log.WithError(err).Error("RPC-API:GetGlobalTaskDetailList failed")
		return &pb.GetTaskDetailListResponse{Status: backend.ErrQueryGlobalTaskList.ErrCode(), Msg: backend.ErrQueryGlobalTaskList.Error()}, nil
	}

	log.Debugf("RPC-API:GetGlobalTaskDetailList succeed, taskList len: {%d}", len(tasks))
	return &pb.GetTaskDetailListResponse{
		Status: 0,
		Msg:    backend.OK,
		Tasks:  tasks,
	}, nil
}

func (svr *Server) GetTaskDetailListByTaskIds(ctx context.Context, req *pb.GetTaskDetailListByTaskIdsRequest) (*pb.GetTaskDetailListResponse, error) {

	if len(req.GetTaskIds()) == 0 {
		return &pb.GetTaskDetailListResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "required taskIds"}, nil
	}

	tasks, err := svr.B.GetTaskDetailListByTaskIds(req.GetTaskIds())
	if nil != err {
		log.WithError(err).Error("RPC-API:GetTaskDetailListByTaskIds failed")
		return &pb.GetTaskDetailListResponse{Status: backend.ErrQueryTaskListByIds.ErrCode(), Msg: backend.ErrQueryTaskListByIds.Error()}, nil
	}

	log.Debugf("RPC-API:GetTaskDetailListByTaskIds succeed, taskIds len: {%d}, taskList len: {%d}", len(req.GetTaskIds()), len(tasks))
	return &pb.GetTaskDetailListResponse{
		Status: 0,
		Msg:    backend.OK,
		Tasks:  tasks,
	}, nil
}

func (svr *Server) GetTaskEventList(ctx context.Context, req *pb.GetTaskEventListRequest) (*pb.GetTaskEventListResponse, error) {

	if "" == strings.Trim(req.GetTaskId(), "") {
		return &pb.GetTaskEventListResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "required taskId"}, nil
	}

	events, err := svr.B.GetTaskEventList(req.GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:GetTaskEventList failed, taskId: {%s}", req.GetTaskId())
		errMsg := fmt.Sprintf("%s, taskId: {%s}", backend.ErrQueryTaskEventList.Error(), req.GetTaskId())
		return &pb.GetTaskEventListResponse{Status: backend.ErrQueryTaskEventList.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:GetTaskEventList succeed, taskId: {%s},  eventList len: {%d}", req.GetTaskId(), len(events))
	return &pb.GetTaskEventListResponse{
		Status:     0,
		Msg:        backend.OK,
		TaskEvents: events,
	}, nil
}

func (svr *Server) GetTaskEventListByTaskIds(ctx context.Context, req *pb.GetTaskEventListByTaskIdsRequest) (*pb.GetTaskEventListResponse, error) {

	if len(req.GetTaskIds()) == 0 {
		return &pb.GetTaskEventListResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "required taskIdx"}, nil
	}

	events, err := svr.B.GetTaskEventListByTaskIds(req.GetTaskIds())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:GetTaskEventListByTaskIds failed, taskId: {%v}", req.GetTaskIds())
		errMsg := fmt.Sprintf("%s, taskId: {%s}", backend.ErrQueryTaskEventList.Error(), req.GetTaskIds())
		return &pb.GetTaskEventListResponse{Status: backend.ErrQueryTaskEventList.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:GetTaskEventListByTaskIds succeed, taskIds: %v,  eventList len: {%d}", req.GetTaskIds(), len(events))
	return &pb.GetTaskEventListResponse{
		Status:     0,
		Msg:        backend.OK,
		TaskEvents: events,
	}, nil
}

func (svr *Server) PublishTaskDeclare(ctx context.Context, req *pb.PublishTaskDeclareRequest) (*pb.PublishTaskDeclareResponse, error) {

	identity, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishTaskDeclare failed, query local identity failed")
		return &pb.PublishTaskDeclareResponse{Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}

	if "" == strings.Trim(req.GetTaskName(), "") {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check taskName failed, taskName is empty")
		return &pb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require taskName"}, nil
	}
	if req.GetUserType() == libtypes.UserType_User_Unknown {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check userType failed, wrong userType: {%s}", req.GetUserType().String())
		return &pb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown userType"}, nil
	}
	if "" == strings.Trim(req.GetUser(), "") {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check user failed, user is empty")
		return &pb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require user"}, nil
	}
	if nil == req.GetSender() {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check taskSender failed, taskSender is empty")
		return &pb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require taskSender"}, nil
	}
	if req.GetSender().GetIdentityId() != identity.GetIdentityId() {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check taskSender failed, the taskSender is not current identity")
		return &pb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "invalid taskSender"}, nil
	}
	if nil == req.GetAlgoSupplier() {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check algoSupplier failed, algoSupplier is empty")
		return &pb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require algoSupplier"}, nil
	}
	if len(req.GetDataSuppliers()) == 0 {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check dataSuppliers failed, dataSuppliers is empty")
		return &pb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require dataSuppliers"}, nil
	}
	if len(req.GetReceivers()) == 0 {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check receivers failed, receivers is empty")
		return &pb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require receivers"}, nil
	}
	if len(req.GetDataPolicyTypes()) == 0 {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check DataPolicyType failed, DataPolicyType is zero value")
		return &pb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown dataPolicyType"}, nil
	}
	if len(req.GetDataPolicyOptions()) == 0 {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check DataPolicyOption failed, DataPolicyOption is empty")
		return &pb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require dataPolicyOption"}, nil
	}
	if len(req.GetPowerPolicyTypes()) == 0 {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check PowerPolicyType failed, PowerPolicyType is zero value")
		return &pb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown powerPolicyType"}, nil
	}
	if len(req.GetPowerPolicyOptions()) == 0 {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check PowerPolicyOption failed, PowerPolicyOption is empty")
		return &pb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require powerPolicyOption"}, nil
	}
	if len(req.GetDataFlowPolicyTypes()) == 0  {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check DataFlowPolicyType failed, DataFlowPolicyType is zero value")
		return &pb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown dataFlowPolicyType"}, nil
	}
	// Maybe the dataFlowPolicyOption is empty,
	// Because dataFlowPolicyType can already represent the way the data flows.
	//
	//if "" == strings.Trim(req.GetDataFlowPolicyOptions(), "") {
	//	log.Errorf("RPC-API:PublishTaskDeclare failed, check DataFlowPolicyOption failed, DataFlowPolicyOption is empty")
	//	return &pb.PublishTaskDeclareResponse{ Status:  backend.ErrRequireParams.ErrCode(), Msg: "require dataFlowPolicyOption"}, nil
	//}
	if req.GetOperationCost() == nil {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check OperationCost failed, OperationCost is empty")
		return &pb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require operationCost"}, nil
	}
	if "" == strings.Trim(req.GetAlgorithmCode(), "") && "" == strings.Trim(req.GetMetaAlgorithmId(), "") {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check AlgorithmCode AND MetaAlgorithmId failed, AlgorithmCode AND MetaAlgorithmId is empty")
		return &pb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require algorithmCode OR metaAlgorithmId"}, nil
	}
	// Maybe the contractExtraParams is empty,
	// Because contractExtraParams is not necessary.
	// Whether contractExtraParams has a value depends on the specific matching algorithmCode requirements.
	//
	//if "" == strings.Trim(req.GetAlgorithmCodeExtraParams(), "") {
	//	log.Errorf("RPC-API:PublishTaskDeclare failed, check ContractExtraParams failed, ContractExtraParams is empty")
	//	return &pb.PublishTaskDeclareResponse{ Status:  backend.ErrRequireParams.ErrCode(), Msg: "require contractExtraParams"}, nil
	//}
	if len(req.GetSign()) == 0 {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check sign failed, sign is empty")
		return &pb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require sign"}, nil
	}

	taskMsg := types.NewTaskMessageFromRequest(req)
	/*from, err := signsuite.Sender(req.GetUserType(), taskMsg.Hash(), req.GetSign())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishTaskDeclare failed, cannot fetch sender from sign, userType: {%s}, user: {%s}",
			req.GetUserType().String(), req.GetUser())
		return &pb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "cannot fetch sender from sign"}, nil
	}
	if from != req.GetUser() {
		log.WithError(err).Errorf("RPC-API:PublishTaskDeclare failed, sender from sign and user is not sameone, userType: {%s}, user: {%s}, sender of sign: {%s}",
			req.GetUserType().String(), req.GetUser(), from)
		return &pb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "the user sign is invalid"}, nil
	}*/
	// Maybe the desc is empty,
	// Because desc is not necessary.
	//
	//if "" == strings.Trim(req.GetDesc(), "") {
	//	log.Errorf("RPC-API:PublishTaskDeclare failed, check Desc failed, Desc is empty")
	//	return &pb.PublishTaskDeclareResponse{ Status:  backend.ErrRequireParams.ErrCode(), Msg: "require desc"}, nil
	//}

	//checkPartyIdCache := make(map[string]struct{}, 0)
	//checkPartyIdCache[req.GetSender().GetPartyId()] = struct{}{}
	//
	//if _, ok := checkPartyIdCache[req.GetAlgoSupplier().GetPartyId()]; ok {
	//	log.Errorf("RPC-API:PublishTaskDeclare failed, check partyId of algoSupplier failed, this partyId has alreay exist, partyId: {%s}",
	//		req.GetAlgoSupplier().GetPartyId())
	//	return &pb.PublishTaskDeclareResponse { Status: backend.ErrPublishTaskMsg.ErrCode(), Msg: fmt.Sprintf("The partyId of the task participants cannot be repeated on algoSupplier, partyId: {%s}", req.GetAlgoSupplier().GetPartyId())}, nil
	//}
	//checkPartyIdCache[req.GetAlgoSupplier().GetPartyId()] = struct{}{}
	//
	//for _, v := range req.GetDataSuppliers() {
	//	if _, ok := checkPartyIdCache[v.GetPartyId()]; ok {
	//		log.Errorf("RPC-API:PublishTaskDeclare failed, check partyId of dataSuppliers failed, this partyId has alreay exist, partyId: {%s}",
	//			v.GetPartyId())
	//		return &pb.PublishTaskDeclareResponse { Status: backend.ErrPublishTaskMsg.ErrCode(), Msg: fmt.Sprintf("The partyId of the task participants cannot be repeated on dataSuppliers, partyId: {%s}", v.GetPartyId())}, nil
	//	}
	//	checkPartyIdCache[v.GetPartyId()] = struct{}{}
	//}
	//
	//// check partyId of powerSuppliers
	//powerPartyIds, err := policy.FetchPowerPartyIds(req.GetPowerPolicyTypes(), req.GetPowerPolicyOptions())
	//if nil != err {
	//	log.WithError(err).Errorf("not fetch partyIds from task powerPolicy")
	//	return &pb.PublishTaskDeclareResponse { Status: backend.ErrPublishTaskMsg.ErrCode(), Msg: "not fetch partyIds from task powerPolicy" }, nil
	//}
	//for _, partyId := range powerPartyIds {
	//	if _, ok := checkPartyIdCache[partyId]; ok {
	//		log.Errorf("RPC-API:PublishTaskDeclare failed, check partyId of powerSuppliers failed, this partyId has alreay exist, partyId: {%s}",
	//			partyId)
	//		return &pb.PublishTaskDeclareResponse { Status: backend.ErrPublishTaskMsg.ErrCode(), Msg: fmt.Sprintf("The partyId of the task participants cannot be repeated on powerSuppliers, partyId: {%s}", partyId)}, nil
	//	}
	//	checkPartyIdCache[partyId] = struct{}{}
	//}
	//
	//// check partyId of receivers
	//for _, v := range req.GetReceivers() {
	//	if _, ok := checkPartyIdCache[v.GetPartyId()]; ok {
	//		log.Errorf("RPC-API:PublishTaskDeclare failed, check partyId of receiver failed, this partyId has alreay exist, partyId: {%s}",
	//			v.GetPartyId())
	//		return &pb.PublishTaskDeclareResponse { Status: backend.ErrPublishTaskMsg.ErrCode(), Msg: fmt.Sprintf("The partyId of the task participants cannot be repeated on receivers, partyId: {%s}", v.GetPartyId())}, nil
	//	}
	//	checkPartyIdCache[v.GetPartyId()] = struct{}{}
	//}

	// generate and store taskId
	taskId := taskMsg.GenTaskId()

	if err = svr.B.SendMsg(taskMsg); nil != err {
		log.WithError(err).Errorf("RPC-API:PublishTaskDeclare failed, send task msg failed, taskId: {%s}",
			taskId)
		errMsg := fmt.Sprintf("%s, taskId: {%s}", backend.ErrPublishTaskMsg.Error(), taskId)
		return &pb.PublishTaskDeclareResponse{Status: backend.ErrPublishTaskMsg.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:PublishTaskDeclare succeed, userType: {%s}, user: {%s}, publish taskId: {%s}", req.GetUserType(), req.GetUser(), taskId)
	return &pb.PublishTaskDeclareResponse{
		Status: 0,
		Msg:    backend.OK,
		TaskId: taskId,
	}, nil
}

func (svr *Server) TerminateTask(ctx context.Context, req *pb.TerminateTaskRequest) (*libtypes.SimpleResponse, error) {

	_, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:TerminateTask failed, query local identity failed")
		return &libtypes.SimpleResponse{Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}

	if req.GetUserType() == libtypes.UserType_User_Unknown {
		return &libtypes.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown userType"}, nil
	}
	if "" == req.GetUser() {
		return &libtypes.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require user"}, nil
	}
	if "" == req.GetTaskId() {
		return &libtypes.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require taskId"}, nil
	}
	if len(req.GetSign()) == 0 {
		return &libtypes.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require sign"}, nil
	}

	taskTerminateMsg := types.NewTaskTerminateMsg(req.GetUserType(), req.GetUser(), req.GetTaskId(), req.GetSign())
	from, err := signsuite.Sender(req.GetUserType(), taskTerminateMsg.Hash(), req.GetSign())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:TerminateTask failed, cannot fetch sender from sign, userType: {%s}, user: {%s}",
			req.GetUserType().String(), req.GetUser())
		return &libtypes.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "cannot fetch sender from sign"}, nil
	}
	if from != req.GetUser() {
		log.WithError(err).Errorf("RPC-API:TerminateTask failed, sender from sign and user is not sameone, userType: {%s}, user: {%s}, sender of sign: {%s}",
			req.GetUserType().String(), req.GetUser(), from)
		return &libtypes.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "the user sign is invalid"}, nil
	}

	task, err := svr.B.GetLocalTask(req.GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:TerminateTask failed, query local task failed")
		return &libtypes.SimpleResponse{Status: backend.ErrTerminateTaskMsg.ErrCode(), Msg: "query local task failed"}, nil
	}

	// check user
	if task.GetInformation().GetUser() != req.GetUser() ||
		task.GetInformation().GetUserType() != req.GetUserType() {
		log.WithError(err).Errorf("terminate task user and publish task user must be same, taskId: {%s}",
			task.GetInformation().GetTaskId())
		return &libtypes.SimpleResponse{Status: backend.ErrTerminateTaskMsg.ErrCode(), Msg: fmt.Sprintf("terminate task user and publish task user must be same, taskId: {%s}",
			task.GetInformation().GetTaskId())}, nil
	}

	if err = svr.B.SendMsg(taskTerminateMsg); nil != err {
		log.WithError(err).Errorf("RPC-API:TerminateTask failed, send task terminate msg failed, taskId: {%s}",
			req.GetTaskId())

		errMsg := fmt.Sprintf("%s, taskId: {%s}", backend.ErrTerminateTaskMsg.Error(), req.GetTaskId())
		return &libtypes.SimpleResponse{Status: backend.ErrTerminateTaskMsg.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:TerminateTask succeed, userType: {%s}, user: {%s}, taskId: {%s}", req.GetUserType(), req.GetUser(), req.GetTaskId())
	return &libtypes.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}

func (svr *Server) GenerateObServerProxyWalletAddress(ctx context.Context, req *emptypb.Empty) (*pb.GenerateObServerProxyWalletAddressResponse, error) {

	wallet, err := svr.B.GenerateObServerProxyWalletAddress()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:GenerateObServerProxyWalletAddress failed")
		return &pb.GenerateObServerProxyWalletAddressResponse{Status: backend.ErrGenerateObserverProxyWallet.ErrCode(), Msg: backend.ErrGenerateObserverProxyWallet.Error()}, nil
	}
	log.Debugf("RPC-API:GenerateObServerProxyWalletAddress Succeed, wallet: {%s}", wallet)
	return &pb.GenerateObServerProxyWalletAddressResponse{
		Status:  0,
		Msg:     backend.OK,
		Address: wallet,
	}, nil
}

func (svr *Server) EstimateTaskGas(ctx context.Context, req *pb.EstimateTaskGasRequest) (*pb.EstimateTaskGasResponse, error) {
	gasLimit, gasPrice, err := svr.B.EstimateTaskGas(req.GetDataTokenAddresses())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:EstimateTaskGas failed")
		return &pb.EstimateTaskGasResponse{Status: backend.ErrEstimateTaskGas.ErrCode(), Msg: backend.ErrEstimateTaskGas.Error()}, nil
	}
	log.Debugf("RPC-API:EstimateGas Succeed: gas {%d}", gasLimit)
	return &pb.EstimateTaskGasResponse{
		Status:   0,
		Msg:      backend.OK,
		GasLimit: gasLimit,
		GasPrice: gasPrice.Uint64(),
	}, nil
}
