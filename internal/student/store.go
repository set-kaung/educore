package student

import "gorm.io/gorm"

func GetAllStudents(db *gorm.DB) ([]StudentData, error) {
	result := make([]StudentData, 0)
	err := db.
		Table("users as u").
		Select("u.id,u.name,u.student_id,d.name as department_name").
		Joins("JOIN departments d on d.id = u.department_id").
		Where("u.role = ?", "student").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
