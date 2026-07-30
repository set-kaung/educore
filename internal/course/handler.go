package course

import (
	"educore/internal"
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

type CourseHandler struct {
	db *gorm.DB
}

func NewCourseHandler(db *gorm.DB) *CourseHandler {
	return &CourseHandler{db: db}
}

// HandleGetBySemester godoc
// @Summary      List course offerings by semester
// @Description  Get all course offerings for a given semester
// @Tags         courses
// @Produce      json
// @Security     BearerAuth
// @Param        semester  query  string  true  "Semester (e.g. Fall2026)"
// @Success      200  {object}  internal.ResponseBody{data=[]CourseOffering}
// @Failure      400  {object}  internal.ResponseBody
// @Router       /courses [get]
func (h *CourseHandler) HandleGetBySemester(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	semester := r.URL.Query().Get("semester")
	if semester == "" {
		return &internal.HTTPError{StatusCode: http.StatusBadRequest, Message: "semester query parameter is required", Err: nil}
	}

	offerings, err := GetOfferingsBySemester(h.db, semester)
	if err != nil {
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "could not get course offerings", Err: err}
	}

	internal.WriteData(w, "", offerings, nil)
	return nil
}

// HandleCreateCourse godoc
// @Summary      Create a new course
// @Description  Add a new course listing. Professor only.
// @Tags         courses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  internal.Course  true  "Course details"
// @Success      200  {object}  internal.ResponseBody{data=internal.Course}
// @Failure      400  {object}  internal.ResponseBody
// @Router       /courses [post]
func (h *CourseHandler) HandleCreateCourse(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	var course internal.Course
	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		return &internal.HTTPError{StatusCode: http.StatusBadRequest, Message: "invalid request body", Err: err}
	}

	if err := CreateCourse(h.db, &course); err != nil {
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "could not create course", Err: err}
	}

	internal.WriteData(w, "", course, nil)
	return nil
}
