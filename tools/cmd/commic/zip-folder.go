package main

import (
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/duke-git/lancet/v2/slice"
	flag "github.com/spf13/pflag"
	"go.uber.org/atomic"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var rootPath *string = flag.StringP("path", "p", "", "The root path to check")

// 在目录下查找文件夹，并压缩成zip包
func main() {
	flag.Parse()
	if len(*rootPath) == 0 {
		panic("you should specify the path to run")
	}
	log.Println("the root path: ", *rootPath)
	var folder = *rootPath

	var err error
	absPath, err := filepath.Abs(folder)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("the absolute root path: ", absPath)
	dirs, err := os.ReadDir(absPath)
	if err != nil {
		log.Fatal(err)
	}
	runtime.GOMAXPROCS(3)

	dirs = slice.Filter(dirs, func(index int, item os.DirEntry) bool {
		return item.IsDir()
	})

	var amount = atomic.NewInt32(int32(len(dirs)))
	for _, d := range dirs {
		asyncHandle(d, folder, amount)
	}
}

func asyncHandle(d os.DirEntry, folder string, amount *atomic.Int32) {
	defer func() {
		var leftCount = amount.Dec()
		println(leftCount, " lefts to deal")
	}()

	if d.IsDir() {
		log.Println("Process:", d.Name())
		var dirPath = filepath.Join(folder, d.Name())
		var zipPath = filepath.Clean(dirPath) + ".zip"
		if err := fileutil.Zip(dirPath, zipPath); err != nil {
			log.Println(d.Name(), "error occurs:", err.Error())
			return
		}
		if err := os.RemoveAll(dirPath); err != nil {
			log.Println(d.Name(), "error occurs:", err.Error())
			return
		}
	} else {
		log.Println(d.Name(), "ignored")
	}

}
