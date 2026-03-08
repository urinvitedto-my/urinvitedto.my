<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useHostStore } from '@/stores/host'
import { formatDate, formatTimeOnly } from '@/utils/date'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import type { HostGuest } from '@/types'

const hostStore = useHostStore()
const {
  events,
  selectedEvent,
  eventsLoading,
  guestsLoading,
  error,
  showAllGuests,
  filteredGuests,
} = storeToRefs(hostStore)

const messageModalGuest = ref<HostGuest | null>(null)

/** Opens the message modal for a guest (only if they have a message). */
function openMessage(guest: HostGuest) {
  if (guest.rsvpMessage) {
    messageModalGuest.value = guest
  }
}

function closeMessage() {
  messageModalGuest.value = null
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

      <div v-if="eventsLoading && !selectedEvent" class="flex items-center justify-center py-20">
        <LoadingSpinner />
      </div>

      <div v-else-if="error" class="text-center py-20">
        <p class="text-red-600 mb-4">{{ error }}</p>
        <button @click="hostStore.fetchEvents()" class="text-accent hover:underline">
          Try again
        </button>
      </div>

      <div v-else class="grid md:grid-cols-3 gap-8">
        <div class="bg-white rounded-lg shadow-sm p-6 md:sticky md:top-24 md:self-start">
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
                <div class="text-sm opacity-75">{{ formatDate(event.startsAt, true) }}</div>
              </div>
            </li>
          </ul>
        </div>

        <div class="md:col-span-2 bg-white rounded-lg shadow-sm p-6">
          <div v-if="!selectedEvent" class="text-gray-500 text-center py-20">
            Select an event to view guests
          </div>
          <template v-else>
            <div class="flex items-center justify-between mb-6">
              <h2 class="text-lg font-semibold text-primary">
                {{ selectedEvent.title }} - Guests
              </h2>
              <label class="flex items-center gap-2 text-sm">
                <input v-model="showAllGuests" type="checkbox" class="rounded" />
                Show all guests
              </label>
            </div>

            <div v-if="guestsLoading" class="flex items-center justify-center py-12">
              <LoadingSpinner size="md" />
            </div>

            <div v-else-if="filteredGuests.length === 0" class="text-gray-500 text-center py-12">
              {{ showAllGuests ? 'No guests yet' : 'No confirmed guests yet' }}
            </div>

            <div v-else class="overflow-x-auto">
              <table class="w-full">
                <thead>
                  <tr class="border-b border-muted">
                    <th class="text-left py-3 px-4 font-medium text-gray-600">Name</th>
                    <th class="text-left py-3 px-4 font-medium text-gray-600">RSVP Date</th>
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
                      <div class="text-gray-400">{{ formatTimeOnly(guest.rsvpAt ?? '') }}</div>
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
                          <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                        </svg>
                        <svg
                          v-else-if="guest.rsvpStatus === 'no'"
                          class="w-5 h-5 text-red-600"
                          fill="none"
                          stroke="currentColor"
                          stroke-width="2.5"
                          viewBox="0 0 24 24"
                        >
                          <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
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

            <!-- Message Modal -->
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
                      <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                      </svg>
                    </button>
                  </div>
                  <p class="text-gray-700 whitespace-pre-wrap">{{ messageModalGuest.rsvpMessage }}</p>
                </div>
              </div>
            </Teleport>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>
