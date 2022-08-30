package yarn

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/datumtechs/datum-network-carrier/carrier"
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	"github.com/datumtechs/datum-network-carrier/pb/common/constant"
	fighterapipb "github.com/datumtechs/datum-network-carrier/pb/fighter/api/data"
	"github.com/datumtechs/datum-network-carrier/pb/fighter/types"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"gotest.tools/assert"
	"io"
	"net"
	"os"
	"testing"
)

const (
	carrierAddress  = "127.0.0.1:8899"
	dataNodeAddress = "127.0.0.1:3344"
	savePath        = "./jobnode_service.go.bak"
)

type DataNode struct {
}

func TestServerDataNodeAndYarn(t *testing.T) {
	go func() {
		listener, err := net.Listen("tcp", dataNodeAddress)
		if err != nil {
			t.Fatalf("listener err %v", err)
		}
		t.Log(dataNodeAddress + " net.Listing...")
		grpcServer := grpc.NewServer()
		fighterapipb.RegisterDataProviderServer(grpcServer, &DataNode{})
		err = grpcServer.Serve(listener)
		if err != nil {
			t.Fatalf("DataNode grpcServer call Serve err %v", err)
		}
	}()
	listener, err := net.Listen("tcp", carrierAddress)
	if err != nil {
		t.Fatalf("listener err %v", err)
	}
	t.Log(carrierAddress + " net.Listing...")
	grpcServer := grpc.NewServer()
	carrierapipb.RegisterYarnServiceServer(grpcServer, &Server{B: &carrier.MockApiBackend{}})
	err = grpcServer.Serve(listener)
	if err != nil {
		t.Fatalf("Carrier grpcServer call Serve err %v", err)
	}
}

func TestDownloadResultData(t *testing.T) {
	conn, err := grpc.Dial(carrierAddress, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("net.Connect err: %v", err)
	}
	defer func() {
		assert.NilError(t, conn.Close())
	}()
	grpcClient := carrierapipb.NewYarnServiceClient(conn)
	stream, err := grpcClient.DownloadTaskResultData(context.Background(), &carrierapipb.DownloadTaskResultDataRequest{
		TaskId:  "task:0x26915db614c60bc03bf0f2239e5528813bc22c0d3e6848fd45a4fa8681a0192e",
		Options: map[string]string{"aaa": "bbb"},
	})
	if err != nil {
		t.Fatalf("Call DownloadData err: %v", err)
	}
	file, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		t.Fatalf("open file fail %v", err)
	}
	//defer func() {
	// 	assert.NilError(t, file.Close())
	//}()
	write := bufio.NewWriter(file)
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Conversations get stream err: %v", err)
		}
		if res.GetStatus() == constant.TaskStatus_Start {
			t.Log("this is status start")
		}
		if res.GetStatus() == constant.TaskStatus_Finished {
			t.Log("this is status Finished")
		}
		_, err = write.Write(res.GetContent())
		assert.NilError(t, err)
		err = write.Flush()
		assert.NilError(t, err)
	}
	defer func() {
		assert.NilError(t, file.Close())
		if err := os.Remove(savePath); err != nil {
			t.Fatalf("remove fail %v", err)
		}
	}()
}

func (dn *DataNode) DownloadData(request *fighterapipb.DownloadRequest, server fighterapipb.DataProvider_DownloadDataServer) error {
	dataPath := request.GetDataPath()
	options := request.GetOptions()
	_, err := json.Marshal(options)
	if err != nil {
		return fmt.Errorf("DownloadData json.Marshal options fail")
	}
	file, err := os.Open(dataPath)
	if err != nil {
		return err
	}
	defer file.Close()

	for {
		var chunk = make([]byte, 1024)
		n, err := file.Read(chunk)
		if err == io.EOF {
			_ = server.Send(&fighterapipb.DownloadReply{
				Data: &fighterapipb.DownloadReply_Status{
					Status: constant.TaskStatus_Finished,
				},
			})
			break
		}
		if err != nil {
			return err
		}
		err = server.Send(&fighterapipb.DownloadReply{
			Data: &fighterapipb.DownloadReply_Content{
				Content: chunk[:n],
			},
		})
		if err != nil {
			_ = server.Send(&fighterapipb.DownloadReply{
				Data: &fighterapipb.DownloadReply_Status{
					Status: constant.TaskStatus_Failed,
				},
			})
		}
	}
	return nil
}

func (dn *DataNode) GetStatus(ctx context.Context, empty *emptypb.Empty) (*fighterapipb.GetStatusReply, error) {
	//TODO implement me
	panic("implement me")
}

func (dn *DataNode) ListData(ctx context.Context, empty *emptypb.Empty) (*fighterapipb.ListDataReply, error) {
	//TODO implement me
	panic("implement me")
}

func (dn *DataNode) UploadData(server fighterapipb.DataProvider_UploadDataServer) error {
	//TODO implement me
	panic("implement me")
}

func (dn *DataNode) BatchUpload(server fighterapipb.DataProvider_BatchUploadServer) error {
	//TODO implement me
	panic("implement me")
}

func (dn *DataNode) DeleteData(ctx context.Context, request *fighterapipb.DownloadRequest) (*fighterapipb.UploadReply, error) {
	//TODO implement me
	panic("implement me")
}

func (dn *DataNode) HandleTaskReadyGo(ctx context.Context, req *types.TaskReadyGoReq) (*types.TaskReadyGoReply, error) {
	//TODO implement me
	panic("implement me")
}

func (dn *DataNode) HandleCancelTask(ctx context.Context, req *types.TaskCancelReq) (*types.TaskCancelReply, error) {
	//TODO implement me
	panic("implement me")
}
