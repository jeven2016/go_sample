package lock

import (
	"runtime"
	"sync"
	"time"
)

//互斥锁
var mutex sync.Mutex = sync.Mutex{}

func MutexSample() {
	i := 0

	println("CPU number=", runtime.NumCPU())
	runtime.GOMAXPROCS(3)
	go func() {
		defer mutex.Unlock()
		println("first running...")
		newI := i

		mutex.Lock()
		time.Sleep(200)
		i = newI + 1
		println("first done")
	}()

	go func() {
		mutex.Lock()
		defer mutex.Unlock()
		println("second running...")
		time.Sleep(300)
		i++
		println("second done")
	}()

	time.Sleep(3 * time.Second)
	defer println("The final i=", i)
}
