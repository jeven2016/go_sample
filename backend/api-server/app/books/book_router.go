package books

import (
	"api/pkg/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterBooks(root *gin.RouterGroup, log *zap.Logger) {
	//gin无法通过Group实现根路径有资源，其下路径还有资源
	catalogsGroup := root.Group("/catalogs").Use(middleware.Auth(log))
	{
		catalogsGroup.GET("", ListCatalogs)
		catalogsGroup.GET("/:catalogId/articles", ListArticles)
	}

	articlesGroup := root.Group("/articles").Use(middleware.Auth(log))
	{
		articlesGroup.GET("", Search)
		articlesGroup.GET("/:articleId", FindArticleById)
	}
}
