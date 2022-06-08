package api

import (
	"github.com/spf13/viper"
	"log"
)

func InitViper() error {
	viper.SetConfigName("api")
	viper.SetConfigType("ini")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/books/")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to load api.ini: %v", err)
		return err
	}

	var config Config

}
