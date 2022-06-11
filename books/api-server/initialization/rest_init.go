package initialization

import (
	"api/api"
	"api/common"
	"github.com/gin-gonic/gin"
)

func SetupEngine(config *common.Config) *gin.Engine {
	engine := gin.Default()
	engine.Use(common.GinLogger(), GinRecovery(true))
	//gin.SetMode(gin.ReleaseMode)
	api.SetupServices()
	registerRoutes(engine)
	return engine
}

func registerRoutes(engine *gin.Engine) {
	engine.GET("/", api.HandleIndex)
	engine.GET("/catalogs", api.ListCatalogs)
}
