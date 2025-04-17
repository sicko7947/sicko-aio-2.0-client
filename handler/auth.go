package handler

import (
	"github.com/gogf/gf/net/ghttp"
	"sicko-aio-2.0-client/auth"
	"sicko-aio-2.0-client/communicator"
)

func Activate(r *ghttp.Request) {
	key := r.Get("key").(string)

	communicator.Config.Settings.Key = key
	code, msg := auth.Login()
	r.Response.WriteJsonExit(map[string]interface{}{
		"code":    code,
		"message": msg,
	})
}

func Deactivate(r *ghttp.Request) {
	auth.Deactivate()
	r.Response.WriteExit()
}
