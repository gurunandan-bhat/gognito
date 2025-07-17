package service

import (
	"context"
	"errors"
	"fmt"
	"gognito/lib/aws"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

type claimsPage struct {
	Title       string
	AccessToken string
	Claims      jwt.MapClaims
}

func (s *Service) handleCallback(w http.ResponseWriter, r *http.Request) error {

	fmt.Println("I got called")

	ctx := context.Background()
	code := r.URL.Query().Get("code")

	// Exchange the authorization code for a token
	rawToken, err := aws.Oauth2Config.Exchange(ctx, code)
	if err != nil {
		return fmt.Errorf("error exchanging token: %w", err)
	}
	tokenString := rawToken.AccessToken

	// Parse the token (do signature verification for your use case in production)
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return fmt.Errorf("error parsing token: %w", err)
	}

	// Check if the token is valid and extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid claims")
	}

	// Prepare data for rendering the template
	pageData := claimsPage{
		Title:       "Cognito Callback with Claims",
		AccessToken: tokenString,
		Claims:      claims,
	}
	s.render(w, "claims.go.html", pageData, http.StatusOK)

	return nil
}
