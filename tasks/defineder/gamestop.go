package defineder

import (
	"github.com/hibiken/asynq"
	"sicko-aio-2.0-client/models"
)

const (
	GamestopMonitor  = `gamestop:monitor:information,restock`
	GamestopCheckout = `gamestop:checkout:v1`
)

func NewGamestopMonitorTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(GamestopMonitor, []byte(taskId))
}

func NewGamestopCheckoutTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(GamestopCheckout, []byte(taskId))
}
