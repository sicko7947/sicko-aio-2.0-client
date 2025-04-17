package handler

import (
	"net/http"

	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gmutex"
	"github.com/google/uuid"
	"github.com/huandu/go-clone"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/tasks"
)

type modifyTaskPayload struct {
	Scrapers []*models.TaskScraper `json:"scrapers,omitempty"`
	Workers  []*models.TaskWorker  `json:"workers,omitempty"`
}

func ModifyTaskGroupScraper(r *ghttp.Request) {

	var p *modifyTaskPayload
	if err := r.Parse(&p); err != nil {
		r.Response.WriteStatus(http.StatusBadRequest)
		r.Response.WriteJsonExit(map[string]bool{"success": false})
		return
	}

	t := clone.Clone(p).(*modifyTaskPayload) // deep copy

	for _, v := range t.Scrapers {
		v.Mutex = gmutex.New()
		taskIDUUID, _ := uuid.Parse(string(v.TaskID))
		tasks.CancelTask(taskIDUUID)

		if communicator.Config.TaskGroups[v.GroupID].TaskScrapers == nil {
			communicator.Config.TaskGroups[v.GroupID].TaskScrapers = make(map[models.TaskID]*models.TaskScraper)
		}

		v2 := clone.Clone(v).(*models.TaskScraper)
		communicator.Config.TaskGroups[v.GroupID].TaskScrapers[v.TaskID] = v
		communicator.TaskScraperObjectGMap.Set(v.TaskID, v2)
	}

	r.Response.WriteJsonExit(map[string]bool{"success": true})
}

func ModifyTaskGroupWorker(r *ghttp.Request) {

	var p *modifyTaskPayload
	if err := r.Parse(&p); err != nil {
		r.Response.WriteStatus(http.StatusBadRequest)
		r.Response.WriteJsonExit(map[string]bool{"success": false})
		return
	}

	t := clone.Clone(p).(*modifyTaskPayload) // deep copy

	for _, v := range t.Workers {
		v.Mutex = gmutex.New()
		taskIDUUID, _ := uuid.Parse(string(v.TaskID))
		tasks.CancelTask(taskIDUUID)

		if communicator.Config.TaskGroups[v.GroupID].TaskWorkers == nil {
			communicator.Config.TaskGroups[v.GroupID].TaskWorkers = make(map[models.TaskID]*models.TaskWorker)
		}

		v2 := clone.Clone(v).(*models.TaskWorker)

		communicator.Config.TaskGroups[v.GroupID].TaskWorkers[v.TaskID] = v
		communicator.TaskWorkerObjectGMap.Set(v.TaskID, v2)
	}

	r.Response.WriteJsonExit(map[string]bool{"success": true})
}
