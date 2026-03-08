<script setup lang="ts">
defineProps<{
  url: string
  location?: string
}>()
</script>

<template>
  <section class="location-photo py-16 px-4">
    <h2 class="text-4xl font-bold text-primary-dark text-center mb-2 font-kaushan">Venue</h2>
    <div v-if="location" class="flex items-center justify-center gap-2 text-gray-600 mb-6">
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
        />
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
        />
      </svg>
      <span class="font-medium">{{ location }}</span>
    </div>
    <div class="max-w-2xl mx-auto relative">
      <div class="absolute inset-y-0 left-0 w-16 bg-linear-to-r from-background to-transparent z-10" />
      <div class="absolute inset-y-0 right-0 w-16 bg-linear-to-l from-background to-transparent z-10" />
      <img
        :src="url"
        alt="Venue photo"
        class="w-full h-80 md:h-130 object-cover rounded-lg"
      />
    </div>
  </section>
</template>
