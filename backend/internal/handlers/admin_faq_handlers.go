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

// ListFAQs handles GET /admin/events/:id/faqs - returns all FAQs.
func (h *Handlers) ListFAQs(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	ctx := r.Context()

	var exists bool
	if err := h.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM events WHERE id = $1)`, eventID).Scan(&exists); err != nil || !exists {
		h.writeError(w, http.StatusNotFound, "not_found", "Event not found")
		return
	}

	rows, err := h.db.Query(ctx, `
		SELECT id, question, answer, order_index, created_at
		FROM event_faqs WHERE event_id = $1
		ORDER BY order_index ASC
	`, eventID)
	if err != nil {
		slog.Error("DB error listing FAQs", "error", err)
		h.writeError(w, http.StatusInternalServerError, "db_error", "Failed to list FAQs")
		return
	}
	defer rows.Close()

	items := []models.AdminFAQ{}
	for rows.Next() {
		var item models.AdminFAQ
		if err := rows.Scan(&item.ID, &item.Question, &item.Answer, &item.OrderIndex, &item.CreatedAt); err != nil {
			slog.Error("Error scanning FAQ", "error", err)
			continue
		}
		items = append(items, item)
	}

	h.writeJSON(w, http.StatusOK, models.AdminFAQsResponse{Items: items})
}

// CreateFAQ handles POST /admin/events/:id/faqs - creates a FAQ.
func (h *Handlers) CreateFAQ(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")

	var req models.CreateFAQRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	req.Question = strings.TrimSpace(req.Question)
	req.Answer = strings.TrimSpace(req.Answer)
	if req.Question == "" {
		h.writeError(w, http.StatusBadRequest, "invalid_question", "Question is required")
		return
	}
	if req.Answer == "" {
		h.writeError(w, http.StatusBadRequest, "invalid_answer", "Answer is required")
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
		_ = h.db.QueryRow(ctx, `SELECT MAX(order_index) FROM event_faqs WHERE event_id = $1`, eventID).Scan(&maxIdx)
		if maxIdx != nil {
			orderIndex = *maxIdx + 1
		}
	}

	var item models.AdminFAQ
	err := h.db.QueryRow(ctx, `
		INSERT INTO event_faqs (event_id, question, answer, order_index)
		VALUES ($1, $2, $3, $4)
		RETURNING id, question, answer, order_index, created_at
	`, eventID, req.Question, req.Answer, orderIndex).Scan(
		&item.ID, &item.Question, &item.Answer, &item.OrderIndex, &item.CreatedAt,
	)
	if err != nil {
		slog.Error("DB error creating FAQ", "error", err)
		h.writeError(w, http.StatusInternalServerError, "db_error", "Failed to create FAQ")
		return
	}

	h.writeJSON(w, http.StatusCreated, item)
}

// UpdateFAQ handles PUT /admin/events/:id/faqs/:itemId - updates a FAQ.
func (h *Handlers) UpdateFAQ(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	itemID := chi.URLParam(r, "itemId")

	var req models.UpdateFAQRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	req.Question = strings.TrimSpace(req.Question)
	req.Answer = strings.TrimSpace(req.Answer)
	if req.Question == "" {
		h.writeError(w, http.StatusBadRequest, "invalid_question", "Question is required")
		return
	}
	if req.Answer == "" {
		h.writeError(w, http.StatusBadRequest, "invalid_answer", "Answer is required")
		return
	}

	orderIndex := 0
	if req.OrderIndex != nil {
		orderIndex = *req.OrderIndex
	}

	ctx := r.Context()

	var item models.AdminFAQ
	err := h.db.QueryRow(ctx, `
		UPDATE event_faqs
		SET question = $1, answer = $2, order_index = $3
		WHERE id = $4 AND event_id = $5
		RETURNING id, question, answer, order_index, created_at
	`, req.Question, req.Answer, orderIndex, itemID, eventID).Scan(
		&item.ID, &item.Question, &item.Answer, &item.OrderIndex, &item.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "not_found", "FAQ not found")
			return
		}
		slog.Error("DB error updating FAQ", "error", err)
		h.writeError(w, http.StatusInternalServerError, "db_error", "Failed to update FAQ")
		return
	}

	h.writeJSON(w, http.StatusOK, item)
}

// DeleteFAQ handles DELETE /admin/events/:id/faqs/:itemId - removes a FAQ.
func (h *Handlers) DeleteFAQ(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	itemID := chi.URLParam(r, "itemId")
	ctx := r.Context()

	result, err := h.db.Exec(ctx,
		`DELETE FROM event_faqs WHERE id = $1 AND event_id = $2`,
		itemID, eventID,
	)
	if err != nil {
		slog.Error("DB error deleting FAQ", "error", err)
		h.writeError(w, http.StatusInternalServerError, "db_error", "Failed to delete FAQ")
		return
	}

	if result.RowsAffected() == 0 {
		h.writeError(w, http.StatusNotFound, "not_found", "FAQ not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
