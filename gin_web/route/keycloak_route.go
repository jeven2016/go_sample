package route

import (
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"log"
	"net/http"
	"net/url"
)

func ConfigKcRoute(engine *gin.Engine) {

	engine.GET("/backend/auth/:realm", func(ctx *gin.Context) {
		realm := ctx.Param("realm")

		if len(realm) == 0 {
			ctx.String(http.StatusOK, "非法参数，没有指定realm")
			log.Println("bad realm:", realm)
			return
		}

		state := uuid.NewString()
		nonce := uuid.NewString()

		baseUrl := "http://localhost:8080/realms/zhongfu/protocol/openid-connect/auth"
		values := url.Values{}
		values.Add("client_id", "web1")
		values.Add("redirect_uri", "http://localhost:9999/oauth/")
		values.Add("state", state)
		values.Add("nonce", nonce)
		values.Add("response_mode", "query")
		values.Add("response_type", "code")
		values.Add("scope", "openid+email")

		params, err := url.QueryUnescape(values.Encode())
		if err != nil {
			ctx.String(http.StatusInternalServerError, "内部错误")
			log.Println("An error occurs(/backend/auth/:realm):", err.Error())
		}
		//ctx.Request.URL.Path = baseUrl + "?" + params
		//engine.HandleContext(ctx)
		ctx.Redirect(http.StatusFound, baseUrl+"?"+params)
	})

	engine.GET("/oauth", func(ctx *gin.Context) {
		values := ctx.Request.URL.Query()

		//获取授权码
		code := values.Get("code")

		//申请token
		//client := gocloak.NewClient("http://localhost:8080/", gocloak.SetAuthAdminRealms("admin/realms"), gocloak.SetAuthRealms("realms"))
		//context.Background()
		tokenUrl := "http://localhost:8080/realms/zhongfu/protocol/openid-connect/token"

		client := resty.New()

		// POST JSON string
		// No need to set content type, if you have client level setting
		resp, err := client.R().SetQueryParams(map[string]string{
			"code":      code,
			"client_id": "web1",
		}).SetHeader("Content-Type", "application/json").Post(tokenUrl)

		if err != nil {
			log.Println("error /oauth:", err.Error())
		}

		log.Println(resp)
	})
}
