<script setup lang="ts">
import type { MonetaryGiftsConfig } from '@/types'

defineProps<{
  config: MonetaryGiftsConfig
}>()

/** Fetches the image as a blob and triggers a browser download. */
async function downloadQr(url: string, label: string) {
  const res = await fetch(url)
  const blob = await res.blob()
  const a = document.createElement('a')
  a.href = URL.createObjectURL(blob)
  a.download = `${label}-qr-code.png`
  a.click()
  URL.revokeObjectURL(a.href)
}
</script>

<template>
  <section class="monetary-gifts py-16 px-4">
    <div class="max-w-3xl mx-auto text-center">
      <h2 class="text-4xl font-bold text-heading mb-4 font-kaushan">Monetary Gifts</h2>

      <div v-if="config.accounts?.length" class="space-y-6">
        <p class="text-xs text-gray-400 mt-1">Click image to download</p>
        <div
          v-for="(account, index) in config.accounts"
          :key="index"
          class="border-b border-muted/50 last:border-0"
        >
          <div class="font-semibold text-primary">{{ account.method }}</div>
          <div v-if="account.qrCodeUrl" class="my-4">
            <img
              :src="account.qrCodeUrl"
              :alt="`${account.method} QR Code`"
              class="w-60 h-60 mx-auto rounded-lg cursor-pointer hover:scale-105 transition-transform duration-200"
              @click="downloadQr(account.qrCodeUrl!, account.method)"
            />
          </div>
          <div v-if="account.number" class="text-gray-600">{{ account.number }}</div>
          <div class="text-sm text-gray-500">{{ account.name }}</div>
        </div>
      </div>
    </div>
  </section>
</template>
