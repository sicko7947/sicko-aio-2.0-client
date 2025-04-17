package utils

import (
	"github.com/sicko7947/sickocommon"
	"sicko-aio-2.0-client/communicator"
)

func GetSyncProxy(proxies []string) (proxy string) {
	proxy = ""
	if len(proxies) > 0 && len(proxies[0]) > 0 {
		proxy = proxies[0]
	} else {
		if proxylist := communicator.Config.Proxies["Default"]; proxylist != nil {
			proxy = sickocommon.GetProxy(proxylist).String()
		}
	}
	return proxy
}
