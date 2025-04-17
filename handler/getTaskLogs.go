package handler

import (
	"github.com/gogf/gf/net/ghttp"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/models"
)

func GetTaskLogs(r *ghttp.Request) {
	taskID := models.TaskID(r.Get("taskId").(string))

	if logs := communicator.TaskLogsGMap.Get(taskID); logs != nil {
		r.Response.WriteJsonExit(logs)
	}
}
