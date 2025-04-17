package defineder

import (
	"github.com/hibiken/asynq"
	"sicko-aio-2.0-client/models"
)

const (
	SnkrsLaunchEntry        = `snkrs:launch:v1`
	SnkrsLaunchEntryPrepare = `snkrs:launch:prepare`
)

func NewSnkrsLaunchEntryPrepare(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(SnkrsLaunchEntryPrepare, []byte(taskId))
}

func NewSnkrsLaunchEntry(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(SnkrsLaunchEntry, []byte(taskId))
}
