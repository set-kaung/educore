package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

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

	conf := Config{DSN: dsn, JWTSecret: jwtSecret, AuthType: os.Getenv("AUTH_TYPE")}

	mux, err := Setup(conf)
	if err != nil {
		slog.Error("failed to create server", "reason", err.Error())
		return
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
