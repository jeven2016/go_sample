package common

type Config struct {
	General General `mapstructure:"general,omitempty"`
	Nexus   Nexus   `mapstructure:"nexus,omitempty"`
	Jfrog   Jfrog   `mapstructure:"jfrog,omitempty"`
}

type General struct {
	LogLevel      string `mapstructure:"log_level,omitempty"`
	LogPath       string `mapstructure:"log_path,omitempty"`
	OutputConsole bool   `mapstructure:"output_console,omitempty"`
	ServiceName   string `mapstructure:"service_name,omitempty"`
	QueueSize     int    `mapstructure:"internal_queue_size,omitempty"`
}

type Nexus struct {
	BaseUrl         string `mapstructure:"baseUrl,omitempty"`
	Repository      string `mapstructure:"repository,omitempty"`
	DownloadAssets  bool   `mapstructure:"download_assets,omitempty"`
	AssetsDirectory string `mapstructure:"assets_directory,omitempty"`
	Username        string `mapstructure:"username,omitempty"`
	Password        string `mapstructure:"password,omitempty"`
}

type Jfrog struct {
}
