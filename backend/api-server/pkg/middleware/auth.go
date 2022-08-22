package middleware

import (
	"api/pkg/common"
	"api/pkg/dto"
	"api/pkg/global"
	"context"
	"github.com/Nerzal/gocloak/v11"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

// Auth , check if the request has valid Authorization header
func Auth(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearHeader := c.GetHeader(common.Authorization)
		if len(bearHeader) == 0 || !strings.Contains(bearHeader, common.Bearer) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Result{Code: http.StatusUnauthorized, Message: "unauthorized user"})
			return
		}

		//validate the access token
		items := strings.Split(bearHeader, common.Bearer)
		if len(items) > 1 {
			token := strings.Trim(items[1], " ")
			result, active := global.GlobalApp.AuthClient.CheckAccessToken(token)
			if !active {
				c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Result{Code: http.StatusUnauthorized, Message: "unauthorized user"})
			}
			data, err := convertor.ToJson(result)
			if err != nil {
				log.Warn("the data isn't in json format", zap.String("result", result.String()))
			}
			log.Info("Valid token checked", zap.String("url", c.Request.URL.String()), zap.String("result", data))

			var resourceName = "admin resource"
			permissions, err := global.GlobalApp.AuthClient.Client.GetUserPermissions(context.Background(), token, "jeven", gocloak.GetUserPermissionParams{
				ResourceID: &resourceName,
			})
			for p := range permissions {
				json, _ := convertor.ToJson(p)
				log.Info("p", zap.String("json", json))
			}

			//获取用户所有的资源权限token，然后introspect
			partyToken2, err := global.GlobalApp.AuthClient.Client.GetRequestingPartyToken(context.Background(), token, "jeven", gocloak.RequestingPartyTokenOptions{
				Audience: gocloak.StringP("api-server"),
			})

			partyPermissions, err := global.GlobalApp.AuthClient.Client.
				GetRequestingPartyPermissions(context.Background(), partyToken2.AccessToken, "jeven",
					gocloak.RequestingPartyTokenOptions{})

			for _, client := range *partyPermissions {
				var str = client.String()
				log.Info("per", zap.String("str", str))
			}

			log.Info("partyToken2", zap.String("partyToken2", partyToken2.IDToken))
			return
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Result{Code: http.StatusUnauthorized, Message: "unauthorized user"})
	}
}
