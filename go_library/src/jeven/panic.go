package jeven

import "reflect"

// MyError 自定义的异常，只需要实现Error方法
type MyError struct {
	Message string
}

func (e *MyError) Error() string {
	return e.Message
}

func Top() {
	defer func() {
		//如果出错，进行捕捉处理
		if err := recover(); err != nil {
			println(reflect.TypeOf(err) == reflect.TypeOf(MyError{}))

			// 如果是MyError类型，则获取其中的值
			if er, ok := err.(MyError); ok {
				println(er.Message, er.Error())
			}
		}
	}()

	inner()
	println("will I execute?")
}

func inner() {
	println("will throw a error")
	//panic(errors.New("inner error"))
	panic(MyError{"inner error"}) //抛出自定义error
}

func RunPanic() {
	Top()
}
