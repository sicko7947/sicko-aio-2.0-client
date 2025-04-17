package defineder

import (
	"github.com/hibiken/asynq"
	"sicko-aio-2.0-client/models"
)

const (
	LouisVuittonMonitor         = `louisvuitton:monitor:information,restock,checklive`
	LouisVuittonCheckoutPrepare = `louisvuitton:checkout,prepare`
	LouisVuittonCheckout        = `louisvuitton:checkout:v1`
	LouisVuittonSync            = `louisvuitton:login:sync`
)

func NewLouisVuittonMonitorTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(LouisVuittonMonitor, []byte(taskId))
}

func NewLouisVuittonCheckoutPrepare(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(LouisVuittonCheckoutPrepare, []byte(taskId))
}

func NewLouisVuittonCheckoutTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(LouisVuittonCheckout, []byte(taskId))
}

func NewLouisVuittonSyncTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(LouisVuittonSync, []byte(taskId))
}
