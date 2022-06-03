package route

import (
	"bytes"
	"encoding/json"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	clientId     = "web2"
	clientSecret = ""
	//clientSecret = "MHp08OG1qxFZ7ByaonKe8I9SiYkMuplH"
)

func baseUrl(realm string) string {
	return "http://localhost:8080/realms/" + realm + "/protocol/openid-connect"
}

func ConfigKcRoute(engine *gin.Engine) {

	//接受浏览器传递的realm参数，返回302并使浏览器跳转至keycloak登录页面
	engine.GET("/backend/auth/:realm", func(ctx *gin.Context) {
		realm := ctx.Param("realm")

		if len(realm) == 0 {
			ctx.String(http.StatusOK, "非法参数，没有指定realm")
			log.Println("bad realm:", realm)
			return
		}

		//创建URL ： 302：申请code的URL
		state := uuid.NewString()
		nonce := uuid.NewString()

		baseUrl := baseUrl(realm) + "/auth"
		values := url.Values{}
		values.Add("client_id", clientId)
		if len(clientSecret) > 0 {
			values.Add("client_secret", clientSecret)
		}
		values.Add("redirect_uri", "http://localhost:9999/oauth?realm="+realm)
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
		ctx.Redirect(http.StatusFound, baseUrl+"?"+params)
	})

	//根据授权码申请token
	engine.GET("/oauth", func(ctx *gin.Context) {
		values := ctx.Request.URL.Query()

		//获取授权码
		code := values.Get("code")
		realm := values.Get("realm")

		if len(code) == 0 {
			ctx.String(http.StatusBadRequest, "No valid code provided")
			log.Println("failed to get code")
			return
		}

		if len(realm) == 0 {
			ctx.String(http.StatusBadRequest, "No valid realm provided")
			log.Println("failed to get realm")
			return
		}

		log.Println("realm=", realm, "code=", code)

		//申请token
		tokenUrl := "http://localhost:8080/realms/" + realm + "/protocol/openid-connect/token"
		values = url.Values{}
		values.Add("grant_type", "authorization_code")
		values.Add("code", code)
		values.Add("client_id", clientId)
		if len(clientSecret) > 0 {
			values.Add("client_secret", clientSecret)
		}
		//redirect_uri
		//REQUIRED, if the "redirect_uri" parameter was included in the
		//authorization request as described in Section 4.1.1, and their
		//values MUST be identical.
		values.Add("redirect_uri", "http://localhost:9999/oauth?realm="+realm)

		req, _ := http.NewRequest("POST", tokenUrl, bytes.NewBufferString(values.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Cache-Control", "no-cache")

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			log.Println("error /oauth:", err.Error())
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Failed to read body", err)
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		println("raw response:", string(body))
		//jsonData, err := convertor.ToJson(string(body))
		if err != nil {
			log.Println("Failed to read body", err)
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		stringBody := string(body)
		var v interface{}
		json.Unmarshal([]byte(stringBody), &v)
		data := v.(map[string]interface{})
		accessToken := data["access_token"]
		refreshTtoken := data["refresh_token"]
		idToken := data["id_token"]

		delete(data, "access_token")
		delete(data, "refresh_token")
		delete(data, "id_token")

		tokenInfo, _ := convertor.ToJson(data)
		ctx.SetCookie("access_token", accessToken.(string), 0, "/", "localhost", false, false)
		ctx.SetCookie("refresh_token", refreshTtoken.(string), 0, "/", "localhost", false, false)
		ctx.SetCookie("id_token", idToken.(string), 0, "/", "localhost", false, false)
		ctx.SetCookie("token_info", tokenInfo, 0, "/", "localhost", false, false)

		//跳转页面首页
		ctx.Redirect(http.StatusFound, "http://localhost:3000/platform/customer/no-customers")
	})

	//获取token中的用户信息
	engine.GET("/realms/:realm/userinfo", func(context *gin.Context) {
		token := context.GetHeader("token")

		realm := context.Param("realm")

		if len(realm) == 0 {
			context.String(http.StatusBadRequest, "No valid realm provided")
			log.Println("failed to get realm")
			return
		}

		//查看用户信息
		tokenUrl := baseUrl(realm) + "/userinfo"
		values := url.Values{}
		values.Add("client_id", clientId)
		if len(clientSecret) > 0 {
			values.Add("client_secret", clientSecret)
		}
		values.Add("access_token", token)

		paramString := values.Encode()
		req, _ := http.NewRequest("POST", tokenUrl, bytes.NewBufferString(paramString))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Cache-Control", "no-cache")
		//req.Header.Add("Authorization", "bearer "+token)
		//req.Header.Add("Accept", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			log.Println("error /oauth:", err.Error())
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Failed to read body", err)
			context.String(http.StatusInternalServerError, err.Error())
			return
		}
		stringBody := string(body)
		context.JSON(http.StatusOK, stringBody)
	})

	//退出登录
	engine.GET("/logout", func(context *gin.Context) {
		println("callback now")
		context.String(http.StatusOK, "ok now~~~~~")
	})

	///token?action=check
	///token?action=refresh
	engine.GET("/token?action=check", func(context *gin.Context) {
		println("callback now")
		context.String(http.StatusOK, "ok now~~~~~")
	})

}
