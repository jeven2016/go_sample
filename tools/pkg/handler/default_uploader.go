package handler

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/duke-git/lancet/v2/fileutil"
	"go.uber.org/zap"

	"move-repository/pkg/common"
)

var jfrogUrl string
var onceJfrog = sync.Once{}

// Uploader : upload interface
type Uploader interface {
	Upload()
}

// GetPackagePath : A callback for retrieving the real absolute file path of a specific package
type GetPackagePath func(config *common.Config, metaData *common.AssetMetaData) string

type DefaultUploader struct {
	Config    *common.Config
	Logger    *zap.Logger
	AssetChan chan common.AssetMetaData
}

func NewUploader(ctx context.Context) *DefaultUploader {
	config := ctx.Value("config").(*common.Config)
	logger := ctx.Value("logger").(*zap.Logger)

	return &DefaultUploader{
		config, logger, make(chan common.AssetMetaData, config.General.QueueSize),
	}
}

func (d *DefaultUploader) UploadAssets(getPackagePath GetPackagePath) {
	if len(d.Config.Jfrog.BaseUrl) == 0 {
		panic("the baseUrl of jfrog is mandatory, you should define it in config file")
	}

	for item := range d.AssetChan {
		var assetUrl = makeJfrogUrl(d.Config.Jfrog) + "/" + item.Path
		sourceFilePath := getPackagePath(d.Config, &item)
		if !fileutil.IsExist(sourceFilePath) {
			continue
		}

		fileBytes, err := ioutil.ReadFile(sourceFilePath)

		response, err := restyClient(d.Config.General.UploadTimeout).R().
			SetBody(fileBytes).
			SetHeader("X-JFrog-Art-Api", d.Config.Jfrog.ApiKey).
			Put(assetUrl)

		if err != nil || response.StatusCode() != http.StatusCreated {
			d.Logger.Warn("failed to upload", zap.String("assetUrl", assetUrl),
				zap.Error(err), zap.Int("status", response.StatusCode()))
			return
		}
		d.Logger.Info("upload successfully", zap.String("assetUrl", assetUrl))
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
