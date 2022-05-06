package main

import (
	"fmt"
	routine "samples/src/goroutine"
	"samples/src/jeven"
	"samples/src/jeven/sub"
	aliasSub "samples/src/jeven/sub/sub2"
	"time"
)

func init() {
	println("main.go_sample init()")
}

//闭包
func specialFunc(num int) func() int {
	inner := num
	return func() int {
		inner++
		fmt.Printf("now inner=%v\n", inner)
		return inner
	}
}

func main() {

	sf := specialFunc(1)
	sf()
	sf()
	sf()

	sub.Hello()

	jeven.HelloJeven()
	aliasSub.Hello()

	total := func(num1 int, num2 int) int {
		return num1 + num2
	}(1, 2)

	fmt.Printf("\ntotal value=%v", total)

	jeven.OperateString()

	jeven.CatchError(1, 0)
	jeven.CatchError(2, 1)

	jeven.Constructor()
	jeven.ConstructPerson()
	jeven.ConstructPointer()

	c := jeven.Car{
		Name: "wzj",
		Date: time.Now(),
	}

	println(c.GetDesc())
	c.ToString()

	//--------------extend ------------
	p := jeven.Parent{Name: "parent"}
	child := jeven.Child{Parent: jeven.Parent{Name: "child"}, Age: 77}

	fmt.Printf("\nparent is %v\n", p)
	fmt.Printf("child is %v\n", child)
	child.ToString()

	//--------------goroutine
	jeven.MainFunc()

	//------------new --------
	jeven.TryNewSample()

	//-----------interface--------
	jeven.DogSays()

	jeven.RunPanic()

	time.Sleep(500)
	//--------- go routine
	routine.RunSimpleGoRoutine()
	routine.RunSimpleGoRoutine2()

	time.Sleep(500)

	println("==============")
	routine.TryLoopRoutine()
	routine.TryLoopRoutine_correct()

	time.Sleep(500)

	println("=========Pub-Sub============")

	routine.RunPubSub()

}
