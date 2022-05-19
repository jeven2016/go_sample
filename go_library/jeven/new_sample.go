package jeven

type MyType string

// add method to *MyType
func (selfRef *MyType) say(msg string) {
	println(*selfRef, msg)
}

func (selfRef *MyType) add(msg string) {
	println(*selfRef, msg)
}

func TryNewSample() {
	my := new(MyType)
	my.say("say")
	my.add("msg")
	println("the final msg is", *my)
}
