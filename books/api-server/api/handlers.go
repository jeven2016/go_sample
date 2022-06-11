package api

import (
	"api/global"
	"api/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

var once sync.Once
var catalogService *service.CatalogService

func SetupServices() {
	once.Do(func() {
		log := global.Log
		catalogService = service.New(log)
	})

}

func HandleIndex(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"msg": "show me the money"})
}

func ListCatalogs(context *gin.Context) {
	list := catalogService.List()
	context.JSON(http.StatusOK, gin.H{"payload": list})
}
