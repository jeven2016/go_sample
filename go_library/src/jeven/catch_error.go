package jeven

func CatchError(num1 int, num2 int) int {
	defer func() {
		err := recover()
		if err != nil {
			println("error =", err)
		} else {
			println("no error")
		}
	}()

	println("preparing to run.....")
	return num1 / num2
}
