package main

import (
	"context"
	"educore/internal"
	"educore/internal/auth"
	"educore/internal/course"
	"educore/internal/department"
	"educore/internal/professor"
	"educore/internal/publicapi"
	"educore/internal/recommendedbook"
	"educore/internal/semestercourse"
	"educore/internal/semesters"
	"educore/internal/student"
	"educore/internal/textbook"
	"educore/internal/web"
	"log/slog"
	"net/http"
	"strings"
	"time"

	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	DSN            string
	JWTSecret      string
	OLUserAgent    string
	StaticDir      string
	BasePath       string
	ADClientID     string
	ADTenantID     string
	ADClientSecret string
	ADRedirectURI  string
}

func Setup(conf Config) (http.Handler, error) {
	db, err := ConnectAndMigrateDatabase(conf.DSN)
	if err != nil {
		return nil, err
	}

	staticFiles, err := web.NewStaticFiles(conf.StaticDir, conf.BasePath)
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
	studentOnly := protected.Append(RoleRequired("student"))
	apiKeyOnly := chain.Append(publicapi.APIKeyMiddleware(db))

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	provider, err := auth.NewOIDCProvider(ctx, conf.ADTenantID, conf.ADClientID, conf.ADClientSecret, conf.ADRedirectURI)
	cancel()
	if err != nil {
		return nil, err
	}
	oidcHandler := &auth.OIDCHandler{Provider: provider, DB: db, JWTSecret: conf.JWTSecret, BasePath: conf.BasePath}

	sh := student.NewStudentHandler(db)
	ah := auth.AuthHandler{BasePath: conf.BasePath}
	th := textbook.NewTextbookHandler(conf.OLUserAgent)
	ch := course.NewCourseHandler(db)
	prh := professor.NewHandler(db)

	sch := semestercourse.NewHandler(db)
	rbh := recommendedbook.NewHandler(db)
	ph := publicapi.NewHandler(db)
	dh := department.NewHandler(db)
	smh := semesters.NewHandler(db)

	mux.Handle("GET /health", chain.Chain(http.HandlerFunc(HealthCheck)))

	mux.Handle("GET /favicon.ico", staticFiles.RootFile("favicon.ico"))
	mux.Handle("GET /logo.png", staticFiles.RootFile("logo.png"))
	mux.Handle("GET /MS_Logo.jpg", staticFiles.RootFile("MS_Logo.jpg"))
	mux.Handle("GET /{$}", staticFiles.Page("index.html"))
	mux.Handle("GET /login", staticFiles.Page("login.html"))
	mux.Handle("GET /students", staticFiles.Page("students.html"))
	mux.Handle("GET /my-courses", staticFiles.Page("my-courses.html"))
	mux.Handle("GET /course-offerings", staticFiles.Page("course-offerings.html"))
	mux.Handle("GET /add-course", staticFiles.Page("add-course.html"))
	mux.Handle("GET /add-offering", staticFiles.Page("add-offering.html"))
	mux.Handle("GET /recommended-books", staticFiles.Page("recommended-books.html"))
	mux.Handle("GET /offering-detail", staticFiles.Page("offering-detail.html"))
	mux.Handle("GET /css/", staticFiles.Assets("css"))
	mux.Handle("GET /js/", staticFiles.Assets("js"))

	mux.Handle("POST /api/logout", api.Chain(internal.HandlerFunc(ah.HandleLogout)))
	mux.Handle("GET /api/session", protected.Chain(internal.HandlerFunc(ah.HandleSession)))
	mux.Handle("GET /api/my/courses", protected.Chain(internal.HandlerFunc(sch.HandleGetEnrolled)))
	mux.Handle("GET /api/my/taught-courses", professorOnly.Chain(internal.HandlerFunc(sch.HandleGetTaught)))
	mux.Handle("GET /api/departments", chain.Chain(internal.HandlerFunc(dh.HandleListDepartments)))
	mux.Handle("GET /api/semesters", protected.Chain(internal.HandlerFunc(smh.HandleList)))

	mux.Handle("GET /setup", staticFiles.Page("setup.html"))
	mux.Handle("GET /auth/login", chain.Chain(http.HandlerFunc(oidcHandler.HandleADLogin)))
	mux.Handle("GET /auth/callback", chain.Chain(http.HandlerFunc(oidcHandler.HandleADCallback)))
	mux.Handle("GET /api/setup/context", api.Append(jwtAuth.SetupMiddleware()).Chain(internal.HandlerFunc(oidcHandler.HandleSetupContext)))
	mux.Handle("POST /api/setup", api.Append(jwtAuth.SetupMiddleware()).Chain(internal.HandlerFunc(oidcHandler.HandleSetup)))

	mux.Handle("GET /api/students", professorOnly.Chain(internal.HandlerFunc(sh.HandleGetAllStudents)))
	mux.Handle("GET /api/courses", protected.Chain(internal.HandlerFunc(ch.HandleListCourses)))
	mux.Handle("POST /api/courses", professorOnly.Chain(internal.HandlerFunc(ch.HandleCreateCourse)))
	mux.Handle("GET /api/professors", protected.Chain(internal.HandlerFunc(prh.HandleListProfessors)))
	mux.Handle("GET /api/semester-courses", protected.Chain(internal.HandlerFunc(sch.HandleGetBySemester)))
	mux.Handle("GET /api/semester-courses/{id}", protected.Chain(internal.HandlerFunc(sch.HandleGetDetail)))
	mux.Handle("POST /api/semester-courses", professorOnly.Chain(internal.HandlerFunc(sch.HandleCreate)))
	mux.Handle("POST /api/semester-courses/{id}/enroll", studentOnly.Chain(internal.HandlerFunc(sch.HandleEnroll)))

	mux.Handle("GET /api/semester-courses/{id}/books", protected.Chain(internal.HandlerFunc(rbh.HandleList)))
	mux.Handle("POST /api/semester-courses/{id}/books", professorOnly.Chain(internal.HandlerFunc(rbh.HandleAdd)))
	mux.Handle("DELETE /api/books/{id}", professorOnly.Chain(internal.HandlerFunc(rbh.HandleRemove)))

	mux.Handle("GET /api/textbooks", chain.Chain(internal.HandlerFunc(th.HandleSearch)))

	mux.Handle("GET /public/students/{student_id}/departments/{department_name}", apiKeyOnly.Chain(internal.HandlerFunc(ph.HandleVerifyDepartmentEnrollment)))
	mux.Handle("POST /admin/api-keys", adminOnly.Chain(internal.HandlerFunc(ph.HandleGrantKey)))
	mux.Handle("DELETE /admin/api-keys/{key}", adminOnly.Chain(internal.HandlerFunc(ph.HandleRevokeKey)))

	base := strings.Trim(conf.BasePath, "/")
	prefix := ""
	if base != "" {
		prefix = "/" + base
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if prefix != "" {
			if !strings.HasPrefix(r.URL.Path, prefix+"/") {
				if r.URL.Path == prefix {
					http.Redirect(w, r, prefix+"/", http.StatusMovedPermanently)
					return
				}
				notFound(w, r, staticFiles)
				return
			}
			r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
			r.URL.RawPath = ""
		}
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
	name := strings.Trim(r.URL.Path, "/")
	if name != "" && static.ServeAsset(w, r, name) {
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
	db, err := gorm.Open(gormmysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(internal.Models...); err != nil {
		return nil, err
	}
	slog.Info("database", "status", "migrate and connected")
	return db, nil
}
