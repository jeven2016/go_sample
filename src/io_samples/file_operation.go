package io_samples

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

func FileAction() {

	//新建文件
	file, _ := os.Create("./tempCreated.md")
	defer file.Close()

	file.WriteString("hello\n")
	file.Write([]byte("Nanjing"))
	file.WriteAt([]byte("what"), 20)
	file.Sync() //flush

	//一次性读文件, 与ioutil.ReadFile等价
	bytes, err := os.ReadFile("./go.mod")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("file content=", string(bytes))

	//缓冲读取文件，以byte数组方式读入
	cache := make([]byte, 2048)
	var buf []byte
	openFile, _ := os.Open("./go.mod")
	defer openFile.Close()
	for {
		n, err := openFile.Read(cache)
		if err == io.EOF {
			//如果读到文件末尾就退出
			break
		}
		if err != nil {
			log.Fatalln("cannot read file:", err)
		}

		//读出的数据添加到切片中去
		buf = append(buf, cache[:n]...)
	}
	println("Read file through cache: ", string(buf))

	//通过读取每一行的方式读取文件,读出的内容不包括换行符
	goFile, _ := os.Open("./go.mod")
	defer goFile.Close()
	reader := bufio.NewReader(goFile)
	var strBuf strings.Builder
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			//如果读到文件末尾就退出
			break
		}
		if err != nil {
			log.Fatalln("cannot read file through NewReader:", err)
		}
		strBuf.WriteString(string(line))
		strBuf.WriteString("\n")
	}
	println("Read file through form reader: ", strBuf.String())

	//通过Writer去写入文件内容,覆盖
	file2, _ := os.Create("./tempCreated.md")

	writer := bufio.NewWriter(file2)
	writer.WriteString(`
		hello~~~ Nanjing2
    `)
	writer.Flush()
	file2.Close()

	//移动文件
	err = os.Rename("./tempCreated.md", "./src/tempCreated.md")
	if err != nil {
		log.Fatalln(err)
	}

	//拷贝文件，使用io. 也可以使用ioutil, Read & Write实现
	src, _ := os.Open("./go.mod")
	//添加os.O_APPEND后，会反复追加文件内容
	dst, _ := os.OpenFile("./tempCopied.md", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer src.Close()
	defer dst.Close()
	_, error := io.Copy(dst, src)
	if error != nil {
		log.Fatalln(error)

	}

}
