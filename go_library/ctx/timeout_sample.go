package ctx

import (
	"context"
	"fmt"
	"sync"
	"time"
)

//-----------------------------------------------------------------------------------------
//创建一个可以生产-消费的队列
var channel = make(chan string, 20)

var wgr = sync.WaitGroup{}

// WithTimeout 测试通知生产消费的携程全部退出
func WithTimeout() {
	println("============================================")
	wgr.Add(2)

	//创建一个3秒后超时的context, 也可以提前退出
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	//context.WithDeadline(context.Background(), time.Time{ })

	go func(ctx context.Context) {
		defer wgr.Done()

		i := 0
		for {
			select {
			case <-ctx.Done():
				println("Cancel go routine1")
				return
			default:
				//生产消息
				channel <- fmt.Sprintf("msg %d", i)
				println("go1 generate", i)
				i += 1
				time.Sleep(1 * time.Second)
			}
		}
	}(ctx)

	go func(ctx context.Context) {
		defer wgr.Done()
		for {
			select {
			case <-ctx.Done():
				println("Cancel go routine2")
				return
			case msg := <-channel:
				//消费消息
				println("consume msg: ", msg)
				time.Sleep(2)
			}
		}

	}(ctx)

	//运行九秒后停止生产和消费
	time.Sleep(5 * time.Second)
	//cancelFunc()
	//println("cancel main now")

	//等待所有的goroutine退出
	wgr.Wait()
}
