package models

import (
	"net/http"

	"github.com/sicko7947/sickocommon"
)

type BrowserCheckoutPayload struct {
	CookieMap map[string]*http.Cookie
	UserAgent string
	Proxy     *sickocommon.Proxy
}
