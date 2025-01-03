package config

type ProxyConfig struct {
	Listen int64  `mapstructure:"listen"`
	Target string `mapstructure:"target"`
}
