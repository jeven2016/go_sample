package middleware

import (
	"net/http"
	"strings"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"api/pkg/common"
	"api/pkg/dto"
	"api/pkg/global"
)

// Auth , check if the request has valid Authorization header
func Auth(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !global.GlobalApp.AuthClient.Config.EnableAuth {
			return
		}

		bearHeader := c.GetHeader(common.Authorization)
		if len(bearHeader) == 0 || !strings.Contains(bearHeader, common.Bearer) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Result{Code: http.StatusUnauthorized, Message: "unauthorized user"})
			return
		}

		// validate the access token
		items := strings.Split(bearHeader, common.Bearer)
		if len(items) > 1 {
			token := strings.Trim(items[1], " ")
			result, active := global.GlobalApp.AuthClient.Introspect(token)
			if !active {
				c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Result{Code: http.StatusUnauthorized, Message: "unauthorized user"})
				return
			}
			data, err := convertor.ToJson(result)
			if err != nil {
				log.Warn("the data isn't in json format", zap.String("result", result.String()))
			}
			log.Info("Valid token checked", zap.String("url", c.Request.URL.String()), zap.String("result", data))
			// permissions, err := global.GlobalApp.AuthClient.Client.GetUserPermissions(context.Background(), token,
			//	global.GlobalApp.AuthClient.Config.KeycloakRealm, gocloak.GetUserPermissionParams{})
			// for p := range permissions {
			//	json, _ := convertor.ToJson(p)
			//	log.Info("p", zap.String("json", json))
			// }
			//
			// path := c.Request.URL.Path
			//
			// log.Info("path" + path)
			// 根据用户的accessToken获取rptToken, 获取用户所有的资源权限
			rptToken, err := global.GlobalApp.AuthClient.GetRptToken(token)
			checkToken, valid := global.GlobalApp.AuthClient.Introspect(rptToken.AccessToken)

			if !valid {
				log.Warn("The RPT token is failed to retrospect", zap.String("rpt token", rptToken.AccessToken))
			} else {
				for _, p := range *checkToken.Permissions {
					var str = p.String()
					log.Info("per", zap.String("str", str))
				}
			}

			tokenNew, claims, err := global.GlobalApp.AuthClient.DecodeAccessToken(token)
			if err == nil {
				log.Info("tokenNew", zap.Bool("tn", tokenNew.Valid))
				log.Info("token2", zap.Error(claims.Valid()))
			}

			return
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Result{Code: http.StatusUnauthorized, Message: "unauthorized user"})
	}
}
