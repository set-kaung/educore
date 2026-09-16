package department

import (
	"educore/internal"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

type DepartmentData struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func ListDepartments(db *gorm.DB) ([]DepartmentData, error) {
	var departments []internal.Department
	err := db.Order("name").Find(&departments).Error
	if err != nil {
		return nil, err
	}

	result := make([]DepartmentData, 0, len(departments))
	for _, d := range departments {
		result = append(result, DepartmentData{ID: d.ID, Name: d.Name})
	}
	return result, nil
}
