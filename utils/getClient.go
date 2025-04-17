package utils

import (
	"time"

	"github.com/gogf/gf/os/gcache"
	"github.com/sicko7947/sickocommon"
	"sicko-aio-2.0-client/communicator"
	"sicko-aio-2.0-client/models"
	"sicko-aio-2.0-client/utils/psychoclient"
)

var cache *gcache.Cache

func init() {
	cache = gcache.New()
}

func InitOrGetSession(proxyGroupName models.ProxyGroupName) (psychoclient.Session, error) {
	proxyGroup := communicator.Config.Proxies[proxyGroupName]

	v, err := cache.GetOrSetFunc(
		proxyGroupName,
		func() (interface{}, error) {
			return producer(proxyGroup)
		}, 0,
	)

	if err != nil {
		return nil, err
	}

	ch := v.(chan psychoclient.Session)
	return <-ch, nil
}

func PutSession(name models.ProxyGroupName, s psychoclient.Session) error {
	s.Close()
	v, err := cache.Get(name)
	if err != nil {
		return err
	}
	ch := v.(chan psychoclient.Session)

	if len(ch) < cap(ch) {
		ch <- s
	}
	return nil
}

func producer(proxyGroup []string) (interface{}, error) {
	ch := make(chan psychoclient.Session, len(proxyGroup))
	go func() {
		var ttl time.Duration
		for {
			for _, v := range proxyGroup {
				switch {
				case float64(len(ch))/float64(cap(ch)) <= 0.2:
					ttl = 0 * time.Millisecond
				case float64(len(ch))/float64(cap(ch)) <= 0.4:
					ttl = 20 * time.Millisecond
				case float64(len(ch))/float64(cap(ch)) <= 0.8:
					ttl = 50 * time.Millisecond
				case float64(len(ch))/float64(cap(ch)) >= 0.8:
					ttl = 100 * time.Millisecond
				case len(ch) < cap(ch):
					continue
				}

				if sesh, err := psychoclient.NewSession(&psychoclient.SessionBuilder{
					Proxy: sickocommon.GetSingleProxy(v).String(),
				}); err == nil {
					ch <- sesh
				}
				time.Sleep(ttl)
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return ch, nil
}
