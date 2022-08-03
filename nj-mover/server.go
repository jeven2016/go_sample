package main

import (
	"context"
	"github.com/duke-git/lancet/v2/convertor"
	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
	"nj-mover/pkg/common"
	"nj-mover/pkg/download"
	"runtime"
)

var configPath *string = flag.StringP("config", "c", "", "The path of config file")

func main() {
	flag.Parse()
	config, _ := common.SetupViper(*configPath)

	//log初始化
	logger := common.SetupLog(*config)
	defer logger.Sync()

	if json, err := convertor.ToJson(config); err == nil {
		logger.Info("the configuration parsed", zap.String("content", json))
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, "config", config)
	ctx = context.WithValue(ctx, "logger", logger)

	runtime.GOMAXPROCS(5)
	download.Download(ctx)

	select {}
}
