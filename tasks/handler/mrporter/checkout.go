package mrporter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"

	http "github.com/zMrKrabz/fhttp"

	"github.com/PuerkitoBio/goquery"
	"github.com/hibiken/asynq"
	"github.com/sicko7947/sickocommon"
	"github.com/tidwall/gjson"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/successHandler"
	"sicko-aio-2.0-client/tasks"
	"sicko-aio-2.0-client/utils/notification"
	"sicko-aio-2.0-client/utils/psychoclient"
	"sicko-aio-2.0-client/utils/redis"
)

type checkoutPayload struct {
	billinAddressID          int64
	creditCardToken          string
	creditCardType           string
	creditCardExpireMonth    string
	creditCardExpireYear     string
	creditCardLastFourDigits string

	applyaccountinfo bool
	usepaypal        bool

	paReq              string
	paRes              string
	md                 string
	redirectUrl        string
	headers            map[string]string
	checkoutProxyGroup []string
	session            psychoclient.Session

	taskID           models.TaskID
	taskGroupSetting *models.TaskGroupSetting
	worker           *models.TaskWorker
}

func (p checkoutPayload) doAddToCart() *models.Error {
	communicator.ModifyTaskStatus(&models.Message{Code: 200, GroupID: p.worker.GroupID, TaskID: p.taskID, Status: "PENDING", Message: "Adding To Cart"})
	endpoint := fmt.Sprintf("https://www.mrporter.com/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", p.taskGroupSetting.Locale)

	data, _ := json.Marshal(map[string]interface{}{})

	reqID, _ := p.session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "POST",
		Headers:  p.headers,
		Payload:  bytes.NewBuffer(data),
	})

	res, _, err := p.session.Do(reqID)
	if err != nil {
		return err
	}
	switch res.StatusCode {
	case 200, 201, 202:
		return nil
	default:
		return &models.Error{Error: errors.New("ERROR_ADDING_TO_CART"), Code: res.StatusCode, Message: "Error Adding To Cart"}
	}
}

func (p *checkoutPayload) doVerifyCreditCard() *models.Error {
	communicator.ModifyTaskStatus(&models.Message{Code: 200, GroupID: p.worker.GroupID, TaskID: p.taskID, Status: "PENDING", Message: "Submitting CC Details"})

	endpoint := "https://xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

	data, _ := json.Marshal(map[string]interface{}{})

	session, err := psychoclient.NewSession(&psychoclient.SessionBuilder{
		Proxy: sickocommon.GetProxy(p.checkoutProxyGroup).String(),
	})
	if err != nil {
		return err
	}

	reqId, _ := session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "POST",
		Headers:  p.headers,
		Payload:  bytes.NewBuffer(data),
	})

	res, respBody, err := session.Do(reqId)
	if err != nil {
		return err
	}
	result := gjson.ParseBytes(respBody)
	switch res.StatusCode {
	case 200, 201, 202:
		if obj := result.Get("cardId"); obj.Exists() {
			p.creditCardToken = result.Get("cardId").String()
			return nil
		}
		fallthrough
	default:
		return &models.Error{Error: errors.New("ERROR_GETTING_CREDITCARD_TOKEN"), Code: res.StatusCode, Message: "Error Getting CC Token"}
	}
}

func (p *checkoutPayload) doApplyAccountInfo() *models.Error {
	communicator.ModifyTaskStatus(&models.Message{Code: 200, GroupID: p.worker.GroupID, TaskID: p.taskID, Status: "PENDING", Message: "Applying Account Info"})

	endpoint := fmt.Sprintf("https://xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", p.taskGroupSetting.Locale)

	data, _ := json.Marshal(make(map[string]string))

	reqID, _ := p.session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "PUT",
		Headers:  p.headers,
		Payload:  bytes.NewBuffer(data),
	})

	res, _, err := p.session.Do(reqID)
	if err != nil {
		return err
	}
	switch res.StatusCode {
	case 200, 201, 202:
		fallthrough
	default:
		return &models.Error{Error: errors.New("ERROR_APPLYING_ACCOUNT_INFORMATION"), Code: res.StatusCode, Message: "Error Applying Account Information"}
	}
}

func (p *checkoutPayload) doPostPaymentInstructions() *models.Error {
	communicator.ModifyTaskStatus(&models.Message{Code: 200, GroupID: p.worker.GroupID, TaskID: p.taskID, Status: "PENDING", Message: "Submitting Payment Info"})

	endpoint := fmt.Sprintf("https://xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", p.taskGroupSetting.Locale)

	var data []byte

	switch p.usepaypal {
	case true:

	case false:

	}

	reqID, _ := p.session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "POST",
		Headers:  p.headers,
		Payload:  bytes.NewBuffer(data),
	})

	res, _, err := p.session.Do(reqID)
	if err != nil {
		return err
	}
	switch res.StatusCode {
	case 201:
		return nil
	default:
		return &models.Error{Error: errors.New("ERROR_SUBMITTING_PAYMENT_METHOD"), Code: res.StatusCode, Message: "Error Submitting Payment Method"}
	}
}

func (p *checkoutPayload) doSubmitPayment() (int, *models.Error) {
	communicator.ModifyTaskStatus(&models.Message{Code: 200, GroupID: p.worker.GroupID, TaskID: p.taskID, Status: "PENDING", Message: "Submitting Payment"})

	endpoint := fmt.Sprintf("https://xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", p.taskGroupSetting.Locale)
	data, _ := json.Marshal(map[string]string{})

	reqID, _ := p.session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "POST",
		Headers:  p.headers,
		Payload:  bytes.NewBuffer(data),
	})

	res, respBody, err := p.session.Do(reqID)
	if err != nil {
		return 0, err
	}
	result := gjson.ParseBytes(respBody)
	switch res.StatusCode {

	case 201: // successfully checked out
		p.worker.TaskInfo.OrderID = result.Get("shoppingCart.customerOrderNumber").String()
		return 201, nil

	case 202: // create 3D checkout request
		obj := result.Get("redirectParameters.0")
		if obj.Exists() {
			p.paReq = obj.Get("PaReq").String()
			p.md = obj.Get("MD").String()
			p.worker.TaskInfo.RedirectURL = obj.Get("redirectUrl").String()
		}
		return 202, nil

	case 400: // error submitting checkout (mostly payment failed)
		obj := result.Get("errors")
		if obj.Exists() {
			errorMessage := obj.Get("0.errorMessage").String()
			return 0, &models.Error{Error: errors.New("ERROR_SUBMITTING_ORDER"), Code: 509, Message: errorMessage}
		}
		fallthrough

	default: // some other unexpected error
		return 0, &models.Error{Error: errors.New("ERROR_SUBMITTING_PAYMENT"), Code: res.StatusCode, Message: "Error Submitting Payment"}
	}
}

func (p *checkoutPayload) doVerify3D() *models.Error {
	communicator.ModifyTaskStatus(&models.Message{Code: 200, GroupID: p.worker.GroupID, TaskID: p.taskID, Status: "PENDING", Message: "Verifying 3D"})

	reqPayload := &bytes.Buffer{}
	writer := multipart.NewWriter(reqPayload)
	data := map[string]string{
		"PaReq":   p.paReq,
		"TermUrl": fmt.Sprintf("https://wxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxe", p.taskGroupSetting.Locale),
		"MD":      p.md,
	}

	for k, v := range data {
		writer.WriteField(k, v)
	}
	writer.Close()

	headers := map[string]string{
		"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"accept-language":           "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7,es;q=0.6",
		"cache-control":             "no-cache",
		"content-type":              writer.FormDataContentType(),
		"dnt":                       "1",
		"sec-fetch-dest":            "document",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-site":            "same-origin",
		"sec-fetch-user":            "?1",
		"upgrade-insecure-requests": "1",
	}

	reqID, _ := p.session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: p.redirectUrl,
		Method:   "POST",
		Headers:  headers,
		Payload:  reqPayload,
	})

	res, respBody, err := p.session.Do(reqID)
	if err != nil {
		return err
	}

	document, e := goquery.NewDocumentFromReader(bytes.NewReader(respBody)) // Load the HTML document
	if e != nil {
		return &models.Error{Error: errors.New("ERROR_GETTING_3D_RESPONSE"), Code: 500, Message: "ERROR_PARSING_RESPONSE_BODY"}
	}
	switch res.StatusCode {
	case 202:
		paRes, exist := document.Find("input").First().Attr("value")
		if exist {
			p.paRes = paRes
			return nil
		}
		fallthrough
	default:
		return &models.Error{Error: errors.New("ERROR_GETTING_3D_RESPONSE"), Code: res.StatusCode, Message: "Error Getting 3D Verification Response"}
	}
}

func (p *checkoutPayload) doPutPaymentInstructions() *models.Error {
	communicator.ModifyTaskStatus(&models.Message{Code: 200, GroupID: p.worker.GroupID, TaskID: p.taskID, Status: "PENDING", Message: "Submitting Payment Instructions"})

	endpoint := fmt.Sprintf("https://xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", p.taskGroupSetting.Locale)
	data, _ := json.Marshal(map[string]interface{}{})

	reqID, _ := p.session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "PUT",
		Headers:  p.headers,
		Payload:  bytes.NewBuffer(data),
	})

	res, _, err := p.session.Do(reqID)
	if err != nil {
		return err
	}
	switch res.StatusCode {
	case 200:
		return nil
	default:
		return &models.Error{Error: errors.New("ERROR_SUBMITTING_ORDER"), Code: res.StatusCode, Message: "Error Submitting Order"}
	}
}

func (p *checkoutPayload) setupTask() (err *models.Error) {

	p.applyaccountinfo = false
	p.usepaypal = true

	// setup task group settings
	p.taskGroupSetting = communicator.Config.TaskGroups[p.worker.GroupID].TaskGroupSetting

	p.checkoutProxyGroup = (communicator.Config.Proxies)[p.worker.CheckoutProxyGroupName]
	if p.session, err = psychoclient.NewSession(&psychoclient.SessionBuilder{
		Proxy: sickocommon.GetProxy(p.checkoutProxyGroup).String(),
	}); err != nil {
		return err
	}

	// setup session cookie
	useragent, cookieMap, err := redis.GetCookie2FromRedis("akamai.mrporter.com")
	if err != nil {
		return err
	}
	p.session.SetCookies(cookieMap)

	p.worker.TaskInfo.CreditCardToken = sickocommon.NikeUUID()

	p.headers = map[string]string{
		"accept":              "*/*",
		"accept-encoding":     "gzip, deflate, br",
		"accept-language":     "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7,es;q=0.6",
		"application-name":    "Blue lobster",
		"application-version": "4.157.0",
		"content-type":        "text/plain;charset=UTF-8",
		"dnt":                 "1",
		"label":               "addToCart",
		"user-agent":          useragent,
	}

	if p.worker.TaskInfo.Account != nil {
		p.worker.TaskInfo.Email = p.worker.TaskInfo.Account.Email
		p.headers["x-ubertoken"] = p.worker.TaskInfo.Account.AccessToken
		p.session.SetCookies(map[string]*http.Cookie{
			"ubertoken": {
				Name:  "ubertoken",
				Value: p.worker.TaskInfo.Account.AccessToken,
			},
		})
	}

	switch {
	case p.worker.TaskInfo.Profile == nil:
		return &models.Error{Error: errors.New("ERROR_ASSIGN_PROFILE"), Code: 500, Message: "Error Assign Profile"}
	case p.worker.TaskInfo.Account == nil:
		return &models.Error{Error: errors.New("ERROR_ASSIGN_ACCOUNT"), Code: 500, Message: "Error Assign Account"}
	}

	return nil
}

// HandleMrporterCheckoutTaskWithOptions : HandleMrporterCheckoutTaskWithOptions
func HandleMrporterCheckoutTaskWithOptions(ctx context.Context, task *asynq.Task) error {

	// Get basic task info
	p := &checkoutPayload{
		taskID: models.TaskID(task.Payload()),
	}

	worker := communicator.TaskWorkerObjectGMap.Get(p.taskID)
	if worker == nil {
		return errors.New("error starting worker")
	}

	p.worker = worker.(*models.TaskWorker)
	if !p.worker.Mutex.TryLock() {
		return errors.New("too many tasks assigned to one worker")
	}

	defer func() {
		if p.worker.Mutex.IsLocked() {
			p.worker.Mutex.Unlock()
		}
	}()

	done := make(chan bool)
	errCh := make(chan *models.Error, 1)
	defer close(errCh)
	defer close(done)
	go func() {
		var err *models.Error
		defer func() {
			recover()
		}()

		err = p.setupTask()
		if tasks.SafeSend(errCh, err) {
			return
		}

		err = p.doAddToCart()
		if tasks.SafeSend(errCh, err) {
			return
		}

		err = p.doApplyAccountInfo()
		if tasks.SafeSend(errCh, err) {
			return
		}

		if !p.usepaypal {
			err = p.doVerifyCreditCard()
			if tasks.SafeSend(errCh, err) {
				return
			}
		}

		err = p.doPostPaymentInstructions()
		if tasks.SafeSend(errCh, err) {
			return
		}

		code, err := p.doSubmitPayment()
		if err != nil { // First submit payment is for creating 3D checkout request
			if tasks.SafeSend(errCh, err) {
				return
			}

			if code == 202 && !p.usepaypal {
				err = p.doVerify3D()
				if tasks.SafeSend(errCh, err) {
					return
				}

				err = p.doPutPaymentInstructions()
				if tasks.SafeSend(errCh, err) {
					return
				}

				_, err = p.doSubmitPayment() // final submit payment is for checking order
				if tasks.SafeSend(errCh, err) {
					return
				}
			}
		}
		done <- true
	}()

	for {
		select {
		case <-ctx.Done():
			communicator.ModifyTaskStatus(&models.Message{ // send task cancellation status to frontend
				Code:    0,
				GroupID: p.worker.GroupID,
				TaskID:  p.taskID,
				Status:  "CANCELLED",
				Message: "Stopped",
			})
			return ctx.Err() // cancelation signal received, abandon this work.
		case err := <-errCh:
			if err == nil {
				continue
			}
			communicator.ModifyTaskStatus(&models.Message{ // send error message to frontend
				Code:    err.Code,
				GroupID: p.worker.GroupID,
				TaskID:  p.taskID,
				Status:  "ARCHIVED",
				Message: err.Message,
			})
			return err.Error
		case <-done:
			status := "COMPLETED"
			if p.worker.Restart {
				status = "ARCHIVED"
			}
			communicator.ModifyTaskStatus(&models.Message{ // send success checkout message to frontend
				GroupID: p.worker.GroupID,
				TaskID:  p.taskID,
				Code:    200,
				Status:  status,
				Message: "Checked Out!",
			})

			// send success checkout message to webhooks
			go successHandler.HandlerSuccess(p.taskGroupSetting, p.worker)
			go notification.PushNotification(p.taskGroupSetting, p.worker)
			return nil
		}
	}
}
