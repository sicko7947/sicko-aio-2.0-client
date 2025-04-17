package macys

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	http "github.com/zMrKrabz/fhttp"

	"github.com/hibiken/asynq"
	"github.com/sicko7947/sickocommon"
	"github.com/tidwall/gjson"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/tasks"
	"sicko-aio-2.0-client/utils/psychoclient"
	"sicko-aio-2.0-client/utils/redis"
)

func (p *checkoutPayload) doAddToCart() *models.Error {
	endpoint := "https://www.macys.com/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

	data, _ := json.Marshal(map[string]interface{}{})

	reqID, _ := p.session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "POST",
		Headers:  p.headers,
		Payload:  bytes.NewBuffer(data),
	})

	res, _, err := p.session.Do(reqID)
	if err != nil {
		return err
	}
	switch res.StatusCode {
	case 200:

		fallthrough
	default:
		return &models.Error{Error: errors.New("ERROR_ADDING_TO_CART"), Code: res.StatusCode, Message: "Error Adding to Cart"}
	}
}

func (p checkoutPayload) doContinueCheckout() *models.Error {
	endpoint := fmt.Sprintf("https://www.macys.com/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", p.checkoutID)

	data, _ := json.Marshal(map[string]string{})

	reqID, _ := p.session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "POST",
		Headers:  p.headers,
		Payload:  bytes.NewBuffer(data),
	})

	res, respBody, err := p.session.Do(reqID)
	if err != nil {
		return err
	}
	result := gjson.ParseBytes(respBody)
	switch res.StatusCode {
	case 200:
		if result.Get("bag").Exists() {
			return nil
		}
		fallthrough
	default:
		return &models.Error{Error: errors.New("ERROR_START_CHECKOUT"), Code: 509, Message: "Error Start Checkout Session"}
	}
}

func (p checkoutPayload) doSubmitPayment() *models.Error {
	endpoint := "https://www.macys.com/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

	reqID, _ := p.session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "GET",
		Headers:  p.headers,
		Payload:  nil,
	})

	res, respBody, err := p.session.Do(reqID)
	if err != nil {
		return err
	}
	result := gjson.Parse(string(respBody))
	switch res.StatusCode {
	case 200:
		if result.Get("order").Exists() {
			return nil
		}
		fallthrough
	default:
		return &models.Error{Error: errors.New("ERROR_GETTING_ORDER_DETAILS"), Code: res.StatusCode, Message: "Error Getting Order Details"}
	}
}

func (p checkoutPayload) doSubmitOrder() *models.Error {
	endpoint := "https://www.macys.com/xxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

	data, _ := json.Marshal(map[string]string{})

	reqID, _ := p.session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "POST",
		Headers:  p.headers,
		Payload:  bytes.NewBuffer(data),
	})

	res, respBody, err := p.session.Do(reqID)
	if err != nil {
		return err
	}

	result := gjson.ParseBytes(respBody)
	fmt.Println(result.String())
	switch res.StatusCode {
	case 200:
		errorObj := result.Get("responsiveConfirm.error.globalErrors.0.code")
		if !errorObj.Exists() {
			return nil
		}
		return &models.Error{Error: errors.New("ERROR_SUBMITTING_ORDER"), Code: 509, Message: errorObj.String()}
	default:
		return &models.Error{Error: errors.New("ERROR_SUBMITTING_ORDER"), Code: 500, Message: "Error Submitting Order"}
	}
}

func (p *checkoutPayload) setupTask() *models.Error {

	// setup task group settings
	p.taskGroupSetting = communicator.Config.TaskGroups[p.worker.GroupID].TaskGroupSetting

	// override quantity limit with max cart setting
	if p.worker.MaxCart && p.worker.Product.QuantityLimit > 0 {
		p.worker.Quantity = p.worker.Product.QuantityLimit
	}

	checkoutProxyGroup := (communicator.Config.Proxies)[p.worker.CheckoutProxyGroupName]

	if session, err := psychoclient.NewSession(&psychoclient.SessionBuilder{
		Proxy: sickocommon.GetProxy(checkoutProxyGroup).String(),
	}); err != nil {
		return err
	} else {
		p.session = session
	}

	// setup session cookie
	useragent, cookieMap, err := redis.GetCookie2FromRedis("akamai.macys.com")
	if err != nil {
		return err
	}
	p.session.SetCookies(cookieMap)
	p.session.SetCookies(map[string]*http.Cookie{
		"currency": {
			Name:  "currency",
			Value: "USD",
		},
		"shippingCountry": {
			Name:  "shippingCountry",
			Value: "US",
		},
	})
	p.headers = map[string]string{
		"accept":             "application/json, text/javascript, */*; q=0.01",
		"accept-language":    "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7,es;q=0.6",
		"content-type":       "application/json",
		"dnt":                "1",
		"origin":             "https://www.macys.com",
		"user-agent":         useragent,
		"x-macys-request-id": sickocommon.NikeUUID(),
		"x-requested-with":   "XMLHttpRequest",
		"cache-control":      "no-cache",
	}

	return nil
}

// HandleMacyCheckoutTaskWithOptions : HandleMacyCheckoutTaskWithOptions
func HandleMacyCheckoutTaskWithOptions(ctx context.Context, task *asynq.Task) error {

	// set base payload
	p := &checkoutPayload{
		taskID: models.TaskID(task.Payload()),
	}

	worker := communicator.TaskWorkerObjectGMap.Get(p.taskID)
	if worker == nil {
		return errors.New("error starting worker")
	}

	p.worker = worker.(*models.TaskWorker)
	if !p.worker.Mutex.TryLock() {
		return errors.New("too many tasks assigned to one worker")
	}

	defer func() {
		if p.worker.Mutex.IsLocked() {
			p.worker.Mutex.Unlock()
		}
	}()

	done := make(chan bool)
	c := make(chan *models.Error, 1)
	defer close(c)
	defer close(done)
	go func() {
		if err := p.setupTask(); err != nil {
			tasks.SafeSend(c, err)
			return
		}

		if err := p.doAddToCart(); err != nil {
			tasks.SafeSend(c, err)
			return
		}

		if err := p.doSubmitPayment(); err != nil {
			tasks.SafeSend(c, err)
			return
		}

		if err := p.doSubmitOrder(); err != nil {
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
		communicator.ModifyTaskStatus(&models.Message{ // send error message to frontend
			Code:    err.Code,
			GroupID: p.worker.GroupID,
			TaskID:  p.taskID,
			Status:  "ARCHIVED",
			Message: err.Message,
		})
		return err.Error
	case <-done:
		status := "COMPLETED"
		if p.worker.Restart {
			status = "ARCHIVED"
		}
		communicator.ModifyTaskStatus(&models.Message{ // send success checkout message to frontend
			GroupID: p.worker.GroupID,
			TaskID:  p.taskID,
			Code:    200,
			Status:  status,
			Message: "Checked Out!",
		})

		// send success checkout message to webhooks
		// webhook.SendDiscordLegacyWebhook(p.taskGroupSetting, p.worker)
		// webhook.SendSlackLegacyWebhook(p.taskGroupSetting, p.worker)
		return nil
	}
}
