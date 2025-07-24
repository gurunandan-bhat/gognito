package service

import (
	"context"
	"errors"
	"fmt"
	"gognito/lib/aws"
	"net/http"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

type claimsPage struct {
	Title        string
	AccessToken  string
	RefreshToken string
	Claims       jwt.MapClaims
	Name         string
	Email        string
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
	rawToken, err := aws.Oauth2Config.Exchange(ctx, code, oauth2.VerifierOption(codeVerifier))
	if err != nil {
		return fmt.Errorf("error exchanging token: %w", err)
	}

	// Get claims from access token
	accessTokenStr := rawToken.AccessToken
	kf, err := keyfunc.NewDefaultCtx(
		context.Background(),
		[]string{s.Config.AWS.KeySetURL},
	)
	if err != nil {
		return fmt.Errorf("error fetching JW Keyset from Amazon: %w", err)
	}

	// Parse the token (do signature verification for your use case in production)
	// See: https://www.angelospanag.me/blog/verifying-a-json-web-token-from-cognito-in-go-and-gin
	// The article checks every claim for reasonable value but we only check for algorith and
	// expiration time and issuer
	token, err := jwt.Parse(
		accessTokenStr,
		kf.Keyfunc,
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(s.Config.AWS.IssuerURL),
	)
	if err != nil {
		return fmt.Errorf("error parsing token: %w", err)
	}

	// Check if the token is valid and extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid claims")
	}

	// Get claims from id token
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
		Name     string `json:"name,omitempty"`
		Email    string `json:"email"`
		Verified bool   `json:"email_verified"`
	}
	if err := idToken.Claims(&idClaims); err != nil {
		return fmt.Errorf("error extracting Claims: %w", err)
	}

	refreshTokenStr := rawToken.RefreshToken

	auth := AuthInfo{
		Name:      idClaims.Name,
		Email:     idClaims.Email,
		Expires:   time.Now().Add(time.Duration(rawToken.ExpiresIn) * time.Second),
		LogoutURL: logoutURL,
	}
	fmt.Printf("%+v\n", auth)

	if err := s.setSessionVar(r, w, "authInfo", auth); err != nil {
		return fmt.Errorf("unable to set auth value in session: %w", err)
	}

	// Prepare data for rendering the template
	pageData := claimsPage{
		Title:        "Cognito Callback with Claims",
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
		Name:         idClaims.Name,
		Email:        idClaims.Email,
		Claims:       claims,
	}

	// fmt.Printf("%+v\n", pageData)
	return s.render(w, r, "claims.go.html", pageData, http.StatusOK)
}
