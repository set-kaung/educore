package main

import (
	"educore/internal"
	"educore/internal/student"
	"net/http"

	"gorm.io/gorm"
)

type StudentHandler struct {
	db *gorm.DB
}

func (sh StudentHandler) HandleGetAllStudents(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	studentData, err := student.GetAllStudents(sh.db)
	if err != nil {
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "could not get student data", Err: err}
	}
	internal.WriteData(w, "", studentData, nil)
	return nil
}
