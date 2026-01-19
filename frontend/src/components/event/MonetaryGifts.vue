<script setup lang="ts">
import type { MonetaryGiftsConfig } from '@/types'

defineProps<{
  config: MonetaryGiftsConfig
}>()
</script>

<template>
  <section class="monetary-gifts bg-[#ececec] py-12 px-4">
    <div class="max-w-3xl mx-auto text-center">
      <h2 class="text-2xl font-bold text-[#14213d] mb-4">Gift</h2>

      <div class="bg-white rounded-lg shadow-sm p-8">
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
        <div v-if="config.accounts?.length" class="space-y-4">
          <div
            v-for="(account, index) in config.accounts"
            :key="index"
            class="bg-[#ececec] rounded-lg p-4"
          >
            <div class="font-semibold text-[#14213d]">{{ account.method }}</div>
            <div class="text-gray-600">{{ account.number }}</div>
            <div class="text-sm text-gray-500">{{ account.name }}</div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped></style>
