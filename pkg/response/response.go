package response

import (
	"encoding/json"
	"imohamedsheta/gocrud/pkg/logger"
	"net/http"
)

// Struct for standard JSON response
type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// writes a JSON response with proper headers and error handling
func Json(w http.ResponseWriter, message string, data any, code int) {
	resp := &Response{
		Message: message,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Log().Error("Failed to write response: " + err.Error())

		// If encoding the original response failed, send a fallback JSON error
		ServerErrorJson(w)
		return
	}

	w.WriteHeader(code)
}

func ServerErrorJson(w http.ResponseWriter) {
	ErrorJson(w, "Server Error", http.StatusInternalServerError)
}

func ErrorJson(w http.ResponseWriter, message string, code int) {
	resp := &Response{
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Log().Error("Failed to write response:  " + err.Error())

		// If encoding the original response failed, send a fallback JSON error
		ServerErrorJson(w)
		return
	}

	w.WriteHeader(code)
}

func ValidationErrorJson(w http.ResponseWriter, validationErrors map[string]string) {
	Json(w, "Validation Error", validationErrors, http.StatusUnprocessableEntity)
}
