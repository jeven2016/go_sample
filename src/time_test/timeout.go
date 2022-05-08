package time_test

import "time"

func TimeoutFunc() {
	time.Sleep(3 * time.Second)

	//3秒后执行一次
	timer := time.NewTimer(3 * time.Second)
	select {
	case <-timer.C: //超时时执行
		println("只会执行一次")
		timer.Stop()
	}

	//执行两次，不重复创建timer
	timer2 := time.NewTimer(3 * time.Second)
	times := 0

	for {
		select {
		case _, ok := <-timer2.C:
			if ok {
				println("index", times)
				times++
				timer2.Reset(3 * time.Second)
				println("rest timer3")
			}

		}
		if times > 1 {
			stopped := timer2.Stop()
			if stopped {
				println("timer2 is stopped")
			}
		}
	}

}
