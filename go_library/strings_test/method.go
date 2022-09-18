package strings_test

func calc(numA int, numB int) int {
	return numB + numA
}

func calc2(numA int, numB int) (int, int, int) {
	return numB + numA, numA, numB
}

// func strings_test() {
//	fmt.Printf("number is %v", calc(10, 20))
//
//	numberA, numberB, numberC := calc2(10, 20)
//
//	fmt.Printf("numberA=%v, numberB=%v, numberC=%v", numberA, numberB, numberC)
// }
