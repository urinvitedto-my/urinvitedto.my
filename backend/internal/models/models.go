// Package models defines data types for the API.
package models

import (
	"encoding/json"
	"time"
)

// Event represents a single event.
type Event struct {
	ID                string          `json:"id"`
	Type              string          `json:"type"`
	Slug              string          `json:"slug"`
	Title             string          `json:"title"`
	Description       *string         `json:"description,omitempty"`
	IsPublic          bool            `json:"isPublic"`
	CoverImageURL     *string         `json:"coverImageUrl,omitempty"`
	LocationPhotoURL  *string         `json:"locationPhotoUrl,omitempty"`
	MusicURL          *string         `json:"musicUrl,omitempty"`
	StartsAt          *time.Time      `json:"startsAt,omitempty"`
	Location          *string         `json:"location,omitempty"`
	CustomContent     json.RawMessage `json:"customContent,omitempty"`
	EnabledComponents json.RawMessage `json:"enabledComponents,omitempty"`
	CreatedAt         time.Time       `json:"createdAt"`
}

// EventSummary is a lightweight event response for the summary endpoint.
type EventSummary struct {
	ID            string     `json:"id"`
	Type          string     `json:"type"`
	Slug          string     `json:"slug"`
	Title         string     `json:"title"`
	IsPublic      bool       `json:"isPublic"`
	CoverImageURL *string    `json:"coverImageUrl,omitempty"`
	MusicURL      *string    `json:"musicUrl,omitempty"`
	StartsAt      *time.Time `json:"startsAt,omitempty"`
	Location      *string    `json:"location,omitempty"`
}

// Host represents a celebrant/host of an event.
type Host struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

// Guest represents an individual invitee.
type Guest struct {
	ID          string     `json:"id"`
	DisplayName string     `json:"displayName"`
	RsvpStatus  string     `json:"rsvpStatus"`
	RsvpMessage *string    `json:"rsvpMessage,omitempty"`
	RsvpAt      *time.Time `json:"rsvpAt,omitempty"`
}

// Invite represents an invite with its guests.
type Invite struct {
	ID     string  `json:"id"`
	Label  *string `json:"label,omitempty"`
	Guests []Guest `json:"guests"`
}

// ScheduleItem represents a timeline activity.
type ScheduleItem struct {
	ID          string    `json:"id"`
	Time        time.Time `json:"time"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	OrderIndex  int       `json:"orderIndex"`
}

// FAQ represents a question and answer.
type FAQ struct {
	ID         string `json:"id"`
	Question   string `json:"question"`
	Answer     string `json:"answer"`
	OrderIndex int    `json:"orderIndex"`
}

// GalleryItem represents a photo or video.
type GalleryItem struct {
	ID         string  `json:"id"`
	MediaType  string  `json:"mediaType"`
	MediaURL   string  `json:"mediaUrl"`
	Caption    *string `json:"caption,omitempty"`
	OrderIndex int     `json:"orderIndex"`
}

// Gift represents a gift suggestion.
type Gift struct {
	ID          string  `json:"id"`
	GiftType    string  `json:"giftType"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Link        *string `json:"link,omitempty"`
	OrderIndex  int     `json:"orderIndex"`
}

// EventDetailsResponse is the full response for the details endpoint.
type EventDetailsResponse struct {
	Event                Event          `json:"event"`
	Hosts                []Host         `json:"hosts"`
	Schedule             []ScheduleItem `json:"schedule"`
	FAQs                 []FAQ          `json:"faqs"`
	Gallery              []GalleryItem  `json:"gallery"`
	Gifts                []Gift         `json:"gifts"`
	Invite               *Invite        `json:"invite,omitempty"`
	ConfirmedGuestsCount int            `json:"confirmedGuestsCount"`
}

// ConfirmedGuestsResponse is the response for the confirmed-guests endpoint.
type ConfirmedGuestsResponse struct {
	Guests []ConfirmedGuest `json:"guests"`
	Count  int              `json:"count"`
}

// ConfirmedGuest is a minimal guest representation for the confirmed list.
type ConfirmedGuest struct {
	DisplayName string `json:"displayName"`
}

// RSVPRequest is the request body for the RSVP endpoint.
type RSVPRequest struct {
	InviteCode string  `json:"inviteCode"`
	GuestID    string  `json:"guestId"`
	Status     string  `json:"status"`
	Message    *string `json:"message,omitempty"`
}

// RSVPResponse is the response for the RSVP endpoint.
type RSVPResponse struct {
	ID          string     `json:"id"`
	DisplayName string     `json:"displayName"`
	RsvpStatus  string     `json:"rsvpStatus"`
	RsvpMessage *string    `json:"rsvpMessage,omitempty"`
	RsvpAt      *time.Time `json:"rsvpAt,omitempty"`
}

// ErrorResponse is a standard error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// --- Admin/Host Models ---

// BaseEvent contains shared fields for list views.
type BaseEvent struct {
	ID        string     `json:"id"`
	Type      string     `json:"type"`
	Slug      string     `json:"slug"`
	Title     string     `json:"title"`
	IsPublic  bool       `json:"isPublic"`
	StartsAt  *time.Time `json:"startsAt,omitempty"`
	Location  *string    `json:"location,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
}

// CreateEventRequest is the request body for creating an event.
type CreateEventRequest struct {
	Type     string  `json:"type"`
	Slug     string  `json:"slug"`
	Title    string  `json:"title"`
	IsPublic bool    `json:"isPublic"`
	StartsAt *string `json:"startsAt,omitempty"`
	Location *string `json:"location,omitempty"`
}

// AdminEvent is an event with hosts for admin list view.
type AdminEvent struct {
	BaseEvent
	Description      *string     `json:"description,omitempty"`
	CoverImageURL    *string     `json:"coverImageUrl,omitempty"`
	LocationPhotoURL *string     `json:"locationPhotoUrl,omitempty"`
	MusicURL         *string     `json:"musicUrl,omitempty"`
	Hosts            []AdminHost `json:"hosts"`
}

// UpdateEventRequest is the request body for updating an event.
type UpdateEventRequest struct {
	Type             string  `json:"type"`
	Slug             string  `json:"slug"`
	Title            string  `json:"title"`
	Description      *string `json:"description"`
	IsPublic         bool    `json:"isPublic"`
	StartsAt         *string `json:"startsAt"`
	Location         *string `json:"location"`
	CoverImageURL    *string `json:"coverImageUrl"`
	LocationPhotoURL *string `json:"locationPhotoUrl"`
	MusicURL         *string `json:"musicUrl"`
}

// AdminHost is a host with email for admin views.
type AdminHost struct {
	ID           string  `json:"id"`
	DisplayName  string  `json:"displayName"`
	ContactEmail *string `json:"contactEmail,omitempty"`
	AuthUserID   *string `json:"authUserId,omitempty"`
}

// AddHostRequest is the request body for adding a host to an event.
type AddHostRequest struct {
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
}

// AdminEventsResponse is the response for listing events.
type AdminEventsResponse struct {
	Events []AdminEvent `json:"events"`
}

// AdminInvite is an invite with guests for admin views.
type AdminInvite struct {
	ID         string       `json:"id"`
	InviteCode string       `json:"inviteCode"`
	Label      *string      `json:"label,omitempty"`
	CreatedAt  time.Time    `json:"createdAt"`
	Guests     []AdminGuest `json:"guests"`
}

// AdminGuest is a guest with full details for admin views.
type AdminGuest struct {
	ID          string     `json:"id"`
	DisplayName string     `json:"displayName"`
	RsvpStatus  string     `json:"rsvpStatus"`
	RsvpMessage *string    `json:"rsvpMessage,omitempty"`
	RsvpAt      *time.Time `json:"rsvpAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
}

// AdminInvitesResponse is the response for listing invites.
type AdminInvitesResponse struct {
	Invites []AdminInvite `json:"invites"`
}

// CreateInviteRequest is the request body for creating an invite.
type CreateInviteRequest struct {
	Label *string `json:"label"`
}

// AddGuestRequest is the request body for adding a guest to an invite.
type AddGuestRequest struct {
	DisplayName string `json:"displayName"`
}

// UpdateGuestRequest is the request body for updating a guest.
type UpdateGuestRequest struct {
	DisplayName string `json:"displayName"`
	RsvpStatus  string `json:"rsvpStatus"`
}

// AdminScheduleItem is a schedule item with createdAt for admin views.
type AdminScheduleItem struct {
	ID          string    `json:"id"`
	Time        time.Time `json:"time"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	OrderIndex  int       `json:"orderIndex"`
	CreatedAt   time.Time `json:"createdAt"`
}

// AdminScheduleResponse is the response for listing schedule items.
type AdminScheduleResponse struct {
	Items []AdminScheduleItem `json:"items"`
}

// CreateScheduleItemRequest is the request body for creating a schedule item.
type CreateScheduleItemRequest struct {
	Time        string  `json:"time"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	OrderIndex  *int    `json:"orderIndex"`
}

// UpdateScheduleItemRequest is the request body for updating a schedule item.
type UpdateScheduleItemRequest struct {
	Time        string  `json:"time"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	OrderIndex  *int    `json:"orderIndex"`
}

// AdminFAQ is a FAQ item with createdAt for admin views.
type AdminFAQ struct {
	ID         string    `json:"id"`
	Question   string    `json:"question"`
	Answer     string    `json:"answer"`
	OrderIndex int       `json:"orderIndex"`
	CreatedAt  time.Time `json:"createdAt"`
}

// AdminFAQsResponse is the response for listing FAQs.
type AdminFAQsResponse struct {
	Items []AdminFAQ `json:"items"`
}

// CreateFAQRequest is the request body for creating a FAQ.
type CreateFAQRequest struct {
	Question   string `json:"question"`
	Answer     string `json:"answer"`
	OrderIndex *int   `json:"orderIndex"`
}

// UpdateFAQRequest is the request body for updating a FAQ.
type UpdateFAQRequest struct {
	Question   string `json:"question"`
	Answer     string `json:"answer"`
	OrderIndex *int   `json:"orderIndex"`
}

// AdminGift is a gift item with createdAt for admin views.
type AdminGift struct {
	ID          string    `json:"id"`
	GiftType    string    `json:"giftType"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	Link        *string   `json:"link,omitempty"`
	OrderIndex  int       `json:"orderIndex"`
	CreatedAt   time.Time `json:"createdAt"`
}

// AdminGiftsResponse is the response for listing gifts.
type AdminGiftsResponse struct {
	Items []AdminGift `json:"items"`
}

// CreateGiftRequest is the request body for creating a gift.
type CreateGiftRequest struct {
	GiftType    string  `json:"giftType"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Link        *string `json:"link"`
	OrderIndex  *int    `json:"orderIndex"`
}

// UpdateGiftRequest is the request body for updating a gift.
type UpdateGiftRequest struct {
	GiftType    string  `json:"giftType"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Link        *string `json:"link"`
	OrderIndex  *int    `json:"orderIndex"`
}

// AdminGalleryItem is a gallery item with createdAt for admin views.
type AdminGalleryItem struct {
	ID         string    `json:"id"`
	MediaType  string    `json:"mediaType"`
	MediaURL   string    `json:"mediaUrl"`
	Caption    *string   `json:"caption,omitempty"`
	OrderIndex int       `json:"orderIndex"`
	CreatedAt  time.Time `json:"createdAt"`
}

// AdminGalleryResponse is the response for listing gallery items.
type AdminGalleryResponse struct {
	Items []AdminGalleryItem `json:"items"`
}

// CreateGalleryItemRequest is the request body for creating a gallery item.
type CreateGalleryItemRequest struct {
	MediaType  string  `json:"mediaType"`
	MediaURL   string  `json:"mediaUrl"`
	Caption    *string `json:"caption"`
	OrderIndex *int    `json:"orderIndex"`
}

// UpdateGalleryItemRequest is the request body for updating a gallery item.
type UpdateGalleryItemRequest struct {
	Caption    *string `json:"caption"`
	OrderIndex *int    `json:"orderIndex"`
}

// HostEvent is an alias for BaseEvent for host dashboard view.
type HostEvent = BaseEvent

// HostEventsResponse is the response for host events endpoint.
type HostEventsResponse struct {
	Events []HostEvent `json:"events"`
}

// HostGuestsResponse is the response for host guests endpoint.
type HostGuestsResponse struct {
	Guests []AdminGuest `json:"guests"`
}
