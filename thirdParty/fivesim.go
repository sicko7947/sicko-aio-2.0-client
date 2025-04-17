package thirdParty

import (
	fivesim "github.com/sicko7947/5sim-go"
)

const (
	APIKEY = ""
)

var (
	client = fivesim.NewClient(APIKEY, "")
)

func GetFivesimClient() fivesim.Client {
	return client
}
