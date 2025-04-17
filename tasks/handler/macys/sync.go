package macys

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/sicko7947/sickocommon"
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

	// setup session cookie
	useragent, cookieMap, err := redis.GetCookie2FromRedis("akamai.macys.com")
	if err != nil {
		return nil, err
	}
	sesh.SetCookies(cookieMap)

	endpoint := `https://xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx`
	form := url.Values{}
	data := map[string]string{}
	for key, value := range data {
		form.Set(key, value)
	}

	reqId, _ := sesh.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "POST",
		Headers: map[string]string{
			"accept":             "application/json, text/javascript, */*; q=0.01",
			"accept-language":    "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7,es;q=0.6",
			"content-type":       "application/json",
			"dnt":                "1",
			"origin":             "https://www.macys.com",
			"user-agent":         useragent,
			"x-macys-request-id": sickocommon.NikeUUID(),
			"x-requested-with":   "XMLHttpRequest",
			"cache-control":      "no-cache",
		},
		Payload: strings.NewReader(form.Encode()),
	})

	res, respBody, err := sesh.Do(reqId, false)
	if err != nil {
		return nil, err
	}
	switch res.StatusCode {
	case 200, 201, 202:
		result := gjson.ParseBytes(respBody)
		if accessTokenObj := result.Get("access_token"); accessTokenObj.Exists() {
			accessToken := accessTokenObj.String()
			refreshToken := result.Get("refresh_token").String()
			lastSyncTime := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

			return &models.Account{
				Email:        src.Email,
				Password:     src.Password,
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
				Status:       "synced",
				LastSyncTime: lastSyncTime,
			}, nil
		}
		fallthrough
	case 401:
		return &models.Account{
			Email:    src.Email,
			Password: src.Password,
		}, nil
	default:
		return nil, &models.Error{Error: errors.New("ERROR_SYNCING_ACCOUNT"), Code: res.StatusCode, Message: "Error Syncing Account"}
	}
}
