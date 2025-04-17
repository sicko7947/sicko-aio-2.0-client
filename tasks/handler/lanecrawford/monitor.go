package lanecrawford

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/hibiken/asynq"
	"github.com/sicko7947/sickocommon"
	"github.com/tidwall/gjson"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/constants"
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/tasks"
	"sicko-aio-2.0-client/utils"
	"sicko-aio-2.0-client/utils/psychoclient"
)

func (p *monitorPayload) doMonitor(sesh psychoclient.Session, errCh chan *models.Error) {
	if sesh == nil {
		tasks.SafeSend(errCh, &models.Error{Error: errors.New("ERROR_GETTING_SESSION"), Code: 500, Message: "Error Getting Session"})
		return
	}

	endpoint := "https://www.lanecrawford.com.hk/xxxxxxxxxxxxxxxxxxxxxxxx"

	form := url.Values{}
	data := map[string]string{}
	for key, value := range data {
		form.Set(key, value)
	}

	reqId, _ := sesh.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "POST",
		Headers: map[string]string{
			"accept":           "*/*",
			"accept-language":  "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7,es;q=0.6",
			"content-type":     "application/x-www-form-urlencoded",
			"dnt":              "1",
			"user-agent":       p.headers["user-agent"],
			"x-requested-with": "XMLHttpRequest",
			"cache-control":    "no-cache",
		},
		Payload: strings.NewReader(form.Encode()),
	})

	res, respBody, err := sesh.Do(reqId)
	if err != nil {
		tasks.SafeSend(errCh, err)
		return
	}

	switch res.StatusCode {
	case 200:
		result := gjson.Get(string(respBody), "atgResponse")

		p.scraper.Product.ProductName = result.Get("essentialInfo.displayName").String()
		p.scraper.Product.ImageURL = result.Get("essentialInfo.heroImage").String()
		p.scraper.Product.Price = "HK$" + result.Get("priceInfo.hkSalePrice").String()

		if !result.Get("isInStoreOnly").Bool() && !result.Get("isCollectAtStoreOnly").Bool() {

			p.scraper.Product.Status = "ACTIVE"
			newMap := make(map[string]*models.SizeSkuMap)
			result.Get("sizeSwatchList.#(hasStock==true)#").ForEach(func(key, value gjson.Result) bool {
				size := value.Get("size").String()
				skuId := value.Get("skuId").String()
				if sickocommon.CheckSliceContains(p.scraper.DesireSizes, "RA") {
					newMap[size] = &models.SizeSkuMap{
						SkuId: skuId,
					}
				} else if sickocommon.CheckSliceContains(p.scraper.DesireSizes, size) {
					newMap[size] = &models.SizeSkuMap{
						SkuId: skuId,
					}
				}
				return true
			})

			p.scraper.Product.SizeSkuMap = newMap
			if len(newMap) > 0 {
				p.inStock <- true
			} else {
				p.inStock <- false
			}
		}
		tasks.SafeSend(errCh, nil)
	default:
		tasks.SafeSend(errCh, &models.Error{Error: errors.New("ERROR_WAITING_FOR_RESTOCK"), Code: res.StatusCode, Message: "Error Waiting For Restock"})
	}
}

// HandleLanecrawfordMonitorTaskWithOptions : HandleLanecrawfordMonitorTaskWithOptions
func HandleLanecrawfordMonitorTaskWithOptions(ctx context.Context, task *asynq.Task) error {
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

	// init session pool
	// monitorProxyGroup := (communicator.Config.Proxies)[p.scraper.ScraperProxyGroupName]

	// setup request headers
	p.headers = map[string]string{
		"accept":           "application/json, text/javascript, */*; q=0.01",
		"accept-language":  "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7,es;q=0.6",
		"dnt":              "1",
		"sec-fetch-dest":   "empty",
		"sec-fetch-mode":   "cors",
		"sec-fetch-site":   "same-origin",
		"x-requested-with": "XMLHttpRequest",
		"cache-control":    "no-cache",
		"user-agent":       gofakeit.RandomString(constants.ChromeUAList),
	}

	communicator.ModifyTaskStatus(&models.Message{GroupID: p.scraper.GroupID, TaskID: p.taskID, Status: "PENDING", Code: 200, Message: "SCRAPING..."})
	// start scrapers according to scraper number set by user
	for i := 0; i < p.scraper.ScraperNum; i++ {
		go func(ctx context.Context) {
			delay := time.NewTicker(time.Duration(p.scraper.ScraperDelay) * time.Millisecond)
			for {
				errCh := make(chan *models.Error, 4)

				sesh, err := utils.InitOrGetSession(p.scraper.ScraperProxyGroupName)
				if err != nil {
					continue
				}

				// go p.doProductInformation(sesh, errCh)
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

	<-ctx.Done()
	communicator.ModifyTaskStatus(&models.Message{ // send task cancellation status to frontend
		GroupID: p.scraper.GroupID,
		TaskID:  p.taskID,
		Status:  "CANCELLED",
		Message: "Stopped",
	})
	return ctx.Err() // cancelation signal received, abandon this work.
}
