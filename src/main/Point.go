package main

import (
	"fmt"
)

func PointTst() {
	var intVal int32 = 56
	fmt.Printf("value is=%v", &intVal)

	var pointer *int32 = &intVal

	fmt.Printf("pointer's addres=%v\n", pointer)
	fmt.Printf("pointer's value=%v\n", *pointer)

	if count := 30; count < 30 {
		fmt.Printf("count is %v", count)
	} else if count == 40 {
		fmt.Printf("count equals to 40")
	} else {
		fmt.Printf("nothing matches")
	}
}
