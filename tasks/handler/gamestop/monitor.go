package gamestop

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/hibiken/asynq"
	"github.com/tidwall/gjson"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/tasks"
	"sicko-aio-2.0-client/utils"
	"sicko-aio-2.0-client/utils/psychoclient"
)

func (p *monitorPayload) doMonitor(sesh psychoclient.Session, errCh chan *models.Error) {
	endpoint := fmt.Sprintf(`https://www.gamestop.com/xxxxxxxxxxxxx/%v`, gofakeit.LetterN(20))

	reqID, _ := sesh.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "GET",
		Headers:  p.headers,
		Payload:  nil,
	})

	res, respBody, err := sesh.Do(reqID)
	if err != nil {
		tasks.SafeSend(errCh, err)
		return
	}
	result := gjson.Get(string(respBody), "product")
	switch res.StatusCode {
	case 200:
		p.online = result.Get("online").Bool()

		tasks.SafeSend(errCh, nil)
	default:
		tasks.SafeSend(errCh, &models.Error{Error: nil, Code: 500, Message: "ERROR_GETTING_PRODUCT"})
	}
}

// HandleGamestopMonitorTaskWithOptions : HandleGamestopMonitorTaskWithOptions
func HandleGamestopMonitorTaskWithOptions(ctx context.Context, task *asynq.Task) error {
	// Get basic task info

	p := &monitorPayload{}
	p.taskID = models.TaskID(task.Payload())
	if scraper := communicator.TaskScraperObjectGMap.Get(p.taskID); scraper != nil {
		p.scraper = scraper.(*models.TaskScraper)
		if !p.scraper.Mutex.TryLock() {
			return errors.New("too many tasks assigned to one worker")
		}
		defer func() {
			if p.scraper.Mutex.IsLocked() {
				p.scraper.Mutex.Unlock()
			}
		}()
	} else {
		return errors.New("error starting worker")
	}

	// setup task group settings
	p.taskGroupSetting = communicator.Config.TaskGroups[p.scraper.GroupID].TaskGroupSetting

	// init session pool
	// monitorProxyGroup := (communicator.Config.Proxies)[p.scraper.ScraperProxyGroupName]

	p.headers = map[string]string{
		"accept":           "*/*",
		"accept-language":  "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7,es;q=0.6",
		"dnt":              "1",
		"referer":          "https://www.gamestop.com",
		"sec-ch-ua-mobile": "?0",
		"sec-fetch-dest":   "empty",
		"sec-fetch-mode":   "cors",
		"sec-fetch-site":   "same-origin",
		"user-agent":       "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36",
		"x-requested-with": "XMLHttpRequest",
		"cache-control":    "no-cache",
	}

	communicator.ModifyTaskStatus(&models.Message{GroupID: p.scraper.GroupID, TaskID: p.taskID, Status: "PENDING", Message: "SCRAPING..."})

	stopCh := make(chan *models.Error, 1)
	defer close(stopCh)
	for i := 0; i < p.scraper.ScraperNum; i++ {
		go func() {
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
				p.scrapes++
				p.scraper.Product.Scrapes = p.scrapes
				communicator.TaskScraperObjectGMap.Set(p.taskID, p.scraper)
				utils.PutSession(p.scraper.ScraperProxyGroupName, sesh)
				<-delay.C
			}
		}()
	}

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if p.inStock && p.online && p.readyToOrder {
					tasks.AddToQueue(p.scraper.GroupID, p.scraper.Product, p.scraper.TriggerNum)
					<-time.NewTicker(time.Duration(p.scraper.TriggerDelay) * time.Millisecond).C
				}
			}
		}
	}(ctx)

	<-ctx.Done()
	communicator.ModifyTaskStatus(&models.Message{ // send task cancellation status to frontend
		GroupID: p.scraper.GroupID,
		TaskID:  p.taskID,
		Status:  "CANCELLED",
		Message: "Stopped",
	})
	return ctx.Err() // cancelation signal received, abandon this work.
}
