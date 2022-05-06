package jeven

import (
	"sync"
)

// like CountLatchDown
var group = sync.WaitGroup{}

func running() {
	defer group.Done()
	println("goroutine running")
}

func MainFunc() {
	group.Add(2)
	go running()
	go running()

	//等计数器为0执行下面代码
	group.Wait()
	println("finish")
}
