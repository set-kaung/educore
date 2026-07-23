package main

import (
	"educore/internal"
	"log/slog"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	DSN       string
	JWTSecret string
}

func Setup(conf Config) (http.Handler, error) {
	db, err := ConnectAndMigrateDatabase(conf.DSN)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	chain := NewRouteChainer(RequestLogMiddleWare)

	mux.Handle("GET /health", chain.Chain(http.HandlerFunc(HealthCheck)))

	protected := chain.Append(NewJWTAuth(conf.JWTSecret).Middleware())

	sh := StudentHandler{db: db}
	ah := AuthHandler{db: db, jwtSecret: conf.JWTSecret}

	mux.Handle("POST /login", chain.Chain(internal.HandlerFunc(ah.HandleLogin)))
	mux.Handle("GET /student", protected.Chain(internal.HandlerFunc(sh.HandleGetAllStudents)))

	return mux, nil
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	internal.WriteData(w, "success", "service is healthy", nil)
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
