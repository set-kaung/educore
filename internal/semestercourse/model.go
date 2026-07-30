package semestercourse

import (
	"errors"
	"time"
)

var ErrInvalidScheduleTime = errors.New("invalid schedule time")
var ErrScheduleConflict = errors.New("professor has a conflicting schedule")

type Schedule struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type CourseOffering struct {
	Name          string `json:"name"`
	CourseCode    string `json:"course_code"`
	Section       string `json:"section"`
	Semester      string `json:"semester"`
	ProfessorName string `json:"professor_name"`
}

type CreateRequest struct {
	Semester  string     `json:"semester"`
	CourseID  uint       `json:"course_id"`
	Section   string     `json:"section"`
	TeachBy   uint       `json:"taught_by"`
	Schedules []Schedule `json:"schedules"`
}

type CreateResponse struct {
	SemesterCourseID uint       `json:"semester_course_id"`
	Semester         string     `json:"semester"`
	CourseID         uint       `json:"course_id"`
	Section          string     `json:"section"`
	TeachBy          uint       `json:"taught_by"`
	Schedules        []Schedule `json:"schedules"`
}
