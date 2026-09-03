package auth

import (
	"educore/internal"
	"errors"

	"gorm.io/gorm"
)

func FindUserByADObjectID(db *gorm.DB, oid string) (uint, string, bool, error) {
	var user internal.User
	err := db.Where("ad_object_id = ?", oid).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, "", false, nil
		}
		return 0, "", false, err
	}
	return user.ID, user.Role, true, nil
}

func CreateUserForOIDC(db *gorm.DB, name, username, studentID string, departmentID uint, oid string) (internal.User, error) {
	user := internal.User{
		Name:         name,
		Username:     username,
		StudentID:    studentID,
		DepartmentID: departmentID,
		ADObjectID:   oid,
		Role:         "student",
	}
	err := db.Create(&user).Error
	return user, err
}
