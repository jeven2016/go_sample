package api

import (
	"api/common"
	"api/entity"
	"api/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

var once sync.Once
var catalogService *service.CatalogService
var articleService *service.ArticleService

func SetupServices(app *common.App) {
	once.Do(func() {
		catalogService = service.NewCatalogService(app)
		articleService = service.NewArticleService(app)
	})

}

func HandleIndex(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"msg": "show me the money"})
}

func ListCatalogs(context *gin.Context) {
	catalogs, err := catalogService.List()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"payload": []entity.BookCatalog{},
		})
		return
	}
	list := catalogs
	context.JSON(http.StatusOK, gin.H{"payload": list})
}

func ListArticles(context *gin.Context) {
	catalogId := context.Param("catalogId")
	if len(catalogId) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "the catalog id is required",
		})
		return
	}
	articleEnttiy, err := articleService.List(catalogId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An internal error encountered",
			"payload": articleEnttiy,
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{"payload": articleEnttiy})
}

func FindArticleById(context *gin.Context) *entity.Article {

}
