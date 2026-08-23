package web

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"path"
	"strings"

	"educore/internal"
	"educore/internal/auth"
)

type PageData struct {
	Title         string
	Authenticated bool
	Role          string
	CSRFToken     string
	Flash         string
	Error         string
	Data          any
}

type ErrorView struct {
	Code    int
	Status  string
	Message string
}

type Renderer struct {
	pages     map[string]*template.Template
	fragments map[string]*template.Template
}

func NewRenderer(files fs.FS) (*Renderer, error) {
	root, err := fs.Sub(files, "templates")
	if err != nil {
		return nil, fmt.Errorf("locate templates: %w", err)
	}

	partials, err := fs.Glob(root, "partials/*.html")
	if err != nil {
		return nil, fmt.Errorf("list partials: %w", err)
	}
	pages, err := fs.Glob(root, "pages/*.html")
	if err != nil {
		return nil, fmt.Errorf("list pages: %w", err)
	}

	r := &Renderer{
		pages:     make(map[string]*template.Template, len(pages)),
		fragments: make(map[string]*template.Template, len(partials)),
	}

	for _, page := range pages {
		set := append([]string{"layouts/base.html"}, partials...)
		set = append(set, page)
		tmpl, err := template.ParseFS(root, set...)
		if err != nil {
			return nil, fmt.Errorf("parse page %s: %w", page, err)
		}
		r.pages[templateName(page)] = tmpl
	}

	for _, partial := range partials {
		tmpl, err := template.ParseFS(root, partial)
		if err != nil {
			return nil, fmt.Errorf("parse fragment %s: %w", partial, err)
		}
		r.fragments[templateName(partial)] = tmpl
	}

	return r, nil
}

func (r *Renderer) Page(w http.ResponseWriter, status int, name string, data PageData) {
	tmpl, ok := r.pages[name]
	if !ok {
		slog.Error("unknown page template", "page", name)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	r.execute(w, status, name, tmpl, data)
}

func (r *Renderer) Fragment(w http.ResponseWriter, status int, name string, data any) {
	tmpl, ok := r.fragments[name]
	if !ok {
		slog.Error("unknown fragment template", "fragment", name)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	r.executeNamed(w, status, name, tmpl, data)
}

func (r *Renderer) Error(w http.ResponseWriter, status int) {
	statusText := http.StatusText(status)
	r.Page(w, status, "error", PageData{
		Title: statusText,
		Data: ErrorView{
			Code:    status,
			Status:  statusText,
			Message: friendlyMessage(status),
		},
	})
}

func (r *Renderer) H(h func(http.ResponseWriter, *http.Request) *internal.HTTPError) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if herr := h(w, req); herr != nil {
			if herr.StatusCode >= 500 {
				slog.Error("internal error", "message", herr.Message, "error", herr.Err)
			}
			r.Error(w, herr.StatusCode)
		}
	})
}

func (r *Renderer) execute(w http.ResponseWriter, status int, name string, tmpl *template.Template, data any) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		slog.Error("execute template", "template", name, "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(buf.Bytes())
}

func (r *Renderer) executeNamed(w http.ResponseWriter, status int, name string, tmpl *template.Template, data any) {
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		slog.Error("execute template", "template", name, "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(buf.Bytes())
}

func BaseData(r *http.Request, title string) PageData {
	data := PageData{Title: title, CSRFToken: CSRFToken(r)}
	if claims, ok := auth.TryGetClaims(r); ok {
		data.Authenticated = true
		data.Role = claims.Role
	}
	return data
}

func templateName(file string) string {
	return strings.TrimSuffix(path.Base(file), ".html")
}

func friendlyMessage(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "The request could not be processed."
	case http.StatusUnauthorized:
		return "Please log in to continue."
	case http.StatusForbidden:
		return "You do not have permission to do that."
	case http.StatusNotFound:
		return "The page you requested does not exist."
	default:
		return "Something went wrong on our side. Please try again."
	}
}
