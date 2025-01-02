package ioc

import (
	"github.com/sqkam/goproxy/config"
	"github.com/sqkam/goproxy/pkg/proxy"
)

func InitProxyServer(conf *config.ProxyConfig) proxy.Service {
	return proxy.NewDefaultServer(conf)
}
