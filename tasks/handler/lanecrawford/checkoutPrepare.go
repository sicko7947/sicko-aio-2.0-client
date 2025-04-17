package lanecrawford

import (
	"bytes"
	"context"
	"errors"
	"math/rand"

	"github.com/PuerkitoBio/goquery"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/hibiken/asynq"
	"github.com/huandu/go-clone"
	"github.com/sicko7947/sickocommon"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/constants"
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/tasks"
	"sicko-aio-2.0-client/utils"
	"sicko-aio-2.0-client/utils/psychoclient"
)

func (p *checkoutPayload) doGetDynSessConf() *models.Error {
	communicator.ModifyTaskStatus(&models.Message{Code: 200, GroupID: p.worker.GroupID, TaskID: p.taskID, Status: "PENDING", Message: "Generating Session"})

	endpoint := "https://www.lanecrawford.com.hk/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

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
	switch res.StatusCode {
	case 200, 201, 202:
		// Load the HTML document
		document, err := goquery.NewDocumentFromReader(bytes.NewReader(respBody))

		if err != nil {
			return &models.Error{Error: nil, Code: 500, Message: "ERROR_PARSING_RESPONSE_BODY"}
		}

		if _dynSessConf, exist := document.Find(`form[id=globalAddItemForm]>div input[name=_dynSessConf]`).Attr("value"); exist {
			if len(_dynSessConf) > 0 {
				p.worker.TaskInfo.Cookies = _dynSessConf
				return nil
			}
		}
		fallthrough
	default:
		return &models.Error{Error: errors.New("ERROR_GETTING_SESSION"), Code: res.StatusCode, Message: "Error Getting Session"}
	}
}

func (p *checkoutPayload) setup() *models.Error {

	// setup profiles & accounts & giftcards
	profileGroup := (communicator.Config.Profiles)[p.worker.ProfileGroupName]
	accountGroup := (communicator.Config.Accounts)[p.worker.AccountGroupName]

	// setup proxy groups
	checkoutProxyGroup := (communicator.Config.Proxies)[p.worker.CheckoutProxyGroupName]

	// setup session
	if session, err := psychoclient.NewSession(&psychoclient.SessionBuilder{
		Proxy: sickocommon.GetProxy(checkoutProxyGroup).String(),
	}); err != nil {
		return err
	} else {
		p.session = session
	}

	p.headers = map[string]string{
		"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"accept-language":           "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7,es;q=0.6",
		"cache-control":             "no-cache",
		"dnt":                       "1",
		"host":                      "www.lanecrawford.com.hk",
		"sec-ch-ua-mobile":          "?0",
		"sec-fetch-dest":            "document",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-site":            "same-origin",
		"sec-fetch-user":            "?1",
		"upgrade-insecure-requests": "1",
		"user-agent":                gofakeit.RandomString(constants.ChromeUAList),
	}

	// setup taskInfos
	if p.worker.TaskInfo == nil {
		p.worker.TaskInfo = new(models.WorkerTaskInfo)
	}

	for _, v := range profileGroup {
		p.worker.TaskInfo.Profile = utils.JigProfile(clone.Clone(v).(*models.Profile))
		break
	}

	if len(accountGroup) > 0 {
		index := rand.Intn(len(accountGroup))
		account := accountGroup[index]
		if account != nil {
			p.worker.TaskInfo.Account = account
			p.worker.TaskInfo.Email = account.Email
		}
	}

	switch {
	case p.worker.TaskInfo.Profile == nil:
		return &models.Error{Error: errors.New("ERROR_ASSIGN_PROFILE"), Code: 500, Message: "Error Assign Profile"}
	case len(accountGroup) > 0 && p.worker.TaskInfo.Account == nil:
		return &models.Error{Error: errors.New("ERROR_ASSIGN_ACCOUNT"), Code: 500, Message: "Error Assign Account"}
	}

	return nil
}

// HandleLanecrawfordCheckoutPrepareTaskWithOptions : HandleLanecrawfordCheckoutPrepareTaskWithOptions
func HandleLanecrawfordCheckoutPrepareTaskWithOptions(ctx context.Context, task *asynq.Task) error {
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
	errCh := make(chan *models.Error, 1)
	defer close(errCh)
	defer close(done)
	go func() {
		defer recover()

		// further setup task info
		if err := p.setup(); err != nil {
			tasks.SafeSend(errCh, err)
			return
		}

		if err := p.doGetDynSessConf(); err != nil {
			tasks.SafeSend(errCh, err)
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
		return nil // cancelation signal received, abandon this work.
	case err := <-errCh:
		communicator.ModifyTaskStatus(&models.Message{ // send error message to frontend
			Code:    err.Code,
			GroupID: p.worker.GroupID,
			TaskID:  p.taskID,
			Status:  "PENDING",
			Message: err.Message,
		})
		return err.Error
	case <-done:
		communicator.TaskWorkerObjectGMap.Set(p.taskID, p.worker) // save the taskGroup
		communicator.ModifyTaskStatus(&models.Message{            // send success checkout message to frontend
			Code:    200,
			GroupID: p.worker.GroupID,
			TaskID:  p.taskID,
			Status:  "ARCHIVED",
			Message: "Waiting for Execution",
		})
		return nil
	}
}
