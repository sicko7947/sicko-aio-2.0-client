package mrporter

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"

	"github.com/tidwall/gjson"
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/utils"
	"sicko-aio-2.0-client/utils/psychoclient"
	"sicko-aio-2.0-client/utils/redis"
)

func Sync(src *models.Account, proxies ...string) (dst *models.Account, err *models.Error) {
	proxy := utils.GetSyncProxy(proxies)

	var sesh psychoclient.Session
	sesh, err = psychoclient.NewSession(&psychoclient.SessionBuilder{
		Proxy: proxy,
	})
	if err != nil {
		return nil, err
	}

	useragent, cookieMap, err := redis.GetCookie2FromRedis("akamai.mrporter.com")
	if err != nil {
		return nil, err
	}
	sesh.SetCookies(cookieMap)

	endpoint := `https://xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx`
	data, _ := json.Marshal(map[string]string{
		"logonId":       src.Email,
		"logonPassword": src.Password,
	})
	reqId, _ := sesh.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "POST",
		Headers: map[string]string{
			"accept":              "*/*",
			"accept-encoding":     "gzip, deflate, br",
			"accept-language":     "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7,es;q=0.6",
			"application-name":    "myaccount",
			"application-version": "5.573.0",
			"content-type":        "application/json",
			"dnt":                 "1",
			"label":               "login",
			"origin":              "https://www.mrporter.com",
			"user-agent":          useragent,
		},
		Payload: bytes.NewBuffer(data),
	})

	res, respBody, err := sesh.Do(reqId)
	if err != nil {
		return nil, err
	}
	result := gjson.ParseBytes(respBody)
	switch res.StatusCode {
	case 200, 201, 202:
		accessToken := result.Get("Ubertoken").String()
		lastSyncTime := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

		return &models.Account{
			Email:        src.Email,
			Password:     src.Password,
			AccessToken:  accessToken,
			RefreshToken: ``,
			Status:       "synced",
			LastSyncTime: lastSyncTime,
		}, nil
	default:
		return nil, &models.Error{Error: errors.New("ERROR_LOGIN_IN"), Code: res.StatusCode, Message: "Error login in"}
	}
}
