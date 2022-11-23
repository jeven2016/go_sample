package clients

import (
	// "api/pkg/common"
	"context"

	"github.com/golang-jwt/jwt/v4"

	"api/pkg/common"

	"github.com/Nerzal/gocloak/v11"
	"go.uber.org/zap"
)

type AuthClient struct {
	Client gocloak.GoCloak
	Config *common.AuthConfig
	Log    *zap.Logger
}

func (c *AuthClient) StartInit() {
	authCfg := c.Config

	if !authCfg.EnableAuth {
		c.Log.Warn("Authentication is not enabled")
		return
	}

	client := gocloak.NewClient(c.Config.KeycloakUrl, gocloak.SetAuthRealms("realms"),
		gocloak.SetAuthAdminRealms("admin/realms"))
	c.Client = client
}

// Introspect Normal case:
// UI calls {backend_api}/auth/login
// Backend API redirects UI client to {keycloak}/realms/myrealm/protocol/openid-connect/auth/ with redirect URL as {backend_api}/auth/callback
// Authentication happens and the handler at {backend_api}/auth/callback gets the authorization_code.
// Above handler makes request to {{keycloak_url}}/realms/{{realm}}/protocol/openid-connect/token to get the access_token
// Set-cookie and returns to UI client.
func (c *AuthClient) Introspect(accessToken string) (*gocloak.RetrospecTokenResult, bool) {
	cfg := c.Config
	// You will need to use the public key for your keycloak realm, which you can either grab for every authentication request, or just cache.
	// To get the public key, you can use the gocloak function GetCerts,
	// and parse the keys list you get back for the signing method your keycloak server is using.
	// You can cache this since it shouldn't really change. You should be able to validate tokens without making any network calls to the keycloak server, if you wanted.
	// Another way to validate tokens is to send the token to the keycloak introspection endpoint, and let keycloak check the token.
	// It just may take a little longer because it's a network call, depending on where your backend server and keycloak server are hosted.
	// You can use the gocloak function RetrospectToken to send the token to keycloak for introspection.
	result, err := c.Client.RetrospectToken(context.TODO(), accessToken, cfg.ClientId, cfg.ClientSecret, cfg.KeycloakRealm)
	if err != nil {
		c.Log.Warn("invalid access token to parse", zap.String("accessToken", accessToken), zap.Error(err))
		return result, false
	}

	if !*result.Active {
		c.Log.Warn("the access token isn't active", zap.String("accessToken", accessToken), zap.Error(err))
		return result, false
	}
	return result, true
}

func (c *AuthClient) GetRptToken(accessToken string) (*gocloak.JWT, error) {
	cfg := c.Config

	// Audience : client id of resource server
	result, err := c.Client.GetRequestingPartyToken(context.Background(), accessToken, cfg.KeycloakRealm,
		gocloak.RequestingPartyTokenOptions{
			Audience: gocloak.StringP(cfg.ClientId),
		})
	if err != nil {
		c.Log.Warn("failed to issue RPT token", zap.String("accessToken", accessToken), zap.Error(err))
	}
	return result, err
}

func (c *AuthClient) DecodeAccessToken(accessToken string) (*jwt.Token, *jwt.MapClaims, error) {
	token, claims, err := c.Client.DecodeAccessToken(context.Background(), accessToken, c.Config.KeycloakRealm)

	if err != nil {
		c.Log.Warn("failed to issue RPT token", zap.String("accessToken", accessToken), zap.Error(err))
	}
	return token, claims, err
}
