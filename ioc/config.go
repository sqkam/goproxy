package ioc

import (
	"github.com/spf13/viper"
)

type ProxyConfig struct {
	Listen int64
	Target string
}

func InitConfig() *ProxyConfig {
	var conf ProxyConfig
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
