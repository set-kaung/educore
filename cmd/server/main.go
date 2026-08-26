package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

//	@title			EduCore API
//	@version		1.0
//	@description	Course-registration backend for CSX 4110.

//	@host		localhost:8080
//	@schemes	http

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						x-api-key
func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		slog.Error("failed to start server",
			"reason", "DATABASE_URL environment variable is empty",
		)
		return
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		slog.Error("failed to start server",
			"reason", "JWT_SECRET environment variable is empty",
		)
		return
	}

	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = "web/public"
	}

	conf := Config{
		DSN:         dsn,
		JWTSecret:   jwtSecret,
		AuthType:    os.Getenv("AUTH_TYPE"),
		OLUserAgent: os.Getenv("OPENLIBRARY_USER_AGENT"),
		StaticDir:   staticDir,
	}

	mux, err := Setup(conf)
	if err != nil {
		slog.Error("failed to create server", "reason", err.Error())
		return
	}

	if abs, err := filepath.Abs(staticDir); err == nil {
		slog.Info("serving static files", "dir", abs)
	}

	slog.Info("starting HTTP server",
		"port", port,
	)

	err = http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
	if err != nil {
		slog.Error("HTTP server stopped",
			"error", err,
		)
	}
}
