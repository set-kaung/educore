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

func DeleteCourse(db *gorm.DB, id uint) error {
	res := db.Delete(&internal.Course{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
