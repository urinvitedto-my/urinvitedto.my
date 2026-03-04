import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  adminListEvents,
  adminCreateEvent,
  adminUpdateEvent,
  adminDeleteEvent,
  adminAddHost,
  adminDeleteHost,
  adminListInvites,
  adminCreateInvite,
  adminDeleteInvite,
  adminAddGuest,
  adminUpdateGuest,
  adminDeleteGuest,
  adminListSchedule,
  adminCreateScheduleItem,
  adminUpdateScheduleItem,
  adminDeleteScheduleItem,
  adminListFAQs,
  adminCreateFAQ,
  adminUpdateFAQ,
  adminDeleteFAQ,
  adminListGifts,
  adminCreateGift,
  adminUpdateGift,
  adminDeleteGift,
  adminListGallery,
  adminCreateGalleryItem,
  adminUpdateGalleryItem,
  adminDeleteGalleryItem,
  adminGetCustomContent,
  adminUpdateCustomContent,
  adminGetEnabledComponents,
  adminUpdateEnabledComponents,
} from '@/services/api'
import type {
  EventType,
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
  ComponentConfig,
} from '@/types'
import { errorMsg } from '@/utils/error'

export const useAdminStore = defineStore('admin', () => {
  // --- Events ---
  const events = ref<AdminEvent[]>([])
  const loading = ref(false)
  const error = ref('')

  const selectedEvent = computed(() =>
    events.value.find((e) => e.id === selectedEventId.value) ?? null,
  )
  const selectedEventId = ref<string | null>(null)

  // --- Sub-entity data, keyed by eventId ---
  const invites = ref<Record<string, AdminInvite[]>>({})
  const schedule = ref<Record<string, AdminScheduleItem[]>>({})
  const faqs = ref<Record<string, AdminFAQ[]>>({})
  const gifts = ref<Record<string, AdminGift[]>>({})
  const gallery = ref<Record<string, AdminGalleryItem[]>>({})
  const customContent = ref<Record<string, CustomContent>>({})
  const enabledComponents = ref<Record<string, ComponentConfig[]>>({})

  // Per sub-entity loading/error keyed as "type:eventId"
  const subLoading = ref<Record<string, boolean>>({})
  const subError = ref<Record<string, string>>({})

  function subKey(type: string, eventId: string) {
    return `${type}:${eventId}`
  }

  function isSubLoading(type: string, eventId: string) {
    return subLoading.value[subKey(type, eventId)] ?? false
  }

  function getSubError(type: string, eventId: string) {
    return subError.value[subKey(type, eventId)] ?? ''
  }

  // ========================
  // Event CRUD
  // ========================

  /** Loads all events from the admin API. */
  async function fetchEvents() {
    loading.value = true
    error.value = ''
    try {
      const data = await adminListEvents()
      events.value = data.events
    } catch (e: unknown) {
      error.value = errorMsg(e, 'Failed to load events')
      throw e
    } finally {
      loading.value = false
    }
  }

  /** Creates a new event and prepends it to the list. */
  async function createEvent(data: {
    type: EventType
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

  /** Updates an existing event in-place. */
  async function updateEvent(
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
    const updated = await adminUpdateEvent(eventId, data)
    const idx = events.value.findIndex((e) => e.id === eventId)
    if (idx !== -1) events.value[idx] = updated
    return updated
  }

  /** Deletes an event and removes it from the list. */
  async function deleteEvent(eventId: string) {
    await adminDeleteEvent(eventId)
    events.value = events.value.filter((e) => e.id !== eventId)
    // Clean up sub-entity data
    delete invites.value[eventId]
    delete schedule.value[eventId]
    delete faqs.value[eventId]
    delete gifts.value[eventId]
    delete gallery.value[eventId]
    delete customContent.value[eventId]
    delete enabledComponents.value[eventId]
  }

  // ========================
  // Host CRUD
  // ========================

  /** Adds a host to an event and updates local state. */
  async function addHost(
    eventId: string,
    data: { email: string; displayName: string },
  ): Promise<AdminHost> {
    const newHost = await adminAddHost(eventId, data)
    const event = events.value.find((e) => e.id === eventId)
    if (event) event.hosts.push(newHost)
    return newHost
  }

  /** Removes a host from an event and updates local state. */
  async function deleteHost(eventId: string, hostId: string) {
    await adminDeleteHost(eventId, hostId)
    const event = events.value.find((e) => e.id === eventId)
    if (event) event.hosts = event.hosts.filter((h) => h.id !== hostId)
  }

  // ========================
  // Invite CRUD
  // ========================

  function getInvites(eventId: string): AdminInvite[] {
    return invites.value[eventId] ?? []
  }

  async function fetchInvites(eventId: string) {
    const key = subKey('invites', eventId)
    subLoading.value[key] = true
    subError.value[key] = ''
    try {
      const data = await adminListInvites(eventId)
      invites.value[eventId] = data.invites
    } catch (e: unknown) {
      subError.value[key] = errorMsg(e, 'Failed to load invites')
    } finally {
      subLoading.value[key] = false
    }
  }

  async function createInvite(eventId: string, data: { label?: string | null }) {
    const newInvite = await adminCreateInvite(eventId, data)
    if (!invites.value[eventId]) invites.value[eventId] = []
    invites.value[eventId].unshift(newInvite)
    return newInvite
  }

  async function deleteInvite(eventId: string, inviteId: string) {
    await adminDeleteInvite(eventId, inviteId)
    if (invites.value[eventId]) {
      invites.value[eventId] = invites.value[eventId].filter((i) => i.id !== inviteId)
    }
  }

  async function addGuestToInvite(
    eventId: string,
    inviteId: string,
    data: { displayName: string },
  ) {
    const newGuest = await adminAddGuest(eventId, inviteId, data)
    const invite = invites.value[eventId]?.find((i) => i.id === inviteId)
    if (invite) invite.guests.push(newGuest)
    return newGuest
  }

  async function updateGuestInInvite(
    eventId: string,
    inviteId: string,
    guestId: string,
    data: { displayName: string; rsvpStatus: 'pending' | 'yes' | 'no' },
  ) {
    const updated = await adminUpdateGuest(eventId, guestId, data)
    const invite = invites.value[eventId]?.find((i) => i.id === inviteId)
    if (invite) {
      const idx = invite.guests.findIndex((g) => g.id === guestId)
      if (idx !== -1) invite.guests[idx] = updated
    }
    return updated
  }

  async function deleteGuestFromInvite(eventId: string, inviteId: string, guestId: string) {
    await adminDeleteGuest(eventId, guestId)
    const invite = invites.value[eventId]?.find((i) => i.id === inviteId)
    if (invite) {
      invite.guests = invite.guests.filter((g) => g.id !== guestId)
    }
  }

  // ========================
  // Schedule CRUD
  // ========================

  function getSchedule(eventId: string): AdminScheduleItem[] {
    return schedule.value[eventId] ?? []
  }

  async function fetchSchedule(eventId: string) {
    const key = subKey('schedule', eventId)
    subLoading.value[key] = true
    subError.value[key] = ''
    try {
      const data = await adminListSchedule(eventId)
      schedule.value[eventId] = data.items
    } catch (e: unknown) {
      subError.value[key] = errorMsg(e, 'Failed to load schedule')
    } finally {
      subLoading.value[key] = false
    }
  }

  async function createScheduleItem(
    eventId: string,
    data: { time: string; title: string; description?: string | null },
  ) {
    const newItem = await adminCreateScheduleItem(eventId, data)
    if (!schedule.value[eventId]) schedule.value[eventId] = []
    schedule.value[eventId].push(newItem)
    schedule.value[eventId].sort((a, b) => a.orderIndex - b.orderIndex)
    return newItem
  }

  async function updateScheduleItem(
    eventId: string,
    itemId: string,
    data: { time: string; title: string; description?: string | null; orderIndex?: number | null },
  ) {
    const updated = await adminUpdateScheduleItem(eventId, itemId, data)
    const items = schedule.value[eventId]
    if (items) {
      const idx = items.findIndex((i) => i.id === itemId)
      if (idx !== -1) items[idx] = updated
      items.sort((a, b) => a.orderIndex - b.orderIndex)
    }
    return updated
  }

  async function deleteScheduleItem(eventId: string, itemId: string) {
    await adminDeleteScheduleItem(eventId, itemId)
    if (schedule.value[eventId]) {
      schedule.value[eventId] = schedule.value[eventId].filter((i) => i.id !== itemId)
    }
  }

  // ========================
  // FAQ CRUD
  // ========================

  function getFAQs(eventId: string): AdminFAQ[] {
    return faqs.value[eventId] ?? []
  }

  async function fetchFAQs(eventId: string) {
    const key = subKey('faqs', eventId)
    subLoading.value[key] = true
    subError.value[key] = ''
    try {
      const data = await adminListFAQs(eventId)
      faqs.value[eventId] = data.items
    } catch (e: unknown) {
      subError.value[key] = errorMsg(e, 'Failed to load FAQs')
    } finally {
      subLoading.value[key] = false
    }
  }

  async function createFAQItem(
    eventId: string,
    data: { question: string; answer: string },
  ) {
    const newItem = await adminCreateFAQ(eventId, data)
    if (!faqs.value[eventId]) faqs.value[eventId] = []
    faqs.value[eventId].push(newItem)
    faqs.value[eventId].sort((a, b) => a.orderIndex - b.orderIndex)
    return newItem
  }

  async function updateFAQItem(
    eventId: string,
    itemId: string,
    data: { question: string; answer: string; orderIndex?: number | null },
  ) {
    const updated = await adminUpdateFAQ(eventId, itemId, data)
    const items = faqs.value[eventId]
    if (items) {
      const idx = items.findIndex((i) => i.id === itemId)
      if (idx !== -1) items[idx] = updated
      items.sort((a, b) => a.orderIndex - b.orderIndex)
    }
    return updated
  }

  async function deleteFAQItem(eventId: string, itemId: string) {
    await adminDeleteFAQ(eventId, itemId)
    if (faqs.value[eventId]) {
      faqs.value[eventId] = faqs.value[eventId].filter((i) => i.id !== itemId)
    }
  }

  // ========================
  // Gift CRUD
  // ========================

  function getGifts(eventId: string): AdminGift[] {
    return gifts.value[eventId] ?? []
  }

  async function fetchGifts(eventId: string) {
    const key = subKey('gifts', eventId)
    subLoading.value[key] = true
    subError.value[key] = ''
    try {
      const data = await adminListGifts(eventId)
      gifts.value[eventId] = data.items
    } catch (e: unknown) {
      subError.value[key] = errorMsg(e, 'Failed to load gifts')
    } finally {
      subLoading.value[key] = false
    }
  }

  async function createGiftItem(
    eventId: string,
    data: {
      giftType: 'physical' | 'monetary'
      title: string
      description?: string | null
      link?: string | null
    },
  ) {
    const newItem = await adminCreateGift(eventId, data)
    if (!gifts.value[eventId]) gifts.value[eventId] = []
    gifts.value[eventId].push(newItem)
    gifts.value[eventId].sort((a, b) => a.orderIndex - b.orderIndex)
    return newItem
  }

  async function updateGiftItem(
    eventId: string,
    itemId: string,
    data: {
      giftType: 'physical' | 'monetary'
      title: string
      description?: string | null
      link?: string | null
      orderIndex?: number | null
    },
  ) {
    const updated = await adminUpdateGift(eventId, itemId, data)
    const items = gifts.value[eventId]
    if (items) {
      const idx = items.findIndex((i) => i.id === itemId)
      if (idx !== -1) items[idx] = updated
      items.sort((a, b) => a.orderIndex - b.orderIndex)
    }
    return updated
  }

  async function deleteGiftItem(eventId: string, itemId: string) {
    await adminDeleteGift(eventId, itemId)
    if (gifts.value[eventId]) {
      gifts.value[eventId] = gifts.value[eventId].filter((i) => i.id !== itemId)
    }
  }

  // ========================
  // Gallery CRUD
  // ========================

  function getGallery(eventId: string): AdminGalleryItem[] {
    return gallery.value[eventId] ?? []
  }

  async function fetchGallery(eventId: string) {
    const key = subKey('gallery', eventId)
    subLoading.value[key] = true
    subError.value[key] = ''
    try {
      const data = await adminListGallery(eventId)
      gallery.value[eventId] = data.items
    } catch (e: unknown) {
      subError.value[key] = errorMsg(e, 'Failed to load gallery')
    } finally {
      subLoading.value[key] = false
    }
  }

  async function createGalleryItem(
    eventId: string,
    data: {
      mediaType: 'photo' | 'video'
      mediaUrl: string
      caption?: string | null
      orderIndex?: number | null
    },
  ) {
    const newItem = await adminCreateGalleryItem(eventId, data)
    if (!gallery.value[eventId]) gallery.value[eventId] = []
    gallery.value[eventId].push(newItem)
    gallery.value[eventId].sort((a, b) => a.orderIndex - b.orderIndex)
    return newItem
  }

  async function updateGalleryItem(
    eventId: string,
    itemId: string,
    data: { caption?: string | null; orderIndex?: number | null },
  ) {
    const updated = await adminUpdateGalleryItem(eventId, itemId, data)
    const items = gallery.value[eventId]
    if (items) {
      const idx = items.findIndex((i) => i.id === itemId)
      if (idx !== -1) items[idx] = updated
      items.sort((a, b) => a.orderIndex - b.orderIndex)
    }
    return updated
  }

  async function deleteGalleryItemFromStore(eventId: string, itemId: string) {
    await adminDeleteGalleryItem(eventId, itemId)
    if (gallery.value[eventId]) {
      gallery.value[eventId] = gallery.value[eventId].filter((i) => i.id !== itemId)
    }
  }

  // ========================
  // Custom Content
  // ========================

  function getCustomContent(eventId: string): CustomContent | null {
    return customContent.value[eventId] ?? null
  }

  async function fetchCustomContent(eventId: string) {
    const key = subKey('customContent', eventId)
    subLoading.value[key] = true
    subError.value[key] = ''
    try {
      customContent.value[eventId] = await adminGetCustomContent(eventId)
    } catch (e: unknown) {
      subError.value[key] = errorMsg(e, 'Failed to load custom content')
    } finally {
      subLoading.value[key] = false
    }
  }

  async function saveCustomContent(eventId: string, data: CustomContent) {
    const saved = await adminUpdateCustomContent(eventId, data)
    customContent.value[eventId] = saved
    return saved
  }

  // ========================
  // Enabled Components
  // ========================

  function getEnabledComponents(eventId: string): ComponentConfig[] {
    return enabledComponents.value[eventId] ?? []
  }

  async function fetchEnabledComponents(eventId: string) {
    const key = subKey('enabledComponents', eventId)
    subLoading.value[key] = true
    subError.value[key] = ''
    try {
      const data = await adminGetEnabledComponents(eventId)
      enabledComponents.value[eventId] = data.components?.length
        ? data.components.sort((a, b) => a.order - b.order)
        : []
    } catch (e: unknown) {
      subError.value[key] = errorMsg(e, 'Failed to load component config')
    } finally {
      subLoading.value[key] = false
    }
  }

  async function saveEnabledComponents(eventId: string, data: EnabledComponents) {
    const saved = await adminUpdateEnabledComponents(eventId, data)
    if (saved.components?.length) {
      enabledComponents.value[eventId] = saved.components.sort((a, b) => a.order - b.order)
    }
    return saved
  }

  // ========================
  // Reorder helper (shared by schedule, faqs, gifts, gallery)
  // ========================

  /** Swaps order of two items by calling update on both, then re-sorts. */
  async function swapOrder<T extends { id: string; orderIndex: number }>(
    items: T[],
    itemId: string,
    direction: 'up' | 'down',
    updateFn: (id: string, orderIndex: number) => Promise<T>,
  ) {
    const idx = items.findIndex((i) => i.id === itemId)
    if (idx === -1) return

    const swapIdx = direction === 'up' ? idx - 1 : idx + 1
    if (swapIdx < 0 || swapIdx >= items.length) return

    const current = items[idx]!
    const swap = items[swapIdx]!

    const [updatedCurrent, updatedSwap] = await Promise.all([
      updateFn(current.id, swap.orderIndex),
      updateFn(swap.id, current.orderIndex),
    ])

    items[idx] = updatedCurrent
    items[swapIdx] = updatedSwap
    items.sort((a, b) => a.orderIndex - b.orderIndex)
  }

  /** Clears all admin state. Called on logout. */
  function $reset() {
    events.value = []
    selectedEventId.value = null
    loading.value = false
    error.value = ''
    invites.value = {}
    schedule.value = {}
    faqs.value = {}
    gifts.value = {}
    gallery.value = {}
    customContent.value = {}
    enabledComponents.value = {}
    subLoading.value = {}
    subError.value = {}
  }

  return {
    // Events
    events,
    selectedEventId,
    selectedEvent,
    loading,
    error,
    fetchEvents,
    createEvent,
    updateEvent,
    deleteEvent,
    // Hosts
    addHost,
    deleteHost,
    // Invites
    getInvites,
    fetchInvites,
    createInvite,
    deleteInvite,
    addGuestToInvite,
    updateGuestInInvite,
    deleteGuestFromInvite,
    // Schedule
    getSchedule,
    fetchSchedule,
    createScheduleItem,
    updateScheduleItem,
    deleteScheduleItem,
    // FAQs
    getFAQs,
    fetchFAQs,
    createFAQItem,
    updateFAQItem,
    deleteFAQItem,
    // Gifts
    getGifts,
    fetchGifts,
    createGiftItem,
    updateGiftItem,
    deleteGiftItem,
    // Gallery
    getGallery,
    fetchGallery,
    createGalleryItem,
    updateGalleryItem,
    deleteGalleryItemFromStore,
    // Custom Content
    getCustomContent,
    fetchCustomContent,
    saveCustomContent,
    // Enabled Components
    getEnabledComponents,
    fetchEnabledComponents,
    saveEnabledComponents,
    // Helpers
    isSubLoading,
    getSubError,
    swapOrder,
    $reset,
  }
})
