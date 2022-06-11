package main

import (
	"api/common"
	"api/initialization"
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
	config, err := initialization.SetupViper(*configPath)
	if err != nil {
		panic(fmt.Sprintf("Program exited, %v", err))
	}

	//log初始化
	logger := common.SetupLog(*config)
	defer logger.Sync()

	if json, err := convertor.ToJson(config); err == nil {
		logger.Info("the configuration parsed", zap.String("content", json))
	}

	initialization.SetupMongodb(config)

	//web初始化
	engine := initialization.SetupEngine(config)

	bind := fmt.Sprintf("%v:%v", config.ApiServerConfig.BindAddress, config.ApiServerConfig.Port)
	if err := engine.Run(bind); err != nil {
		panic(errors.New(fmt.Sprintf("failed to start: %v", err)))
	}
}
