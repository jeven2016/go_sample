package download

import (
	"context"
	"encoding/json"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"nj-mover/pkg/common"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const urlSeparator = "/"

var once = sync.Once{}

var client *resty.Client

func Download(ctx context.Context) {
	config := ctx.Value("config").(*common.Config)
	var assetsChan = make(chan common.Item, config.General.QueueSize)
	go parsePages(ctx, assetsChan, config)
	go downloadAssets(ctx, assetsChan, config)
}

func parsePages(ctx context.Context, itemChan chan<- common.Item, config *common.Config) {
	logger := ctx.Value("logger").(*zap.Logger)
	nexusCfg := config.Nexus

	var assetCount, pages = 0, 0
	var continuationToken string //the token to fetch the next list of assets

	defer func() {
		close(itemChan)
		logger.Info("Finished for parsing all pages for this repository",
			zap.String("repository", nexusCfg.Repository),
			zap.Int("assetCount", assetCount),
			zap.Int("pageCount", pages))
	}()

	for {
		//if pages > 0 && len(continuationToken) == 0 {
		if pages >= 1 {
			break
		}
		url := makeUrl(nexusCfg, continuationToken)

		res, err := restyClient().R().SetBasicAuth(nexusCfg.Username, nexusCfg.Password).Get(url)
		if err != nil {
			logger.Warn("failed to get list of components",
				zap.String("repository", nexusCfg.Repository),
				zap.Error(err))
		}

		var assets = &common.Assets{}
		err = json.Unmarshal(res.Body(), assets)
		if err != nil {
			logger.Warn("failed to convert json string to assets data",
				zap.String("repository", nexusCfg.Repository),
				zap.Error(err))
		}

		// add items into chanel
		if items := assets.Items; items != nil {
			for _, it := range *items {
				itemChan <- it
				assetCount++
			}
		}

		pages++
		continuationToken = assets.ContinuationToken
		logger.Info("Completed for parsing this pages",
			zap.String("repository", nexusCfg.Repository),
			zap.Int("page", pages),
			zap.Int("currentAssetCount", assetCount))
	}

}

func makeUrl(nexusCfg common.Nexus, continuationToken string) string {
	if len(nexusCfg.BaseUrl) == 0 {
		panic("the baseUrl of nexus is mandatory, you should define it in config file")
	}

	//http://localhost:8081/service/rest/v1/components?repository=npm-proxy
	var url = strings.TrimRight(nexusCfg.BaseUrl, urlSeparator) + "/service/rest/v1/components?repository=" + nexusCfg.Repository
	if continuationToken != "" {
		url += "&continuationToken=" + continuationToken
	}
	return url
}

func restyClient() *resty.Client {
	once.Do(func() {
		client = resty.New().SetTimeout(20 * time.Second)
	})
	return client
}

func downloadAssets(ctx context.Context, itemChan <-chan common.Item, config *common.Config) {
	logger := ctx.Value("logger").(*zap.Logger)

	logger.Info("waiting for downloading assets", zap.String("repository", config.Nexus.Repository))
	for item := range itemChan {
		directory, err := ensureDirectory(config, item.Repository, logger)
		if err != nil {
			return
		}
		for _, asset := range *item.Assets {
			fileName := path.Base(asset.Path)
			_, err := restyClient().R().SetOutput(filepath.Join(directory, fileName)).Get(asset.DownloadUrl)
			if err != nil {
				logger.Error("failed to download asset",
					zap.String("repository", config.Nexus.Repository),
					zap.String("url", asset.DownloadUrl))
				continue
			}

			logger.Info("asset downloaded", zap.String("repository", config.Nexus.Repository),
				zap.String("path", asset.Path))
		}

	}
}

func ensureDirectory(config *common.Config, repository string, logger *zap.Logger) (string, error) {
	path := filepath.Join(config.Nexus.AssetsDirectory, repository)
	exists := fileutil.IsExist(path)
	if !exists {
		if err := fileutil.CreateDir(path); err != nil {
			logger.Error("failed to create directory",
				zap.String("directory", path),
				zap.Error(err))
			return "", err
		}
	}
	return path, nil
}
