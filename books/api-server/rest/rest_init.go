package rest

import (
	"api/common"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

func InitEngine(config *common.Config) {
	engine := gin.Default()
	engine.Use(common.GinLogger(), GinRecovery(true))
	//gin.SetMode(gin.ReleaseMode)
	registerRoutes(engine)

	bind := fmt.Sprintf("%v:%v", config.ApiServerConfig.BindAddress, config.ApiServerConfig.Port)
	if err := engine.Run(bind); err != nil {
		panic(errors.New(fmt.Sprintf("failed to start: %v", err)))
	}

}

func registerRoutes(engine *gin.Engine) {
	engine.GET("/", handleIndex)
}
