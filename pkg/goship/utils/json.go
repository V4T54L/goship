package utils

import (
	"encoding/json"
	"net/http"
)

// WriteJSONError is a helper to write a standardized JSON error response.
func WriteJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
