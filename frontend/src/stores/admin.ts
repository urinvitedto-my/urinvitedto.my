import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  adminListEvents,
  adminCreateEvent,
  adminUpdateEvent,
  adminDeleteEvent,
  adminAddHost,
  adminDeleteHost,
} from '@/services/api'
import type { AdminEvent, AdminHost } from '@/types'

export const useAdminStore = defineStore('admin', () => {
  const events = ref<AdminEvent[]>([])
  const selectedEventId = ref<string | null>(null)
  const loading = ref(false)
  const error = ref('')

  const selectedEvent = computed(() =>
    events.value.find((e) => e.id === selectedEventId.value) ?? null,
  )

  /**
   * Loads all events from the admin API.
   */
  async function fetchEvents() {
    loading.value = true
    error.value = ''
    try {
      const data = await adminListEvents()
      events.value = data.events
    } catch (e: any) {
      error.value = e.message || 'Failed to load events'
      throw e
    } finally {
      loading.value = false
    }
  }

  /**
   * Creates a new event and prepends it to the list.
   */
  async function createEvent(data: {
    type: string
    slug: string
    title: string
    isPublic: boolean
    startsAt?: string
    location?: string
  }): Promise<AdminEvent> {
    const newEvent = await adminCreateEvent(data)
    events.value.unshift(newEvent)
    return newEvent
  }

  /**
   * Updates an existing event in-place.
   */
  async function updateEvent(
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
      musicUrl?: string | null
    },
  ): Promise<AdminEvent> {
    const updated = await adminUpdateEvent(eventId, data)
    const idx = events.value.findIndex((e) => e.id === eventId)
    if (idx !== -1) {
      events.value[idx] = updated
    }
    return updated
  }

  /**
   * Deletes an event and removes it from the list.
   */
  async function deleteEvent(eventId: string) {
    await adminDeleteEvent(eventId)
    events.value = events.value.filter((e) => e.id !== eventId)
  }

  /**
   * Adds a host to an event and updates local state.
   */
  async function addHost(
    eventId: string,
    data: { email: string; displayName: string },
  ): Promise<AdminHost> {
    const newHost = await adminAddHost(eventId, data)
    const event = events.value.find((e) => e.id === eventId)
    if (event) {
      event.hosts.push(newHost)
    }
    return newHost
  }

  /**
   * Removes a host from an event and updates local state.
   */
  async function deleteHost(eventId: string, hostId: string) {
    await adminDeleteHost(eventId, hostId)
    const event = events.value.find((e) => e.id === eventId)
    if (event) {
      event.hosts = event.hosts.filter((h) => h.id !== hostId)
    }
  }

  return {
    events,
    selectedEventId,
    loading,
    error,
    selectedEvent,
    fetchEvents,
    createEvent,
    updateEvent,
    deleteEvent,
    addHost,
    deleteHost,
  }
})
