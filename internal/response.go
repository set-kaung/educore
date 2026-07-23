package internal

import (
	"encoding/json"
	"net/http"
)

type ResponseBody struct {
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	Message    string `json:"message,omitempty"`
	Data       any    `json:"data,omitempty"`
}

func writeResponse(w http.ResponseWriter, statusCode int, message string, data any, headers http.Header) error {
	body := ResponseBody{
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
		Message:    message,
		Data:       data,
	}
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(body)
}

func WriteServerError(w http.ResponseWriter, headers http.Header) error {
	return writeResponse(w, http.StatusInternalServerError, "", nil, headers)
}

func WriteData(w http.ResponseWriter, messsage string, data any, headers http.Header) error {
	return writeResponse(w, http.StatusOK, messsage, data, headers)
}

func WriteError(w http.ResponseWriter, statusCode int, message string, headers http.Header) error {
	return writeResponse(w, statusCode, message, nil, headers)
}
