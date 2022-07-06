package rest

import (
	common2 "api/pkg/common"
	"api/pkg/dto"
	"api/pkg/entity"
	service2 "api/pkg/service"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"sync"
)

var once sync.Once
var catalogService *service2.CatalogService
var articleService *service2.ArticleService

func SetupServices(app *common2.App) {
	once.Do(func() {
		catalogService = service2.NewCatalogService(app)
		articleService = service2.NewArticleService(app)
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
	var articlePageRequest dto.PageRequest
	err := c.ShouldBindQuery(&articlePageRequest)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, gin.H{"error": errs.Translate(common2.Trans)})
		return
	}

	catalogId := c.Param("catalogId")
	if len(catalogId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the catalog id is required",
		})
		return
	}

	resp, err := articleService.List(catalogId, &articlePageRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}

//FindArticleById  根据id获取Article
func FindArticleById(c *gin.Context) {
	articleId := c.Param("articleId")
	articleEntity, err := articleService.FindById(articleId)
	if err != nil {
		if errors.Is(err, common2.NotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "An internal error encountered",
		})
		return
	}
	c.JSON(http.StatusOK, articleEntity)
}

func Search(c *gin.Context) {
	var articlePageRequest dto.PageRequest
	err := c.ShouldBindQuery(&articlePageRequest)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, gin.H{"error": errs.Translate(common2.Trans)})
		return
	}
	if len(articlePageRequest.Search) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "search参数未设置"})
		return
	}
	resp, err := articleService.Search(&articlePageRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}
