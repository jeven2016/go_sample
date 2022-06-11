package initialization

import (
	"api/common"
	"api/global"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"time"
)

func SetupMongodb(config *common.Config) {
	// 设置MongoDB连接地址
	mongodbConfig := config.MongoConfig

	global.Log.Info("The mongo url", zap.String("url", mongodbConfig.Uri))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 连接MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbConfig.Uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		// 延迟释放连接
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// 检测MongoDB是否连接成功
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	global.Log.Info("Connecting mongoDB successfully")

	global.Db = client.Database(mongodbConfig.Database)

	//// 获取numbers集合
	//collection := client.Database("testing").Collection("numbers")
	//
	//// 插入一个文档
	//res, err := collection.InsertOne(ctx, bson.D{{"name", "pi"}, {"value", 3.14159}})
	//id := res.InsertedID
	//fmt.Println("新增文档Id=", id)
}
