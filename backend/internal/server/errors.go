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
		Error:   message,
		Code:    code,
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
