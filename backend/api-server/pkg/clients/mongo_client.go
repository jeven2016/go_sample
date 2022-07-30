package clients

import (
	"api/pkg/common"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"time"
)

type MongoClient struct {
	Client *mongo.Client
	Db     *mongo.Database
	Log    *zap.Logger
	Config *common.MongoConfig
}

func (c *MongoClient) StartInit() error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 连接MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(c.Config.Uri))
	if err != nil {
		c.Log.Error("Cannot connect to mongodb", zap.Error(err))
		return err
	}

	// 检测MongoDB是否连接成功
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		c.Log.Error("Failed to Ping for mongodb", zap.Error(err))
		return err
	}

	//初始化全局Db
	db := client.Database(c.Config.Database)

	c.Db = db
	c.Client = client
	return err
}
