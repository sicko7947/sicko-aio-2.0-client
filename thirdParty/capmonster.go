package thirdParty

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gogf/gf/container/gmap"
	"github.com/sicko7947/sickocommon"
	"github.com/tidwall/gjson"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/utils/psychoclient"
)

const (
	NoCaptchaTaskProxyless   string = "NoCaptchaTaskProxyless"
	RecaptchaV3TaskProxyless string = "RecaptchaV3TaskProxyless"
)

type CaptchaTask struct {
	TaskType            string `json:"type"`
	WebsiteURL          string `json:"websiteURL"`
	WebsiteKey          string `json:"websiteKey"`
	RecaptchaDataSValue string `json:"recaptchaDataSValue,omitempty"`
	ProxyType           string `json:"proxyType"`
	ProxyAddress        string `json:"proxyAddress"`
	ProxyPort           int64  `json:"proxyPort"`
	ProxyLogin          string `json:"proxyLogin,omitempty"`
	ProxyPassword       string `json:"proxyPassword,omitempty"`
	UserAgent           string `json:"userAgent,omitempty"`
	Cookies             string `json:"cookies,omitempty"`
}

type Capmonster interface {
	CreateTask(task *CaptchaTask) (taskId int, err error)
	GetTaskResult(taskId int) (capResponse string, err error)
	GetBalance()
}

type capmonster struct {
	clientKey string
	tasks     *gmap.IntStrMap
	session   psychoclient.Session
}

func New(clientKey string) Capmonster {
	checkoutProxyGroup := (communicator.Config.Proxies)["sicko"]
	sesh, _ := psychoclient.NewSession(&psychoclient.SessionBuilder{
		Proxy: sickocommon.GetProxy(checkoutProxyGroup).String(),
	})
	return &capmonster{
		session:   sesh,
		clientKey: clientKey,
		tasks:     gmap.NewIntStrMap(true),
	}
}

func (c *capmonster) CreateTask(t *CaptchaTask) (taskId int, err error) {
	var data []byte
	switch t.TaskType {
	case NoCaptchaTaskProxyless:
		data, _ = json.Marshal(map[string]interface{}{
			"clientKey": c.clientKey,
			"task": map[string]string{
				"type":       NoCaptchaTaskProxyless,
				"websiteURL": t.WebsiteURL,
				"websiteKey": t.WebsiteKey,
			},
		})
	case RecaptchaV3TaskProxyless:
		data, _ = json.Marshal(map[string]interface{}{
			"clientKey": c.clientKey,
			"task": map[string]interface{}{
				"type":       NoCaptchaTaskProxyless,
				"websiteURL": t.WebsiteURL,
				"websiteKey": t.WebsiteKey,
				"minScore":   0.3,
				"pageAction": "login",
			},
		})

	}

	reqId, _ := c.session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: "https://api.capmonster.cloud/createTask",
		Method:   "POST",
		Payload:  bytes.NewBuffer(data),
	})
	res, respBody, e := c.session.Do(reqId)
	if e != nil {
		return 0, e.Error
	}
	result := gjson.Parse(string(respBody))
	switch res.StatusCode {
	case 200, 201, 202:

		if errorId := result.Get("errorId").Int(); errorId == 0 {
			taskId := int(result.Get("taskId").Int())
			c.tasks.Set(taskId, "processing")
			return taskId, nil
		}
		fallthrough
	default:
		return 0, fmt.Errorf(`error creating %v type task`, t.TaskType)
	}
}

func (c *capmonster) GetTaskResult(taskId int) (capResponse string, err error) {
	data, _ := json.Marshal(map[string]interface{}{
		"clientKey": c.clientKey,
		"taskId":    taskId,
	})

	reqId, _ := c.session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: "https://api.capmonster.cloud/getTaskResult",
		Method:   "POST",
		Payload:  bytes.NewBuffer(data),
	})

	resCh := make(chan string, 1)
	errCh := make(chan error, 1)
	defer close(resCh)
	defer close(errCh)

	go func(resCh chan<- string, errCh chan<- error) {
		defer func() {
			recover()
		}()

		for {
			res, respBody, err := c.session.Do(reqId, false)
			if err != nil {
				errCh <- err.Error
				return
			}
			result := gjson.Parse(string(respBody))
			switch res.StatusCode {
			case 200, 201, 202:
				if res := result.Get("solution.gRecaptchaResponse").String(); len(res) > 0 {
					resCh <- res
					return
				}
				fallthrough
			default:
				errCh <- fmt.Errorf(`error getting captcha response - status code %v`, res.StatusCode)
				return
			}
		}
	}(resCh, errCh)

	select {
	case res := <-resCh:
		return res, nil
	case err := <-errCh:
		return "", err
	case <-time.After(2 * time.Minute):
		return "", errors.New("error getting captcha response - process timed out")
	}
}

func (c *capmonster) GetBalance() {}
