package handler

import (
	"context"
	"encoding/json"
	"github.com/duke-git/lancet/v2/fileutil"
	"go.uber.org/zap"
	"io/ioutil"
	"move-repository/pkg/common"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var jfrogUrl string
var onceJfrog = sync.Once{}

func Upload(ctx context.Context) {
	config := ctx.Value("config").(*common.Config)
	logger := ctx.Value("logger").(*zap.Logger)

	var metaDataChan = make(chan common.AssetMetaData, config.General.QueueSize)

	go loadJsonFiles(config, logger, metaDataChan)

	for i := 0; i < config.General.UploadTasks; i++ {
		go uploadAssets(config, logger, metaDataChan)
	}
}

func loadJsonFiles(config *common.Config, logger *zap.Logger, mdChan chan<- common.AssetMetaData) {
	directory := config.General.AssetsDirectory
	nexusRepo := config.Nexus.Repository
	files, err := filepath.Glob(filepath.Join(directory, nexusRepo, "*.json"))
	if err != nil {
		logger.Error("failed to load meta data(.json)", zap.String("director", directory),
			zap.Error(err))
		os.Exit(0)
	}

	for _, jsonPath := range files {
		go func(jp string) {
			jsonString, err := fileutil.ReadFileToString(jp)
			if err != nil {
				logger.Warn("failed to load json file", zap.String("file", jp), zap.Error(err))
				return
			}

			var metaData common.AssetMetaData
			err = json.Unmarshal([]byte(jsonString), &metaData)
			if err != nil {
				logger.Warn("failed to unmarshal json file",
					zap.String("file", jp),
					zap.String("json", jsonString),
					zap.Error(err))
				return
			}

			mdChan <- metaData
		}(jsonPath)
	}

}

func uploadAssets(config *common.Config, logger *zap.Logger, mdChan <-chan common.AssetMetaData) {
	if len(config.Jfrog.BaseUrl) == 0 {
		panic("the baseUrl of jfrog is mandatory, you should define it in config file")
	}

	for item := range mdChan {
		var assetUrl = makeJfrogUrl(config.Jfrog) + "/" + item.Path
		sourceFilePath := filepath.Join(config.General.AssetsDirectory, config.Nexus.Repository, item.Name)
		if !fileutil.IsExist(sourceFilePath) {
			continue
		}

		fileBytes, err := ioutil.ReadFile(sourceFilePath)

		response, err := restyClient(config.General.UploadTimeout).R().
			SetBody(fileBytes).
			SetHeader("X-JFrog-Art-Api", config.Jfrog.ApiKey).
			Put(assetUrl)

		if err != nil || response.StatusCode() != http.StatusCreated {
			logger.Warn("failed to upload", zap.String("assetUrl", assetUrl),
				zap.Error(err), zap.Int("status", response.StatusCode()))
			return
		}
		logger.Info("upload successfully", zap.String("assetUrl", assetUrl))
		os.Exit(0)
	}
}

func makeJfrogUrl(jfrog common.Jfrog) string {
	onceJfrog.Do(func() {
		if len(jfrog.BaseUrl) == 0 {
			panic("the baseUrl of jfrog is mandatory, you should define it in config file")
		}
		jfrogUrl = strings.TrimRight(jfrog.BaseUrl, urlSeparator) + "/" + jfrog.Repository
	})

	return jfrogUrl
}
