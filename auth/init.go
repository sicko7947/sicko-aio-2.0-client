package auth

import (
	"time"

	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/utils/psychoclient"
)

var (
	hasLogin bool
)

func init() {
	if communicator.DEV_ENV {
		hasLogin = true
		return
	}
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		for {
			if !hasLogin {
				continue
			}
			polling()
			<-ticker.C
		}
	}()
}

func getIpAddress() string {
	_, respBody, err := psychoclient.NewClient(new(psychoclient.SessionBuilder)).RoundTripNewRequest(&psychoclient.RequestBuilder{
		Endpoint: "https://api.ipify.org/",
		Method:   "GET",
	})
	if err != nil {
		return ""
	}
	return string(respBody)
}
