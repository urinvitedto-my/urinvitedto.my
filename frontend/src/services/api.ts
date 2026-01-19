// API service for backend calls
import type {
  EventSummary,
  EventDetailsResponse,
  ConfirmedGuestsResponse,
  RSVPRequest,
  RSVPResponse,
} from '@/types'

const API_BASE = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

/**
 * Fetches event summary (basic info to determine public/private).
 */
export async function getEventSummary(type: string, slug: string): Promise<EventSummary> {
  const res = await fetch(`${API_BASE}/api/v1/events/${type}/${slug}/summary`)
  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.message || 'Failed to fetch event summary')
  }
  return res.json()
}

/**
 * Fetches full event details. For private events, include invite code.
 */
export async function getEventDetails(
  type: string,
  slug: string,
  inviteCode?: string
): Promise<EventDetailsResponse> {
  const url = new URL(`${API_BASE}/api/v1/events/${type}/${slug}/details`)
  if (inviteCode) {
    url.searchParams.set('invite', inviteCode)
  }

  const res = await fetch(url.toString())
  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.message || 'Failed to fetch event details')
  }
  return res.json()
}

/**
 * Fetches list of confirmed guests.
 */
export async function getConfirmedGuests(
  type: string,
  slug: string
): Promise<ConfirmedGuestsResponse> {
  const res = await fetch(`${API_BASE}/api/v1/events/${type}/${slug}/confirmed-guests`)
  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.message || 'Failed to fetch confirmed guests')
  }
  return res.json()
}

/**
 * Submits RSVP for a guest.
 */
export async function submitRSVP(
  type: string,
  slug: string,
  data: RSVPRequest
): Promise<RSVPResponse> {
  const res = await fetch(`${API_BASE}/api/v1/events/${type}/${slug}/rsvp`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.message || 'Failed to submit RSVP')
  }
  return res.json()
}
