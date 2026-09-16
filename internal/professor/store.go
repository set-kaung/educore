package professor

import (
	"educore/internal"

	"gorm.io/gorm"
)

func ListProfessors(db *gorm.DB) ([]ProfessorData, error) {
	var result []ProfessorData
	err := db.Model(&internal.User{}).
		Select("id, name").
		Where("role = ?", "professor").
		Order("name asc").
		Scan(&result).Error
	return result, err
}
