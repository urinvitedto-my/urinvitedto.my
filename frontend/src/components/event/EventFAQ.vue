<script setup lang="ts">
import { ref } from 'vue'
import type { FAQ } from '@/types'

defineProps<{
  faqs: FAQ[]
}>()

const openItems = ref<Set<string>>(new Set())

/**
 * Toggles FAQ item open/closed.
 */
function toggle(id: string) {
  if (openItems.value.has(id)) {
    openItems.value.delete(id)
  } else {
    openItems.value.add(id)
  }
}
</script>

<template>
  <section class="event-faq py-16 px-4">
    <div class="max-w-3xl mx-auto">
      <h2 class="text-2xl font-bold text-[#2c2c2c] text-center mb-8">FAQ</h2>

      <div class="space-y-4">
        <div
          v-for="faq in faqs"
          :key="faq.id"
          class="border border-[#e5e5e5] bg-white/80 backdrop-blur rounded-xl overflow-hidden shadow-sm"
        >
          <button
            @click="toggle(faq.id)"
            class="w-full flex items-center justify-between px-6 py-4 text-left hover:bg-black/5 transition-colors"
          >
            <span class="font-medium text-[#2c2c2c]">{{ faq.question }}</span>
            <svg
              :class="[
                'w-5 h-5 text-gray-500 transition-transform',
                openItems.has(faq.id) ? 'rotate-180' : '',
              ]"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
          </button>
          <div
            v-show="openItems.has(faq.id)"
            class="px-6 pb-4 text-gray-600"
          >
            {{ faq.answer }}
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped></style>
