package goroutine

import (
	"strconv"
	"time"
)

/**
* select就是用来监听和channel有关的IO操作，当 IO 操作发生时，触发相应的动作。
* 如果有一个或多个IO操作可以完成，则Go运行时系统会随机的选择一个执行，否则的话，如果有default分支，则执行default分支语句，
* 如果连default都没有，则select语句会一直阻塞，直到至少有一个IO操作可以进行.
 */

var select_chan = make(chan int, 10)
var string_chan = make(chan string, 10)

func selectConsum() {
	for {
		println("receiving....")

		//当两个通道有写入操作即触发对应的操作。
		//当没有数据写入，会阻塞在chanel处
		select {
		case intVa, ok := <-select_chan:
			if ok {
				println("Get int value:", intVa)
			}

		case strVal := <-string_chan:
			println("Get string value", strVal)

			//去掉default，确保可以阻塞中
			//default:
			//	println("Something got wrong...")
		}
	}
}

func TrySelectChan() {
	go selectConsum()

	time.Sleep(5 * time.Second)
	for i := 0; i < 10; i++ {
		select_chan <- i
		string_chan <- strconv.Itoa(i) + "++"
	}

	time.Sleep(5 * time.Second)
}
