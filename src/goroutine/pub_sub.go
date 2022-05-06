package goroutine

import (
	"math/rand"
	"strconv"
	"time"
)

//创建一个容量为3的chanel
var chanel = make(chan string, 3)

func pub() {
	//println("push a int value")
	chanel <- strconv.Itoa(rand.Intn(32))
}

func sub() {
	println("sub started...")
	for {
		println("sub consume a msg :", <-chanel)
	}
}

func sub2() {
	println("sub2 started...")
	for {
		println("sub2 consume a msg :", <-chanel)
	}
}

func RunPubSub() {
	go sub()
	go sub2()
	pub()
	pub()
	pub()
	pub()
	pub()
	time.Sleep(3 * time.Second)
}
