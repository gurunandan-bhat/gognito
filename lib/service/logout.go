package service

import (
	"fmt"
	"net/http"
)

func (s *Service) logout(w http.ResponseWriter, r *http.Request) error {

	if err := s.setSessionVar(r, "authInfo", AuthInfo{}); err != nil {
		return fmt.Errorf("error delting auth session: %w", err)
	}

	http.Redirect(w, r, "/", http.StatusFound)

	return nil
}
