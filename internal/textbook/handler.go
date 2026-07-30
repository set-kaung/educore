package textbook

import (
	"educore/internal"
	"educore/internal/openlibrary"
	"net/http"
)

type TextbookHandler struct {
	lib *openlibrary.Client
}

func NewTextbookHandler(userAgent string) *TextbookHandler {
	return &TextbookHandler{lib: openlibrary.NewClient(userAgent)}
}

func (h *TextbookHandler) HandleSearch(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	query := r.URL.Query().Get("q")
	if query == "" {
		return &internal.HTTPError{StatusCode: http.StatusBadRequest, Message: "query parameter 'q' is required", Err: nil}
	}

	result, err := h.lib.Search(query)
	if err != nil {
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "failed to search textbooks", Err: err}
	}

	internal.WriteData(w, "success", result, nil)
	return nil
}
