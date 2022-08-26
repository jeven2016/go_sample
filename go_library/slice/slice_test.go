package slice

import (
	"container/list"
	"encoding/json"
	"log"
	"testing"
)

func TestSliceAndPointer(t *testing.T) {
	//pointerArray()

	//testListJson()

	testPointerAddr()
}

func other(val *int) {
	log.Printf("val=%p", &val)
}

func testPointerAddr() {
	var intVal2 = 333

	var addrInt = &intVal2
	log.Printf("val begin=%p", addrInt)
	other(addrInt)

}

func testPointPerson() {
	var array = []Person{
		Person{
			Name: "internal",
		},
	}

	var name = "pp"
	var name2 = &name
	var p2 = &Person2{
		Name: name2,
		List: array,
	}
	log.Printf("%p", p2)
}

func testListJson() {
	p := Person{
		Name: "w",
	}
	p2 := Person{
		Name: "w",
	}

	p3 := Person{
		Name: "w",
	}

	var list = list.New()
	list.PushBack(p)
	list.PushBack(p2)

	p3.List = *list

	marshal, err := json.Marshal(p3)
	if err != nil {
		panic(err)
	}
	log.Printf("%v", string(marshal))
}

type Person struct {
	Name string

	List list.List
}

type Person2 struct {
	Name *string

	List []Person
}

func pointerArray() {
	var person = Person{
		Name: "test",
	}
	var array = make([]Person, 3)
	array2 := append(array, person)

	log.Printf("array=%p", &array)
	log.Printf("array2=%p", &array2)
	log.Printf("person=%p", &person)
	testArray(array2)
}

func testArray(arr []Person) {
	log.Printf("array=%p", &arr)
	log.Printf("person=%p", &arr[0])
}
