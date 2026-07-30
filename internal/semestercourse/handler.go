package semestercourse

import (
	"educore/internal"
	"encoding/json"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// HandleCreate godoc
// @Summary      Create semester course schedule
// @Description  Schedule a course for a semester with sections and timeslots. Checks professor conflicts.
// @Tags         semester-courses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  CreateRequest  true  "Semester course details"
// @Success      200  {object}  internal.ResponseBody{data=CreateResponse}
// @Failure      400  {object}  internal.ResponseBody
// @Failure      409  {object}  internal.ResponseBody
// @Router       /semester-courses [post]
func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &internal.HTTPError{StatusCode: http.StatusBadRequest, Message: "invalid request body", Err: err}
	}

	res, err := Create(h.db, req)
	if err != nil {
		if errors.Is(err, ErrScheduleConflict) {
			return &internal.HTTPError{StatusCode: http.StatusConflict, Message: "professor has a conflicting schedule", Err: err}
		}
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "could not create semester course", Err: err}
	}

	internal.WriteData(w, "", res, nil)
	return nil
}
