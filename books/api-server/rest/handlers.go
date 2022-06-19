package rest

import (
	"api/common"
	"api/dto"
	"api/entity"
	"api/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	var articlePageRequest *dto.PageRequest
	err := c.ShouldBindQuery(articlePageRequest)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, gin.H{"error": errs.Translate(common.Trans)})
		return
	}

	catalogId := c.Param("catalogId")
	if len(catalogId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the catalog id is required",
		})
		return
	}

	resp, err := articleService.List(catalogId, articlePageRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, gin.H{"payload": resp})
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
