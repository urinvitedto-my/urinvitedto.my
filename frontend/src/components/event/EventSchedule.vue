<script setup lang="ts">
import type { ScheduleItem } from '@/types'
import { formatTimeOnly } from '@/utils/date'

defineProps<{
  items: ScheduleItem[]
}>()
</script>

<template>
  <section class="event-schedule py-16 px-4">
    <div class="max-w-2xl mx-auto">
      <h2 class="text-4xl font-bold text-heading text-center mb-8 font-kaushan">Schedule</h2>

      <div class="relative pl-6">
        <!-- Timeline line -->
        <div class="absolute left-1 top-2 bottom-2 w-px bg-muted"></div>

        <div class="space-y-5">
          <div
            v-for="item in items"
            :key="item.id"
            class="relative flex items-start justify-between gap-4"
          >
            <!-- Timeline dot -->
            <div
              class="absolute -left-6 top-1.5 w-[9px] h-[9px] bg-accent rounded-full ring-2 ring-guest-bg"
            ></div>

            <div class="flex-1">
              <h3 class="font-semibold text-heading">
                {{ item.title }}
              </h3>
              <p v-if="item.description" class="text-gray-600 text-sm mt-1">
                {{ item.description }}
              </p>
            </div>
            <span class="text-accent text-sm font-semibold shrink-0 pt-0.5">
              {{ formatTimeOnly(item.time) }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
