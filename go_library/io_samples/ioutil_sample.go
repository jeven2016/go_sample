package io_samples

import (
	"io/fs"
	"io/ioutil"
	"log"
	"os"
)

func Ioutil_sample() {
	//返回目录下的文件列表
	fileInfos, err := ioutil.ReadDir(".")
	if err != nil {
		log.Println("error occurs：", err)
	}

	for _, f := range fileInfos {
		println(f.Name())
	}

	//一次性读取文件内容
	file, err := ioutil.ReadFile("./go.mod")
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	println("content is :", string(file))

	//创建临时目录
	dir, err := ioutil.TempDir(".", "tempDir")
	defer os.RemoveAll(dir)
	if err != nil {
		log.Fatal(err)
	}
	println("temp dir created:", dir)

	//创建临时file
	tempFile, err := ioutil.TempFile(".", "tempFile")
	defer os.Remove(tempFile.Name())
	if err != nil {
		log.Fatalln(err)
	}
	println("temp fie created:", tempFile.Name())

	//写入文件
	err = ioutil.WriteFile("tempFile", []byte("hello,nanjing"), fs.FileMode(fs.ModeAppend|fs.ModePerm))
	if err != nil {
		println("the final error", err)
	}
}
