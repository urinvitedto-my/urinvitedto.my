<script setup lang="ts">
defineProps<{
  guests: { displayName: string }[]
  count: number
}>()
</script>

<template>
  <section class="confirmed-guests py-16 px-4">
    <div class="max-w-3xl mx-auto">
      <h2 class="text-4xl font-bold text-guest-bg text-center mb-2 font-kaushan">Who's Coming</h2>
      <p class="text-guest-bg text-center mb-8">{{ count }} guest{{ count !== 1 ? 's' : '' }} attending</p>

      <div v-if="guests.length === 0" class="text-center text-guest-bg py-8">
        No confirmed guests yet
      </div>

      <div v-else class="flex flex-wrap justify-center gap-3">
        <div
          v-for="(guest, index) in guests"
          :key="index"
          class="px-4 py-2 text-guest-bg"
        >
          {{ guest.displayName }}
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.text-outline {
  text-shadow:
    -1px -1px 0 rgba(255, 255, 255, 0.4),
     1px -1px 0 rgba(255, 255, 255, 0.4),
    -1px  1px 0 rgba(255, 255, 255, 0.4),
     1px  1px 0 rgba(255, 255, 255, 0.4);
}
</style>
