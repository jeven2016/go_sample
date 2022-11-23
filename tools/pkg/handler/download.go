package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"move-repository/pkg/common"
)

const urlSeparator = "/"

var once = sync.Once{}

var client *resty.Client

// Download 从nexus上下载仓库下的包
func Download(ctx context.Context) {
	config := ctx.Value("config").(*common.Config)
	logger := ctx.Value("logger").(*zap.Logger)

	baseUrl, repository := getBaseUrlAndRepository(config.Nexus)
	ctx = context.WithValue(ctx, "repository", repository)
	ctx = context.WithValue(ctx, "baseUrl", baseUrl)

	var assetsChan = make(chan common.Item, config.General.QueueSize)

	directory, err := ensureDirectory(config, logger, ctx)
	if err != nil {
		return
	}

	go parsePages(assetsChan, config, logger, ctx)

	for i := 0; i < config.General.UploadTasks; i++ {
		go downloadAssets(directory, assetsChan, config, logger, ctx)
	}
}

func parsePages(itemChan chan<- common.Item, config *common.Config, logger *zap.Logger, ctx context.Context) {
	// var nexusCfg = config.Nexus
	logger.Info("Parsing the pages......")
	var assetCount, pages = 0, 0
	var continuationToken string // the token to fetch the next list of assets

	defer func() {
		close(itemChan)
		logger.Info("Finished for parsing all pages for this repository",
			zap.Int("assetsToDownload", assetCount),
			zap.Int("pageCount", pages))
	}()

	repository := getRepositoryName(ctx)
	baseUrl := ctx.Value("baseUrl").(string)

	for {
		if pages > 0 && len(continuationToken) == 0 {
			break
		}
		url := makeUrl(baseUrl, repository, continuationToken)

		res, err := restyClient(config.General.UploadTimeout).
			R().
			// SetBasicAuth(nexusCfg.Username, nexusCfg.Password).
			Get(url)
		if err != nil {
			logger.Warn("failed to get list of components",
				zap.String("repository", repository),
				zap.Error(err))
		}

		var assets = &common.Assets{}
		var resp = res.Body()
		err = json.Unmarshal(resp, assets)
		if err != nil {
			logger.Warn("failed to convert json string to assets data",
				zap.String("repository", repository),
				zap.String("response", convertor.ToString(resp)),
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

		if len(assets.ContinuationToken) == 0 {
			logger.Warn("assets.ContinuationToken is blank and that means no next page need to check", zap.Int("page", pages))
		}

		logger.Info("Completed for parsing this pages",
			zap.String("repository", repository),
			zap.Int("page", pages),
			zap.Int("currentAssetCount", assetCount))
	}

}

func getBaseUrlAndRepository(nexusCfg common.Nexus) (string, string) {
	if len(nexusCfg.RepositoryUrl) == 0 {
		panic("the baseUrl of nexus is required, you should define it in config file")
	}
	result, err := url.Parse(nexusCfg.RepositoryUrl)
	if err != nil {
		panic("Invalid repository url provided: " + nexusCfg.RepositoryUrl)
	}

	baseUrl := fmt.Sprintf("%s://%s", result.Scheme, result.Host)

	// get the repository name from the url
	var pureUrl = strings.TrimRight(result.Path, urlSeparator)
	var repository = pureUrl[strings.LastIndex(pureUrl, "/")+1:]

	return baseUrl, repository
}

func makeUrl(baseUrl string, repository string, continuationToken string) string {
	// http://localhost:8081/service/rest/v1/components?repository=npm-proxy
	var url = baseUrl + "/service/rest/v1/components?repository=" + repository
	if continuationToken != "" {
		url += "&continuationToken=" + continuationToken
	}
	return url
}

func restyClient(timeout int) *resty.Client {
	once.Do(func() {
		client = resty.New().SetTimeout(time.Duration(timeout) * time.Second)
	})
	return client
}

func downloadAssets(directory string, itemChan <-chan common.Item, config *common.Config, logger *zap.Logger, ctx context.Context) {
	logger.Info("waiting for downloading assets", zap.String("repository", getRepositoryName(ctx)))
	for item := range itemChan {
		for _, asset := range *item.Assets {
			downloadAsset(asset, directory, logger, config, ctx)
		}
	}
}

func genFileName(filePath string) string {
	fileName := assetName(filePath)

	// check if group prefix exists
	if strings.HasPrefix(filePath, "@") {
		group := strings.Split(filePath, "/")[0]
		fileName = group + "___" + fileName
	}

	return fileName
}

func downloadAsset(asset common.Asset, directory string, logger *zap.Logger, config *common.Config, ctx context.Context) {
	fileName := genFileName(asset.Path)
	repo := getRepositoryName(ctx)

	if fileutil.IsExist(filepath.Join(directory, fileName)) {
		logger.Info("Ignored(asset exists)", zap.String("file", fileName))
		return
	}

	// save metadata
	err := writeMedata(asset, directory, fileName, assetName(asset.Path))

	if err != nil {
		logger.Error("failed to write meta data",
			zap.String("repository", repo),
			zap.String("url", asset.Npm.Name),
			zap.Error(err))
		return
	}

	// save asset
	_, err = restyClient(config.General.UploadTimeout).R().SetOutput(filepath.Join(directory, fileName)).Get(asset.DownloadUrl)
	if err != nil {
		logger.Error("failed to handler ast",
			zap.String("repository", repo),
			zap.String("url", asset.DownloadUrl))
		return
	}

	logger.Info("ast downloaded", zap.String("repository", repo),
		zap.String("path", asset.Path))
}

func writeMedata(asset common.Asset, directory string, fileName string, pureName string) error {
	metaData, _ := convertor.ToJson(common.AssetMetaData{
		Name: pureName,
		Path: asset.Path,
	})

	filePath := filepath.Join(directory, fileName) + ".json"
	err := ioutil.WriteFile(filePath, []byte(metaData), 0664)
	return err
}

func ensureDirectory(config *common.Config, logger *zap.Logger, ctx context.Context) (string, error) {
	dir, exists := directoryExists(config, getRepositoryName(ctx))
	if !exists {
		abs, err := filepath.Abs(dir)
		if err != nil {
			logger.Error("invalid path defined in config file",
				zap.String("assets_directory", dir),
				zap.Error(err))
			return "", err
		}

		if err := os.MkdirAll(abs, os.ModePerm); err != nil {
			logger.Error("failed to create directory",
				zap.String("directory", dir),
				zap.Error(err))
			return "", err
		}

		logger.Info("Directory is created", zap.String("directory", dir))
	}
	return dir, nil
}

func getRepositoryName(ctx context.Context) string {
	return ctx.Value("repository").(string)
}

func directoryExists(config *common.Config, repository string) (string, bool) {
	dir := filepath.Join(config.Nexus.AssetsDirectory, repository)
	exists := fileutil.IsExist(dir)
	return dir, exists
}

func assetName(pathString string) string {
	return path.Base(pathString)
}
