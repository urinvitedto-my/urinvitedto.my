<script setup lang="ts">
import type { LocationDetails } from '@/types'

defineProps<{
  locationDetails: LocationDetails
  address?: string
}>()
</script>

<template>
  <section class="event-map bg-white py-12 px-4">
    <div class="max-w-3xl mx-auto">
      <h2 class="text-2xl font-bold text-[#14213d] text-center mb-6">Venue</h2>

      <!-- Map embed -->
      <div v-if="locationDetails.mapEmbedUrl" class="aspect-video rounded-lg overflow-hidden mb-6">
        <iframe
          :src="locationDetails.mapEmbedUrl"
          class="w-full h-full border-0"
          loading="lazy"
          referrerpolicy="no-referrer-when-downgrade"
          allowfullscreen
        ></iframe>
      </div>

      <!-- Address -->
      <div v-if="address" class="text-center mb-6">
        <div class="flex items-center justify-center gap-2 text-gray-700">
          <svg class="w-5 h-5 text-[#fca311]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
          </svg>
          <span>{{ address }}</span>
        </div>
      </div>

      <!-- Additional info -->
      <div class="space-y-4 text-gray-600">
        <div v-if="locationDetails.parkingInfo" class="flex gap-3">
          <svg class="w-5 h-5 text-[#14213d] shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
          </svg>
          <p>{{ locationDetails.parkingInfo }}</p>
        </div>
        <div v-if="locationDetails.accessibilityNotes" class="flex gap-3">
          <svg class="w-5 h-5 text-[#14213d] shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <p>{{ locationDetails.accessibilityNotes }}</p>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped></style>
