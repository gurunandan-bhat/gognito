package service

import (
	"net/http"
	"time"
)

type Middleware func(serviceHandler) serviceHandler

func (s *Service) logMiddleware(next serviceHandler) serviceHandler {

	return func(w http.ResponseWriter, r *http.Request) error {
		s.Logger.Info("First Logger: Logging before")
		if err := next(w, r); err != nil {
			return err
		}
		s.Logger.Info("First Logger: Looging after")
		return nil
	}
}

func (s *Service) logAnotherMiddleware(next serviceHandler) serviceHandler {

	return func(w http.ResponseWriter, r *http.Request) error {
		s.Logger.Info("Second Logger: Logging before")
		if err := next(w, r); err != nil {
			return err
		}
		s.Logger.Info("Second Logger: Looging after")
		return nil
	}
}

func (s *Service) validateAuth(next serviceHandler) serviceHandler {

	return func(w http.ResponseWriter, r *http.Request) error {

		// Check if we have a session variable called "authInfo"
		auth, err := s.getSessionVar(r, "authInfo")
		if err != nil {
			return err
		}

		authData, ok := auth.(AuthInfo)
		if ok || time.Now().After(authData.Expires) {
			http.Redirect(w, r, "/login", http.StatusFound)
		}

		return next(w, r)
	}
}
