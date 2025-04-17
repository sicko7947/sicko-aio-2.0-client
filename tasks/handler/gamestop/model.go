package gamestop

import (
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/utils/psychoclient"
)

type monitorPayload struct {
	taskID models.TaskID

	scrapes      int64
	inStock      bool
	online       bool
	readyToOrder bool
	headers      map[string]string

	scraper          *models.TaskScraper
	taskGroupSetting *models.TaskGroupSetting
}

type checkoutPayload struct {
	taskID models.TaskID

	headers map[string]string
	session psychoclient.Session

	worker           *models.TaskWorker
	taskGroupSetting *models.TaskGroupSetting
}
