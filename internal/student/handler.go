package student

import (
	"educore/internal"

	"net/http"

	"gorm.io/gorm"
)

type StudentHandler struct {
	db *gorm.DB
}

func NewStudentHandler(db *gorm.DB) *StudentHandler {
	return &StudentHandler{db: db}
}

func (sh StudentHandler) HandleGetAllStudents(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	studentData, err := GetAllStudents(sh.db)
	if err != nil {
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "could not get student data", Err: err}
	}
	internal.WriteData(w, "", studentData, nil)
	return nil
}
