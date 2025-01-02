package ioc

import (
	"github.com/spf13/viper"
	"github.com/sqkam/goproxy/config"
)

func InitConfig() *config.ProxyConfig {
	var conf config.ProxyConfig
	v := viper.New()
	v.SetConfigFile("./config.yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(&conf); err != nil {
		panic(err)
	}

	return &conf
}
