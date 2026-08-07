package publicapi

import (
	"educore/internal"

	"gorm.io/gorm"
)

func ValidateAPIKey(db *gorm.DB, key string) (bool, error) {
	var count int64
	err := db.Model(&internal.ApiKey{}).
		Where("key = ? AND valid = ?", key, true).
		Count(&count).Error
	return count > 0, err
}

func IsStudentInDepartment(db *gorm.DB, studentID string, departmentName string) (bool, error) {
	var count int64
	err := db.Table("students as s").
		Joins("JOIN departments d ON d.id = s.department_id").
		Where("s.student_id = ? AND d.name = ?", studentID, departmentName).
		Count(&count).Error
	return count > 0, err
}
