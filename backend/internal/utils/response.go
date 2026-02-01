// Package utils provides shared utility functions.
package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/urinvitedto-my/backend/internal/models"
)

// WriteJSON writes a JSON response with the given status code.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("Error encoding JSON", "error", err)
	}
}

// WriteError writes a JSON error response.
func WriteError(w http.ResponseWriter, status int, err, msg string) {
	WriteJSON(w, status, models.ErrorResponse{Error: err, Message: msg})
}
