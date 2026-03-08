import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  getEventSummary,
  getEventDetails,
  getConfirmedGuests,
  submitRSVP as apiSubmitRSVP,
} from '@/services/api'
import type {
  EventType,
  EventSummary,
  EventDetailsResponse,
  ConfirmedGuestsResponse,
  RSVPRequest,
  RSVPResponse,
  ComponentConfig,
} from '@/types'
import { errorMsg } from '@/utils/error'

export const useEventStore = defineStore('event', () => {
  const eventSummary = ref<EventSummary | null>(null)
  const eventDetails = ref<EventDetailsResponse | null>(null)
  const confirmedGuests = ref<ConfirmedGuestsResponse | null>(null)
  const loading = ref(false)
  const error = ref('')

  const event = computed(() => eventDetails.value?.event ?? null)
  const hosts = computed(() => eventDetails.value?.hosts ?? [])
  const schedule = computed(() => eventDetails.value?.schedule ?? [])
  const faqs = computed(() => eventDetails.value?.faqs ?? [])
  const gallery = computed(() => eventDetails.value?.gallery ?? [])
  const gifts = computed(() => eventDetails.value?.gifts ?? [])
  const invite = computed(() => eventDetails.value?.invite ?? null)
  const customContent = computed(() => eventDetails.value?.event?.customContent ?? null)
  const enabledComponents = computed(() => eventDetails.value?.event?.enabledComponents ?? null)

  /**
   * Ordered list of enabled components, with a sensible default fallback.
   * Each custom section gets its own entry (CustomSection:<id>) so they
   * can be positioned independently among other components.
   */
  const orderedComponents = computed<ComponentConfig[]>(() => {
    const sections = eventDetails.value?.event?.customContent?.customSections ?? []

    if (!enabledComponents.value?.components) {
      const defaults: ComponentConfig[] = [
        { name: 'EventDetails', enabled: true, order: 1 },
        { name: 'LocationPhoto', enabled: true, order: 2 },
        { name: 'CountdownTimer', enabled: true, order: 3 },
        { name: 'EventMap', enabled: true, order: 4 },
        { name: 'EventSchedule', enabled: true, order: 5 },
        { name: 'EventGallery', enabled: true, order: 6 },
        { name: 'AttireGuide', enabled: true, order: 7 },
        { name: 'EventFAQ', enabled: true, order: 8 },
        { name: 'MonetaryGifts', enabled: true, order: 9 },
        { name: 'GiftGuide', enabled: true, order: 10 },
      ]
      sections.forEach((s, i) => {
        defaults.push({ name: `CustomSection:${s.id}`, enabled: true, order: 11 + i })
      })
      return defaults
    }

    return enabledComponents.value.components
      .filter((c) => c.enabled)
      .sort((a, b) => a.order - b.order)
  })

  /**
   * Fetches the lightweight event summary (used by landing page).
   */
  async function fetchSummary(type: EventType, slug: string) {
    loading.value = true
    error.value = ''
    try {
      eventSummary.value = await getEventSummary(type, slug)
    } catch (e: unknown) {
      error.value = errorMsg(e, 'Failed to fetch event summary')
      throw e
    } finally {
      loading.value = false
    }
  }

  /**
   * Fetches full event details. Optionally include invite code for private events.
   */
  async function fetchDetails(type: EventType, slug: string, inviteCode?: string) {
    loading.value = true
    error.value = ''
    try {
      eventDetails.value = await getEventDetails(type, slug, inviteCode)
    } catch (e: unknown) {
      error.value = errorMsg(e, 'Failed to fetch event details')
      throw e
    } finally {
      loading.value = false
    }
  }

  /**
   * Fetches the confirmed guest list for display.
   */
  async function fetchConfirmedGuests(type: EventType, slug: string) {
    try {
      confirmedGuests.value = await getConfirmedGuests(type, slug)
    } catch (e: unknown) {
      error.value = errorMsg(e, 'Failed to fetch confirmed guests')
      throw e
    }
  }

  /**
   * Submits an RSVP response for a guest, then patches local state
   * so the UI updates immediately without a full re-fetch.
   */
  async function submitRSVP(
    type: EventType,
    slug: string,
    data: RSVPRequest,
  ): Promise<RSVPResponse> {
    const response = await apiSubmitRSVP(type, slug, data)

    if (eventDetails.value?.invite) {
      const guest = eventDetails.value.invite.guests.find(
        (g) => g.id === data.guestId,
      )
      if (guest) {
        guest.rsvpStatus = data.status
        guest.rsvpMessage = data.message
      }
    }

    return response
  }

  /**
   * Clears all event state (e.g. when navigating away).
   */
  function $reset() {
    eventSummary.value = null
    eventDetails.value = null
    confirmedGuests.value = null
    loading.value = false
    error.value = ''
  }

  return {
    eventSummary,
    eventDetails,
    confirmedGuests,
    loading,
    error,
    event,
    hosts,
    schedule,
    faqs,
    gallery,
    gifts,
    invite,
    customContent,
    enabledComponents,
    orderedComponents,
    fetchSummary,
    fetchDetails,
    fetchConfirmedGuests,
    submitRSVP,
    $reset,
  }
})
