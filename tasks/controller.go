package tasks

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynq/tools/asynq/cmd"
	"github.com/huandu/go-clone"
	"github.com/tidwall/gjson"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/tasks/defineder"
)

const (
	redisURL = "localhost:6379" // constant redis addr
)

// global variables
var (
	r           = cmd.CreateRDB()
	Client      *asynq.Client
	redisClient asynq.RedisClientOpt
)

func init() {
	redisClient = asynq.RedisClientOpt{Addr: redisURL}
	Client = asynq.NewClient(redisClient) // initiate client connect to redis for distributing tasks
}

// CancelTask : send cancellation signal to rdb
func CancelTask(id uuid.UUID) error {
	err := r.PublishCancelation(id.String())
	if err != nil {
		return err
	}
	return nil
}

// SafeSend : check if channel is closed and safely send data to channel
func SafeSend(ch chan *models.Error, value *models.Error) (closed bool) {
	defer func() {
		if recover() != nil {
			closed = true
		}
	}()
	ch <- value  // panic if ch is closed
	return false // <=> closed = false; return
}

// AddToQueue : Add tasks to queue from functions
func AddToQueue(groupID models.TaskGroupID, p *models.ScraperTaskProduct, triggerNum int) {
	product := clone.Clone(p).(*models.ScraperTaskProduct) // deeop copy obj incase unwanted memory sharing error

	count := 0

	done := make(chan bool, 1)
	msg := make(chan *models.Message, triggerNum)
	defer close(done)
	defer close(msg)
	go func() {
		defer func() {
			done <- true
		}()

		data, _ := communicator.TaskMessageGMap.MarshalJSON()
		taskMessageMap := gjson.ParseBytes(data).Map()

		for k, v := range taskMessageMap {
			if v.Get("status").String() == "ARCHIVED" && models.TaskGroupID(v.Get("groupId").String()) == groupID {

				if count > triggerNum { // check if exceeds trigger number
					break
				}
				count++

				msg <- &models.Message{ // Change task status to pending incase for re-enque
					Code:    200,
					GroupID: groupID,
					TaskID:  models.TaskID(k),
					Status:  "PENDING",
					Message: "Received Job",
				}

				if taskworkerObj := communicator.TaskWorkerObjectGMap.Get(models.TaskID(k)); taskworkerObj != nil {
					taskworker := taskworkerObj.(*models.TaskWorker)
					if taskworker.Mutex.IsLocked() {
						count--
						continue
					}

					// initialize Product for empty taskworker product
					if taskworker.Product == nil {
						taskworker.Product = &models.WorkerTaskProduct{}
					}
					// initialize TaskInfo for empty taskworker info
					if taskworker.TaskInfo == nil {
						taskworker.TaskInfo = &models.WorkerTaskInfo{}
					}

					taskworker.Product.StyleColor = product.StyleColor
					taskworker.Product.QuantityLimit = product.QuantityLimit
					taskworker.Product.ImageURL = product.ImageURL
					taskworker.Product.ProductName = product.ProductName
					taskworker.Product.ProductDescription = product.ProductDescription
					taskworker.Product.ProductID = product.ProductID
					taskworker.Product.Price = product.Price
					taskworker.Product.Other = product.Other
					for size, sizeSkuMap := range product.SizeSkuMap { // set size and skuid to worker's product
						taskworker.Product.Size = size
						taskworker.Product.SkuID = sizeSkuMap.SkuId
						taskworker.Product.Gtin = sizeSkuMap.Gtin
						break
					}
					communicator.TaskWorkerObjectGMap.Set(taskworker.TaskID, taskworker)

					var taskInvoke *asynq.Task
					switch communicator.Config.TaskGroups[groupID].TaskGroupSetting.TaskType { // Define task types

					// NIKE
					case defineder.NikeCheckoutLegacy:
						taskInvoke = defineder.NewNikeCheckoutLegacyTask(taskworker.TaskID)
					case defineder.NikeCheckoutLegacyV2:
						taskInvoke = defineder.NewNikeCheckoutLegacyV2Task(taskworker.TaskID)
					case defineder.NikeCheckoutLegacyV3:
						taskInvoke = defineder.NewNikeCheckoutLegacyV3Task(taskworker.TaskID)
					case defineder.NikeCheckoutReserveStock:
						taskInvoke = defineder.NewNikeCheckoutReserveStockTask(taskworker.TaskID)
					case defineder.NikeCheckoutV2:
						taskInvoke = defineder.NewNikeCheckoutV2Task(taskworker.TaskID)
					case defineder.NikeCheckoutV3:
						taskInvoke = defineder.NewNikeCheckoutV3Task(taskworker.TaskID)

					// SNKRS
					case defineder.SnkrsLaunchEntry:
						taskInvoke = defineder.NewSnkrsLaunchEntry(taskworker.TaskID)

					// MRPORTER
					case defineder.MrPorterCheckout:
						taskInvoke = defineder.NewMrPorterCheckoutTask(taskworker.TaskID)

					// NETAPORTER
					case defineder.NetAPorterCheckout:
						taskInvoke = defineder.NewNetAPorterCheckoutTask(taskworker.TaskID)

					// LUISAVIAROMA
					case defineder.LuisaviaromaCheckout:
						taskInvoke = defineder.NewLuisaviaromaCheckoutTask(taskworker.TaskID)

					// SUPPLYSTORE
					case defineder.SupplyStoreCheckout:
						taskInvoke = defineder.NewSupplyStoreCheckoutTask(taskworker.TaskID)

					// MACYS
					case defineder.MacysCheckout:
						taskInvoke = defineder.NewMacysCheckoutTask(taskworker.TaskID)

					// SNEAKERBOY
					case defineder.SneakerboyCheckout:
						taskInvoke = defineder.NewSneakerboyCheckoutTask(taskworker.TaskID)

					// PACSUN
					case defineder.PacsunCheckout:
						taskInvoke = defineder.NewPacsunCheckoutTask(taskworker.TaskID)

					// SSENSE
					case defineder.SsenseCheckout:
						taskInvoke = defineder.NewSsenseCheckoutTask(taskworker.TaskID)

					// NEW BALANCE
					case defineder.NewBalanaceCheckout:
						taskInvoke = defineder.NewNewBalanaceCheckoutTask(taskworker.TaskID)

					// LANECRAWFORD
					case defineder.LanecrawfordCheckout:
						taskInvoke = defineder.NewLanecrawfordCheckoutTask(taskworker.TaskID)

					// TAF
					case defineder.TafCheckout:
						taskInvoke = defineder.NewTafCheckoutTask(taskworker.TaskID)

					// TI
					case defineder.TexasInstrumentsCheckout:
						taskInvoke = defineder.NewTexasInstrumentsCheckoutTask(taskworker.TaskID)
					}

					// parse the start time for worker
					startTime, err := time.Parse("2006-01-02T15:04:05.000Z", taskworker.StartTime)
					wait := err == nil && startTime.UTC().After(time.Now())
					if wait {
						msg <- &models.Message{ // update the task message with starting time
							Code:    200,
							GroupID: groupID,
							TaskID:  models.TaskID(k),
							Status:  "PENDING",
							Message: fmt.Sprintf(`Starting at %s`, taskworker.StartTime),
						}
					}
					go func() {
						if wait {
							<-time.After(time.Until(startTime.UTC()))
						}
						Client.Enqueue( // enqueue task
							taskInvoke,
							asynq.MaxRetry(-1),
							asynq.Deadline(time.Now().AddDate(3, 0, 0)),
						)
					}()
				}
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		case m := <-msg:
			communicator.ModifyTaskStatus(m)
		}
	}
}
