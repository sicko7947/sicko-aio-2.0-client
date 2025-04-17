package handler

import (
	"net/http"

	"github.com/gogf/gf/net/ghttp"
	"sicko-aio-2.0-client/captchaPool"
	"sicko-aio-2.0-client/models"
)

type captcha struct {
	Token string
	Site  models.CATEGORY
}

// CancelTasks : Cancel Selected Tasks
func StoreCaptchaToken(r *ghttp.Request) {

	// Parsing Payload
	var p *captcha
	if err := r.Parse(&p); err != nil {
		r.Response.WriteStatus(http.StatusBadRequest)
		r.Response.WriteJsonExit(map[string]bool{"success": false})
		return
	}

	switch p.Site {
	case models.ADIDAS:
		captchaPool.PutAdidasCaptchaToken(p.Token)
	case models.SUPPLYSTORE:
		captchaPool.PutSupplyStoreCaptchaToken(p.Token)
	}
	r.Response.WriteJsonExit(map[string]bool{"success": true})
}
