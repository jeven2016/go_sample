package goroutine

import (
	"runtime"
	"time"
)

func RunSimpleGoRoutine() {
	go println("RunSimpleGoRoutine goroutine running...")

	println("RunSimpleGoRoutine end ")
	time.Sleep(3000)
}

func RunSimpleGoRoutine2() {
	go println("RunSimpleGoRoutine2 routine running...")

	//暂停当前M给其他M机会执行，
	runtime.Gosched()
	println("RunSimpleGoRoutine2 end ")
}
