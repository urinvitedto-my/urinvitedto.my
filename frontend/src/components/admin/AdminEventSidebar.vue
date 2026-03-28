<script setup lang="ts">
import { ref, computed } from 'vue'
import { formatDate } from '@/utils/date'
import type { AdminEvent, EventType } from '@/types'

const props = defineProps<{
  events: AdminEvent[]
  selectedEventId: string | null
}>()

const emit = defineEmits<{ select: [eventId: string] }>()

const search = ref('')
const typeFilter = ref<EventType | 'all'>('all')
const visibilityFilter = ref<'all' | 'public' | 'private'>('all')
const sortBy = ref<'date' | 'title' | 'created'>('created')
const sortDir = ref<'asc' | 'desc'>('desc')

/** Filters and sorts events based on current controls. */
const filteredEvents = computed(() => {
  let list = props.events

  const q = search.value.toLowerCase().trim()
  if (q) {
    list = list.filter(
      (e) => e.title.toLowerCase().includes(q) || e.slug.toLowerCase().includes(q),
    )
  }

  if (typeFilter.value !== 'all') {
    list = list.filter((e) => e.type === typeFilter.value)
  }

  if (visibilityFilter.value === 'public') {
    list = list.filter((e) => e.isPublic)
  } else if (visibilityFilter.value === 'private') {
    list = list.filter((e) => !e.isPublic)
  }

  const sorted = [...list].sort((a, b) => {
    let cmp = 0
    switch (sortBy.value) {
      case 'date':
        cmp = (a.startsAt ?? '').localeCompare(b.startsAt ?? '')
        break
      case 'title':
        cmp = a.title.localeCompare(b.title)
        break
      case 'created':
        cmp = a.createdAt.localeCompare(b.createdAt)
        break
    }
    return sortDir.value === 'asc' ? cmp : -cmp
  })

  return sorted
})

function toggleSortDir() {
  sortDir.value = sortDir.value === 'asc' ? 'desc' : 'asc'
}

function pendingCount(event: AdminEvent): number {
  return event.guestCount - event.rsvpYes - event.rsvpNo
}
</script>

<template>
  <div class="admin-sidebar flex flex-col h-full">
    <!-- Controls (pinned, never scrolls) -->
    <div class="shrink-0 flex flex-col gap-3 pb-3">
      <!-- Search -->
      <input
        v-model="search"
        type="text"
        placeholder="Search by title or slug..."
        class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-accent focus:outline-none"
      />

      <!-- Filters + Sort -->
      <div class="flex flex-wrap items-center gap-2">
        <select
          v-model="typeFilter"
          class="px-2 py-1.5 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-accent focus:outline-none"
        >
          <option value="all">All types</option>
          <option value="wedding">Wedding</option>
          <option value="birthday">Birthday</option>
          <option value="party">Party</option>
        </select>

        <select
          v-model="visibilityFilter"
          class="px-2 py-1.5 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-accent focus:outline-none"
        >
          <option value="all">All</option>
          <option value="public">Public</option>
          <option value="private">Private</option>
        </select>

        <select
          v-model="sortBy"
          class="px-2 py-1.5 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-accent focus:outline-none"
        >
          <option value="created">Created</option>
          <option value="date">Event date</option>
          <option value="title">Title</option>
        </select>

        <button
          @click="toggleSortDir"
          class="px-2 py-1.5 border border-gray-300 rounded-lg text-sm hover:bg-gray-50 transition-colors"
          :title="sortDir === 'asc' ? 'Ascending' : 'Descending'"
        >
          {{ sortDir === 'asc' ? '↑' : '↓' }}
        </button>
      </div>
    </div>

    <!-- Scrollable Event Cards -->
    <div class="flex-1 min-h-0 overflow-y-auto max-h-[calc(100vh-20rem)] lg:max-h-[calc(100vh-20rem)]">
      <div class="flex flex-col gap-3">
        <p v-if="filteredEvents.length === 0" class="text-sm text-gray-400 text-center py-4">
          No events match your filters
        </p>

        <div
          v-for="event in filteredEvents"
          :key="event.id"
          @click="emit('select', event.id)"
          class="rounded-lg border-2 p-3 cursor-pointer transition-all hover:shadow-md"
          :class="selectedEventId === event.id
            ? 'border-accent bg-amber-50 shadow-sm'
            : 'border-gray-200 bg-white hover:border-gray-300'"
        >
          <div class="flex items-center gap-2 mb-1">
            <span class="inline-block px-1.5 py-0.5 bg-accent text-black text-xs font-medium rounded capitalize">
              {{ event.type }}
            </span>
            <span class="text-xs text-gray-400">{{ event.isPublic ? 'Public' : 'Private' }}</span>
          </div>

          <h4 class="text-sm font-semibold text-primary truncate">{{ event.title }}</h4>

          <p class="text-xs text-gray-500 truncate mt-0.5">
            {{ formatDate(event.startsAt) }} · {{ event.location || 'No location' }}
          </p>

          <div class="flex items-center gap-2 mt-1.5 text-xs text-gray-500">
            <span>{{ event.guestCount }} guest{{ event.guestCount !== 1 ? 's' : '' }}</span>
            <span class="text-gray-300">|</span>
            <span class="text-green-600">{{ event.rsvpYes }} yes</span>
            <span class="text-red-500">{{ event.rsvpNo }} no</span>
            <span class="text-gray-400">{{ pendingCount(event) }} pending</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
