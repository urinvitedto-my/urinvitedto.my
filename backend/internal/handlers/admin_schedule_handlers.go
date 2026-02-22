package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/urinvitedto-my/backend/internal/models"
)

// ListSchedule handles GET /admin/events/:id/schedule - returns all schedule items.
func (h *Handlers) ListSchedule(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	ctx := r.Context()

	var exists bool
	if err := h.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM events WHERE id = $1)`, eventID).Scan(&exists); err != nil || !exists {
		h.writeError(w, http.StatusNotFound, "not_found", "Event not found")
		return
	}

	rows, err := h.db.Query(ctx, `
		SELECT id, time, title, description, order_index, created_at
		FROM event_schedule WHERE event_id = $1
		ORDER BY order_index ASC, time ASC
	`, eventID)
	if err != nil {
		slog.Error("DB error listing schedule", "error", err)
		h.writeError(w, http.StatusInternalServerError, "db_error", "Failed to list schedule")
		return
	}
	defer rows.Close()

	items := []models.AdminScheduleItem{}
	for rows.Next() {
		var item models.AdminScheduleItem
		if err := rows.Scan(&item.ID, &item.Time, &item.Title, &item.Description, &item.OrderIndex, &item.CreatedAt); err != nil {
			slog.Error("Error scanning schedule item", "error", err)
			continue
		}
		items = append(items, item)
	}

	h.writeJSON(w, http.StatusOK, models.AdminScheduleResponse{Items: items})
}

// CreateScheduleItem handles POST /admin/events/:id/schedule - creates a schedule item.
func (h *Handlers) CreateScheduleItem(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")

	var req models.CreateScheduleItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		h.writeError(w, http.StatusBadRequest, "invalid_title", "Title is required")
		return
	}

	parsedTime, err := time.Parse(time.RFC3339, req.Time)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_time", "Time must be a valid ISO 8601 date string")
		return
	}

	ctx := r.Context()

	var exists bool
	if err := h.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM events WHERE id = $1)`, eventID).Scan(&exists); err != nil || !exists {
		h.writeError(w, http.StatusNotFound, "not_found", "Event not found")
		return
	}

	// auto-assign order_index if not provided
	orderIndex := 0
	if req.OrderIndex != nil {
		orderIndex = *req.OrderIndex
	} else {
		var maxIdx *int
		_ = h.db.QueryRow(ctx, `SELECT MAX(order_index) FROM event_schedule WHERE event_id = $1`, eventID).Scan(&maxIdx)
		if maxIdx != nil {
			orderIndex = *maxIdx + 1
		}
	}

	var item models.AdminScheduleItem
	err = h.db.QueryRow(ctx, `
		INSERT INTO event_schedule (event_id, time, title, description, order_index)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, time, title, description, order_index, created_at
	`, eventID, parsedTime, req.Title, req.Description, orderIndex).Scan(
		&item.ID, &item.Time, &item.Title, &item.Description, &item.OrderIndex, &item.CreatedAt,
	)
	if err != nil {
		slog.Error("DB error creating schedule item", "error", err)
		h.writeError(w, http.StatusInternalServerError, "db_error", "Failed to create schedule item")
		return
	}

	h.writeJSON(w, http.StatusCreated, item)
}

// UpdateScheduleItem handles PUT /admin/events/:id/schedule/:itemId - updates a schedule item.
func (h *Handlers) UpdateScheduleItem(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	itemID := chi.URLParam(r, "itemId")

	var req models.UpdateScheduleItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		h.writeError(w, http.StatusBadRequest, "invalid_title", "Title is required")
		return
	}

	parsedTime, err := time.Parse(time.RFC3339, req.Time)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_time", "Time must be a valid ISO 8601 date string")
		return
	}

	ctx := r.Context()

	// default to 0 if not provided
	orderIndex := 0
	if req.OrderIndex != nil {
		orderIndex = *req.OrderIndex
	}

	var item models.AdminScheduleItem
	err = h.db.QueryRow(ctx, `
		UPDATE event_schedule
		SET time = $1, title = $2, description = $3, order_index = $4
		WHERE id = $5 AND event_id = $6
		RETURNING id, time, title, description, order_index, created_at
	`, parsedTime, req.Title, req.Description, orderIndex, itemID, eventID).Scan(
		&item.ID, &item.Time, &item.Title, &item.Description, &item.OrderIndex, &item.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "not_found", "Schedule item not found")
			return
		}
		slog.Error("DB error updating schedule item", "error", err)
		h.writeError(w, http.StatusInternalServerError, "db_error", "Failed to update schedule item")
		return
	}

	h.writeJSON(w, http.StatusOK, item)
}

// DeleteScheduleItem handles DELETE /admin/events/:id/schedule/:itemId - removes a schedule item.
func (h *Handlers) DeleteScheduleItem(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	itemID := chi.URLParam(r, "itemId")
	ctx := r.Context()

	result, err := h.db.Exec(ctx,
		`DELETE FROM event_schedule WHERE id = $1 AND event_id = $2`,
		itemID, eventID,
	)
	if err != nil {
		slog.Error("DB error deleting schedule item", "error", err)
		h.writeError(w, http.StatusInternalServerError, "db_error", "Failed to delete schedule item")
		return
	}

	if result.RowsAffected() == 0 {
		h.writeError(w, http.StatusNotFound, "not_found", "Schedule item not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
