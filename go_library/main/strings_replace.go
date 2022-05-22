package main

import (
	"path/filepath"
	"regexp"
)

func ReplaceString() {
	desc := "[Thz.la]gtj-060.mp4"
	reg := regexp.MustCompile("(\\[.*?])|(.*?原版首发_)|(.*?@)")
	println(reg.ReplaceAllString(desc, ""))
	println(reg.ReplaceAllString("xxfhd.com_原版首发_FSDSS-015.mp4", ""))
	println(reg.ReplaceAllString("xxfhd.com_原版首发@FSDSS-016.mp4", ""))

	path := "d:\\sdf\\sdf\\dfgk\\xx.mp4"
	slash := filepath.FromSlash(path)
	println(filepath.Dir(slash))
}
