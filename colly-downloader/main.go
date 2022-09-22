package main

import (
	"runtime"

	"colly-downloader/pkg/yzs8"
)

func main() {

	runtime.GOMAXPROCS(4)
	// yzs8.Start()

	yzs8.ConvertCatalog()
	select {}
}
