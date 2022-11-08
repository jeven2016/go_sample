package main

import (
	"context"
	"runtime"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/slice"
	flag "github.com/spf13/pflag"
	"go.uber.org/zap"

	"move-repository/pkg/common"
	"move-repository/pkg/handler"
)

var configPath *string = flag.StringP("config", "c", "/home/jujucom/Desktop/workspace/projects/go_samples/tools/conf/conf.ini", "The path of config file")
var command *string = flag.StringP("command", "m", "upload", "The supported command: download, upload ")
var repoType *string = flag.StringP("source-repository-type", "t", "nexus", "The type of repository that has packages stored and need to upload into JFrog, it can only be: nexus or verdaccio")

// supported repositories
var repos = []string{"nexus", "verdaccio"}

func main() {
	flag.Parse()
	config, _ := common.SetupViper(*configPath)

	if len(*command) == 0 {
		panic("you should specify the command to run: download or upload")
	}

	validateRepoType(repoType)

	// log初始化
	logger := common.SetupLog(*config)
	defer logger.Sync()

	if json, err := convertor.ToJson(config); err == nil {
		logger.Info("the configuration parsed", zap.String("content", json))
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, "config", config)
	ctx = context.WithValue(ctx, "logger", logger)

	runtime.GOMAXPROCS(5)

	switch *command {
	case "download":
		handler.Download(ctx)
	case "upload":
		if *repoType == "verdaccio" {
			handler.NewVerdaccioUploader(ctx).Upload()
		} else {
			handler.NewNexusUploader(ctx).Upload()
		}
	}
	select {}
}

func validateRepoType(repoType *string) {
	if !slice.Contain(repos, *repoType) {
		panic("unsupported type of repository, it should be : nexus or verdaccio")
	}
}

func show(h *handler.Uploader) {

}
