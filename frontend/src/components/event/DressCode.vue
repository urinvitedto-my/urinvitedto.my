<script setup lang="ts">
import type { DressCode } from '@/types'

defineProps<{
  dressCode: DressCode
}>()
</script>

<template>
  <section class="dress-code py-16 px-4">
    <div class="max-w-3xl mx-auto text-center">
      <h2 class="text-4xl font-bold text-heading mb-4 font-kaushan">Dress Code</h2>

      <h3 class="text-xl font-semibold text-heading mb-2">
        {{ dressCode.title }}
      </h3>

      <p class="text-gray-600 mb-4">{{ dressCode.description }}</p>

      <p v-if="dressCode.notes" class="text-sm text-gray-500 italic mb-4">
        {{ dressCode.notes }}
      </p>

      <ul v-if="dressCode.examples?.length" class="flex flex-wrap justify-center gap-2">
        <li
          v-for="example in dressCode.examples"
          :key="example"
          class="px-3 py-1 bg-surface text-gray-700 rounded-full text-sm"
        >
          {{ example }}
        </li>
      </ul>
    </div>
  </section>
</template>
