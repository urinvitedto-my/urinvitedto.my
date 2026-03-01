<script setup lang="ts">
import type { DressCode } from '@/types'

defineProps<{
  dressCode: DressCode
}>()
</script>

<template>
  <section class="dress-code py-16 px-4">
    <div class="max-w-3xl mx-auto text-center">
      <h2 class="text-2xl font-bold text-[#2c2c2c] mb-4">Dress Code</h2>

      <div class="bg-white/80 backdrop-blur rounded-xl border border-[#e5e5e5]/50 shadow-sm p-8">
        <div class="w-16 h-16 bg-[#14213d] rounded-full flex items-center justify-center mx-auto mb-4">
          <svg class="w-8 h-8 text-[#fca311]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01" />
          </svg>
        </div>

        <h3 class="text-xl font-semibold text-[#2c2c2c] mb-2">
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
            class="px-3 py-1 bg-[#ececec] text-gray-700 rounded-full text-sm"
          >
            {{ example }}
          </li>
        </ul>
      </div>
    </div>
  </section>
</template>

<style scoped></style>
