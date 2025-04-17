package defineder

import (
	"github.com/hibiken/asynq"
	"sicko-aio-2.0-client/models"
)

const (
	FinishlineCreateAccountTask = `finishline:create,account`
	FinishlineEnterRaffleTask   = `finishline:login,raffle`
	FinishlineRevealEntryTask   = `finishline:reveal,raffle`
)

func NewFinishlineCreateAccountTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(FinishlineCreateAccountTask, []byte(taskId))
}

func NewFinishlineEnterRaffleTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(FinishlineEnterRaffleTask, []byte(taskId))
}

func NewFinishlineRevealEntryTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(FinishlineRevealEntryTask, []byte(taskId))
}
