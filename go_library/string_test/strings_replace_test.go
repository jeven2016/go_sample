package string_test_test

import (
	"regexp"
	"testing"
)

func TestStringReplacing(t *testing.T) {
	//desc := "[Thz.la]gtj-060.mp4"
	//reg := regexp.MustCompile("(\\[.*?])|(.*?原版首发_)|(.*?@)")
	//println(reg.ReplaceAllString(desc, ""))
	//println(reg.ReplaceAllString("xxfhd.com_原版首发_FSDSS-015.mp4", ""))
	//println(reg.ReplaceAllString("xxfhd.com_原版首发@FSDSS-016.mp4", ""))
	//
	//path := "d:\\sdf\\sdf\\dfgk\\xx.mp4"
	//slash := filepath.FromSlash(path)
	//println(filepath.Dir(slash))
	//
	//value := "this is a test in this country"
	//trimRight := strings.TrimRight(value, "this")
	//println(trimRight)

	var url = "http://192.168.159.129:8080/realms/master"
	reg2 := regexp.MustCompile("(https?://.*?)/.*")
	submatch := reg2.FindStringSubmatch(url)
	for _, s := range submatch {
		println(s)
	}
}
