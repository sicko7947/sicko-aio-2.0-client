package handler

import (
	"net/http"
	"time"

	"github.com/gogf/gf/net/ghttp"
	"github.com/google/uuid"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/tasks"
)

type cancelTaskPayload struct {
	GroupID models.TaskGroupID `json:"groupId"`
	IDs     []models.TaskID    `json:"ids"`
}

// CancelTasks : Cancel Selected Tasks
func CancelTasks(r *ghttp.Request) {

	// Parsing Payload
	var p *cancelTaskPayload
	if err := r.Parse(&p); err != nil {
		r.Response.WriteStatus(http.StatusBadRequest)
		r.Response.WriteJson(map[string]bool{"success": false})
		return
	}

	for _, taskID := range p.IDs {
		ticker := time.NewTicker(3 * time.Millisecond)
		taskIDUUID, _ := uuid.Parse(string(taskID))
		go func() {
			for i := 0; i < 10; i++ { // retry 10 times in case log lagging
				tasks.CancelTask(taskIDUUID)
				communicator.ModifyTaskStatus(&models.Message{ // send task cancellation status to frontend
					Code:    0,
					GroupID: p.GroupID,
					TaskID:  taskID,
					Status:  "CANCELLED",
					Message: "Stopped",
				})
				<-time.NewTicker(100 * time.Millisecond).C
			}
		}()
		<-ticker.C
	}

	r.Response.WriteJsonExit(map[string]bool{
		"success": true,
	})
}
