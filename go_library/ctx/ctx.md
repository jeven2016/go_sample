
## Context  
context可以用来在goroutine之间传递上下文信息，相同的context可以传递给运行在不同goroutine中的函数，上下文对于多个goroutine同时使用是安全的，context包定义了上下文类型，可以使用background、TODO创建一个上下文，在函数调用链之间传播context，也可以使用WithDeadline、WithTimeout、WithCancel 或 WithValue 创建的修改副本替换它，听起来有点绕，其实总结起就是一句话：context的作用就是在不同的goroutine之间同步请求特定的数据、取消信号以及处理请求的截止日期。

考虑下面这种情况：假如主协程中有多个任务1, 2, …m，主协程对这些任务有超时控制；而其中任务1又有多个子任务1, 2, …n，任务1对这些子任务也有自己的超时控制，那么这些子任务既要感知主协程的取消信号，也需要感知任务1的取消信号。

如果还是使用done channel的用法，我们需要定义两个done channel，子任务们需要同时监听这两个done channel。嗯，这样其实好像也还行哈。但是如果层级更深，如果这些子任务还有子任务，那么使用done channel的方式将会变得非常繁琐且混乱。

我们需要一种优雅的方案来实现这样一种机制：

上层任务取消后，所有的下层任务都会被取消；
中间某一层的任务取消后，只会将当前任务的下层任务取消，而不会影响上层的任务以及同级任务。
这个时候context就派上用场了。我们首先看看context的结构设计和实现原理。

###
https://zhuanlan.zhihu.com/p/136664236