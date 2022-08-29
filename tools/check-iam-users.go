package main

import (
	"context"
	"github.com/Nerzal/gocloak/v11"
	flag "github.com/spf13/pflag"
	"move-repository/pkg/department"
)

var adminUser *string = flag.StringP("admin-user", "u", "", "The username")
var adminPwd *string = flag.StringP("admin-password", "p", "", "The password")
var kcUri *string = flag.StringP("keycloak-uri", "i", "", "The uri of keycloak")
var kcRealm *string = flag.StringP("keycloak-realm", "r", "", "The realm of keycloak")
var clientId *string = flag.StringP("client-id", "c", "", "The client's id")
var secret *string = flag.StringP("client-secret", "s", "", "The client's secret")
var oaUserFilePath *string = flag.StringP("source-oa-user-file", "s", "/home/jujucom/Desktop/workspace/projects/go_samples/tools/conf/oa-user.json", "The source json file exported from OA")

func main() {
	users := department.Import[department.OaUser](oaUserFilePath)

	client := gocloak.NewClient(*kcUri, gocloak.SetAuthRealms("realms"),
		gocloak.SetAuthAdminRealms("admin/realms"))

	token, err := client.GetToken(context.Background(), *kcRealm, gocloak.TokenOptions{
		ClientID:     clientId,
		ClientSecret: secret,
		Username:     adminUser,
		Password:     adminPwd,
	})
	department.HandleError(err)

	oaUserMap := make(map[string]string)
	for u := range *users {
		oaUserMap[u]
	}

	kcUsers, err := client.GetUsers(context.Background(), token.AccessToken, *kcRealm, gocloak.GetUsersParams{})
	department.HandleError(err)
}
