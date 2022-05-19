package main

import "fmt"

//same functions
func sum(number1 int, number2 int) (int, int) {
	return number1 + number2, number1
}

func sum2(number1 int, number2 int) (sum int, num int) {
	sum = number2 + number1
	num = number1
	return
}

type myFunc func(int, int)

func RunFunc() {
	sumA, sumB := sum(1, 2)
	fmt.Printf("sum=%v, %v\n", sumA, sumB)

	sumA, sumB = sum2(1, 2)
	fmt.Printf("sum=%v, %v\n", sumA, sumB)

	testFunc := sum
	sumA, sumB = testFunc(1, 2)
	fmt.Printf("sum=%v, %v\n", sumA, sumB)
}
