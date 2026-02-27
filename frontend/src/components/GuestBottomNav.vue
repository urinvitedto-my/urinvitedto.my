<script setup lang="ts">
const navItems = [
  { id: 'section-top', label: 'Top' },
  { id: 'section-countdown', label: 'Countdown' },
  { id: 'section-gallery', label: 'Gallery' },
  { id: 'section-venue', label: 'Venue' },
  { id: 'section-rsvp', label: 'RSVP' },
]

/**
 * Smooth-scrolls to a section by its element id, or to the top if not found.
 */
function scrollToSection(sectionId: string) {
  if (sectionId === 'section-top') {
    window.scrollTo({ top: 0, behavior: 'smooth' })
    return
  }

  const el = document.getElementById(sectionId)
  if (el) {
    el.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }
}
</script>

<template>
  <nav class="fixed bottom-0 left-0 right-0 z-50 guest-bottom-nav pb-[env(safe-area-inset-bottom)]">
    <div class="flex items-center justify-evenly max-w-lg mx-auto pt-20 pb-4">
        <!-- 1: Envelope / Invitation (top of page) -->
        <button
          v-for="item in navItems"
          :key="item.id"
          @click="scrollToSection(item.id)"
          :aria-label="item.label"
          class="flex items-center justify-center w-16 h-16 rounded-full active:bg-gray-100 transition-colors cursor-pointer"
        >
          <!-- Envelope -->
          <svg v-if="item.id === 'section-top'" class="w-11 h-11" viewBox="0 0 24 24" fill="#14213d">
            <path d="M20 4H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm-.4 4.25l-6.54 4.09c-.65.41-1.47.41-2.12 0L4.4 8.25a.85.85 0 11.9-1.44L12 11l5.7-4.19a.85.85 0 01.9 1.44z" />
          </svg>

          <!-- Calendar with checkmark -->
          <svg v-else-if="item.id === 'section-countdown'" class="w-11 h-11" viewBox="0 0 24 24" fill="#14213d">
            <path d="M19 3h-1V2a1 1 0 00-2 0v1H8V2a1 1 0 00-2 0v1H5C3.9 3 3 3.9 3 5v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm0 16H5V8h14v11z" />
            <path d="M10.3 16.7l-2.5-2.5a1 1 0 011.4-1.4l1.1 1.1 3.3-3.3a1 1 0 011.4 1.4l-4.7 4.7z" />
          </svg>

          <!-- Heart -->
          <svg v-else-if="item.id === 'section-gallery'" class="w-11 h-11" viewBox="0 0 24 24" fill="#14213d">
            <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z" />
          </svg>

          <!-- Navigation / send arrow -->
          <svg v-else-if="item.id === 'section-venue'" class="w-11 h-11" viewBox="0 0 24 24" fill="#14213d">
            <path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z" />
          </svg>

          <!-- Contact / ID card -->
          <svg v-else-if="item.id === 'section-rsvp'" class="w-11 h-11" viewBox="0 0 24 24" fill="#14213d">
            <path d="M2 3h20c1.1 0 2 .9 2 2v14c0 1.1-.9 2-2 2H2c-1.1 0-2-.9-2-2V5c0-1.1.9-2 2-2zm7 6a2.5 2.5 0 100 5 2.5 2.5 0 000-5zm-5 10h10v-1.25c0-1.66-3.33-2.5-5-2.5s-5 .84-5 2.5V19zm13-8h4v1.5h-4V11zm0 3h4v1.5h-4V14zm0-6h4v1.5h-4V8z" />
          </svg>
        </button>
    </div>
  </nav>
</template>

<style scoped>
.guest-bottom-nav {
  background: linear-gradient(to bottom, transparent 0%, rgba(255, 255, 255, 0.85) 40%, white 100%);
}
</style>
