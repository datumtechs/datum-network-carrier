package backend

import (
	"context"
	"errors"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
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
		return nil, NewRpcBizErr(ErrGetNodeTaskListStr)
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
		Msg:      OK,
		TaskList: arr,
	}, nil
}

func (svr *TaskServiceServer) GetTaskEventList(ctx context.Context, req *pb.GetTaskEventListRequest) (*pb.GetTaskEventListResponse, error) {

	events, err := svr.B.GetTaskEventList(req.TaskId)
	if nil != err {
		return nil, NewRpcBizErr(ErrGetNodeTaskEventListStr)
	}

	return &pb.GetTaskEventListResponse{
		Status:        0,
		Msg:           OK,
		TaskEventList: types.ConvertTaskEventArrToPB(events),
	}, nil
}

func (svr *TaskServiceServer) PublishTaskDeclare(ctx context.Context, req *pb.PublishTaskDeclareRequest) (*pb.PublishTaskDeclareResponse, error) {
	if req == nil || req.Owner == nil {
		return nil, errors.New("required owner")
	}
	if req.Owner.MetaDataInfo == nil {
		return nil, errors.New("required metadataInfo")
	}
	if req.Owner.MemberInfo == nil {
		return nil, errors.New("required memberInfo of owner")
	}
	if req.OperationCost == nil {
		return nil, errors.New("require operationCost")
	}
	if req.Receivers == nil {
		return nil, errors.New("require receivers")
	}
	if req.Partners == nil {
		return nil, errors.New("require partners")
	}
	taskMsg := types.NewTaskMessageFromRequest(req)
	/*taskMsg.Data.TaskName = req.TaskName
	taskMsg.Data.CreateAt = uint64(time.Now().UnixNano())
	taskMsg.Data.Owner.Name = req.Owner.MemberInfo.Name
	taskMsg.Data.Owner.NodeId = req.Owner.MemberInfo.NodeId
	taskMsg.Data.Owner.IdentityId = req.Owner.MemberInfo.IdentityId
	taskMsg.Data.Owner.MetaData.ColumnIndexList = req.Owner.MetaDataInfo.ColumnIndexList
	taskMsg.Data.Owner.MetaData.MetaDataId = req.Owner.MetaDataInfo.MetaDataId*/

	partners := make([]*types.TaskSupplier, len(req.Partners))
	for i, v := range req.Partners {
		partner := &types.TaskSupplier{
			NodeAlias: &types.NodeAlias{
				Name:       v.MemberInfo.Name,
				NodeId:     v.MemberInfo.NodeId,
				IdentityId: v.MemberInfo.IdentityId,
			},
			MetaData: &types.SupplierMetaData{
				MetaDataId:      v.MetaDataInfo.MetaDataId,
				ColumnIndexList: v.MetaDataInfo.ColumnIndexList,
			},
		}
		partners[i] = partner
	}
	taskMsg.Data.Partners = partners

	receivers := make([]*types.TaskResultReceiver, len(req.Receivers))
	for i, v := range req.Receivers {

		providers := make([]*types.NodeAlias, len(v.Providers))
		for j, val := range v.Providers {
			provider := &types.NodeAlias{
				Name:       val.Name,
				NodeId:     val.NodeId,
				IdentityId: val.IdentityId,
			}
			providers[j] = provider
		}

		receiver := &types.TaskResultReceiver{
			NodeAlias: &types.NodeAlias{
				Name:       v.MemberInfo.Name,
				NodeId:     v.MemberInfo.NodeId,
				IdentityId: v.MemberInfo.IdentityId,
			},
			Providers: providers,
		}

		receivers[i] = receiver
	}
	taskMsg.Data.Receivers = receivers

	/*taskMsg.Data.CalculateContractCode = req.CalculateContractcode
	taskMsg.Data.DataSplitContractCode = req.DatasplitContractcode
	taskMsg.Data.OperationCost = &types.TaskOperationCost{
		Processor: req.OperationCost.CostProcessor,
		Mem:       req.OperationCost.CostMem,
		Bandwidth: req.OperationCost.CostBandwidth,
		Duration:  req.OperationCost.Duration,
	}*/
	taskId := taskMsg.GetTaskId()

	err := svr.B.SendMsg(taskMsg)
	if nil != err {
		return nil, NewRpcBizErr(ErrSendTaskMsgStr)
	}
	return &pb.PublishTaskDeclareResponse{
		Status: 0,
		Msg:    OK,
		TaskId: taskId,
	}, nil
}
