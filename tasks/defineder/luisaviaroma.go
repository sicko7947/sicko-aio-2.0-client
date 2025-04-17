package defineder

import (
	"github.com/hibiken/asynq"
	"sicko-aio-2.0-client/models"
)

const (
	LuisaviaromaMonitor         = `luisaviaroma:monitor:information,restock,checklive`
	LuisaviaromaCheckoutPrepare = `Luisaviaroma:checkout,prepare`
	LuisaviaromaCheckout        = `Luisaviaroma:checkout:v1`
	LuisaviaromaSync            = `Luisaviaroma:login:sync`
)

func NewLuisaviaromaMonitorTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(LuisaviaromaMonitor, []byte(taskId))
}

func NewLuisaviaromaCheckoutPrepare(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(LuisaviaromaCheckoutPrepare, []byte(taskId))
}

func NewLuisaviaromaCheckoutTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(LuisaviaromaCheckout, []byte(taskId))
}

func NewLuisaviaromaSyncTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(LuisaviaromaSync, []byte(taskId))
}
