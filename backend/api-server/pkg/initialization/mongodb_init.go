package initialization

import (
	"api/pkg/common"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"time"
)

// SetupMongodb 初始化Mongodb
func SetupMongodb(config *common.Config, log *zap.Logger) (*mongo.Client, *mongo.Database, error) {
	// 设置MongoDB连接地址
	mongodbConfig := config.MongoConfig

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 连接MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbConfig.Uri))
	if err != nil {
		log.Error("Cannot connect to mongodb", zap.Error(err))
		return client, nil, err
	}

	// 检测MongoDB是否连接成功
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Error("Failed to Ping for mongodb", zap.Error(err))
		return client, nil, err
	}

	//初始化全局Db
	db := client.Database(mongodbConfig.Database)
	return client, db, nil
}
