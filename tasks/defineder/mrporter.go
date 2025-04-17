package defineder

import (
	"github.com/hibiken/asynq"
	"sicko-aio-2.0-client/models"
)

const (
	MrPorterMonitor         = `mrporter:monitor:information,restock,checklive`
	MrPorterCheckoutPrepare = `mrporter:checkout,prepare`
	MrPorterCheckout        = `mrporter:checkout:v1`
	MrPorterSync            = `mrporter:login:sync`
)

func NewMrPorterMonitorTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(MrPorterMonitor, []byte(taskId))
}

func NewMrPorterCheckoutPrepare(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(MrPorterCheckoutPrepare, []byte(taskId))
}

func NewMrPorterCheckoutTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(MrPorterCheckout, []byte(taskId))
}

func NewMrPorterSyncTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(MrPorterSync, []byte(taskId))
}
