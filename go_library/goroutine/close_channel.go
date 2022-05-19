package goroutine

import (
	"runtime"
	"strconv"
)

//创建一个容量为3的chanel
var chanel2 = make(chan string, 3)

func send() {
	for i := 0; i < 3; i++ {
		chanel <- strconv.Itoa(i)
	}
	println("msg sent")
	//发送完后就关闭channel， 其中的数据仍然可以被消费到
	close(chanel2)
	println("channel closed")
}

func receive() {
	println("receiver started...")
	//1) 如果写端没有写数据，也没有关闭。<-ch; 会阻塞
	//2）如果写端写数据， value 保存 <-ch 读到的数据。 ok 被设置为 true
	//3）如果写端关闭。 value 为数据类型默认值。ok 被设置为 false
	for {
		//channel没有关闭的情况下，如果没有数据会阻塞
		if msg, ok := <-chanel2; ok {
			println("receiver consumes a msg :", msg)
		} else {
			//chanel关闭则中断
			break
		}
	}
}

func RunCloseChannel() {
	go send()
	go receive()

	runtime.Gosched()
}
