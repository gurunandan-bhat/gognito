package service

import (
	"fmt"
	"gognito/lib/config"
	"gognito/lib/model"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	mysqlstore "github.com/danielepintore/gorilla-sessions-mysql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/csrf"
)

type Service struct {
	Config       *config.Config
	Model        *model.Model
	Muxer        *chi.Mux
	SessionStore *mysqlstore.MysqlStore
	Template     map[string]*template.Template
	Logger       *slog.Logger
}

func NewService(cfg *config.Config) (*Service, error) {

	mux := chi.NewRouter()

	// force a redirect to https:// in production
	if cfg.InProduction {
		mux.Use(middleware.SetHeader(
			"Strict-Transport-Security",
			"max-age=63072000; includeSubDomains",
		))
	}

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Recoverer)

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	csrfMiddleware := csrf.Protect(
		[]byte(cfg.Security.CSRFKey),
		csrf.Secure(cfg.InProduction),
		csrf.SameSite(csrf.SameSiteStrictMode),
	)
	mux.Use(csrfMiddleware)

	model, err := model.NewModel(cfg)
	if err != nil {
		log.Fatalf("error initializing database connection: %s", err)
	}

	dbStore, err := model.NewDbSessionStore(cfg)
	if err != nil {
		log.Fatalf("error initializing db store: %s", err)
	}

	// Static file handler
	filesDir := http.Dir(filepath.Join(cfg.AppRoot, "assets"))
	fs := http.FileServer(filesDir)
	mux.Handle("/assets/*", http.StripPrefix("/assets", fs))

	template, err := newTemplateCache(filepath.Join(cfg.AppRoot, "templates"))
	if err != nil {
		log.Fatalf("Cannot build template cache: %s", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	s := &Service{
		Config:       cfg,
		SessionStore: dbStore,
		Model:        model,
		Muxer:        mux,
		Template:     template,
		Logger:       logger,
	}

	s.setRoutes()

	return s, nil
}

func (s *Service) setRoutes() {

	s.Muxer.Method(http.MethodGet, "/", serviceHandler(s.index))
	s.Muxer.Method(http.MethodGet, "/about", serviceHandler((s.about)))
	s.Muxer.Method(http.MethodGet, "/action", serviceHandler((s.action)))
	s.Muxer.Method(http.MethodGet, "/another-action", serviceHandler((s.anotherAction)))
	s.Muxer.Method(http.MethodGet, "/aws", serviceHandler(s.initAWS))
	s.Muxer.Method(http.MethodGet, "/login", serviceHandler(s.login))
	s.Muxer.Method(http.MethodPost, "/auth-status", serviceHandler(s.handleCallback))
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
