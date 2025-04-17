package successHandler

import (
	"context"
	"io"
	"time"

	"github.com/huandu/go-clone"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/models"
	auth_service "sicko-aio-2.0-client/proto/auth"
	"sicko-aio-2.0-client/tasks/defineder"
	"sicko-aio-2.0-client/utils/grpc"
)

func RetrieveSuccess() ([]*auth_service.StreamRetrieveSuccessItemsResponse_SuccessItem, error) {
	// Setup grpc connection and streams
	conn, err := grpc.GetAuthServerGrpcCon()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	streamClient := auth_service.NewAuthStreamClient(conn.Value())
	retrieveSuccessStream, err := streamClient.RetrieveSuccess(context.Background())
	if err != nil {
		return nil, err
	}

	// send request
	err = retrieveSuccessStream.Send(&auth_service.StreamRetrieveSuccessItemsRequest{
		KeyId: communicator.Config.Settings.Key,
	})
	if err != nil && err != io.EOF {
		return nil, err
	}

	resCh := make(chan []*auth_service.StreamRetrieveSuccessItemsResponse_SuccessItem, 1)
	defer close(resCh)
	go func() {
		for {
			if retrieveSuccessStream.Context().Err() != nil {
				return
			}

			res, err := retrieveSuccessStream.Recv()
			if err != nil {
				return
			}

			resCh <- res.GetSuccessItems()
			retrieveSuccessStream.Context().Done()
		}
	}()

	return <-resCh, nil
}

func HandlerSuccess(taskGroupSetting *models.TaskGroupSetting, worker *models.TaskWorker) error {
	i := communicator.SuccessCountGmap.Get(string(worker.GroupID))
	communicator.SuccessCountGmap.Set(string(worker.GroupID), i+1)

	newWorker := clone.Clone(worker).(*models.TaskWorker)

	// Setup grpc connection and streams
	conn, err := grpc.GetAuthServerGrpcCon()
	if err != nil {
		return err
	}
	defer conn.Close()

	streamClient := auth_service.NewAuthStreamClient(conn.Value())
	successHandlerStream, err := streamClient.HandleSuccessCheckout(context.Background())
	if err != nil {
		return err
	}

	// setup payload
	payload := &auth_service.StreamHandleSuccessCheckoutRequest{
		KeyId: communicator.Config.Settings.Key,
		Setup: &auth_service.SuccessSetup{
			Category:        string(taskGroupSetting.Category),
			Region:          taskGroupSetting.Country,
			TaskType:        taskGroupSetting.TaskType,
			MonitorMode:     defineder.NikeMonitorGRPC,
			Timestamp:       time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
			UsePsychoCookie: false,
		},
		Product: &auth_service.SuccessProduct{
			MerchGroup:         taskGroupSetting.MerchGroup,
			ProductSku:         newWorker.Product.StyleColor,
			ProductName:        newWorker.Product.ProductName,
			ProductDescription: newWorker.Product.ProductDescription,
			Size:               newWorker.Product.Size,
			Price:              newWorker.Product.Price,
			Quantity:           int64(newWorker.Quantity),
			ImageUrl:           newWorker.Product.ImageURL,
			RedirectUrl:        newWorker.TaskInfo.RedirectURL,
		},
	}

	if newWorker.TaskInfo != nil && newWorker.TaskInfo.Profile != nil {
		payload.Product.OrderNumber = newWorker.TaskInfo.OrderID
		payload.Product.ProfileName = newWorker.TaskInfo.Profile.ProfileName
		payload.Product.Email = newWorker.TaskInfo.Email
		payload.Product.DiscountCode = newWorker.TaskInfo.Discount

		if newWorker.TaskInfo.Account != nil {
			payload.Product.Account = newWorker.TaskInfo.Account.Email + ":" + newWorker.TaskInfo.Account.Password
		}

		for _, v := range newWorker.TaskInfo.GiftCardGroup {
			payload.Product.GiftCards += v.CardNumber + ","
		}
	}

	// send request
	err = successHandlerStream.Send(payload)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}
