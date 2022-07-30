package main

import (
	"api/pkg/clients"
	"api/pkg/global"
	"api/pkg/initialization"
	"context"
	"errors"
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
	flag "github.com/spf13/pflag"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var configPath *string = flag.StringP("config", "c", "", "The path of config file")

func main() {
	// 解析传入的参数, 运行时可以添加 -h查看细节
	flag.Parse()

	//读取配置文件
	config, _ := initialization.SetupViper(*configPath)

	//创建一个全局的App
	app := new(global.App)

	//log初始化
	logger := initialization.SetupLog(*config)
	defer logger.Sync()

	//设置logger
	app.Log = logger

	if json, err := convertor.ToJson(config); err == nil {
		app.Log.Info("the configuration parsed", zap.String("content", json))
	}

	//初始化redis
	redisClient := &clients.RedisClient{Log: logger, Config: &config.RedisConfig}
	err := redisClient.StartInit()
	if err == nil {
		app.Log.Info("Connecting to redis successfully")
		app.RedisClient = redisClient
	}
	defer closeRedis(redisClient, app)

	//初始化Mongodb
	mongoClient := &clients.MongoClient{Log: logger, Config: &config.MongoConfig}
	err = mongoClient.StartInit()
	if err == nil {
		app.Log.Info("Connecting to mongodb successfully")
		app.MongoClient = mongoClient
	}
	defer closeMongodb(app, err, mongoClient.Client)

	//校验错误提示国际化
	if err := initialization.InitTrans("zh"); err != nil {
		app.Log.Warn("Failed to initialize i18 resources")
	}

	//初始化OIDC client
	authClient := &clients.AuthClient{Log: logger, Config: &config.AuthConfig}
	err = authClient.StartInit()
	if err == nil {
		app.Log.Info("Initializing client for keycloak successfully")
		app.MongoClient = mongoClient
	}
	app.AuthClient = authClient

	//web初始化
	engine := initialization.SetupEngine(config, app)

	bind := fmt.Sprintf("%v:%v", config.ApiServerConfig.BindAddress, config.ApiServerConfig.Port)
	if err := engine.Run(bind); err != nil {
		panic(errors.New(fmt.Sprintf("Server fails to start: %v", err)))
	}
}

func closeMongodb(app *global.App, err error, client *mongo.Client) {
	// 延迟释放连接
	app.Log.Info("closing mongodb connection")
	if err = client.Disconnect(context.TODO()); err != nil {
		app.Log.Error("Failed to disconnect mongodb", zap.Error(err))
	} else {
		app.Log.Info("Mongodb disconnected now")
	}
}

func closeRedis(redisClient *clients.RedisClient, app *global.App) {
	if redisClient.Client == nil {
		return
	}
	app.Log.Info("closing mongodb connection")
	if err := redisClient.Client.Close(); err != nil {
		app.Log.Error("Failed to disconnect redis", zap.Error(err))
	} else {
		app.Log.Info("Redis disconnected now")
	}
}
