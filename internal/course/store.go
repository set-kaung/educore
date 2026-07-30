package course

import (
	"educore/internal"

	"gorm.io/gorm"
)

func CreateCourse(db *gorm.DB, c *internal.Course) error {
	return db.Create(c).Error
}
