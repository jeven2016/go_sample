package jeven

type Parent struct {
	Name string
}

// Child 看作Child继承了Parent
type Child struct {
	Parent
	Age int
}

// ToString 关联到Parent类型, 在Child中也可以调用Parent上的函数
func (pt *Parent) ToString() {
	println("Name=", pt.Name)
}
