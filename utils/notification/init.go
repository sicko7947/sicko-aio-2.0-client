package notification

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/sicko7947/sickocommon"
	"github.com/withmandala/go-log"
	"sicko-aio-2.0-client/utils/psychoclient"
)

type webhook struct {
	WebhookUrl string
	Data       interface{}
}

var (
	logger   *log.Logger
	session  psychoclient.Session
	webhooks = make(chan webhook)
)

func webhookScheduler() {
	var webhooksSent = 0
	var wg sync.WaitGroup
	for webhook := range webhooks {
		wg.Add(1)
		webhooksSent++
		go discordSend(webhook, &wg)
		if webhooksSent > 2 {
			wg.Wait()
			webhooksSent--
		}
	}
}

func init() {
	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	sickocommon.PathCheckAndCreate("logs", 0777)
	f, err := os.Create(fmt.Sprintf("logs/success_%s.log", timestamp))
	if err != nil {
		return
	}
	logger = log.New(f).WithoutTimestamp()

	session, _ = psychoclient.NewSession(&psychoclient.SessionBuilder{
		UseDefaultClient: true,
	})

	go webhookScheduler()
}
