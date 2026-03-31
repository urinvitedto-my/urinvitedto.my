<script setup lang="ts">
import { ref, computed } from 'vue'

const PREVIEW_LIMIT = 4

const props = defineProps<{
  guests: { displayName: string }[]
  count: number
}>()

const expanded = ref(false)

const hasMore = computed(() => props.guests.length > PREVIEW_LIMIT)

const visibleGuests = computed(() =>
  expanded.value ? props.guests : props.guests.slice(0, PREVIEW_LIMIT),
)
</script>

<template>
  <section class="confirmed-guests pb-16 px-4">
    <div class="max-w-3xl mx-auto bg-primary-dark/60 p-4 rounded-lg shadow-lg shadow-primary-dark">
      <h2 class="text-4xl font-bold text-guest-bg text-center mb-2 font-kaushan text-shadow-md text-shadow-black/35">
        Who's Coming
      </h2>
      <p class="text-white text-center text-lg mb-8 text-shadow-md text-shadow-black/35">
        {{ count }} guest{{ count !== 1 ? 's' : '' }} attending
      </p>

      <div v-if="guests.length === 0" class="text-center text-white py-8 text-shadow-md text-shadow-black/35">
        No confirmed guests yet
      </div>

      <template v-else>
        <!-- Desktop: show preview names, Mobile: hidden until expanded -->
        <div v-if="!expanded" class="hidden sm:flex flex-wrap justify-center gap-3">
          <div
            v-for="(guest, index) in visibleGuests"
            :key="index"
            class="px-4 py-2 text-white text-shadow-md text-shadow-black/35"
          >
            {{ guest.displayName }}
          </div>
        </div>

        <!-- Expanded: full list on all screens -->
        <div v-if="expanded" class="flex flex-wrap justify-center gap-3">
          <div
            v-for="(guest, index) in guests"
            :key="index"
            class="px-4 py-2 text-white text-shadow-md text-shadow-black/35"
          >
            {{ guest.displayName }}
          </div>
        </div>

        <div class="text-center mt-4">
          <button
            @click="expanded = !expanded"
            class="text-white/80 hover:text-white text-sm underline underline-offset-4 transition-colors cursor-pointer text-shadow-md text-shadow-black/35"
          >
            {{ expanded ? 'Show less' : `View all ${guests.length} guests` }}
          </button>
        </div>
      </template>
    </div>
  </section>
</template>

<style scoped>
.text-outline {
  text-shadow:
    -1px -1px 0 rgba(255, 255, 255, 0.4),
    1px -1px 0 rgba(255, 255, 255, 0.4),
    -1px 1px 0 rgba(255, 255, 255, 0.4),
    1px 1px 0 rgba(255, 255, 255, 0.4);
}
</style>
