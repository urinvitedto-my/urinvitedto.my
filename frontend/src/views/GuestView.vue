<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getEventDetails, getConfirmedGuests } from '@/services/api'
import type { EventDetailsResponse, ConfirmedGuestsResponse, ComponentConfig } from '@/types'

// event components
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

const props = defineProps<{
  type: string
  slug: string
}>()

const route = useRoute()
const loading = ref(true)
const error = ref('')
const eventData = ref<EventDetailsResponse | null>(null)
const confirmedGuests = ref<ConfirmedGuestsResponse | null>(null)

// get invite code from query
const inviteCode = computed(() => {
  const code = route.query.invite
  return typeof code === 'string' ? code.toUpperCase() : ''
})

// ordered components based on enabled_components config
const orderedComponents = computed(() => {
  if (!eventData.value?.event.enabledComponents?.components) {
    // default order if not configured
    return [
      { name: 'EventDetails', enabled: true, order: 1 },
      { name: 'LocationPhoto', enabled: true, order: 2 },
      { name: 'CountdownTimer', enabled: true, order: 3 },
      { name: 'EventMap', enabled: true, order: 4 },
      { name: 'EventSchedule', enabled: true, order: 5 },
      { name: 'EventGallery', enabled: true, order: 6 },
      { name: 'DressCode', enabled: true, order: 7 },
      { name: 'EventFAQ', enabled: true, order: 8 },
      { name: 'MonetaryGifts', enabled: true, order: 9 },
      { name: 'GiftGuide', enabled: true, order: 10 },
      { name: 'CustomSections', enabled: true, order: 11 },
    ] as ComponentConfig[]
  }
  return eventData.value.event.enabledComponents.components
    .filter((c) => c.enabled)
    .sort((a, b) => a.order - b.order)
})

onMounted(async () => {
  await loadEventData()
})

/**
 * Loads event details and confirmed guests.
 */
async function loadEventData() {
  loading.value = true
  error.value = ''

  try {
    const code = inviteCode.value || undefined
    const [details, guests] = await Promise.all([
      getEventDetails(props.type, props.slug, code),
      getConfirmedGuests(props.type, props.slug),
    ])
    eventData.value = details
    confirmedGuests.value = guests
  } catch (e: any) {
    error.value = e.message || 'Failed to load event'
  } finally {
    loading.value = false
  }
}

/**
 * Checks if a component should be rendered.
 */
function shouldRender(name: string): boolean {
  return orderedComponents.value.some((c) => c.name === name)
}
</script>

<template>
  <div class="guest-view min-h-screen pt-16">
    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-20">
      <div class="animate-spin rounded-full h-12 w-12 border-4 border-[#fca311] border-t-transparent"></div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="max-w-md mx-auto py-20 px-4 text-center">
      <p class="text-red-600 mb-4">{{ error }}</p>
      <RouterLink :to="`/${type}/${slug}`" class="text-[#fca311] hover:underline">
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
          :hosts="eventData.hosts"
        />

        <LocationPhoto
          v-else-if="comp.name === 'LocationPhoto' && eventData.event.locationPhotoUrl"
          :url="eventData.event.locationPhotoUrl"
          :location="eventData.event.location"
        />

        <CountdownTimer
          v-else-if="comp.name === 'CountdownTimer' && eventData.event.startsAt"
          :target-date="eventData.event.startsAt"
          :custom-message="eventData.event.customContent?.countdownTimer?.customMessage"
        />

        <EventMap
          v-else-if="comp.name === 'EventMap' && eventData.event.customContent?.locationDetails"
          :location-details="eventData.event.customContent.locationDetails"
          :address="eventData.event.location"
        />

        <EventSchedule
          v-else-if="comp.name === 'EventSchedule' && eventData.schedule.length"
          :items="eventData.schedule"
        />

        <EventGallery
          v-else-if="comp.name === 'EventGallery' && eventData.gallery.length"
          :items="eventData.gallery"
        />

        <DressCode
          v-else-if="comp.name === 'DressCode' && eventData.event.customContent?.dressCode"
          :dress-code="eventData.event.customContent.dressCode"
        />

        <EventFAQ
          v-else-if="comp.name === 'EventFAQ' && eventData.faqs.length"
          :faqs="eventData.faqs"
        />

        <MonetaryGifts
          v-else-if="comp.name === 'MonetaryGifts' && eventData.event.customContent?.monetaryGifts?.enabled"
          :config="eventData.event.customContent.monetaryGifts"
        />

        <GiftGuide
          v-else-if="comp.name === 'GiftGuide' && eventData.gifts.length"
          :gifts="eventData.gifts"
        />

        <template v-else-if="comp.name === 'CustomSections' && eventData.event.customContent?.customSections?.length">
          <CustomSection
            v-for="section in eventData.event.customContent.customSections"
            :key="section.id"
            :section="section"
          />
        </template>
      </template>

      <!-- Fixed Sections (always at bottom) -->
      <InviteRSVP
        v-if="eventData.invite"
        :invite="eventData.invite"
        :type="type"
        :slug="slug"
        :invite-code="inviteCode"
        @rsvp-updated="loadEventData"
      />

      <ConfirmedGuests
        v-if="confirmedGuests"
        :guests="confirmedGuests.guests"
        :count="confirmedGuests.count"
      />
    </template>
  </div>
</template>

<style scoped></style>
