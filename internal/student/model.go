package student

type StudentData struct {
	Username       string `gorm:"column:username"`
	StudentID      string `gorm:"column:student_id"`
	DepartmentName string `gorm:"column:department_name"`
}
