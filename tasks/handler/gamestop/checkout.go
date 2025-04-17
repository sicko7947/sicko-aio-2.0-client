package gamestop

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/sicko7947/sickocommon"
	"github.com/tidwall/gjson"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/tasks"
	"sicko-aio-2.0-client/utils/psychoclient"
)

func (p checkoutPayload) doAddToCart() *models.Error {
	endpoint := "https://www.gamestop.com/xxxxxx"

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

	res, respBody, err := p.session.Do(reqID)

	if err != nil {
		return err
	}

	result := gjson.Parse(string(respBody))
	switch res.StatusCode {
	case 200, 201, 202:
		// check if atc was successful
		if errors := result.Get("error").Bool(); !errors {
			return nil
		}
		fallthrough
	default:
		return &models.Error{Error: errors.New("ERROR_ADDING_TO_CART"), Code: res.StatusCode, Message: "Error Adding To Cart"}
	}
}

func (p *checkoutPayload) doSubmitShipping() *models.Error {
	endpoint := "https://www.gamestop.com/xxx"

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

	res, _, err := p.session.Do(reqID)
	if err != nil {
		return err
	}
	switch res.StatusCode {
	case 200, 201, 202:
		return nil
	default:
		return &models.Error{Error: errors.New("ERROR_SUBMITTING_SHIPPING_METHOD"), Code: 509, Message: "Error Submitting Shipping Method"}
	}
}

func (p *checkoutPayload) doSubmitPayment() *models.Error {
	endpoint := "https://www.gamestop.com/xxxx"

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

	res, respBody, err := p.session.Do(reqID)
	if err != nil {
		return err
	}
	result := gjson.Parse(string(respBody))
	switch res.StatusCode {
	case 200, 201, 202:
		// check if atc was successful
		if errors := result.Get("error").Bool(); !errors {
			return nil
		}
		fallthrough
	default:
		return &models.Error{Error: errors.New("ERROR_SUBMITTING_PAYMENT"), Code: 509, Message: "Payment Invalid"}
	}
}

func (p *checkoutPayload) setupTask() {
	profileGroup := (communicator.Config.Profiles)[p.worker.ProfileGroupName]

	checkoutProxyGroup := (communicator.Config.Proxies)[p.worker.CheckoutProxyGroupName]
	p.session, _ = psychoclient.NewSession(&psychoclient.SessionBuilder{
		Proxy: sickocommon.GetProxy(checkoutProxyGroup).String(),
	})

	p.worker.TaskInfo.CreditCardToken = sickocommon.NikeUUID()

	p.headers = map[string]string{
		"accept":           "*/*",
		"accept-language":  "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7,es;q=0.6",
		"content-type":     "application/x-www-form-urlencoded",
		"dnt":              "1",
		"origin":           "https://www.gamestop.com",
		"referer":          "https://www.gamestop.com",
		"sec-ch-ua-mobile": "?0",
		"sec-fetch-dest":   "empty",
		"sec-fetch-mode":   "cors",
		"sec-fetch-site":   "same-origin",
		"user-agent":       "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36",
		"x-requested-with": "XMLHttpRequest",
		"cache-control":    "no-cache",
	}

	var profile *models.Profile
	for _, v := range profileGroup {
		profile = v
		break
	}
	switch {
	case profile != nil:
		p.worker.TaskInfo.Profile = profile
	}
}

// HandleGamestopCheckoutTaskWithOptions : HandleGamestopCheckoutTaskWithOptions
func HandleGamestopCheckoutTaskWithOptions(ctx context.Context, task *asynq.Task) error {
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

	p.setupTask() // further setup task info

	done := make(chan bool)
	c := make(chan *models.Error, 1)
	defer close(c)
	defer close(done)
	go func() {

		if err := p.doAddToCart(); err != nil {
			tasks.SafeSend(c, err)
			return
		}

		if err := p.doSubmitShipping(); err != nil {
			tasks.SafeSend(c, err)
			return
		}

		if err := p.doSubmitPayment(); err != nil {
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
			Status:  "PENDING",
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
