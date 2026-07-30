package semestercourse

import "time"

type Schedule struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
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
