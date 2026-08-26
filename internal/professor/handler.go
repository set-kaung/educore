package professor

import (
	"educore/internal"
	"net/http"

	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// HandleListProfessors godoc
// @Summary      List professors
// @Description  Get all professors by name
// @Tags         professors
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  internal.ResponseBody{data=[]ProfessorData}
// @Router       /api/professors [get]
func (h Handler) HandleListProfessors(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	professors, err := ListProfessors(h.db)
	if err != nil {
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "could not list professors", Err: err}
	}
	internal.WriteData(w, "", professors, nil)
	return nil
}
