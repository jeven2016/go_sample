package jeven

// State 自定义一个类型
type State int

//通过使用 const 来定义一连串的常量，并借助 iota 常量计数器，来快速的为数值类型的常量连续赋值，
//非常方便。虽然没有了 enum 关键字，在这种情况下发现，是多余的，枚举本质上就是常量的组合。
const (
	Begin = iota
	Running
	End
)

//给这个类型结构体添加新的方法
func (s State) String() string {
	switch s {
	case Begin:
		return "Begin"
	case Running:
		return "Running"
	case End:
		return "End"
	}
	return "NA"
}

func RunState() {
	var state State = 1
	println(state.String())
}
