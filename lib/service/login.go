package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"gognito/lib/aws"
	"io"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

var codeVerifier string
var stateLength = 32

func (s *Service) login(w http.ResponseWriter, r *http.Request) error {

	authInfo, _ := s.getSessionVar(r, "authInfo")
	if authInfo != nil {
		authData := authInfo.(*AuthInfo)
		if time.Now().Before(authData.Expires) {
			http.Redirect(w, r, "/", http.StatusFound)
			return nil
		}
	}

	state := s.Config.AWS.State // Replace with a secure random string in production
	if err := aws.AuthInit(s.Config); err != nil {
		return fmt.Errorf("error initialing auth config: %w", err)
	}

	buf := make([]byte, stateLength)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return fmt.Errorf("error generating %d random bytes: %v", stateLength, err)
	}
	codeVerifier = hex.EncodeToString(buf)
	sha2 := sha256.New()
	_, err = io.WriteString(sha2, codeVerifier)
	if err != nil {
		return fmt.Errorf("error encoding verifier to SHA256: %w", err)
	}
	codeChallenge := base64.RawURLEncoding.EncodeToString(sha2.Sum(nil))

	urlStr := aws.Oauth2Config.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
	)

	http.Redirect(w, r, urlStr, http.StatusFound)

	return nil
}
