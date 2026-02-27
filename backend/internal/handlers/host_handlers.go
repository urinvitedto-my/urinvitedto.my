package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/urinvitedto-my/backend/internal/models"
)

// GetHostEvents handles GET /host/events - returns events for the authenticated host.
func (h *Handlers) GetHostEvents(w http.ResponseWriter, r *http.Request) {
	// get email from context (set by auth middleware)
	email, ok := r.Context().Value(UserEmailKey).(string)
	if !ok || email == "" {
		h.writeError(w, http.StatusUnauthorized, "unauthorized", "Not authenticated")
		return
	}

	ctx := r.Context()
	emailLower := strings.ToLower(strings.TrimSpace(email))
	slog.Info("Fetching host events", "email", emailLower)

	// fetch events where user is a host (by auth_user_id OR contact_email)
	rows, err := h.db.Query(ctx, `
		SELECT DISTINCT e.id, e.type, e.slug, e.title, e.is_public, e.starts_at, e.location, e.created_at
		FROM events e
		INNER JOIN hosts h ON h.event_id = e.id
		LEFT JOIN auth.users u ON u.id = h.auth_user_id
		WHERE LOWER(u.email) = $1 OR LOWER(h.contact_email) = $1
		ORDER BY e.starts_at DESC NULLS LAST
	`, emailLower)
	if err != nil {
		slog.Error("DB error fetching host events", "error", err)
		h.writeError(
			w,
			http.StatusInternalServerError,
			"db_error",
			"Failed to fetch events",
		)
		return
	}
	defer rows.Close()

	events := []models.HostEvent{}
	for rows.Next() {
		var e models.HostEvent
		if err := rows.Scan(
			&e.ID,
			&e.Type,
			&e.Slug,
			&e.Title,
			&e.IsPublic,
			&e.StartsAt,
			&e.Location,
			&e.CreatedAt,
		); err != nil {
			slog.Error("Error scanning event", "error", err)
			continue
		}
		events = append(events, e)
	}

	slog.Info("Host events found", "email", emailLower, "count", len(events))
	h.writeJSON(w, http.StatusOK, models.HostEventsResponse{Events: events})
}

// GetHostGuests handles GET /host/events/:eventId/guests - returns guests for an event if the user is a host.
func (h *Handlers) GetHostGuests(w http.ResponseWriter, r *http.Request) {
	email, ok := r.Context().Value(UserEmailKey).(string)
	if !ok || email == "" {
		h.writeError(w, http.StatusUnauthorized, "unauthorized", "Not authenticated")
		return
	}

	eventID := chi.URLParam(r, "eventId")
	if eventID == "" {
		h.writeError(w, http.StatusBadRequest, "bad_request", "Event ID required")
		return
	}

	ctx := r.Context()
	emailLower := strings.ToLower(strings.TrimSpace(email))

	// verify user is a host of this event (auth_user_id OR contact_email)
	var isHost bool
	err := h.db.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM hosts h
			LEFT JOIN auth.users u ON u.id = h.auth_user_id
			WHERE h.event_id = $1 AND (LOWER(u.email) = $2 OR LOWER(h.contact_email) = $2)
		)
	`, eventID, emailLower).Scan(&isHost)
	if err != nil || !isHost {
		h.writeError(w, http.StatusForbidden, "forbidden", "Not a host of this event")
		return
	}

	rows, err := h.db.Query(ctx, `
		SELECT id, display_name, rsvp_status, rsvp_message, rsvp_at, created_at
		FROM guests WHERE event_id = $1 ORDER BY display_name ASC
	`, eventID)
	if err != nil {
		slog.Error("DB error fetching host guests", "error", err)
		h.writeError(w, http.StatusInternalServerError, "db_error", "Failed to fetch guests")
		return
	}
	defer rows.Close()

	guests := []models.AdminGuest{}
	for rows.Next() {
		var g models.AdminGuest
		if err := rows.Scan(&g.ID, &g.DisplayName, &g.RsvpStatus, &g.RsvpMessage, &g.RsvpAt, &g.CreatedAt); err != nil {
			slog.Error("Error scanning guest", "error", err)
			continue
		}
		guests = append(guests, g)
	}

	h.writeJSON(w, http.StatusOK, models.HostGuestsResponse{Guests: guests})
}
