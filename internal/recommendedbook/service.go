package recommendedbook

import (
	"educore/internal"
	"errors"
	"strings"

	"gorm.io/gorm"
)

func ListBooks(db *gorm.DB, semesterCourseID uint) ([]Book, error) {
	return listBooks(db, semesterCourseID)
}

func AddBook(db *gorm.DB, teacherID, semesterCourseID uint, req AddRequest) (Book, error) {
	sc, err := getSemesterCourse(db, semesterCourseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Book{}, ErrSemesterCourseNotFound
		}
		return Book{}, err
	}
	if sc.TeachBy != teacherID {
		return Book{}, ErrNotTeacher
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return Book{}, errors.New("book title is required")
	}

	key := strings.TrimSpace(req.OpenLibraryKey)
	exists, err := exists(db, semesterCourseID, key)
	if err != nil {
		return Book{}, err
	}
	if exists {
		return Book{}, ErrBookExists
	}

	return createBook(db, semesterCourseID, internal.RecommendedBook{
		SemesterCourseID: semesterCourseID,
		Title:            strings.TrimSpace(req.Title),
		Author:           strings.TrimSpace(req.Author),
		Isbn:             strings.TrimSpace(req.ISBN),
		CoverI:           req.CoverI,
		OpenLibraryKey:   key,
	})
}

func RemoveBook(db *gorm.DB, teacherID, bookID uint) error {
	var book internal.RecommendedBook
	err := db.First(&book, bookID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrBookNotFound
		}
		return err
	}

	sc, err := getSemesterCourse(db, book.SemesterCourseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSemesterCourseNotFound
		}
		return err
	}
	if sc.TeachBy != teacherID {
		return ErrNotTeacher
	}

	return deleteBook(db, book.ID)
}
