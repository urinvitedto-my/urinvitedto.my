<script setup lang="ts">
defineProps<{
  url: string
  location?: string
}>()
</script>

<template>
  <section class="location-photo">
    <div
      class="w-full h-64 md:h-96 bg-cover bg-center relative"
      :style="{ backgroundImage: `url(${url})` }"
    >
      <div class="absolute inset-0 bg-black/30 flex items-end">
        <div v-if="location" class="p-6 text-white">
          <div class="flex items-center gap-2">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            <span class="font-medium">{{ location }}</span>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped></style>
