package lanecrawford

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/hibiken/asynq"
	"github.com/sicko7947/sickocommon"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/constants"
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/successHandler"
	"sicko-aio-2.0-client/tasks"
	"sicko-aio-2.0-client/utils/notification"
	"sicko-aio-2.0-client/utils/psychoclient"
)

func (p *checkoutPayload) doAddToCart() *models.Error {
	communicator.ModifyTaskStatus(&models.Message{Code: 200, GroupID: p.worker.GroupID, TaskID: p.taskID, Status: "PENDING", Message: "Adding To Cart"})

	endpoint := "https://www.lanecrawford.com.hk/xxxxxxxxxxxxxx"

	form := url.Values{}
	data := map[string]string{}

	for key, value := range data {
		form.Set(key, value)
	}

	reqID, _ := p.session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "POST",
		Headers:  p.headers,
		Payload:  strings.NewReader(form.Encode()),
	})

	res, respBody, err := p.session.RoundTrip(reqID)
	if err != nil {
		return err
	}
	switch res.StatusCode {
	case 200, 201, 202, 302:
		// check if atc was successful
		fmt.Println(string(respBody))
		return nil
	default:
		return &models.Error{Error: errors.New("ERROR_ADDING_TO_CART"), Code: res.StatusCode, Message: "Error Adding To Cart"}
	}
}

func (p *checkoutPayload) doSubmitShippingInfo() *models.Error {
	communicator.ModifyTaskStatus(&models.Message{Code: 200, GroupID: p.worker.GroupID, TaskID: p.taskID, Status: "PENDING", Message: "Submitting Shipping Info"})

	endpoint := "https://www.lanecrawford.com.hk/xxxxxxxxxxxxxx"

	form := url.Values{}
	data := map[string]string{}

	for key, value := range data {
		form.Set(key, value)
	}

	reqID, _ := p.session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "POST",
		Headers:  p.headers,
		Payload:  strings.NewReader(form.Encode()),
	})

	res, _, err := p.session.RoundTrip(reqID)
	if err != nil {
		return err
	}
	switch res.StatusCode {
	case 200, 201, 202, 302, 203, 204:
		fmt.Println(res)

		// check if atc was successful
		return nil
	default:
		return &models.Error{Error: errors.New("ERROR_SUBMITTING_INFO"), Code: res.StatusCode, Message: "Error Submitting Shipping"}
	}
}

func (p *checkoutPayload) doSubmitPaymentInfo() *models.Error {
	communicator.ModifyTaskStatus(&models.Message{Code: 200, GroupID: p.worker.GroupID, TaskID: p.taskID, Status: "PENDING", Message: "Submitting Payment Info (1)"})

	endpoint := "https://www.lanecrawford.com.hk/xxxxxxxxxxxxxx"

	form := url.Values{}
	data := map[string]string{}

	for key, value := range data {
		form.Set(key, value)
	}

	reqID, _ := p.session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "POST",
		Headers:  p.headers,
		Payload:  strings.NewReader(form.Encode()),
	})

	res, _, err := p.session.RoundTrip(reqID)
	if err != nil {
		return err
	}
	switch res.StatusCode {
	case 200, 201, 202, 302, 203, 204:
		fmt.Println(res)

		// check if atc was successful
		return nil
	default:
		return &models.Error{Error: errors.New("ERROR_SUBMITTING_INFO"), Code: res.StatusCode, Message: "Error Submitting Payment Info (1)"}
	}
}

func (p *checkoutPayload) doSubmitPaymentInfo2() *models.Error {
	communicator.ModifyTaskStatus(&models.Message{Code: 200, GroupID: p.worker.GroupID, TaskID: p.taskID, Status: "PENDING", Message: "Submitting Payment Info (2)"})

	endpoint := "https://www.lanecrawford.com.hk/xxxxxxxxxxxxxx"

	form := url.Values{}
	data := map[string]string{}

	for key, value := range data {
		form.Set(key, value)
	}

	reqID, _ := p.session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "POST",
		Headers:  p.headers,
		Payload:  strings.NewReader(form.Encode()),
	})

	res, _, err := p.session.RoundTrip(reqID)
	if err != nil {
		return err
	}
	switch res.StatusCode {
	case 200, 201, 202, 302, 203, 204:
		fmt.Println(res)

		// check if atc was successful
		return nil
	default:
		return &models.Error{Error: errors.New("ERROR_SUBMITTING_INFO_2"), Code: res.StatusCode, Message: "Error Submitting Payment Info (2)"}
	}
}

func (p *checkoutPayload) doCreatingOrder() *models.Error {
	communicator.ModifyTaskStatus(&models.Message{Code: 200, GroupID: p.worker.GroupID, TaskID: p.taskID, Status: "PENDING", Message: "Creating Order"})

	endpoint := "https://www.lanecrawford.com.hk/xxxxxxxxxxxxxx"

	form := url.Values{}
	data := map[string]string{}

	for key, value := range data {
		form.Set(key, value)
	}

	reqID, _ := p.session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "POST",
		Headers:  p.headers,
		Payload:  strings.NewReader(form.Encode()),
	})

	res, _, err := p.session.RoundTrip(reqID)
	if err != nil {
		return err
	}
	switch res.StatusCode {
	case 200, 201, 202, 302, 203, 204:
		fmt.Println(res)

		// check if atc was successful
		return nil
	default:
		return &models.Error{Error: errors.New("ERROR_CREATING_ORDER"), Code: res.StatusCode, Message: "Error Creating Order"}
	}
}

func (p *checkoutPayload) doAuthorizePayment() *models.Error {
	communicator.ModifyTaskStatus(&models.Message{Code: 200, GroupID: p.worker.GroupID, TaskID: p.taskID, Status: "PENDING", Message: "Authorizing Payment"})

	endpoint := "https://www.lanecrawford.com/xxxxxxxxxxxxx"

	reqID, _ := p.session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "GET",
		Headers:  p.headers,
		Payload:  nil,
	})

	res, respBody, err := p.session.RoundTrip(reqID)
	if err != nil {
		return err
	}

	fmt.Println(res)
	fmt.Println(string(respBody))
	switch res.StatusCode {
	case 200, 201, 202, 302:
		fmt.Println(res)
		fmt.Println(string(respBody))
		// check if atc was successful
		return nil
	default:
		return &models.Error{Error: errors.New("ERROR_CREATING_ORDER"), Code: res.StatusCode, Message: "Error Authorizing Payment"}
	}
}

func (p *checkoutPayload) setupTask() *models.Error {

	checkoutProxyGroup := (communicator.Config.Proxies)[p.worker.CheckoutProxyGroupName]
	if session, err := psychoclient.NewSession(&psychoclient.SessionBuilder{
		FollowRedirects: true,
		Proxy:           sickocommon.GetProxy(checkoutProxyGroup).String(),
	}); err != nil {
		return err
	} else {
		p.session = session
	}

	p.headers = map[string]string{
		"accept":           "application/json, text/javascript, */*; q=0.01",
		"x-requested-with": "XMLHttpRequest",
		"sec-ch-ua-mobile": "?0",
		"dnt":              "1",
		"cache-control":    "no-cache",
		"content-type":     "application/x-www-form-urlencoded",
		"user-agent":       gofakeit.RandomString(constants.ChromeUAList),
	}

	return nil
}

// HandleLanecrawfordCheckoutTaskWithOptions : HandleLanecrawfordCheckoutTaskWithOptions
func HandleLanecrawfordCheckoutTaskWithOptions(ctx context.Context, task *asynq.Task) error {
	// Get basic task info
	p := &checkoutPayload{}
	p.taskID = models.TaskID(task.Payload())
	if worker := communicator.TaskWorkerObjectGMap.Get(p.taskID); worker != nil {
		p.worker = worker.(*models.TaskWorker)
		if !p.worker.Mutex.TryLock() {
			return errors.New("too many tasks assigned to one worker")
		}
		defer func() {
			if p.worker.Mutex.IsLocked() {
				p.worker.Mutex.Unlock()
			}
		}()
	} else {
		return errors.New("error starting worker")
	}

	// setup task group settings
	p.taskGroupSetting = communicator.Config.TaskGroups[p.worker.GroupID].TaskGroupSetting

	done := make(chan bool)
	c := make(chan *models.Error, 1)
	defer close(c)
	defer close(done)
	go func() {
		defer recover()

		if err := p.setupTask(); err != nil {
			tasks.SafeSend(c, err)
			return
		}

		if err := p.doGetDynSessConf(); err != nil {
			tasks.SafeSend(c, err)
			return
		}

		if err := p.doAddToCart(); err != nil {
			tasks.SafeSend(c, err)
			return
		}

		if err := p.doSubmitShippingInfo(); err != nil {
			tasks.SafeSend(c, err)
			return
		}

		if err := p.doSubmitPaymentInfo(); err != nil {
			tasks.SafeSend(c, err)
			return
		}

		if err := p.doSubmitPaymentInfo2(); err != nil {
			tasks.SafeSend(c, err)
			return
		}

		if err := p.doCreatingOrder(); err != nil {
			tasks.SafeSend(c, err)
			return
		}

		if err := p.doAuthorizePayment(); err != nil {
			tasks.SafeSend(c, err)
			return
		}

		done <- true
	}()

	select {
	case <-ctx.Done():
		communicator.ModifyTaskStatus(&models.Message{ // send task cancellation status to frontend
			Code:    0,
			GroupID: p.worker.GroupID,
			TaskID:  p.taskID,
			Status:  "CANCELLED",
			Message: "Stopped",
		})
		return ctx.Err() // cancelation signal received, abandon this work.
	case err := <-c:
		fmt.Println(err)
		communicator.ModifyTaskStatus(&models.Message{ // send error message to frontend
			Code:    err.Code,
			GroupID: p.worker.GroupID,
			TaskID:  p.taskID,
			Status:  "ARCHIVED",
			Message: err.Message,
		})
		return err.Error
	case <-done:
		communicator.ModifyTaskStatus(&models.Message{ // send success checkout message to frontend
			Code:    200,
			GroupID: p.worker.GroupID,
			TaskID:  p.taskID,
			Status:  "COMPLETED",
			Message: "Checked Out!",
		})

		// send success checkout message to webhooks
		go successHandler.HandlerSuccess(p.taskGroupSetting, p.worker)
		go notification.PushNotification(p.taskGroupSetting, p.worker)
		return nil
	}
}
