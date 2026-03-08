<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useHostStore } from '@/stores/host'
import { formatDate, formatTimeOnly } from '@/utils/date'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import type { HostGuest, AdminInvite } from '@/types'

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
  showAllGuests,
  filteredGuests,
} = storeToRefs(hostStore)

const messageModalGuest = ref<HostGuest | null>(null)
const showInvitesModal = ref(false)
const copiedInviteId = ref<string | null>(null)

/** Opens the message modal for a guest (only if they have a message). */
function openMessage(guest: HostGuest) {
  if (guest.rsvpMessage) {
    messageModalGuest.value = guest
  }
}

function closeMessage() {
  messageModalGuest.value = null
}

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

/** Builds the event landing URL for the selected event. */
const eventUrl = computed(() => {
  if (!selectedEvent.value) return ''
  const origin = window.location.origin
  return `${origin}/${selectedEvent.value.type}/${selectedEvent.value.slug}`
})

/** Formats host names for the invite message (e.g. "Jester & Sione's"). */
function formatHostNames(): string {
  const names = eventHosts.value.map((h) => h.displayName)
  if (names.length === 0) return ''
  if (names.length === 1) return `${names[0]}'s`
  const last = names.pop()
  return `${names.join(', ')} & ${last}'s`
}

/** Generates the invite message text for a given invite. */
function buildInviteMessage(invite: AdminInvite): string {
  const eventType = selectedEvent.value?.type ?? 'event'
  const isPrivate = selectedEvent.value && !selectedEvent.value.isPublic
  const hostLabel = formatHostNames()
  const label = invite.label || 'there'

  const lines = [
    `Hi ${label}!`,
    '',
    `We're excited to let you know that you are invited to`,
    `${hostLabel} ${eventType}!`,
    '',
    `We would love for you to be part of this special day.`,
    '',
    `View your invitation here:`,
    eventUrl.value,
  ]

  if (isPrivate) {
    lines.push('', `Your invite code: ${invite.inviteCode}`)
  }

  return lines.join('\n')
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
  <div class="host-dashboard-view min-h-screen pt-24 pb-8 px-4">
    <div class="max-w-6xl mx-auto">
      <div class="mb-8">
        <h1 class="text-2xl font-bold text-primary">Host Dashboard</h1>
      </div>

      <div
        v-if="eventsLoading && !selectedEvent"
        class="flex items-center justify-center py-20"
      >
        <LoadingSpinner />
      </div>

      <div v-else-if="error" class="text-center py-20">
        <p class="text-red-600 mb-4">{{ error }}</p>
        <button @click="hostStore.fetchEvents()" class="text-accent hover:underline">
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
              </div>
            </li>
          </ul>
        </div>

        <div class="md:col-span-2 bg-white rounded-lg shadow-sm p-6">
          <div v-if="!selectedEvent" class="text-gray-500 text-center py-20">
            Select an event to view guests
          </div>
          <template v-else>
            <div class="flex flex-wrap items-center justify-between gap-3 mb-6">
              <h2 class="text-lg font-semibold text-primary">
                {{ selectedEvent.title }} - Guests
              </h2>
              <div class="flex items-center gap-4">
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
                <label class="flex items-center gap-2 text-sm">
                  <input v-model="showAllGuests" type="checkbox" class="rounded" />
                  Show all guests
                </label>
              </div>
            </div>

            <div v-if="guestsLoading" class="flex items-center justify-center py-12">
              <LoadingSpinner size="md" />
            </div>

            <div
              v-else-if="filteredGuests.length === 0"
              class="text-gray-500 text-center py-12"
            >
              {{ showAllGuests ? 'No guests yet' : 'No confirmed guests yet' }}
            </div>

            <div v-else class="overflow-x-auto">
              <table class="w-full">
                <thead>
                  <tr class="border-b border-muted">
                    <th class="text-left py-3 px-4 font-medium text-gray-600">Name</th>
                    <th class="text-left py-3 px-4 font-medium text-gray-600">
                      RSVP Date
                    </th>
                    <th></th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="guest in filteredGuests"
                    :key="guest.id"
                    class="border-b border-surface"
                  >
                    <td class="py-3 px-4">{{ guest.displayName }}</td>
                    <td class="py-3 px-4 text-gray-600 text-sm">
                      <div>{{ formatDate(guest.rsvpAt) }}</div>
                      <div class="text-gray-400">
                        {{ formatTimeOnly(guest.rsvpAt ?? '') }}
                      </div>
                    </td>
                    <td class="py-3 px-4 text-right">
                      <button
                        @click="openMessage(guest)"
                        :class="[
                          'shrink-0',
                          guest.rsvpMessage ? 'cursor-pointer' : 'cursor-default',
                        ]"
                        :title="guest.rsvpMessage ? 'View message' : undefined"
                        type="button"
                      >
                        <svg
                          v-if="guest.rsvpStatus === 'yes'"
                          class="w-5 h-5 text-green-600"
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
                          v-else-if="guest.rsvpStatus === 'no'"
                          class="w-5 h-5 text-red-600"
                          fill="none"
                          stroke="currentColor"
                          stroke-width="2.5"
                          viewBox="0 0 24 24"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            d="M6 18L18 6M6 6l12 12"
                          />
                        </svg>
                        <span
                          v-else
                          class="inline-block w-5 h-5 rounded-full border-2 border-gray-300"
                        />
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- Guest Message Modal -->
            <Teleport to="body">
              <div
                v-if="messageModalGuest"
                class="fixed inset-0 z-50 flex items-center justify-center bg-black/40"
                @click.self="closeMessage"
              >
                <div class="bg-white rounded-lg shadow-xl max-w-md w-full mx-4 p-6">
                  <div class="flex items-center justify-between mb-4">
                    <h3 class="text-lg font-semibold text-primary">
                      {{ messageModalGuest.displayName }}
                    </h3>
                    <button
                      @click="closeMessage"
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
                  <p class="text-gray-700 whitespace-pre-wrap">
                    {{ messageModalGuest.rsvpMessage }}
                  </p>
                </div>
              </div>
            </Teleport>

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
