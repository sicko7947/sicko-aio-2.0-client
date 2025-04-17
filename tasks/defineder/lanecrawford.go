package defineder

import (
	"github.com/hibiken/asynq"
	"sicko-aio-2.0-client/models"
)

const (
	LanecrawfordMonitor         = `lanecrawford:monitor:information,restock`
	LanecrawfordCheckout        = `lanecrawford:checkout:v1`
	LanecrawfordCheckoutPrepare = `lanecrawford:checkout:prepare`
)

func NewLanecrawfordMonitorTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(LanecrawfordMonitor, []byte(taskId))
}

func NewLanecrawfordCheckoutPrepareTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(LanecrawfordCheckoutPrepare, []byte(taskId))
}

func NewLanecrawfordCheckoutTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(LanecrawfordCheckout, []byte(taskId))
}
