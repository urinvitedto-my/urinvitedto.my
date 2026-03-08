<script setup lang="ts">
import type { LocationDetails } from '@/types'

defineProps<{
  locationDetails: LocationDetails
}>()
</script>

<template>
  <section class="event-map py-16 px-4">
    <div class="max-w-3xl mx-auto">
      <h2 class="text-4xl font-bold text-primary-dark text-center mb-6 font-kaushan">Map</h2>

      <!-- Map embed -->
      <div
        v-if="locationDetails.mapEmbedUrl"
        class="aspect-video rounded-lg overflow-hidden mb-6"
      >
        <iframe
          :src="locationDetails.mapEmbedUrl"
          title="Venue map"
          class="w-full h-full border-0"
          loading="lazy"
          referrerpolicy="no-referrer-when-downgrade"
          allowfullscreen
        ></iframe>
      </div>

      <!-- Additional info -->
      <div class="space-y-4 text-gray-600">
        <div v-if="locationDetails.parkingInfo" class="flex gap-3">
          <svg
            class="w-5 h-5 text-primary shrink-0 mt-0.5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4"
            />
          </svg>
          <p>{{ locationDetails.parkingInfo }}</p>
        </div>
        <div v-if="locationDetails.accessibilityNotes" class="flex gap-3">
          <svg
            class="w-5 h-5 text-primary shrink-0 mt-0.5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <p>{{ locationDetails.accessibilityNotes }}</p>
        </div>
      </div>
    </div>
  </section>
</template>
