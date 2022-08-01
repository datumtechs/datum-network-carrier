package yarn

import (
	"bufio"
	"fmt"
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"github.com/datumtechs/datum-network-carrier/pb/common/constant"
	fighterapipb "github.com/datumtechs/datum-network-carrier/pb/fighter/api/data"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"net"
	"os"
	"testing"
)

func startDownloadTaskResultDataService() {
	address := "192.168.21.143:8899"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Println("listener err ", err)
	}
	log.Println(address + " net.Listing...")
	grpcServer := grpc.NewServer()
	carrierapipb.RegisterYarnServiceServer(grpcServer, &YarnServiceServer{})
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Println("grpcServer call Serve err ", err)
	}
}
func TestServer_DownloadTaskResultData(t *testing.T) {
	startDownloadTaskResultDataService()
}

type YarnServiceServer struct{}

func (s *YarnServiceServer) DownloadTaskResultData(req *carrierapipb.DownloadTaskResultDataRequest, server carrierapipb.YarnService_DownloadTaskResultDataServer) error {
	log.Info("taskId is " + req.GetTaskId())
	dataNodeIp := "192.168.10.154"
	dataNodePort := 8700
	dataPath := "/home/user1/data/data_root/insurance_predict_partyB_20220628-105125.csv"
	var dataServerAddress = fmt.Sprintf("%s:%d", dataNodeIp, dataNodePort)
	log.Info("dataServerAddress is ", dataServerAddress)
	conn, err := grpc.Dial(dataServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("net.Connect err: %v", err)
	}
	defer conn.Close()
	grpcClient := fighterapipb.NewDataProviderClient(conn)
	stream, err := grpcClient.DownloadData(context.Background(), &fighterapipb.DownloadRequest{
		DataPath: dataPath,
		//Options:  map[string]string{"compress": "zip"},
	})
	if err != nil {
		log.Fatalf("Call DownloadData err: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Conversations get stream err: %v", err)
		}
		if res.GetStatus() == constant.TaskStatus_Start {
			err = server.Send(&carrierapipb.DownloadTaskResultDataResponse{
				Data: &carrierapipb.DownloadTaskResultDataResponse_Content{
					Content: res.GetContent(),
				},
			})
		} else {
			err = server.Send(&carrierapipb.DownloadTaskResultDataResponse{
				Data: &carrierapipb.DownloadTaskResultDataResponse_Status{
					Status: res.GetStatus(),
				},
			})
		}
		if err != nil {
			return err
		}
	}
	return nil
}
func TestServer_MockDownloadTaskResultDataClient(t *testing.T) {
	mockDownloadTaskResultDataClient()
}
func mockDownloadTaskResultDataClient() {
	var dataServerAddress = "192.168.21.143:8899"
	conn, err := grpc.Dial(dataServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("net.Connect err: %v", err)
	}
	defer conn.Close()
	grpcClient := carrierapipb.NewYarnServiceClient(conn)
	stream, err := grpcClient.DownloadTaskResultData(context.Background(), &carrierapipb.DownloadTaskResultDataRequest{
		TaskId: "this is TaskId,123456",
	})
	if err != nil {
		log.Fatalf("Call DownloadData err: %v", err)
	}
	filePath := "./golang.csv"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("open file fail", err)
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Conversations get stream err: %v", err)
		}
		if res.GetStatus() == constant.TaskStatus_Start {
			log.Println("this is status start")
		}
		if res.GetStatus() == constant.TaskStatus_Finished {
			log.Println("this is status Finished")
		}
		write.Write(res.GetContent())
		write.Flush()
	}
}
func (s *YarnServiceServer) GetNodeInfo(ctx context.Context, empty *emptypb.Empty) (*carrierapipb.GetNodeInfoResponse, error) {
	panic("implement me")
}
func (s *YarnServiceServer) SetSeedNode(ctx context.Context, request *carrierapipb.SetSeedNodeRequest) (*carrierapipb.SetSeedNodeResponse, error) {
	panic("implement me")
}
func (s *YarnServiceServer) DeleteSeedNode(ctx context.Context, request *carrierapipb.DeleteSeedNodeRequest) (*carriertypespb.SimpleResponse, error) {
	panic("implement me")
}
func (s *YarnServiceServer) GetSeedNodeList(ctx context.Context, empty *emptypb.Empty) (*carrierapipb.GetSeedNodeListResponse, error) {
	panic("implement me")
}
func (s *YarnServiceServer) SetDataNode(ctx context.Context, request *carrierapipb.SetDataNodeRequest) (*carrierapipb.SetDataNodeResponse, error) {
	panic("implement me")
}
func (s *YarnServiceServer) UpdateDataNode(ctx context.Context, request *carrierapipb.UpdateDataNodeRequest) (*carrierapipb.SetDataNodeResponse, error) {
	panic("implement me")
}
func (s *YarnServiceServer) DeleteDataNode(ctx context.Context, request *carrierapipb.DeleteRegisteredNodeRequest) (*carriertypespb.SimpleResponse, error) {
	panic("implement me")
}
func (s *YarnServiceServer) GetDataNodeList(ctx context.Context, empty *emptypb.Empty) (*carrierapipb.GetRegisteredNodeListResponse, error) {
	panic("implement me")
}
func (s *YarnServiceServer) SetJobNode(ctx context.Context, request *carrierapipb.SetJobNodeRequest) (*carrierapipb.SetJobNodeResponse, error) {
	panic("implement me")
}
func (s *YarnServiceServer) UpdateJobNode(ctx context.Context, request *carrierapipb.UpdateJobNodeRequest) (*carrierapipb.SetJobNodeResponse, error) {
	panic("implement me")
}
func (s *YarnServiceServer) DeleteJobNode(ctx context.Context, request *carrierapipb.DeleteRegisteredNodeRequest) (*carriertypespb.SimpleResponse, error) {
	panic("implement me")
}
func (s *YarnServiceServer) GetJobNodeList(ctx context.Context, empty *emptypb.Empty) (*carrierapipb.GetRegisteredNodeListResponse, error) {
	panic("implement me")
}
func (s *YarnServiceServer) ReportTaskEvent(ctx context.Context, request *carrierapipb.ReportTaskEventRequest) (*carriertypespb.SimpleResponse, error) {
	panic("implement me")
}
func (s *YarnServiceServer) ReportTaskResourceUsage(ctx context.Context, request *carrierapipb.ReportTaskResourceUsageRequest) (*carriertypespb.SimpleResponse, error) {
	panic("implement me")
}
func (s *YarnServiceServer) ReportUpFileSummary(ctx context.Context, request *carrierapipb.ReportUpFileSummaryRequest) (*carriertypespb.SimpleResponse, error) {
	panic("implement me")
}
func (s *YarnServiceServer) ReportTaskResultFileSummary(ctx context.Context, request *carrierapipb.ReportTaskResultFileSummaryRequest) (*carriertypespb.SimpleResponse, error) {
	panic("implement me")
}
func (s *YarnServiceServer) QueryAvailableDataNode(ctx context.Context, request *carrierapipb.QueryAvailableDataNodeRequest) (*carrierapipb.QueryAvailableDataNodeResponse, error) {
	panic("implement me")
}
func (s *YarnServiceServer) QueryFilePosition(ctx context.Context, request *carrierapipb.QueryFilePositionRequest) (*carrierapipb.QueryFilePositionResponse, error) {
	panic("implement me")
}
func (s *YarnServiceServer) GetTaskResultFileSummary(ctx context.Context, request *carrierapipb.GetTaskResultFileSummaryRequest) (*carrierapipb.GetTaskResultFileSummaryResponse, error) {
	panic("implement me")
}
func (s *YarnServiceServer) GetTaskResultFileSummaryList(ctx context.Context, empty *emptypb.Empty) (*carrierapipb.GetTaskResultFileSummaryListResponse, error) {
	panic("implement me")
}
func (s *YarnServiceServer) GenerateObServerProxyWalletAddress(ctx context.Context, empty *emptypb.Empty) (*carrierapipb.GenerateObServerProxyWalletAddressResponse, error) {
	panic("implement me")
}
