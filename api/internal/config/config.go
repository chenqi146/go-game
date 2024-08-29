package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	HeartbeatPeriod int64
	Redis           struct {
		Host string
		Type string
		Pass string
		Tls  bool
	}
}
