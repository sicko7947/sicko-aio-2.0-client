package models

import "net/http"

type Cookie struct {
	Useragent string                  `json:"useragent"`
	CookieMap map[string]*http.Cookie `json:"cookieMap"`
}
