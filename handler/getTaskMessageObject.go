package handler

import (
	"github.com/gogf/gf/net/ghttp"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/models"
)

func GetTaskMessageObject(r *ghttp.Request) {
	groupID := models.TaskGroupID(r.Get("groupId").(string))

	obj := make(map[models.TaskID]map[string]interface{}) // init obj

	communicator.TaskMessageGMap.Iterator(func(k, v interface{}) bool {
		value := v.(map[string]interface{})
		if value["groupId"].(models.TaskGroupID) == groupID {
			obj[k.(models.TaskID)] = value
		}
		return true
	})

	r.Response.WriteJsonExit(obj)
}
