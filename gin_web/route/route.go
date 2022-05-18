package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterRouter(engine *gin.Engine) {

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

}
