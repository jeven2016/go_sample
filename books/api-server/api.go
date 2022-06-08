package main

import (
	"api/common"
	"fmt"
	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
)

var configPath *string = flag.StringP("config", "c", "", "The path of defined config file")

func main() {
	// 解析传入的参数, 运行时可以添加 -h查看细节
	flag.Parse()

	println("the config path is", *configPath)
	config, err := common.InitViper(*configPath)
	if err != nil {
		panic(fmt.Sprintf("Program exited, %v", err))
	}

	logger := common.InitLog(*config)
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			logger.Error("Some logs aren't flushed into log file")
		}
	}(logger)

}
