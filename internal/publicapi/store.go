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
	err := db.Table("users as u").
		Joins("JOIN departments d ON d.id = u.department_id").
		Where("u.student_id = ? AND u.role = ? AND d.name = ?", studentID, "student", departmentName).
		Count(&count).Error
	return count > 0, err
}
