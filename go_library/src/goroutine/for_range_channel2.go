package goroutine

import (
	"runtime"
	"strconv"
)

//创建一个容量为3的chanel
var chanelFr = make(chan string, 3)

//限制为单项通道，表示通道只能接受消息. 单通道无法通过代码转换成双通道
func sendFr(ch chan<- string) {
	for i := 0; i < 4; i++ {
		ch <- strconv.Itoa(i)
	}
	println("chanelFr msg sent")
	//发送完后就关闭channel， 其中的数据仍然可以被消费到
	close(ch)
	println("chanelFr closed")
}

//限制为单项通道，表示只能从通道获取数据，不能写入
func receiveFr(ch <-chan string) {
	println("receiveFr started...")

	//可以通过for range访问通道
	for msg := range ch {
		println("receiveFr consumes a msg :", msg)
	}
}

func RunSingleChannelFr() {
	go sendFr(chanelFr)
	go receiveFr(chanelFr)

	runtime.Gosched()
}
