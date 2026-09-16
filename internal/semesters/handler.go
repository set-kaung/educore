package semesters

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

// HandleList godoc
// @Summary      List semesters
// @Description  Get all semesters, with the current one flagged
// @Tags         semesters
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  internal.ResponseBody{data=[]SemesterData}
// @Router       /api/semesters [get]
func (h *Handler) HandleList(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	semesters, err := List(h.db)
	if err != nil {
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "could not list semesters", Err: err}
	}
	internal.WriteData(w, "", semesters, nil)
	return nil
}
