package lock

import (
	"strconv"
	"sync"
	"time"
)

var rwLock sync.RWMutex

func RwLock() {
	strVal := ""

	//read
	go func() {
		for i := 0; i < 3; i++ {
			time.Sleep(200)
			rwLock.RLock() // 阻塞等待锁
			println("read value=", strVal)
			rwLock.RUnlock()
		}
	}()

	//write
	go func() {
		for i := 0; i < 3; i++ {
			time.Sleep(400)
			rwLock.Lock() // 阻塞等待锁
			strVal += strconv.Itoa(i)
			println("write value=", strVal)
			rwLock.Unlock()
			time.Sleep(200)
		}
	}()

	time.Sleep(10 * time.Second)
}
