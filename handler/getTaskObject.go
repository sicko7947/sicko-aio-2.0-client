package handler

import (
	"github.com/gogf/gf/net/ghttp"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/models"
)

func GetTaskObject(r *ghttp.Request) {
	groupID := models.TaskGroupID(r.Get("groupId").(string))

	obj := &models.TaskGroup{ // init obj
		GroupID:          groupID,
		TaskScrapers:     make(map[models.TaskID]*models.TaskScraper),
		TaskWorkers:      make(map[models.TaskID]*models.TaskWorker),
		TaskGroupSetting: communicator.Config.TaskGroups[groupID].TaskGroupSetting,
	}

	done := make(chan bool, 2)
	go func(done chan bool) { // iterate through task scraper object gmap
		communicator.TaskScraperObjectGMap.Iterator(func(k, v interface{}) bool {
			value := v.(*models.TaskScraper)
			if value.GroupID == groupID {
				obj.TaskScrapers[k.(models.TaskID)] = value
			}
			return true
		})
		done <- true
	}(done)

	go func(done chan bool) { // iterate through task worker object gmap
		communicator.TaskWorkerObjectGMap.Iterator(func(k, v interface{}) bool {
			value := v.(*models.TaskWorker)
			if value.GroupID == groupID {
				obj.TaskWorkers[k.(models.TaskID)] = value
			}
			return true
		})
		done <- true
	}(done)

	// wait for both goroutine to finish
	<-done
	<-done

	r.Response.WriteJsonExit(obj)
}
