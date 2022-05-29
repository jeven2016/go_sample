package read_ini

import (
	"github.com/duke-git/lancet/v2/convertor"
	ini "gopkg.in/ini.v1"
	"log"
)

type IniSetting struct {
	From                []string `ini:"from"`
	To                  string   `ini:"to"`
	FileExtension       []string `ini:"file_extension,omitempty"`
	FileMinSize         string   `ini:"file_min_size"`
	CheckPicture        string   `ini:"check_picture"`
	PictureExtension    []string `ini:"picture_extension,omitempty"`
	PicMinSize          string   `ini:"pic_min_size"`
	CreateRootDirectory bool     `ini:"create_root_directory"`
}

type Config struct {
	Name string `ini:"name"`

	Setting IniSetting `ini:"setting"`
}

func ReadIniFile() {
	//第一种方式，先load文件再转换
	data, err := ini.Load("./config.ini")
	if err != nil {
		log.Fatalln(err)
	}

	config := Config{}
	err = data.MapTo(&config)
	if err != nil {
		log.Fatalln(err)
	}
	println(convertor.ToJson(config))

	//直接从源文件转换
	err = ini.MapTo(&config, "./config.ini")
	if err != nil {
		log.Fatalln(err)
	}
	println(convertor.ToJson(config))

}
