package notification

import (
	"fmt"

	"github.com/go-toast/toast"
	"github.com/huandu/go-clone"
	"sicko-aio-2.0-client/models"
)

func PushNotification(taskGroupSetting *models.TaskGroupSetting, worker *models.TaskWorker) {
	newWorker := clone.Clone(worker).(*models.TaskWorker)
	if newWorker == nil || newWorker.TaskInfo == nil || newWorker.Product == nil {
		return
	}

	notification := toast.Notification{
		AppID:   "SICKO AIO 2.0",
		Title:   "Successfully Checked Out !!!",
		Message: fmt.Sprintf("[%s][%s] %s\nSize: %s", taskGroupSetting.Category, taskGroupSetting.Country, newWorker.Product.ProductName, newWorker.Product.Size),
		Icon:    "C:\\Users\\sicko\\Desktop\\aaaa.png",
	}
	notification.Push()

	switch taskGroupSetting.Category {
	case models.NIKE:
		switch taskGroupSetting.MerchGroup {
		case "XP", "XA":
			sendDiscordNikeLegacyCheckoutWebhook(taskGroupSetting, newWorker)
			sendSlackNikeLegacyCheckoutWebhook(taskGroupSetting, newWorker)
		default:
			sendDiscordNikeACOCheckoutWebhook(taskGroupSetting, newWorker)
			sendSlackNikeACOCheckoutWebhook(taskGroupSetting, newWorker)
		}

	case models.ADIDAS:
	case models.FINISHLINE:
	case models.GAMESTOP:
	case models.MACYS:
	case models.SSENSE:
		sendDiscordSsenseCheckoutWebhook(taskGroupSetting, newWorker)
		sendSlackSsenseCheckoutWebhook(taskGroupSetting, newWorker)
	case models.LUISAVIAROMA:
		sendDiscordLuisaviaromaCheckoutWebhook(taskGroupSetting, newWorker)
		sendSlackLuisaviaromaCheckoutWebhook(taskGroupSetting, newWorker)
	case models.MRPORTER:
		sendDiscordMrPorterCheckoutWebhook(taskGroupSetting, newWorker)
		sendSlackMrPorterCheckoutWebhook(taskGroupSetting, newWorker)
	case models.PACSUN:
		sendDiscordPacsunCheckoutWebhook(taskGroupSetting, newWorker)
		sendSlackPacsunCheckoutWebhook(taskGroupSetting, newWorker)
	case models.SNEAKERBOY:
		sendDiscordSneakerboyCheckoutWebhook(taskGroupSetting, newWorker)
		sendSlackSneakerboyCheckoutWebhook(taskGroupSetting, newWorker)
	case models.SUPPLYSTORE:
	case models.YEEZYSUPPLY:
	case models.TAF:
		sendDiscordTafCheckoutWebhook(taskGroupSetting, newWorker)
		sendSlackTafCheckoutWebhook(taskGroupSetting, newWorker)
	case models.NEWBALANCE:
		sendDiscordNewBalanceCheckoutWebhook(taskGroupSetting, newWorker)
		sendSlackNewBalanceCheckoutWebhook(taskGroupSetting, newWorker)
	}
}
