package main

import (
	"educore/internal"
	"educore/internal/auth"
	"educore/internal/auth/authenticators"
	"educore/internal/course"
	"educore/internal/professor"
	"educore/internal/publicapi"
	"educore/internal/semestercourse"
	"educore/internal/student"
	"educore/internal/textbook"
	"educore/internal/web"
	"log/slog"
	"net/http"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	DSN         string
	JWTSecret   string
	AuthType    string
	OLUserAgent string
	StaticDir   string
}

func Setup(conf Config) (http.Handler, error) {
	db, err := ConnectAndMigrateDatabase(conf.DSN)
	if err != nil {
		return nil, err
	}

	staticFiles, err := web.NewStaticFiles(conf.StaticDir)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	chain := NewRouteChainer(RequestLogMiddleWare)
	api := chain.Append(web.SameOrigin)

	jwtAuth := auth.NewJWTAuth(conf.JWTSecret)
	protected := api.Append(jwtAuth.Middleware())
	professorOnly := protected.Append(RoleRequired("professor", "admin"))
	adminOnly := protected.Append(RoleRequired("admin"))
	apiKeyOnly := chain.Append(publicapi.APIKeyMiddleware(db))

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
	prh := professor.NewHandler(db)

	sch := semestercourse.NewHandler(db)
	ph := publicapi.NewHandler(db)

	mux.Handle("GET /health", chain.Chain(http.HandlerFunc(HealthCheck)))

	mux.Handle("GET /{$}", staticFiles.Page("index.html"))
	mux.Handle("GET /login", staticFiles.Page("login.html"))
	mux.Handle("GET /students", staticFiles.Page("students.html"))
	mux.Handle("GET /my-courses", staticFiles.Page("my-courses.html"))
	mux.Handle("GET /course-offerings", staticFiles.Page("course-offerings.html"))
	mux.Handle("GET /add-course", staticFiles.Page("add-course.html"))
	mux.Handle("GET /add-offering", staticFiles.Page("add-offering.html"))
	mux.Handle("GET /css/", staticFiles.Assets("css"))
	mux.Handle("GET /js/", staticFiles.Assets("js"))

	mux.Handle("POST /api/login", api.Chain(internal.HandlerFunc(ah.HandleLogin)))
	mux.Handle("POST /api/logout", api.Chain(internal.HandlerFunc(ah.HandleLogout)))
	mux.Handle("GET /api/session", protected.Chain(internal.HandlerFunc(ah.HandleSession)))
	mux.Handle("GET /api/my/courses", protected.Chain(internal.HandlerFunc(sch.HandleGetEnrolled)))

	mux.Handle("GET /api/students", professorOnly.Chain(internal.HandlerFunc(sh.HandleGetAllStudents)))
	mux.Handle("GET /api/courses", protected.Chain(internal.HandlerFunc(ch.HandleListCourses)))
	mux.Handle("POST /api/courses", professorOnly.Chain(internal.HandlerFunc(ch.HandleCreateCourse)))
	mux.Handle("GET /api/professors", protected.Chain(internal.HandlerFunc(prh.HandleListProfessors)))
	mux.Handle("GET /api/semesters", protected.Chain(internal.HandlerFunc(sch.HandleListSemesters)))
	mux.Handle("GET /api/semester-courses", protected.Chain(internal.HandlerFunc(sch.HandleGetBySemester)))
	mux.Handle("POST /api/semester-courses", professorOnly.Chain(internal.HandlerFunc(sch.HandleCreate)))

	mux.Handle("GET /api/textbooks", chain.Chain(internal.HandlerFunc(th.HandleSearch)))

	mux.Handle("GET /public/students/{student_id}/departments/{department_name}", apiKeyOnly.Chain(internal.HandlerFunc(ph.HandleVerifyDepartmentEnrollment)))
	mux.Handle("POST /admin/api-keys", adminOnly.Chain(internal.HandlerFunc(ph.HandleGrantKey)))
	mux.Handle("DELETE /admin/api-keys/{key}", adminOnly.Chain(internal.HandlerFunc(ph.HandleRevokeKey)))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, pattern := mux.Handler(r); pattern == "" {
			notFound(w, r, staticFiles)
			return
		}
		mux.ServeHTTP(w, r)
	}), nil
}

func notFound(w http.ResponseWriter, r *http.Request, static *web.StaticFiles) {
	slog.Warn("no route matched", "method", r.Method, "path", r.URL.Path)
	if isAPIRequest(r) {
		internal.WriteError(w, http.StatusNotFound, "route not found", nil)
		return
	}
	static.ServePage(w, r, "404.html", http.StatusNotFound)
}

func isAPIRequest(r *http.Request) bool {
	prefixes := []string{"/api/", "/public/", "/admin/", "/health"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(r.URL.Path, prefix) || r.URL.Path == strings.TrimSuffix(prefix, "/") {
			return true
		}
	}
	return false
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
