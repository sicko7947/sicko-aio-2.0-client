package server

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"sicko-aio-2.0-client/handler"
)

// BackendServer : server thats runs at the backend operating the bot
func BackendServer() {
	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(MiddlewareCORS)

		group.POST("/login", handler.Activate)
		group.POST("/deactivate", handler.Deactivate)

		group.PATCH("/tasklogs", handler.GetTaskLogs)
		group.PATCH("/taskobj", handler.GetTaskObject)
		group.PATCH("/taskmsgobj", handler.GetTaskMessageObject)

		group.GET("/config", handler.GetConfig)
		group.PATCH("/config", handler.ModifyConfig)

		group.POST("/account/sync", handler.SyncAccount)

		group.GET("/success", handler.GetSuccessItems)

		group.PUT("/captcha", handler.StoreCaptchaToken)

		group.PATCH("/taskgroup/modify", handler.ModifyTaskGroup)
		group.PATCH("/taskgroup/delete", handler.DeleteTaskGroup)
		group.PATCH("/taskgroup/scraper/modify", handler.ModifyTaskGroupScraper)
		group.PATCH("/taskgroup/worker/modify", handler.ModifyTaskGroupWorker)
		group.PATCH("/task/cancel", handler.CancelTasks)
		group.PATCH("/task/delete", handler.DeleteTasks)

		group.PATCH("/taskgroup/scraper/start", handler.StartScraperTasks)
		group.PATCH("/taskgroup/worker/start", handler.StartWorkerTasks)
	})
	s.BindHandler("/ws", handler.InitialWebsocketConnection)
	s.SetLogStdout(false)
	s.SetPort(26667)
	s.Run()
}

func MiddlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}
