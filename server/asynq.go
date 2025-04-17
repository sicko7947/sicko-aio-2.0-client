package server

import (
	"log"
	"time"

	"github.com/hibiken/asynq"
	"sicko-aio-2.0-client/tasks/defineder"
	"sicko-aio-2.0-client/tasks/handler/adidas"
	"sicko-aio-2.0-client/tasks/handler/lanecrawford"
	"sicko-aio-2.0-client/tasks/handler/louisvuitton"
	"sicko-aio-2.0-client/tasks/handler/macys"
	"sicko-aio-2.0-client/tasks/handler/mrporter"
	nike_monitor "sicko-aio-2.0-client/tasks/handler/nike/monitor"
)

const (
	redisURL      = "localhost:6379"
	concurrentNum = 99999
)

// WorkerServer : Start a worker server
func WorkerServer() {
	r := asynq.RedisClientOpt{Addr: redisURL}
	srv := asynq.NewServer(r, asynq.Config{
		Concurrency: concurrentNum,
		RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
			return 1 * time.Microsecond
		},
	})

	mux := asynq.NewServeMux()
	// Task handlers

	// NIKE
	mux.HandleFunc(defineder.NikeMonitorGRPC, nike_monitor.HandleMonitorTaskWithOptions)

	mux.HandleFunc(defineder.NikeCheckoutPrepare, nike_checkout.HandleCheckoutPrepareTaskWithOptions)
	mux.HandleFunc(defineder.NikeCheckoutV2, nike_checkout.HandleCheckoutV2TaskWithOptions)
	mux.HandleFunc(defineder.NikeCheckoutV3, nike_checkout.HandleCheckoutV3TaskWithOptions)

	mux.HandleFunc(defineder.NikeCheckoutLegacy, nike_checkout.HandleCheckoutLegacyTaskWithOptions)
	mux.HandleFunc(defineder.NikeCheckoutLegacyV2, nike_checkout.HandleCheckoutLegacyV2TaskWithOptions)
	mux.HandleFunc(defineder.NikeCheckoutLegacyV3, nike_checkout.HandleCheckoutLegacyV3TaskWithOptions)
	mux.HandleFunc(defineder.NikeCheckoutReserveStock, nike_checkout.HandleReserveStockTaskWithOptions)

	// SNKRS
	mux.HandleFunc(defineder.SnkrsLaunchEntry, snkrs.HandleLaunchEntryTaskWithOptions)
	mux.HandleFunc(defineder.SnkrsLaunchEntryPrepare, snkrs.HandleLaunchPrepareTaskWithOptions)

	// ADIDAS
	mux.HandleFunc(defineder.AdidasTask, adidas.HandleAdidasTaskWithOptions)

	// MACYS
	mux.HandleFunc(defineder.MacysMonitor, macys.HandleMacysMonitorTaskWithOptions)
	mux.HandleFunc(defineder.MacysCheckoutPrepare, macys.HandleMacyCheckoutPrepareTaskWithOptions)
	mux.HandleFunc(defineder.MacysCheckout, macys.HandleMacyCheckoutTaskWithOptions)

	// MRPORTER
	mux.HandleFunc(defineder.MrPorterMonitor, mrporter.HandleMrporterMonitorTaskWithOptions)
	mux.HandleFunc(defineder.MrPorterCheckoutPrepare, mrporter.HandleMrPorterCheckoutPrepareTaskWithOptions)
	mux.HandleFunc(defineder.MrPorterCheckout, mrporter.HandleMrporterCheckoutTaskWithOptions)

	// LOUVIS VUITTON
	mux.HandleFunc(defineder.LouisVuittonMonitor, louisvuitton.HandleLouisVuittonMonitorTaskWithOptions)
	mux.HandleFunc(defineder.LouisVuittonCheckoutPrepare, louisvuitton.HandleLouisVuittonCheckoutPrepareTaskWithOptions)
	mux.HandleFunc(defineder.LouisVuittonCheckout, louisvuitton.HandleLouisVuittonCheckoutTaskWithOptions)

	// LANECRAWFORD
	mux.HandleFunc(defineder.LanecrawfordMonitor, lanecrawford.HandleLanecrawfordMonitorTaskWithOptions)
	mux.HandleFunc(defineder.LanecrawfordCheckoutPrepare, lanecrawford.HandleLanecrawfordCheckoutPrepareTaskWithOptions)
	mux.HandleFunc(defineder.LanecrawfordCheckout, lanecrawford.HandleLanecrawfordCheckoutTaskWithOptions)

	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}
