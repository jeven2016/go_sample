package time_test

import (
	"fmt"
	"time"
)

func TimeAfter() {
	//等待3秒后继续执行
	timerAfter := time.After(3 * time.Second)
	select {
	case _, ok := <-timerAfter:
		println("ok?", ok)
	}
}

//time interval
func Time_ticket() {
	tick := time.NewTicker(3 * time.Second)

	i := 0
	for t := range tick.C {
		if i > 1 {
			tick.Stop()
			break
		}
		fmt.Printf("t=%v, t=%T\n", t, t)
		i++
	}
}

// Time_after time After
func Time_after() {
	t := time.After(3 * time.Second)
	<-t
	println("Runs after 3 seconds from Time_after")

	//第二种方法
	time.AfterFunc(3*time.Second, func() {
		println("Runs again after 3 seconds from AfterFunc")
	})
	time.Sleep(4 * time.Second)
}
