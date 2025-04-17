package auth

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/sicko7947/sickocommon"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/models"
	auth_service "sicko-aio-2.0-client/proto/auth"
	"sicko-aio-2.0-client/utils/grpc"
)

func Login() (code int64, message string) {
	if communicator.DEV_ENV {
		return 200, "OK"
	}
	// Setup grpc connection and streams
	conn, err := grpc.GetAuthServerGrpcCon()
	if err != nil {
		grpc.ResetConn()
		return 409, "Error Activating Your Key"
	}
	defer conn.Close()

	authStreamClient := auth_service.NewAuthStreamClient(conn.Value())
	authStream, err := authStreamClient.Auth(context.Background())
	if err != nil {
		grpc.ResetConn()
		return 409, "Error Activating Your Key"
	}

	authStream.Send(&auth_service.StreamAuthRequest{
		Key:       communicator.Config.Settings.Key,
		Ipaddress: getIpAddress(),
		CpuId:     sickocommon.GetCpuID(),
		Timestamp: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
	})

	res, err := authStream.Recv()
	if err != nil {
		grpc.ResetConn()
		return 409, "Error Activating Your Key"
	}

	if res.GetCode() == 200 {
		hasLogin = true
	}
	return res.GetCode(), res.GetMessage()
}

func polling() *models.Error {
	// Setup grpc connection and streams
	conn, err := grpc.GetAuthServerGrpcCon()
	if err != nil {
		return &models.Error{Code: 409, Message: "Error getting connection"}
	}
	defer conn.Close()

	pollingStreamClient := auth_service.NewAuthStreamClient(conn.Value())
	pollingStream, err := pollingStreamClient.Polling(context.Background())
	if err != nil {
		return &models.Error{Code: 409, Message: "Error getting connection"}
	}

	pollingStream.Send(&auth_service.StreamPollingRequest{
		Key:       communicator.Config.Settings.Key,
		Ipaddress: getIpAddress(),
		CpuId:     sickocommon.GetCpuID(),
		Timestamp: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
	})
	res, err := pollingStream.Recv()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		os.Exit(499)
		return &models.Error{Code: 409, Message: "Error Getting Polling Response"}
	}

	// set token after successful login
	switch res.GetCode() {
	case 200:
		return nil
	default:
		os.Exit(499)
		return &models.Error{
			Code:    int(res.GetCode()),
			Message: res.GetMessage(),
		}
	}
}

func Deactivate() *models.Error {
	// Setup grpc connection and streams
	conn, err := grpc.GetAuthServerGrpcCon()
	if err != nil {
		return &models.Error{Code: 409, Message: "Error getting connection"}
	}
	defer conn.Close()

	authStreamClient := auth_service.NewAuthStreamClient(conn.Value())
	authStream, err := authStreamClient.Deactivate(context.Background())
	if err != nil {
		return &models.Error{Code: 409, Message: "Error getting connection"}
	}

	authStream.Send(&auth_service.StreamDeactivateRequest{
		Key:       communicator.Config.Settings.Key,
		Ipaddress: getIpAddress(),
		CpuId:     sickocommon.GetCpuID(),
		Timestamp: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
	})

	os.Exit(125)
	return nil
}
