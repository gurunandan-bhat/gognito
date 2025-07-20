package service

import (
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

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	mux.Use(newSlogger(cfg, logger))

	model, err := model.NewModel(cfg)
	if err != nil {
		log.Fatalf("error initializing model: %s", err)
	}

	sessionStore, err := newDbSessionStore(cfg, model)
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

	s := &Service{
		Config:       cfg,
		SessionStore: sessionStore,
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
	s.Muxer.Method(http.MethodGet, "/auth-status", serviceHandler(s.handleCallback))
	s.Muxer.Method(http.MethodGet, "/form", serviceHandler(s.validateAuth(s.form)))
}
