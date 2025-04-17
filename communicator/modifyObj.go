package communicator

import (
	"time"

	"sicko-aio-2.0-client/models"
)

// ModifyTaskStatus : Modify task status by task id & group id
func ModifyTaskStatus(m *models.Message) {

	msg := map[string]interface{}{
		"groupId": m.GroupID, // set current task group Id
		"status":  m.Status,  // set current task status
		"code":    m.Code,    // set current task status code
		"message": m.Message, // set current task message
	}

	// Update temp task message
	TempMessageGmap.Set(m.TaskID, msg)

	// Update current task message
	TaskMessageGMap.Set(m.TaskID, msg)

	// Update logs to Logs gmap
	if taskLogsObj := TaskLogsGMap.GetOrSet(m.TaskID, []*models.Message{
		{
			TimeStamp: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
			Code:      m.Code,
			Status:    m.Status,
			Message:   m.Message,
		},
	}); taskLogsObj != nil {
		newTaskLogs := taskLogsObj.([]*models.Message)
		if len(newTaskLogs) > 200 { // check if task logs reached 200, if so only take the latest 200, otherwise continue
			newTaskLogs = newTaskLogs[200:]
		}
		newTaskLogs = append(newTaskLogs, &models.Message{ // appending new status into slice
			TimeStamp: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"), // setup timestamp
			Status:    m.Status,                                            // add status
			Code:      m.Code,                                              // add status code
			Message:   m.Message,                                           // add task message
		})
		TaskLogsGMap.Set(m.TaskID, newTaskLogs) // save new task logs to map
	}
}
