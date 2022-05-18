package entity

type Person struct {
	Id   string `json:"-"` //JSON 不序列化
	Name string `json:"name"`
	Age  string `json:"age"`
	Desc string `json:"description", omitempty` //变量和Json里面的字段名不一样,且如果为类型零值或空值，序列化时忽略该字段
}
