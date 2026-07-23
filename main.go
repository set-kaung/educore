package main

import (
	"educore/internal"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		slog.Error("failed to start server",
			"reason", "PORT environment variable is empty",
		)
		return
	}
	mux := http.NewServeMux()

	mux.Handle("/", internal.HandlerFunc(func(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
		internal.WriteData(w, "success", "hello", nil)
		return nil
	}))
	slog.Info("starting HTTP server",
		"port", port,
	)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
	if err != nil {
		slog.Error("HTTP server stopped",
			"error", err,
		)
	}
}
