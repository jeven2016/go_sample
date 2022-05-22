package entity

type Person struct {
	Ignored string `json:"-" xml:"i-gnored"` //JSON 不序列化
	Name    string `json:"name" xml:"n_ame"`
	Age     int32  `json:"age" xml:"a-ge"`
	Desc    string `json:"description" xml:"d-esc"` //变量和Json里面的字段名不一样,且如果为类型零值或空值，序列化时忽略该字段
}
