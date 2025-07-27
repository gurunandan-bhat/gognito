package service

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

func newTemplateCache(templateRoot string) (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}
	pages, err := filepath.Glob(templateRoot + "/pages/*.go.html")
	if err != nil {
		return nil, fmt.Errorf("error generating list of templates in pages: %w", err)
	}

	for _, page := range pages {

		name := filepath.Base(page)
		files := []string{
			templateRoot + "/common/base.go.html",
			templateRoot + "/common/head.go.html",
			templateRoot + "/common/top-menu.go.html",
			templateRoot + "/common/footer.go.html",
			templateRoot + "/common/js-includes.go.html",
			templateRoot + "/common/auth-state.go.html",
			page,
		}
		tSet, err := template.ParseFiles(files...)
		if err != nil {
			return nil, fmt.Errorf("error creating template set for %s: %w", page, err)
		}
		tSet, err = tSet.ParseGlob(templateRoot + "/includes/*.go.html")
		if err != nil {
			return nil, fmt.Errorf("error parsing included templates set for %s: %w", page, err)
		}

		cache[name] = tSet
	}

	return cache, nil
}

func (s *Service) render(w http.ResponseWriter, r *http.Request, template string, data any, status int) error {

	// Check whether that template exists in the cache
	tmpl, ok := s.Template[template]
	if !ok {
		return fmt.Errorf("template %s is not available in the cache", template)
	}

	authState, err := s.getSessionVar(r, "authInfo")
	if err != nil {
		return fmt.Errorf("render error fetching auth state: %w", err)
	}

	if authState == nil {
		authState = AuthInfo{}
	}

	data = struct {
		Common   any
		PageData any
	}{
		Common:   authState,
		PageData: data,
	}

	var b bytes.Buffer
	if err := tmpl.ExecuteTemplate(&b, "base", data); err != nil {
		return fmt.Errorf("error executing template %s: %w", template, err)
	}

	w.WriteHeader(status)
	w.Header().Add("Content-Type", "text/html")
	if _, err := w.Write(b.Bytes()); err != nil {
		return fmt.Errorf("error writing response n render: %w", err)
	}

	return nil
}
