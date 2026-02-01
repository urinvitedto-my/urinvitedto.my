// API service for backend calls
import type {
  EventSummary,
  EventDetailsResponse,
  ConfirmedGuestsResponse,
  RSVPRequest,
  RSVPResponse,
  AdminEvent,
  AdminHost,
} from '@/types'
import { supabase } from './supabase'

const API_BASE = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

/**
 * Gets the current session's access token for authenticated requests.
 */
async function getAuthToken(): Promise<string | null> {
  const { data } = await supabase.auth.getSession()
  return data.session?.access_token || null
}

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

// --- Admin API Functions ---

/**
 * Fetches all events for admin dashboard.
 */
export async function adminListEvents(): Promise<{ events: AdminEvent[] }> {
  const token = await getAuthToken()
  if (!token) throw new Error('Not authenticated')

  const res = await fetch(`${API_BASE}/api/v1/admin/events`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.message || 'Failed to list events')
  }
  return res.json()
}

/**
 * Creates a new event.
 */
export async function adminCreateEvent(data: {
  type: string
  slug: string
  title: string
  isPublic: boolean
  startsAt?: string
  location?: string
}): Promise<AdminEvent> {
  const token = await getAuthToken()
  if (!token) throw new Error('Not authenticated')

  const res = await fetch(`${API_BASE}/api/v1/admin/events`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.message || 'Failed to create event')
  }
  return res.json()
}

/**
 * Adds a host to an event.
 */
export async function adminAddHost(
  eventId: string,
  data: { email: string; displayName: string }
): Promise<AdminHost> {
  const token = await getAuthToken()
  if (!token) throw new Error('Not authenticated')

  const res = await fetch(`${API_BASE}/api/v1/admin/events/${eventId}/hosts`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.message || 'Failed to add host')
  }
  return res.json()
}

/**
 * Removes a host from an event.
 */
export async function adminDeleteHost(eventId: string, hostId: string): Promise<void> {
  const token = await getAuthToken()
  if (!token) throw new Error('Not authenticated')

  const res = await fetch(`${API_BASE}/api/v1/admin/events/${eventId}/hosts/${hostId}`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.message || 'Failed to delete host')
  }
}

// --- Host API Functions ---

export interface HostEvent {
  id: string
  type: string
  slug: string
  title: string
  isPublic: boolean
  startsAt?: string
  location?: string
  createdAt: string
}

/**
 * Fetches events for the authenticated host.
 */
export async function getHostEvents(): Promise<{ events: HostEvent[] }> {
  const token = await getAuthToken()
  if (!token) throw new Error('Not authenticated')

  const res = await fetch(`${API_BASE}/api/v1/host/events`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.message || 'Failed to fetch events')
  }
  return res.json()
}
