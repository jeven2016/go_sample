package route

import (
	"gin_web/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterRouter(engine *gin.Engine) {
	engine.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", gin.H{
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

	engine.POST("/users", func(context *gin.Context) {
		var p entity.Person
		if err := context.ShouldBindJSON(&p); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"code": "bad request", "error": err})
			return
		}
		context.JSON(http.StatusCreated, p)
	})

}
