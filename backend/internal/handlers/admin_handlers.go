package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/urinvitedto-my/backend/internal/models"
)

var slugRegex = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// ListEvents handles GET /admin/events - returns all events with hosts.
func (h *Handlers) ListEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rows, err := h.db.Query(ctx, `
		SELECT id, type, slug, title, description, is_public,
			cover_image_url, location_photo_url, music_url, starts_at, location, created_at
		FROM events ORDER BY created_at DESC
	`)
	if err != nil {
		slog.Error("DB error listing events", "error", err)
		h.writeError(
			w,
			http.StatusInternalServerError,
			"db_error",
			"Failed to list events",
		)
		return
	}
	defer rows.Close()

	events := []models.AdminEvent{}
	for rows.Next() {
		var e models.AdminEvent
		if err := rows.Scan(
			&e.ID,
			&e.Type,
			&e.Slug,
			&e.Title,
			&e.Description,
			&e.IsPublic,
			&e.CoverImageURL,
			&e.LocationPhotoURL,
			&e.MusicURL,
			&e.StartsAt,
			&e.Location,
			&e.CreatedAt,
		); err != nil {
			slog.Error("Error scanning event", "error", err)
			continue
		}
		events = append(events, e)
	}

	// fetch hosts for each event
	for i := range events {
		hosts, err := h.fetchAdminHosts(ctx, events[i].ID)
		if err != nil {
			slog.Error("Error fetching hosts", "eventId", events[i].ID, "error", err)
		}
		events[i].Hosts = hosts
	}

	h.writeJSON(w, http.StatusOK, models.AdminEventsResponse{Events: events})
}

// CreateEvent handles POST /admin/events - creates a new event.
func (h *Handlers) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req models.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	// validate type
	req.Type = strings.ToLower(strings.TrimSpace(req.Type))
	if req.Type != "wedding" && req.Type != "birthday" && req.Type != "party" {
		h.writeError(
			w,
			http.StatusBadRequest,
			"invalid_type",
			"Type must be wedding, birthday, or party",
		)
		return
	}

	// validate slug
	req.Slug = strings.ToLower(strings.TrimSpace(req.Slug))
	if req.Slug == "" {
		h.writeError(w, http.StatusBadRequest, "invalid_slug", "Slug is required")
		return
	}
	if !slugRegex.MatchString(req.Slug) {
		h.writeError(
			w,
			http.StatusBadRequest,
			"invalid_slug",
			"Slug must be lowercase alphanumeric with hyphens only",
		)
		return
	}

	// validate title
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		h.writeError(w, http.StatusBadRequest, "invalid_title", "Title is required")
		return
	}

	ctx := r.Context()

	// parse startsAt if provided
	var startsAt *time.Time
	if req.StartsAt != nil && *req.StartsAt != "" {
		t, err := time.Parse(time.RFC3339, *req.StartsAt)
		if err != nil {
			h.writeError(
				w,
				http.StatusBadRequest,
				"invalid_date",
				"Invalid date format. Use ISO 8601 (e.g., 2024-06-15T14:00:00Z)",
			)
			return
		}
		startsAt = &t
	}

	var event models.AdminEvent
	err := h.db.QueryRow(ctx, `
		INSERT INTO events (type, slug, title, is_public, starts_at, location)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, type, slug, title, is_public, starts_at, location, created_at
	`, req.Type, req.Slug, req.Title, req.IsPublic, startsAt, req.Location).Scan(
		&event.ID, &event.Type, &event.Slug, &event.Title, &event.IsPublic,
		&event.StartsAt, &event.Location, &event.CreatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "unique constraint") {
			h.writeError(
				w,
				http.StatusConflict,
				"duplicate_slug",
				"An event with this type and slug already exists",
			)
			return
		}
		slog.Error("DB error creating event", "error", err)
		h.writeError(
			w,
			http.StatusInternalServerError,
			"db_error",
			"Failed to create event",
		)
		return
	}

	event.Hosts = []models.AdminHost{}
	h.writeJSON(w, http.StatusCreated, event)
}

// AddHost handles POST /admin/events/:id/hosts - adds a host to an event.
func (h *Handlers) AddHost(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")

	var req models.AddHostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.DisplayName = strings.TrimSpace(req.DisplayName)

	if req.Email == "" {
		h.writeError(w, http.StatusBadRequest, "invalid_email", "Email is required")
		return
	}
	if req.DisplayName == "" {
		h.writeError(
			w,
			http.StatusBadRequest,
			"invalid_name",
			"Display name is required",
		)
		return
	}

	ctx := r.Context()

	// verify event exists
	var exists bool
	err := h.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM events WHERE id = $1)`, eventID).
		Scan(&exists)
	if err != nil || !exists {
		h.writeError(w, http.StatusNotFound, "not_found", "Event not found")
		return
	}

	// look up auth user by email from Supabase auth.users
	var authUserID *string
	err = h.db.QueryRow(ctx, `SELECT id::text FROM auth.users WHERE email = $1`, req.Email).
		Scan(&authUserID)
	if err == pgx.ErrNoRows {
		// user doesn't exist in auth yet - that's okay, we'll store null
		authUserID = nil
	} else if err != nil {
		slog.Error("DB error looking up auth user", "error", err)
		h.writeError(
			w,
			http.StatusInternalServerError,
			"db_error",
			"Failed to look up user",
		)
		return
	}

	// insert host
	var host models.AdminHost
	err = h.db.QueryRow(ctx, `
		INSERT INTO hosts (event_id, display_name, contact_email, auth_user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, display_name, contact_email, auth_user_id
	`, eventID, req.DisplayName, req.Email, authUserID).Scan(
		&host.ID, &host.DisplayName, &host.ContactEmail, &host.AuthUserID,
	)
	if err != nil {
		slog.Error("DB error adding host", "error", err)
		h.writeError(
			w,
			http.StatusInternalServerError,
			"db_error",
			"Failed to add host",
		)
		return
	}

	h.writeJSON(w, http.StatusCreated, host)
}

// DeleteHost handles DELETE /admin/events/:id/hosts/:hostId - removes a host.
func (h *Handlers) DeleteHost(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	hostID := chi.URLParam(r, "hostId")

	ctx := r.Context()

	result, err := h.db.Exec(
		ctx,
		`DELETE FROM hosts WHERE id = $1 AND event_id = $2`,
		hostID,
		eventID,
	)
	if err != nil {
		slog.Error("DB error deleting host", "error", err)
		h.writeError(
			w,
			http.StatusInternalServerError,
			"db_error",
			"Failed to delete host",
		)
		return
	}

	if result.RowsAffected() == 0 {
		h.writeError(w, http.StatusNotFound, "not_found", "Host not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateEvent handles PUT /admin/events/:id - updates an existing event.
func (h *Handlers) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")

	var req models.UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	req.Type = strings.ToLower(strings.TrimSpace(req.Type))
	if req.Type != "wedding" && req.Type != "birthday" && req.Type != "party" {
		h.writeError(
			w,
			http.StatusBadRequest,
			"invalid_type",
			"Type must be wedding, birthday, or party",
		)
		return
	}

	req.Slug = strings.ToLower(strings.TrimSpace(req.Slug))
	if req.Slug == "" {
		h.writeError(w, http.StatusBadRequest, "invalid_slug", "Slug is required")
		return
	}
	if !slugRegex.MatchString(req.Slug) {
		h.writeError(
			w,
			http.StatusBadRequest,
			"invalid_slug",
			"Slug must be lowercase alphanumeric with hyphens only",
		)
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		h.writeError(w, http.StatusBadRequest, "invalid_title", "Title is required")
		return
	}

	ctx := r.Context()

	var startsAt *time.Time
	if req.StartsAt != nil && *req.StartsAt != "" {
		t, err := time.Parse(time.RFC3339, *req.StartsAt)
		if err != nil {
			h.writeError(
				w,
				http.StatusBadRequest,
				"invalid_date",
				"Invalid date format. Use ISO 8601 (e.g., 2024-06-15T14:00:00Z)",
			)
			return
		}
		startsAt = &t
	}

	var event models.AdminEvent
	err := h.db.QueryRow(ctx, `
		UPDATE events SET
			type = $1, slug = $2, title = $3, description = $4, is_public = $5,
			starts_at = $6, location = $7, cover_image_url = $8, location_photo_url = $9,
			music_url = $10
		WHERE id = $11
		RETURNING id, type, slug, title, description, is_public,
			cover_image_url, location_photo_url, music_url, starts_at, location, created_at
	`, req.Type, req.Slug, req.Title, req.Description, req.IsPublic,
		startsAt, req.Location, req.CoverImageURL, req.LocationPhotoURL, req.MusicURL, eventID,
	).Scan(
		&event.ID, &event.Type, &event.Slug, &event.Title, &event.Description,
		&event.IsPublic, &event.CoverImageURL, &event.LocationPhotoURL, &event.MusicURL,
		&event.StartsAt, &event.Location, &event.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "not_found", "Event not found")
			return
		}
		if strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "unique constraint") {
			h.writeError(
				w,
				http.StatusConflict,
				"duplicate_slug",
				"An event with this type and slug already exists",
			)
			return
		}
		slog.Error("DB error updating event", "error", err)
		h.writeError(
			w,
			http.StatusInternalServerError,
			"db_error",
			"Failed to update event",
		)
		return
	}

	hosts, err := h.fetchAdminHosts(ctx, event.ID)
	if err != nil {
		slog.Error("Error fetching hosts", "eventId", event.ID, "error", err)
	}
	event.Hosts = hosts

	h.writeJSON(w, http.StatusOK, event)
}

// DeleteEvent handles DELETE /admin/events/:id - deletes an event and all related data.
func (h *Handlers) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	ctx := r.Context()

	result, err := h.db.Exec(ctx, `DELETE FROM events WHERE id = $1`, eventID)
	if err != nil {
		slog.Error("DB error deleting event", "error", err)
		h.writeError(
			w,
			http.StatusInternalServerError,
			"db_error",
			"Failed to delete event",
		)
		return
	}

	if result.RowsAffected() == 0 {
		h.writeError(w, http.StatusNotFound, "not_found", "Event not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// fetchAdminHosts retrieves hosts with email for admin views.
func (h *Handlers) fetchAdminHosts(
	ctx context.Context,
	eventID string,
) ([]models.AdminHost, error) {
	rows, err := h.db.Query(ctx, `
		SELECT id, display_name, contact_email, auth_user_id
		FROM hosts WHERE event_id = $1
	`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hosts := []models.AdminHost{}
	for rows.Next() {
		var host models.AdminHost
		if err := rows.Scan(
			&host.ID,
			&host.DisplayName,
			&host.ContactEmail,
			&host.AuthUserID,
		); err != nil {
			return nil, err
		}
		hosts = append(hosts, host)
	}
	return hosts, nil
}
