package service

import (
	"gognito/lib/config"
	"gognito/lib/model"
	"net/http"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
)

type AuthInfo struct {
	Name      string
	Email     string
	Roles     []string
	Expires   time.Time
	LogoutURL string
}

func newDbSessionStore(cfg *config.Config, m *model.Model) *scs.SessionManager {

	sessMgr := scs.New()
	sessMgr.Store = mysqlstore.New(m.DbHandle.DB)

	sessMgr.Lifetime = 1 * time.Hour
	sessMgr.Cookie.Name = cfg.Session.CookieName
	sessMgr.Cookie.HttpOnly = true
	sessMgr.Cookie.Path = "/"
	sessMgr.Cookie.Persist = true
	sessMgr.Cookie.SameSite = http.SameSiteStrictMode
	sessMgr.Cookie.Secure = cfg.InProduction
	sessMgr.Cookie.Partitioned = false

	return sessMgr
}

func (s *Service) setSessionVar(r *http.Request, key string, value any) error {

	return nil
}

func (s *Service) getSessionVar(r *http.Request, key string) (any, error) {

	return nil, nil
}
