package handler

import (
	"encoding/json"
	"sync"

	"github.com/gogf/gf/net/ghttp"
	"github.com/huandu/go-clone"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/models"
)

var ws *ghttp.WebSocket
var outChan = make(chan interface{}, 1000)

const (
	INIT_TASK_MESSAGE_OBJECTS = "INIT_TASK_MESSAGE_OBJECTS" // get all current task message objects
	GET_TASK_MESSAGE_OBJECTS  = "GET_TASK_MESSAGE_OBJECTS"  // get latest task message
	GET_TASK_OBJECTS          = "GET_TASK_OBJECTS"          // get all current task objects
	GET_TASK_GROUP_OBJECTS    = "GET_TASK_GROUP_OBJECTS"    // get all current task group objects
	GET_TASK_LOGS             = "GET_TASK_LOGS"             // get current task log history
	INIT_CONFIG               = "INIT_CONFIG"               // get current config
	UPDATE_CONFIG             = "UPDATE_CONFIG"             // update the config
	SYNC_ACCOUNTS             = "SYNC_ACCOUNTS"             // sync accounts
)

type websocketPayload struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Data string `json:"data,omitempty"`
}

func InitialWebsocketConnection(r *ghttp.Request) {
	ws, _ = r.WebSocket()
	ws.EnableWriteCompression(true)

	go func() {
		for {
			var payload websocketPayload
			err := ws.ReadJSON(&payload)
			if err != nil {
				return
			}

			switch payload.Type {

			case INIT_CONFIG:
				go sendConfig()
			case UPDATE_CONFIG:
				go updateConfig(payload.Data)

			case GET_TASK_GROUP_OBJECTS:
				go sendTaskGroupObject()
			case GET_TASK_OBJECTS:
				go sendTaskObject(payload.ID)
			case INIT_TASK_MESSAGE_OBJECTS:
				go sendTaskMsgObject(payload.ID)
			case GET_TASK_MESSAGE_OBJECTS:
				go sendLatestTaskMsgObject(payload.ID)
			case GET_TASK_LOGS:
				go sendTaskLogs(payload.ID)

			case SYNC_ACCOUNTS:
			}
		}
	}()

	go writeLoop(outChan)

}

func writeLoop(outChan chan interface{}) {
	defer func() {
		if recover() != nil {
			writeLoop(outChan)
		}
	}()

	for {
		obj, ok := <-outChan
		if !ok {
			ws.Close()
			return
		}
		ws.WriteJSON(obj)
	}
}

func sendTaskGroupObject() {
	obj := make(map[string]map[string]interface{})
	for k, v := range communicator.Config.TaskGroups {
		obj[string(k)] = map[string]interface{}{
			"scraperCount":     len(v.TaskScrapers),
			"workerCount":      len(v.TaskWorkers),
			"successCount":     communicator.SuccessCountGmap.Get(string(k)),
			"taskGroupSetting": v.TaskGroupSetting,
		}
	}

	data, _ := json.Marshal(obj)
	outChan <- &websocketPayload{
		Type: "GET_TASK_GROUP_OBJECTS",
		Data: string(data),
	}
}

func sendTaskObject(groupId string) {
	groupID := models.TaskGroupID(groupId)

	type taskGroup struct {
		GroupID          models.TaskGroupID      `json:"groupId"`
		TaskGroupSetting models.TaskGroupSetting `json:"taskGroupSetting"`
		TaskWorkers      map[string]interface{}  `json:"taskWorkers"`
		TaskScrapers     map[string]interface{}  `json:"taskScrapers"`
	}

	var wg sync.WaitGroup
	taskScrapers := make(map[string]interface{})
	taskWorkers := make(map[string]interface{})

	wg.Add(2)
	go func(done func()) {
		defer done()
		for taskId, v := range communicator.TaskScraperObjectGMap.MapStrAny() {
			scraper := v.(*models.TaskScraper)
			if scraper.GroupID == groupID {
				taskScrapers[taskId] = scraper
			}
		}
	}(wg.Done)
	go func(done func()) {
		defer done()
		for taskId, v := range communicator.TaskWorkerObjectGMap.MapStrAny() {
			worker := v.(*models.TaskWorker)
			if worker.GroupID == groupID {
				taskWorkers[taskId] = worker
			}
		}
	}(wg.Done)
	wg.Wait()

	if communicator.Config.TaskGroups[groupID] == nil {
		return
	}

	data, _ := json.Marshal(taskGroup{ // init obj
		GroupID:          groupID,
		TaskScrapers:     taskScrapers,
		TaskWorkers:      taskWorkers,
		TaskGroupSetting: *(communicator.Config.TaskGroups[groupID].TaskGroupSetting),
	})

	outChan <- &websocketPayload{
		Type: "GET_TASK_OBJECTS",
		ID:   groupId,
		Data: string(data),
	}
}

func sendTaskMsgObject(groupId string) {
	groupID := models.TaskGroupID(groupId)

	data := make(map[string]map[string]interface{}) // init obj

	for taskId, v := range communicator.TaskMessageGMap.MapStrAny() {
		value := v.(map[string]interface{})
		if value["groupId"].(models.TaskGroupID) == groupID {
			data[taskId] = value
		}
	}

	payload, _ := json.Marshal(data)
	outChan <- &websocketPayload{
		Type: "INIT_TASK_MESSAGE_OBJECTS",
		ID:   groupId,
		Data: string(payload),
	}
}

func sendLatestTaskMsgObject(groupId string) {
	groupID := models.TaskGroupID(groupId)

	data := make(map[string]map[string]interface{})
	for taskId, v := range communicator.TempMessageGmap.MapStrAny() {
		value := v.(map[string]interface{})
		if value["groupId"].(models.TaskGroupID) == groupID {
			data[taskId] = value
		}
	}

	payload, _ := json.Marshal(data)
	outChan <- &websocketPayload{
		Type: "GET_TASK_MESSAGE_OBJECTS",
		ID:   groupId,
		Data: string(payload),
	}
	communicator.TempMessageGmap.Clear()
}

func sendTaskLogs(taskId string) {
	taskID := models.TaskID(taskId)

	if logs := communicator.TaskLogsGMap.Get(taskID); logs != nil {
		payload, _ := json.Marshal(logs)
		outChan <- &websocketPayload{
			Type: "GET_TASK_LOGS",
			ID:   taskId,
			Data: string(payload),
		}
	}
}

func sendConfig() {
	config := clone.Clone(communicator.Config).(*models.Config) // deep copy
	payload, _ := json.Marshal(config)
	outChan <- &websocketPayload{
		Type: "UPDATE_CONFIG",
		Data: string(payload),
	}
}

func updateConfig(data string) {
	var config *models.Config
	err := json.Unmarshal([]byte(data), &config)
	if err != nil {
		return
	}

	// overide config with new config
	if config.TaskGroups != nil {
		communicator.Config.TaskGroups = config.TaskGroups
	}
	if config.Accounts != nil {
		communicator.Config.Accounts = config.Accounts
	}

	if config.Profiles != nil {
		communicator.Config.Profiles = config.Profiles
	}

	if config.Proxies != nil {
		communicator.Config.Proxies = config.Proxies
	}

	if config.GiftCards != nil {
		communicator.Config.GiftCards = config.GiftCards
	}

	if config.Discounts != nil {
		communicator.Config.Discounts = config.Discounts
	}

	if config.Settings != nil {
		communicator.Config.Settings = config.Settings
	}

	sendConfig()
}
