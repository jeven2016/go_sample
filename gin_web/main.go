package main

import (
	"gin_web/route"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	engine.LoadHTMLGlob("./html/*")      //html路径
	engine.Static("/static", "./static") //静态资源映射，放置图片等静态资源

	route.RegisterRouter(engine)
	//port is 8080 by default
	if err := engine.Run(":8080"); err != nil {
		println("error occurs", err.Error())
	}
}