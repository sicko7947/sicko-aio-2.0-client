package monitor

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"github.com/shimingyah/pool"
	"github.com/sicko7947/sickocommon"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/models"
	grpc_service "sicko-aio-2.0-client/proto/rpc"
	"sicko-aio-2.0-client/tasks"
	"sicko-aio-2.0-client/utils/grpc"
)

type monitorPayload struct {
	taskID models.TaskID

	inStock           chan bool
	isCooling         bool
	monitorProxyGroup []string

	scraper          *models.TaskScraper
	taskGroupSetting models.TaskGroupSetting

	conn                     pool.Conn
	streamClient             grpc_service.StreamClient
	productInformationStream grpc_service.Stream_ProductInformationClient
	productCheckliveStream   grpc_service.Stream_ProductCheckLiveClient
	productMonitorStream     grpc_service.Stream_ProductMonitorClient
}

func (p *monitorPayload) doGetMonitorInformation(c chan *models.Error) {
	for {
		if p.productInformationStream == nil || p.productInformationStream.Context().Err() != nil {
			if closed := tasks.SafeSend(c, &models.Error{Error: errors.New("GET_STREAM_ERROR"), Code: 409, Message: "ProductInformation get stream err"}); closed {
				return
			}
			continue
		}

		res, err := p.productInformationStream.Recv()
		if err != nil {
			p.productInformationStream.Context().Done()
			if closed := tasks.SafeSend(c, &models.Error{Error: errors.New("GET_STREAM_ERROR"), Code: 409, Message: fmt.Sprintf("ProductInformation get stream err: %v", err)}); closed {
				return
			}
			continue
		}

		e := res.GetErrors()
		if e != nil {
			if closed := tasks.SafeSend(c, &models.Error{Error: errors.New("ERROR_GETTING_PRODUCT_INFORMATION"), Code: int(e.GetCode()), Message: e.GetMessage()}); closed {
				return
			}
		}

		if len(res.GetObjects()) > 0 {
			if len(p.scraper.Product.Status) == 0 {
				p.scraper.Product.Status = res.GetObjects()[p.scraper.Product.StyleColor].GetStatus()
			}

			if len(p.scraper.Product.ProductName) == 0 {
				p.scraper.Product.ProductName = res.GetObjects()[p.scraper.Product.StyleColor].GetProductName()
			}

			if len(p.scraper.Product.ProductDescription) == 0 {
				p.scraper.Product.ProductDescription = res.GetObjects()[p.scraper.Product.StyleColor].GetProductDescription()
			}

			if len(p.scraper.Product.ProductID) == 0 {
				p.scraper.Product.ProductID = res.GetObjects()[p.scraper.Product.StyleColor].GetProductId()
			}

			if p.scraper.Product.QuantityLimit == 0 {
				p.scraper.Product.QuantityLimit = int(res.GetObjects()[p.scraper.Product.StyleColor].GetQuantityLimit())
			}

			if len(p.scraper.Product.LaunchID) == 0 {
				p.scraper.Product.LaunchID = res.GetObjects()[p.scraper.Product.StyleColor].GetStatus()
			}

			if len(p.scraper.Product.Price) == 0 {
				p.scraper.Product.Price = res.GetObjects()[p.scraper.Product.StyleColor].GetPrice()
			}

			if len(p.scraper.Product.PublishType) == 0 {
				p.scraper.Product.PublishType = res.GetObjects()[p.scraper.Product.StyleColor].GetPublishType()
			}

			if p.scraper.Product.SizeSkuMap == nil {
				sizeSkus := res.GetObjects()[p.scraper.Product.StyleColor].GetSizeSkus()

				sizes := make(map[string]*models.SizeSkuMap)
				for k, v := range sizeSkus {
					sizes[k] = &models.SizeSkuMap{
						SkuId: v.SkuId,
						Gtin:  v.Gtin,
					}
				}
				p.scraper.Product.SizeSkuMap = sizes
			}
		}
	}
}

func (p *monitorPayload) doGetCheckLive(c chan *models.Error) {
	for {
		if p.productCheckliveStream == nil || p.productCheckliveStream.Context().Err() != nil {
			if closed := tasks.SafeSend(c, &models.Error{Error: errors.New("GET_STREAM_ERROR"), Code: 409, Message: "ProductCheckLive get stream err"}); closed {
				return
			}
			continue
		}

		res, err := p.productCheckliveStream.Recv()
		if err != nil {
			p.productCheckliveStream.Context().Done()
			if closed := tasks.SafeSend(c, &models.Error{Error: errors.New("GET_STREAM_ERROR"), Code: 409, Message: fmt.Sprintf("ProductCheckLive get stream err: %v", err)}); closed {
				return
			}
			continue
		}

		if e := res.GetErrors(); e != nil {
			if closed := tasks.SafeSend(c, &models.Error{Error: errors.New("ERROR_CHECKING_PRODUCT_LIVE"), Code: int(e.GetCode()), Message: e.GetMessage()}); closed {
				return
			}
		} else {
			p.scraper.Product.IsLive = res.GetLive()

			if status := res.GetStatus(); len(status) > 0 {
				p.scraper.Product.Status = status
			}
			if quantityLimit := res.GetQuantityLimit(); quantityLimit > 0 {
				p.scraper.Product.QuantityLimit = int(quantityLimit)
			}
			if publishType := res.GetPublishType(); len(publishType) > 0 {
				p.scraper.Product.PublishType = publishType
			}
		}
	}
}

func (p *monitorPayload) doGetMonitor(c chan *models.Error) {
	for {
		if p.productMonitorStream == nil || p.productMonitorStream.Context().Err() != nil {
			if closed := tasks.SafeSend(c, &models.Error{Error: errors.New("GET_STREAM_ERROR"), Code: 409, Message: "ProductMonitor get stream err"}); closed {
				return
			}
			continue
		}

		res, err := p.productMonitorStream.Recv()
		if err != nil {
			p.productMonitorStream.Context().Done()
			if closed := tasks.SafeSend(c, &models.Error{Error: errors.New("GET_STREAM_ERROR"), Code: 409, Message: fmt.Sprintf("ProductMonitor get stream err: %v", err)}); closed {
				return
			}
			continue
		}

		if e := res.GetErrors(); e != nil {
			if closed := tasks.SafeSend(c, &models.Error{Error: errors.New("ERROR_WAITING_FOR_RESTOCK"), Code: int(e.GetCode()), Message: e.GetMessage()}); closed {
				return
			}
		} else {
			// check and update data
			if productDescription := res.GetProductDescription(); len(productDescription) > 0 && productDescription != p.scraper.Product.ProductDescription {
				p.scraper.Product.ProductDescription = productDescription
			}

			if productName := res.GetProductName(); len(productName) > 0 && productName != p.scraper.Product.ProductName {
				p.scraper.Product.ProductName = productName
			}

			// if productId := res.GetProductId(); productId != p.scraper.Product.ProductID {
			// 	p.scraper.Product.ProductID = productId
			// }

			if quantityLimit := res.GetQuantityLimit(); quantityLimit > 0 && int(quantityLimit) != p.scraper.Product.QuantityLimit {
				p.scraper.Product.QuantityLimit = int(quantityLimit)
			}

			if sizeSkuMap := res.GetSizeSkuMap(); len(sizeSkuMap) > 0 {
				newMap := make(map[string]*models.SizeSkuMap)
				for k, v := range sizeSkuMap {
					if sickocommon.CheckSliceContains(p.scraper.DesireSizes, "RA") {
						newMap[k] = &models.SizeSkuMap{
							SkuId: v.SkuId,
							Gtin:  v.Gtin,
						}
					} else if sickocommon.CheckSliceContains(p.scraper.DesireSizes, k) {
						newMap[k] = &models.SizeSkuMap{
							SkuId: v.SkuId,
							Gtin:  v.Gtin,
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
		}
	}
}

func (p monitorPayload) doMonitorInformation(c chan *models.Error) {
	switch {
	case p.productInformationStream == nil || p.productInformationStream.Context().Err() != nil:
		tasks.SafeSend(c, &models.Error{Error: errors.New("GET_STREAM_ERROR"), Code: 409, Message: "ProductInformation get stream err"})
		return
	case p.scraper.Product == nil:
		tasks.SafeSend(c, &models.Error{Error: errors.New("INTERNAL_SERVER_ERROR"), Code: 500, Message: "Scraper Setting Error"})
		return
	}

	// Setup grpc connection and streams
	err := p.productInformationStream.Send(&grpc_service.StreamProductInformationRequest{
		Country:    p.taskGroupSetting.Country,
		MerchGroup: p.taskGroupSetting.MerchGroup,
		Language:   p.taskGroupSetting.Language,
		Ids: []string{
			p.scraper.Product.StyleColor,
		},
		Proxy: []string{
			sickocommon.GetProxy(p.monitorProxyGroup).String(), sickocommon.GetProxy(p.monitorProxyGroup).String(), sickocommon.GetProxy(p.monitorProxyGroup).String(), sickocommon.GetProxy(p.monitorProxyGroup).String(), sickocommon.GetProxy(p.monitorProxyGroup).String(), sickocommon.GetProxy(p.monitorProxyGroup).String(), sickocommon.GetProxy(p.monitorProxyGroup).String(),
		},
	})
	if err != nil && err != io.EOF {
		tasks.SafeSend(c, &models.Error{Error: errors.New("ERROR_GETTING_PRODUCT_INFORMATION"), Code: 409, Message: "Error Getting Product Information"})
	}
}

func (p monitorPayload) doCheckLive(c chan *models.Error) {
	switch {
	case p.productCheckliveStream == nil || p.productCheckliveStream.Context().Err() != nil:
		tasks.SafeSend(c, &models.Error{Error: errors.New("GET_STREAM_ERROR"), Code: 409, Message: "ProductInformation get stream err"})
		return
	case p.scraper.Product == nil, len(p.scraper.Product.ProductID) == 0:
		tasks.SafeSend(c, &models.Error{Error: errors.New("PRODUCT_NOT_FOUND"), Code: 409, Message: "Error Checking Product Status"})
		return
	}

	err := p.productCheckliveStream.Send(&grpc_service.StreamProductCheckLiveRequest{
		ProductId: p.scraper.Product.ProductID,
		Proxy:     sickocommon.GetProxy(p.monitorProxyGroup).String(),
	})
	if err != nil && err != io.EOF {
		tasks.SafeSend(c, &models.Error{Error: errors.New("ERROR_CHECKING_PRODUCT_LIVE"), Code: 409, Message: "Error Checking Product Live"})
	}
}

func (p *monitorPayload) doMonitor(c chan *models.Error) {
	switch {
	case p.productMonitorStream == nil || p.productMonitorStream.Context().Err() != nil:
		tasks.SafeSend(c, &models.Error{Error: errors.New("GET_STREAM_ERROR"), Code: 409, Message: "ProductInformation get stream err"})
		return
	case p.scraper.Product == nil, len(p.scraper.Product.ProductID) == 0:
		tasks.SafeSend(c, &models.Error{Error: errors.New("PRODUCT_NOT_FOUND"), Code: 409, Message: "Error Waiting For Restock"})
		return
	}

	err := p.productMonitorStream.Send(&grpc_service.StreamProductMonitorGraphqlRequest{
		ProductId: p.scraper.Product.ProductID,
		Country:   p.taskGroupSetting.Country,
		Locale:    p.taskGroupSetting.Locale,
		Proxy:     sickocommon.GetProxy(p.monitorProxyGroup).String(),
	})
	switch err {
	case nil:
		tasks.SafeSend(c, nil)
	default:
		tasks.SafeSend(c, &models.Error{Error: errors.New("ERROR_WAITING_FOR_RESTOCK"), Code: 409, Message: "Error Waiting For Restock"})
	}
}

func (p *monitorPayload) maintainGrpcConnection(ctx context.Context) {
	defer func() {
		if p.conn != nil {
			p.conn.Close()
		}
	}()
	t := time.NewTicker(20 * time.Millisecond)
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			switch {
			case p.conn == nil, p.conn.Value().GetState().String() != "READY":
				conn, err := grpc.GetFunctionServerGrpcCon()
				if err != nil {
					communicator.ModifyTaskStatus(&models.Message{GroupID: p.scraper.GroupID, TaskID: p.taskID, Status: "PENDING", Code: 500, Message: "Error Getting Stream Connection"})
				}
				p.conn = conn
				p.streamClient = grpc_service.NewStreamClient(p.conn.Value())
				if p.productInformationStream, err = p.streamClient.ProductInformation(ctx); err != nil {
					p.conn = nil
				}
				if p.productCheckliveStream, err = p.streamClient.ProductCheckLive(ctx); err != nil {
					p.conn = nil
				}
				if p.productMonitorStream, err = p.streamClient.ProductMonitor(ctx); err != nil {
					p.conn = nil
				}
			case p.productInformationStream.Context().Err() != nil:
				productInformationStream, err := p.streamClient.ProductInformation(ctx)
				if err != nil {
					p.conn = nil
					break
				}
				p.productInformationStream = productInformationStream

			case p.productCheckliveStream.Context().Err() != nil:
				productCheckliveStream, err := p.streamClient.ProductCheckLive(ctx)
				if err != nil {
					p.conn = nil
					break
				}
				p.productCheckliveStream = productCheckliveStream

			case p.productMonitorStream.Context().Err() != nil:
				productMonitorStream, err := p.streamClient.ProductMonitor(ctx)
				if err != nil {
					p.conn = nil
					break
				}
				p.productMonitorStream = productMonitorStream
			}
		}
	}
}

func (p *monitorPayload) setupMonitorGRPCTask() error {

	scraper := communicator.TaskScraperObjectGMap.Get(p.taskID)
	if scraper == nil {
		return errors.New("error starting scraper")
	}

	p.scraper = scraper.(*models.TaskScraper)
	if !p.scraper.Mutex.TryLock() {
		return errors.New("too many tasks assigned to one worker")
	}

	p.scraper.Product = &models.ScraperTaskProduct{
		StyleColor: p.scraper.Product.StyleColor,
	}

	// setup task group settings
	p.taskGroupSetting = *(communicator.Config.TaskGroups[p.scraper.GroupID].TaskGroupSetting)

	// setup monitor proxy group
	p.monitorProxyGroup = (communicator.Config.Proxies)[p.scraper.ScraperProxyGroupName]

	return nil
}

// HandleMonitorTaskWithOptions : HandleMonitorGraphqlTaskWithOptions
func HandleMonitorTaskWithOptions(ctx context.Context, task *asynq.Task) (err error) {

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

	if err := p.setupMonitorGRPCTask(); err != nil {
		return err
	}

	// setup & maintain grpc connection & streams
	go p.maintainGrpcConnection(ctx)

	// listener on stream receiving data
	ch := make(chan *models.Error, 3)
	defer close(ch)
	go p.doGetMonitorInformation(ch)
	go p.doGetCheckLive(ch)
	go p.doGetMonitor(ch)
	go func() {
		for {
			err, ok := <-ch
			if err != nil {
				communicator.ModifyTaskStatus(&models.Message{GroupID: p.scraper.GroupID, TaskID: p.taskID, Status: "PENDING", Code: err.Code, Message: err.Message})
			}
			if !ok {
				return
			}
		}
	}()

	// stream send data
	for i := 0; i < p.scraper.ScraperNum; i++ { // start scrapers according to scraper number set by user
		go func(ctx context.Context) {
			for {
				c := make(chan *models.Error, 3)

				go p.doMonitorInformation(c)
				go p.doCheckLive(c)
				go p.doMonitor(c)

			L:
				for {
					select {
					case <-ctx.Done():
						return
					case err := <-c: // controlling errors
						if err != nil {
							communicator.ModifyTaskStatus(&models.Message{GroupID: p.scraper.GroupID, TaskID: p.taskID, Status: "PENDING", Code: err.Code, Message: err.Message})
							break L
						}
						communicator.ModifyTaskStatus(&models.Message{GroupID: p.scraper.GroupID, TaskID: p.taskID, Status: "PENDING", Code: 200, Message: "SCRAPING..."})
						break L
					}
				}

				close(c)
				p.scraper.Product.Scrapes++
				communicator.TaskScraperObjectGMap.Set(p.taskID, p.scraper) // save the taskGroup
				<-time.NewTicker(time.Duration(p.scraper.ScraperDelay) * time.Millisecond).C
			}
		}(ctx)
	}

	// listener for instock product
	go func(ctx context.Context) {
		triggerDelay := time.NewTicker(time.Duration(p.scraper.TriggerDelay) * time.Millisecond)
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
						switch p.scraper.CheckLive { // check both product and stock live then activate normal task
						case true:
							if p.scraper.Product.IsLive {
								tasks.AddToQueue(p.scraper.GroupID, p.scraper.Product, p.scraper.TriggerNum)
							}
						default:
							tasks.AddToQueue(p.scraper.GroupID, p.scraper.Product, p.scraper.TriggerNum)
						}
						p.isCooling = true
						<-triggerDelay.C
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
