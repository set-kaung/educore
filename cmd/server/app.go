package main

import (
	"educore/internal"
	"educore/internal/auth"
	"educore/internal/auth/authenticators"
	"educore/internal/course"
	"educore/internal/semestercourse"
	"educore/internal/student"
	"educore/internal/textbook"
	"log/slog"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	DSN         string
	JWTSecret   string
	AuthType    string
	OLUserAgent string
}

func Setup(conf Config) (http.Handler, error) {
	db, err := ConnectAndMigrateDatabase(conf.DSN)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	chain := NewRouteChainer(RequestLogMiddleWare)

	mux.Handle("GET /health", chain.Chain(http.HandlerFunc(HealthCheck)))

	protected := chain.Append(auth.NewJWTAuth(conf.JWTSecret).Middleware())
	professorOnly := protected.Append(RoleRequired("professor", "admin"))

	var authenticator auth.Authenticator
	switch conf.AuthType {
	case "ad":
		authenticator = authenticators.NewADAuthenticator("", "", "")
	default:
		authenticator = authenticators.NewMockAuthenticator()
	}

	sh := student.NewStudentHandler(db)
	ah := auth.AuthHandler{Authenticator: authenticator, JWTSecret: conf.JWTSecret}
	th := textbook.NewTextbookHandler(conf.OLUserAgent)
	ch := course.NewCourseHandler(db)
	sch := semestercourse.NewHandler(db)

	mux.Handle("POST /login", chain.Chain(internal.HandlerFunc(ah.HandleLogin)))
	mux.Handle("GET /student", protected.Chain(internal.HandlerFunc(sh.HandleGetAllStudents)))
	mux.Handle("GET /textbooks", chain.Chain(internal.HandlerFunc(th.HandleSearch)))
	mux.Handle("GET /courses", protected.Chain(internal.HandlerFunc(ch.HandleGetBySemester)))
	mux.Handle("POST /courses", professorOnly.Chain(internal.HandlerFunc(ch.HandleCreateCourse)))
	mux.Handle("POST /semester-courses", professorOnly.Chain(internal.HandlerFunc(sch.HandleCreate)))

	return mux, nil
}

// HealthCheck godoc
// @Summary      Health check
// @Success      200  {object}  internal.ResponseBody
// @Router       /health [get]
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
