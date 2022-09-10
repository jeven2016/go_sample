package jeven

type Animal interface {
	Say(msg string) string
	ToString() string
}
type Dog struct {
	Msg string
}

type Cc struct {
	Dog
}

func (d *Cc) Say(msg string) string {
	d.Msg = d.Msg + " " + msg
	println("says:", d.Msg)
	return d.Msg
}

// Say 这里使用了指针，为了修改里面的变量值，此时是*Dog实现了接口，而不是Dog
// 一个结构体提供了接口定义的所有方法即认为实现的接口
func (d *Dog) Say(msg string) string {
	d.Msg = d.Msg + " " + msg
	println("says:", d.Msg)
	return d.Msg
}

func (d Dog) ToString() string {
	return "ToString: " + d.Msg
}

// //////cc

func animalSays(a Animal) {
	println(a.ToString())
}

func DogSays() {
	dog := new(Dog)
	dog.Msg = "dog"

	// 不能使用下面这个方法创建，因为say使用了指针，*Dog实现了接口的方法而不是Doga， 且nimalSays里面需要指针类型的，所以这个会报错
	// dog := Dog("dog")

	dog.Say("something")
	animalSays(dog)
	// var b Animal = Cc{} //不可，Cc没有直接实现接口，不可赋值
	// fmt.Printf("b is %v", b)
}
