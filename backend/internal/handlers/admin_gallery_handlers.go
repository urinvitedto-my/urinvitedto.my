package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/urinvitedto-my/backend/internal/models"
)

// ListGallery handles GET /admin/events/:id/gallery - returns all gallery items.
func (h *Handlers) ListGallery(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	ctx := r.Context()

	var exists bool
	if err := h.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM events WHERE id = $1)`, eventID).Scan(&exists); err != nil || !exists {
		h.writeError(w, http.StatusNotFound, "not_found", "Event not found")
		return
	}

	rows, err := h.db.Query(ctx, `
		SELECT id, media_type, media_url, caption, order_index, created_at
		FROM event_gallery WHERE event_id = $1
		ORDER BY order_index ASC
	`, eventID)
	if err != nil {
		slog.Error("DB error listing gallery", "error", err)
		h.writeError(w, http.StatusInternalServerError, "db_error", "Failed to list gallery")
		return
	}
	defer rows.Close()

	items := []models.AdminGalleryItem{}
	for rows.Next() {
		var item models.AdminGalleryItem
		if err := rows.Scan(&item.ID, &item.MediaType, &item.MediaURL, &item.Caption, &item.OrderIndex, &item.CreatedAt); err != nil {
			slog.Error("Error scanning gallery item", "error", err)
			continue
		}
		items = append(items, item)
	}

	h.writeJSON(w, http.StatusOK, models.AdminGalleryResponse{Items: items})
}

// CreateGalleryItem handles POST /admin/events/:id/gallery - creates a gallery item.
func (h *Handlers) CreateGalleryItem(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")

	var req models.CreateGalleryItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	req.MediaType = strings.ToLower(strings.TrimSpace(req.MediaType))
	req.MediaURL = strings.TrimSpace(req.MediaURL)
	if req.MediaType != "photo" && req.MediaType != "video" {
		h.writeError(w, http.StatusBadRequest, "invalid_type", "Media type must be photo or video")
		return
	}
	if req.MediaURL == "" {
		h.writeError(w, http.StatusBadRequest, "invalid_url", "Media URL is required")
		return
	}

	ctx := r.Context()

	var exists bool
	if err := h.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM events WHERE id = $1)`, eventID).Scan(&exists); err != nil || !exists {
		h.writeError(w, http.StatusNotFound, "not_found", "Event not found")
		return
	}

	orderIndex := 0
	if req.OrderIndex != nil {
		orderIndex = *req.OrderIndex
	} else {
		var maxIdx *int
		_ = h.db.QueryRow(ctx, `SELECT MAX(order_index) FROM event_gallery WHERE event_id = $1`, eventID).Scan(&maxIdx)
		if maxIdx != nil {
			orderIndex = *maxIdx + 1
		}
	}

	var item models.AdminGalleryItem
	err := h.db.QueryRow(ctx, `
		INSERT INTO event_gallery (event_id, media_type, media_url, caption, order_index)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, media_type, media_url, caption, order_index, created_at
	`, eventID, req.MediaType, req.MediaURL, req.Caption, orderIndex).Scan(
		&item.ID, &item.MediaType, &item.MediaURL, &item.Caption, &item.OrderIndex, &item.CreatedAt,
	)
	if err != nil {
		slog.Error("DB error creating gallery item", "error", err)
		h.writeError(w, http.StatusInternalServerError, "db_error", "Failed to create gallery item")
		return
	}

	h.writeJSON(w, http.StatusCreated, item)
}

// UpdateGalleryItem handles PUT /admin/events/:id/gallery/:itemId - updates a gallery item.
func (h *Handlers) UpdateGalleryItem(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	itemID := chi.URLParam(r, "itemId")

	var req models.UpdateGalleryItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	orderIndex := 0
	if req.OrderIndex != nil {
		orderIndex = *req.OrderIndex
	}

	ctx := r.Context()

	var item models.AdminGalleryItem
	err := h.db.QueryRow(ctx, `
		UPDATE event_gallery
		SET caption = $1, order_index = $2
		WHERE id = $3 AND event_id = $4
		RETURNING id, media_type, media_url, caption, order_index, created_at
	`, req.Caption, orderIndex, itemID, eventID).Scan(
		&item.ID, &item.MediaType, &item.MediaURL, &item.Caption, &item.OrderIndex, &item.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "not_found", "Gallery item not found")
			return
		}
		slog.Error("DB error updating gallery item", "error", err)
		h.writeError(w, http.StatusInternalServerError, "db_error", "Failed to update gallery item")
		return
	}

	h.writeJSON(w, http.StatusOK, item)
}

// DeleteGalleryItem handles DELETE /admin/events/:id/gallery/:itemId - removes a gallery item.
func (h *Handlers) DeleteGalleryItem(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	itemID := chi.URLParam(r, "itemId")
	ctx := r.Context()

	result, err := h.db.Exec(ctx,
		`DELETE FROM event_gallery WHERE id = $1 AND event_id = $2`,
		itemID, eventID,
	)
	if err != nil {
		slog.Error("DB error deleting gallery item", "error", err)
		h.writeError(w, http.StatusInternalServerError, "db_error", "Failed to delete gallery item")
		return
	}

	if result.RowsAffected() == 0 {
		h.writeError(w, http.StatusNotFound, "not_found", "Gallery item not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
