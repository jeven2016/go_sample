package interview

import (
	"fmt"
	"testing"
)

func TestFor(t *testing.T) {
	var array = [...]int{1, 2, 3}
	//for range 的时候它的地址会发生变化么？
	//在 for a,b := range c 遍历中， a 和 b 在内存中只会存在一份，即之后每次循环时遍历到的数据都是以值覆盖的方式
	//赋给 a 和 b，a，b 的内存地址始终不变。由于有这个特性，for 循环里面如果开协程，不要直接把 a 或者 b 的地址传给协程。
	//解决办法：在每次循环时，创建一个临时变量。
	for a, b := range array {
		println(fmt.Sprintf("a address=%v", &a))
		println(fmt.Sprintf("b address=%v", &b))
	}

}

func TestDefer(t *testing.T) {
	//go defer，多个 defer 的顺序，defer 在什么时机会修改返回值？
	//作用：defer延迟函数，释放资源，收尾工作；如释放锁，关闭文件，关闭链接；捕获panic;
	//避坑指南：defer函数紧跟在资源打开后面，否则defer可能得不到执行，导致内存泄露。
	//多个 defer 调用顺序是 LIFO（后入先出），defer后的操作可以理解为压入栈中
	//defer，return，return value（函数返回值） 执行顺序：首先return，其次return value，
	//最后defer。defer可以修改函数最终返回值，修改时机：有名返回值或者函数返回指针
	var f = func() (i int) {
		defer func() {
			i++
			fmt.Println("defer2:", i)
		}()
		defer func() {
			i++
			fmt.Println("defer1:", i)
		}()
		return i //或者直接写成return
	}
	println(f())

	var f2 = func() *int {
		var i int
		defer func() {
			i++
			fmt.Println("defer2:", i)
		}()
		defer func() {
			i++
			fmt.Println("defer1:", i)
		}()
		return &i
	}

	println(*f2())
}
