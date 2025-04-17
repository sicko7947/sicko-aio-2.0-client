package main

import (
	_ "net/http/pprof"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"sicko-aio-2.0-client/server"
	"sicko-aio-2.0-client/utils/redis"
)

func init() {
	redis.ClearOldAsynqData()
	go func() {
		for {
			redis.CheckExpireCookieInRedis()
			time.Sleep(1 * time.Minute)
		}
	}()
}

func startBackend() {
	go server.WorkerServer()
	go server.BackendServer()
}

func main() {

	gofakeit.Seed(0)
	startBackend()
	select {}
}
