<script setup lang="ts">
import type { Gift } from '@/types'

defineProps<{
  gifts: Gift[]
}>()
</script>

<template>
  <section class="gift-guide py-16 px-4">
    <div class="max-w-3xl mx-auto">
      <h2 class="text-2xl font-bold text-heading text-center mb-8">Gift Guide</h2>

      <div class="grid md:grid-cols-2 gap-4">
        <div
          v-for="gift in gifts"
          :key="gift.id"
          class="bg-white/80 backdrop-blur border border-muted/50 shadow-sm rounded-xl p-6"
        >
          <div class="flex items-start gap-4">
            <div class="w-10 h-10 bg-primary rounded-full flex items-center justify-center shrink-0">
              <svg class="w-5 h-5 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7" />
              </svg>
            </div>
            <div class="flex-1">
              <h3 class="font-semibold text-primary mb-1">{{ gift.title }}</h3>
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
      </div>
    </div>
  </section>
</template>
