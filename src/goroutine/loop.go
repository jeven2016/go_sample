package goroutine

import "runtime"

func TryLoopRoutine() {
	names := []string{"shanghai", "nanjing", "beijing"}

	//循环结束name已经是最后一个值，此时还没有执行goroutine
	for _, val := range names {
		go func() {
			println("TryLoopRoutine>>hello", val)
		}()
	}
	//等待goroutine执行，否则会提前结束
	runtime.Gosched()
}

func TryLoopRoutine_correct() {
	names := []string{"shanghai", "nanjing", "beijing"}

	//进行值拷贝，当协程进行时不会引用最后一个字符串
	for _, val := range names {
		go func(name string) {
			println("TryLoopRoutine_correct>> hello", name)
		}(val)
	}
	runtime.Gosched()
}
