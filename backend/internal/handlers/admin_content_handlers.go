package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

// GetCustomContent handles GET /admin/events/:id/custom-content.
func (h *Handlers) GetCustomContent(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	ctx := r.Context()

	var raw json.RawMessage
	err := h.db.QueryRow(ctx,
		`SELECT COALESCE(custom_content, '{}') FROM events WHERE id = $1`, eventID,
	).Scan(&raw)
	if err != nil {
		if err == pgx.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "not_found", "Event not found")
			return
		}
		slog.Error("DB error reading custom_content", "error", err)
		h.writeError(w, http.StatusInternalServerError, "db_error", "Failed to read custom content")
		return
	}

	h.writeJSON(w, http.StatusOK, raw)
}

// UpdateCustomContent handles PUT /admin/events/:id/custom-content.
func (h *Handlers) UpdateCustomContent(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")

	var body json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_body", "Invalid JSON")
		return
	}

	ctx := r.Context()

	var saved json.RawMessage
	err := h.db.QueryRow(ctx,
		`UPDATE events SET custom_content = $1 WHERE id = $2
		 RETURNING custom_content`, body, eventID,
	).Scan(&saved)
	if err != nil {
		if err == pgx.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "not_found", "Event not found")
			return
		}
		slog.Error("DB error updating custom_content", "error", err)
		h.writeError(w, http.StatusInternalServerError, "db_error", "Failed to update custom content")
		return
	}

	h.writeJSON(w, http.StatusOK, saved)
}

// GetEnabledComponents handles GET /admin/events/:id/enabled-components.
func (h *Handlers) GetEnabledComponents(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	ctx := r.Context()

	var raw json.RawMessage
	err := h.db.QueryRow(ctx,
		`SELECT COALESCE(enabled_components, '{}') FROM events WHERE id = $1`, eventID,
	).Scan(&raw)
	if err != nil {
		if err == pgx.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "not_found", "Event not found")
			return
		}
		slog.Error("DB error reading enabled_components", "error", err)
		h.writeError(w, http.StatusInternalServerError, "db_error", "Failed to read enabled components")
		return
	}

	h.writeJSON(w, http.StatusOK, raw)
}

// UpdateEnabledComponents handles PUT /admin/events/:id/enabled-components.
func (h *Handlers) UpdateEnabledComponents(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")

	var body json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_body", "Invalid JSON")
		return
	}

	ctx := r.Context()

	var saved json.RawMessage
	err := h.db.QueryRow(ctx,
		`UPDATE events SET enabled_components = $1 WHERE id = $2
		 RETURNING enabled_components`, body, eventID,
	).Scan(&saved)
	if err != nil {
		if err == pgx.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "not_found", "Event not found")
			return
		}
		slog.Error("DB error updating enabled_components", "error", err)
		h.writeError(w, http.StatusInternalServerError, "db_error", "Failed to update enabled components")
		return
	}

	h.writeJSON(w, http.StatusOK, saved)
}
