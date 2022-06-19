package initialization

import (
	"api/common"
	"api/rest"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func SetupEngine(config *common.Config, app *common.App) *gin.Engine {
	var engine = gin.Default()
	engine.Use(GinLogger(app.Log), GinRecovery(true, app.Log))
	//gin.SetMode(gin.ReleaseMode)
	rest.SetupServices(app)
	registerRoutes(engine, app.Log)
	return engine
}

func registerRoutes(engine *gin.Engine, log *zap.Logger) {
	root := engine.Group("/api/v1", func(context *gin.Context) {
		log.Info("A request incoming")
	})
	catalog := root.Group("catalogs")
	{
		catalog.GET("/", rest.ListCatalogs)
		catalog.GET("/:catalogId/articles", rest.ListArticles)
	}
	root.GET("/articles/:articleId", rest.FindArticleById)
	engine.NoRoute(func(ctx *gin.Context) { ctx.JSON(http.StatusNotFound, gin.H{}) })

}
