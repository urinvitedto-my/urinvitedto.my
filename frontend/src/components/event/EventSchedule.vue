<script setup lang="ts">
import type { ScheduleItem } from '@/types'

defineProps<{
  items: ScheduleItem[]
}>()

/**
 * Formats time for display.
 */
function formatTime(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleTimeString('en-US', {
    hour: 'numeric',
    minute: '2-digit',
  })
}
</script>

<template>
  <section class="event-schedule bg-[#ececec] py-12 px-4">
    <div class="max-w-3xl mx-auto">
      <h2 class="text-2xl font-bold text-[#14213d] text-center mb-8">Schedule</h2>

      <div class="relative">
        <!-- Timeline line -->
        <div class="absolute left-4 md:left-1/2 top-0 bottom-0 w-0.5 bg-[#e5e5e5] transform md:-translate-x-1/2"></div>

        <!-- Items -->
        <div class="space-y-8">
          <div
            v-for="(item, index) in items"
            :key="item.id"
            class="relative flex items-start gap-6 md:gap-0"
          >
            <!-- Timeline dot -->
            <div class="absolute left-4 md:left-1/2 w-4 h-4 bg-[#fca311] rounded-full transform -translate-x-1/2 mt-1.5 z-10"></div>

            <!-- Content -->
            <div
              :class="[
                'ml-10 md:ml-0 bg-white rounded-lg shadow-sm p-6 flex-1',
                index % 2 === 0 ? 'md:mr-[52%]' : 'md:ml-[52%]',
              ]"
            >
              <div class="text-[#fca311] font-semibold mb-1">
                {{ formatTime(item.time) }}
              </div>
              <h3 class="text-lg font-semibold text-[#14213d] mb-2">
                {{ item.title }}
              </h3>
              <p v-if="item.description" class="text-gray-600">
                {{ item.description }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped></style>
