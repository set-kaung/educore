package main

import (
	"educore/internal"
	"educore/internal/auth"
	"educore/internal/auth/authenticators"
	"educore/internal/course"
	"educore/internal/publicapi"
	"educore/internal/semestercourse"
	"educore/internal/student"
	"educore/internal/textbook"
	"educore/internal/web"
	"educore/internal/web/pages"
	assets "educore/web"
	"io/fs"
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

	renderer, err := web.NewRenderer(assets.Files)
	if err != nil {
		return nil, err
	}

	staticRoot, err := fs.Sub(assets.Files, "static")
	if err != nil {
		return nil, err
	}

	mux.Handle("GET /health", chain.Chain(http.HandlerFunc(HealthCheck)))

	jwtAuth := auth.NewJWTAuth(conf.JWTSecret)
	csrf := web.NewCSRF(conf.JWTSecret)

	protected := chain.Append(jwtAuth.Middleware())
	professorOnly := protected.Append(RoleRequired("professor", "admin"))
	adminOnly := protected.Append(RoleRequired("admin"))
	apiKeyOnly := chain.Append(publicapi.APIKeyMiddleware(db))

	guestUI := chain.Append(csrf.Middleware).Append(jwtAuth.OptionalMiddleware())
	authedUI := guestUI.Append(jwtAuth.PageMiddleware("/login"))

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
	ph := publicapi.NewHandler(db)

	lh := pages.NewLoginHandler(authenticator, conf.JWTSecret, renderer)
	sph := pages.NewStudentsHandler(db, renderer)

	mux.Handle("GET /{$}", guestUI.Chain(internal.HandlerFunc(lh.Home)))
	mux.Handle("GET /static/", chain.Chain(http.StripPrefix("/static/", http.FileServerFS(staticRoot))))
	mux.Handle("GET /login", guestUI.Append(web.RedirectIfAuthenticated("/students")).Chain(internal.HandlerFunc(lh.Show)))
	mux.Handle("POST /session", guestUI.Chain(internal.HandlerFunc(lh.Submit)))
	mux.Handle("POST /logout", chain.Chain(internal.HandlerFunc(lh.Logout)))
	mux.Handle("GET /students", authedUI.Chain(renderer.H(sph.Show)))
	mux.Handle("GET /students/search", authedUI.Chain(renderer.H(sph.Search)))

	mux.Handle("POST /login", chain.Chain(internal.HandlerFunc(ah.HandleLogin)))
	mux.Handle("GET /student", protected.Chain(internal.HandlerFunc(sh.HandleGetAllStudents)))
	mux.Handle("GET /textbooks", chain.Chain(internal.HandlerFunc(th.HandleSearch)))
	mux.Handle("GET /semester-courses", protected.Chain(internal.HandlerFunc(sch.HandleGetBySemester)))
	mux.Handle("POST /courses", professorOnly.Chain(internal.HandlerFunc(ch.HandleCreateCourse)))
	mux.Handle("POST /semester-courses", professorOnly.Chain(internal.HandlerFunc(sch.HandleCreate)))
	mux.Handle("GET /public/students/{student_id}/departments/{department_name}", apiKeyOnly.Chain(internal.HandlerFunc(ph.HandleVerifyDepartmentEnrollment)))
	mux.Handle("POST /admin/api-keys", adminOnly.Chain(internal.HandlerFunc(ph.HandleGrantKey)))
	mux.Handle("DELETE /admin/api-keys/{key}", adminOnly.Chain(internal.HandlerFunc(ph.HandleRevokeKey)))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, pattern := mux.Handler(r); pattern == "" {
			renderer.Error(w, http.StatusNotFound)
			return
		}
		mux.ServeHTTP(w, r)
	}), nil
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
