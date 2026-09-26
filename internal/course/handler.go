package course

import (
	"educore/internal"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type CourseHandler struct {
	db *gorm.DB
}

func NewCourseHandler(db *gorm.DB) *CourseHandler {
	return &CourseHandler{db: db}
}

// HandleListCourses godoc
// @Summary      List courses
// @Description  Get all course listings
// @Tags         courses
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  internal.ResponseBody{data=[]internal.Course}
// @Router       /api/courses [get]
func (h *CourseHandler) HandleListCourses(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	courses, err := ListCourses(h.db)
	if err != nil {
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "could not list courses", Err: err}
	}

	res := make([]CourseResponse, 0, len(courses))
	for _, c := range courses {
		res = append(res, ToCourseResponse(c))
	}
	internal.WriteData(w, "", res, nil)
	return nil
}

// HandleCreateCourse godoc
// @Summary      Create a new course
// @Description  Add a new course listing. Professor only.
// @Tags         courses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  CourseRequest  true  "Course details"
// @Success      200  {object}  internal.ResponseBody{data=CourseResponse}
// @Failure      400  {object}  internal.ResponseBody
// @Router       /api/courses [post]
func (h *CourseHandler) HandleCreateCourse(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	var req CourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &internal.HTTPError{StatusCode: http.StatusBadRequest, Message: "invalid request body", Err: err}
	}

	c := internal.Course{
		Name:       req.Name,
		CourseCode: req.CourseCode,
	}

	if err := CreateCourse(h.db, &c); err != nil {
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "could not create course", Err: err}
	}

	internal.WriteData(w, "", ToCourseResponse(c), nil)
	return nil
}

// HandleDeleteCourse godoc
// @Summary      Delete a course listing
// @Description  Soft-delete a course listing. Professor only.
// @Tags         courses
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "Course ID"
// @Success      200  {object}  internal.ResponseBody
// @Failure      400  {object}  internal.ResponseBody
// @Failure      404  {object}  internal.ResponseBody
// @Router       /api/courses/{id} [delete]
func (h *CourseHandler) HandleDeleteCourse(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		return &internal.HTTPError{StatusCode: http.StatusBadRequest, Message: "invalid course id", Err: nil}
	}

	if err := DeleteCourse(h.db, uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &internal.HTTPError{StatusCode: http.StatusNotFound, Message: "course not found", Err: err}
		}
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "could not delete course", Err: err}
	}

	internal.WriteData(w, "course deleted", nil, nil)
	return nil
}
