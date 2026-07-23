package main

import (
	"educore/internal"
	"log/slog"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	DSN string
}

func Setup(conf Config) (*http.ServeMux, error) {
	db, err := ConnectAndMigrateDatabase(conf.DSN)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	mux.Handle("GET /health", internal.HandlerFunc(func(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
		internal.WriteData(w, "success", "service is healthy", nil)
		return nil
	}))

	sh := StudentHandler{db: db}

	mux.Handle("GET /student", internal.HandlerFunc(sh.HandleGetAllStudents))

	return mux, nil
}

func ConnectAndMigrateDatabase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {

		return nil, err
	}
	if err := db.AutoMigrate(internal.Models...); err != nil {
		return nil, err
	}
	slog.Info("database", "status", "migrate and connected")

	return db, nil
}
