// API service for backend calls
import type {
  EventType,
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
  AdminFAQ,
  AdminGift,
  AdminGalleryItem,
  CustomContent,
  EnabledComponents,
  HostEvent,
  HostGuest,
} from '@/types'
import { supabase } from './supabase'
import router from '@/router'

const API_BASE = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

/** Parses response JSON, throws with backend message on non-ok status. */
async function handleResponse<T>(res: Response, fallback: string): Promise<T> {
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.message || fallback)
  }
  if (res.status === 204 || res.headers.get('content-length') === '0') {
    return undefined as T
  }
  return res.json()
}

/**
 * Gets the current session's access token from the auth store.
 * Reads the reactive store instead of calling supabase.auth.getSession()
 * to avoid navigator.locks contention across tabs.
 */
async function getAuthToken(): Promise<string | null> {
  const { useAuthStore } = await import('@/stores/auth')
  const authStore = useAuthStore()
  return authStore.session?.access_token || null
}

/** Wrapper for authenticated fetch. Handles 401 by signing out + redirecting. */
async function authFetch(url: string, options: RequestInit = {}): Promise<Response> {
  const token = await getAuthToken()
  if (!token) throw new Error('Not authenticated')

  const headers = new Headers(options.headers)
  headers.set('Authorization', `Bearer ${token}`)

  const res = await fetch(url, { ...options, headers })

  if (res.status === 401) {
    await supabase.auth.signOut({ scope: 'local' })
    router.push('/host/login')
    throw new Error('Session expired')
  }

  return res
}

// --- Public Event API ---

/** Fetches event summary (basic info to determine public/private). */
export async function getEventSummary(type: EventType, slug: string): Promise<EventSummary> {
  const res = await fetch(`${API_BASE}/api/v1/events/${type}/${slug}/summary`)
  return handleResponse(res, 'Failed to fetch event summary')
}

/** Fetches full event details. For private events, include invite code. */
export async function getEventDetails(
  type: EventType,
  slug: string,
  inviteCode?: string,
): Promise<EventDetailsResponse> {
  const url = new URL(`${API_BASE}/api/v1/events/${type}/${slug}/details`)
  if (inviteCode) url.searchParams.set('invite', inviteCode)

  const res = await fetch(url.toString())
  return handleResponse(res, 'Failed to fetch event details')
}

/** Fetches list of confirmed guests. */
export async function getConfirmedGuests(
  type: EventType,
  slug: string,
): Promise<ConfirmedGuestsResponse> {
  const res = await fetch(`${API_BASE}/api/v1/events/${type}/${slug}/confirmed-guests`)
  return handleResponse(res, 'Failed to fetch confirmed guests')
}

/** Submits RSVP for a guest. */
export async function submitRSVP(
  type: EventType,
  slug: string,
  data: RSVPRequest,
): Promise<RSVPResponse> {
  const res = await fetch(`${API_BASE}/api/v1/events/${type}/${slug}/rsvp`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return handleResponse(res, 'Failed to submit RSVP')
}

// --- Auth API ---

/** Fetches the current user's email and admin status from the backend. */
export async function getMe(): Promise<{ email: string; isAdmin: boolean }> {
  const res = await authFetch(`${API_BASE}/api/v1/auth/me`)
  return handleResponse(res, 'Failed to fetch user info')
}

// --- Admin Event API ---

/** Fetches all events for admin dashboard. */
export async function adminListEvents(): Promise<{ events: AdminEvent[] }> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events`)
  return handleResponse(res, 'Failed to list events')
}

/** Creates a new event. */
export async function adminCreateEvent(data: {
  type: EventType
  slug: string
  title: string
  isPublic: boolean
  startsAt?: string
  location?: string
}): Promise<AdminEvent> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return handleResponse(res, 'Failed to create event')
}

/** Updates an existing event. */
export async function adminUpdateEvent(
  eventId: string,
  data: {
    type: EventType
    slug: string
    title: string
    description?: string | null
    isPublic: boolean
    startsAt?: string | null
    location?: string | null
    coverImageUrl?: string | null
    locationPhotoUrl?: string | null
    musicUrl?: string | null
  },
): Promise<AdminEvent> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return handleResponse(res, 'Failed to update event')
}

/** Deletes an event and all related data. */
export async function adminDeleteEvent(eventId: string): Promise<void> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}`, {
    method: 'DELETE',
  })
  await handleResponse<void>(res, 'Failed to delete event')
}

/** Adds a host to an event. */
export async function adminAddHost(
  eventId: string,
  data: { email: string; displayName: string },
): Promise<AdminHost> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}/hosts`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return handleResponse(res, 'Failed to add host')
}

/** Removes a host from an event. */
export async function adminDeleteHost(eventId: string, hostId: string): Promise<void> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}/hosts/${hostId}`, {
    method: 'DELETE',
  })
  await handleResponse<void>(res, 'Failed to delete host')
}

// --- Admin Invite/Guest API ---

/** Fetches all invites with guests for an event. */
export async function adminListInvites(eventId: string): Promise<{ invites: AdminInvite[] }> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}/invites`)
  return handleResponse(res, 'Failed to list invites')
}

/** Creates a new invite with auto-generated code. */
export async function adminCreateInvite(
  eventId: string,
  data: { label?: string | null },
): Promise<AdminInvite> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}/invites`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return handleResponse(res, 'Failed to create invite')
}

/** Deletes an invite and all its guests. */
export async function adminDeleteInvite(eventId: string, inviteId: string): Promise<void> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}/invites/${inviteId}`, {
    method: 'DELETE',
  })
  await handleResponse<void>(res, 'Failed to delete invite')
}

/** Adds a guest to an invite. */
export async function adminAddGuest(
  eventId: string,
  inviteId: string,
  data: { displayName: string },
): Promise<AdminGuest> {
  const res = await authFetch(
    `${API_BASE}/api/v1/admin/events/${eventId}/invites/${inviteId}/guests`,
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    },
  )
  return handleResponse(res, 'Failed to add guest')
}

/** Updates a guest's name or RSVP status. */
export async function adminUpdateGuest(
  eventId: string,
  guestId: string,
  data: { displayName: string; rsvpStatus: 'pending' | 'yes' | 'no' },
): Promise<AdminGuest> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}/guests/${guestId}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return handleResponse(res, 'Failed to update guest')
}

/** Removes a guest. */
export async function adminDeleteGuest(eventId: string, guestId: string): Promise<void> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}/guests/${guestId}`, {
    method: 'DELETE',
  })
  await handleResponse<void>(res, 'Failed to delete guest')
}

// --- Admin Schedule API ---

/** Fetches all schedule items for an event. */
export async function adminListSchedule(
  eventId: string,
): Promise<{ items: AdminScheduleItem[] }> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}/schedule`)
  return handleResponse(res, 'Failed to list schedule')
}

/** Creates a new schedule item. */
export async function adminCreateScheduleItem(
  eventId: string,
  data: { time: string; title: string; description?: string | null; orderIndex?: number | null },
): Promise<AdminScheduleItem> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}/schedule`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return handleResponse(res, 'Failed to create schedule item')
}

/** Updates a schedule item. */
export async function adminUpdateScheduleItem(
  eventId: string,
  itemId: string,
  data: { time: string; title: string; description?: string | null; orderIndex?: number | null },
): Promise<AdminScheduleItem> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}/schedule/${itemId}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return handleResponse(res, 'Failed to update schedule item')
}

/** Deletes a schedule item. */
export async function adminDeleteScheduleItem(
  eventId: string,
  itemId: string,
): Promise<void> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}/schedule/${itemId}`, {
    method: 'DELETE',
  })
  await handleResponse<void>(res, 'Failed to delete schedule item')
}

// --- Admin FAQ API ---

/** Fetches all FAQs for an event. */
export async function adminListFAQs(eventId: string): Promise<{ items: AdminFAQ[] }> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}/faqs`)
  return handleResponse(res, 'Failed to list FAQs')
}

/** Creates a new FAQ. */
export async function adminCreateFAQ(
  eventId: string,
  data: { question: string; answer: string; orderIndex?: number | null },
): Promise<AdminFAQ> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}/faqs`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return handleResponse(res, 'Failed to create FAQ')
}

/** Updates a FAQ. */
export async function adminUpdateFAQ(
  eventId: string,
  itemId: string,
  data: { question: string; answer: string; orderIndex?: number | null },
): Promise<AdminFAQ> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}/faqs/${itemId}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return handleResponse(res, 'Failed to update FAQ')
}

/** Deletes a FAQ. */
export async function adminDeleteFAQ(eventId: string, itemId: string): Promise<void> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}/faqs/${itemId}`, {
    method: 'DELETE',
  })
  await handleResponse<void>(res, 'Failed to delete FAQ')
}

// --- Admin Gift API ---

/** Fetches all gifts for an event. */
export async function adminListGifts(eventId: string): Promise<{ items: AdminGift[] }> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}/gifts`)
  return handleResponse(res, 'Failed to list gifts')
}

/** Creates a new gift. */
export async function adminCreateGift(
  eventId: string,
  data: {
    giftType: 'physical' | 'monetary'
    title: string
    description?: string | null
    link?: string | null
    orderIndex?: number | null
  },
): Promise<AdminGift> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}/gifts`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return handleResponse(res, 'Failed to create gift')
}

/** Updates a gift. */
export async function adminUpdateGift(
  eventId: string,
  itemId: string,
  data: {
    giftType: 'physical' | 'monetary'
    title: string
    description?: string | null
    link?: string | null
    orderIndex?: number | null
  },
): Promise<AdminGift> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}/gifts/${itemId}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return handleResponse(res, 'Failed to update gift')
}

/** Deletes a gift. */
export async function adminDeleteGift(eventId: string, itemId: string): Promise<void> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}/gifts/${itemId}`, {
    method: 'DELETE',
  })
  await handleResponse<void>(res, 'Failed to delete gift')
}

// --- Admin Gallery API ---

/** Fetches all gallery items for an event. */
export async function adminListGallery(eventId: string): Promise<{ items: AdminGalleryItem[] }> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}/gallery`)
  return handleResponse(res, 'Failed to list gallery')
}

/** Creates a new gallery item. */
export async function adminCreateGalleryItem(
  eventId: string,
  data: {
    mediaType: 'photo' | 'video'
    mediaUrl: string
    caption?: string | null
    orderIndex?: number | null
  },
): Promise<AdminGalleryItem> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}/gallery`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return handleResponse(res, 'Failed to create gallery item')
}

/** Updates a gallery item (caption and order). */
export async function adminUpdateGalleryItem(
  eventId: string,
  itemId: string,
  data: { caption?: string | null; orderIndex?: number | null },
): Promise<AdminGalleryItem> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}/gallery/${itemId}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return handleResponse(res, 'Failed to update gallery item')
}

/** Deletes a gallery item. */
export async function adminDeleteGalleryItem(eventId: string, itemId: string): Promise<void> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}/gallery/${itemId}`, {
    method: 'DELETE',
  })
  await handleResponse<void>(res, 'Failed to delete gallery item')
}

// --- Admin Custom Content API ---

/** Fetches the custom_content JSONB for an event. */
export async function adminGetCustomContent(eventId: string): Promise<CustomContent> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}/custom-content`)
  return handleResponse(res, 'Failed to get custom content')
}

/** Replaces the custom_content JSONB for an event. */
export async function adminUpdateCustomContent(
  eventId: string,
  data: CustomContent,
): Promise<CustomContent> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}/custom-content`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return handleResponse(res, 'Failed to update custom content')
}

// --- Admin Enabled Components API ---

/** Fetches the enabled_components JSONB for an event. */
export async function adminGetEnabledComponents(eventId: string): Promise<EnabledComponents> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}/enabled-components`)
  return handleResponse(res, 'Failed to get enabled components')
}

/** Replaces the enabled_components JSONB for an event. */
export async function adminUpdateEnabledComponents(
  eventId: string,
  data: EnabledComponents,
): Promise<EnabledComponents> {
  const res = await authFetch(`${API_BASE}/api/v1/admin/events/${eventId}/enabled-components`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return handleResponse(res, 'Failed to update enabled components')
}

// --- Host API ---

/** Fetches events for the authenticated host. */
export async function getHostEvents(): Promise<{ events: HostEvent[] }> {
  const res = await authFetch(`${API_BASE}/api/v1/host/events`)
  return handleResponse(res, 'Failed to fetch events')
}

/** Fetches guests for an event (host must be linked). */
export async function getHostGuests(eventId: string): Promise<{ guests: HostGuest[] }> {
  const res = await authFetch(`${API_BASE}/api/v1/host/events/${eventId}/guests`)
  return handleResponse(res, 'Failed to fetch guests')
}
