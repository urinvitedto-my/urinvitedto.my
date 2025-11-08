// Package handlers
package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Handlers bundles shared deps for HTTP handlers.
type Handlers struct {
	db *pgxpool.Pool
}

// New creates a new Handlers instance.
func New(db *pgxpool.Pool) *Handlers {
	return &Handlers{db: db}
}

// writeJSON writes a JSON response to the response writer.
func (h *Handlers) writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)

	err := enc.Encode(v)
	if err != nil {
		slog.Error("Error encoding JSON", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
