package service

import (
	"fmt"
	"gognito/lib/aws"
	"net/http"

	"golang.org/x/oauth2"
)

func (s *Service) login(w http.ResponseWriter, r *http.Request) error {

	state := s.Config.AWS.State // Replace with a secure random string in production
	if err := aws.AuthInit(s.Config); err != nil {
		return fmt.Errorf("error initialing auth config: %w", err)
	}

	urlStr := aws.Oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOnline)
	// url, err := url.Parse(urlStr)
	// if err != nil {
	// 	return fmt.Errorf("error parsing url: %w", err)
	// }
	// values := url.Query()
	// fmt.Println("URL queries:", values)

	http.Redirect(w, r, urlStr, http.StatusFound)

	return nil
}
