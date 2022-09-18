package strings_test

import "fmt"

/**
GO 中没有枚举， 但利用 const 实现一个枚举
*/
const (
	red   = 0
	green = 1
	brown = 3
)

/**
iota 是一个常量计数器， 数值递增
*/
const (
	left  = iota
	right = iota
	top   = iota
)

func testEnum(value int) bool {
	return value == red
}
func testIota() {
	fmt.Println(left, right, top)
}

func testConst() {
	const constVal int = 333
	const constVal2 int = 444

	const constVal3, constVal4 = constVal, constVal2

	fmt.Println(constVal3, constVal4)
}
