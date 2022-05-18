package main

import (
	"gin_web/route"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	route.RegisterRouter(engine)
	//port is 8080 by default
	if err := engine.Run(":8080"); err != nil {
		println("error occurs", err.Error())
	}
}
