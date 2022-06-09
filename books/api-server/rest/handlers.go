package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func handleIndex(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"msg": "show me the money"})
}
