package string_test

import (
	"fmt"
	"strconv"
)

func convertString() {
	var boolVal bool = true
	fmt.Printf("boolValue=%v  type=%T\n", boolVal, boolVal)

	// 返回两个参数，不关心第二个error返回，
	boolVal, _ = strconv.ParseBool("true")
	fmt.Printf("parsed from string boolValue=%v type=%T\n", boolVal, boolVal)

	var intValue int64
	intValue, _ = strconv.ParseInt("19", 10, 64)
	fmt.Printf("int 64 value=%v", intValue)

	var value32 int32 = int32(intValue)
	fmt.Printf("int 32 value=%v", value32)
}

func PublicMethod() {
	convertString()
}

func RunConvert() {
	convertString()
}
