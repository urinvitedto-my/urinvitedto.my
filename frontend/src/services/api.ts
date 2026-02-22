// API service for backend calls
import type {
  EventSummary,
  EventDetailsResponse,
  ConfirmedGuestsResponse,
  RSVPRequest,
  RSVPResponse,
  AdminEvent,
  AdminHost,
  AdminInvite,
  AdminGuest,
  AdminScheduleItem,
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

/**
 * Updates an existing event.
 */
export async function adminUpdateEvent(
  eventId: string,
  data: {
    type: string
    slug: string
    title: string
    description?: string | null
    isPublic: boolean
    startsAt?: string | null
    location?: string | null
    coverImageUrl?: string | null
    locationPhotoUrl?: string | null
  }
): Promise<AdminEvent> {
  const token = await getAuthToken()
  if (!token) throw new Error('Not authenticated')

  const res = await fetch(`${API_BASE}/api/v1/admin/events/${eventId}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.message || 'Failed to update event')
  }
  return res.json()
}

/**
 * Deletes an event and all related data.
 */
export async function adminDeleteEvent(eventId: string): Promise<void> {
  const token = await getAuthToken()
  if (!token) throw new Error('Not authenticated')

  const res = await fetch(`${API_BASE}/api/v1/admin/events/${eventId}`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.message || 'Failed to delete event')
  }
}

// --- Admin Invite/Guest API Functions ---

/**
 * Fetches all invites with guests for an event.
 */
export async function adminListInvites(eventId: string): Promise<{ invites: AdminInvite[] }> {
  const token = await getAuthToken()
  if (!token) throw new Error('Not authenticated')

  const res = await fetch(`${API_BASE}/api/v1/admin/events/${eventId}/invites`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.message || 'Failed to list invites')
  }
  return res.json()
}

/**
 * Creates a new invite with auto-generated code.
 */
export async function adminCreateInvite(
  eventId: string,
  data: { label?: string | null }
): Promise<AdminInvite> {
  const token = await getAuthToken()
  if (!token) throw new Error('Not authenticated')

  const res = await fetch(`${API_BASE}/api/v1/admin/events/${eventId}/invites`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.message || 'Failed to create invite')
  }
  return res.json()
}

/**
 * Deletes an invite and all its guests.
 */
export async function adminDeleteInvite(eventId: string, inviteId: string): Promise<void> {
  const token = await getAuthToken()
  if (!token) throw new Error('Not authenticated')

  const res = await fetch(`${API_BASE}/api/v1/admin/events/${eventId}/invites/${inviteId}`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.message || 'Failed to delete invite')
  }
}

/**
 * Adds a guest to an invite.
 */
export async function adminAddGuest(
  eventId: string,
  inviteId: string,
  data: { displayName: string }
): Promise<AdminGuest> {
  const token = await getAuthToken()
  if (!token) throw new Error('Not authenticated')

  const res = await fetch(`${API_BASE}/api/v1/admin/events/${eventId}/invites/${inviteId}/guests`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.message || 'Failed to add guest')
  }
  return res.json()
}

/**
 * Updates a guest's name or RSVP status.
 */
export async function adminUpdateGuest(
  eventId: string,
  guestId: string,
  data: { displayName: string; rsvpStatus: string }
): Promise<AdminGuest> {
  const token = await getAuthToken()
  if (!token) throw new Error('Not authenticated')

  const res = await fetch(`${API_BASE}/api/v1/admin/events/${eventId}/guests/${guestId}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.message || 'Failed to update guest')
  }
  return res.json()
}

/**
 * Removes a guest.
 */
export async function adminDeleteGuest(eventId: string, guestId: string): Promise<void> {
  const token = await getAuthToken()
  if (!token) throw new Error('Not authenticated')

  const res = await fetch(`${API_BASE}/api/v1/admin/events/${eventId}/guests/${guestId}`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.message || 'Failed to delete guest')
  }
}

// --- Admin Schedule API Functions ---

/**
 * Fetches all schedule items for an event.
 */
export async function adminListSchedule(eventId: string): Promise<{ items: AdminScheduleItem[] }> {
  const token = await getAuthToken()
  if (!token) throw new Error('Not authenticated')

  const res = await fetch(`${API_BASE}/api/v1/admin/events/${eventId}/schedule`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.message || 'Failed to list schedule')
  }
  return res.json()
}

/**
 * Creates a new schedule item.
 */
export async function adminCreateScheduleItem(
  eventId: string,
  data: { time: string; title: string; description?: string | null; orderIndex?: number | null }
): Promise<AdminScheduleItem> {
  const token = await getAuthToken()
  if (!token) throw new Error('Not authenticated')

  const res = await fetch(`${API_BASE}/api/v1/admin/events/${eventId}/schedule`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.message || 'Failed to create schedule item')
  }
  return res.json()
}

/**
 * Updates a schedule item.
 */
export async function adminUpdateScheduleItem(
  eventId: string,
  itemId: string,
  data: { time: string; title: string; description?: string | null; orderIndex?: number | null }
): Promise<AdminScheduleItem> {
  const token = await getAuthToken()
  if (!token) throw new Error('Not authenticated')

  const res = await fetch(`${API_BASE}/api/v1/admin/events/${eventId}/schedule/${itemId}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  })
  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.message || 'Failed to update schedule item')
  }
  return res.json()
}

/**
 * Deletes a schedule item.
 */
export async function adminDeleteScheduleItem(eventId: string, itemId: string): Promise<void> {
  const token = await getAuthToken()
  if (!token) throw new Error('Not authenticated')

  const res = await fetch(`${API_BASE}/api/v1/admin/events/${eventId}/schedule/${itemId}`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.message || 'Failed to delete schedule item')
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
