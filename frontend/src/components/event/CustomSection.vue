<script setup lang="ts">
import type { CustomSection } from '@/types'

defineProps<{
  section: CustomSection
}>()
</script>

<template>
  <section
    class="custom-section py-16 px-4"
    :style="section.bgColor ? { backgroundColor: section.bgColor } : {}"
  >
    <div class="max-w-3xl mx-auto">
      <h2 v-if="section.title" class="text-2xl font-bold text-heading text-center mb-6">
        {{ section.title }}
      </h2>

      <!-- Admin-authored HTML content (only editable by admins, no sanitization needed) -->
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
