package main

import (
	"educore/internal"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return
	}

	if err := db.AutoMigrate(internal.Models...); err != nil {
		slog.Error("failed to migrate database", "error", err)
		return
	}
	slog.Info("database migration complete")

	mux := http.NewServeMux()

	mux.Handle("/health", internal.HandlerFunc(func(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
		internal.WriteData(w, "success", "service is healthy", nil)
		return nil
	}))
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
