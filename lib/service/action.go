package service

import (
	"net/http"
)

type actionPageData struct {
	Title   string
	Message string
}

func (s *Service) action(w http.ResponseWriter, r *http.Request) error {

	data := actionPageData{
		Title:   "Action",
		Message: "This is the Action Page",
	}

	return s.render(w, r, "index.go.html", data, http.StatusOK)
}
