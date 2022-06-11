package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleIndex(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"msg": "show me the money"})
}

func ListCatalogs(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"msg": "show me the money"})
}
