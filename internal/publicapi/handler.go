package publicapi

import (
	"educore/internal"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// HandleVerifyDepartmentEnrollment godoc
// @Summary      Verify student department enrollment
// @Description  Peer API: verify whether a student belongs to a specific department. Requires x-api-key header.
// @Tags         public-api
// @Produce      json
// @Security     ApiKeyAuth
// @Param        student_id      path  string  true  "Student ID"
// @Param        department_name path  string  true  "Department name"
// @Success      200  {object}  internal.ResponseBody{data=VerifyDepartmentResponse}
// @Failure      401  {object}  internal.ResponseBody
// @Router       /public/students/{student_id}/departments/{department_name} [get]
func (h *Handler) HandleVerifyDepartmentEnrollment(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	studentID := r.PathValue("student_id")
	departmentName := r.PathValue("department_name")

	enrolled, err := IsStudentInDepartment(h.db, studentID, departmentName)
	if err != nil {
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "could not verify enrollment", Err: err}
	}

	internal.WriteData(w, "", VerifyDepartmentResponse{
		StudentID:  studentID,
		Department: departmentName,
		Enrolled:   enrolled,
	}, nil)
	return nil
}

// HandleGrantKey godoc
// @Summary      Grant a new API key
// @Description  Admin only: generate and issue a new API key for peer service integration.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
//
//	@Param        body  body  GrantKeyRequest  true  "Partner details"
//
// @Success      200  {object}  internal.ResponseBody{data=GrantKeyResponse}
// @Failure      400  {object}  internal.ResponseBody
// @Router       /admin/api-keys [post]
func (h *Handler) HandleGrantKey(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	return h.grantKey(w, r)
}

func (h *Handler) grantKey(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	var req GrantKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &internal.HTTPError{StatusCode: http.StatusBadRequest, Message: "invalid request body", Err: err}
	}

	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Email) == "" {
		return &internal.HTTPError{StatusCode: http.StatusBadRequest, Message: "name and email are required", Err: nil}
	}

	key := uuid.New().String()
	if err := h.db.Create(&internal.ApiKey{Name: req.Name, Email: req.Email, Key: key, Valid: true}).Error; err != nil {
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "could not create api key", Err: err}
	}

	internal.WriteData(w, "", GrantKeyResponse{Name: req.Name, Email: req.Email, Key: key}, nil)
	return nil
}

// HandleRevokeKey godoc
// @Summary      Revoke an API key
// @Description  Admin only: revoke an existing API key by its key string.
// @Tags         admin
// @Produce      json
// @Security     BearerAuth
// @Param        key  path  string  true  "API key to revoke"
// @Success      200  {object}  internal.ResponseBody
// @Failure      400  {object}  internal.ResponseBody
// @Router       /admin/api-keys/{key} [delete]
func (h *Handler) HandleRevokeKey(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	key := r.PathValue("key")
	if key == "" {
		return &internal.HTTPError{StatusCode: http.StatusBadRequest, Message: "key is required", Err: nil}
	}

	if err := h.db.Model(&internal.ApiKey{}).Where("key = ?", key).Update("valid", false).Error; err != nil {
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "could not revoke api key", Err: err}
	}

	internal.WriteData(w, "", "api key revoked", nil)
	return nil
}
