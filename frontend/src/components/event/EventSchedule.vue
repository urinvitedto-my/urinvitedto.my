<script setup lang="ts">
import type { ScheduleItem } from '@/types'
import { formatTimeOnly } from '@/utils/date'

defineProps<{
  items: ScheduleItem[]
}>()
</script>

<template>
  <section class="event-schedule py-16 px-4">
    <div class="max-w-3xl mx-auto">
      <h2 class="text-2xl font-bold text-heading text-center mb-8">Schedule</h2>

      <div class="relative">
        <!-- Timeline line -->
        <div class="absolute left-4 md:left-1/2 top-0 bottom-0 w-0.5 bg-gray-300 transform md:-translate-x-1/2"></div>

        <!-- Items -->
        <div class="space-y-8">
          <div
            v-for="(item, index) in items"
            :key="item.id"
            class="relative flex items-start gap-6 md:gap-0"
          >
            <!-- Timeline dot -->
            <div class="absolute left-4 md:left-1/2 w-4 h-4 bg-accent rounded-full transform -translate-x-1/2 mt-1.5 z-10"></div>

            <!-- Content -->
            <div
              :class="[
                'ml-10 md:ml-0 bg-white/80 backdrop-blur rounded-xl border border-muted/50 shadow-sm p-6 flex-1',
                index % 2 === 0 ? 'md:mr-[52%]' : 'md:ml-[52%]',
              ]"
            >
              <div class="text-accent font-semibold mb-1">
                {{ formatTimeOnly(item.time) }}
              </div>
              <h3 class="text-lg font-semibold text-heading mb-2">
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
