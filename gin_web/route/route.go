package route

import (
	"gin_web/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func RegisterRouter(engine *gin.Engine) {

	//全局捕捉错误
	engine.Use(gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		if err != nil {
			msg := err.(string)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": msg})
		}
		c.Abort()
	}))

	engine.GET("/panic", func(context *gin.Context) {
		panic("panic now, haha ~~")
	})

	engine.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", gin.H{
			"message": "who are you?",
		})
	})

	engine.GET("/multi.html", func(context *gin.Context) {
		context.HTML(http.StatusOK, "multi.html", gin.H{
			"message": "who are you?",
		})
	})

	// http://localhost:8080/hello?name=XXX
	engine.GET("/hello", func(context *gin.Context) {
		if query, exists := context.GetQuery("name"); exists {
			context.Writer.Write([]byte("hello, " + query))
		}
	})

	engine.GET("/hello/person/:id", func(context *gin.Context) {
		param := context.Param("id")
		context.String(http.StatusOK, "param id is "+param)
	})

	v1 := engine.Group("v1")
	{
		//中间件拦截器
		v1.Use(func(context *gin.Context) {
			println("before:=======")
			println("path", context.FullPath())
			context.Next()
			//context.Abort()
			println("end:====== ")
		})
		v1.GET("message", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{
				"version": "v1",
				"message": "v1 message",
			})
		})
		v1.GET("status", func(context *gin.Context) {
			contentType := context.ContentType()
			context.JSON(http.StatusOK, gin.H{
				"version":     "v1",
				"status":      "ok",
				"contentType": contentType,
			})
		})
	}

	engine.GET("/users", func(context *gin.Context) {
		p := entity.Person{
			Ignored: "id",
			Name:    "wzj",
			Age:     18,
			Desc:    "say something to you",
		}
		context.JSON(http.StatusCreated, p)
	})

	//为了能够更方便的获取请求相关参数，提高开发效率，我们可以使用ShouldBind，它能够基于请求自动提取JSON，Form表单，Query等类型的值，
	//并把值绑定到指定的结构体对象
	engine.POST("/users", func(context *gin.Context) {
		var p entity.Person
		if err := context.ShouldBindJSON(&p); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"code": "bad request", "error": err})
			return
		}
		context.JSON(http.StatusCreated, p)
	})

	//异步调用
	engine.GET("/async", func(context *gin.Context) {
		//copy 能在其他goroutine中使用的拷贝
		c := context.Copy()
		go func() {
			time.Sleep(5 * time.Second)
			println("task finished after waiting for 5 seconds", c.Request.URL.String())
		}()
		c.Status(http.StatusAccepted)
	})

	//上传单个文件
	engine.POST("/upload", func(context *gin.Context) {
		file, err := context.FormFile("file")
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get the uploaded file"})
			return
		}
		err = context.SaveUploadedFile(file, "./uploadFile")
		if err != nil {
			println(err)
		}
		context.JSON(http.StatusOK, gin.H{"msg": "ok"})
	})

	//上传多个文件
	engine.POST("/multi-upload", func(context *gin.Context) {
		form, err := context.MultipartForm()
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get the uploaded file"})
			return
		}
		for _, file := range form.File["file"] {
			err := context.SaveUploadedFile(file, "./"+file.Filename)
			if err != nil {
				context.Status(http.StatusInternalServerError)
				println(err)
			}
		}

		context.JSON(http.StatusOK, gin.H{"msg": "ok"})
	})

	// xml data
	engine.POST("/xml", func(context *gin.Context) {
		var person entity.Person
		err := context.ShouldBindXML(&person)
		if err != nil {
			context.XML(http.StatusBadRequest, gin.H{"message": "data isn't in xml format", "error": err})
			//context.Abort()
			return
		}
		context.XML(http.StatusCreated, person)
	})

}
