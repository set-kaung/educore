package semestercourse

import (
	"errors"
)

var ErrInvalidScheduleTime = errors.New("invalid schedule time")
var ErrInvalidWeekday = errors.New("invalid weekday")
var ErrInvalidSemester = errors.New("invalid semester")
var ErrScheduleConflict = errors.New("professor has a conflicting schedule")
var ErrSemesterCourseNotFound = errors.New("semester course not found")
var ErrAlreadyEnrolled = errors.New("already enrolled in this course")
var ErrEnrollmentConflict = errors.New("schedule conflicts with an enrolled course")
var ErrSemesterNotCurrent = errors.New("enrollment is only allowed for the current semester")

var validWeekdays = map[string]bool{
	"Monday":    true,
	"Tuesday":   true,
	"Wednesday": true,
	"Thursday":  true,
	"Friday":    true,
	"Saturday":  true,
	"Sunday":    true,
}

type Schedule struct {
	Weekday   string `json:"weekday"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type CourseOffering struct {
	SemesterCourseID uint   `json:"semester_course_id"`
	Name             string `json:"name"`
	CourseCode       string `json:"course_code"`
	Section          string `json:"section"`
	Semester         string `json:"semester"`
	ProfessorName    string `json:"professor_name"`
	Schedule         string `json:"schedule"`
}

type CreateRequest struct {
	SemesterID uint       `json:"semester_id"`
	CourseID   uint       `json:"course_id"`
	Section    string     `json:"section"`
	TeachBy    uint       `json:"taught_by"`
	Schedules  []Schedule `json:"schedules"`
}

type CreateResponse struct {
	SemesterCourseID uint       `json:"semester_course_id"`
	Semester         string     `json:"semester"`
	CourseID         uint       `json:"course_id"`
	Section          string     `json:"section"`
	TeachBy          uint       `json:"taught_by"`
	Schedules        []Schedule `json:"schedules"`
}

type EnrollResponse struct {
	EnrollmentID     uint   `json:"enrollment_id"`
	SemesterCourseID uint   `json:"semester_course_id"`
	Semester         string `json:"semester"`
	CourseID         uint   `json:"course_id"`
	Section          string `json:"section"`
}
