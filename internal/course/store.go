package course

import (
	"educore/internal"

	"gorm.io/gorm"
)

func GetOfferingsBySemester(db *gorm.DB, semester string) ([]CourseOffering, error) {
	var result []CourseOffering
	err := db.
		Table("semester_courses as sc").
		Select("c.name, c.course_code, sc.section, sc.semester, p.name as professor_name").
		Joins("JOIN courses c ON c.id = sc.course_id").
		Joins("JOIN professors p ON p.id = sc.taught_by").
		Where("sc.semester = ?", semester).
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func CreateCourse(db *gorm.DB, c *internal.Course) error {
	return db.Create(c).Error
}
