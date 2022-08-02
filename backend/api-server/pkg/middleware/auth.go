package middleware

import (
	"api/pkg/dto"
	"api/pkg/global"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func Auth(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearHeader := c.GetHeader("Authorization")
		if len(bearHeader) == 0 || !strings.Contains(bearHeader, "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Result{Code: http.StatusUnauthorized, Message: "unauthorized user"})
			return
		}

		items := strings.Split(bearHeader, "Bearer")
		if len(items) > 1 {
			token := strings.Trim(items[1], " ")
			_, active := global.GlobalApp.AuthClient.CheckAccessToken(token)
			if !active {
				c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Result{Code: http.StatusUnauthorized, Message: "unauthorized user"})
			}
			log.Info("checked token")
			return
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Result{Code: http.StatusUnauthorized, Message: "unauthorized user"})
	}
}
