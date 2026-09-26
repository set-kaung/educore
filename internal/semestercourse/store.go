package semestercourse

import (
	"educore/internal"
	"time"

	"gorm.io/gorm"
)

const scheduleAgg = "COALESCE(GROUP_CONCAT(CONCAT(scs.weekday, ' ', TIME_FORMAT(scs.start_time, '%H:%i'), '-', TIME_FORMAT(scs.end_time, '%H:%i')) ORDER BY scs.start_time, scs.end_time SEPARATOR ', '), '')"

const offeringSelect = "sc.id as semester_course_id, c.name, c.course_code, sc.section, sem.name as semester, p.name as professor_name, " +
	scheduleAgg + " as schedule"

const offeringGroup = "sc.id, c.name, c.course_code, sc.section, sem.name, p.name"

const scheduleJoin = "JOIN semester_course_schedules scs ON scs.semester_course_id = sc.id AND scs.deleted_at IS NULL"

func createSemesterCourse(db *gorm.DB, sc internal.SemesterCourse) (internal.SemesterCourse, error) {
	err := db.Create(&sc).Error
	return sc, err
}

func createSchedule(db *gorm.DB, schedule internal.SemesterCourseSchedule) error {
	return db.Create(&schedule).Error
}

func createEnrollment(db *gorm.DB, enrollment internal.Enrollment) (internal.Enrollment, error) {
	err := db.Create(&enrollment).Error
	return enrollment, err
}

func getSemesterCourse(db *gorm.DB, id uint) (internal.SemesterCourse, error) {
	var sc internal.SemesterCourse
	err := db.First(&sc, id).Error
	return sc, err
}

func DeleteSemesterCourse(db *gorm.DB, id uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var sc internal.SemesterCourse
		if err := tx.First(&sc, id).Error; err != nil {
			return err
		}
		if err := tx.Delete(&internal.SemesterCourse{}, id).Error; err != nil {
			return err
		}
		return tx.Delete(&internal.SemesterCourseSchedule{}, "semester_course_id = ?", id).Error
	})
}

func IsEnrolled(db *gorm.DB, studentID, semesterCourseID uint) (bool, error) {
	var count int64
	err := db.Model(&internal.Enrollment{}).
		Joins("JOIN semester_courses sc ON sc.id = enrollments.semester_course_id AND sc.deleted_at IS NULL").
		Where("enrollments.user_id = ? AND enrollments.semester_course_id = ? AND enrollments.deleted_at IS NULL", studentID, semesterCourseID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func HasStudentScheduleConflict(db *gorm.DB, studentID, semesterCourseID uint) (bool, error) {
	var count int64
	err := db.Table("semester_course_schedules as ts").
		Joins("JOIN semester_courses tsc ON tsc.id = ts.semester_course_id AND tsc.deleted_at IS NULL").
		Joins("JOIN enrollments e ON e.user_id = ? AND e.deleted_at IS NULL", studentID).
		Joins("JOIN semester_courses esc ON esc.id = e.semester_course_id AND esc.semester_id = tsc.semester_id AND esc.id != tsc.id AND esc.deleted_at IS NULL").
		Joins("JOIN semester_course_schedules es ON es.semester_course_id = esc.id AND es.weekday = ts.weekday AND es.start_time < ts.end_time AND es.end_time > ts.start_time AND es.deleted_at IS NULL").
		Where("tsc.id = ? AND ts.deleted_at IS NULL", semesterCourseID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func HasConflict(db *gorm.DB, professorID, semesterID uint, weekday string, from, to time.Time) (bool, error) {
	var count int64
	err := db.Table("semester_course_schedules as scs").
		Joins("JOIN semester_courses sc ON sc.id = scs.semester_course_id AND sc.deleted_at IS NULL").
		Where("scs.deleted_at IS NULL AND sc.semester_id = ? AND sc.taught_by = ? AND scs.weekday = ? AND scs.start_time < ? AND scs.end_time > ?",
			semesterID, professorID, weekday, to, from).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetOfferingsBySemester(db *gorm.DB, semesterID uint) ([]CourseOffering, error) {
	var result []CourseOffering
	err := db.
		Table("semester_courses as sc").
		Select(offeringSelect).
		Joins("JOIN courses c ON c.id = sc.course_id").
		Joins("JOIN semesters sem ON sem.id = sc.semester_id").
		Joins("JOIN users p ON p.id = sc.taught_by").
		Joins(scheduleJoin).
		Where("sc.semester_id = ? AND sc.deleted_at IS NULL AND c.deleted_at IS NULL AND sem.deleted_at IS NULL AND p.deleted_at IS NULL", semesterID).
		Group(offeringGroup).
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
		Select(offeringSelect).
		Joins("JOIN semester_courses sc ON sc.id = e.semester_course_id").
		Joins("JOIN courses c ON c.id = sc.course_id").
		Joins("JOIN semesters sem ON sem.id = sc.semester_id").
		Joins("JOIN users p ON p.id = sc.taught_by").
		Joins(scheduleJoin).
		Where("e.user_id = ? AND e.deleted_at IS NULL AND sc.deleted_at IS NULL AND c.deleted_at IS NULL AND sem.deleted_at IS NULL AND p.deleted_at IS NULL", studentID).
		Group(offeringGroup).
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetOfferingDetail(db *gorm.DB, id uint) (*CourseOffering, error) {
	var result CourseOffering
	err := db.
		Table("semester_courses as sc").
		Select(offeringSelect).
		Joins("JOIN courses c ON c.id = sc.course_id").
		Joins("JOIN semesters sem ON sem.id = sc.semester_id").
		Joins("JOIN users p ON p.id = sc.taught_by").
		Joins(scheduleJoin).
		Where("sc.id = ? AND sc.deleted_at IS NULL AND c.deleted_at IS NULL AND sem.deleted_at IS NULL AND p.deleted_at IS NULL", id).
		Group(offeringGroup).
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func GetTaughtCourses(db *gorm.DB, professorID uint) ([]CourseOffering, error) {
	var result []CourseOffering
	err := db.
		Table("semester_courses as sc").
		Select(offeringSelect).
		Joins("JOIN courses c ON c.id = sc.course_id").
		Joins("JOIN semesters sem ON sem.id = sc.semester_id").
		Joins("JOIN users p ON p.id = sc.taught_by").
		Joins(scheduleJoin).
		Where("sc.taught_by = ? AND sc.deleted_at IS NULL AND c.deleted_at IS NULL AND sem.deleted_at IS NULL AND p.deleted_at IS NULL", professorID).
		Group(offeringGroup).
		Order("sem.name DESC, c.name").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
