package semesters

import (
	"educore/internal"
	"time"
)

type SemesterData struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	IsCurrent bool       `json:"is_current"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}

func toData(s internal.Semester) SemesterData {
	return SemesterData{
		ID:        s.ID,
		Name:      s.Name,
		IsCurrent: s.IsCurrent,
		StartDate: s.StartDate,
		EndDate:   s.EndDate,
	}
}
