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
  const guestFilter = ref<'all' | 'yes' | 'no' | 'pending'>('all')

  /** Guests filtered by the active status filter. */
  const filteredGuests = computed(() => {
    if (guestFilter.value === 'all') return guests.value
    return guests.value.filter((g) => g.rsvpStatus === guestFilter.value)
  })

  const totalCount = computed(() => guests.value.length)
  const yesCount = computed(() => guests.value.filter((g) => g.rsvpStatus === 'yes').length)
  const noCount = computed(() => guests.value.filter((g) => g.rsvpStatus === 'no').length)
  const pendingCount = computed(
    () => guests.value.filter((g) => g.rsvpStatus === 'pending').length,
  )

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
    guestFilter.value = 'all'
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
    guestFilter,
    filteredGuests,
    totalCount,
    yesCount,
    noCount,
    pendingCount,
    fetchEvents,
    selectEvent,
    fetchInvites,
    $reset,
  }
})
