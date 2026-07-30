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

// HandleGetAllStudents godoc
// @Summary      List students
// @Description  Get all students with department names
// @Tags         students
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  internal.ResponseBody{data=[]StudentData}
// @Router       /student [get]
func (sh StudentHandler) HandleGetAllStudents(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	studentData, err := GetAllStudents(sh.db)
	if err != nil {
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "could not get student data", Err: err}
	}
	internal.WriteData(w, "", studentData, nil)
	return nil
}
