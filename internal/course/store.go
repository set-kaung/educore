package course

import (
	"educore/internal"

	"gorm.io/gorm"
)

func CreateCourse(db *gorm.DB, c *internal.Course) error {
	return db.Create(c).Error
}

func ListCourses(db *gorm.DB) ([]internal.Course, error) {
	var courses []internal.Course
	err := db.Order("name asc").Find(&courses).Error
	return courses, err
}
