package internal

import (
	"log/slog"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request) *HTTPError

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(w, r)
	if err == nil {
		return
	}

	if err.StatusCode >= 500 {
		slog.Error("internal error", "error message", err.Message)
		slog.Info("support ticket created")
		WriteServerError(w, nil)
		return
	}

	WriteError(w, err.StatusCode, err.Message, nil)
}
