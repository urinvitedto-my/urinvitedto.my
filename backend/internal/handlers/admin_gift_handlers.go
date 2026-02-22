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

// ListGifts handles GET /admin/events/:id/gifts - returns all gifts.
func (h *Handlers) ListGifts(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	ctx := r.Context()

	var exists bool
	if err := h.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM events WHERE id = $1)`, eventID).Scan(&exists); err != nil || !exists {
		h.writeError(w, http.StatusNotFound, "not_found", "Event not found")
		return
	}

	rows, err := h.db.Query(ctx, `
		SELECT id, gift_type, title, description, link, order_index, created_at
		FROM event_gifts WHERE event_id = $1
		ORDER BY order_index ASC
	`, eventID)
	if err != nil {
		slog.Error("DB error listing gifts", "error", err)
		h.writeError(w, http.StatusInternalServerError, "db_error", "Failed to list gifts")
		return
	}
	defer rows.Close()

	items := []models.AdminGift{}
	for rows.Next() {
		var item models.AdminGift
		if err := rows.Scan(&item.ID, &item.GiftType, &item.Title, &item.Description, &item.Link, &item.OrderIndex, &item.CreatedAt); err != nil {
			slog.Error("Error scanning gift", "error", err)
			continue
		}
		items = append(items, item)
	}

	h.writeJSON(w, http.StatusOK, models.AdminGiftsResponse{Items: items})
}

// CreateGift handles POST /admin/events/:id/gifts - creates a gift.
func (h *Handlers) CreateGift(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")

	var req models.CreateGiftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	req.GiftType = strings.ToLower(strings.TrimSpace(req.GiftType))
	req.Title = strings.TrimSpace(req.Title)
	if req.GiftType != "physical" && req.GiftType != "monetary" {
		h.writeError(w, http.StatusBadRequest, "invalid_type", "Gift type must be physical or monetary")
		return
	}
	if req.Title == "" {
		h.writeError(w, http.StatusBadRequest, "invalid_title", "Title is required")
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
		_ = h.db.QueryRow(ctx, `SELECT MAX(order_index) FROM event_gifts WHERE event_id = $1`, eventID).Scan(&maxIdx)
		if maxIdx != nil {
			orderIndex = *maxIdx + 1
		}
	}

	var item models.AdminGift
	err := h.db.QueryRow(ctx, `
		INSERT INTO event_gifts (event_id, gift_type, title, description, link, order_index)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, gift_type, title, description, link, order_index, created_at
	`, eventID, req.GiftType, req.Title, req.Description, req.Link, orderIndex).Scan(
		&item.ID, &item.GiftType, &item.Title, &item.Description, &item.Link, &item.OrderIndex, &item.CreatedAt,
	)
	if err != nil {
		slog.Error("DB error creating gift", "error", err)
		h.writeError(w, http.StatusInternalServerError, "db_error", "Failed to create gift")
		return
	}

	h.writeJSON(w, http.StatusCreated, item)
}

// UpdateGift handles PUT /admin/events/:id/gifts/:itemId - updates a gift.
func (h *Handlers) UpdateGift(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	itemID := chi.URLParam(r, "itemId")

	var req models.UpdateGiftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	req.GiftType = strings.ToLower(strings.TrimSpace(req.GiftType))
	req.Title = strings.TrimSpace(req.Title)
	if req.GiftType != "physical" && req.GiftType != "monetary" {
		h.writeError(w, http.StatusBadRequest, "invalid_type", "Gift type must be physical or monetary")
		return
	}
	if req.Title == "" {
		h.writeError(w, http.StatusBadRequest, "invalid_title", "Title is required")
		return
	}

	orderIndex := 0
	if req.OrderIndex != nil {
		orderIndex = *req.OrderIndex
	}

	ctx := r.Context()

	var item models.AdminGift
	err := h.db.QueryRow(ctx, `
		UPDATE event_gifts
		SET gift_type = $1, title = $2, description = $3, link = $4, order_index = $5
		WHERE id = $6 AND event_id = $7
		RETURNING id, gift_type, title, description, link, order_index, created_at
	`, req.GiftType, req.Title, req.Description, req.Link, orderIndex, itemID, eventID).Scan(
		&item.ID, &item.GiftType, &item.Title, &item.Description, &item.Link, &item.OrderIndex, &item.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "not_found", "Gift not found")
			return
		}
		slog.Error("DB error updating gift", "error", err)
		h.writeError(w, http.StatusInternalServerError, "db_error", "Failed to update gift")
		return
	}

	h.writeJSON(w, http.StatusOK, item)
}

// DeleteGift handles DELETE /admin/events/:id/gifts/:itemId - removes a gift.
func (h *Handlers) DeleteGift(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	itemID := chi.URLParam(r, "itemId")
	ctx := r.Context()

	result, err := h.db.Exec(ctx,
		`DELETE FROM event_gifts WHERE id = $1 AND event_id = $2`,
		itemID, eventID,
	)
	if err != nil {
		slog.Error("DB error deleting gift", "error", err)
		h.writeError(w, http.StatusInternalServerError, "db_error", "Failed to delete gift")
		return
	}

	if result.RowsAffected() == 0 {
		h.writeError(w, http.StatusNotFound, "not_found", "Gift not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
