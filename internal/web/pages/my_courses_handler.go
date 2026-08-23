package pages

import (
	"net/http"

	"educore/internal"
	"educore/internal/auth"
	"educore/internal/semestercourse"
	"educore/internal/web"

	"gorm.io/gorm"
)

type MyCoursesHandler struct {
	DB     *gorm.DB
	Render *web.Renderer
}

func NewMyCoursesHandler(db *gorm.DB, renderer *web.Renderer) *MyCoursesHandler {
	return &MyCoursesHandler{DB: db, Render: renderer}
}

func (h *MyCoursesHandler) Show(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	claims := auth.GetClaims(r)

	courses, err := semestercourse.GetEnrolledCourses(h.DB, claims.UserID)
	if err != nil {
		return &internal.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Message:    "could not get enrolled courses",
			Err:        err,
		}
	}

	data := web.BaseData(r, "My Courses")
	data.Data = courses
	h.Render.Page(w, http.StatusOK, "my_courses", data)
	return nil
}
