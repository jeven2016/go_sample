package yzs8

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

func initLog() *zap.Logger {
	log := SetupLog()
	return log
}

func CreateMongoClient(log *zap.Logger) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 连接MongoDB
	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI("mongodb://db_user:db_pwd@localhost:27017/books?retryWrites=true&w=majority&authSource=admin&maxPoolSize=20"))

	if err != nil {
		log.Error("Cannot connect to mongodb", zap.Error(err))
		return nil, err
	}

	// 检测MongoDB是否连接成功
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Error("Failed to Ping for mongodb", zap.Error(err))
		return nil, err
	}

	return client, nil
}

func RedisClient(log *zap.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         "192.168.1.14:6379",
		Password:     "pwd",
		DB:           1,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Error("Cannot connect to redis", zap.Error(err))
		return nil, err
	}
	return client, nil
}
