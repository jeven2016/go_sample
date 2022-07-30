package middleware

import (
	"api/pkg/dto"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func Auth(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearHeader := c.GetHeader("Authorization")
		if len(bearHeader) == 0 || !strings.Contains(bearHeader, "bear") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Result{Code: http.StatusUnauthorized, Message: "unauthorized user"})
			return
		}
	}
}
