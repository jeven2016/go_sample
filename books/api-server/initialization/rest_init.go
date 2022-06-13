package initialization

import (
	"api/api"
	"api/common"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupEngine(config *common.Config, app *common.App) *gin.Engine {
	engine := gin.Default()
	engine.Use(GinLogger(app.Log), GinRecovery(true, app.Log))
	//gin.SetMode(gin.ReleaseMode)
	api.SetupServices(app)
	registerRoutes(engine, app.Log)
	return engine
}

func registerRoutes(engine *gin.Engine, log *zap.Logger) {
	root := engine.Group("/api/v1", func(context *gin.Context) {
		log.Info("A request incoming")
	})
	root.GET("/", api.HandleIndex)
	catalog := root.Group("catalogs", api.ListCatalogs)
	catalog.GET(":catalogId/articles", api.ListArticles)
}
