package service

import (
	"net/http"
)

type AboutPageData struct {
	Title   string
	Message string
}

func (s *Service) about(w http.ResponseWriter, r *http.Request) error {

	data := AboutPageData{
		Title:   "About",
		Message: "This is the About Page",
	}

	return s.render(w, r, "index.go.html", data, http.StatusOK)
}
