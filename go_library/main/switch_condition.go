package main

import (
	"fmt"
	"math/rand"
)

func getVal() int {
	intVal := rand.Int()
	return intVal
}

func SwitchCase() {
	val := 22

	switch val {
	case 10, 30, 40:
		fmt.Println("val is 10")

	case 20:
		fmt.Println("val is 20")

	default:
		fmt.Println("val cannot distinguish")

	}

	switch val := getVal(); {
	case val > 10:
		println("val > 10")

	case val <= 10:
		println("val <=10")

	}
}
