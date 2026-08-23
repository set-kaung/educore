package semestercourse

import (
	"educore/internal"
	"time"

	"gorm.io/gorm"
)

func createSemesterCourse(db *gorm.DB, sc internal.SemesterCourse) (internal.SemesterCourse, error) {
	err := db.Create(&sc).Error
	return sc, err
}

func createSchedule(db *gorm.DB, schedule internal.SemesterCourseSchedule) error {
	return db.Create(&schedule).Error
}

func HasConflict(db *gorm.DB, professorID uint, semester string, from, to time.Time) (bool, error) {
	var count int64
	err := db.Table("semester_course_schedules as scs").
		Joins("JOIN semester_courses sc ON sc.id = scs.semester_course_id").
		Where("sc.semester = ? AND sc.taught_by = ? AND scs.from < ? AND scs.to > ?",
			semester, professorID, to, from).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

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

func GetEnrolledCourses(db *gorm.DB, studentID uint) ([]CourseOffering, error) {
	var result []CourseOffering
	err := db.
		Table("enrollments as e").
		Select("c.name, c.course_code, sc.section, sc.semester, p.name as professor_name").
		Joins("JOIN semester_courses sc ON sc.id = e.semester_course_id").
		Joins("JOIN courses c ON c.id = sc.course_id").
		Joins("JOIN professors p ON p.id = sc.taught_by").
		Where("e.student_id = ? AND e.deleted_at IS NULL", studentID).
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
