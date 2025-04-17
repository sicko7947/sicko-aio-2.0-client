package adidas

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/utils/psychoclient"
)

func (p *taskPayload) doGetAvailability() *models.Error {
	endpoint := fmt.Sprintf("https://%s/xxxx%s/xxxxxx?%v", p.taskGroupSetting.Domain, p.worker.Product.StyleColor, rand.Float64())

	reqID, _ := p.session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "GET",
		Headers:  p.headers,
		Payload:  nil,
	})

	res, respBody, err := p.session.Do(reqID)
	if err != nil {
		return err
	}
	result := gjson.Parse(string(respBody))
	switch res.StatusCode {
	case 200:
		if obj := result.Get("variation_list.#(availability>0)#").Array(); len(obj) > 0 {
			o := obj[rand.Intn(len(obj))]
			p.worker.Product.SkuID = o.Get("sku").String()
			p.worker.Product.Size = o.Get("size").String()
			return nil
		}
		return nil
	default:
		return &models.Error{Error: errors.New("ERROR_GETTING_SIZES"), Code: 508, Message: "Error Getting Sizes"}
	}
}

func (p *taskPayload) doGetWaitingRoomConfig() *models.Error {
	endpoint := fmt.Sprintf("https://www.adidas.com/hxxxx/%s/xxxx", p.taskGroupSetting.Country)

	reqID, _ := p.session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "GET",
		Headers:  p.headers,
		Payload:  nil,
	})

	res, respBody, err := p.session.Do(reqID)

	if err != nil {
		return err
	}

	result := gjson.Parse(string(respBody))
	switch res.StatusCode {
	case 200:
		if obj := result.Get("statusMessageKey"); obj.Exists() {
			if obj.String() == "sale_started" {
				return nil
			}
		}
		return &models.Error{Error: nil, Code: 302, Message: "SALE_NOT_LIVE"}
	default:
		return &models.Error{Error: errors.New("ERROR_GETTING_WAITINGROOM_CONFIG"), Code: 508, Message: "Error Getting WaitingRoom Config"}
	}
}

func (p *taskPayload) doGetQueue() *models.Error {
	endpoint := fmt.Sprintf("https://%s/xxxxxxx", p.taskGroupSetting.Domain)

	reqID, _ := p.session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "GET",
		Headers:  p.headers,
		Payload:  nil,
	})

	res, _, err := p.session.Do(reqID)
	if err != nil {
		return err
	}
	switch res.StatusCode {
	case 200, 201, 202:
		return nil
	case 418:
		return &models.Error{Error: errors.New("STILL_IN_QUEUE"), Code: 418, Message: "Still in Queue"}
	default:
		return &models.Error{Error: errors.New("ERROR_GETTING_QUEUE"), Code: 508, Message: "Error Getting Queue"}
	}
}

func (p *taskPayload) doATC() *models.Error {
	endpoint := fmt.Sprintf("https://%s/xxxxxxx", p.taskGroupSetting.Domain)

	data, _ := json.Marshal([]map[string]interface{}{})

	reqID, _ := p.session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "POST",
		Headers:  p.headers,
		Payload:  bytes.NewBuffer(data),
	})

	res, respBody, err := p.session.Do(reqID)

	if err != nil {
		return err
	}

	result := gjson.Parse(string(respBody))
	switch res.StatusCode {
	case 200, 201, 202:
		p.headers["authorization"] = res.Header.Get("authorization")
		p.basketID = result.Get("basketId").String()
		p.consentVersion = result.Get("consentVersion").String()
		return nil
	default:
		return &models.Error{Error: errors.New("ERROR_ADDING_TO_CART"), Code: 509, Message: "Error Adding To Cart"}
	}
}

func (p *taskPayload) doGetShippingMethod() *models.Error {
	endpoint := fmt.Sprintf(`https://%s/xxxxx/%s/xxxxxxx`, p.taskGroupSetting.Domain, p.basketID)

	reqID, _ := p.session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "GET",
		Headers:  p.headers,
		Payload:  nil,
	})

	res, respBody, err := p.session.Do(reqID)

	if err != nil {
		return err
	}

	result := gjson.Parse(string(respBody))
	switch res.StatusCode {
	case 200, 201, 202:
		p.shippingID = result.Get("0.id").String()
		p.shipmentID = result.Get("0.shipmentId").String()
		p.carrierCode = result.Get("0.carrierCode").String()
		p.carrierServiceCode = result.Get("0.carrierServiceCode").String()

		collectionTimestampFrom := time.Unix(result.Get("0.collection.from").Int(), 0).Format(time.RFC3339)
		collectionTimestampTo := time.Unix(result.Get("0.collection.to").Int(), 0).Format(time.RFC3339)

		deliveryTimestampFrom := time.Unix(result.Get("0.delivery.from").Int(), 0).Format(time.RFC3339)
		deliveryTimestampTo := time.Unix(result.Get("0.delivery.to").Int(), 0).Format(time.RFC3339)

		p.collectionPeriod = collectionTimestampFrom + "," + collectionTimestampTo
		p.deliveryPeriod = deliveryTimestampFrom + "," + deliveryTimestampTo
		return nil
	default:
		return &models.Error{Error: errors.New("ERROR_ADDING_TO_CART"), Code: 509, Message: "Error Adding To Cart"}
	}
}

func (p *taskPayload) doSubmitAddress() *models.Error {
	endpoint := fmt.Sprintf("https://%s/xxxxx%s", p.taskGroupSetting.Domain, p.basketID)
	data, _ := json.Marshal(map[string]interface{}{
		"customer": map[string]interface{}{
			"email":             p.worker.TaskInfo.Profile.Email,
			"receiveSmsUpdates": false,
		},
		"billingAddress": map[string]string{
			"country":        p.worker.TaskInfo.Profile.ShippingCountry,
			"firstName":      p.worker.TaskInfo.Profile.ShippingFirstName,
			"lastName":       p.worker.TaskInfo.Profile.ShippingLastName,
			"zipcode":        p.worker.TaskInfo.Profile.ShippingZipCode,
			"address1":       p.worker.TaskInfo.Profile.ShippingAddress1,
			"houseNumber":    "0",
			"colony":         p.worker.TaskInfo.Profile.ShippingCity,
			"suburb":         p.worker.TaskInfo.Profile.ShippingCity,
			"city":           p.worker.TaskInfo.Profile.ShippingCity,
			"countyProvince": p.worker.TaskInfo.Profile.ShippingState,
			"phoneNumber":    p.worker.TaskInfo.Profile.ShippingPhoneNumber,
		},
		"methodList": []map[string]string{
			{
				"id":                 p.shippingID,
				"shipmentId":         p.shipmentID,
				"carrierCode":        p.carrierCode,
				"carrierServiceCode": p.carrierServiceCode,
				"shipNode":           p.shipNode,
				"collectionPeriod":   p.collectionPeriod,
				"deliveryPeriod":     p.deliveryPeriod,
			},
		},
		"newsletterSubscription": true,
		"consentVersion":         p.consentVersion,
	})

	reqID, _ := p.session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "PATCH",
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
		return &models.Error{Error: errors.New("ERROR_SUBMITTING_ADDRESS"), Code: 509, Message: "Error Submitting Address"}
	}
}

func (p *taskPayload) doSubmitOrder() *models.Error {
	endpoint := fmt.Sprintf("https://%s/", p.taskGroupSetting.Domain)

	expireMonth, _ := strconv.Atoi(p.worker.TaskInfo.Profile.ExpireMonth)
	expireYear, _ := strconv.Atoi(p.worker.TaskInfo.Profile.ExpireYear)
	data, _ := json.Marshal(map[string]interface{}{
		"basketId":            p.basketID,
		"encryptedInstrument": "xxxxxxxxxxxxxxxxxxxx",
		"fingerprint":         "xxxxxxxxxxxxxxxxxxxx",
		"paymentInstrument": map[string]interface{}{
			"holder":          p.worker.TaskInfo.Profile.ShippingFirstName + " " + p.worker.TaskInfo.Profile.ShippingLastName,
			"expirationMonth": expireMonth,
			"expirationYear":  expireYear,
			"lastFour":        p.worker.TaskInfo.Profile.CreditCardNumber[len(p.worker.TaskInfo.Profile.CreditCardNumber)-4:],
			"paymentMethodId": "CREDIT_CARD",
			"cardType":        p.worker.TaskInfo.Profile.CardType,
		},
	})

	reqID, _ := p.session.BuildRequest(&psychoclient.RequestBuilder{
		Endpoint: endpoint,
		Method:   "PATCH",
		Headers:  p.headers,
		Payload:  bytes.NewBuffer(data),
	})

	res, respBody, err := p.session.Do(reqID)

	if err != nil {
		return err
	}

	result := gjson.Parse(string(respBody))
	switch res.StatusCode {
	case 200, 201, 202:
		return nil
	default:
		errorMessage := result.Get("message").String()
		return &models.Error{Error: errors.New("ERROR_SUBMITTING_ORDER"), Code: 509, Message: errorMessage}
	}

}
