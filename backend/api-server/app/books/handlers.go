package books

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"api/app/books/dto"
	"api/app/books/service"
	common "api/pkg/common"
	"api/pkg/global"
)

var once sync.Once
var catalogService *service.CatalogService
var articleService *service.ArticleService

func SetupServices(app *global.App) {
	once.Do(func() {
		catalogService = service.NewCatalogService(app)
		articleService = service.NewArticleService(app)
	})

}

func ListCatalogs(context *gin.Context) {
	catalogResponse, err := catalogService.List()
	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusOK, catalogResponse)
}

func ListArticles(c *gin.Context) {
	var articlePageRequest dto.PageRequest
	err := c.ShouldBindQuery(&articlePageRequest)
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

	resp, err := articleService.List(catalogId, &articlePageRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// FindArticleById  根据id获取Article
func FindArticleById(c *gin.Context) {
	articleId := c.Param("articleId")
	articleEntity, err := articleService.FindById(articleId)
	if err != nil {
		if errors.Is(err, common.NotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "An handler error encountered",
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
		c.JSON(http.StatusBadRequest, gin.H{"error": errs.Translate(common.Trans)})
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
