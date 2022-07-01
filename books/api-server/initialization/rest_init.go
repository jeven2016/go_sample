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
	//gin无法通过Group实现根路径有资源，其下路径还有资源
	root.GET("/catalogs", rest.ListCatalogs)
	root.GET("/catalogs/:catalogId/articles", rest.ListArticles)

	root.GET("/articles", rest.Search)
	root.GET("/articles/:articleId", rest.FindArticleById)
	engine.NoRoute(func(ctx *gin.Context) { ctx.JSON(http.StatusNotFound, gin.H{}) })

}
