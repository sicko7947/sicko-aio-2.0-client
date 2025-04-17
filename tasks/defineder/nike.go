package defineder

import (
	"github.com/hibiken/asynq"
	"sicko-aio-2.0-client/models"
)

const (
	NikeMonitorGtin          = `nike:monitorV1:information,restock(gtin),checklive`
	NikeMonitorDelivery      = `nike:monitorV1:information,restock(legacy),checklive`
	NikeMonitorGraphQL       = `nike:monitorV2:information,restock,checklive`
	NikeMonitorGRPC          = `nike:monitorV3:information,restock,checklive`
	NikeCheckoutLegacy       = `nike:checkout:legacy`
	NikeCheckoutLegacyV2     = `nike:checkout:legacy:v2`
	NikeCheckoutLegacyV3     = `nike:checkout:legacy:v3`
	NikeCheckoutReserveStock = `nike:checkout:reserve`
	NikeCheckoutPrepare      = `nike:checkout:prepare`
	NikeCheckoutV2           = `nike:checkout:v2`
	NikeCheckoutV3           = `nike:checkout:v3`

	NikeSync    = `nike:login:sync`
	NikeRefresh = `nike:login:refresh`
)

func NewNikeMonitorGtinTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(NikeMonitorGtin, []byte(taskId))
}

func NewNikeMonitorDeliveryTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(NikeMonitorDelivery, []byte(taskId))
}

func NewnNikeMonitorGraphQLTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(NikeMonitorGraphQL, []byte(taskId))
}

func NewnNikeMonitorGRPCTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(NikeMonitorGRPC, []byte(taskId))
}

func NewNikeCheckoutLegacyV2Task(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(NikeCheckoutLegacyV2, []byte(taskId))
}

func NewNikeCheckoutLegacyV3Task(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(NikeCheckoutLegacyV3, []byte(taskId))
}

func NewNikeCheckoutReserveStockTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(NikeCheckoutReserveStock, []byte(taskId))
}

func NewNikeCheckoutLegacyTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(NikeCheckoutLegacy, []byte(taskId))
}

func NewNikeCheckoutPrepare(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(NikeCheckoutPrepare, []byte(taskId))
}

func NewNikeCheckoutV2Task(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(NikeCheckoutV2, []byte(taskId))
}

func NewNikeCheckoutV3Task(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(NikeCheckoutV3, []byte(taskId))
}

func NewNikeLoginSyncTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(NikeSync, []byte(taskId))
}

func NewNikeLoginRefreshTask(taskId models.TaskID) *asynq.Task {
	return asynq.NewTask(NikeRefresh, []byte(taskId))
}
