package service

import "net/http"

func (s *Service) logout(w http.ResponseWriter, r *http.Request) error {

	s.setSessionVar(r, w, "authInfo", AuthInfo{})

	http.Redirect(w, r, "/", http.StatusFound)

	return nil
}
