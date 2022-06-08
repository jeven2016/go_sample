## Project description
**Books online** 

### Packages installed
```shell
go get github.com/duke-git/lancet/v2

go get github.com/gin-gonic/gin

#配置文件读取, Viper是适用于Go应用程序的完整配置解决方案。它被设计用于在应用程序中工作，并且可以处理所有类型的配置需求和格式。
go get github.com/spf13/viper

# 高性能日志
go get -u go.uber.org/zap
go get github.com/natefinch/lumberjack
```

#### Zap log
zap提供了两种类型的日志记录器Logger 和Sugared Logger

区别是：
在每一微秒和每一次内存分配都很重要的上下文中，使用Logger。它比Sugared Logger更快，内存分配次数也更少，但它只支持强类型的结构化日志记录。
在性能不需要太考虑时，使用SugaredLogger。并且支持结构化和printf风格的日志记录。
通过调用zap.NewProduction()/zap.NewDevelopment()或者zap.NewExample()创建一个Logger。
上面的每一个函数都将创建一个Logger。唯一的区别在于它将记录的信息不同。例如：production logger默认记录调用函数信息、日期和时间等。
通过Logger 调用INFO、ERROR等。
* 对比总结  
Example和Production使用的是json格式输出，development使用行的形式输出
Development
从警告级别向上打印到堆栈中来跟踪
始终打印包/文件/行（方法）
在行尾添加任何额外字段作为json字符串
以大写形式打印级别名称
以毫秒为单位打印ISO8601格式的时间戳

* Production  
调试级别消息不记录
Error,panic级别的记录，会在堆栈中跟踪文件，warn不会
始终将调用者添加到文件中
以时间戳格式打印日期
以小写形式打印级别名称

* 使用lumberjack进行日志切割归档  
因为zap本身不支持切割归档日志文件，为了添加日志切割归档功能，我们将使用第三方库lumberjack来实现