package web

import (
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// StaticFiles serves HTML pages and assets from a directory on disk.
type StaticFiles struct {
	dir string
}

func NewStaticFiles(dir string) (*StaticFiles, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fs.ErrInvalid
	}
	return &StaticFiles{dir: dir}, nil
}

// Page returns a handler that serves the given file from the static root,
// used to map URLs like /login to login.html.
func (s *StaticFiles) Page(name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.ServePage(w, r, name, http.StatusOK)
	})
}

// ServePage writes the given file with an explicit HTTP status code,
// used e.g. to serve 404.html with status 404.
func (s *StaticFiles) ServePage(w http.ResponseWriter, r *http.Request, name string, status int) {
	full, err := s.resolve(name)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data, err := os.ReadFile(full)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}

// Assets serves files under the named directory of the static root
func (s *StaticFiles) Assets(mount string) http.Handler {
	mount = strings.Trim(mount, "/")
	prefix := mount + "/"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.Trim(path.Clean(r.URL.Path), "/")
		if name == mount || !strings.HasPrefix(name, prefix) {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Cache-Control", "public, max-age=300")
		s.serveFile(w, r, name)
	})
}

func (s *StaticFiles) serveFile(w http.ResponseWriter, r *http.Request, name string) {
	full, err := s.resolve(name)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	file, err := os.Open(full)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil || !info.Mode().IsRegular() {
		http.NotFound(w, r)
		return
	}

	http.ServeContent(w, r, filepath.Base(full), info.ModTime(), file)
}

func (s *StaticFiles) resolve(name string) (string, error) {
	clean := path.Clean("/" + name)
	full := filepath.Join(s.dir, filepath.FromSlash(clean))

	rootAbs, err := filepath.Abs(s.dir)
	if err != nil {
		return "", err
	}
	rootReal, err := filepath.EvalSymlinks(rootAbs)
	if err != nil {
		return "", err
	}
	fileReal, err := filepath.EvalSymlinks(full)
	if err != nil {
		return "", err
	}
	if fileReal != rootReal && !strings.HasPrefix(fileReal, rootReal+string(filepath.Separator)) {
		return "", fs.ErrPermission
	}

	info, err := os.Stat(fileReal)
	if err != nil || !info.Mode().IsRegular() {
		return "", fs.ErrNotExist
	}
	return fileReal, nil
}
