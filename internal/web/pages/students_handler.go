package pages

import (
	"net/http"
	"strings"

	"educore/internal"
	"educore/internal/auth"
	"educore/internal/student"
	"educore/internal/web"

	"gorm.io/gorm"
)

type StudentsView struct {
	Students []student.StudentData
	Query    string
}

type StudentsHandler struct {
	DB     *gorm.DB
	Render *web.Renderer
}

func NewStudentsHandler(db *gorm.DB, renderer *web.Renderer) *StudentsHandler {
	return &StudentsHandler{DB: db, Render: renderer}
}

func (h *StudentsHandler) Show(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	if claims, ok := auth.TryGetClaims(r); ok && claims.Role == "student" {
		http.Redirect(w, r, "/my-courses", http.StatusSeeOther)
		return nil
	}

	students, query, herr := h.list(r)
	if herr != nil {
		return herr
	}

	data := web.BaseData(r, "Students")
	data.Data = StudentsView{Students: students, Query: query}
	h.Render.Page(w, http.StatusOK, "students", data)
	return nil
}

func (h *StudentsHandler) Search(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	students, _, herr := h.list(r)
	if herr != nil {
		return herr
	}

	h.Render.Fragment(w, http.StatusOK, "student_rows", students)
	return nil
}

func (h *StudentsHandler) list(r *http.Request) ([]student.StudentData, string, *internal.HTTPError) {
	all, err := student.GetAllStudents(h.DB)
	if err != nil {
		return nil, "", &internal.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Message:    "could not get student data",
			Err:        err,
		}
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	return filterStudents(all, query), query, nil
}

func filterStudents(students []student.StudentData, query string) []student.StudentData {
	if query == "" {
		return students
	}

	needle := strings.ToLower(query)
	filtered := make([]student.StudentData, 0, len(students))
	for _, s := range students {
		if containsAny(needle, s.Username, s.StudentID, s.DepartmentName) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

func containsAny(needle string, haystacks ...string) bool {
	for _, hay := range haystacks {
		if strings.Contains(strings.ToLower(hay), needle) {
			return true
		}
	}
	return false
}
