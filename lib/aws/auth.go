package aws

import (
	"context"
	"fmt"
	"gognito/lib/config"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

var Oauth2Config = oauth2.Config{}

func AuthInit(cfg *config.Config) error {

	provider, err := oidc.NewProvider(context.Background(), cfg.AWS.IssuerURL)
	if err != nil {
		return fmt.Errorf("error creating OIDC provider %w", err)
	}

	// Set up OAuth2 config
	Oauth2Config = oauth2.Config{
		ClientID:     cfg.AWS.ClientID,
		ClientSecret: cfg.AWS.ClientSecret,
		RedirectURL:  cfg.AWS.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
	}

	return nil
}
