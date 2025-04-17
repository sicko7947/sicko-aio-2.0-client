package macys

import (
	"context"
	"errors"
	"time"

	"github.com/hibiken/asynq"
	"github.com/sicko7947/sickocommon"
	"github.com/tidwall/gjson"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/tasks"
	"sicko-aio-2.0-client/utils"
	"sicko-aio-2.0-client/utils/psychoclient"
)

func (p *monitorPayload) doGetInformation(sesh psychoclient.Session, errCh chan *models.Error) {
	endpoint := "https://www.macys.com/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

	reqID, _ := sesh.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "GET",
		Headers:  p.headers,
		Payload:  nil,
	})

	res, _, err := sesh.Do(reqID)
	if err != nil {
		tasks.SafeSend(errCh, err)
		return
	}
	switch res.StatusCode {
	case 200:
		// Load the HTML document

		tasks.SafeSend(errCh, &models.Error{Error: nil, Code: 404, Message: "Waiting For Product"})
	default:
		tasks.SafeSend(errCh, &models.Error{Error: errors.New("ERROR_GETTING_PRODUCT"), Code: res.StatusCode, Message: "Error Getting Product"})
		return
	}
}

func (p *monitorPayload) doMonitor(sesh psychoclient.Session, errCh chan *models.Error) {
	if p.scraper.Product == nil || len(p.scraper.Product.ProductID) == 0 {
		tasks.SafeSend(errCh, &models.Error{Error: errors.New("ERROR_WAITING_FOR_RESTOCK"), Code: 404, Message: "Product Not Found"})
		return
	}

	endpoint := "https://www.macys.com/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" + p.scraper.Product.ProductID

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
	result := gjson.ParseBytes(respBody)
	switch res.StatusCode {
	case 200:
		if productObj := result.Get("product.0"); productObj.Exists() {
			p.scraper.Product.Status = productObj.Get("status").String()
			p.scraper.Product.ProductURL = "https://www.macys.com" + productObj.Get("identifier.productUrl").String()
			p.scraper.Product.ProductID = productObj.Get("identifier.productId").String()
			p.scraper.Product.ProductName = productObj.Get("detail.name").String()
			p.scraper.Product.ProductDescription = productObj.Get("detail.bulletText.7").String()
			p.scraper.Product.QuantityLimit = int(productObj.Get("detail.maxQuantity").Int())

			tempMap := make(map[string]string)
			productObj.Get("relationships.upcs").ForEach(func(key, value gjson.Result) bool {
				if value.Get("availability.available").Bool() {
					tempMap[value.Get("traits.sizes.selectedSize").String()] = key.String()
				}
				return true
			})

			sizeMap := make(map[string]string)
			productObj.Get("traits.sizes.sizeMap").ForEach(func(key, value gjson.Result) bool {
				sizeMap[value.Get("name").String()] = tempMap[key.String()]
				return true
			})

			if len(sizeMap) > 0 {
				newMap := make(map[string]*models.SizeSkuMap)
				for k, v := range sizeMap {
					if sickocommon.CheckSliceContains(p.scraper.DesireSizes, "RA") {
						newMap[k] = &models.SizeSkuMap{
							SkuId: v,
						}

					} else if sickocommon.CheckSliceContains(p.scraper.DesireSizes, k) {
						newMap[k] = &models.SizeSkuMap{
							SkuId: v,
						}
					}
				}

				if len(newMap) > 0 {
					p.inStock <- true
				} else {
					p.inStock <- false
				}
				p.scraper.Product.SizeSkuMap = newMap
			}

			tasks.SafeSend(errCh, nil)
			return
		}
		fallthrough
	default:
		tasks.SafeSend(errCh, &models.Error{Error: errors.New("ERROR_WAITING_FOR_RESTOCK"), Code: res.StatusCode, Message: "Error Waiting For Restock"})
	}
}

// HandleMacysMonitorTaskWithOptions : HandleMacysMonitorTaskWithOptions
func HandleMacysMonitorTaskWithOptions(ctx context.Context, task *asynq.Task) error {
	// Get basic task info

	p := &monitorPayload{
		inStock: make(chan bool),
		taskID:  models.TaskID(task.Payload()),
	}
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
		"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"accept-language":           "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7,es;q=0.6",
		"cache-control":             "no-cache",
		"dnt":                       "1",
		"sec-fetch-dest":            "document",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-site":            "same-origin",
		"sec-fetch-user":            "?1",
		"upgrade-insecure-requests": "1",
		"user-agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36",
	}

	stopCh := make(chan *models.Error, 1)
	defer close(stopCh)
	for i := 0; i < p.scraper.ScraperNum; i++ {
		go func(ctx context.Context) {
			delay := time.NewTicker(time.Duration(p.scraper.ScraperDelay) * time.Millisecond)
			for {
				errCh := make(chan *models.Error, 3)

				sesh, err := utils.InitOrGetSession(p.scraper.ScraperProxyGroupName)
				if err != nil {
					continue
				}

				go p.doGetInformation(sesh, errCh)
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
		ticker := time.NewTicker(time.Duration(p.scraper.TriggerDelay) * time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				return
			case instock := <-p.inStock:
				if instock {
					tasks.AddToQueue(p.scraper.GroupID, p.scraper.Product, p.scraper.TriggerNum)
					<-ticker.C
				}
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
