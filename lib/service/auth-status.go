package service

import (
	"context"
	"errors"
	"fmt"
	"gognito/lib/aws"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
)

type claimsPage struct {
	Title        string
	AccessToken  string
	RefreshToken string
	Claims       jwt.MapClaims
	Email        string
	CurrVal      int
}

func (s *Service) handleCallback(w http.ResponseWriter, r *http.Request) error {

	ctx := context.Background()
	code := r.URL.Query().Get("code")

	// Check no one tampered with the request
	state := r.URL.Query().Get("state")
	if state != s.Config.AWS.State {
		return errors.New("stae was modified!")
	}

	// Exchange the authorization code for a token
	rawToken, err := aws.Oauth2Config.Exchange(ctx, code)
	if err != nil {
		return fmt.Errorf("error exchanging token: %w", err)
	}
	accessTokenStr := rawToken.AccessToken
	refreshTokenStr := rawToken.RefreshToken

	// Now extract id token
	rawIDToken, ok := rawToken.Extra("id_token").(string)
	if !ok {
		return errors.New("no id token found")
	}

	p, err := oidc.NewProvider(ctx, s.Config.AWS.IssuerURL)
	if err != nil {
		return fmt.Errorf("error creating provider: %w", err)
	}
	verifier := p.Verifier(&oidc.Config{ClientID: s.Config.AWS.ClientID})
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return fmt.Errorf("error verifying ID token: %w", err)
	}
	var idClaims struct {
		Email    string `json:"email"`
		Verified bool   `json:"email_verified"`
	}
	if err := idToken.Claims(&idClaims); err != nil {
		return fmt.Errorf("error extracting Claims: %w", err)
	}

	// Parse the token (do signature verification for your use case in production)
	token, _, err := new(jwt.Parser).ParseUnverified(accessTokenStr, jwt.MapClaims{})
	if err != nil {
		return fmt.Errorf("error parsing token: %w", err)
	}

	// Check if the token is valid and extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid claims")
	}
	fmt.Printf("Type of roles is %T\n", claims["cognito:groups"])

	// auth := AuthInfo{
	// 	Email:        idClaims.Email,
	// 	AccessToken:  accessTokenStr,
	// 	RefreshToken: refreshTokenStr,
	// 	Expires:      rawToken.Expiry,
	// }
	// fmt.Printf("%+v\n", auth)

	// if err := s.setSessionVar(r, w, "authInfo", auth); err != nil {
	// 	return fmt.Errorf("unable to set auth value in session: %w", err)
	// }

	// Prepare data for rendering the template
	pageData := claimsPage{
		Title:        "Cognito Callback with Claims",
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
		// Email:        idClaims.Email,
		Claims: claims,
	}

	return s.render(w, "claims.go.html", pageData, http.StatusOK)
}
