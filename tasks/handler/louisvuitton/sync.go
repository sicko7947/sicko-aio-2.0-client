package louisvuitton

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"

	"github.com/tidwall/gjson"
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/utils/psychoclient"
	"sicko-aio-2.0-client/utils/redis"
)

func Sync(src *models.Account, proxy string) (dst *models.Account, err *models.Error) {
	useragent, cookieMap, err := redis.GetCookie2FromRedis("akamai.mrporter.com")
	if err != nil {
		return nil, err
	}

	var sesh psychoclient.Session
	sesh, err = psychoclient.NewSession(&psychoclient.SessionBuilder{
		Proxy: proxy,
	})
	if err != nil {
		return nil, err
	}

	sesh.SetCookies(cookieMap)

	endpoint := `xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx`
	data, _ := json.Marshal(map[string]string{
		"logonId":       src.Email,
		"logonPassword": src.Password,
	})
	reqId, _ := sesh.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "POST",
		Headers: map[string]string{
			"accept":              "*/*",
			"accept-language":     "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7,es;q=0.6",
			"application-name":    "myaccount",
			"application-version": "5.510.0",
			"content-type":        "application/json",
			"dnt":                 "1",
			"label":               "login",
			"sec-ch-ua-mobile":    "?0",
			"sec-fetch-dest":      "empty",
			"sec-fetch-mode":      "cors",
			"sec-fetch-site":      "same-origin",
			"x-ibm-client-id":     "0b1e2c22-581d-435b-9cde-70bc52cba701",
			"cache-control":       "no-cache",
			"user-agent":          useragent,
		},
		Payload: bytes.NewBuffer(data),
	})

	res, respBody, err := sesh.Do(reqId)
	if err != nil {
		return nil, err
	}
	result := gjson.Parse(string(respBody))
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
