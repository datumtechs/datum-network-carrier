package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/datumtechs/datum-network-carrier/rpc/backend"
	"github.com/datumtechs/datum-network-carrier/rpc/backend/task"
	"github.com/datumtechs/datum-network-carrier/signsuite"
	"github.com/datumtechs/datum-network-carrier/types"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

func (svr *Server) PublishWorkFlowDeclare(ctx context.Context, req *carrierapipb.PublishWorkFlowDeclareRequest) (*carrierapipb.PublishWorkFlowDeclareResponse, error) {
	identity, err := svr.B.GetNodeIdentity()
	reqStr, _ := json.Marshal(req)
	log.Infof("PublishWorkFlowDeclare req is : %s", reqStr)
	//  check policy include ring
	if checkWorkflowTaskListReferTo(req) {
		log.Errorf("RPC-API:PublishWorkFlowDeclare failed, checkWorkflowTaskListReferTo return is true,{%s}", req.Policy)
		return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "taskList refer to exits ring"}, nil
	}
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishWorkFlowDeclare failed, query local identity failed, workflowName {%s}", req.GetWorkflowName())
		return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrQueryNodeIdentity.ErrCode(), Msg: backend.ErrQueryNodeIdentity.Error()}, nil
	}
	if req.GetPolicyType() == commonconstantpb.WorkFlowPolicyType_Unknown_Policy {
		log.Errorf("RPC-API:PublishWorkFlowDeclare failed,workflow policy type unknown, {%s},workflow name {%s}", req.GetPolicyType().String(), req.GetWorkflowName())
		return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown policy type"}, nil
	}
	if req.GetPolicy() == "" {
		log.Errorf("RPC-API:PublishWorkFlowDeclare failed, check taskName failed, workflow policy is empty, workflowName {%s}", req.GetWorkflowName())
		return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require policy"}, nil
	}
	for _, v := range req.GetTaskList() {
		if "" == strings.Trim(v.GetTaskName(), "") {
			log.Errorf("RPC-API:PublishWorkFlowDeclare failed, check taskName failed, taskName is empty")
			return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require taskName"}, nil
		}
		if v.GetUserType() == commonconstantpb.UserType_User_Unknown {
			log.Errorf("RPC-API:PublishWorkFlowDeclare failed, check userType failed, wrong userType: {%s}, taskName: {%s}", v.GetUserType().String(), v.GetTaskName())
			return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown userType"}, nil
		}
		if "" == strings.Trim(v.GetUser(), "") {
			log.Errorf("RPC-API:PublishWorkFlowDeclare failed, check user failed, user is empty,taskName: {%s}", v.GetTaskName())
			return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require user"}, nil
		}
		if nil == v.GetSender() {
			log.Errorf("RPC-API:PublishWorkFlowDeclare failed, check taskSender failed, taskSender is empty,taskName: {%s}", v.GetTaskName())
			return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require taskSender"}, nil
		}
		if v.GetSender().GetIdentityId() != identity.GetIdentityId() {
			log.Errorf("RPC-API:PublishWorkFlowDeclare failed, check taskSender failed, the taskSender is not current identity,taskName: {%s}", v.GetTaskName())
			return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "invalid taskSender"}, nil
		}
		if nil == v.GetAlgoSupplier() {
			log.Errorf("RPC-API:PublishWorkFlowDeclare failed, check algoSupplier failed, algoSupplier is empty,taskName: {%s}", v.GetTaskName())
			return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require algoSupplier"}, nil
		}
		if len(v.GetDataSuppliers()) == 0 {
			log.Errorf("RPC-API:PublishWorkFlowDeclare failed, check dataSuppliers failed, dataSuppliers is empty,taskName: {%s}", v.GetTaskName())
			return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require dataSuppliers"}, nil
		}
		if len(v.GetReceivers()) == 0 {
			log.Errorf("RPC-API:PublishWorkFlowDeclare failed, check receivers failed, receivers is empty,taskName: {%s}", v.GetTaskName())
			return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require receivers"}, nil
		}
		// about powerPolicy
		if len(v.GetPowerPolicyTypes()) == 0 {
			log.Errorf("RPC-API:PublishWorkFlowDeclare failed, check PowerPolicyType failed, powerPolicyTypes len is %d,taskName: {%s}", len(v.GetPowerPolicyTypes()), v.GetTaskName())
			return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown powerPolicyTypes"}, nil
		}
		if len(v.GetPowerPolicyOptions()) == 0 {
			log.Errorf("RPC-API:PublishWorkFlowDeclare failed, check PowerPolicyOption failed, powerPolicyOptions len is %d,taskName: {%s}", len(v.GetPowerPolicyOptions()), v.GetTaskName())
			return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require powerPolicyOptions"}, nil
		}
		if len(v.GetPowerPolicyTypes()) != len(v.GetPowerPolicyOptions()) {
			log.Errorf("RPC-API:PublishWorkFlowDeclare failed, invalid powerPolicys len,taskName: {%s}", v.GetTaskName())
			return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "invalid powerPolicys len"}, nil
		}
		// about receiverPolicy
		if len(v.GetReceiverPolicyTypes()) == 0 {
			log.Errorf("RPC-API:PublishWorkFlowDeclare failed, check PowerPolicyType failed, receiverPolicyTypes len is %d,taskName: {%s}", len(v.GetReceiverPolicyTypes()), v.GetTaskName())
			return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown receiverPolicyTypes"}, nil
		}
		if len(v.GetReceiverPolicyOptions()) == 0 {
			log.Errorf("RPC-API:PublishWorkFlowDeclare failed, check PowerPolicyOption failed, receiverPolicyOptions len is %d,taskName: {%s}", len(v.GetReceiverPolicyOptions()), v.GetTaskName())
			return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require receiverPolicyOptions"}, nil
		}
		if len(v.GetReceiverPolicyTypes()) != len(v.GetReceiverPolicyOptions()) || len(v.GetReceiverPolicyTypes()) != len(v.GetReceivers()) {
			log.Errorf("RPC-API:PublishWorkFlowDeclare failed, invalid receiverPolicys len,taskName: {%s}", v.GetTaskName())
			return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "invalid receiverPolicys len"}, nil
		}

		if len(v.GetDataFlowPolicyTypes()) == 0 {
			log.Errorf("RPC-API:PublishWorkFlowDeclare failed, check DataFlowPolicyType failed, DataFlowPolicyTypes len is %d,taskName: {%s}", len(v.GetDataFlowPolicyTypes()), v.GetTaskName())
			return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "unknown dataFlowPolicyTypes"}, nil
		}
		if len(v.GetDataFlowPolicyOptions()) == 0 {
			log.Errorf("RPC-API:PublishWorkFlowDeclare failed, check DataFlowPolicyOption failed, DataFlowPolicyOptions len is %d,taskName: {%s}", len(v.GetDataFlowPolicyOptions()), v.GetTaskName())
			return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require dataFlowPolicyOptions"}, nil
		}
		if len(v.GetDataFlowPolicyTypes()) != len(v.GetDataFlowPolicyOptions()) {
			log.Errorf("RPC-API:PublishWorkFlowDeclare failed, invalid dataFlowPolicys len,taskName: {%s}", v.GetTaskName())
			return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "invalid dataFlowPolicys len"}, nil
		}

		if v.GetOperationCost() == nil {
			log.Errorf("RPC-API:PublishWorkFlowDeclare failed, check OperationCost failed, OperationCost is empty,taskName: {%s}", v.GetTaskName())
			return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require operationCost"}, nil
		}
		if "" == strings.Trim(v.GetAlgorithmCode(), "") && "" == strings.Trim(v.GetMetaAlgorithmId(), "") {
			log.Errorf("RPC-API:PublishWorkFlowDeclare failed, check AlgorithmCode AND MetaAlgorithmId failed, AlgorithmCode AND MetaAlgorithmId is empty,taskName: {%s}", v.GetTaskName())
			return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require algorithmCode OR metaAlgorithmId"}, nil
		}
		err = task.CheckTaskReqPowerPolicy(v.GetPowerPolicyTypes(), v.GetPowerPolicyOptions())
		if err != nil {
			log.Errorf("RPC-API:PublishWorkFlowDeclare failed,CheckTaskReqPowerPolicy,taskName: {%s}", v.GetTaskName())
			return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: err.Error()}, nil
		}
	}
	if len(req.GetSign()) == 0 {
		log.Errorf("RPC-API:PublishWorkFlowDeclare failed, check sign failed, sign is empty, workflowName {%s}", req.GetWorkflowName())
		return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "require sign"}, nil
	}
	workflowMsg := types.NewWorkFlowMessageFromRequest(req)
	from, _, err := signsuite.Sender(req.GetUserType(), workflowMsg.Hash(), req.GetSign())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishWorkFlowDeclare failed, cannot fetch sender from sign, userType: {%s}, user: {%s}",
			req.GetUserType().String(), req.GetUser())
		return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "cannot fetch sender from sign"}, nil
	}
	if from != req.GetUser() {
		log.WithError(err).Errorf("RPC-API:PublishWorkFlowDeclare failed, sender from sign and user is not sameone, userType: {%s}, user: {%s}, sender of sign: {%s}",
			req.GetUserType().String(), req.GetUser(), from)
		return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrRequireParams.ErrCode(), Msg: "the user sign is invalid"}, nil
	}
	workflowId := workflowMsg.GenWorkflowId()
	if err = svr.B.SendMsg(workflowMsg); nil != err {
		log.WithError(err).Errorf("RPC-API:PublishWorkFlowDeclareResponse failed, send task msg failed, workflowId: {%s}", workflowId)
		errMsg := fmt.Sprintf("%s, workflowId: {%s}", backend.ErrPublishTaskMsg.Error(), workflowId)
		return &carrierapipb.PublishWorkFlowDeclareResponse{Status: backend.ErrPublishTaskMsg.ErrCode(), Msg: errMsg}, nil
	}
	log.Infof("PublishWorkFlowDeclare workflowName %s,workflowId %s", workflowMsg.Data.GetWorkflowName(), workflowId)
	return &carrierapipb.PublishWorkFlowDeclareResponse{
		Status: 0,
		Msg:    backend.OK,
		Id:     workflowId,
	}, nil
}

func (svr *Server) QueryWorkFlowStatus(ctx context.Context, req *carrierapipb.QueryWorkStatusRequest) (*carrierapipb.QueryWorkStatusResponse, error) {
	return svr.B.GetWorkflowStatus(req.GetWorkflowIds())
}

func (svr *Server) QueryAllWorkFlowDetails(ctx context.Context, req *emptypb.Empty) (*carrierapipb.QueryWorkflowDetailsResponse, error) {
	return svr.B.GetAllWorkflowDetails()
}

func checkWorkflowTaskListReferTo(req *carrierapipb.PublishWorkFlowDeclareRequest) bool {
	policy := req.GetPolicy()
	var wp *types.WorkflowPolicy
	if err := json.Unmarshal([]byte(policy), &wp); err != nil {
		log.Errorf("checkWorkflowTaskListReferTo Unmarshal fail, policy is %s", policy)
		return false
	}
	referToGraph := make(map[string][]string, 0)
	for _, v := range *wp {
		referTo := make([]string, 0)
		for _, s := range v.Reference {
			referTo = append(referTo, s.Target)
		}
		if _, ok := referToGraph[v.Origin]; ok {
			log.Warnf("taskName %s repet,please check", v.Origin)
			return true
		}
		referToGraph[v.Origin] = referTo
	}
	dependencyOrder := directedAcyclicGraphTopologicalSort(referToGraph)
	if len(dependencyOrder) == 0 {
		log.Warnf("there is a circular dependency on the workflow policy, please check")
		return true
	}
	//updateTaskListOrder := make([]*carrierapipb.PublishTaskDeclareRequest, 0)
	//for index := range dependencyOrder {
	//	taskNameOrder := dependencyOrder[len(dependencyOrder)-index-1]
	//	for _, taskDetail := range req.GetTaskList() {
	//		if taskNameOrder == taskDetail.GetTaskName() {
	//			updateTaskListOrder = append(updateTaskListOrder, taskDetail)
	//		}
	//	}
	//}
	//log.Debugf("updateTaskListOrder result:%v", updateTaskListOrder)
	//req.TaskList = updateTaskListOrder
	return false
}

func directedAcyclicGraphTopologicalSort(graph map[string][]string) []string {
	/*
		graph := make(map[string][]string, 0)
		graph["aTask"] = []string{"bTask", "cTask", "eTask"}
		graph["bTask"] = []string{"dTask"}
		graph["cTask"] = []string{"dTask"}
		graph["dTask"] = []string{}
		graph["eTask"] = []string{"cTask", "dTask"}

		return -> [aTask eTask cTask bTask dTask]
	*/
	// Initialize the in-degree of all vertices as 0
	inDegrees := make(map[string]uint32, 0)
	for k, _ := range graph {
		inDegrees[k] = 0
	}

	for _, v := range graph {
		for _, n := range v {
			//Calculate the in-degree of each vertex
			inDegrees[n] += 1
		}
	}
	//Filter vertices with in-degree 0
	inDegreesEqualZero := make([]string, 0)
	for k, v := range inDegrees {
		if v == 0 {
			inDegreesEqualZero = append(inDegreesEqualZero, k)
		}
	}
	var seq []string
	for {
		if len(inDegreesEqualZero) == 0 {
			break
		}
		// Delete from last by default
		theLastOne := inDegreesEqualZero[len(inDegreesEqualZero)-1]
		seq = append(seq, theLastOne)
		inDegreesEqualZero = inDegreesEqualZero[:len(inDegreesEqualZero)-1]
		for _, v := range graph[theLastOne] {
			//  remove all references to it
			inDegrees[v] -= 1
			if inDegrees[v] == 0 {
				//  Filter again vertices with in-degree 0
				inDegreesEqualZero = append(inDegreesEqualZero, v)
			}
		}
	}
	if len(seq) == len(inDegrees) {
		// If there are vertices with non-zero in-degree after the loop ends, there is a cycle in the graph
		return seq
	} else {
		log.Warnf("directedAcyclicGraphTopologicalSort include ring,please check %v", graph)
		return nil
	}
}
