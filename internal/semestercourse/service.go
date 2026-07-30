package semestercourse

import (
	"educore/internal"

	"gorm.io/gorm"
)

func CreateSemesterCourse(db *gorm.DB, req CreateRequest) (CreateResponse, error) {
	for _, sched := range req.Schedules {
		if sched.From.IsZero() || sched.To.IsZero() || !sched.To.After(sched.From) {
			return CreateResponse{}, ErrInvalidScheduleTime
		}

		conflict, err := HasConflict(db, req.TeachBy, req.Semester, sched.From, sched.To)
		if err != nil {
			return CreateResponse{}, err
		}
		if conflict {
			return CreateResponse{}, ErrScheduleConflict
		}
	}

	sc, err := createSemesterCourse(db, internal.SemesterCourse{
		Semester: req.Semester,
		CourseID: req.CourseID,
		Section:  req.Section,
		TeachBy:  req.TeachBy,
	})
	if err != nil {
		return CreateResponse{}, err
	}

	for _, sched := range req.Schedules {
		to := sched.To
		err := createSchedule(db, internal.SemesterCourseSchedule{
			SemesterCourseID: sc.ID,
			From:             sched.From,
			To:               &to,
		})
		if err != nil {
			return CreateResponse{}, err
		}
	}

	return CreateResponse{
		SemesterCourseID: sc.ID,
		Semester:         sc.Semester,
		CourseID:         sc.CourseID,
		Section:          sc.Section,
		TeachBy:          sc.TeachBy,
		Schedules:        req.Schedules,
	}, nil
}
