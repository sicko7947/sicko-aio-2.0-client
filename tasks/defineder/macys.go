package defineder

import (
	"github.com/hibiken/asynq"
	"sicko-aio-2.0-client/models"
)

const (
	MacysMonitor         = `maycs:monitor:information,restock`
	MacysCheckoutPrepare = `maycs:checkout:prepare`
	MacysCheckout        = `maycs:checkout`
)

func NewMacysMonitorTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(MacysMonitor, []byte(taskId))
}

func NewMacysCheckoutPrepareTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(MacysCheckoutPrepare, []byte(taskId))
}

func NewMacysCheckoutTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(MacysCheckout, []byte(taskId))
}
