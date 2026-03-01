<script setup lang="ts">
import type { Event, Host } from '@/types'

defineProps<{
  event: Event
  hosts: Host[]
}>()

/**
 * Formats date for display.
 */
function formatDate(dateStr?: string): string {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('en-US', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: 'numeric',
    minute: '2-digit',
  })
}
</script>

<template>
  <section class="event-details bg-white py-12 px-4">
    <div class="max-w-3xl mx-auto text-center">
      <!-- Title -->
      <h1 class="text-3xl md:text-5xl font-bold text-[#14213d] mb-6">
        {{ event.title }}
      </h1>

      <!-- Date & Location -->
      <div class="space-y-3 text-gray-600">
        <div v-if="event.startsAt" class="flex items-center justify-center gap-2">
          <svg class="w-5 h-5 text-[#fca311]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          <span>{{ formatDate(event.startsAt) }}</span>
        </div>
        <div v-if="event.location" class="flex items-center justify-center gap-2">
          <svg class="w-5 h-5 text-[#fca311]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
          <span>{{ event.location }}</span>
        </div>
      </div>

      <!-- Description -->
      <p v-if="event.description" class="mt-8 text-gray-700 whitespace-pre-wrap max-w-2xl mx-auto">
        {{ event.description }}
      </p>
    </div>
  </section>
</template>

<style scoped></style>
