<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useHostStore } from '@/stores/host'
import { formatDate } from '@/utils/date'
import { useInviteMessage } from '@/composables/useInviteMessage'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import type { AdminInvite } from '@/types'

const hostStore = useHostStore()
const {
  events,
  selectedEvent,
  eventsLoading,
  guestsLoading,
  invitesLoading,
  invites,
  eventHosts,
  error,
  guestFilter,
  filteredGuests,
  totalCount,
  yesCount,
  noCount,
  pendingCount,
} = storeToRefs(hostStore)

type StatusFilter = 'all' | 'yes' | 'no' | 'pending'

const filterButtons: { key: StatusFilter; label: string }[] = [
  { key: 'all', label: 'All' },
  { key: 'yes', label: 'Confirmed' },
  { key: 'no', label: 'Declined' },
  { key: 'pending', label: 'Pending' },
]

/** Pending RSVPs for an event row (guests minus yes/no). */
function hostPendingCount(event: { guestCount: number; rsvpYes: number; rsvpNo: number }) {
  return event.guestCount - event.rsvpYes - event.rsvpNo
}

/** Returns the count for each filter button. */
function filterCount(key: StatusFilter): number {
  switch (key) {
    case 'all':
      return totalCount.value
    case 'yes':
      return yesCount.value
    case 'no':
      return noCount.value
    case 'pending':
      return pendingCount.value
  }
}

/** Returns the RSVP badge CSS classes for a status. */
function rsvpBadgeClass(status: string): string {
  switch (status) {
    case 'yes':
      return 'bg-green-100 text-green-700'
    case 'no':
      return 'bg-red-100 text-red-700'
    default:
      return 'bg-gray-100 text-gray-600'
  }
}

const showInvitesModal = ref(false)
const copiedInviteId = ref<string | null>(null)

const { buildInviteMessage } = useInviteMessage(selectedEvent, eventHosts)

/** Opens the invite messages modal, fetching invites if needed. */
async function openInvitesModal() {
  showInvitesModal.value = true
  if (selectedEvent.value) {
    await hostStore.fetchInvites(selectedEvent.value.id)
  }
}

function closeInvitesModal() {
  showInvitesModal.value = false
  copiedInviteId.value = null
}

/** Copies the invite message to clipboard. */
async function copyMessage(invite: AdminInvite) {
  try {
    await navigator.clipboard.writeText(buildInviteMessage(invite))
    copiedInviteId.value = invite.id
    setTimeout(() => {
      if (copiedInviteId.value === invite.id) copiedInviteId.value = null
    }, 2000)
  } catch {
    // fallback: silent fail
  }
}

onMounted(() => {
  hostStore.fetchEvents()
})
</script>

<template>
  <div class="host-dashboard-view min-h-screen pt-24 pb-8 px-4 md:mx-16">
    <div class="max-w-8xl mx-auto">
      <div class="mb-8">
        <h1 class="text-2xl font-bold text-primary">Host Dashboard</h1>
      </div>

      <div
        v-if="eventsLoading && !selectedEvent"
        class="flex items-center justify-center py-20"
      >
        <LoadingSpinner />
      </div>

      <div
        v-else-if="error"
        class="flex flex-col items-center justify-center text-center min-h-[60vh]"
      >
        <h2 class="text-2xl font-bold text-primary mb-2">Something went wrong</h2>
        <p class="text-base text-gray-500 mb-6 max-w-sm">{{ error }}</p>
        <button
          @click="hostStore.fetchEvents()"
          class="bg-accent text-black font-semibold px-8 py-3 rounded-lg text-base hover:bg-accent-dark transition-colors"
        >
          Try again
        </button>
      </div>

      <div v-else class="grid md:grid-cols-3 gap-8">
        <div
          class="bg-white rounded-lg shadow-sm p-6 md:sticky md:top-24 md:self-start"
        >
          <h2 class="text-lg font-semibold text-primary mb-4">Your Events</h2>
          <div v-if="events.length === 0" class="text-gray-500 text-center py-8">
            No events found
          </div>
          <ul v-else class="space-y-2 max-h-[60vh] overflow-y-auto">
            <li v-for="event in events" :key="event.id">
              <div
                role="button"
                tabindex="0"
                @click="hostStore.selectEvent(event)"
                @keydown.enter="hostStore.selectEvent(event)"
                :class="[
                  'w-full text-left px-4 py-3 rounded-lg transition-colors cursor-pointer',
                  selectedEvent?.id === event.id
                    ? 'bg-accent text-black'
                    : 'bg-surface hover:bg-muted',
                ]"
              >
                <div class="flex items-center justify-between">
                  <div class="font-medium">{{ event.title }}</div>
                  <RouterLink
                    :to="`/${event.type}/${event.slug}`"
                    class="text-xs underline opacity-75 hover:opacity-100"
                    @click.stop
                  >
                    View page
                  </RouterLink>
                </div>
                <div class="text-sm opacity-75">
                  {{ formatDate(event.startsAt, true) }}
                </div>
                <div
                  class="flex items-center gap-2 mt-1.5 text-xs opacity-90"
                  :class="selectedEvent?.id === event.id ? 'text-black/80' : 'text-gray-600'"
                >
                  <span
                    >{{ event.guestCount }} guest{{
                      event.guestCount !== 1 ? 's' : ''
                    }}</span
                  >
                  <span :class="selectedEvent?.id === event.id ? 'text-black/40' : 'text-gray-300'"
                    >|</span
                  >
                  <span class="text-green-600">{{ event.rsvpYes }} yes</span>
                  <span class="text-red-500">{{ event.rsvpNo }} no</span>
                  <span class="text-gray-500">{{ hostPendingCount(event) }} pending</span>
                </div>
              </div>
            </li>
          </ul>
        </div>

        <div class="md:col-span-2 bg-white rounded-lg shadow-sm p-6">
          <div v-if="!selectedEvent" class="text-gray-500 text-center py-20">
            Select an event to view guests
          </div>
          <template v-else>
            <div class="flex flex-wrap items-center justify-between gap-3 mb-4">
              <h2 class="text-lg font-semibold text-primary">
                {{ selectedEvent.title }} - Guests
              </h2>
              <button
                @click="openInvitesModal"
                class="inline-flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium rounded-lg bg-accent text-black hover:bg-accent/80 transition-colors"
                type="button"
              >
                <svg
                  class="w-4 h-4"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
                  />
                </svg>
                Invite Messages
              </button>
            </div>

            <div v-if="guestsLoading" class="flex items-center justify-center py-12">
              <LoadingSpinner size="md" />
            </div>

            <template v-else>
              <!-- Filter pills -->
              <div class="flex flex-wrap gap-2 mb-4">
                <button
                  v-for="btn in filterButtons"
                  :key="btn.key"
                  @click="guestFilter = btn.key"
                  class="px-3 py-1.5 rounded-full text-xs font-medium transition-colors"
                  :class="
                    guestFilter === btn.key
                      ? 'bg-primary text-white'
                      : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                  "
                >
                  {{ btn.label }} ({{ filterCount(btn.key) }})
                </button>
              </div>

              <!-- Empty states -->
              <p
                v-if="totalCount === 0"
                class="text-sm text-gray-400 py-4"
              >
                No guests yet.
              </p>

              <p
                v-else-if="filteredGuests.length === 0"
                class="text-sm text-gray-400 py-4"
              >
                No guests match this filter.
              </p>

              <!-- Guest cards -->
              <div v-else class="space-y-2">
                <div
                  v-for="guest in filteredGuests"
                  :key="guest.id"
                  class="bg-gray-50 rounded-lg px-4 py-3 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-1"
                >
                  <div class="min-w-0 flex-1">
                    <div class="flex items-center gap-2 flex-wrap">
                      <span class="text-sm font-semibold text-primary">{{
                        guest.displayName
                      }}</span>
                      <span
                        class="text-xs px-1.5 py-0.5 rounded capitalize"
                        :class="rsvpBadgeClass(guest.rsvpStatus)"
                      >
                        {{ guest.rsvpStatus }}
                      </span>
                    </div>
                    <p
                      v-if="guest.rsvpMessage"
                      class="text-xs text-gray-500 italic mt-0.5 wrap-break-word"
                    >
                      "{{ guest.rsvpMessage }}"
                    </p>
                  </div>

                  <div class="flex items-center gap-3 text-xs text-gray-400 shrink-0">
                    <span v-if="guest.rsvpAt">{{ formatDate(guest.rsvpAt, true) }}</span>
                  </div>
                </div>
              </div>
            </template>

            <!-- Invite Messages Modal -->
            <Teleport to="body">
              <div
                v-if="showInvitesModal"
                class="fixed inset-0 z-50 flex items-center justify-center bg-black/40"
                @click.self="closeInvitesModal"
              >
                <div
                  class="bg-white rounded-lg shadow-xl w-full max-w-2xl mx-4 max-h-[85vh] flex flex-col"
                >
                  <div
                    class="flex items-center justify-between p-6 pb-4 border-b border-muted"
                  >
                    <h3 class="text-lg font-semibold text-primary">Invite Messages</h3>
                    <button
                      @click="closeInvitesModal"
                      class="text-gray-400 hover:text-gray-600 transition-colors"
                      type="button"
                    >
                      <svg
                        class="w-5 h-5"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        viewBox="0 0 24 24"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          d="M6 18L18 6M6 6l12 12"
                        />
                      </svg>
                    </button>
                  </div>

                  <div class="flex-1 overflow-y-auto p-6 pt-4">
                    <div
                      v-if="invitesLoading"
                      class="flex items-center justify-center py-12"
                    >
                      <LoadingSpinner size="md" />
                    </div>

                    <div
                      v-else-if="invites.length === 0"
                      class="text-gray-500 text-center py-12"
                    >
                      No invites found for this event
                    </div>

                    <div v-else class="space-y-4">
                      <div
                        v-for="invite in invites"
                        :key="invite.id"
                        class="border border-muted rounded-lg p-4"
                      >
                        <div class="flex items-start justify-between gap-3 mb-2">
                          <div class="min-w-0">
                            <div class="font-medium text-primary truncate">
                              {{ invite.label || `Invite ${invite.inviteCode}` }}
                            </div>
                            <div class="text-xs text-gray-500 mt-0.5">
                              Code:
                              <span class="font-mono font-medium">{{
                                invite.inviteCode
                              }}</span>
                              <span v-if="invite.guests.length" class="ml-2">
                                &middot;
                                {{ invite.guests.map((g) => g.displayName).join(', ') }}
                              </span>
                            </div>
                          </div>
                          <button
                            @click="copyMessage(invite)"
                            class="shrink-0 inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium rounded-md transition-colors"
                            :class="
                              copiedInviteId === invite.id
                                ? 'bg-green-100 text-green-700'
                                : 'bg-surface text-gray-700 hover:bg-muted'
                            "
                            type="button"
                          >
                            <svg
                              v-if="copiedInviteId === invite.id"
                              class="w-3.5 h-3.5"
                              fill="none"
                              stroke="currentColor"
                              stroke-width="2.5"
                              viewBox="0 0 24 24"
                            >
                              <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                d="M5 13l4 4L19 7"
                              />
                            </svg>
                            <svg
                              v-else
                              class="w-3.5 h-3.5"
                              fill="none"
                              stroke="currentColor"
                              stroke-width="2"
                              viewBox="0 0 24 24"
                            >
                              <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
                              />
                            </svg>
                            {{ copiedInviteId === invite.id ? 'Copied!' : 'Copy' }}
                          </button>
                        </div>
                        <pre
                          class="text-sm text-gray-600 bg-surface rounded-md p-3 whitespace-pre-wrap font-sans leading-relaxed"
                          >{{ buildInviteMessage(invite) }}</pre
                        >
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </Teleport>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>
