package service

import (
	"encoding/gob"
	"fmt"
	"gognito/lib/config"
	"gognito/lib/model"
	"net/http"
	"time"

	mysqlstore "github.com/danielepintore/gorilla-sessions-mysql"
)

type AuthInfo struct {
	Email        string
	Roles        []string
	AccessToken  string
	RefreshToken string
	Expires      time.Time
}

func newDbSessionStore(cfg *config.Config, m *model.Model) (*mysqlstore.MysqlStore, error) {

	keyPair := mysqlstore.KeyPair{
		AuthenticationKey: []byte(cfg.Session.AuthenticationKey),
		EncryptionKey:     []byte(cfg.Session.EncryptionKey),
	}

	// register so gorilla can save complex data structures
	gob.Register(&AuthInfo{})

	cleanupAfter := 60 * time.Minute
	return mysqlstore.NewMysqlStore(
		m.DbHandle.DB,
		"mdbsession",
		[]mysqlstore.KeyPair{keyPair},
		mysqlstore.WithCleanupInterval(cleanupAfter),
		mysqlstore.WithHttpOnly(true),
		mysqlstore.WithSameSite(http.SameSiteLaxMode),
		mysqlstore.WithMaxAge(cfg.Session.MaxAgeHours*3600),
		mysqlstore.WithSecure(cfg.InProduction),
	)
}

func (s *Service) getSessionVar(r *http.Request, name string) (any, error) {

	sessionName := s.Config.Session.Name
	session, err := s.SessionStore.Get(r, sessionName)
	if err != nil {
		return nil, fmt.Errorf("error fetching session %s: %w", sessionName, err)
	}
	return session.Values[name], nil
}

func (s *Service) setSessionVar(r *http.Request, w http.ResponseWriter, name string, value any) error {

	sessionName := s.Config.Session.Name
	session, err := s.SessionStore.Get(r, sessionName)
	if err != nil {
		return fmt.Errorf("error fetching session %s: %w", sessionName, err)
	}

	session.Values[name] = value
	return session.Save(r, w)
}
