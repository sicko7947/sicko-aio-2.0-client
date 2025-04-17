package macys

import (
	"context"
	"errors"
	"math/rand"
	"time"

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

func (p *checkoutPayload) setupCheckoutPrepare() *models.Error {

	// setup profiles & accounts & giftcards
	profileGroup := (communicator.Config.Profiles)[p.worker.ProfileGroupName]
	accountGroup := (communicator.Config.Accounts)[p.worker.AccountGroupName]

	// setup proxy groups
	checkoutProxyGroup := (communicator.Config.Proxies)[p.worker.CheckoutProxyGroupName]
	proxy := sickocommon.GetProxy(checkoutProxyGroup).String()

	// setup request sessions
	if sesh, err := psychoclient.NewSession(&psychoclient.SessionBuilder{
		Proxy: proxy,
	}); err != nil {
		return err
	} else {
		p.session = sesh
	}

	// setup taskInfos
	if p.worker.TaskInfo == nil {
		p.worker.TaskInfo = new(models.WorkerTaskInfo)
	}

	p.headers = map[string]string{
		"accept":             "application/json, text/javascript, */*; q=0.01",
		"accept-language":    "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7,es;q=0.6",
		"content-type":       "application/json",
		"dnt":                "1",
		"origin":             "https://www.macys.com",
		"user-agent":         gofakeit.RandomString(constants.ChromeUAList),
		"x-macys-request-id": sickocommon.NikeUUID(),
		"x-requested-with":   "XMLHttpRequest",
		"cache-control":      "no-cache",
	}

	for _, v := range profileGroup {
		p.worker.TaskInfo.Profile = utils.JigProfile(clone.Clone(v).(*models.Profile))
		p.worker.TaskInfo.Email = p.worker.TaskInfo.Profile.Email
		break
	}

	if len(accountGroup) > 0 {
		index := rand.Intn(len(accountGroup))
		account := accountGroup[index]

		lastSyncTime, errParseLastSyncTime := time.Parse("2006-01-02T15:04:05.000Z", account.LastSyncTime)

		switch { // if account never synced befor
		case account.Status == "", errParseLastSyncTime != nil:
			account, err := Sync(account, proxy)
			if err != nil {
				return err
			}
			accountGroup[index] = account
			p.worker.TaskInfo.Account = account
		case time.Now().After(lastSyncTime.Add(time.Minute * 178)): // access token expired. require refresh
			account, err := Sync(account, proxy)
			if err != nil {
				return err
			}
			accountGroup[index] = account
			p.worker.TaskInfo.Account = account
		default: // access token is valid
			p.worker.TaskInfo.Account = account
		}
	}

	switch {
	case p.worker.TaskInfo.Profile == nil:
		return &models.Error{Error: errors.New("ERROR_ASSIGN_TASKINFO"), Code: 500, Message: "Error Assign Profile"}
	case len(accountGroup) > 0 && p.worker.TaskInfo.Account == nil:
		return &models.Error{Error: errors.New("ERROR_ASSIGN_TASKINFO"), Code: 500, Message: "Error Assign Account"}
	}

	return nil
}

// HandleMacyCheckoutPrepareTaskWithOptions : HandleMacyCheckoutPrepareTaskWithOptions
func HandleMacyCheckoutPrepareTaskWithOptions(ctx context.Context, task *asynq.Task) error {

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

	communicator.ModifyTaskStatus(&models.Message{GroupID: p.worker.GroupID, TaskID: p.taskID, Message: "Preparing for Checkout", Status: "PENDING"})

	done := make(chan bool)
	errCh := make(chan *models.Error, 1)
	defer close(errCh)
	defer close(done)
	go func() {
		// var err *models.Error
		defer func() {
			recover()
		}()

		if err := p.setupCheckoutPrepare(); err != nil {
			tasks.SafeSend(errCh, err)
			return
		}

		done <- true
	}()

	for {
		select {
		case <-ctx.Done():
			communicator.ModifyTaskStatus(&models.Message{ // send task cancellation status to frontend
				Code:    400,
				GroupID: p.worker.GroupID,
				TaskID:  p.taskID,
				Status:  "CANCELLED",
				Message: "Stopped",
			})
			return nil // cancelation signal received, abandon this work.
		case err := <-errCh:
			if err == nil {
				continue
			}
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
}
