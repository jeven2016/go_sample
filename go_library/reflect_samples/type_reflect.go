package reflect_samples

import (
	"fmt"
	"reflect"
)

type MyStruct struct {
	Name string
	Age  int32 `type:"int" author:"wang" `
}

func ReflectType() {
	ptr := &MyStruct{Name: "who", Age: 33}

	//名称
	printType(32)
	printType(int64(44))
	printType(MyStruct{})
	printType(ptr)

	//类型: uint
	printType(32)
	printType(77)
	printType(int64(44))
	printType(MyStruct{})
	printType(ptr)

	//遍历结构体所有属性(不能是pointer)
	ptrType := reflect.TypeOf(*ptr)
	for i := 0; i < ptrType.NumField(); i++ {
		field := ptrType.Field(i)
		println(field.Type.Name(), field.Name, field.Index, field.Tag, field.PkgPath, field.IsExported())
	}

	//获取字段属性，以及Tag
	if field, exists := ptrType.FieldByName("Age"); exists {
		tag := field.Tag.Get("author")
		print("Age tag :", tag)
	}

}

func printType(val interface{}) {
	valType := reflect.TypeOf(val)
	fmt.Printf("name=%v, kind=%v\n", valType.Name(), valType.Kind())
}
