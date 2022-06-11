package initialization

import (
	"api/common"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

var config = new(common.Config)

func SetupViper(configFilePath string) (*common.Config, error) {
	if len(configFilePath) == 0 {
		viper.SetConfigName("api")
		viper.SetConfigType("ini")
		viper.AddConfigPath("./conf")
		viper.AddConfigPath("/etc/books/")
	} else {
		viper.SetConfigFile(configFilePath)
	}
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to load api.ini: %v", err)
		return nil, err
	}

	if err := viper.Unmarshal(config); err != nil {
		log.Fatalf("Failed to convert the data in config.ini: %v", err)
		return nil, err
	}

	// 监控配置文件变化
	viper.WatchConfig()
	// 注意！！！配置文件发生变化后要同步到全局变量Conf
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Config file is changed, reload it now")
		if err := viper.Unmarshal(config); err != nil {
			panic(fmt.Errorf("The config cann't be updated:%s \n", err))
		}
	})

	return config, nil
}
