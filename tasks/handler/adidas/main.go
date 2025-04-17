package adidas

import (
	"context"
	"errors"
	"time"

	"github.com/hibiken/asynq"
	"github.com/sicko7947/sickocommon"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/tasks"
	"sicko-aio-2.0-client/utils/psychoclient"
)

func (p *taskPayload) doQueue() {

Loop:
	for {
		waitingRoomConfigChan, availabilityChan, queueChan := make(chan *models.Error, 1), make(chan *models.Error, 1), make(chan *models.Error, 1)

		go func() {
			if err := p.doGetWaitingRoomConfig(); err != nil { // get waiting room config
				tasks.SafeSend(waitingRoomConfigChan, err)
			}
		}()

		go func() {
			if err := p.doGetAvailability(); err != nil { // get product availability
				tasks.SafeSend(availabilityChan, err)
			}
		}()

		go func() {
			err := p.doGetQueue() // get current queue status
			tasks.SafeSend(queueChan, err)
		}()

	L:
		for {
			select {
			case err := <-waitingRoomConfigChan:
				communicator.ModifyTaskStatus(&models.Message{GroupID: p.worker.GroupID, TaskID: p.taskID, Status: "PENDING", Code: err.Code, Message: err.Message})
			case err := <-availabilityChan:
				communicator.ModifyTaskStatus(&models.Message{GroupID: p.worker.GroupID, TaskID: p.taskID, Status: "PENDING", Code: err.Code, Message: err.Message})
			case err := <-queueChan:
				if err == nil {
					break Loop
				}
				communicator.ModifyTaskStatus(&models.Message{GroupID: p.worker.GroupID, TaskID: p.taskID, Status: "PENDING", Code: err.Code, Message: err.Message})
				break L
			}
		}

		communicator.TaskWorkerObjectGMap.Set(p.taskID, p.worker)
		<-time.NewTicker(5000 * time.Millisecond).C
	}
}

func (p *taskPayload) setupTask() *models.Error {
	// set task group settings
	p.taskGroupSetting = communicator.Config.TaskGroups[p.worker.GroupID].TaskGroupSetting

	// setup task request session
	checkoutProxyGroup := (communicator.Config.Proxies)[p.worker.CheckoutProxyGroupName]
	if session, err := psychoclient.NewSession(&psychoclient.SessionBuilder{
		Proxy: sickocommon.GetProxy(checkoutProxyGroup).String(),
	}); err != nil {
		return err
	} else {
		p.session = session
	}

	// setup task headers
	p.headers = map[string]string{
		"accept":           "application/json, text/plain, */*",
		"accept-encoding":  "gzip, deflate, br",
		"accept-language":  "en-US,en;q=0.9",
		"dnt":              "1",
		"sec-ch-ua-mobile": "?0",
		"sec-fetch-dest":   "empty",
		"sec-fetch-mode":   "cors",
		"sec-fetch-site":   "cross-site",
		"user-agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36",
		"cache-control":    "no-cache",
	}

	return nil
}

// HandleAdidasTaskWithOptions : HandleAdidasTaskWithOptions
func HandleAdidasTaskWithOptions(ctx context.Context, task *asynq.Task) error {
	// Get basic task info

	p := &taskPayload{}
	p.taskID = models.TaskID(task.Payload())
	if worker := communicator.TaskWorkerObjectGMap.Get(p.taskID); worker != nil {
		p.worker = worker.(*models.TaskWorker)
		if !p.worker.Mutex.TryLock() {
			return errors.New("too many tasks assigned to one worker")
		}
		defer func() {
			if p.worker.Mutex.IsLocked() {
				p.worker.Mutex.Unlock()
			}
		}()
	} else {
		return errors.New("error starting worker")
	}

	// further setup task
	if err := p.setupTask(); err != nil {
		communicator.ModifyTaskStatus(&models.Message{GroupID: p.worker.GroupID, TaskID: p.taskID, Status: "PENDING", Code: err.Code, Message: err.Message})
		return err.Error
	}

	c := make(chan *models.Error, 1)
	done := make(chan bool, 1)
	go func() {
	loop:
		for {
			p.doQueue()

			go func() {
				if err := p.doGetShippingMethod(); err != nil {
					tasks.SafeSend(c, err)
				}
			}()

			if err := p.doATC(); err != nil {
				tasks.SafeSend(c, err)
				break
			}

			if err := p.doSubmitAddress(); err != nil {
				tasks.SafeSend(c, err)
				break
			}

			err := p.doSubmitOrder() // need cookie
			if err != nil {
				tasks.SafeSend(c, err)
				break
			} else {
				done <- true
				break loop
			}
		}
	}()

	for {
		select {
		case err := <-c:
			communicator.ModifyTaskStatus(&models.Message{ // send error message to frontend
				GroupID: p.worker.GroupID,
				TaskID:  p.taskID,
				Status:  "PENDING",
				Code:    err.Code,
				Message: err.Message,
			})
		case <-ctx.Done():
			communicator.ModifyTaskStatus(&models.Message{ // send task cancellation status to frontend
				Code:    0,
				GroupID: p.worker.GroupID,
				TaskID:  p.taskID,
				Status:  "CANCELLED",
				Message: "Stopped",
			})
			return ctx.Err() // cancelation signal received, abandon this work.
		case <-done:
			communicator.ModifyTaskStatus(&models.Message{ // send success checkout message to frontend
				Code:    200,
				GroupID: p.worker.GroupID,
				TaskID:  p.taskID,
				Status:  "COMPLETED",
				Message: "Checked Out!",
			})

			// send success checkout message to webhooks
			// webhook.SendDiscordLegacyWebhook(p.taskGroupSetting, p.worker)
			// webhook.SendSlackLegacyWebhook(p.taskGroupSetting, p.worker)
			return nil
		}
	}
}
