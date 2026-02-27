<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useEventStore } from '@/stores/event'
import { getEventDetails } from '@/services/api'

const props = defineProps<{
  type: string
  slug: string
}>()

const router = useRouter()
const eventStore = useEventStore()
const inviteCode = ref('')
const submitting = ref(false)

const loading = ref(true)
const error = ref('')
const eventSummary = computed(() => eventStore.eventSummary)
const eventDetails = computed(() => eventStore.eventDetails)

const isMuted = ref(false)
let audioEl: HTMLAudioElement | null = null
const activationEvents = ['scroll', 'click', 'touchstart', 'keydown'] as const

/**
 * Starts music on first user interaction (works across all browsers).
 */
function startMusic() {
  if (!audioEl) return
  audioEl.play().catch(() => {})
  activationEvents.forEach((evt) => window.removeEventListener(evt, startMusic))
}

/**
 * Toggles background music mute state.
 */
function toggleMute() {
  if (!audioEl) return
  isMuted.value = !isMuted.value
  audioEl.muted = isMuted.value
  if (!isMuted.value) {
    audioEl.play().catch(() => {})
  }
}

/**
 * Sets up the audio element if the event has a music URL.
 */
function initAudio() {
  const musicSrc = eventSummary.value?.musicUrl
  if (!musicSrc) return

  audioEl = new Audio(musicSrc)
  audioEl.loop = true
  audioEl.volume = 0.3
  activationEvents.forEach((evt) =>
    window.addEventListener(evt, startMusic, { once: true }),
  )
}

/**
 * Tears down audio and removes any lingering listeners.
 */
function cleanupAudio() {
  activationEvents.forEach((evt) => window.removeEventListener(evt, startMusic))
  if (audioEl) {
    audioEl.pause()
    audioEl.src = ''
    audioEl = null
  }
}

onMounted(async () => {
  await loadEvent()
  initAudio()
})

onUnmounted(() => {
  cleanupAudio()
  eventStore.$reset()
})

/**
 * Loads event summary, then details if public.
 */
async function loadEvent() {
  loading.value = true
  error.value = ''

  try {
    await eventStore.fetchSummary(props.type, props.slug)

    if (eventStore.eventSummary?.isPublic) {
      await eventStore.fetchDetails(props.type, props.slug)
    }
  } catch (e: any) {
    error.value = e.message || 'Failed to load event'
  } finally {
    loading.value = false
  }
}

/**
 * Handles invite code submission for private events.
 */
async function handleInviteSubmit() {
  if (!inviteCode.value.trim()) {
    error.value = 'Please enter your invite code'
    return
  }

  submitting.value = true
  error.value = ''

  try {
    await getEventDetails(props.type, props.slug, inviteCode.value.toUpperCase())
    router.push({
      name: 'guest',
      params: { type: props.type, slug: props.slug },
      query: { invite: inviteCode.value.toUpperCase() },
    })
  } catch (e: any) {
    error.value = e.message || 'Invalid invite code'
  } finally {
    submitting.value = false
  }
}

/**
 * Formats date for display.
 */
function formatDate(dateStr?: string): string {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('en-US', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: 'numeric',
    minute: '2-digit',
  })
}
</script>

<template>
  <div class="event-landing-view min-h-screen">
    <!-- Music toggle -->
    <button
      v-if="eventSummary?.musicUrl"
      @click="toggleMute"
      class="fixed bottom-6 right-6 z-50 bg-[#14213d] text-white p-3 rounded-full shadow-lg hover:bg-[#1a2a4d] transition-colors"
      :title="isMuted ? 'Unmute music' : 'Mute music'"
    >
      <svg v-if="!isMuted" class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.536 8.464a5 5 0 010 7.072M17.95 6.05a8 8 0 010 11.9M11 5L6 9H2v6h4l5 4V5z" />
      </svg>
      <svg v-else class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5.586 15H4a1 1 0 01-1-1v-4a1 1 0 011-1h1.586l4.707-4.707A1 1 0 0112 5v14a1 1 0 01-1.707.707L5.586 15zM17 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2" />
      </svg>
    </button>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center pt-32 pb-20">
      <div class="animate-spin rounded-full h-12 w-12 border-4 border-[#fca311] border-t-transparent"></div>
    </div>

    <!-- Error -->
    <div v-else-if="error && !eventSummary" class="max-w-md mx-auto pt-32 pb-20 px-4 text-center">
      <p class="text-red-600 mb-4">{{ error }}</p>
      <button @click="loadEvent" class="text-[#fca311] hover:underline">Try again</button>
    </div>

    <!-- Event Content -->
    <template v-else-if="eventSummary">
      <!-- Cover Image -->
      <div
        v-if="eventSummary.coverImageUrl"
        class="w-full h-80 md:h-112 bg-cover bg-center"
        :style="{ backgroundImage: `url(${eventSummary.coverImageUrl})` }"
      >
        <div class="w-full h-full bg-black/40 flex items-center justify-center">
          <h1 class="text-3xl md:text-5xl font-bold text-white text-center px-4">
            {{ eventSummary.title }}
          </h1>
        </div>
      </div>

      <!-- No cover image fallback -->
      <div v-else class="bg-[#14213d] pt-16 pb-16 px-4 text-center">
        <h1 class="text-3xl md:text-5xl font-bold text-white">
          {{ eventSummary.title }}
        </h1>
      </div>

      <!-- Event Info -->
      <div class="max-w-2xl mx-auto py-8 px-4">
        <div class="bg-white rounded-lg shadow-sm p-6 mb-6">
          <div class="flex items-center gap-2 text-gray-600 mb-4">
            <span class="inline-block px-3 py-1 bg-[#fca311] text-black text-sm font-medium rounded-full capitalize">
              {{ eventSummary.type }}
            </span>
          </div>

          <div v-if="eventSummary.startsAt" class="flex items-center gap-3 text-gray-700 mb-3">
            <svg class="w-5 h-5 text-[#14213d]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
            <span>{{ formatDate(eventSummary.startsAt) }}</span>
          </div>

          <div v-if="eventSummary.location" class="flex items-center gap-3 text-gray-700">
            <svg class="w-5 h-5 text-[#14213d]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            <span>{{ eventSummary.location }}</span>
          </div>
        </div>

        <!-- Private Event: Invite Code Entry -->
        <div v-if="!eventSummary.isPublic" class="bg-white rounded-lg shadow-sm p-6">
          <h2 class="text-xl font-semibold text-[#14213d] mb-4 text-center">
            Enter Your Invite Code
          </h2>
          <p class="text-gray-600 text-center mb-6">
            Please enter the 6-character code from your invitation.
          </p>

          <form @submit.prevent="handleInviteSubmit" class="space-y-4">
            <div>
              <input
                v-model="inviteCode"
                type="text"
                maxlength="6"
                placeholder="ABC123"
                class="w-full px-4 py-3 text-center text-2xl tracking-widest font-mono border border-[#e5e5e5] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#fca311] uppercase"
                :disabled="submitting"
              />
            </div>

            <p v-if="error" class="text-red-600 text-center text-sm">{{ error }}</p>

            <button
              type="submit"
              :disabled="submitting"
              class="w-full bg-[#fca311] text-black font-semibold py-3 rounded-lg hover:bg-[#e5930f] transition-colors disabled:opacity-50"
            >
              {{ submitting ? 'Checking...' : 'Open Invite' }}
            </button>
          </form>
        </div>

        <!-- Public Event: Show Details -->
        <div v-else-if="eventDetails" class="bg-white rounded-lg shadow-sm p-6">
          <p v-if="eventDetails.event.description" class="text-gray-700 whitespace-pre-wrap mb-6">
            {{ eventDetails.event.description }}
          </p>
          <p v-else class="text-gray-500 italic mb-6">No additional details available.</p>

          <RouterLink
            :to="{ name: 'guest', params: { type, slug } }"
            class="block w-full bg-[#fca311] text-black font-semibold py-3 rounded-lg hover:bg-[#e5930f] transition-colors text-center"
          >
            View Event
          </RouterLink>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped></style>
