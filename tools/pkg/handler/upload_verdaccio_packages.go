package handler

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	list "github.com/duke-git/lancet/v2/datastructure/list"
	"github.com/duke-git/lancet/v2/fileutil"
	"go.uber.org/zap"

	"move-repository/pkg/common"
)

type VerdaccioUploader struct {
	*DefaultUploader
}

func NewVerdaccioUploader(ctx context.Context) *VerdaccioUploader {
	return &VerdaccioUploader{DefaultUploader: NewUploader(ctx)}
}

func (n *VerdaccioUploader) Upload() {
	go n.LoadJsonFiles()
	for i := 0; i < n.Config.General.UploadTasks; i++ {
		go n.UploadAssets(n.GetPackagePath)
	}
}

func (n *VerdaccioUploader) GetPackagePath(config *common.Config, metaData *common.AssetMetaData) string {
	return filepath.Join(config.General.AssetsDirectory, config.Nexus.Repository, genFileName(metaData.Path))
}

func (n *VerdaccioUploader) LoadJsonFiles() {
	defer close(n.AssetChan)

	storagePath := n.Config.Verdaccio.Storage

	files, err := filterJsonFiles(storagePath)
	if err != nil {
		n.Logger.Error("failed to load 1.json file", zap.String("director", storagePath),
			zap.Error(err))
		os.Exit(0)
	}
	if len(files) == 0 {
		n.Logger.Warn("no .json files found in this directory", zap.String("directory", storagePath))
		return
	}
	for _, jsonPath := range files {
		func(jp *common.PackageJsonInfo) {

			dir := filepath.Dir(jp.FilePath)

			// 获取同目录下package对应的压缩包
			tgzFiles, err := filepath.Glob(filepath.Join(dir, "*.tgz"))
			if err != nil {
				n.Logger.Error("failed to list .tgz files", zap.String("director", dir),
					zap.Error(err))
				return
			}

			jsonString, err := fileutil.ReadFileToString(jp.FilePath)
			if err != nil {
				n.Logger.Warn("failed to load json file", zap.String("file", jp.FilePath), zap.Error(err))
				return
			}

			for _, file := range tgzFiles {
				var uriPath string

				// 解析package.json中内容，获取tgz文件对于的URL path
				// 例如： 匹配查找字符串 "tarball": "https://registry.npmjs.org/@tsconfig/node10/-/node10-1.0.0.tgz",
				// 其中 @tsconfig/node10/-/node10-1.0.0.tgz即为需要的URI path
				regString := fmt.Sprintf("\"tarball\"\\s*:\\s*\"https?://.+?/(.*?/%v)\"", filepath.Base(file))
				var pathReg = regexp.MustCompile(regString)
				matchStrings := pathReg.FindStringSubmatch(jsonString)
				if len(matchStrings) == 3 {
					uriPath = matchStrings[2]
				}

				if len(uriPath) > 0 {
					metaData := common.AssetMetaData{
						Name:             strings.Split(filepath.Base(file), "-")[0], // node10-1.0.0.tgz => node10
						Path:             uriPath,
						AbsoluteFilePath: file,
					}
					n.AssetChan <- metaData
				}

			}

		}(jsonPath)
	}
}

// 每一个需要上传的压缩包都会有一个package.json
func filterJsonFiles(basePath string) ([]*common.PackageJsonInfo, error) {
	var jsonFiles list.List[*common.PackageJsonInfo]

	realPath, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}

	// golang filepath.Glob()无法在多级子目录下递归过滤文件，故只有通过walk方法遍历
	err = filepath.Walk(realPath, func(path string, info fs.FileInfo, err error) error {
		if info.Name() == "package.json" {
			jsonFiles.Push(&common.PackageJsonInfo{FileInfo: info, FilePath: path})
		}
		return nil
	})

	return jsonFiles.Data(), err
}
