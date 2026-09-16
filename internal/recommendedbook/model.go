package recommendedbook

import "errors"

var ErrSemesterCourseNotFound = errors.New("semester course not found")
var ErrBookNotFound = errors.New("recommended book not found")
var ErrBookExists = errors.New("book already recommended")
var ErrNotTeacher = errors.New("only the professor teaching this course may manage its books")

type Book struct {
	ID               uint   `json:"id"`
	SemesterCourseID uint   `json:"semester_course_id"`
	Title            string `json:"title"`
	Author           string `json:"author"`
	ISBN             string `json:"isbn"`
	CoverI           int    `json:"cover_i"`
	OpenLibraryKey   string `json:"open_library_key"`
}

type AddRequest struct {
	Title          string `json:"title"`
	Author         string `json:"author"`
	ISBN           string `json:"isbn"`
	CoverI         int    `json:"cover_i"`
	OpenLibraryKey string `json:"open_library_key"`
}
