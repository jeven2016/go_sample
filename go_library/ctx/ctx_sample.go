package ctx

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func other(ctx context.Context) {
	value := ctx.Value("param")
	println("go other param:", value)

}

func WithValue() {
	// 传递共享值的context
	ctx := context.Background() // 先创建一个空的根context
	subCtx := context.WithValue(ctx, "param", "value")

	go func() {
		value := subCtx.Value("param").(string)
		println("go1 param:", value)
	}()

	go func() {
		value := subCtx.Value("param").(string)
		println("go2 param:", value)
	}()

	go other(ctx)
	runtime.Gosched()
}

// -----------------------------------------------------------------------------------------
// 创建一个可以生产-消费的队列
var msgChan = make(chan string, 20)

var wg = sync.WaitGroup{}

// WithCancel 测试通知生产消费的携程全部退出
func WithCancel() {
	println("============================================")
	wg.Add(2)

	// 创建一个可以取消的context
	ctx, cancelFunc := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		defer wg.Done()

		i := 0
		for {
			select {
			case <-ctx.Done():
				println("Cancel go routine1")
				return
			default:
				// 生产消息
				msgChan <- fmt.Sprintf("msg %d", i)
				println("go1 generate", i)
				i += 1
				time.Sleep(1 * time.Second)
			}
		}
	}(ctx)

	go func(ctx context.Context) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				println("Cancel go routine2")
				return
			case msg := <-msgChan:
				// 消费消息
				println("consume msg: ", msg)
				time.Sleep(2)
			}
		}

	}(ctx)

	// 运行九秒后停止生产和消费
	time.Sleep(5 * time.Second)
	cancelFunc()
	println("cancel strings_test now")

	// 等待所有的goroutine退出
	wg.Wait()
}
