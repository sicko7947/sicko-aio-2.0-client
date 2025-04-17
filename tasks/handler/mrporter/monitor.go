package mrporter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gogf/gf/util/gconv"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"github.com/sicko7947/sickocommon"
	"github.com/tidwall/gjson"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/tasks"
	"sicko-aio-2.0-client/utils/psychoclient"
)

type monitorPayload struct {
	inStock   chan bool
	isCooling bool

	headers map[string]string

	taskID models.TaskID

	scraper          *models.TaskScraper
	taskGroupSetting *models.TaskGroupSetting
}

func (p *monitorPayload) doMonitor(sesh psychoclient.Session, errCh chan *models.Error) {

	endpoint := fmt.Sprintf(`https://xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx_%s`, p.taskGroupSetting.Locale, p.scraper.Product.StyleColor, p.taskGroupSetting.Country)

	data, _ := json.Marshal(make(map[string]string))
	reqID, _ := sesh.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "POST",
		Headers:  p.headers,
		Payload:  bytes.NewBuffer(data),
	})

	res, respBody, err := sesh.Do(reqID)
	if err != nil {
		tasks.SafeSend(errCh, err)
		return
	}
	result := gjson.Get(string(respBody), "products.0")
	switch res.StatusCode {
	case 200:
		// "productColours.visible"  might be the checklive option ?
		p.scraper.Product.ImageURL = fmt.Sprintf("https://xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", p.scraper.Product.StyleColor)

		if productName := result.Get("name"); productName.Exists() {
			p.scraper.Product.ProductName = productName.String()
		}

		if price := result.Get("price.sellingPrice.amount"); price.Exists() {
			p.scraper.Product.Price = result.Get("price.currency.symbol").String() + gconv.String(price.Int()/100)
		}

		if skusObj := result.Get("productColours.0.sKUs"); skusObj.Exists() {
			sizeMap := make(map[string]*models.SizeSkuMap)
			isRA := sickocommon.CheckSliceContains(p.scraper.DesireSizes, "RA")
			skusObj.ForEach(func(key, value gjson.Result) bool {
				if value.Get("buyable").Bool() {
					size := value.Get(`size.schemas.#(name=="US").labels.0`)
					if !size.Exists() {
						size = value.Get("size.centralSizeLabel")
					}

					sizeString := size.String()
					if isRA {
						sizeMap[sizeString] = &models.SizeSkuMap{
							SkuId: value.Get("partNumber").String(),
						}
					} else if sickocommon.CheckSliceContains(p.scraper.DesireSizes, sizeString) {
						sizeMap[sizeString] = &models.SizeSkuMap{
							SkuId: value.Get("partNumber").String(),
						}
					}
				}
				return true
			})

			p.scraper.Product.SizeSkuMap = sizeMap
			if len(sizeMap) > 0 {
				p.inStock <- true
			} else {
				p.inStock <- false
			}
			tasks.SafeSend(errCh, nil)
			return
		}
		tasks.SafeSend(errCh, &models.Error{Error: nil, Code: res.StatusCode, Message: "Stock Not Live"})
	default:
		tasks.SafeSend(errCh, &models.Error{Error: errors.New("ERROR_WAITING_FOR_RESTOCK"), Code: res.StatusCode, Message: "Error Waiting For Restock"})
	}
}

// HandleMrporterMonitorTaskWithOptions : HandleMrporterMonitorTaskWithOptions
func HandleMrporterMonitorTaskWithOptions(ctx context.Context, task *asynq.Task) error {

	// set base monitorPayload
	p := &monitorPayload{
		inStock: make(chan bool),
		taskID:  models.TaskID(task.Payload()),
	}

	defer func() {
		if p.scraper != nil && p.scraper.Mutex.IsLocked() {
			p.scraper.Mutex.Unlock()
		}
	}()

	scraper := communicator.TaskScraperObjectGMap.Get(p.taskID)
	if scraper == nil {
		return errors.New("error starting scraper")
	}

	p.scraper = scraper.(*models.TaskScraper)
	if !p.scraper.Mutex.TryLock() {
		return errors.New("too many tasks assigned to one scraper")
	}

	p.scraper.Product = &models.ScraperTaskProduct{
		StyleColor: p.scraper.Product.StyleColor,
	}

	// setup task group settings
	p.taskGroupSetting = communicator.Config.TaskGroups[p.scraper.GroupID].TaskGroupSetting

	// setup scraper proxy group
	scraperProxyGroup := (communicator.Config.Proxies)[p.scraper.ScraperProxyGroupName]

	// setup request headers
	p.headers = map[string]string{
		"accept":              "*/*",
		"accept-language":     "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7,es;q=0.6",
		"application-name":    "myaccount",
		"application-version": "5.443.0",
		"content-type":        "application/json",
		"label":               "login",
		"sec-ch-ua-mobile":    "?0",
		"sec-fetch-dest":      "empty",
		"sec-fetch-mode":      "cors",
		"sec-fetch-site":      "same-origin",
		"cache-control":       "no-cache",
	}

	stopCh := make(chan *models.Error, 1)
	defer close(stopCh)
	for i := 0; i < p.scraper.ScraperNum; i++ {
		go func(ctx context.Context) {
			for {
				errCh := make(chan *models.Error)

				sesh, _ := psychoclient.NewSession(&psychoclient.SessionBuilder{
					UseDefaultClient: true,
					Proxy:            sickocommon.GetProxy(scraperProxyGroup).String(),
				})

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
				sesh.Close()
				p.scraper.Product.Scrapes++
				communicator.TaskScraperObjectGMap.Set(p.taskID, p.scraper)
				<-time.NewTicker(time.Duration(p.scraper.ScraperDelay) * time.Millisecond).C
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
			GroupID: p.scraper.GroupID,
			TaskID:  p.taskID,
			Status:  "CANCELLED",
			Message: "Stopped",
		})
		return ctx.Err() // cancelation signal received, abandon this work.
	}
}
