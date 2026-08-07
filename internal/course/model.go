package course

import "educore/internal"

type CourseRequest struct {
	Name       string `json:"name"`
	CourseCode string `json:"course_code"`
}

type CourseResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	CourseCode string `json:"course_code"`
}

func ToCourseResponse(c internal.Course) CourseResponse {
	return CourseResponse{
		ID:         c.ID,
		Name:       c.Name,
		CourseCode: c.CourseCode,
	}
}
