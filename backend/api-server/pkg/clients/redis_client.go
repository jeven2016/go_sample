package clients

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
	"go.uber.org/zap"

	"api/pkg/common"
)

type RedisClient struct {
	Client *redis.Client
	Config *common.RedisConfig
	Log    *zap.Logger
}

func (c *RedisClient) StartInit() error {
	var redisCfg = c.Config
	client := redis.NewClient(&redis.Options{
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
		c.Log.Error("Cannot connect to redis", zap.Error(err))
		return err
	}
	c.Client = client
	return nil
}

// func (c *RedisClient) GetList(key string) (any, error) {
//
// 	lRange := c.Client.LRange(context.Background(), key, 0, -1)
//
// }
