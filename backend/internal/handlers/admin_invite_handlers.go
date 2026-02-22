package handlers

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/urinvitedto-my/backend/internal/models"
)

const inviteCodeChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const inviteCodeLength = 6
const maxCodeRetries = 10

// generateInviteCode creates a random 6-char alphanumeric code (A-Z, 0-9).
func generateInviteCode() (string, error) {
	code := make([]byte, inviteCodeLength)
	for i := range code {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(inviteCodeChars))))
		if err != nil {
			return "", err
		}
		code[i] = inviteCodeChars[idx.Int64()]
	}
	return string(code), nil
}

// ListInvites handles GET /admin/events/:id/invites - returns all invites with guests.
func (h *Handlers) ListInvites(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	ctx := r.Context()

	// verify event exists
	var exists bool
	if err := h.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM events WHERE id = $1)`, eventID).Scan(&exists); err != nil || !exists {
		h.writeError(w, http.StatusNotFound, "not_found", "Event not found")
		return
	}

	rows, err := h.db.Query(ctx, `
		SELECT id, invite_code, label, created_at
		FROM invites WHERE event_id = $1
		ORDER BY created_at DESC
	`, eventID)
	if err != nil {
		slog.Error("DB error listing invites", "error", err)
		h.writeError(w, http.StatusInternalServerError, "db_error", "Failed to list invites")
		return
	}
	defer rows.Close()

	invites := []models.AdminInvite{}
	for rows.Next() {
		var inv models.AdminInvite
		if err := rows.Scan(&inv.ID, &inv.InviteCode, &inv.Label, &inv.CreatedAt); err != nil {
			slog.Error("Error scanning invite", "error", err)
			continue
		}
		invites = append(invites, inv)
	}

	// fetch guests for each invite
	for i := range invites {
		guests, err := h.fetchInviteGuests(ctx, invites[i].ID)
		if err != nil {
			slog.Error("Error fetching guests", "inviteId", invites[i].ID, "error", err)
		}
		invites[i].Guests = guests
	}

	h.writeJSON(w, http.StatusOK, models.AdminInvitesResponse{Invites: invites})
}

// CreateInvite handles POST /admin/events/:id/invites - creates an invite with auto-generated code.
func (h *Handlers) CreateInvite(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")

	var req models.CreateInviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// body is optional (label only), so treat decode errors as empty request
		req = models.CreateInviteRequest{}
	}

	ctx := r.Context()

	var exists bool
	if err := h.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM events WHERE id = $1)`, eventID).Scan(&exists); err != nil || !exists {
		h.writeError(w, http.StatusNotFound, "not_found", "Event not found")
		return
	}

	// generate unique invite code with retry
	var invite models.AdminInvite
	for attempt := 0; attempt < maxCodeRetries; attempt++ {
		code, err := generateInviteCode()
		if err != nil {
			slog.Error("Failed to generate invite code", "error", err)
			h.writeError(w, http.StatusInternalServerError, "code_error", "Failed to generate invite code")
			return
		}

		err = h.db.QueryRow(ctx, `
			INSERT INTO invites (event_id, invite_code, label)
			VALUES ($1, $2, $3)
			RETURNING id, invite_code, label, created_at
		`, eventID, code, req.Label).Scan(
			&invite.ID, &invite.InviteCode, &invite.Label, &invite.CreatedAt,
		)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") ||
				strings.Contains(err.Error(), "unique constraint") {
				continue
			}
			slog.Error("DB error creating invite", "error", err)
			h.writeError(w, http.StatusInternalServerError, "db_error", "Failed to create invite")
			return
		}

		invite.Guests = []models.AdminGuest{}
		h.writeJSON(w, http.StatusCreated, invite)
		return
	}

	h.writeError(w, http.StatusInternalServerError, "code_error",
		fmt.Sprintf("Failed to generate unique code after %d attempts", maxCodeRetries))
}

// DeleteInvite handles DELETE /admin/events/:id/invites/:inviteId - removes an invite and its guests.
func (h *Handlers) DeleteInvite(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	inviteID := chi.URLParam(r, "inviteId")
	ctx := r.Context()

	result, err := h.db.Exec(ctx,
		`DELETE FROM invites WHERE id = $1 AND event_id = $2`,
		inviteID, eventID,
	)
	if err != nil {
		slog.Error("DB error deleting invite", "error", err)
		h.writeError(w, http.StatusInternalServerError, "db_error", "Failed to delete invite")
		return
	}

	if result.RowsAffected() == 0 {
		h.writeError(w, http.StatusNotFound, "not_found", "Invite not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddGuest handles POST /admin/events/:id/invites/:inviteId/guests - adds a guest to an invite.
func (h *Handlers) AddGuest(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	inviteID := chi.URLParam(r, "inviteId")

	var req models.AddGuestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	req.DisplayName = strings.TrimSpace(req.DisplayName)
	if req.DisplayName == "" {
		h.writeError(w, http.StatusBadRequest, "invalid_name", "Display name is required")
		return
	}

	ctx := r.Context()

	// verify invite belongs to event
	var inviteExists bool
	err := h.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM invites WHERE id = $1 AND event_id = $2)`,
		inviteID, eventID,
	).Scan(&inviteExists)
	if err != nil || !inviteExists {
		h.writeError(w, http.StatusNotFound, "not_found", "Invite not found")
		return
	}

	var guest models.AdminGuest
	err = h.db.QueryRow(ctx, `
		INSERT INTO guests (event_id, invite_id, display_name)
		VALUES ($1, $2, $3)
		RETURNING id, display_name, rsvp_status, rsvp_message, rsvp_at, created_at
	`, eventID, inviteID, req.DisplayName).Scan(
		&guest.ID, &guest.DisplayName, &guest.RsvpStatus,
		&guest.RsvpMessage, &guest.RsvpAt, &guest.CreatedAt,
	)
	if err != nil {
		slog.Error("DB error adding guest", "error", err)
		h.writeError(w, http.StatusInternalServerError, "db_error", "Failed to add guest")
		return
	}

	h.writeJSON(w, http.StatusCreated, guest)
}

// UpdateGuest handles PUT /admin/events/:id/guests/:guestId - updates a guest's name or RSVP.
func (h *Handlers) UpdateGuest(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	guestID := chi.URLParam(r, "guestId")

	var req models.UpdateGuestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	req.DisplayName = strings.TrimSpace(req.DisplayName)
	if req.DisplayName == "" {
		h.writeError(w, http.StatusBadRequest, "invalid_name", "Display name is required")
		return
	}

	req.RsvpStatus = strings.ToLower(strings.TrimSpace(req.RsvpStatus))
	if req.RsvpStatus != "pending" && req.RsvpStatus != "yes" && req.RsvpStatus != "no" {
		h.writeError(w, http.StatusBadRequest, "invalid_status", "RSVP status must be pending, yes, or no")
		return
	}

	ctx := r.Context()

	var guest models.AdminGuest
	err := h.db.QueryRow(ctx, `
		UPDATE guests SET display_name = $1, rsvp_status = $2
		WHERE id = $3 AND event_id = $4
		RETURNING id, display_name, rsvp_status, rsvp_message, rsvp_at, created_at
	`, req.DisplayName, req.RsvpStatus, guestID, eventID).Scan(
		&guest.ID, &guest.DisplayName, &guest.RsvpStatus,
		&guest.RsvpMessage, &guest.RsvpAt, &guest.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "not_found", "Guest not found")
			return
		}
		slog.Error("DB error updating guest", "error", err)
		h.writeError(w, http.StatusInternalServerError, "db_error", "Failed to update guest")
		return
	}

	h.writeJSON(w, http.StatusOK, guest)
}

// DeleteGuest handles DELETE /admin/events/:id/guests/:guestId - removes a guest.
func (h *Handlers) DeleteGuest(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	guestID := chi.URLParam(r, "guestId")
	ctx := r.Context()

	result, err := h.db.Exec(ctx,
		`DELETE FROM guests WHERE id = $1 AND event_id = $2`,
		guestID, eventID,
	)
	if err != nil {
		slog.Error("DB error deleting guest", "error", err)
		h.writeError(w, http.StatusInternalServerError, "db_error", "Failed to delete guest")
		return
	}

	if result.RowsAffected() == 0 {
		h.writeError(w, http.StatusNotFound, "not_found", "Guest not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// fetchInviteGuests retrieves guests for an invite.
func (h *Handlers) fetchInviteGuests(ctx context.Context, inviteID string) ([]models.AdminGuest, error) {
	rows, err := h.db.Query(ctx, `
		SELECT id, display_name, rsvp_status, rsvp_message, rsvp_at, created_at
		FROM guests WHERE invite_id = $1
		ORDER BY created_at ASC
	`, inviteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	guests := []models.AdminGuest{}
	for rows.Next() {
		var g models.AdminGuest
		if err := rows.Scan(&g.ID, &g.DisplayName, &g.RsvpStatus, &g.RsvpMessage, &g.RsvpAt, &g.CreatedAt); err != nil {
			return nil, err
		}
		guests = append(guests, g)
	}
	return guests, nil
}
