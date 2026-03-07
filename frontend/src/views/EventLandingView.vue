<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useEventStore } from '@/stores/event'
import type { EventType } from '@/types'
import LoadingSpinner from '@/components/LoadingSpinner.vue'

const props = defineProps<{
  type: EventType
  slug: string
}>()

const router = useRouter()
const eventStore = useEventStore()
const inviteCode = ref('')
const submitting = ref(false)

const loading = ref(true)
const error = ref('')
const { eventSummary, eventDetails } = storeToRefs(eventStore)

onMounted(async () => {
  await loadEvent()
})

onUnmounted(() => {
  eventStore.$reset()
})

/** Loads event summary, then details if public. */
async function loadEvent() {
  loading.value = true
  error.value = ''

  try {
    await eventStore.fetchSummary(props.type, props.slug)

    if (eventStore.eventSummary?.isPublic) {
      await eventStore.fetchDetails(props.type, props.slug)
    }
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : 'Failed to load event'
  } finally {
    loading.value = false
  }
}

/** Handles invite code submission for private events. */
async function handleInviteSubmit() {
  if (!inviteCode.value.trim()) {
    error.value = 'Please enter your invite code'
    return
  }

  submitting.value = true
  error.value = ''

  try {
    await eventStore.fetchDetails(props.type, props.slug, inviteCode.value.toUpperCase())
    router.push({
      name: 'guest',
      params: { type: props.type, slug: props.slug },
      query: { invite: inviteCode.value.toUpperCase() },
    })
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : 'Invalid invite code'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="event-landing-view">
    <!-- Loading -->
    <div v-if="loading" class="landing-fullscreen bg-primary flex items-center justify-center">
      <LoadingSpinner />
    </div>

    <!-- Error (no event loaded) -->
    <div v-else-if="error && !eventSummary" class="landing-fullscreen bg-primary flex items-center justify-center">
      <div class="text-center px-4">
        <p class="text-red-400 mb-4">{{ error }}</p>
        <button @click="loadEvent" class="text-accent hover:underline">Try again</button>
      </div>
    </div>

    <!-- Event Content -->
    <div
      v-else-if="eventSummary"
      class="landing-fullscreen bg-cover bg-center"
      :class="{ 'bg-primary': !eventSummary.coverImageUrl }"
      :style="eventSummary.coverImageUrl
        ? { backgroundImage: `url(${eventSummary.coverImageUrl})` }
        : undefined"
    >
      <div class="absolute inset-0 bg-black/40"></div>

      <div class="relative z-10 flex flex-col items-center justify-center h-full px-4 text-center">
        <h1 class="landing-title text-white">
          {{ eventSummary.title }}
        </h1>

        <!-- Public Event: Open Invitation button -->
        <div v-if="eventSummary.isPublic && eventDetails" class="mt-6">
          <RouterLink
            :to="{ name: 'guest', params: { type, slug } }"
            class="landing-btn"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
            </svg>
            Open Invitation
          </RouterLink>
        </div>

        <!-- Private Event: Invite Code Entry -->
        <div v-if="!eventSummary.isPublic" class="mt-6 w-full max-w-sm">
          <form @submit.prevent="handleInviteSubmit" class="flex flex-col items-center gap-4">
            <p class="text-white/80 text-sm font-medium tracking-wide">
              Enter Your Invite Code
            </p>

            <input
              v-model="inviteCode"
              type="text"
              maxlength="6"
              placeholder="ABC123"
              class="w-full px-4 py-3 text-center text-xl tracking-widest font-mono rounded-lg bg-white/15 backdrop-blur border border-white/30 text-white placeholder-white/50 focus:outline-none focus:ring-2 focus:ring-white/50 uppercase"
              :disabled="submitting"
            />

            <button
              type="submit"
              :disabled="submitting"
              class="landing-btn"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
              </svg>
              {{ submitting ? 'Checking Guest List...' : 'Open Invitation' }}
            </button>

            <p v-if="error" class="text-red-300 text-center text-sm">{{ error }}</p>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.landing-fullscreen {
  position: relative;
  width: 100%;
  height: 100vh;
  height: 100dvh;
  overflow: hidden;
}

.landing-title {
  font-family: 'Lavishly Yours', cursive;
  font-size: clamp(4rem, 14vw, 9rem);
  font-weight: 400;
  letter-spacing: normal;
  text-transform: none;
  line-height: 1.3;
}

.landing-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  width: 100%;
  padding: 0.875rem 2rem;
  background-color: color-mix(in srgb, var(--color-accent) 90%, transparent);
  color: black;
  font-size: 1.05rem;
  font-weight: 600;
  letter-spacing: 0.05em;
  border-radius: 0.5rem;
  transition: background-color 0.2s;
  text-decoration: none;
}

.landing-btn:hover {
  background-color: var(--color-accent);
}

.landing-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
