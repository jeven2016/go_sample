package initialization

import (
	"api/pkg/common"
	"context"
	"github.com/go-redis/redis/v9"
	"go.uber.org/zap"
	"time"
)

func SetupRedis(config *common.Config, log *zap.Logger) (client *redis.Client, err error) {
	var redisCfg = config.RedisConfig
	client = redis.NewClient(&redis.Options{
		Addr:         redisCfg.Address,
		Password:     redisCfg.Password,
		DB:           redisCfg.DefaultDb,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  time.Duration(redisCfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(redisCfg.WriteTimeout) * time.Second,
		PoolSize:     redisCfg.PoolSize,
		PoolTimeout:  time.Duration(redisCfg.PoolTimeout) * time.Second,
	})

	if _, err := client.Ping(context.TODO()).Result(); err != nil {
		log.Error("Cannot connect to mongodb", zap.Error(err))
	}
	return client, err
}
