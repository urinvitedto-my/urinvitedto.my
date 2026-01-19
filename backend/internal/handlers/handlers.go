// Package handlers implements HTTP handlers for the API.
package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/urinvitedto-my/backend/internal/models"
)

// Handlers bundles shared deps for HTTP handlers.
type Handlers struct {
	db *pgxpool.Pool
}

// New creates a new Handlers instance.
func New(db *pgxpool.Pool) *Handlers {
	return &Handlers{db: db}
}

// writeJSON writes a JSON response.
func (h *Handlers) writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("Error encoding JSON", "error", err)
	}
}

// writeError writes an error response.
func (h *Handlers) writeError(
	w http.ResponseWriter,
	status int,
	err string,
	msg string,
) {
	h.writeJSON(w, status, models.ErrorResponse{Error: err, Message: msg})
}

// GetEventSummary handles GET /events/:type/:slug/summary
func (h *Handlers) GetEventSummary(w http.ResponseWriter, r *http.Request) {
	eventType := chi.URLParam(r, "type")
	slug := chi.URLParam(r, "slug")

	var event models.EventSummary
	err := h.db.QueryRow(r.Context(), `
		SELECT id, type, slug, title, is_public, cover_image_url, starts_at, location
		FROM events WHERE type = $1 AND slug = $2
	`, eventType, slug).Scan(
		&event.ID, &event.Type, &event.Slug, &event.Title,
		&event.IsPublic, &event.CoverImageURL, &event.StartsAt, &event.Location,
	)

	if err == pgx.ErrNoRows {
		h.writeError(w, http.StatusNotFound, "not_found", "Event not found")
		return
	}
	if err != nil {
		slog.Error("DB error fetching event summary", "error", err)
		h.writeError(
			w,
			http.StatusInternalServerError,
			"db_error",
			"Failed to fetch event",
		)
		return
	}

	h.writeJSON(w, http.StatusOK, event)
}

// GetEventDetails handles GET /events/:type/:slug/details
func (h *Handlers) GetEventDetails(w http.ResponseWriter, r *http.Request) {
	eventType := chi.URLParam(r, "type")
	slug := chi.URLParam(r, "slug")
	inviteCode := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("invite")))

	ctx := r.Context()

	// fetch event
	var event models.Event
	err := h.db.QueryRow(ctx, `
		SELECT id, type, slug, title, description, is_public, cover_image_url,
		       location_photo_url, starts_at, location, custom_content, enabled_components, created_at
		FROM events WHERE type = $1 AND slug = $2
	`, eventType, slug).Scan(
		&event.ID, &event.Type, &event.Slug, &event.Title, &event.Description,
		&event.IsPublic, &event.CoverImageURL, &event.LocationPhotoURL,
		&event.StartsAt, &event.Location, &event.CustomContent, &event.EnabledComponents, &event.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		h.writeError(w, http.StatusNotFound, "not_found", "Event not found")
		return
	}
	if err != nil {
		slog.Error("DB error fetching event", "error", err)
		h.writeError(
			w,
			http.StatusInternalServerError,
			"db_error",
			"Failed to fetch event",
		)
		return
	}

	// private event requires valid invite code
	if !event.IsPublic && inviteCode == "" {
		h.writeError(
			w,
			http.StatusUnauthorized,
			"invite_required",
			"Invite code required for private event",
		)
		return
	}

	resp := models.EventDetailsResponse{Event: event}

	// fetch hosts
	resp.Hosts, err = h.fetchHosts(ctx, event.ID)
	if err != nil {
		slog.Error("DB error fetching hosts", "error", err)
	}

	// fetch schedule
	resp.Schedule, err = h.fetchSchedule(ctx, event.ID)
	if err != nil {
		slog.Error("DB error fetching schedule", "error", err)
	}

	// fetch FAQs
	resp.FAQs, err = h.fetchFAQs(ctx, event.ID)
	if err != nil {
		slog.Error("DB error fetching FAQs", "error", err)
	}

	// fetch gallery
	resp.Gallery, err = h.fetchGallery(ctx, event.ID)
	if err != nil {
		slog.Error("DB error fetching gallery", "error", err)
	}

	// fetch gifts
	resp.Gifts, err = h.fetchGifts(ctx, event.ID)
	if err != nil {
		slog.Error("DB error fetching gifts", "error", err)
	}

	// fetch confirmed guests count
	resp.ConfirmedGuestsCount, err = h.countConfirmedGuests(ctx, event.ID)
	if err != nil {
		slog.Error("DB error counting confirmed guests", "error", err)
	}

	// for private events, validate invite and fetch guests on this invite
	if !event.IsPublic {
		invite, guests, invErr := h.fetchInviteWithGuests(ctx, event.ID, inviteCode)
		if invErr == pgx.ErrNoRows {
			h.writeError(
				w,
				http.StatusUnauthorized,
				"invalid_invite",
				"Invalid invite code",
			)
			return
		}
		if invErr != nil {
			slog.Error("DB error fetching invite", "error", invErr)
			h.writeError(
				w,
				http.StatusInternalServerError,
				"db_error",
				"Failed to validate invite",
			)
			return
		}
		invite.Guests = guests
		resp.Invite = invite
	}

	h.writeJSON(w, http.StatusOK, resp)
}

// GetConfirmedGuests handles GET /events/:type/:slug/confirmed-guests
func (h *Handlers) GetConfirmedGuests(w http.ResponseWriter, r *http.Request) {
	eventType := chi.URLParam(r, "type")
	slug := chi.URLParam(r, "slug")

	ctx := r.Context()

	// fetch event ID
	var eventID string
	err := h.db.QueryRow(ctx, `SELECT id FROM events WHERE type = $1 AND slug = $2`, eventType, slug).
		Scan(&eventID)
	if err == pgx.ErrNoRows {
		h.writeError(w, http.StatusNotFound, "not_found", "Event not found")
		return
	}
	if err != nil {
		slog.Error("DB error fetching event ID", "error", err)
		h.writeError(
			w,
			http.StatusInternalServerError,
			"db_error",
			"Failed to fetch event",
		)
		return
	}

	// fetch confirmed guests
	rows, err := h.db.Query(ctx, `
		SELECT display_name FROM guests
		WHERE event_id = $1 AND rsvp_status = 'yes'
		ORDER BY display_name ASC
	`, eventID)
	if err != nil {
		slog.Error("DB error fetching confirmed guests", "error", err)
		h.writeError(
			w,
			http.StatusInternalServerError,
			"db_error",
			"Failed to fetch guests",
		)
		return
	}
	defer rows.Close()

	guests := []models.ConfirmedGuest{}
	for rows.Next() {
		var g models.ConfirmedGuest
		if err := rows.Scan(&g.DisplayName); err != nil {
			slog.Error("Error scanning guest", "error", err)
			continue
		}
		guests = append(guests, g)
	}

	h.writeJSON(w, http.StatusOK, models.ConfirmedGuestsResponse{
		Guests: guests,
		Count:  len(guests),
	})
}

// PostRSVP handles POST /events/:type/:slug/rsvp
func (h *Handlers) PostRSVP(w http.ResponseWriter, r *http.Request) {
	eventType := chi.URLParam(r, "type")
	slug := chi.URLParam(r, "slug")

	var req models.RSVPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid_body", "Invalid request body")
		return
	}

	req.InviteCode = strings.ToUpper(strings.TrimSpace(req.InviteCode))

	// validate status
	if req.Status != "yes" && req.Status != "no" {
		h.writeError(
			w,
			http.StatusBadRequest,
			"invalid_status",
			"Status must be 'yes' or 'no'",
		)
		return
	}

	ctx := r.Context()

	// fetch event ID
	var eventID string
	err := h.db.QueryRow(ctx, `SELECT id FROM events WHERE type = $1 AND slug = $2`, eventType, slug).
		Scan(&eventID)
	if err == pgx.ErrNoRows {
		h.writeError(w, http.StatusNotFound, "not_found", "Event not found")
		return
	}
	if err != nil {
		slog.Error("DB error fetching event ID", "error", err)
		h.writeError(
			w,
			http.StatusInternalServerError,
			"db_error",
			"Failed to fetch event",
		)
		return
	}

	// validate invite exists
	var inviteID string
	err = h.db.QueryRow(ctx, `
		SELECT id FROM invites WHERE event_id = $1 AND invite_code = $2
	`, eventID, req.InviteCode).Scan(&inviteID)
	if err == pgx.ErrNoRows {
		h.writeError(
			w,
			http.StatusUnauthorized,
			"invalid_invite",
			"Invalid invite code",
		)
		return
	}
	if err != nil {
		slog.Error("DB error validating invite", "error", err)
		h.writeError(
			w,
			http.StatusInternalServerError,
			"db_error",
			"Failed to validate invite",
		)
		return
	}

	// validate guest belongs to this invite
	var guestEventID, guestInviteID string
	err = h.db.QueryRow(ctx, `
		SELECT event_id, invite_id FROM guests WHERE id = $1
	`, req.GuestID).Scan(&guestEventID, &guestInviteID)
	if err == pgx.ErrNoRows {
		h.writeError(w, http.StatusNotFound, "guest_not_found", "Guest not found")
		return
	}
	if err != nil {
		slog.Error("DB error fetching guest", "error", err)
		h.writeError(
			w,
			http.StatusInternalServerError,
			"db_error",
			"Failed to fetch guest",
		)
		return
	}

	if guestEventID != eventID || guestInviteID != inviteID {
		h.writeError(
			w,
			http.StatusForbidden,
			"guest_mismatch",
			"Guest does not belong to this invite",
		)
		return
	}

	// update RSVP
	now := time.Now()
	var resp models.RSVPResponse
	err = h.db.QueryRow(ctx, `
		UPDATE guests
		SET rsvp_status = $1, rsvp_message = $2, rsvp_at = $3
		WHERE id = $4
		RETURNING id, display_name, rsvp_status, rsvp_message, rsvp_at
	`, req.Status, req.Message, now, req.GuestID).Scan(
		&resp.ID, &resp.DisplayName, &resp.RsvpStatus, &resp.RsvpMessage, &resp.RsvpAt,
	)
	if err != nil {
		slog.Error("DB error updating RSVP", "error", err)
		h.writeError(
			w,
			http.StatusInternalServerError,
			"db_error",
			"Failed to update RSVP",
		)
		return
	}

	h.writeJSON(w, http.StatusOK, resp)
}

// fetchHosts retrieves hosts for an event.
func (h *Handlers) fetchHosts(
	ctx context.Context,
	eventID string,
) ([]models.Host, error) {
	rows, err := h.db.Query(ctx, `
		SELECT id, display_name FROM hosts WHERE event_id = $1
	`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hosts := []models.Host{}
	for rows.Next() {
		var host models.Host
		if err := rows.Scan(&host.ID, &host.DisplayName); err != nil {
			return nil, err
		}
		hosts = append(hosts, host)
	}
	return hosts, nil
}

// fetchSchedule retrieves schedule items for an event.
func (h *Handlers) fetchSchedule(
	ctx context.Context,
	eventID string,
) ([]models.ScheduleItem, error) {
	rows, err := h.db.Query(ctx, `
		SELECT id, time, title, description, order_index
		FROM event_schedule WHERE event_id = $1 ORDER BY order_index ASC
	`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []models.ScheduleItem{}
	for rows.Next() {
		var item models.ScheduleItem
		if err := rows.Scan(&item.ID, &item.Time, &item.Title, &item.Description, &item.OrderIndex); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// fetchFAQs retrieves FAQs for an event.
func (h *Handlers) fetchFAQs(
	ctx context.Context,
	eventID string,
) ([]models.FAQ, error) {
	rows, err := h.db.Query(ctx, `
		SELECT id, question, answer, order_index
		FROM event_faqs WHERE event_id = $1 ORDER BY order_index ASC
	`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	faqs := []models.FAQ{}
	for rows.Next() {
		var faq models.FAQ
		if err := rows.Scan(&faq.ID, &faq.Question, &faq.Answer, &faq.OrderIndex); err != nil {
			return nil, err
		}
		faqs = append(faqs, faq)
	}
	return faqs, nil
}

// fetchGallery retrieves gallery items for an event.
func (h *Handlers) fetchGallery(
	ctx context.Context,
	eventID string,
) ([]models.GalleryItem, error) {
	rows, err := h.db.Query(ctx, `
		SELECT id, media_type, media_url, caption, order_index
		FROM event_gallery WHERE event_id = $1 ORDER BY order_index ASC
	`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []models.GalleryItem{}
	for rows.Next() {
		var item models.GalleryItem
		if err := rows.Scan(&item.ID, &item.MediaType, &item.MediaURL, &item.Caption, &item.OrderIndex); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// fetchGifts retrieves gift items for an event.
func (h *Handlers) fetchGifts(
	ctx context.Context,
	eventID string,
) ([]models.Gift, error) {
	rows, err := h.db.Query(ctx, `
		SELECT id, gift_type, title, description, link, order_index
		FROM event_gifts WHERE event_id = $1 ORDER BY order_index ASC
	`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gifts := []models.Gift{}
	for rows.Next() {
		var gift models.Gift
		if err := rows.Scan(&gift.ID, &gift.GiftType, &gift.Title, &gift.Description, &gift.Link, &gift.OrderIndex); err != nil {
			return nil, err
		}
		gifts = append(gifts, gift)
	}
	return gifts, nil
}

// countConfirmedGuests counts guests with rsvp_status = 'yes'.
func (h *Handlers) countConfirmedGuests(
	ctx context.Context,
	eventID string,
) (int, error) {
	var count int
	err := h.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM guests WHERE event_id = $1 AND rsvp_status = 'yes'
	`, eventID).Scan(&count)
	return count, err
}

// fetchInviteWithGuests retrieves an invite and its guests.
func (h *Handlers) fetchInviteWithGuests(
	ctx context.Context,
	eventID, inviteCode string,
) (*models.Invite, []models.Guest, error) {
	var invite models.Invite
	err := h.db.QueryRow(ctx, `
		SELECT id, label FROM invites WHERE event_id = $1 AND invite_code = $2
	`, eventID, inviteCode).Scan(&invite.ID, &invite.Label)
	if err != nil {
		return nil, nil, err
	}

	rows, err := h.db.Query(ctx, `
		SELECT id, display_name, rsvp_status, rsvp_message, rsvp_at
		FROM guests WHERE invite_id = $1 ORDER BY display_name ASC
	`, invite.ID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	guests := []models.Guest{}
	for rows.Next() {
		var g models.Guest
		if err := rows.Scan(&g.ID, &g.DisplayName, &g.RsvpStatus, &g.RsvpMessage, &g.RsvpAt); err != nil {
			return nil, nil, err
		}
		guests = append(guests, g)
	}

	return &invite, guests, nil
}
