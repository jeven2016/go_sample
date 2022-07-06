package main

import (
	"api/pkg/common"
	initialization2 "api/pkg/initialization"
	"context"
	"errors"
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
)

var configPath *string = flag.StringP("config", "c", "", "The path of defined config file")

func main() {
	// 解析传入的参数, 运行时可以添加 -h查看细节
	flag.Parse()

	//读取配置文件
	config, _ := initialization2.SetupViper(*configPath)

	//创建一个全局的App
	app := new(common.App)

	//log初始化
	logger := initialization2.SetupLog(*config)
	defer logger.Sync()

	//设置logger
	app.Log = logger

	if json, err := convertor.ToJson(config); err == nil {
		app.Log.Error("the configuration parsed", zap.String("content", json))
	}

	//初始化redis
	redisClient, err := initialization2.SetupRedis(config, app.Log)
	if err == nil {
		app.Log.Info("Connecting to redis successfully", zap.Error(err))
		app.RedisClient = redisClient
	}
	defer func() {
		if redisClient == nil {
			return
		}
		if err := redisClient.Close(); err != nil {
			app.Log.Error("Cannot disconnect redis", zap.Error(err))
		}
	}()

	//初始化Mongodb
	client, db, err := initialization2.SetupMongodb(config, app.Log)
	if err == nil {
		app.Log.Info("Connecting to mongoDB successfully")
		defer func() {
			// 延迟释放连接
			if err = client.Disconnect(context.TODO()); err != nil {
				app.Log.Error("The mongodb client cannot disconnect.", zap.Error(err))
			}
		}()
		//设置DB
		app.Db = db
	} else {
		app.Log.Error("Cannot initialize mongo connection", zap.Error(err))
	}

	//校验错误提示国际化
	if err := initialization2.InitTrans("zh"); err != nil {
		app.Log.Warn("Failed to initialize i18 resources")
	}

	//web初始化
	engine := initialization2.SetupEngine(config, app)

	bind := fmt.Sprintf("%v:%v", config.ApiServerConfig.BindAddress, config.ApiServerConfig.Port)
	if err := engine.Run(bind); err != nil {
		panic(errors.New(fmt.Sprintf("Server fails to start: %v", err)))
	}
}
