package kc

import (
	"context"
	gocloak "github.com/Nerzal/gocloak/v11"
)

func init() {
	println("init kc now")
}

func KeycloakClient() {
	client := gocloak.NewClient("http://localhost:8080/", gocloak.SetAuthAdminRealms("admin/realms"), gocloak.SetAuthRealms("realms"))
	ctx := context.Background()
	//token, err := client.LoginAdmin(ctx, "admin", "admin", "master")
	token, err := client.Login(ctx, "web1", "vzqjfzQ7qlIphUdHNWJjIycUnla25OTk", "zhongfu", "wzj", "wzj")
	if err != nil {
		panic("Login failed:" + err.Error())
	}

	rptResult, err := client.RetrospectToken(ctx, token.AccessToken, "web1", "", "zhongfu")
	if err != nil {
		panic("Inspection failed:" + err.Error())
	}

	if !*rptResult.Active {
		panic("Token is not active")
	}

	permissions := rptResult.Permissions
	print(permissions)
	// Do something with the permissions ;
}
