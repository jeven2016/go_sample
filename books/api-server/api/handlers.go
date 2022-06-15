package api

import (
	"api/common"
	"api/dto"
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
	catalogList, err := catalogService.List()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"payload": []entity.BookCatalog{},
		})
		return
	}
	list := catalogList
	context.JSON(http.StatusOK, gin.H{"payload": list})
}

func ListArticles(c *gin.Context) {
	var articlePageRequest dto.PageRequest
	err := c.ShouldBindQuery(&articlePageRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	catalogId := c.Param("catalogId")
	if len(catalogId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the catalog id is required",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
	return

	articleEnttiy, err := articleService.List(catalogId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An internal error encountered",
			"payload": articleEnttiy,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"payload": articleEnttiy})
}

func FindArticleById(context *gin.Context) {
	articleId := context.Param("articleId")
	entity, err := articleService.FindById(articleId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "An internal error encountered",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{"payload": entity})
}
