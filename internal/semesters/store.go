package semesters

import (
	"educore/internal"

	"gorm.io/gorm"
)

func List(db *gorm.DB) ([]SemesterData, error) {
	var semesters []internal.Semester
	if err := db.Order("is_current desc, start_date asc, id desc").Find(&semesters).Error; err != nil {
		return nil, err
	}
	result := make([]SemesterData, 0, len(semesters))
	for _, s := range semesters {
		result = append(result, toData(s))
	}
	return result, nil
}

func GetCurrent(db *gorm.DB) (internal.Semester, error) {
	var semester internal.Semester
	err := db.Where("is_current = ?", true).First(&semester).Error
	return semester, err
}

func GetByID(db *gorm.DB, id uint) (internal.Semester, error) {
	var semester internal.Semester
	err := db.First(&semester, id).Error
	return semester, err
}
