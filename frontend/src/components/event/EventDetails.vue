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

      <!-- Dot divider -->
      <svg class="mx-auto mt-12 mb-12 w-[200px] md:w-[400px]" viewBox="0 0 400 6" fill="none">
        <circle v-for="i in 25" :key="i" :cx="4 + (i - 1) * 16.33" cy="3" r="2.5" fill="rgba(255,255,255,0.5)" />
      </svg>

      <!-- Description -->
      <p v-if="event.description" class="mb-2 text-gray-300 whitespace-pre-wrap max-w-2xl mx-auto text-xl md:text-3xl font-bold" style="font-family: 'Cormorant Garamond', 'Playfair Display', 'Georgia', serif;">
        {{ event.description }}
      </p>

      <!-- Date & Location -->
      <div class="space-y-3 text-gray-300 text-base md:text-xl">
        <div v-if="event.startsAt" class="flex items-center justify-center gap-2">
          <span>{{ formatDateFull(event.startsAt) }}</span>
        </div>
      </div>
    </div>
  </section>
</template>
