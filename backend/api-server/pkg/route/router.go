package route

import (
	"api/app/books"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func Register(engine *gin.Engine, log *zap.Logger) {
	root := engine.Group("/api/v1", func(context *gin.Context) {
		log.Info("A request incoming")
	})
	books.RegisterBooks(root, log)
	engine.NoRoute(func(ctx *gin.Context) { ctx.JSON(http.StatusNotFound, gin.H{}) })

}
