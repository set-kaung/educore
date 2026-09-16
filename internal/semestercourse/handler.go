package semestercourse

import (
	"educore/internal"
	"educore/internal/auth"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// HandleGetEnrolled godoc
// @Summary      List the authenticated student's enrolled courses
// @Description  Get all course offerings the current student is enrolled in
// @Tags         semester-courses
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  internal.ResponseBody{data=[]CourseOffering}
// @Router       /api/my/courses [get]
func (h *Handler) HandleGetEnrolled(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	claims := auth.GetClaims(r)

	courses, err := GetEnrolledCourses(h.db, claims.UserID)
	if err != nil {
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "could not get enrolled courses", Err: err}
	}

	internal.WriteData(w, "", courses, nil)
	return nil
}

// HandleGetTaught godoc
// @Summary      List the authenticated professor's taught courses
// @Description  Get all course offerings the current professor teaches
// @Tags         semester-courses
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  internal.ResponseBody{data=[]CourseOffering}
// @Router       /api/my/taught-courses [get]
func (h *Handler) HandleGetTaught(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	claims := auth.GetClaims(r)

	courses, err := GetTaughtCourses(h.db, claims.UserID)
	if err != nil {
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "could not get taught courses", Err: err}
	}

	internal.WriteData(w, "", courses, nil)
	return nil
}

// HandleGetBySemester godoc
// @Summary      List semester course offerings
// @Description  Get all course offerings for a given semester
// @Tags         semester-courses
// @Produce      json
// @Security     BearerAuth
// @Param        semester_id  query  int  true  "Semester ID"
// @Success      200  {object}  internal.ResponseBody{data=[]CourseOffering}
// @Failure      400  {object}  internal.ResponseBody
// @Router       /api/semester-courses [get]
func (h *Handler) HandleGetBySemester(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	semesterID, err := strconv.ParseUint(r.URL.Query().Get("semester_id"), 10, 64)
	if err != nil || semesterID == 0 {
		return &internal.HTTPError{StatusCode: http.StatusBadRequest, Message: "semester_id query parameter is required", Err: nil}
	}

	offerings, err := GetOfferingsBySemester(h.db, uint(semesterID))
	if err != nil {
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "could not get course offerings", Err: err}
	}

	internal.WriteData(w, "", offerings, nil)
	return nil
}

// HandleGetDetail godoc
// @Summary      Get a single course offering
// @Description  Get detailed information about a specific semester course offering
// @Tags         semester-courses
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "Semester course ID"
// @Success      200  {object}  internal.ResponseBody{data=CourseOffering}
// @Failure      400  {object}  internal.ResponseBody
// @Failure      404  {object}  internal.ResponseBody
// @Router       /api/semester-courses/{id} [get]
func (h *Handler) HandleGetDetail(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		return &internal.HTTPError{StatusCode: http.StatusBadRequest, Message: "invalid semester course id", Err: nil}
	}

	offering, err := GetOfferingDetail(h.db, uint(id))
	if err != nil {
		return &internal.HTTPError{StatusCode: http.StatusNotFound, Message: "course offering not found", Err: err}
	}

	internal.WriteData(w, "", offering, nil)
	return nil
}

// HandleCreate godoc
// @Summary      Create semester course schedule
// @Description  Schedule a course for a semester with sections and timeslots. Checks professor conflicts.
// @Tags         semester-courses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  CreateRequest  true  "Semester course details"
// @Success      200  {object}  internal.ResponseBody{data=CreateResponse}
// @Failure      400  {object}  internal.ResponseBody
// @Failure      409  {object}  internal.ResponseBody
// @Router       /api/semester-courses [post]
func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &internal.HTTPError{StatusCode: http.StatusBadRequest, Message: "invalid request body", Err: err}
	}

	res, err := CreateSemesterCourse(h.db, req)
	if err != nil {
		if errors.Is(err, ErrScheduleConflict) {
			return &internal.HTTPError{StatusCode: http.StatusConflict, Message: "professor has a conflicting schedule", Err: err}
		}
		if errors.Is(err, ErrInvalidScheduleTime) {
			return &internal.HTTPError{StatusCode: http.StatusBadRequest, Message: "invalid schedule time", Err: err}
		}
		if errors.Is(err, ErrInvalidWeekday) {
			return &internal.HTTPError{StatusCode: http.StatusBadRequest, Message: "invalid weekday", Err: err}
		}
		if errors.Is(err, ErrInvalidSemester) {
			return &internal.HTTPError{StatusCode: http.StatusBadRequest, Message: "invalid semester", Err: err}
		}
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "could not create semester course", Err: err}
	}

	internal.WriteData(w, "", res, nil)
	return nil
}

// HandleEnroll godoc
// @Summary      Enroll in a course offering
// @Description  Enroll the authenticated student in a semester course offering. Rejects duplicate enrollment and schedule conflicts with courses the student is already enrolled in.
// @Tags         semester-courses
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "Semester course ID"
// @Success      200  {object}  internal.ResponseBody{data=EnrollResponse}
// @Failure      400  {object}  internal.ResponseBody
// @Failure      404  {object}  internal.ResponseBody
// @Failure      409  {object}  internal.ResponseBody
// @Router       /api/semester-courses/{id}/enroll [post]
func (h *Handler) HandleEnroll(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	claims := auth.GetClaims(r)

	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		return &internal.HTTPError{StatusCode: http.StatusBadRequest, Message: "invalid semester course id", Err: err}
	}

	res, err := EnrollInSemesterCourse(h.db, claims.UserID, uint(id))
	if err != nil {
		if errors.Is(err, ErrSemesterCourseNotFound) {
			return &internal.HTTPError{StatusCode: http.StatusNotFound, Message: "course offering not found", Err: err}
		}
		if errors.Is(err, ErrAlreadyEnrolled) {
			return &internal.HTTPError{StatusCode: http.StatusConflict, Message: "already enrolled in this course", Err: err}
		}
		if errors.Is(err, ErrEnrollmentConflict) {
			return &internal.HTTPError{StatusCode: http.StatusConflict, Message: "schedule conflicts with an enrolled course", Err: err}
		}
		if errors.Is(err, ErrSemesterNotCurrent) {
			return &internal.HTTPError{StatusCode: http.StatusConflict, Message: "enrollment is only allowed for the current semester", Err: err}
		}
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "could not enroll in course", Err: err}
	}

	internal.WriteData(w, "enrolled", res, nil)
	return nil
}
