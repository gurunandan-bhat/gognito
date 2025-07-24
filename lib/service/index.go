package service

import (
	"net/http"
)

type IndexPageData struct {
	Title   string
	Message string
}

func (s *Service) index(w http.ResponseWriter, r *http.Request) error {

	data := IndexPageData{
		Title:   "Home",
		Message: "This is the Home Page",
	}

	return s.render(w, r, "index.go.html", data, http.StatusOK)
}
