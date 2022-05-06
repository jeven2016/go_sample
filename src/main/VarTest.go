package main

import "fmt"

func variable() {
	//define variables
	var a int32 = 1 //添加了类型
	var b = 32
	c := 33 //定义一个新的变量

	const ca = "const"

	fmt.Printf(`a=%d`, a)
	fmt.Printf(`b=%d`, b)
	fmt.Printf(`c=%d`, c)
}

func change() {
	a := 1
	b := 2

	//交换a和b
	a, b = b, a

	fmt.Println()
	fmt.Printf("chang: a=%d ", a)
	fmt.Printf("chang: b=%d", b)
}

func RunVar() {
	variable()
	change()
}
