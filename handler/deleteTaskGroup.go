package handler

import (
	"github.com/gogf/gf/net/ghttp"
	"github.com/google/uuid"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/tasks"
)

func DeleteTaskGroup(r *ghttp.Request) {
	groupID := models.TaskGroupID(r.Get("groupId").(string))

	done := make(chan bool, 2)
	go func(done chan bool) { // iterate through task scraper object gmap
		communicator.TaskScraperObjectGMap.Iterator(func(k, v interface{}) bool {
			value := v.(*models.TaskScraper)
			if value.GroupID == groupID {
				taskIDUUID, _ := uuid.Parse(string(value.TaskID))
				tasks.CancelTask(taskIDUUID)
				go communicator.TaskScraperObjectGMap.Remove(value.TaskID)
				go communicator.TaskLogsGMap.Remove(value.TaskID)
			}
			return true
		})
		done <- true
	}(done)

	go func(done chan bool) { // iterate through task worker object gmap
		communicator.TaskWorkerObjectGMap.Iterator(func(k, v interface{}) bool {
			value := v.(*models.TaskWorker)
			if value.GroupID == groupID {
				taskIDUUID, _ := uuid.Parse(string(value.TaskID))
				tasks.CancelTask(taskIDUUID)
				go communicator.TaskWorkerObjectGMap.Remove(value.TaskID)
				go communicator.TaskLogsGMap.Remove(value.TaskID)
			}
			return true
		})
		done <- true
	}(done)

	<-done
	<-done

	delete(communicator.Config.TaskGroups, groupID)
	r.Response.WriteJsonExit(map[string]bool{"success": true})
}
