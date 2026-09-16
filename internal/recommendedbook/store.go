package recommendedbook

import (
	"educore/internal"

	"gorm.io/gorm"
)

func listBooks(db *gorm.DB, semesterCourseID uint) ([]Book, error) {
	var models []internal.RecommendedBook
	err := db.Where("semester_course_id = ?", semesterCourseID).Find(&models).Error
	if err != nil {
		return nil, err
	}
	books := make([]Book, 0, len(models))
	for _, m := range models {
		books = append(books, toBook(m))
	}
	return books, nil
}

func createBook(db *gorm.DB, semesterCourseID uint, book internal.RecommendedBook) (Book, error) {
	if err := db.Create(&book).Error; err != nil {
		return Book{}, err
	}
	return toBook(book), nil
}

func deleteBook(db *gorm.DB, id uint) error {
	return db.Delete(&internal.RecommendedBook{}, id).Error
}

func getSemesterCourse(db *gorm.DB, id uint) (internal.SemesterCourse, error) {
	var sc internal.SemesterCourse
	err := db.First(&sc, id).Error
	return sc, err
}

func exists(db *gorm.DB, semesterCourseID uint, key string) (bool, error) {
	var count int64
	err := db.Model(&internal.RecommendedBook{}).
		Where("semester_course_id = ? AND open_library_key = ?", semesterCourseID, key).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func toBook(m internal.RecommendedBook) Book {
	return Book{
		ID:               m.ID,
		SemesterCourseID: m.SemesterCourseID,
		Title:            m.Title,
		Author:           m.Author,
		ISBN:             m.Isbn,
		CoverI:           m.CoverI,
		OpenLibraryKey:   m.OpenLibraryKey,
	}
}
