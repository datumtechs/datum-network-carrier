package task

import (
	"context"
	"fmt"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

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

	if "" == strings.Trim(req.GetTaskId(), "") {
		return nil, backend.NewRpcBizErr(ErrGetNodeTaskEventList.Code, "require taskId")
	}

	events, err := svr.B.GetTaskEventList(req.GetTaskId())
	if nil != err {
		log.WithError(err).Errorf("RPC-API:QueryTaskEventList failed, taskId: {%s}", req.GetTaskId())
		errMsg := fmt.Sprintf(ErrGetNodeTaskEventList.Msg, req.GetTaskId())
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
		errMsg := fmt.Sprintf(ErrGetNodeTaskEventList.Msg, req.GetTaskIds())
		return nil, backend.NewRpcBizErr(ErrGetNodeTaskEventList.Code, errMsg)
	}
	log.Debugf("RPC-API:QueryTaskEventListByTaskIds succeed, taskId: {%v},  eventList len: {%d}", req.GetTaskIds(), len(events))
	return &pb.GetTaskEventListResponse{
		Status:        0,
		Msg:           backend.OK,
		TaskEventList: events,
	}, nil
}

func taskJsonReq (req *pb.PublishTaskDeclareRequest) string {

	sender := fmt.Sprintf(`{"partyId": "%s", "nodeName": "%s", "nodeId": "%s", "identityId": "%s"}`,
		req.GetSender().GetPartyId(), req.GetSender().GetNodeName(), req.GetSender().GetNodeId(),
		req.GetSender().GetIdentityId())

	algoSupplier := fmt.Sprintf(`{"partyId": "%s", "nodeName": "%s", "nodeId": "%s", "identityId": "%s"}`,
		req.GetAlgoSupplier().GetPartyId(), req.GetAlgoSupplier().GetNodeName(),
		req.GetAlgoSupplier().GetNodeId(), req.GetAlgoSupplier().GetIdentityId())

	dataSuppliers := "[]"
	dataSupplierArr := make([]string, len(req.GetDataSuppliers()))
	for i, dataSupplier := range req.GetDataSuppliers() {

		organization :=  fmt.Sprintf(`{"partyId": "%s", "nodeName": "%s", "nodeId": "%s", "identityId": "%s"}`,
			dataSupplier.GetOrganization().GetPartyId(), dataSupplier.GetOrganization().GetNodeName(),
			dataSupplier.GetOrganization().GetNodeId(), dataSupplier.GetOrganization().GetIdentityId())


		selectedColumns := "[]"
		selectedColumnArr := make([]string, len(dataSupplier.GetMetadataInfo().GetSelectedColumns()))
		for j, selectedColumn := range dataSupplier.GetMetadataInfo().GetSelectedColumns() {
			selectedColumnArr[j] = fmt.Sprint(selectedColumn)
		}
		if len(selectedColumnArr) != 0 {
			selectedColumns = "[" + strings.Join(selectedColumnArr, ",") + "]"
		}

		metadataInfo := fmt.Sprintf(`{"metadataId": "%s", "keyColumn": %d, "selectedColumns": %s}`,
			dataSupplier.GetMetadataInfo().GetMetadataId(), dataSupplier.GetMetadataInfo().GetKeyColumn(), selectedColumns)

		dataSupplierArr[i] = fmt.Sprintf(`{"organization": %s, "metadataInfo": %s}`,
			organization, metadataInfo)
	}
	if len(dataSupplierArr) != 0 {
		dataSuppliers = "[" + strings.Join(dataSupplierArr, ",") + "]"
	}

	powerPartyIds := "[]"
	powerPartyIdArr := make([]string, len(req.GetPowerPartyIds()))
	for i, partyId := range req.GetPowerPartyIds() {
		powerPartyIdArr[i] = `"` + partyId + `"`
	}
	if len(powerPartyIdArr) != 0 {
		powerPartyIds = "[" + strings.Join(powerPartyIdArr, ",") + "]"
	}

	receivers := "[]"
	receiverArr := make([]string, len(req.GetReceivers()))
	for i, receiver := range req.GetReceivers() {
		receiverArr[i] = fmt.Sprintf(`{"partyId": "%s", "nodeName": "%s", "nodeId": "%s", "identityId": "%s"}`,
			receiver.GetPartyId(), receiver.GetNodeName(),
			receiver.GetNodeId(), receiver.GetIdentityId())
	}
	if len(receiverArr) != 0 {
		receivers = "[" + strings.Join(receiverArr, ",") + "]"
	}

	operationCost := fmt.Sprintf(`{"memory": %d, "processor": %d, "bandwidth": %d, "duration": %d}`,
		req.GetOperationCost().GetMemory(), req.GetOperationCost().GetProcessor(), req.GetOperationCost().GetBandwidth(), req.GetOperationCost().GetDuration())

	return fmt.Sprintf(`{"taskName": "%s", "user": "%s", "userType": %d, "sender": %s, "algoSupplier": %s, "dataSuppliers": %s, "powerPartyIds": %s, "receivers": %s, "operationCost": %s, "calculateContractCode": "%s", "dataSplitContractCode": "%s", "contractExtraParams": "%s", "sign": "%s", "desc": "%s"}`,
		req.GetTaskName(), req.GetUser(), req.GetUserType(), sender, algoSupplier, dataSuppliers, powerPartyIds, receivers, operationCost, req.GetCalculateContractCode(), req.GetDataSplitContractCode(), req.GetContractExtraParams(), string(req.GetSign()), req.GetDesc())
}

func (svr *Server) PublishTaskDeclare(ctx context.Context, req *pb.PublishTaskDeclareRequest) (*pb.PublishTaskDeclareResponse, error) {

	//log.Debugf("Received Publish task req: %s", taskJsonReq(req))

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
	if "" == req.GetCalculateContractCode() {
		log.Errorf("RPC-API:PublishTaskDeclare failed, check CalculateContractCode failed, CalculateContractCode is empty")
		return nil, ErrReqCalculateContractCodeForPublishTask
	}

	identity, err := svr.B.GetNodeIdentity()
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishTaskDeclare failed, query local identity failed, can not publish task")
		return nil, ErrPublishTaskDeclare
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

	// add  dataSuppliers
	dataSuppliers := make([]*libtypes.TaskDataSupplier, len(req.GetDataSuppliers()))
	for i, v := range req.GetDataSuppliers() {

		if _, ok := checkPartyIdCache[v.GetOrganization().GetPartyId()]; ok {
			log.Errorf("RPC-API:PublishTaskDeclare failed, check partyId of dataSupplier failed, this partyId has alreay exist, partyId: {%s}",
				v.GetOrganization().GetPartyId())
			return nil, fmt.Errorf("The partyId of the task participants cannot be repeated on dataSuppliers, partyId: {%s}", v.GetOrganization().GetPartyId())
		}
		checkPartyIdCache[v.GetOrganization().GetPartyId()] = struct{}{}

		var (
			identityId string
			internalMetadata bool
		)
		if v.GetOrganization().GetIdentityId() != identity.GetIdentityId() { // Is not current identity
			identityId = v.GetOrganization().GetIdentityId()
		} else {  //  current identity
			internalMetadata, err = svr.B.IsInternalMetadata(v.GetMetadataInfo().GetMetadataId())
			if nil != err {
				log.WithError(err).Errorf("RPC-API:PublishTaskDeclare failed, check metadata whether internal metadata failed, identityId: {%s}, metadataId: {%s}",
					v.GetOrganization().GetIdentityId(), v.GetMetadataInfo().GetMetadataId())

				errMsg := fmt.Sprintf(ErrReqMetadataDetailForPublishTask.Msg,
					v.GetOrganization().GetIdentityId(), v.GetMetadataInfo().GetMetadataId())
				return nil, backend.NewRpcBizErr(ErrReqMetadataDetailForPublishTask.Code, errMsg)
			}
		}


		metadata, err := svr.B.GetMetadataDetail(identityId, v.GetMetadataInfo().GetMetadataId())
		if nil != err {
			log.WithError(err).Errorf("RPC-API:PublishTaskDeclare failed, query metadata of partner failed, identityId: {%s}, metadataId: {%s}",
				v.GetOrganization().GetIdentityId(), v.GetMetadataInfo().GetMetadataId())

			errMsg := fmt.Sprintf(ErrReqMetadataDetailForPublishTask.Msg,
				v.GetOrganization().GetIdentityId(), v.GetMetadataInfo().GetMetadataId())
			return nil, backend.NewRpcBizErr(ErrReqMetadataDetailForPublishTask.Code, errMsg)
		}

		var (
		 	keyColumn *libtypes.MetadataColumn
			selectedColumns []*libtypes.MetadataColumn
		)

		// handle publish metadata columns
		if !internalMetadata {
			colTmp := make(map[uint32]*libtypes.MetadataColumn, len(metadata.GetData().GetMetadataColumns()))
			for _, col := range metadata.GetData().GetMetadataColumns() {
				colTmp[col.GetCIndex()] = col
			}

			if col, ok := colTmp[v.GetMetadataInfo().GetKeyColumn()]; ok {
				keyColumn = col
			} else {
				errMsg := fmt.Sprintf(ErrReqMetadataByKeyColumn.Msg, v.GetOrganization().GetIdentityId(),
					v.GetMetadataInfo().GetMetadataId(), v.GetMetadataInfo().GetKeyColumn())
				log.Errorf("RPC-API:PublishTaskDeclare failed, check columns colIndex of keyColumn, this colIndex is not exist, partyId: {%s}, colIndex: {%d}",
					v.GetOrganization().GetPartyId(), v.GetMetadataInfo().GetKeyColumn())
				return nil, backend.NewRpcBizErr(ErrReqMetadataByKeyColumn.Code, errMsg)
			}

			selectedColumns = make([]*libtypes.MetadataColumn, len(v.GetMetadataInfo().GetSelectedColumns()))

			for j, colIndex := range v.GetMetadataInfo().GetSelectedColumns() {
				if col, ok := colTmp[colIndex]; ok {
					selectedColumns[j] = col
				} else {
					errMsg := fmt.Sprintf(ErrReqMetadataBySelectedColumn.Msg,
						v.GetOrganization().GetIdentityId(), v.GetMetadataInfo().GetMetadataId(), colIndex)
					log.Errorf("RPC-API:PublishTaskDeclare failed, check columns colIndex of selectedColumns, this colIndex is not exist, partyId: {%s}, colIndex: {%d}",
						v.GetOrganization().GetPartyId(), colIndex)
					return nil, backend.NewRpcBizErr(ErrReqMetadataBySelectedColumn.Code, errMsg)
				}
			}
		} else {
			// take it zero value
			keyColumn = &libtypes.MetadataColumn{}
			selectedColumns = make([]*libtypes.MetadataColumn, 0)
		}



		dataSuppliers[i] = &libtypes.TaskDataSupplier{
			Organization: &apicommonpb.TaskOrganization{
				PartyId:    v.GetOrganization().GetPartyId(),
				NodeName:   v.GetOrganization().GetNodeName(),
				NodeId:     v.GetOrganization().GetNodeId(),
				IdentityId: v.GetOrganization().GetIdentityId(),
			},
			MetadataId:      v.GetMetadataInfo().GetMetadataId(),
			MetadataName:    metadata.GetData().GetTableName(),
			KeyColumn:       keyColumn,
			SelectedColumns: selectedColumns,
		}
	}
	// add dataSuppliers
	taskMsg.Data.SetMetadataSupplierArr(dataSuppliers)

	for _, partyId := range req.GetPowerPartyIds() {
		if _, ok := checkPartyIdCache[partyId]; ok {
			log.Errorf("RPC-API:PublishTaskDeclare failed, check partyId of powerSupplier failed, this partyId has alreay exist, partyId: {%s}",
				partyId)
			return nil, fmt.Errorf("The partyId of the task participants cannot be repeated on powerPartyIds, partyId: {%s}", partyId)
		}
		checkPartyIdCache[partyId] = struct{}{}
	}

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

	err = svr.B.SendMsg(taskMsg)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:PublishTaskDeclare failed, send task msg failed, taskId: {%s}",
			taskId)
		errMsg := fmt.Sprintf(ErrSendTaskMsgByTaskId.Msg, taskId)
		return nil, backend.NewRpcBizErr(ErrSendTaskMsgByTaskId.Code, errMsg)
	}
	//log.Debugf("RPC-API:PublishTaskDeclare succeed, taskId: {%s}, taskMsg: %s", taskId, taskMsg.String())
	log.Debugf("RPC-API:PublishTaskDeclare succeed, taskId: {%s}", taskId)
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
		return nil, ErrTerminateTask
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

	err = svr.B.SendMsg(taskTerminateMsg)
	if nil != err {
		log.WithError(err).Errorf("RPC-API:TerminateTask failed, send task terminate msg failed, taskId: {%s}",
			req.GetTaskId())

		errMsg := fmt.Sprintf(ErrSendTaskMsg.Msg, req.GetTaskId())
		return nil, backend.NewRpcBizErr(ErrSendTaskMsg.Code, errMsg)
	}
	log.Debugf("RPC-API:TerminateTask succeed, taskId: {%s}", req.GetTaskId())
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
