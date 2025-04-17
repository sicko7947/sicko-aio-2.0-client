package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gogf/gf/net/ghttp"
	"github.com/hibiken/asynq"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/tasks"
	"sicko-aio-2.0-client/tasks/defineder"
)

var noDeadline = time.Now().AddDate(3, 0, 0)

type starttaskId struct {
	GroupID   models.TaskGroupID `json:"groupId"`
	TaskID    models.TaskID      `json:"taskId"`
	StartTime string             `json:"startTime,omitempty"`
}

// StartScraperTasks : Start Scraper Task
func StartScraperTasks(r *ghttp.Request) {

	// Parsing Payload
	var p []starttaskId
	if err := r.Parse(&p); err != nil {
		r.Response.WriteStatus(http.StatusBadRequest)
		r.Response.WriteJsonExit(map[string]bool{"success": false})
		return
	}

	for _, payload := range p {
		taskGroupSettings := communicator.Config.TaskGroups[payload.GroupID].TaskGroupSetting
		scraper := (communicator.TaskScraperObjectGMap.Get(payload.TaskID)).(*models.TaskScraper)
		if scraper == nil || scraper.Mutex.IsLocked() {
			continue
		}

		switch taskGroupSettings.Category {
		case models.NIKE, models.SNKRS:
			switch scraper.MonitorMode {
			case "V1":
				nikeMonitorTask := defineder.NewNikeMonitorGtinTask(payload.TaskID)
				if _, err := tasks.Client.Enqueue(nikeMonitorTask,
					asynq.MaxRetry(-1),
					asynq.Deadline(noDeadline),
				); err != nil {
					log.Fatal(err)
				}
			case "V2":
				nikeMonitorTask := defineder.NewnNikeMonitorGraphQLTask(payload.TaskID)
				if _, err := tasks.Client.Enqueue(nikeMonitorTask,
					asynq.MaxRetry(-1),
					asynq.Deadline(noDeadline),
				); err != nil {
					log.Fatal(err)
				}
			case "GRPC":
				nikeMonitorTask := defineder.NewNikeMonitorDeliveryTask(payload.TaskID)
				if _, err := tasks.Client.Enqueue(nikeMonitorTask,
					asynq.MaxRetry(-1),
					asynq.Deadline(noDeadline),
				); err != nil {
					log.Fatal(err)
				}
			}

		case models.ADIDAS:
			adidasTask := defineder.NewAdidasTask(payload.TaskID)
			if _, err := tasks.Client.Enqueue(adidasTask,
				asynq.MaxRetry(-1),
				asynq.Deadline(noDeadline),
			); err != nil {
				log.Fatal(err)
			}

		case models.MACYS:
			macysMonitorTask := defineder.NewMacysMonitorTask(payload.TaskID)
			if _, err := tasks.Client.Enqueue(macysMonitorTask,
				asynq.MaxRetry(-1),
				asynq.Deadline(noDeadline),
			); err != nil {
				log.Fatal(err)
			}

		case models.MRPORTER:
			mrporterMonitorTask := defineder.NewMrPorterMonitorTask(payload.TaskID)
			if _, err := tasks.Client.Enqueue(mrporterMonitorTask,
				asynq.MaxRetry(-1),
				asynq.Deadline(noDeadline),
			); err != nil {
				log.Fatal(err)
			}

		case models.LUISAVIAROMA:
			luisaviaromaMonitorTask := defineder.NewLuisaviaromaMonitorTask(payload.TaskID)
			if _, err := tasks.Client.Enqueue(luisaviaromaMonitorTask,
				asynq.MaxRetry(-1),
				asynq.Deadline(noDeadline),
			); err != nil {
				log.Fatal(err)
			}

		case models.LANECRAWFORD:
			lanecrawfordMonitorTask := defineder.NewLanecrawfordMonitorTask(payload.TaskID)
			if _, err := tasks.Client.Enqueue(lanecrawfordMonitorTask,
				asynq.MaxRetry(-1),
				asynq.Deadline(noDeadline),
			); err != nil {
				log.Fatal(err)
			}

		case models.LOUISVUITTON:
			louisVuittonMonitorTask := defineder.NewLouisVuittonMonitorTask(payload.TaskID)
			if _, err := tasks.Client.Enqueue(louisVuittonMonitorTask,
				asynq.MaxRetry(-1),
				asynq.Deadline(noDeadline),
			); err != nil {
				log.Fatal(err)
			}

		case models.GAMESTOP:
			gamestopMonitorTask := defineder.NewGamestopMonitorTask(payload.TaskID)
			if _, err := tasks.Client.Enqueue(gamestopMonitorTask,
				asynq.MaxRetry(-1),
				asynq.Deadline(noDeadline),
			); err != nil {
				log.Fatal(err)
			}

		}
	}
	r.Response.WriteJsonExit(map[string]bool{"success": true})

}

// StartWorkerTasks : Start Worker Tasks
func StartWorkerTasks(r *ghttp.Request) {

	// Parsing Payload
	var p []*starttaskId
	if err := r.Parse(&p); err != nil {
		r.Response.WriteStatus(http.StatusBadRequest)
		r.Response.WriteJsonExit(map[string]bool{"success": false})
		return
	}

	for _, payload := range p {
		taskGroupSettings := communicator.Config.TaskGroups[payload.GroupID].TaskGroupSetting
		worker := (communicator.TaskWorkerObjectGMap.Get(payload.TaskID)).(*models.TaskWorker)
		if worker == nil || worker.Mutex.IsLocked() {
			continue
		}

		switch taskGroupSettings.Category {
		case models.NIKE:
			worker.TaskInfo = nil
			worker.Product = nil

			nikeCheckoutPrepareTask := defineder.NewNikeCheckoutPrepare(payload.TaskID) // start task preperation first, task will get invoked
			if _, err := tasks.Client.Enqueue(nikeCheckoutPrepareTask,
				asynq.MaxRetry(5),
				asynq.Deadline(noDeadline),
			); err != nil {
				log.Fatal(err)
			}

		case models.SNKRS:
			worker.TaskInfo = nil
			worker.Product = nil

			snkrsLaunchEntryPrepareTask := defineder.NewSnkrsLaunchEntryPrepare(payload.TaskID)
			if _, err := tasks.Client.Enqueue(snkrsLaunchEntryPrepareTask,
				asynq.MaxRetry(5),
				asynq.Deadline(noDeadline),
			); err != nil {
				log.Fatal(err)
			}

		case models.MACYS:
			worker.TaskInfo = nil
			worker.Product = nil

			macysCheckoutPrepareTask := defineder.NewMacysCheckoutPrepareTask(payload.TaskID)
			if _, err := tasks.Client.Enqueue(macysCheckoutPrepareTask,
				asynq.MaxRetry(999),
				asynq.Deadline(noDeadline),
			); err != nil {
				log.Fatal(err)
			}

		case models.ADIDAS:
			adidasTask := defineder.NewAdidasTask(payload.TaskID)
			if _, err := tasks.Client.Enqueue(adidasTask,
				asynq.MaxRetry(999),
				asynq.Deadline(noDeadline),
			); err != nil {
				log.Fatal(err)
			}

		case models.GAMESTOP:
			gamestopCheckoutTask := defineder.NewGamestopCheckoutTask(payload.TaskID)
			if _, err := tasks.Client.Enqueue(gamestopCheckoutTask,
				asynq.MaxRetry(999),
				asynq.Deadline(noDeadline),
			); err != nil {
				log.Fatal(err)
			}
		case models.MRPORTER:
			worker.TaskInfo = nil
			worker.Product = nil

			mrporterCheckoutPrepareTask := defineder.NewMrPorterCheckoutPrepare(payload.TaskID)
			if _, err := tasks.Client.Enqueue(mrporterCheckoutPrepareTask,
				asynq.MaxRetry(999),
				asynq.Deadline(noDeadline),
			); err != nil {
				log.Fatal(err)
			}

		case models.LOUISVUITTON:
			worker.TaskInfo = nil
			worker.Product = nil

			louisVuittonCheckoutPrepareTask := defineder.NewMrPorterCheckoutPrepare(payload.TaskID)
			if _, err := tasks.Client.Enqueue(louisVuittonCheckoutPrepareTask,
				asynq.MaxRetry(999),
				asynq.Deadline(noDeadline),
			); err != nil {
				log.Fatal(err)
			}

		case models.LUISAVIAROMA:
			worker.TaskInfo = nil
			worker.Product = nil

			luisaviaromaCheckoutPrepareTask := defineder.NewLuisaviaromaCheckoutPrepare(payload.TaskID)
			if _, err := tasks.Client.Enqueue(luisaviaromaCheckoutPrepareTask,
				asynq.Deadline(noDeadline),
			); err != nil {
				log.Fatal(err)
			}

		case models.LANECRAWFORD:
			worker.TaskInfo = nil
			worker.Product = nil

			lanecrawfordCheckoutPrepareTask := defineder.NewLanecrawfordCheckoutPrepareTask(payload.TaskID)
			if _, err := tasks.Client.Enqueue(lanecrawfordCheckoutPrepareTask,
				asynq.MaxRetry(5),
				asynq.Deadline(noDeadline),
			); err != nil {
				log.Fatal(err)
			}

		}
	}

	r.Response.WriteJsonExit(map[string]bool{"success": true})
}
