package louisvuitton

import (
	"context"
	"errors"
	"math/rand"

	"github.com/hibiken/asynq"
	"github.com/huandu/go-clone"
	"github.com/sicko7947/sickocommon"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/utils"
)

func (p *checkoutPayload) setupCheckoutPrepare() *models.Error {
	communicator.ModifyTaskStatus(&models.Message{Code: 200, GroupID: p.worker.GroupID, TaskID: p.taskID, Status: "PENDING", Message: "Preparing Checkout"})

	// setup profiles & accounts & giftcards
	profileGroup := (communicator.Config.Profiles)[p.worker.ProfileGroupName]
	accountGroup := (communicator.Config.Accounts)[p.worker.AccountGroupName]

	// setup proxy groups
	checkoutProxyGroup := (communicator.Config.Proxies)[p.worker.CheckoutProxyGroupName]
	proxy := sickocommon.GetProxy(checkoutProxyGroup).String()

	// setup taskInfo
	if p.worker.TaskInfo == nil {
		p.worker.TaskInfo = new(models.WorkerTaskInfo)
	}

	// setup profile
	for _, v := range profileGroup {
		p.worker.TaskInfo.Profile = utils.JigProfile(clone.Clone(v).(*models.Profile))
		break
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
	}

	switch {
	case p.worker.TaskInfo.Profile == nil:
		return &models.Error{Error: errors.New("ERROR_ASSIGN_PROFILE"), Code: 500, Message: "Error Assign Profile"}
	case p.worker.TaskInfo.Account == nil:
		return &models.Error{Error: errors.New("ERROR_ASSIGN_ACCOUNT"), Code: 500, Message: "Error Assign Account"}
	}

	return nil
}

func HandleLouisVuittonCheckoutPrepareTaskWithOptions(ctx context.Context, task *asynq.Task) error {
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

	if err := p.setupCheckoutPrepare(); err != nil {
		communicator.ModifyTaskStatus(&models.Message{
			Code:    err.Code,
			GroupID: p.worker.GroupID,
			TaskID:  p.taskID,
			Status:  "ARCHIVED",
			Message: err.Message,
		})
		return err.Error
	}

	communicator.ModifyTaskStatus(&models.Message{
		GroupID: p.worker.GroupID,
		TaskID:  p.taskID,
		Status:  "ARCHIVED",
		Message: "Waiting for Execution",
	})
	return nil
}
