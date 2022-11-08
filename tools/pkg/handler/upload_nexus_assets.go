package handler

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/duke-git/lancet/v2/fileutil"
	"go.uber.org/zap"

	"move-repository/pkg/common"
)

type NexusUploader struct {
	*DefaultUploader
}

func NewNexusUploader(ctx context.Context) *NexusUploader {
	return &NexusUploader{DefaultUploader: NewUploader(ctx)}
}

func (n *NexusUploader) Upload() {
	go n.LoadJsonFiles()
	for i := 0; i < n.Config.General.UploadTasks; i++ {
		go n.UploadAssets(n.GetPackagePath)
	}
}

func (n *NexusUploader) GetPackagePath(config *common.Config, metaData *common.AssetMetaData) string {
	return filepath.Join(config.Jfrog.AssetsDirectory, genFileName(metaData.Path))
}

func (n *NexusUploader) LoadJsonFiles() {
	defer close(n.AssetChan)

	directory := n.Config.Jfrog.AssetsDirectory
	files, err := filepath.Glob(filepath.Join(directory, "*.json"))
	if err != nil {
		n.Logger.Error("failed to load meta data(.json)", zap.String("director", directory),
			zap.Error(err))
		os.Exit(0)
	}

	if len(files) == 0 {
		n.Logger.Warn("no json metadata file found in this directory", zap.String("directory", directory))
		return
	}
	for _, jp := range files {
		jsonString, err := fileutil.ReadFileToString(jp)
		if err != nil {
			n.Logger.Warn("failed to load json file", zap.String("file", jp), zap.Error(err))
			return
		}

		var metaData common.AssetMetaData
		err = json.Unmarshal([]byte(jsonString), &metaData)
		if err != nil {
			n.Logger.Warn("failed to unmarshal json file",
				zap.String("file", jp),
				zap.String("json", jsonString),
				zap.Error(err))
			return
		}

		n.AssetChan <- metaData
	}
	n.Logger.Info("all jsons files loaded")
}
