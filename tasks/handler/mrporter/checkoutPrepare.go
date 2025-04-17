package mrporter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/hibiken/asynq"
	"github.com/huandu/go-clone"
	"github.com/sicko7947/sickocommon"
	http "github.com/zMrKrabz/fhttp"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/constants"
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/tasks"
	"sicko-aio-2.0-client/utils"
	"sicko-aio-2.0-client/utils/psychoclient"
)

func (p *checkoutPayload) setupCheckoutPrepare() *models.Error {
	communicator.ModifyTaskStatus(&models.Message{Code: 200, GroupID: p.worker.GroupID, TaskID: p.taskID, Status: "PENDING", Message: "Preparing Checkout"})

	// setup task group settings
	p.taskGroupSetting = communicator.Config.TaskGroups[p.worker.GroupID].TaskGroupSetting

	// setup profiles & accounts & giftcards
	profileGroup := (communicator.Config.Profiles)[p.worker.ProfileGroupName]
	accountGroup := (communicator.Config.Accounts)[p.worker.AccountGroupName]

	// setup proxy groups
	checkoutProxyGroup := (communicator.Config.Proxies)[p.worker.CheckoutProxyGroupName]
	proxy := sickocommon.GetProxy(checkoutProxyGroup).String()
	if sesh, err := psychoclient.NewSession(&psychoclient.SessionBuilder{
		Proxy: proxy,
	}); err != nil {
		return err
	} else {
		p.session = sesh
	}

	// setup taskInfo
	if p.worker.TaskInfo == nil {
		p.worker.TaskInfo = new(models.WorkerTaskInfo)
	}

	// setup profile
	for _, v := range profileGroup {
		p.worker.TaskInfo.Profile = utils.JigProfile(clone.Clone(v).(*models.Profile))
		break
	}

	// setup request headers
	p.headers = map[string]string{
		"accept":              "*/*",
		"accept-language":     "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7,es;q=0.6",
		"application-name":    "checkout",
		"application-version": "7.283.0",
		"content-type":        "application/json",
		"label":               "setPaymentInstruction",
		"sec-ch-ua-mobile":    "?0",
		"sec-fetch-dest":      "empty",
		"sec-fetch-mode":      "cors",
		"sec-fetch-site":      "same-origin",
		"user-agent":          gofakeit.RandomString(constants.ChromeUAList),
		"x-ibm-client-id":     "0b1e2c22-581d-435b-9cde-70bc52cba701",
		"cache-control":       "no-cache",
	}

	// setup account
	if len(accountGroup) > 0 {
		index := rand.Intn(len(accountGroup))
		account := accountGroup[index]

		switch { // if account never synced befor
		case account.Status == "", account.AccessToken == "":
			account, err := Sync(account, proxy)
			if err != nil {
				return err
			}
			accountGroup[index] = account
			p.worker.TaskInfo.Account = account
		default:
			p.worker.TaskInfo.Account = account
		}
		p.headers["x-ubertoken"] = p.worker.TaskInfo.Account.AccessToken
		p.session.SetCookies(map[string]*http.Cookie{
			"ubertoken": {
				Name:  "ubertoken",
				Value: p.worker.TaskInfo.Account.AccessToken,
			},
		})
	}

	switch {
	case p.worker.TaskInfo.Profile == nil:
		return &models.Error{Error: errors.New("ERROR_ASSIGN_PROFILE"), Code: 500, Message: "Error Assign Profile"}
	case p.worker.TaskInfo.Account == nil:
		return &models.Error{Error: errors.New("ERROR_ASSIGN_ACCOUNT"), Code: 500, Message: "Error Assign Account"}
	}

	return nil
}

func (p *checkoutPayload) doCheckCart() *models.Error {
	communicator.ModifyTaskStatus(&models.Message{Code: 200, GroupID: p.worker.GroupID, TaskID: p.taskID, Status: "PENDING", Message: "Checking Cart"})

	endpoint := fmt.Sprintf("https://xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", p.taskGroupSetting.Locale)

	data, _ := json.Marshal(make(map[string]string))

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
	case 200, 201, 202:

		return nil
	default:
		return &models.Error{Error: errors.New("ERROR_CHECKING_CART"), Code: res.StatusCode, Message: "Error Checking Cart"}
	}
}

func (p *checkoutPayload) doClearCart(data []byte) *models.Error {
	communicator.ModifyTaskStatus(&models.Message{Code: 200, GroupID: p.worker.GroupID, TaskID: p.taskID, Status: "PENDING", Message: "Clearing Cart"})

	endpoint := fmt.Sprintf("https://xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", p.taskGroupSetting.Locale)

	reqID, _ := p.session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "DELETE",
		Headers:  p.headers,
		Payload:  bytes.NewBuffer(data),
	})

	res, _, err := p.session.Do(reqID)
	if err != nil {
		return err
	}
	switch res.StatusCode {
	case 200, 201, 202:
		return nil
	default:
		return &models.Error{Error: errors.New("ERROR_CLEARING_CART"), Code: res.StatusCode, Message: "Error Clearing Cart"}
	}
}

func HandleMrPorterCheckoutPrepareTaskWithOptions(ctx context.Context, task *asynq.Task) error {

	// Get basic task info
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
	errCh := make(chan *models.Error, 1)
	defer close(errCh)
	defer close(done)
	go func() {
		var err *models.Error
		defer func() {
			recover()
		}()

		err = p.setupCheckoutPrepare()
		if tasks.SafeSend(errCh, err) {
			return
		}

		err = p.doCheckCart()
		if tasks.SafeSend(errCh, err) {
			return
		}

		done <- true
	}()

	for {
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
		case err := <-errCh:
			if err == nil {
				continue
			}
			communicator.ModifyTaskStatus(&models.Message{
				Code:    err.Code,
				GroupID: p.worker.GroupID,
				TaskID:  p.taskID,
				Status:  "ARCHIVED",
				Message: err.Message,
			})
			return err.Error
		case <-done:
			communicator.TaskWorkerObjectGMap.Set(p.taskID, p.worker) // save the taskGroup
			communicator.ModifyTaskStatus(&models.Message{
				GroupID: p.worker.GroupID,
				TaskID:  p.taskID,
				Status:  "ARCHIVED",
				Message: "Waiting for Execution",
			})
			return nil
		}
	}
}
