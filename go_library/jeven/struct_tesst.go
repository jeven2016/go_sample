package jeven

import (
	"bytes"
	"encoding/json"
)

type Person struct {
	Name string
	Age  int
	Desc string
}

func ConstructPerson() {
	var man Person
	man.Name = "wzj"
	man.Age = 22
	man.Desc = "hello wang"

	output(man)
}

func ConstructPointer() {
	//they both work
	var man *Person = new(Person)
	man.Name = "cc"
	(*man).Age = 22
	(*man).Desc = "hello cc"

	output(*man)
}

func Constructor() {
	//they both work
	var man Person = Person{Name: "p3", Age: 23, Desc: "what?"}

	output(man)
}

func output(p Person) string {
	bs, _ := json.Marshal(p)
	var out bytes.Buffer
	json.Indent(&out, bs, "", "\t")

	jsonStr := out.String()

	println(jsonStr)
	return jsonStr
}
