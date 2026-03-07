<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useEventStore } from '@/stores/event'
import type { EventType } from '@/types'
import LoadingSpinner from '@/components/LoadingSpinner.vue'

import EventDetails from '@/components/event/EventDetails.vue'
import LocationPhoto from '@/components/event/LocationPhoto.vue'
import CountdownTimer from '@/components/event/CountdownTimer.vue'
import EventMap from '@/components/event/EventMap.vue'
import EventSchedule from '@/components/event/EventSchedule.vue'
import EventGallery from '@/components/event/EventGallery.vue'
import DressCode from '@/components/event/DressCode.vue'
import EventFAQ from '@/components/event/EventFAQ.vue'
import MonetaryGifts from '@/components/event/MonetaryGifts.vue'
import GiftGuide from '@/components/event/GiftGuide.vue'
import CustomSection from '@/components/event/CustomSection.vue'
import InviteRSVP from '@/components/event/InviteRSVP.vue'
import ConfirmedGuests from '@/components/event/ConfirmedGuests.vue'
import GuestBottomNav from '@/components/GuestBottomNav.vue'

const props = defineProps<{
  type: EventType
  slug: string
}>()

const route = useRoute()
const eventStore = useEventStore()

const {
  loading,
  error,
  eventDetails: eventData,
  confirmedGuests,
  orderedComponents,
} = storeToRefs(eventStore)

const inviteCode = computed(() => {
  const code = route.query.invite
  return typeof code === 'string' ? code.toUpperCase() : ''
})

const isMuted = ref(false)
let audioEl: HTMLAudioElement | null = null
const activationEvents = ['scroll', 'click', 'touchstart', 'keydown'] as const

/** Starts music on first user interaction. */
function startMusic() {
  if (!audioEl) return
  audioEl.play().catch(() => {})
  activationEvents.forEach((evt) => window.removeEventListener(evt, startMusic))
}

/** Toggles background music mute state. */
function toggleMute() {
  if (!audioEl) return
  isMuted.value = !isMuted.value
  audioEl.muted = isMuted.value
  if (!isMuted.value) {
    audioEl.play().catch(() => {})
  }
}

/** Sets up the audio element if the event has a music URL. */
function initAudio() {
  const musicSrc = eventData.value?.event.musicUrl
  if (!musicSrc) return

  audioEl = new Audio(musicSrc)
  audioEl.loop = true
  audioEl.volume = 0.3
  activationEvents.forEach((evt) =>
    window.addEventListener(evt, startMusic, { once: true }),
  )
}

/** Tears down audio and removes any lingering listeners. */
function cleanupAudio() {
  activationEvents.forEach((evt) => window.removeEventListener(evt, startMusic))
  if (audioEl) {
    audioEl.pause()
    audioEl.src = ''
    audioEl = null
  }
}

onMounted(async () => {
  await loadEventData()
  initAudio()
})

onUnmounted(() => {
  cleanupAudio()
  eventStore.$reset()
})

/** Loads event details and confirmed guests in parallel. */
async function loadEventData() {
  try {
    const code = inviteCode.value || undefined
    await Promise.all([
      eventStore.fetchDetails(props.type, props.slug, code),
      eventStore.fetchConfirmedGuests(props.type, props.slug),
    ])
  } catch {
    // errors are already set in the event store
  }
}
</script>

<template>
  <div id="section-top" class="guest-view min-h-screen pb-24">
    <!-- Music toggle -->
    <button
      v-if="eventData?.event.musicUrl"
      @click="toggleMute"
      class="fixed right-2 z-100 p-2.5 text-primary transition-opacity hover:opacity-70 bottom-[calc(5rem+env(safe-area-inset-bottom))]"
      :title="isMuted ? 'Unmute music' : 'Mute music'"
      :aria-label="isMuted ? 'Unmute music' : 'Mute music'"
    >
      <svg
        v-if="!isMuted"
        class="w-5 h-5"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M15.536 8.464a5 5 0 010 7.072M17.95 6.05a8 8 0 010 11.9M11 5L6 9H2v6h4l5 4V5z"
        />
      </svg>
      <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M5.586 15H4a1 1 0 01-1-1v-4a1 1 0 011-1h1.586l4.707-4.707A1 1 0 0112 5v14a1 1 0 01-1.707.707L5.586 15zM17 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2"
        />
      </svg>
    </button>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-20">
      <LoadingSpinner />
    </div>

    <!-- Error -->
    <div v-else-if="error" class="max-w-md mx-auto py-20 px-4 text-center">
      <p class="text-red-600 mb-4">{{ error }}</p>
      <RouterLink :to="`/${type}/${slug}`" class="text-accent hover:underline">
        Go back
      </RouterLink>
    </div>

    <!-- Event Content -->
    <template v-else-if="eventData">
      <!-- Dynamic Components -->
      <template v-for="comp in orderedComponents" :key="comp.name">
        <EventDetails
          v-if="comp.name === 'EventDetails'"
          :event="eventData.event"
          :gallery="eventData.gallery"
        />

        <LocationPhoto
          v-else-if="comp.name === 'LocationPhoto' && eventData.event.locationPhotoUrl"
          :url="eventData.event.locationPhotoUrl"
          :location="eventData.event.location"
        />

        <div
          v-else-if="comp.name === 'CountdownTimer' && eventData.event.startsAt"
          id="section-countdown"
        >
          <CountdownTimer
            :target-date="eventData.event.startsAt"
            :custom-message="
              eventData.event.customContent?.countdownTimer?.customMessage
            "
          />
        </div>

        <div
          v-else-if="
            comp.name === 'EventMap' && eventData.event.customContent?.locationDetails
          "
          id="section-venue"
        >
          <EventMap
            :location-details="eventData.event.customContent.locationDetails"
            :address="eventData.event.location"
          />
        </div>

        <EventSchedule
          v-else-if="comp.name === 'EventSchedule' && eventData.schedule.length"
          :items="eventData.schedule"
        />

        <div
          v-else-if="comp.name === 'EventGallery' && eventData.gallery.length"
          id="section-gallery"
        >
          <EventGallery :items="eventData.gallery" />
        </div>

        <DressCode
          v-else-if="
            comp.name === 'DressCode' && eventData.event.customContent?.dressCode
          "
          :dress-code="eventData.event.customContent.dressCode"
        />

        <EventFAQ
          v-else-if="comp.name === 'EventFAQ' && eventData.faqs.length"
          :faqs="eventData.faqs"
        />

        <MonetaryGifts
          v-else-if="
            comp.name === 'MonetaryGifts' &&
            eventData.event.customContent?.monetaryGifts?.enabled
          "
          :config="eventData.event.customContent.monetaryGifts"
        />

        <GiftGuide
          v-else-if="comp.name === 'GiftGuide' && eventData.gifts.length"
          :gifts="eventData.gifts"
        />

        <template
          v-else-if="
            comp.name === 'CustomSections' &&
            eventData.event.customContent?.customSections?.length
          "
        >
          <CustomSection
            v-for="section in eventData.event.customContent.customSections"
            :key="section.id"
            :section="section"
          />
        </template>
      </template>

      <!-- Fixed Sections (always at bottom) -->
      <div v-if="eventData.invite" id="section-rsvp">
        <InviteRSVP
          :invite="eventData.invite"
          :type="type"
          :slug="slug"
          :invite-code="inviteCode"
          @rsvp-updated="loadEventData"
        />
      </div>

      <ConfirmedGuests
        v-if="confirmedGuests"
        :guests="confirmedGuests.guests"
        :count="confirmedGuests.count"
      />

      <GuestBottomNav />
    </template>
  </div>
</template>

<style scoped>
.guest-view {
  font-family: var(--font-gelasio);
  color: var(--color-guest-text);
  background: linear-gradient(
    to bottom,
    var(--color-primary-dark),
    var(--color-guest-bg) 20%,
    var(--color-guest-bg) 80%,
    var(--color-primary-dark)
  );
}

.guest-view :deep(h1),
.guest-view :deep(h2) {
  font-family: 'Kaushan Script', cursive;
  font-weight: 300;
  letter-spacing: 0.15em;
}
</style>
