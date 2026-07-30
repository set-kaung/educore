package semestercourse

import (
	"errors"
	"time"

	"educore/internal"
	"gorm.io/gorm"
)

var ErrScheduleConflict = errors.New("professor has a conflicting schedule")

func Create(db *gorm.DB, req CreateRequest) (CreateResponse, error) {
	var response CreateResponse

	err := db.Transaction(func(tx *gorm.DB) error {
		// Check schedule conflicts for the professor in this semester
		for _, s := range req.Schedules {
			if s.From.IsZero() || s.To.IsZero() || !s.To.After(s.From) {
				return errors.New("invalid schedule time")
			}

			var count int64
			err := tx.Table("semester_course_schedules as scs").
				Joins("JOIN semester_courses sc ON sc.id = scs.semester_course_id").
				Where("sc.semester = ? AND sc.taught_by = ? AND scs.from < ? AND scs.to > ?",
					req.Semester, req.TeachBy, s.To, s.From).
				Count(&count).Error
			if err != nil {
				return err
			}
			if count > 0 {
				return ErrScheduleConflict
			}
		}

		// Create semester course
		sc := internal.SemesterCourse{
			Semester: req.Semester,
			CourseID: req.CourseID,
			Section:  req.Section,
			TeachBy:  req.TeachBy,
		}
		if err := tx.Create(&sc).Error; err != nil {
			return err
		}

		response = CreateResponse{
			SemesterCourseID: sc.ID,
			Semester:         sc.Semester,
			CourseID:         sc.CourseID,
			Section:          sc.Section,
			TeachBy:          sc.TeachBy,
			Schedules:        req.Schedules,
		}

		// Create schedules
		for _, s := range req.Schedules {
			schedule := internal.SemesterCourseSchedule{
				SemesterCourseID: sc.ID,
				From:             s.From,
				To:               &s.To,
			}
			if err := tx.Create(&schedule).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return CreateResponse{}, err
	}
	return response, nil
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
