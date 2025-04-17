package adidas

import (
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/utils/psychoclient"
)

type taskPayload struct {
	taskID models.TaskID

	basketID       string
	consentVersion string

	shippingID         string
	shipmentID         string
	shipNode           string
	carrierCode        string
	deliveryPeriod     string
	collectionPeriod   string
	carrierServiceCode string

	session psychoclient.Session
	headers map[string]string

	worker           *models.TaskWorker
	taskGroupSetting *models.TaskGroupSetting
}
