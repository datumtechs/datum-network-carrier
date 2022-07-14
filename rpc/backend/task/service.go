package task

import (
	"context"
	"fmt"
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/datumtechs/datum-network-carrier/rpc/backend"
	"github.com/datumtechs/datum-network-carrier/signsuite"
	"github.com/datumtechs/datum-network-carrier/types"

	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

func (svr *Server) GetLocalTaskDetailList(ctx context.Context, req *carrierapipb.GetTaskDetailListRequest) (*carrierapipb.GetTaskDetailListResponse, error) {
	pageSize := req.GetPageSize()
	if pageSize == 0 {
		pageSize = backend.DefaultPageSize
	}
	tasks, err := svr.B.GetLocalTaskDetailList(req.GetLastUpdated(), pageSize)
	if nil != err {
		log.WithError(err).Error("RPC-API:GetLocalTaskDetailList failed")
		return &carrierapipb.GetTaskDetailListResponse{Status: backend.ErrQueryLocalTaskList.ErrCode(), Msg: backend.ErrQueryLocalTaskList.Error()}, nil
	}

	log.Debugf("RPC-API:GetLocalTaskDetailList succeed, taskList len: {%d}", len(tasks))
	return &carrierapipb.GetTaskDetailListResponse{
		Status: 0,
		Msg:    backend.OK,
		Tasks:  tasks,
	}, nil
}

func (svr *Server) GetGlobalTaskDetailList(ctx context.Context, req *carrierapipb.GetTaskDetailListRequest) (*carrierapipb.GetTaskDetailListResponse, error) {
	pageSize := req.GetPageSize()
	if pageSize == 0 {
		pageSize = backend.DefaultPageSize
	}
	tasks, err := svr.B.GetGlobalTaskDetailList(req.GetLastUpdated(), pageSize)
	if nil != err {
		log.WithError(err).Error("RPC-API:GetGlobalTaskDetailList failed")
		return &carrierapipb.GetTaskDetailListResponse{Status: backend.ErrQueryGlobalTaskList.ErrCode(), Msg: backend.ErrQueryGlobalTaskList.Error()}, nil
	}

	log.Debugf("RPC-API:GetGlobalTaskDetailList succeed, taskList len: {%d}", len(tasks))
	return &carrierapipb.GetTaskDetailListResponse{
		Status: 0,
		Msg:    backend.OK,
		Tasks:  tasks,
	}, nil
}

func (svr *Server) GetTaskDetailListByTaskIds(ctx context.Context, req *carrierapipb.GetTaskDetailListByTaskIdsRequest) (*carrierapipb.GetTaskDetailListResponse, error) {

	if len(req.GetTaskIds()) == 0 {
		return &carrierapipb.GetTaskDetailListResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "required taskIds"}, nil
	}

	tasks, err := svr.B.GetTaskDetailListByTaskIds(req.GetTaskIds())
	if nil != err {
		log.WithError(err).Error("RPC-API:GetTaskDetailListByTaskIds failed")
		return &carrierapipb.GetTaskDetailListResponse{Status: backend.ErrQueryTaskListByIds.ErrCode(), Msg: backend.ErrQueryTaskListByIds.Error()}, nil
	}

	log.Debugf("RPC-API:GetTaskDetailListByTaskIds succeed, taskIds len: {%d}, taskList len: {%d}", len(req.GetTaskIds()), len(tasks))
	return &carrierapipb.GetTaskDetailListResponse{
		Status: 0,
		Msg:    backend.OK,
		Tasks:  tasks,
	}, nil
}

func (svr *Server) GetTaskEventList(ctx context.Context, req *carrierapipb.GetTaskEventListRequest) (*carrierapipb.GetTaskEventListResponse, error) {

	if "" == strings.Trim(req.GetTaskId(), "") {
		return &carrierapipb.GetTaskEventListResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "required taskId"}, nil
	}

	events, err := svr.B.GetTaskEventList(req.GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:GetTaskEventList failed, taskId: {%s}", req.GetTaskId())
		errMsg := fmt.Sprintf("%s, taskId: {%s}", backend.ErrQueryTaskEventList.Error(), req.GetTaskId())
		return &carrierapipb.GetTaskEventListResponse{Status: backend.ErrQueryTaskEventList.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:GetTaskEventList succeed, taskId: {%s},  eventList len: {%d}", req.GetTaskId(), len(events))
	return &carrierapipb.GetTaskEventListResponse{
		Status:     0,
		Msg:        backend.OK,
		TaskEvents: events,
	}, nil
}

func (svr *Server) GetTaskEventListByTaskIds(ctx context.Context, req *carrierapipb.GetTaskEventListByTaskIdsRequest) (*carrierapipb.GetTaskEventListResponse, error) {

	if len(req.GetTaskIds()) == 0 {
		return &carrierapipb.GetTaskEventListResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "required taskIdx"}, nil
	}

	events, err := svr.B.GetTaskEventListByTaskIds(req.GetTaskIds())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:GetTaskEventListByTaskIds failed, taskId: {%v}", req.GetTaskIds())
		errMsg := fmt.Sprintf("%s, taskId: {%s}", backend.ErrQueryTaskEventList.Error(), req.GetTaskIds())
		return &carrierapipb.GetTaskEventListResponse{Status: backend.ErrQueryTaskEventList.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:GetTaskEventListByTaskIds succeed, taskIds: %v,  eventList len: {%d}", req.GetTaskIds(), len(events))
	return &carrierapipb.GetTaskEventListResponse{
		Status:     0,
		Msg:        backend.OK,
		TaskEvents: events,
	}, nil
}

func (svr *Server) PublishTaskDeclare(ctx context.Context, req *carrierapipb.PublishTaskDeclareRequest) (*carrierapipb.PublishTaskDeclareResponse, error) {

	identity, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishTaskDeclare failed, query local identity failed")
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}

	if "" == strings.Trim(req.GetTaskName(), "") {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check taskName failed, taskName is empty")
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require taskName"}, nil
	}
	if req.GetUserType() == commonconstantpb.UserType_User_Unknown {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check userType failed, wrong userType: {%s}", req.GetUserType().String())
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown userType"}, nil
	}
	if "" == strings.Trim(req.GetUser(), "") {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check user failed, user is empty")
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require user"}, nil
	}
	if nil == req.GetSender() {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check taskSender failed, taskSender is empty")
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require taskSender"}, nil
	}
	if req.GetSender().GetIdentityId() != identity.GetIdentityId() {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check taskSender failed, the taskSender is not current identity")
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "invalid taskSender"}, nil
	}
	if nil == req.GetAlgoSupplier() {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check algoSupplier failed, algoSupplier is empty")
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require algoSupplier"}, nil
	}
	if len(req.GetDataSuppliers()) == 0 {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check dataSuppliers failed, dataSuppliers is empty")
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require dataSuppliers"}, nil
	}
	if len(req.GetReceivers()) == 0 {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check receivers failed, receivers is empty")
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require receivers"}, nil
	}
	// about dataPolicy
	if len(req.GetDataPolicyTypes()) == 0 {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check DataPolicyType failed, dataPolicyTypes len is %d", len(req.GetDataPolicyTypes()))
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown dataPolicyTypes"}, nil
	}
	if len(req.GetDataPolicyOptions()) == 0 {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check DataPolicyOption failed, dataPolicyOptions len is %d", len(req.GetDataPolicyOptions()))
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require dataPolicyOptions"}, nil
	}
	if len(req.GetDataPolicyTypes()) != len(req.GetDataPolicyOptions()) || len(req.GetDataPolicyTypes()) != len(req.GetDataSuppliers()) {
		log.Errorf("RPC-API:PublishTaskDeclare failed, invalid dataPolicys len")
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "invalid dataPolicys len"}, nil
	}
	// about powerPolicy
	if len(req.GetPowerPolicyTypes()) == 0 {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check PowerPolicyType failed, powerPolicyTypes len is %d", len(req.GetPowerPolicyTypes()))
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown powerPolicyTypes"}, nil
	}
	if len(req.GetPowerPolicyOptions()) == 0 {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check PowerPolicyOption failed, powerPolicyOptions len is %d", len(req.GetPowerPolicyOptions()))
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require powerPolicyOptions"}, nil
	}
	if len(req.GetPowerPolicyTypes()) != len(req.GetPowerPolicyOptions()) {
		log.Errorf("RPC-API:PublishTaskDeclare failed, invalid powerPolicys len")
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "invalid powerPolicys len"}, nil
	}
	// about receiverPolicy
	if len(req.GetReceiverPolicyTypes()) == 0 {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check PowerPolicyType failed, receiverPolicyTypes len is %d", len(req.GetReceiverPolicyTypes()))
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown receiverPolicyTypes"}, nil
	}
	if len(req.GetReceiverPolicyOptions()) == 0 {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check PowerPolicyOption failed, receiverPolicyOptions len is %d", len(req.GetReceiverPolicyOptions()))
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require receiverPolicyOptions"}, nil
	}
	if len(req.GetReceiverPolicyTypes()) != len(req.GetReceiverPolicyOptions()) || len(req.GetReceiverPolicyTypes()) != len(req.GetReceivers()) {
		log.Errorf("RPC-API:PublishTaskDeclare failed, invalid receiverPolicys len")
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "invalid receiverPolicys len"}, nil
	}

	if len(req.GetDataFlowPolicyTypes()) == 0 {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check DataFlowPolicyType failed, DataFlowPolicyTypes len is %d", len(req.GetDataFlowPolicyTypes()))
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown dataFlowPolicyTypes"}, nil
	}
	if len(req.GetDataFlowPolicyOptions()) == 0 {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check DataFlowPolicyOption failed, DataFlowPolicyOptions len is %d", len(req.GetDataFlowPolicyOptions()))
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require dataFlowPolicyOptions"}, nil
	}
	if len(req.GetDataFlowPolicyTypes()) != len(req.GetDataFlowPolicyOptions()) {
		log.Errorf("RPC-API:PublishTaskDeclare failed, invalid dataFlowPolicys len")
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "invalid dataFlowPolicys len"}, nil
	}

	if req.GetOperationCost() == nil {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check OperationCost failed, OperationCost is empty")
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require operationCost"}, nil
	}
	if "" == strings.Trim(req.GetAlgorithmCode(), "") && "" == strings.Trim(req.GetMetaAlgorithmId(), "") {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check AlgorithmCode AND MetaAlgorithmId failed, AlgorithmCode AND MetaAlgorithmId is empty")
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require algorithmCode OR metaAlgorithmId"}, nil
	}
	// Maybe the contractExtraParams is empty,
	// Because contractExtraParams is not necessary.
	// Whether contractExtraParams has a value depends on the specific matching algorithmCode requirements.
	//
	//if "" == strings.Trim(req.GetAlgorithmCodeExtraParams(), "") {
	//	log.Errorf("RPC-API:PublishTaskDeclare failed, check ContractExtraParams failed, ContractExtraParams is empty")
	//	return &carrierapipb.PublishTaskDeclareResponse{ Status:  backend.ErrRequireParams.ErrCode(), Msg: "require contractExtraParams"}, nil
	//}
	if len(req.GetSign()) == 0 {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check sign failed, sign is empty")
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require sign"}, nil
	}

	taskMsg := types.NewTaskMessageFromRequest(req)
	/*from, err := signsuite.Sender(req.GetUserType(), taskMsg.Hash(), req.GetSign())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishTaskDeclare failed, cannot fetch sender from sign, userType: {%s}, user: {%s}",
			req.GetUserType().String(), req.GetUser())
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "cannot fetch sender from sign"}, nil
	}
	if from != req.GetUser() {
		log.WithError(err).Errorf("RPC-API:PublishTaskDeclare failed, sender from sign and user is not sameone, userType: {%s}, user: {%s}, sender of sign: {%s}",
			req.GetUserType().String(), req.GetUser(), from)
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "the user sign is invalid"}, nil
	}*/
	// Maybe the desc is empty,
	// Because desc is not necessary.
	//
	//if "" == strings.Trim(req.GetDesc(), "") {
	//	log.Errorf("RPC-API:PublishTaskDeclare failed, check Desc failed, Desc is empty")
	//	return &carrierapipb.PublishTaskDeclareResponse{ Status:  backend.ErrRequireParams.ErrCode(), Msg: "require desc"}, nil
	//}

	checkPartyIdCache := make(map[string]struct{}, 0)
	checkPartyIdCache[req.GetSender().GetPartyId()] = struct{}{}

	if _, ok := checkPartyIdCache[req.GetAlgoSupplier().GetPartyId()]; ok {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check partyId of algoSupplier failed, this partyId has alreay exist, partyId: {%s}",
			req.GetAlgoSupplier().GetPartyId())
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrPublishTaskMsg.ErrCode(), Msg: fmt.Sprintf("The partyId of the task participants cannot be repeated on algoSupplier, partyId: {%s}", req.GetAlgoSupplier().GetPartyId())}, nil
	}
	checkPartyIdCache[req.GetAlgoSupplier().GetPartyId()] = struct{}{}

	for _, v := range req.GetDataSuppliers() {
		if _, ok := checkPartyIdCache[v.GetPartyId()]; ok {
			log.Errorf("RPC-API:PublishTaskDeclare failed, check partyId of dataSuppliers failed, this partyId has alreay exist, partyId: {%s}",
				v.GetPartyId())
			return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrPublishTaskMsg.ErrCode(), Msg: fmt.Sprintf("The partyId of the task participants cannot be repeated on dataSuppliers, partyId: {%s}", v.GetPartyId())}, nil
		}
		checkPartyIdCache[v.GetPartyId()] = struct{}{}
	}

	// check partyId of powerSuppliers
	powerPartyIds, err := svr.B.GetPolicyEngine().FetchPowerPartyIdsFromPowerPolicy(req.GetPowerPolicyTypes(), req.GetPowerPolicyOptions())
	if nil != err {
		log.WithError(err).Errorf("not fetch partyIds from task powerPolicy")
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrPublishTaskMsg.ErrCode(), Msg: "not fetch partyIds from task powerPolicy"}, nil
	}
	for _, partyId := range powerPartyIds {
		if _, ok := checkPartyIdCache[partyId]; ok {
			log.Errorf("RPC-API:PublishTaskDeclare failed, check partyId of powerSuppliers failed, this partyId has alreay exist, partyId: {%s}",
				partyId)
			return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrPublishTaskMsg.ErrCode(), Msg: fmt.Sprintf("The partyId of the task participants cannot be repeated on powerSuppliers, partyId: {%s}", partyId)}, nil
		}
		checkPartyIdCache[partyId] = struct{}{}
	}

	// check partyId of receivers
	for i, v := range req.GetReceivers() {

		partyId, err := svr.B.GetPolicyEngine().FetchReceiverPartyIdByOptionFromReceiverPolicy(req.GetReceiverPolicyTypes()[i], req.GetReceiverPolicyOptions()[i])
		if nil != err {
			log.WithError(err).Errorf("RPC-API:PublishTaskDeclare failed, fetch partyId of receiverPolicy failed, index: {%d}, partyId of receiverPolicy: {%s}",
				i, partyId)
			return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrPublishTaskMsg.ErrCode(),
				Msg: fmt.Sprintf("fetch partyId of receiverPolicy failed, %s, index: {%d}, partyId: {%s}",
					err, i, v.GetPartyId())}, nil
		}

		if partyId != v.GetPartyId() {
			log.Errorf("RPC-API:PublishTaskDeclare failed, partyId of receiverPolicy and receiver is not same, index: {%d}, partyId of receiver: {%s}, partyId of receiverPolicy: {%s}",
				i, v.GetPartyId(), partyId)
			return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrPublishTaskMsg.ErrCode(),
				Msg: fmt.Sprintf("partyId of receiverPolicy and receiver is not same, index: {%d}, partyId of receiver: {%s}, partyId of receiverPolicy: {%s}",
					i, v.GetPartyId(), partyId)}, nil
		}

		if _, ok := checkPartyIdCache[v.GetPartyId()]; ok {
			log.Errorf("RPC-API:PublishTaskDeclare failed, check partyId of receiver failed, this partyId has alreay exist, partyId: {%s}",
				v.GetPartyId())
			return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrPublishTaskMsg.ErrCode(),
				Msg: fmt.Sprintf("The partyId of the task participants cannot be repeated on receivers, partyId: {%s}",
					v.GetPartyId())}, nil
		}
		checkPartyIdCache[v.GetPartyId()] = struct{}{}
	}

	// generate and store taskId
	taskId := taskMsg.GenTaskId()

	if err = svr.B.SendMsg(taskMsg); nil != err {
		log.WithError(err).Errorf("RPC-API:PublishTaskDeclare failed, send task msg failed, taskId: {%s}",
			taskId)
		errMsg := fmt.Sprintf("%s, taskId: {%s}", backend.ErrPublishTaskMsg.Error(), taskId)
		return &carrierapipb.PublishTaskDeclareResponse{Status: backend.ErrPublishTaskMsg.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:PublishTaskDeclare succeed, userType: {%s}, user: {%s}, publish taskId: {%s}", req.GetUserType(), req.GetUser(), taskId)
	return &carrierapipb.PublishTaskDeclareResponse{
		Status: 0,
		Msg:    backend.OK,
		TaskId: taskId,
	}, nil
}

func (svr *Server) TerminateTask(ctx context.Context, req *carrierapipb.TerminateTaskRequest) (*carriertypespb.SimpleResponse, error) {

	identity, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:TerminateTask failed, query local identity failed")
		return &carriertypespb.SimpleResponse{Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}

	if req.GetUserType() == commonconstantpb.UserType_User_Unknown {
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown userType"}, nil
	}
	if "" == req.GetUser() {
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require user"}, nil
	}
	if "" == req.GetTaskId() {
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require taskId"}, nil
	}
	if len(req.GetSign()) == 0 {
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require sign"}, nil
	}

	taskTerminateMsg := types.NewTaskTerminateMsg(req.GetUserType(), req.GetUser(), req.GetTaskId(), req.GetSign())
	from, err := signsuite.Sender(req.GetUserType(), taskTerminateMsg.Hash(), req.GetSign())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:TerminateTask failed, cannot fetch sender from sign, userType: {%s}, user: {%s}",
			req.GetUserType().String(), req.GetUser())
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "cannot fetch sender from sign"}, nil
	}
	if from != req.GetUser() {
		log.WithError(err).Errorf("RPC-API:TerminateTask failed, sender from sign and user is not sameone, userType: {%s}, user: {%s}, sender of sign: {%s}",
			req.GetUserType().String(), req.GetUser(), from)
		return &carriertypespb.SimpleResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "the user sign is invalid"}, nil
	}

	task, err := svr.B.GetLocalTask(req.GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:TerminateTask failed, query local task failed")
		return &carriertypespb.SimpleResponse{Status: backend.ErrTerminateTaskMsg.ErrCode(), Msg: "query local task failed"}, nil
	}

	// The sender of the terminatemsg must be the sender of the task
	if task.GetInformation().GetSender().GetIdentityId() != identity.GetIdentityId() {
		log.Errorf("RPC-API:TerminateTask failed, check taskSender failed, the taskSender is not current identity")
		return &carriertypespb.SimpleResponse{Status: backend.ErrTerminateTaskMsg.ErrCode(), Msg: "invalid taskSender"}, nil
	}

	// check user
	if task.GetInformation().GetUser() != req.GetUser() ||
		task.GetInformation().GetUserType() != req.GetUserType() {
		log.WithError(err).Errorf("terminate task user and publish task user must be same, taskId: {%s}",
			task.GetInformation().GetTaskId())
		return &carriertypespb.SimpleResponse{Status: backend.ErrTerminateTaskMsg.ErrCode(), Msg: fmt.Sprintf("terminate task user and publish task user must be same, taskId: {%s}",
			task.GetInformation().GetTaskId())}, nil
	}

	if err = svr.B.SendMsg(taskTerminateMsg); nil != err {
		log.WithError(err).Errorf("RPC-API:TerminateTask failed, send task terminate msg failed, taskId: {%s}",
			req.GetTaskId())

		errMsg := fmt.Sprintf("%s, taskId: {%s}", backend.ErrTerminateTaskMsg.Error(), req.GetTaskId())
		return &carriertypespb.SimpleResponse{Status: backend.ErrTerminateTaskMsg.ErrCode(), Msg: errMsg}, nil
	}
	log.Debugf("RPC-API:TerminateTask succeed, userType: {%s}, user: {%s}, taskId: {%s}", req.GetUserType(), req.GetUser(), req.GetTaskId())
	return &carriertypespb.SimpleResponse{
		Status: 0,
		Msg:    backend.OK,
	}, nil
}

func (svr *Server) GenerateObServerProxyWalletAddress(ctx context.Context, req *emptypb.Empty) (*carrierapipb.GenerateObServerProxyWalletAddressResponse, error) {

	wallet, err := svr.B.GenerateObServerProxyWalletAddress()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:GenerateObServerProxyWalletAddress failed")
		return &carrierapipb.GenerateObServerProxyWalletAddressResponse{Status: backend.ErrGenerateObserverProxyWallet.ErrCode(), Msg: backend.ErrGenerateObserverProxyWallet.Error()}, nil
	}
	log.Debugf("RPC-API:GenerateObServerProxyWalletAddress Succeed, wallet: {%s}", wallet)
	return &carrierapipb.GenerateObServerProxyWalletAddressResponse{
		Status:  0,
		Msg:     backend.OK,
		Address: wallet,
	}, nil
}

func (svr *Server) EstimateTaskGas(ctx context.Context, req *carrierapipb.EstimateTaskGasRequest) (*carrierapipb.EstimateTaskGasResponse, error) {
	gasLimit, gasPrice, err := svr.B.EstimateTaskGas(req.GetTaskSponsorAddress(), req.GetTkItems())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:EstimateTaskGas failed")
		return &carrierapipb.EstimateTaskGasResponse{Status: backend.ErrEstimateTaskGas.ErrCode(), Msg: backend.ErrEstimateTaskGas.Error()}, nil
	}
	log.Debugf("RPC-API:EstimateGas Succeed: gas {%d}", gasLimit)
	return &carrierapipb.EstimateTaskGasResponse{
		Status:   0,
		Msg:      backend.OK,
		GasLimit: gasLimit,
		GasPrice: gasPrice.Uint64(),
	}, nil
}
