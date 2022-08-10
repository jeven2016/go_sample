package clients

import (
	"api/pkg/common"
	"errors"

	//"api/pkg/common"
	"context"
	"github.com/Nerzal/gocloak/v11"
	"go.uber.org/zap"
)

type AuthClient struct {
	Client gocloak.GoCloak
	Config *common.AuthConfig
	Log    *zap.Logger
}

func (c *AuthClient) StartInit() error {
	authCfg := c.Config

	if !authCfg.EnableAuth {
		return errors.New("auth: OpenID connect feature is not enabled")
	}

	client := gocloak.NewClient(c.Config.KeycloakUrl, gocloak.SetAuthRealms("realms"),
		gocloak.SetAuthAdminRealms("admin/realms"))
	c.Client = client
	return nil
}

func (c *AuthClient) CheckAccessToken(accessToken string) (*gocloak.RetrospecTokenResult, bool) {
	cfg := c.Config
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
