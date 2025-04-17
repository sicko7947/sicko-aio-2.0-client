package lanecrawford

import (
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/utils/psychoclient"
)

type monitorPayload struct {
	taskID models.TaskID

	inStock   chan bool
	isCooling bool

	scrapes int64
	headers map[string]string

	scraper          *models.TaskScraper
	taskGroupSetting *models.TaskGroupSetting
}

type checkoutPayload struct {
	taskID models.TaskID

	_dynSessConf string

	headers map[string]string
	session psychoclient.Session

	worker           *models.TaskWorker
	taskGroupSetting *models.TaskGroupSetting
}
