package common

type Config struct {
	LogLevel      string `mapstructure:"log_level"`
	LogPath       string `mapstructure:"log_path"`
	OutputConsole bool   `mapstructure:"output_console"`
	Dev           bool   `mapstructure:"dev"`
	ServiceName   string `mapstructure:"service_name"`
	Port          int16  `mapstructure:"port"`
}
