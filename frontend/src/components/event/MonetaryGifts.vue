<script setup lang="ts">
import type { MonetaryGiftsConfig } from '@/types'

defineProps<{
  config: MonetaryGiftsConfig
}>()
</script>

<template>
  <section class="monetary-gifts py-16 px-4">
    <div class="max-w-3xl mx-auto text-center">
      <h2 class="text-4xl font-bold text-heading mb-4 font-kaushan">Monetary Gifts</h2>

      <!-- QR Code -->
      <div v-if="config.qrCodeUrl" class="mb-6">
        <img
          :src="config.qrCodeUrl"
          alt="Payment QR Code"
          class="w-48 h-48 mx-auto rounded-lg"
        />
      </div>

      <!-- Instructions -->
      <p v-if="config.instructions" class="text-gray-600 mb-6">
        {{ config.instructions }}
      </p>

      <!-- Account details -->
      <div v-if="config.accounts?.length" class="space-y-3">
        <div
          v-for="(account, index) in config.accounts"
          :key="index"
          class="py-3 border-b border-muted/50 last:border-0"
        >
          <div class="font-semibold text-primary">{{ account.method }}</div>
          <div class="text-gray-600">{{ account.number }}</div>
          <div class="text-sm text-gray-500">{{ account.name }}</div>
        </div>
      </div>
    </div>
  </section>
</template>
