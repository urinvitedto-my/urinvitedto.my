<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import type { Event, GalleryItem } from '@/types'
import { formatDateFull } from '@/utils/date'

const props = defineProps<{
  event: Event
  gallery: GalleryItem[]
}>()

const currentIndex = ref(0)

const photos = computed(() =>
  props.gallery.filter((item) => item.mediaType === 'photo'),
)

/** Splits title into lines around "&" for stacked display. */
const titleParts = computed(() => {
  if (!props.event.title.includes('&')) return null
  return props.event.title.split('&').map((part) => part.trim())
})

let slideshowTimer: ReturnType<typeof setInterval> | null = null

onMounted(() => {
  if (photos.value.length > 1) {
    slideshowTimer = setInterval(() => {
      currentIndex.value = (currentIndex.value + 1) % photos.value.length
    }, 5000)
  }
})

onUnmounted(() => {
  if (slideshowTimer) clearInterval(slideshowTimer)
})
</script>

<template>
  <section class="event-details relative overflow-hidden py-10 md:py-24 lg:py-30 px-4">
    <!-- Slideshow background -->
    <template v-if="photos.length">
      <img
        v-for="(photo, idx) in photos"
        :key="photo.id"
        :src="photo.mediaUrl"
        :alt="photo.caption || ''"
        class="absolute inset-0 w-full h-full object-cover transition-opacity duration-1000"
        :class="idx === currentIndex ? 'opacity-100' : 'opacity-0'"
        aria-hidden="true"
      />
      <div class="absolute inset-0 bg-black/60" aria-hidden="true" />
    </template>

    <div class="relative max-w-3xl mx-auto text-center">
      <!-- Title -->
      <h1 class="text-8xl md:text-8xl lg:text-9xl font-normal text-white mb-6" style="font-family: 'Lavishly Yours', cursive; text-transform: none; letter-spacing: normal;">
        <template v-if="titleParts">
          <span v-for="(part, i) in titleParts" :key="i">
            {{ part }}<br v-if="i < titleParts.length - 1" /><span v-if="i < titleParts.length - 1" class="block text-3xl md:text-4xl">&amp;</span>
          </span>
        </template>
        <template v-else>{{ event.title }}</template>
      </h1>

      <!-- Description -->
      <p v-if="event.description" class="mb-8 text-gray-300 whitespace-pre-wrap max-w-2xl mx-auto text-xl font-bold" style="font-family: 'Cormorant Garamond', 'Playfair Display', 'Georgia', serif;">
        {{ event.description }}
      </p>

      <!-- Date & Location -->
      <div class="space-y-3 text-gray-300">
        <div v-if="event.startsAt" class="flex items-center justify-center gap-2">
          <svg class="w-5 h-5 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          <span>{{ formatDateFull(event.startsAt) }}</span>
        </div>
        <div v-if="event.location" class="flex items-center justify-center gap-2">
          <svg class="w-5 h-5 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
          <span>{{ event.location }}</span>
        </div>
      </div>
    </div>
  </section>
</template>
