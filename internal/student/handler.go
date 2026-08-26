package student

import (
	"educore/internal"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type StudentHandler struct {
	db *gorm.DB
}

func NewStudentHandler(db *gorm.DB) *StudentHandler {
	return &StudentHandler{db: db}
}

// HandleGetAllStudents godoc
// @Summary      List students
// @Description  Get all students with department names, optionally filtered by name, ID or department. Professor or admin only.
// @Tags         students
// @Produce      json
// @Security     BearerAuth
// @Param        q  query  string  false  "Filter by name, ID or department"
// @Success      200  {object}  internal.ResponseBody{data=[]StudentData}
// @Router       /api/students [get]
func (sh StudentHandler) HandleGetAllStudents(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	studentData, err := GetAllStudents(sh.db)
	if err != nil {
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "could not get student data", Err: err}
	}

	internal.WriteData(w, "", FilterStudents(studentData, strings.TrimSpace(r.URL.Query().Get("q"))), nil)
	return nil
}

// FilterStudents returns the students matching every whitespace-separated
// term against username, student ID or department name.
func FilterStudents(students []StudentData, query string) []StudentData {
	if query == "" {
		return students
	}

	filtered := make([]StudentData, 0, len(students))
	for _, s := range students {
		if matches(s, query) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

func matches(s StudentData, query string) bool {
	for _, term := range strings.Fields(strings.ToLower(query)) {
		if !containsAny(term, s.Username, s.StudentID, s.DepartmentName) {
			return false
		}
	}
	return true
}

func containsAny(needle string, haystacks ...string) bool {
	for _, hay := range haystacks {
		if strings.Contains(strings.ToLower(hay), needle) {
			return true
		}
	}
	return false
}
