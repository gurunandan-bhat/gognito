package service

import (
	"gognito/lib/aws"
	"net/http"
)

func (s *Service) initAWS(w http.ResponseWriter, r *http.Request) error {

	data := actionPageData{
		Title:   "AWS Log",
		Message: "This is the AWS Log Page",
	}

	aws.Init()

	return s.render(w, r, "index.go.html", data, http.StatusOK)
}
