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

	url := aws.Oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOnline)

	fmt.Println("URL", url)
	http.Redirect(w, r, url, http.StatusFound)

	return nil
}
