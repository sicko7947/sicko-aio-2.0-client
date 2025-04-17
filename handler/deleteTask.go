package handler

import (
	"net/http"

	"github.com/gogf/gf/net/ghttp"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/models"
)

type deleteTaskPayload struct {
	GroupID    models.TaskGroupID `json:"groupId"`
	ScraperIds []models.TaskID    `json:"scraperIds"`
	WorkerIds  []models.TaskID    `json:"workerIds"`
}

// CancelTasks : Cancel Selected Tasks
func DeleteTasks(r *ghttp.Request) {

	// Parsing Payload
	var p *deleteTaskPayload
	if err := r.Parse(&p); err != nil {
		r.Response.WriteStatus(http.StatusBadRequest)
		r.Response.WriteJsonExit(map[string]bool{"success": false})
		return
	}

	group := communicator.Config.TaskGroups[p.GroupID]

	for _, taskID := range p.ScraperIds {
		communicator.TaskScraperObjectGMap.Remove(taskID)
		communicator.TaskMessageGMap.Remove(taskID)
		communicator.TaskLogsGMap.Remove(taskID)

		delete(group.TaskScrapers, taskID)
	}

	for _, taskID := range p.WorkerIds {
		communicator.TaskWorkerObjectGMap.Remove(taskID)
		communicator.TaskMessageGMap.Remove(taskID)
		communicator.TaskLogsGMap.Remove(taskID)

		delete(group.TaskWorkers, taskID)
	}

	r.Response.WriteJsonExit(map[string]bool{"success": true})
}
