<script setup lang="ts">
import { ref } from 'vue'
import type { GalleryItem } from '@/types'

defineProps<{
  items: GalleryItem[]
}>()

const lightboxOpen = ref(false)
const lightboxIndex = ref(0)

/**
 * Returns the gallery item at the current lightbox index.
 */
function currentItem(items: GalleryItem[]): GalleryItem {
  return items[lightboxIndex.value]!
}

/**
 * Opens lightbox at given index.
 */
function openLightbox(index: number) {
  lightboxIndex.value = index
  lightboxOpen.value = true
}

/**
 * Closes lightbox.
 */
function closeLightbox() {
  lightboxOpen.value = false
}

/**
 * Navigates lightbox.
 */
function navigate(direction: 1 | -1, total: number) {
  lightboxIndex.value = (lightboxIndex.value + direction + total) % total
}
</script>

<template>
  <section class="event-gallery bg-white py-12 px-4">
    <div class="max-w-5xl mx-auto">
      <h2 class="text-2xl font-bold text-[#14213d] text-center mb-8">Gallery</h2>

      <!-- Grid -->
      <div class="grid grid-cols-2 md:grid-cols-3 gap-4">
        <div
          v-for="(item, index) in items"
          :key="item.id"
          class="aspect-square overflow-hidden rounded-lg cursor-pointer group"
          @click="openLightbox(index)"
        >
          <img
            v-if="item.mediaType === 'photo'"
            :src="item.mediaUrl"
            :alt="item.caption || 'Gallery image'"
            class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
            loading="lazy"
          />
          <div
            v-else
            class="w-full h-full bg-[#14213d] flex items-center justify-center group-hover:bg-[#1a2a4d] transition-colors"
          >
            <svg class="w-12 h-12 text-white" fill="currentColor" viewBox="0 0 24 24">
              <path d="M8 5v14l11-7z" />
            </svg>
          </div>
        </div>
      </div>
    </div>

    <!-- Lightbox -->
    <div
      v-if="lightboxOpen"
      class="fixed inset-0 bg-black/90 z-50 flex items-center justify-center"
      @click.self="closeLightbox"
    >
      <!-- Close button -->
      <button
        @click="closeLightbox"
        class="absolute top-4 right-4 text-white hover:text-[#fca311] transition-colors"
      >
        <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>

      <!-- Navigation -->
      <button
        v-if="items.length > 1"
        @click="navigate(-1, items.length)"
        class="absolute left-4 text-white hover:text-[#fca311] transition-colors"
      >
        <svg class="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </button>
      <button
        v-if="items.length > 1"
        @click="navigate(1, items.length)"
        class="absolute right-4 text-white hover:text-[#fca311] transition-colors"
      >
        <svg class="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </button>

      <!-- Content -->
      <div class="max-w-4xl max-h-[80vh] mx-4">
        <img
          v-if="currentItem(items).mediaType === 'photo'"
          :src="currentItem(items).mediaUrl"
          :alt="currentItem(items).caption || 'Gallery image'"
          class="max-w-full max-h-[80vh] object-contain"
        />
        <video
          v-else
          :src="currentItem(items).mediaUrl"
          controls
          class="max-w-full max-h-[80vh]"
        ></video>
        <p v-if="currentItem(items).caption" class="text-white text-center mt-4">
          {{ currentItem(items).caption }}
        </p>
      </div>
    </div>
  </section>
</template>

<style scoped></style>
