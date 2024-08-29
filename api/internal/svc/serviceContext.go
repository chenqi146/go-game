package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go-game/api/internal/config"
	"strings"
)

type ServiceContext struct {
	Config      config.Config
	RedisClient *redis.Redis
	WordPairs   []string
}

func NewServiceContext(c config.Config, w string) *ServiceContext {
	return &ServiceContext{
		Config: c,
		RedisClient: redis.MustNewRedis(redis.RedisConf{
			Host: c.Redis.Host,
			Type: c.Redis.Type,
			Pass: c.Redis.Pass,
			Tls:  c.Redis.Tls,
		}),
		WordPairs: strings.Fields(w),
	}
}
