package read_ini

import (
	"gopkg.in/ini.v1"
	"log"
	"os"
)

func Generate() {
	obj := &Config{
		Name: "testConfig",

		Setting: IniSetting{
			From:             []string{"/home/jujucom/Desktop/backup2/folder/"},
			To:               "/home/jujucom/Desktop/backup2/dest/",
			FileExtension:    []string{".mp4", ".mkv", ".avi", ".rmvb"},
			FileMinSize:      "100KB",
			CheckPicture:     "true",
			PictureExtension: []string{".png", ".jpg"},

			PicMinSize:          "50KB",
			CreateRootDirectory: true,
		},
	}
	config := ini.Empty()
	err := ini.ReflectFrom(config, &obj)

	if err != nil {
		log.Fatalln(err)
	}

	//将对象转换输出到控制台
	config.WriteToIndent(os.Stdout, " ")

	//保存到目录下
	//config.SaveToIndent("./file.tmp", " ")

}
