import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getHostEvents, getHostGuests, getHostInvites } from '@/services/api'
import type { HostEvent, HostGuest, AdminInvite, Host } from '@/types'
import { errorMsg } from '@/utils/error'

export const useHostStore = defineStore('host', () => {
  const events = ref<HostEvent[]>([])
  const selectedEvent = ref<HostEvent | null>(null)
  const guests = ref<HostGuest[]>([])
  const invites = ref<AdminInvite[]>([])
  const eventHosts = ref<Host[]>([])
  const eventsLoading = ref(false)
  const guestsLoading = ref(false)
  const invitesLoading = ref(false)
  const error = ref('')
  const showAllGuests = ref(false)

  /** Only confirmed guests (or all, based on toggle). */
  const filteredGuests = computed(() => {
    if (showAllGuests.value) return guests.value
    return guests.value.filter((g) => g.rsvpStatus === 'yes')
  })

  /** Fetches events for the authenticated host. */
  async function fetchEvents() {
    eventsLoading.value = true
    error.value = ''
    try {
      const data = await getHostEvents()
      events.value = data.events
    } catch (e: unknown) {
      error.value = errorMsg(e, 'Failed to load events')
    } finally {
      eventsLoading.value = false
    }
  }

  /** Selects an event and loads its guests. */
  async function selectEvent(event: HostEvent) {
    selectedEvent.value = event
    guestsLoading.value = true
    error.value = ''
    try {
      const data = await getHostGuests(event.id)
      guests.value = data.guests || []
    } catch (e: unknown) {
      error.value = errorMsg(e, 'Failed to load guests')
      guests.value = []
    } finally {
      guestsLoading.value = false
    }
  }

  /** Fetches invites for the selected event. */
  async function fetchInvites(eventId: string) {
    invitesLoading.value = true
    try {
      const data = await getHostInvites(eventId)
      invites.value = data.invites || []
      eventHosts.value = data.hosts || []
    } catch (e: unknown) {
      error.value = errorMsg(e, 'Failed to load invites')
      invites.value = []
      eventHosts.value = []
    } finally {
      invitesLoading.value = false
    }
  }

  /** Clears all host state. */
  function $reset() {
    events.value = []
    selectedEvent.value = null
    guests.value = []
    invites.value = []
    eventHosts.value = []
    eventsLoading.value = false
    guestsLoading.value = false
    invitesLoading.value = false
    error.value = ''
    showAllGuests.value = false
  }

  return {
    events,
    selectedEvent,
    guests,
    invites,
    eventHosts,
    eventsLoading,
    guestsLoading,
    invitesLoading,
    error,
    showAllGuests,
    filteredGuests,
    fetchEvents,
    selectEvent,
    fetchInvites,
    $reset,
  }
})
