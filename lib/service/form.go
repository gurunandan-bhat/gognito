package service

import "net/http"

func (s *Service) form(w http.ResponseWriter, r *http.Request) error {

	s.Logger.Info("Reached handler")

	data := struct {
		Title string
	}{Title: "Example Form"}

	return s.render(w, r, "protected.go.html", data, http.StatusOK)
}
