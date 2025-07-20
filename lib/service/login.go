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

	"golang.org/x/oauth2"
)

var codeVerifier string

func (s *Service) login(w http.ResponseWriter, r *http.Request) error {

	state := s.Config.AWS.State // Replace with a secure random string in production
	if err := aws.AuthInit(s.Config); err != nil {
		return fmt.Errorf("error initialing auth config: %w", err)
	}
	randLength := 32
	buf := make([]byte, randLength)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return fmt.Errorf("error generating %d random bytes: %v", randLength, err)
	}
	codeVerifier = hex.EncodeToString(buf)
	sha2 := sha256.New()
	io.WriteString(sha2, codeVerifier)
	codeChallenge := base64.RawURLEncoding.EncodeToString(sha2.Sum(nil))

	urlStr := aws.Oauth2Config.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
	)

	// url, err := url.Parse(urlStr)
	// if err != nil {
	// 	return fmt.Errorf("error parsing url: %w", err)
	// }
	// values := url.Query()
	// fmt.Println("URL queries:", values)

	http.Redirect(w, r, urlStr, http.StatusFound)

	return nil
}
