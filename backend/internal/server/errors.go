package server

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorResponse represents a standardized API error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}

// respondWithError sends a standardized error response
func respondWithError(w http.ResponseWriter, code int, message string, err error) {
	log.Printf("API Error [%d]: %s - %v", code, message, err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	errResp := ErrorResponse{
		Error: message,
		Code:  code,
	}

	if err != nil {
		errResp.Message = err.Error()
	}

	json.NewEncoder(w).Encode(errResp)
}

// respondWithJSON sends a successful JSON response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// respondWithNotFound sends a standardized 404 Not Found response
// resource: the type of resource that was not found (e.g., "download", "config")
func respondWithNotFound(w http.ResponseWriter, resource string, err error) {
	message := "Resource not found: " + resource
	respondWithError(w, http.StatusNotFound, message, err)
}

// respondWithBadRequest sends a standardized 400 Bad Request response
// reason: explanation of why the request was invalid
func respondWithBadRequest(w http.ResponseWriter, reason string, err error) {
	respondWithError(w, http.StatusBadRequest, reason, err)
}

// respondWithValidationError sends a standardized validation error response (400)
// field: the field that failed validation
func respondWithValidationError(w http.ResponseWriter, field string, err error) {
	message := "Validation failed: " + field
	respondWithError(w, http.StatusBadRequest, message, err)
}

// respondWithInternalError sends a standardized 500 Internal Server Error response
// operation: the operation that failed (e.g., "create download", "list downloads")
func respondWithInternalError(w http.ResponseWriter, operation string, err error) {
	message := "Failed to " + operation
	respondWithError(w, http.StatusInternalServerError, message, err)
}
