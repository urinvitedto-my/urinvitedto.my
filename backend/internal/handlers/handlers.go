// Package handlers implements HTTP handlers for the API.
package handlers

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/urinvitedto-my/backend/internal/models"
	"github.com/urinvitedto-my/backend/internal/utils"
)

// ContextKey is a typed key for context values.
type ContextKey string

// UserEmailKey is the context key for storing authenticated user email.
const UserEmailKey ContextKey = "userEmail"

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
	utils.WriteJSON(w, status, v)
}

// writeError writes an error response.
func (h *Handlers) writeError(
	w http.ResponseWriter,
	status int,
	err string,
	msg string,
) {
	utils.WriteError(w, status, err, msg)
}

// IsAdmin checks if the given email belongs to an admin.
func (h *Handlers) IsAdmin(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := h.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM admins WHERE email = $1)`, email).
		Scan(&exists)
	return exists, err
}

// getEventID fetches event ID by type and slug.
func (h *Handlers) getEventID(
	ctx context.Context,
	eventType, slug string,
) (string, error) {
	var eventID string
	err := h.db.QueryRow(ctx, `
		SELECT id FROM events WHERE type = $1 AND slug = $2
	`, eventType, slug).Scan(&eventID)
	return eventID, err
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
		if err := rows.Scan(
			&item.ID,
			&item.Time,
			&item.Title,
			&item.Description,
			&item.OrderIndex,
		); err != nil {
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
		if err := rows.Scan(
			&faq.ID,
			&faq.Question,
			&faq.Answer,
			&faq.OrderIndex,
		); err != nil {
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
		if err := rows.Scan(
			&item.ID,
			&item.MediaType,
			&item.MediaURL,
			&item.Caption,
			&item.OrderIndex,
		); err != nil {
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
		if err := rows.Scan(
			&gift.ID,
			&gift.GiftType,
			&gift.Title,
			&gift.Description,
			&gift.Link,
			&gift.OrderIndex,
		); err != nil {
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
		if err := rows.Scan(
			&g.ID,
			&g.DisplayName,
			&g.RsvpStatus,
			&g.RsvpMessage,
			&g.RsvpAt,
		); err != nil {
			return nil, nil, err
		}
		guests = append(guests, g)
	}

	return &invite, guests, nil
}
