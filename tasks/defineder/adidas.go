package defineder

import (
	"github.com/hibiken/asynq"
	"sicko-aio-2.0-client/models"
)

const (
	AdidasTask = `adidas:queue,checkout,restock`
)

func NewAdidasTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(AdidasTask, []byte(taskId))
}
