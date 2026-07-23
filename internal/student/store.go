package student

import "gorm.io/gorm"

func GetAllStudents(db *gorm.DB) ([]StudentData, error) {
	var result []StudentData
	err := db.
		Table("students as st").
		Select("st.username,st.student_id,d.name as department_name").
		Joins("JOIN departments d on d.id = st.department_id").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
