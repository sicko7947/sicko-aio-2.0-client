package macys

import (
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/utils/psychoclient"
)

type checkoutPayload struct {
	taskID models.TaskID

	checkoutID string
	headers    map[string]string
	session    psychoclient.Session

	taskGroupSetting *models.TaskGroupSetting
	worker           *models.TaskWorker
}

type monitorPayload struct {
	taskID models.TaskID

	scrapes int64
	inStock chan bool

	headers map[string]string

	scraper          *models.TaskScraper
	taskGroupSetting *models.TaskGroupSetting
}
