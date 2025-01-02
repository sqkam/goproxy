//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/sqkam/goproxy/ioc"
	"github.com/sqkam/goproxy/pkg/proxy"
)

func InitProxyServer() proxy.Service {
	panic(wire.Build(ioc.InitConfig, ioc.InitProxyServer))
}
