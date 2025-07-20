package service

import "net/http"

func (s *Service) form(w http.ResponseWriter, r *http.Request) error {

	s.Logger.Info("Reached handler")

	data := struct {
		AuthState AuthInfo
		Title     string
	}{Title: "Example Form"}

	return s.render(w, "protected.go.html", data, http.StatusOK)
}
