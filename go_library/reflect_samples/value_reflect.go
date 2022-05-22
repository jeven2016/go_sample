package reflect_samples

import "reflect"

// ReflectValue only run with go run value_reflect.go
func ReflectValue() {
	intVal := int32(64)
	floatVal := float32(4.0)

	intValType := reflect.ValueOf(intVal)
	floatValType := reflect.ValueOf(floatVal)

	//-----------------------------int32------------------------------------
	//是否可以转int
	println("int32 can be int =", intValType.CanInt())

	//是否可以转指针型
	println("int32 can Addr =", intValType.CanAddr())

	println("intVal Kind and Name:", intValType.Type().Kind(), intValType.Type().Name())

	if intValType.Kind() == reflect.Struct {
		println("int32 NumField:", intValType.NumField())
	}

	//-----------------------------float32------------------------------------
	println("float32 can be int =", floatValType.CanInt())

	//是否可以转指针型
	println("float32 can Addr =", floatValType.CanAddr())

	println("float32 Kind and Name:", floatValType.Type().Kind(), floatValType.Type().Name())

	if intValType.Kind() == reflect.Struct {
		println("float32 NumField:", floatValType.NumField())
	}
}
