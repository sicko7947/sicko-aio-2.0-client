package louisvuitton

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/constants"
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/tasks"
	"sicko-aio-2.0-client/utils"
	"sicko-aio-2.0-client/utils/psychoclient"
)

type monitorPayload struct {
	inStock   chan bool
	isCooling bool
	scrapes   int64

	headers map[string]string

	taskID models.TaskID

	scraper          *models.TaskScraper
	taskGroupSetting *models.TaskGroupSetting
}

func (p monitorPayload) doMonitor(sesh psychoclient.Session, errCh chan *models.Error) {
	defer recover()

	if sesh == nil {
		tasks.SafeSend(errCh, &models.Error{Error: errors.New("ERROR_GETTING_SESSION"), Code: 500, Message: "Error Getting Session"})
		return
	}

	endpoint := `https://xxxxxxxxxxxxxxxxxxx`

	data, _ := json.Marshal(map[string]interface{}{})

	reqID, _ := sesh.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "POST",
		Headers:  p.headers,
		Payload:  bytes.NewBuffer(data),
	})

	res, _, err := sesh.Do(reqID)
	if err != nil {
		tasks.SafeSend(errCh, err)
		return
	}
	fmt.Println(res)
	switch res.StatusCode {
	case 200:

		tasks.SafeSend(errCh, nil)
		return
	default:
		tasks.SafeSend(errCh, &models.Error{Error: errors.New("ERROR_GETTING_PRODUCT_INFORMATION"), Code: res.StatusCode, Message: "Error Getting Product Information"})
	}
}

// HandleLouisVuittonMonitorTaskWithOptions : HandleLouisVuittonMonitorTaskWithOptions
func HandleLouisVuittonMonitorTaskWithOptions(ctx context.Context, task *asynq.Task) error {
	// Get basic task info

	p := &monitorPayload{}
	p.scrapes = 0
	p.inStock = make(chan bool)
	p.taskID = models.TaskID(task.Payload())
	if scraper := communicator.TaskScraperObjectGMap.Get(p.taskID); scraper != nil {
		p.scraper = scraper.(*models.TaskScraper)
		if !p.scraper.Mutex.TryLock() {
			return errors.New("too many tasks assigned to one scraper")
		}
		defer func() {
			if p.scraper.Mutex.IsLocked() {
				p.scraper.Mutex.Unlock()
			}
		}()
	} else {
		return errors.New("error starting scraper")
	}

	// setup task group settings
	p.taskGroupSetting = communicator.Config.TaskGroups[p.scraper.GroupID].TaskGroupSetting

	p.headers = map[string]string{
		"accept":            "*/*",
		"accept-language":   "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7,es;q=0.6",
		"content-type":      "application/json",
		"dnt":               "1",
		"origin":            fmt.Sprintf("https://%s", p.taskGroupSetting.Domain),
		"qb-source-package": "@qubit/recommendations",
		"referer":           fmt.Sprintf("https://%s", p.taskGroupSetting.Domain),
		"user-agent":        gofakeit.RandomString(constants.ChromeUAList),
		"cache-control":     "no-cache",
	}

	// setup monitor proxy group
	// monitorProxyGroup := (communicator.Config.Proxies)[p.scraper.ScraperProxyGroupName]

	stopCh := make(chan *models.Error, 1)
	defer close(stopCh)
	for i := 0; i < p.scraper.ScraperNum; i++ {
		go func(ctx context.Context) {
			delay := time.NewTicker(time.Duration(p.scraper.ScraperDelay) * time.Millisecond)
			for {
				errCh := make(chan *models.Error, 1)
				sesh, err := utils.InitOrGetSession(p.scraper.ScraperProxyGroupName)
				if err != nil {
					continue
				}

				go p.doMonitor(sesh, errCh)

			L:
				for {
					select {
					case <-ctx.Done():
						return
					case err := <-errCh: // controlling errors
						if err != nil {
							communicator.ModifyTaskStatus(&models.Message{GroupID: p.scraper.GroupID, TaskID: p.taskID, Status: "PENDING", Code: err.Code, Message: err.Message})
						} else {
							communicator.ModifyTaskStatus(&models.Message{GroupID: p.scraper.GroupID, TaskID: p.taskID, Status: "PENDING", Code: 200, Message: "SCRAPING..."})
						}
						break L
					}
				}

				close(errCh)
				p.scraper.Product.Scrapes++
				communicator.TaskScraperObjectGMap.Set(p.taskID, p.scraper)
				utils.PutSession(p.scraper.ScraperProxyGroupName, sesh)
				<-delay.C
			}
		}(ctx)
	}

	// listener for instock product
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case isInStock := <-p.inStock:
				if p.isCooling {
					break
				}

				go func() {
					if isInStock {
						tasks.AddToQueue(p.scraper.GroupID, p.scraper.Product, p.scraper.TriggerNum)
						p.isCooling = true
						<-time.NewTicker(time.Duration(p.scraper.TriggerDelay) * time.Millisecond).C
						p.isCooling = false
					}
				}()
			}
		}
	}(ctx)

	select {
	case <-ctx.Done():
		communicator.ModifyTaskStatus(&models.Message{ // send task cancellation status to frontend
			Code:    0,
			GroupID: p.scraper.GroupID,
			TaskID:  p.taskID,
			Status:  "CANCELLED",
			Message: "Stopped",
		})
		return ctx.Err() // cancelation signal received, abandon this work.
	}
}
