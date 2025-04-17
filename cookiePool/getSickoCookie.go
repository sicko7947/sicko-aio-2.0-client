package cookiePool

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sicko7947/sickocommon"
	"github.com/tidwall/gjson"
	"sicko-aio-2.0-client/models"
	auth_service "sicko-aio-2.0-client/proto/auth"
	"sicko-aio-2.0-client/tasks"
)

func getCookieData(resChan chan<- *models.Cookie, errChan chan *models.Error) {
	for {
		if cookieDataStream == nil {
			continue
		}

		res, err := cookieDataStream.Recv()
		if err != nil {
			tasks.SafeSend(errChan, &models.Error{Error: errors.New("GET_STREAM_ERROR"), Code: 409, Message: fmt.Sprintf("RequestCookieData get stream err: %v", err)})
			return
		}

		encodedData := res.GetData()

		dataString, err := sickocommon.LZDecompress(encodedData, "")
		if err != nil {
			tasks.SafeSend(errChan, &models.Error{Error: errors.New("ERROR_PARSING_COOKIES"), Code: 500, Message: "Error Parsing Cookie Data"})
			return
		}

		result := gjson.Parse(dataString)
		useragent := result.Get("useragent").String()
		cookieMap := make(map[string]*http.Cookie)

		json.Unmarshal([]byte(result.Get("cookies").String()), &cookieMap)

		resChan <- &models.Cookie{
			Useragent: useragent,
			CookieMap: cookieMap,
		}
		return
	}
}

func GetCookieDate() (*models.Cookie, *models.Error) {
	if cookieDataStream == nil {
		return nil, &models.Error{Error: errors.New("GET_STREAM_ERROR"), Code: 409, Message: "RequestCookieData get stream err"}
	}
	go cookieDataStream.Send(&auth_service.StreamGetCookieDataRequest{})

	resChan := make(chan *models.Cookie, 1)
	errChan := make(chan *models.Error, 1)
	defer close(resChan)
	defer close(errChan)

	go getCookieData(resChan, errChan)

	select {
	case err := <-errChan:
		return nil, err
	case res := <-resChan:
		return res, nil
	case <-time.After(5 * time.Second):
		return nil, &models.Error{Error: errors.New("ERROR_GETTING_SICKO_COOKIE"), Code: 500, Message: "Error Getting Sicko Cookie - Request Timeout"}
	}
}
