package service

import "net/http"

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
