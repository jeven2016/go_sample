package common

import "io/fs"

type Config struct {
	General   General   `mapstructure:"general,omitempty"`
	Nexus     Nexus     `mapstructure:"nexus,omitempty"`
	Jfrog     Jfrog     `mapstructure:"jfrog,omitempty"`
	Verdaccio Verdaccio `mapstructure:"verdaccio,omitempty"`
}

type General struct {
	AssetsDirectory string `mapstructure:"assets_directory,omitempty"`
	LogLevel        string `mapstructure:"log_level,omitempty"`
	LogPath         string `mapstructure:"log_path,omitempty"`
	OutputConsole   bool   `mapstructure:"output_console,omitempty"`
	QueueSize       int    `mapstructure:"internal_queue_size,omitempty"`
	UploadTimeout   int    `mapstructure:"upload_timeout_seconds,omitempty"`
	UploadTasks     int    `mapstructure:"number_of_upload_tasks,omitempty"`
}

type Nexus struct {
	BaseUrl        string `mapstructure:"base_url,omitempty"`
	Repository     string `mapstructure:"repository,omitempty"`
	DownloadAssets bool   `mapstructure:"download_assets,omitempty"`
	Username       string `mapstructure:"username,omitempty"`
	Password       string `mapstructure:"password,omitempty"`
}

type Jfrog struct {
	BaseUrl    string `mapstructure:"base_url,omitempty"`
	ApiKey     string `mapstructure:"api_key,omitempty"`
	Username   string `mapstructure:"username,omitempty"`
	Repository string `mapstructure:"repository,omitempty"`
}

type Verdaccio struct {
	Storage string `mapstructure:"storage,omitempty"`
}

type PackageJsonInfo struct {
	FileInfo fs.FileInfo
	FilePath string
}
