<script setup lang="ts">
import type { Gift } from '@/types'

defineProps<{
  gifts: Gift[]
}>()
</script>

<template>
  <section class="gift-guide py-16 px-4">
    <div class="max-w-3xl mx-auto">
      <h2 class="text-4xl font-bold text-heading text-center mb-8 font-kaushan">Gift Guide</h2>

      <div class="space-y-6">
        <div
          v-for="gift in gifts"
          :key="gift.id"
          class="pb-6 border-b border-muted/50 last:border-0 last:pb-0 text-center"
        >
          <h3 class="text-xl font-normal text-primary mb-1">{{ gift.title }}</h3>
          <p v-if="gift.description" class="text-gray-600 text-sm mb-2">
            {{ gift.description }}
          </p>
          <a
            v-if="gift.link"
            :href="gift.link"
            target="_blank"
            rel="noopener noreferrer"
            class="text-accent hover:underline text-sm inline-flex items-center gap-1"
          >
            View
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
            </svg>
          </a>
        </div>
      </div>
    </div>
  </section>
</template>
