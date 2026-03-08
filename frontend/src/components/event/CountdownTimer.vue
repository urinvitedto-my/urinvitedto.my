<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps<{
  targetDate: string
  customMessage?: string
}>()

const now = ref(Date.now())
let interval: ReturnType<typeof setInterval> | null = null

const timeLeft = computed(() => {
  const target = new Date(props.targetDate).getTime()
  const diff = target - now.value

  if (diff <= 0) {
    if (interval) {
      clearInterval(interval)
      interval = null
    }
    return { days: 0, hours: 0, minutes: 0, seconds: 0, expired: true }
  }

  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
  const seconds = Math.floor((diff % (1000 * 60)) / 1000)

  return { days, hours, minutes, seconds, expired: false }
})

onMounted(() => {
  interval = setInterval(() => {
    now.value = Date.now()
  }, 1000)
})

onUnmounted(() => {
  if (interval) clearInterval(interval)
})
</script>

<template>
  <section class="countdown-timer py-16 px-4">
    <div class="max-w-3xl mx-auto text-center">
      <p v-if="customMessage" class="text-lg mb-6 text-gray-500">
        {{ customMessage }}
      </p>

      <div v-if="timeLeft.expired" class="text-2xl font-bold text-primary-dark">
        The event has started!
      </div>

      <div v-else class="grid grid-cols-4 gap-4 max-w-md mx-auto">
        <div class="text-center">
          <div class="text-4xl md:text-6xl font-bold text-primary-dark countdown-number">
            {{ timeLeft.days }}
          </div>
          <div class="text-sm text-primary-dark mt-2 uppercase tracking-wider">Days</div>
        </div>
        <div class="text-center">
          <div class="text-4xl md:text-6xl font-bold text-primary-dark countdown-number">
            {{ timeLeft.hours }}
          </div>
          <div class="text-sm text-primary-dark mt-2 uppercase tracking-wider">Hours</div>
        </div>
        <div class="text-center">
          <div class="text-4xl md:text-6xl font-bold text-primary-dark countdown-number">
            {{ timeLeft.minutes }}
          </div>
          <div class="text-sm text-primary-dark mt-2 uppercase tracking-wider">Minutes</div>
        </div>
        <div class="text-center">
          <div class="text-4xl md:text-6xl font-bold text-primary-dark countdown-number">
            {{ timeLeft.seconds }}
          </div>
          <div class="text-sm text-primary-dark mt-2 uppercase tracking-wider">Seconds</div>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.countdown-number {
  text-shadow:
    -1px -1px 0 rgba(255, 255, 255, 0.4),
     1px -1px 0 rgba(255, 255, 255, 0.4),
    -1px  1px 0 rgba(255, 255, 255, 0.4),
     1px  1px 0 rgba(255, 255, 255, 0.4);
}
</style>
