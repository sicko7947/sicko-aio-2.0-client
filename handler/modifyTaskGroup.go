package handler

import (
	"net/http"
	"strings"

	"github.com/gogf/gf/net/ghttp"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/constants"
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/tasks/defineder"
)

// ModifyTaskGroup : Modify TaskpGroup
func ModifyTaskGroup(r *ghttp.Request) {

	// Parsing Payload
	var taskGroup *models.TaskGroup
	if err := r.Parse(&taskGroup); err != nil {
		r.Response.WriteStatus(http.StatusBadRequest)
		r.Response.WriteJsonExit(map[string]bool{"success": false})
		return
	}

	switch taskGroup.TaskGroupSetting.Category {
	case models.NIKE:
		countryInfo := constants.GetNikeCountryInfo(taskGroup.TaskGroupSetting.Country) // getting basic task country information
		taskGroup.TaskGroupSetting.MerchGroup = countryInfo.MerchGroup
		taskGroup.TaskGroupSetting.Language = countryInfo.Language
		taskGroup.TaskGroupSetting.Currency = countryInfo.Currency
		taskGroup.TaskGroupSetting.Locale = countryInfo.Locale

		// distribute different task types by merchgroups
		switch countryInfo.MerchGroup {
		case "XP", "XA":
			taskGroup.TaskGroupSetting.TaskType = defineder.NikeCheckoutLegacy
		case "MX":
			taskGroup.TaskGroupSetting.TaskType = defineder.NikeCheckoutLegacyV2
		case "JP":
			taskGroup.TaskGroupSetting.TaskType = defineder.NikeCheckoutV2
		case "US", "EU", "CN":
			taskGroup.TaskGroupSetting.TaskType = defineder.NikeCheckoutV3
		}
		communicator.Config.TaskGroups[taskGroup.GroupID] = taskGroup

	case models.SNKRS:
		countryInfo := constants.GetNikeCountryInfo(taskGroup.TaskGroupSetting.Country) // getting basic task country information
		taskGroup.TaskGroupSetting.MerchGroup = countryInfo.MerchGroup
		taskGroup.TaskGroupSetting.Language = countryInfo.Language
		taskGroup.TaskGroupSetting.Currency = countryInfo.Currency
		taskGroup.TaskGroupSetting.Locale = countryInfo.Locale

		taskGroup.TaskGroupSetting.TaskType = defineder.SnkrsLaunchEntry
		communicator.Config.TaskGroups[taskGroup.GroupID] = taskGroup

	case models.ADIDAS:
		countryInfo := constants.GetAdidasCountryInfo(taskGroup.TaskGroupSetting.Country) // getting basic task country information
		taskGroup.TaskGroupSetting.Domain = countryInfo.Domain
		taskGroup.TaskGroupSetting.TaskType = defineder.AdidasTask
		communicator.Config.TaskGroups[taskGroup.GroupID] = taskGroup

	case models.MRPORTER:
		taskGroup.TaskGroupSetting.Locale = strings.ToLower(taskGroup.TaskGroupSetting.Country)
		taskGroup.TaskGroupSetting.TaskType = defineder.MrPorterCheckout
		communicator.Config.TaskGroups[taskGroup.GroupID] = taskGroup

	case models.LUISAVIAROMA:
		countryInfo := constants.GetNikeCountryInfo(taskGroup.TaskGroupSetting.Country) // getting basic task country information
		taskGroup.TaskGroupSetting.Currency = countryInfo.Currency
		taskGroup.TaskGroupSetting.Language = countryInfo.Language
		taskGroup.TaskGroupSetting.TaskType = defineder.LuisaviaromaCheckout
		communicator.Config.TaskGroups[taskGroup.GroupID] = taskGroup

	case models.MACYS:
		taskGroup.TaskGroupSetting.TaskType = defineder.MacysCheckout
		communicator.Config.TaskGroups[taskGroup.GroupID] = taskGroup

	case models.GAMESTOP:
		taskGroup.TaskGroupSetting.TaskType = defineder.GamestopMonitor
		communicator.Config.TaskGroups[taskGroup.GroupID] = taskGroup

	case models.LANECRAWFORD:
		taskGroup.TaskGroupSetting.TaskType = defineder.LanecrawfordCheckout
		communicator.Config.TaskGroups[taskGroup.GroupID] = taskGroup

	case models.LOUISVUITTON:
		countryInfo := constants.GetLouisVuittonCountryInfo(taskGroup.TaskGroupSetting.Country) // getting basic task country information
		taskGroup.TaskGroupSetting.Domain = countryInfo.Domain
		taskGroup.TaskGroupSetting.TaskType = defineder.LouisVuittonCheckout
		communicator.Config.TaskGroups[taskGroup.GroupID] = taskGroup

	}
	r.Response.WriteJsonExit(map[string]bool{"success": true})
}
