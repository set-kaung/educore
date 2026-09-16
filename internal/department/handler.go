package department

import (
	"educore/internal"
	"net/http"
)

// HandleListDepartments godoc
// @Summary      List departments
// @Description  Get all departments, seeding a default set when none exist
// @Tags         departments
// @Produce      json
// @Success      200  {object}  internal.ResponseBody{data=[]DepartmentData}
// @Failure      500  {object}  internal.ResponseBody
// @Router       /api/departments [get]
func (h *Handler) HandleListDepartments(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	departments, err := ListDepartments(h.db)
	if err != nil {
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "could not list departments", Err: err}
	}

	internal.WriteData(w, "", departments, nil)
	return nil
}
