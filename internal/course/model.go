package course

type CourseOffering struct {
	Name          string `gorm:"column:name"`
	CourseCode    string `gorm:"column:course_code"`
	Section       string `gorm:"column:section"`
	Semester      string `gorm:"column:semester"`
	ProfessorName string `gorm:"column:professor_name"`
}
