package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/net/ghttp"
	"sicko-aio-2.0-client/constants"
	"sicko-aio-2.0-client/utils/psychoclient"
)

type testProxyPayload struct {
	Store   string   `json:"store"`
	Url     string   `json:"url,omitempty"`
	Proxies []string `json:"proxies"`
}

type result struct {
	ip     string
	speed  int64
	status int
}

var proxyResult *gmap.StrAnyMap

func init() {
	proxyResult = gmap.NewStrAnyMap(true)
}

// TestProxies : TestProxies
func TestProxies(r *ghttp.Request) {
	// Parsing Payload
	var payload *testProxyPayload
	if err := r.Parse(&payload); err != nil {
		r.Response.WriteStatus(http.StatusBadRequest)
		r.Response.WriteJsonExit(map[string]bool{"success": false})
		return
	}

	go func() {
		c := make(chan *result, len(payload.Proxies))
		defer close(c)
		count := 0

		for _, proxy := range payload.Proxies {
			go func() {
				defer recover()

				start := time.Now()
				res, _, err := psychoclient.NewClient(&psychoclient.SessionBuilder{
					Proxy: proxy,
				}).RoundTripNewRequest(&psychoclient.RequestBuilder{
					Endpoint: payload.Url,
					Method:   "GET",
					Headers: map[string]string{
						"accept":          "application/json",
						"content-type":    "application/json; charset=UTF-8",
						"accept-language": "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7,es;q=0.6",
						"cache-control":   "no-cache",
						"user-agent":      gofakeit.RandomString(constants.ChromeUAList),
					},
					Payload: nil,
				})
				if err != nil {
					c <- &result{
						ip:     proxy,
						speed:  0,
						status: 400,
					}
					return
				}

				duration := time.Since(start)

				c <- &result{
					ip:     proxy,
					speed:  duration.Milliseconds(),
					status: res.StatusCode,
				}
			}()
			count++
		}

		for {
			select {
			case res := <-c:
				fmt.Println(res)
				proxyResult.Set(res.ip, map[string]interface{}{
					"speed":  fmt.Sprint(res.speed) + "ms",
					"status": res.status,
				})

			}
			if count == len(payload.Proxies) {
				break
			}
		}
	}()

	r.Response.WriteJsonExit(map[string]bool{"success": true})
}

func GetProxiesTestResult(r *ghttp.Request) {
	result := []map[string]interface{}{}

	proxyResult.Iterator(func(k string, v interface{}) bool {
		vv := v.(map[string]interface{})
		result = append(result, map[string]interface{}{
			"proxy":  k,
			"status": vv["status"],
			"speed":  vv["speed"],
		})
		return true
	})
	r.Response.WriteJsonExit(result)
}
