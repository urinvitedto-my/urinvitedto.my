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
