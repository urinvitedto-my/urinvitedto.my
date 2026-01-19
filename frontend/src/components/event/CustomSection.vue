<script setup lang="ts">
import type { CustomSection } from '@/types'

defineProps<{
  section: CustomSection
}>()
</script>

<template>
  <section
    class="custom-section py-12 px-4"
    :style="{ backgroundColor: section.bgColor || '#ececec' }"
  >
    <div class="max-w-3xl mx-auto">
      <h2 class="text-2xl font-bold text-[#14213d] text-center mb-6">
        {{ section.title }}
      </h2>

      <!-- HTML content -->
      <div
        class="prose max-w-none mb-6"
        v-html="section.content"
      ></div>

      <!-- Optional image -->
      <img
        v-if="section.image"
        :src="section.image"
        :alt="section.title"
        class="w-full rounded-lg"
        loading="lazy"
      />
    </div>
  </section>
</template>

<style scoped></style>
