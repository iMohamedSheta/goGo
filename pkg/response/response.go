package response

import (
	"encoding/json"
	"imohamedsheta/gocrud/pkg/logger"
	"net/http"
)

// Struct for standard JSON response
type Response struct {
	Message   string `json:"message"`
	Data      any    `json:"data,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
}

// writes a JSON response with proper headers and error handling
func Json(w http.ResponseWriter, message string, data any, code int) {
	resp := &Response{
		Message: message,
		Data:    data,
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Log().Error("Failed to write response: " + err.Error())

		// If encoding the original response failed, send a fallback JSON error
		ServerErrorJson(w)
		return
	}
}

func ServerErrorJson(w http.ResponseWriter) {
	ErrorJson(w, "Internal Server Error", "server_error", http.StatusInternalServerError)
}

func ErrorJson(w http.ResponseWriter, message string, errorCode string, code int) {
	resp := &Response{
		ErrorCode: errorCode,
		Message:   message,
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Log().Error("Failed to write response:  " + err.Error())

		// If encoding the original response failed, send a fallback JSON error
		ServerErrorJson(w)
		return
	}
}

func ValidationErrorJson(w http.ResponseWriter, validationErrors map[string]string) {
	resp := &Response{
		ErrorCode: "validation_error",
		Message:   "Validation Error",
		Data:      validationErrors,
	}

	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Log().Error("Failed to write response:  " + err.Error())

		// If encoding the original response failed, send a fallback JSON error
		ServerErrorJson(w)
		return
	}
}
