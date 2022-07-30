## Project description

**Books online**

### 工程标准化参考

https://github.com/golang-standards/project-layout/blob/master/README_zh.md

### Packages installed

```shell
go get github.com/duke-git/lancet/v2

go get github.com/gin-gonic/gin

#配置文件读取, Viper是适用于Go应用程序的完整配置解决方案。它被设计用于在应用程序中工作，并且可以处理所有类型的配置需求和格式。
go get github.com/spf13/viper

# 高性能日志
go get -u go.uber.org/zap
go get github.com/natefinch/lumberjack

# gin 校验
https://github.com/go-playground/validator

# hot reloading
go get github.com/silenceper/gowatch
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

### Mongodb

reference: https://www.mongodb.com/docs/drivers/go/current/fundamentals/connection/  
https://learnku.com/articles/61966
https://blog.51cto.com/u_12970189/2547519
https://onejav.com/actress/Eimi%20Fukada?page=2
https://www.cartoon18.com/

```shell
go get go.mongodb.org/mongo-driver/mongo
```

在GO中使用BSON对象
MongoDB中的JSON文档以称为BSON（二进制编码的JSON）的二进制表示形式存储。与其他将JSON数据存储为简单字符串和数字的数据库不同，
BSON编码扩展了JSON表示形式，例如int，long，date，float point和decimal128。这使应用程序更容易可靠地处理，排序和比较数据。
Go Driver有两种系列用于表示BSON数据：D系列类型和Raw系列类型。
D系列包括四种类型：

D：BSON文档。此类型应用在顺序很重要的场景下，例如MongoDB命令。
M：无序map。除不保留顺序外，与D相同。
A：一个BSON数组。
E：D中的单个元素。

打包： https://blog.csdn.net/yhflyl/article/details/120649170

### Gin普通字段的校验

```shell
import "github.com/astaxie/beego/validation"

func (a *Article) GetArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")
	var data *models.Article
	code := codes.InvalidParams
	if !valid.HasErrors() {
		data = a.Service.GetArticle(id)
		code = codes.SUCCESS
	} else {
		for _, err := range valid.Errors {
			a.Log.Info("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}
	RespData(c, http.StatusOK, code, data)
}

```

reference:
https://juejin.cn/post/7012155588280844301

### 使用国内源  
设置代理：  
```shell
go env -w GOPROXY=https://goproxy.cn,direct
```

### 集成keycloak
https://zhuanlan.zhihu.com/p/488194876
```shell

```