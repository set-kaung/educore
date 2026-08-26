package semestercourse

import (
	"educore/internal"
	"strings"
	"time"

	"gorm.io/gorm"
)

const timeLayout = "15:04"

func parseWeekday(w string) (string, bool) {
	w = strings.TrimSpace(w)
	if w == "" {
		return "", false
	}
	w = strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
	if !validWeekdays[w] {
		return "", false
	}
	return w, true
}

func parseTimeOfDay(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	// Parse against a fixed valid date so MySQL TIME comparisons don't see year 0.
	t, err := time.Parse("2006-01-02 "+timeLayout, "2000-01-01 "+s)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

func CreateSemesterCourse(db *gorm.DB, req CreateRequest) (CreateResponse, error) {
	parsedSchedules := make([]internal.SemesterCourseSchedule, 0, len(req.Schedules))

	for _, sched := range req.Schedules {
		weekday, ok := parseWeekday(sched.Weekday)
		if !ok {
			return CreateResponse{}, ErrInvalidWeekday
		}

		start, ok := parseTimeOfDay(sched.StartTime)
		if !ok {
			return CreateResponse{}, ErrInvalidScheduleTime
		}

		end, ok := parseTimeOfDay(sched.EndTime)
		if !ok {
			return CreateResponse{}, ErrInvalidScheduleTime
		}

		if !end.After(start) {
			return CreateResponse{}, ErrInvalidScheduleTime
		}

		conflict, err := HasConflict(db, req.TeachBy, req.Semester, weekday, start, end)
		if err != nil {
			return CreateResponse{}, err
		}
		if conflict {
			return CreateResponse{}, ErrScheduleConflict
		}

		parsedSchedules = append(parsedSchedules, internal.SemesterCourseSchedule{
			SemesterCourseID: 0,
			Weekday:          weekday,
			StartTime:        start,
			EndTime:          end,
		})
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

	for _, sched := range parsedSchedules {
		sched.SemesterCourseID = sc.ID
		if err := createSchedule(db, sched); err != nil {
			return CreateResponse{}, err
		}
	}

	responseSchedules := make([]Schedule, 0, len(parsedSchedules))
	for _, sched := range parsedSchedules {
		responseSchedules = append(responseSchedules, Schedule{
			Weekday:   sched.Weekday,
			StartTime: sched.StartTime.Format(timeLayout),
			EndTime:   sched.EndTime.Format(timeLayout),
		})
	}

	return CreateResponse{
		SemesterCourseID: sc.ID,
		Semester:         sc.Semester,
		CourseID:         sc.CourseID,
		Section:          sc.Section,
		TeachBy:          sc.TeachBy,
		Schedules:        responseSchedules,
	}, nil
}
