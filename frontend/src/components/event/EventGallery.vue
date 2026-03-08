<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import type { GalleryItem } from '@/types'

const props = defineProps<{
  items: GalleryItem[]
}>()

// Lightbox
const lightboxOpen = ref(false)
const lightboxIndex = ref(0)

// Triple the array: [copy1, copy2, copy3]. We live in copy2.
const cloned = computed(() => [...props.items, ...props.items, ...props.items])

// Absolute position in the cloned array. Starts at the beginning of copy2.
const trackIndex = ref(props.items.length)
const animated = ref(true)

let autoTimer: ReturnType<typeof setInterval> | null = null
let resetTimer: ReturnType<typeof setTimeout> | null = null

const isHovered = ref(false)

const MOBILE_BREAKPOINT = 768
const visibleCount = ref(typeof window !== 'undefined' && window.innerWidth >= MOBILE_BREAKPOINT ? 3 : 2)

/** Updates visible item count based on screen width. */
function updateVisibleCount() {
  visibleCount.value = window.innerWidth >= MOBILE_BREAKPOINT ? 3 : 2
}

// Pointer drag state (covers mouse, touch, trackpad)
const isDragging = ref(false)
const dragDelta = ref(0)
const clickBlocked = ref(false)
let startX = 0
let startY = 0
let isHorizontalDrag = false

/**
 * Moves the track by one step, then schedules a silent reset if needed.
 */
function step(dir: 1 | -1) {
  animated.value = true
  trackIndex.value += dir
  scheduleReset()
}

/**
 * After the 700ms CSS transition completes, silently teleport back
 * to the equivalent position in copy2 if we drifted into copy1 or copy3.
 */
function scheduleReset() {
  if (resetTimer) clearTimeout(resetTimer)
  resetTimer = setTimeout(async () => {
    const len = props.items.length
    let jumped = false
    if (trackIndex.value >= len * 2) {
      animated.value = false
      trackIndex.value -= len
      jumped = true
    } else if (trackIndex.value < len) {
      animated.value = false
      trackIndex.value += len
      jumped = true
    }
    if (!jumped) return
    await nextTick()
    requestAnimationFrame(() => requestAnimationFrame(() => {
      animated.value = true
    }))
  }, 750)
}

/**
 * Unified pointer handler - covers mouse, touch, and trackpad.
 */
function onPointerDown(e: PointerEvent) {
  isDragging.value = true
  clickBlocked.value = false
  isHorizontalDrag = false
  startX = e.clientX
  startY = e.clientY
  dragDelta.value = 0
  ;(e.currentTarget as HTMLElement).setPointerCapture(e.pointerId)
}

function onPointerMove(e: PointerEvent) {
  if (!isDragging.value) return
  const dx = e.clientX - startX
  const dy = e.clientY - startY

  // On first significant movement, decide if horizontal or vertical
  if (!isHorizontalDrag && (Math.abs(dx) > 5 || Math.abs(dy) > 5)) {
    if (Math.abs(dy) > Math.abs(dx)) {
      // Vertical scroll — stop tracking so the page can scroll
      isDragging.value = false
      dragDelta.value = 0
      return
    }
    isHorizontalDrag = true
  }

  if (isHorizontalDrag) {
    e.preventDefault()
    dragDelta.value = dx
    if (Math.abs(dx) > 5) clickBlocked.value = true
  }
}

function onPointerUp(e: PointerEvent) {
  if (!isDragging.value) return
  const dx = e.clientX - startX
  isDragging.value = false
  animated.value = true
  dragDelta.value = 0

  if (isHorizontalDrag && Math.abs(dx) > 40) {
    step(dx < 0 ? 1 : -1)
  }
  // If not enough drag, dragDelta animates back to 0 (snap back)
}

/**
 * Opens lightbox, resolving the real item index from the cloned array.
 * Blocked if the user just finished a drag.
 */
function handleItemClick(clonedIdx: number) {
  if (clickBlocked.value) {
    clickBlocked.value = false
    return
  }
  lightboxIndex.value = clonedIdx % props.items.length
  lightboxOpen.value = true
}

/** Closes lightbox. */
function closeLightbox() {
  lightboxOpen.value = false
}

/** Navigates lightbox. */
function navigateLightbox(dir: 1 | -1) {
  const total = props.items.length
  lightboxIndex.value = (lightboxIndex.value + dir + total) % total
}

/** Returns the gallery item at the current lightbox index. */
function currentItem(): GalleryItem {
  return props.items[lightboxIndex.value]!
}

onMounted(() => {
  trackIndex.value = props.items.length
  updateVisibleCount()
  window.addEventListener('resize', updateVisibleCount)
  autoTimer = setInterval(() => { if (!isHovered.value) step(1) }, 4000)
})

onUnmounted(() => {
  window.removeEventListener('resize', updateVisibleCount)
  if (autoTimer) clearInterval(autoTimer)
  if (resetTimer) clearTimeout(resetTimer)
})
</script>

<template>
  <section class="event-gallery py-16 px-4">
    <div class="max-w-5xl mx-auto">
      <h2 class="text-4xl font-bold text-primary-dark text-center mb-8 font-kaushan">Our Prenup</h2>

      <!-- Carousel -->
      <div
        class="relative"
        @mouseenter="isHovered = true"
        @mouseleave="isHovered = false"
      >
        <!-- Left arrow -->
        <button
          v-if="items.length > 1"
          @click="step(-1)"
          aria-label="Previous image"
          class="absolute left-0 top-1/2 -translate-y-1/2 -translate-x-3 z-10 w-9 h-9 flex items-center justify-center rounded-full bg-white/80 backdrop-blur shadow border border-muted/50 text-heading hover:text-accent transition-colors"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>

        <!-- Track -->
        <div
          class="overflow-hidden cursor-grab active:cursor-grabbing select-none"
          style="touch-action: pan-y;"
          @pointerdown="onPointerDown"
          @pointermove="onPointerMove"
          @pointerup="onPointerUp"
          @pointercancel="onPointerUp"
        >
          <div
            class="flex"
            :style="{
              width: `${cloned.length / visibleCount * 100}%`,
              transform: `translateX(calc(-${trackIndex / cloned.length * 100}% + ${dragDelta}px))`,
              transition: isDragging ? 'none' : (animated ? 'transform 0.7s ease-in-out' : 'none'),
            }"
          >
            <div
              v-for="(item, idx) in cloned"
              :key="idx"
              class="px-2 group"
              :style="{ width: `${100 / cloned.length}%` }"
              @click="handleItemClick(idx)"
            >
              <div class="aspect-4/5 overflow-hidden rounded-xl shadow-sm border border-muted/50 bg-white/80 backdrop-blur">
                <img
                  v-if="item.mediaType === 'photo'"
                  :src="item.mediaUrl"
                  :alt="item.caption || 'Gallery image'"
                  class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500 pointer-events-none"
                  loading="lazy"
                  draggable="false"
                />
                <div
                  v-else
                  class="w-full h-full bg-primary flex items-center justify-center group-hover:bg-primary-dark transition-colors"
                >
                  <svg class="w-12 h-12 text-accent pointer-events-none" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M8 5v14l11-7z" />
                  </svg>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Right arrow -->
        <button
          v-if="items.length > 1"
          @click="step(1)"
          aria-label="Next image"
          class="absolute right-0 top-1/2 -translate-y-1/2 translate-x-3 z-10 w-9 h-9 flex items-center justify-center rounded-full bg-white/80 backdrop-blur shadow border border-muted/50 text-heading hover:text-accent transition-colors"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Lightbox -->
    <div
      v-if="lightboxOpen"
      class="fixed inset-0 bg-black/90 z-50 flex items-center justify-center"
      @click.self="closeLightbox"
    >
      <button
        @click="closeLightbox"
        aria-label="Close lightbox"
        class="absolute top-4 right-4 text-white hover:text-accent transition-colors"
      >
        <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>

      <button
        v-if="items.length > 1"
        @click="navigateLightbox(-1)"
        aria-label="Previous image"
        class="absolute left-4 text-white hover:text-accent transition-colors"
      >
        <svg class="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </button>
      <button
        v-if="items.length > 1"
        @click="navigateLightbox(1)"
        aria-label="Next image"
        class="absolute right-4 text-white hover:text-accent transition-colors"
      >
        <svg class="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </button>

      <div class="max-w-4xl max-h-[80vh] mx-4">
        <img
          v-if="currentItem().mediaType === 'photo'"
          :src="currentItem().mediaUrl"
          :alt="currentItem().caption || 'Gallery image'"
          class="max-w-full max-h-[80vh] object-contain"
        />
        <video
          v-else
          :src="currentItem().mediaUrl"
          controls
          class="max-w-full max-h-[80vh]"
        ></video>
        <p v-if="currentItem().caption" class="text-white text-center mt-4">
          {{ currentItem().caption }}
        </p>
      </div>
    </div>
  </section>
</template>
