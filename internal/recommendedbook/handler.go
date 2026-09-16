package recommendedbook

import (
	"educore/internal"
	"educore/internal/auth"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// HandleList godoc
// @Summary      List recommended books for a course offering
// @Description  Get all books a professor has recommended for the given semester course
// @Tags         recommended-books
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "Semester course ID"
// @Success      200  {object}  internal.ResponseBody{data=[]Book}
// @Router       /api/semester-courses/{id}/books [get]
func (h *Handler) HandleList(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	semesterCourseID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || semesterCourseID == 0 {
		return &internal.HTTPError{StatusCode: http.StatusBadRequest, Message: "invalid semester course id", Err: nil}
	}

	books, err := ListBooks(h.db, uint(semesterCourseID))
	if err != nil {
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "could not list recommended books", Err: err}
	}

	internal.WriteData(w, "", books, nil)
	return nil
}

// HandleAdd godoc
// @Summary      Add a recommended book to a course offering
// @Description  Professor adds a book to their course offering's recommended reading list
// @Tags         recommended-books
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  int         true  "Semester course ID"
// @Param        body  body  AddRequest  true  "Book details"
// @Success      200   {object}  internal.ResponseBody{data=Book}
// @Failure      400   {object}  internal.ResponseBody
// @Failure      403   {object}  internal.ResponseBody
// @Failure      404   {object}  internal.ResponseBody
// @Failure      409   {object}  internal.ResponseBody
// @Router       /api/semester-courses/{id}/books [post]
func (h *Handler) HandleAdd(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	claims := auth.GetClaims(r)

	semesterCourseID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || semesterCourseID == 0 {
		return &internal.HTTPError{StatusCode: http.StatusBadRequest, Message: "invalid semester course id", Err: nil}
	}

	var req AddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &internal.HTTPError{StatusCode: http.StatusBadRequest, Message: "invalid request body", Err: err}
	}

	book, err := AddBook(h.db, claims.UserID, uint(semesterCourseID), req)
	if err != nil {
		if errors.Is(err, ErrSemesterCourseNotFound) {
			return &internal.HTTPError{StatusCode: http.StatusNotFound, Message: "course offering not found", Err: err}
		}
		if errors.Is(err, ErrNotTeacher) {
			return &internal.HTTPError{StatusCode: http.StatusForbidden, Message: "only the professor teaching this course may manage its books", Err: err}
		}
		if errors.Is(err, ErrBookExists) {
			return &internal.HTTPError{StatusCode: http.StatusConflict, Message: "book is already on the recommended reading list", Err: err}
		}
		if err.Error() == "book title is required" {
			return &internal.HTTPError{StatusCode: http.StatusBadRequest, Message: "book title is required", Err: err}
		}
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "could not add recommended book", Err: err}
	}

	internal.WriteData(w, "book added to recommended reading list", book, nil)
	return nil
}

// HandleRemove godoc
// @Summary      Remove a recommended book
// @Description  Professor removes a recommended book from their course offering
// @Tags         recommended-books
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "Recommended book ID"
// @Success      200  {object}  internal.ResponseBody
// @Failure      403  {object}  internal.ResponseBody
// @Failure      404  {object}  internal.ResponseBody
// @Router       /api/books/{id} [delete]
func (h *Handler) HandleRemove(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	claims := auth.GetClaims(r)

	bookID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || bookID == 0 {
		return &internal.HTTPError{StatusCode: http.StatusBadRequest, Message: "invalid book id", Err: nil}
	}

	err = RemoveBook(h.db, claims.UserID, uint(bookID))
	if err != nil {
		if errors.Is(err, ErrBookNotFound) {
			return &internal.HTTPError{StatusCode: http.StatusNotFound, Message: "recommended book not found", Err: err}
		}
		if errors.Is(err, ErrSemesterCourseNotFound) {
			return &internal.HTTPError{StatusCode: http.StatusNotFound, Message: "course offering not found", Err: err}
		}
		if errors.Is(err, ErrNotTeacher) {
			return &internal.HTTPError{StatusCode: http.StatusForbidden, Message: "only the professor teaching this course may manage its books", Err: err}
		}
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "could not remove recommended book", Err: err}
	}

	internal.WriteData(w, "book removed from recommended reading list", nil, nil)
	return nil
}
