package jeven

import "strings"

func HelloJeven() {
	println("hello jeven")
}

func init() {
	println("jeven_main.go_sample init()")
}

func OperateString() {
	str := "hello中国"
	println("len=", len(str))

	for i, value := range str {
		println("index=", i, "\tvalue=", string(value))
	}

	indexVal := strings.Index(str, "hello")
	println("indexVal =", indexVal)

	contains := strings.Contains(str, "中国")
	println("contains =", contains)

	splitSlice := strings.Split(str, "llo")
	println("split=", splitSlice, "len=", len(splitSlice))

	for k, v := range splitSlice {
		println("k=", k, "v=", v)
	}

	//changed by pointer
	pt := &str
	*pt = "h2"
	println("current str=", str)
}
