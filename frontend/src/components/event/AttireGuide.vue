<script setup lang="ts">
import { computed } from 'vue'
import type { AttireGuide } from '@/types'

const props = defineProps<{
  attireGuide: AttireGuide
}>()

/** Splits description into segments, bolding parts that match the title. */
const descriptionParts = computed(() => {
  const { title, description } = props.attireGuide
  if (!title || !description.includes(title)) {
    return [{ text: description, bold: false }]
  }

  const parts: { text: string; bold: boolean }[] = []
  let remaining = description
  let idx = remaining.indexOf(title)

  while (idx !== -1) {
    if (idx > 0) parts.push({ text: remaining.slice(0, idx), bold: false })
    parts.push({ text: title, bold: true })
    remaining = remaining.slice(idx + title.length)
    idx = remaining.indexOf(title)
  }

  if (remaining) parts.push({ text: remaining, bold: false })
  return parts
})
</script>

<template>
  <section class="attire-guide py-16 px-4">
    <div class="max-w-3xl mx-auto text-center">
      <h2 class="text-4xl font-bold text-primary-dark mb-4 font-kaushan">Attire Guide</h2>

      <p class="text-gray-600 mb-4 text-xl">
        <template v-for="(part, i) in descriptionParts" :key="i">
          <strong v-if="part.bold" class="block my-1">{{ part.text }}</strong>
          <template v-else>{{ part.text.trim() }}</template>
        </template>
      </p>

      <img
        v-if="attireGuide.imageUrl"
        :src="attireGuide.imageUrl"
        alt="Attire motif"
        class="mx-auto max-w-sm w-full object-cover"
      />
    </div>
  </section>
</template>
