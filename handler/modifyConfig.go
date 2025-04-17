package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gogf/gf/net/ghttp"
	"github.com/huandu/go-clone"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/models"
)

// ModifyConfig : Modify Config
func ModifyConfig(r *ghttp.Request) {
	var p *models.Config

	err := json.Unmarshal(r.GetBody(), &p)
	if err != nil {
		r.Response.WriteStatusExit(http.StatusBadRequest)
		return
	}

	// overide config with new config
	if p.TaskGroups != nil {
		communicator.Config.TaskGroups = p.TaskGroups
	}
	if p.Accounts != nil {
		communicator.Config.Accounts = p.Accounts
	}

	if p.Profiles != nil {
		communicator.Config.Profiles = p.Profiles
	}

	if p.Proxies != nil {
		communicator.Config.Proxies = p.Proxies
	}

	if p.GiftCards != nil {
		communicator.Config.GiftCards = p.GiftCards
	}

	if p.Discounts != nil {
		communicator.Config.Discounts = p.Discounts
	}

	if p.Settings != nil {
		communicator.Config.Settings = p.Settings
	}

	t := clone.Clone(communicator.Config).(*models.Config) // deep copy
	r.Response.WriteJsonExit(t)
}

func GetConfig(r *ghttp.Request) {
	t := clone.Clone(communicator.Config).(*models.Config) // deep copy
	r.Response.WriteJsonExit(t)
}
