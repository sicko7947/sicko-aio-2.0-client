package handler

import (
	"net/http"

	"github.com/gogf/gf/net/ghttp"
	"sicko-aio-2.0-client/successHandler"
)

func GetSuccessItems(r *ghttp.Request) {

	items, err := successHandler.RetrieveSuccess()
	if err != nil {
		r.Response.WriteStatus(http.StatusBadRequest)
		r.Response.WriteExit()
		return
	}
	r.Response.WriteJsonExit(items)
}
