<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import { useAdminStore } from '@/stores/admin'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import AdminCreateEventForm from '@/components/admin/AdminCreateEventForm.vue'
import AdminEventSidebar from '@/components/admin/AdminEventSidebar.vue'
import AdminEventDetail from '@/components/admin/AdminEventDetail.vue'

const authStore = useAuthStore()
const adminStore = useAdminStore()

const { isAdmin } = storeToRefs(authStore)
const {
  events,
  loading: eventsLoading,
  error: eventsError,
  selectedEventId,
} = storeToRefs(adminStore)

const showCreateForm = ref(false)

const selectedEvent = computed(
  () => events.value.find((e) => e.id === selectedEventId.value) ?? null,
)

/** Whether the mobile detail view is showing (an event is selected on small screens). */
const mobileShowDetail = computed(() => !!selectedEventId.value)

function selectEvent(eventId: string) {
  selectedEventId.value = eventId
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

/** Goes back to event list on mobile. */
function handleBack() {
  selectedEventId.value = null
}

function handleDeleted() {
  selectedEventId.value = null
}

onMounted(async () => {
  try {
    await adminStore.fetchEvents()
  } catch {
    // error is in adminStore.error
  }
})

watch(isAdmin, async (newVal) => {
  if (newVal && adminStore.events.length === 0) {
    try {
      await adminStore.fetchEvents()
    } catch {
      // error is in adminStore.error
    }
  }
})
</script>

<template>
  <div class="admin-view min-h-screen pt-24 pb-8 px-4 md:mx-16">
    <div class="max-w-8xl mx-auto">
      <!--
        MOBILE (<lg): toggle between list and detail views.
        When no event selected -> show header + list.
        When event selected -> show back button + detail.
      -->

      <!-- Mobile: Detail View -->
      <template v-if="mobileShowDetail && selectedEvent">
        <div class="lg:hidden">
          <button
            @click="handleBack"
            class="flex items-center gap-1 text-sm text-primary hover:underline mb-4"
          >
            ← Back to events
          </button>

          <div class="bg-white rounded-lg shadow-sm p-4 sm:p-6">
            <AdminEventDetail
              :key="selectedEvent.id"
              :event="selectedEvent"
              @deleted="handleDeleted"
            />
          </div>
        </div>
      </template>

      <!-- Mobile: List View (hidden when detail is showing) -->
      <div :class="{ 'hidden lg:block': mobileShowDetail }">
        <!-- Header -->
        <div class="flex items-center justify-between mb-6">
          <h1 class="text-2xl font-bold text-primary">Admin Dashboard</h1>
          <button
            @click="showCreateForm = !showCreateForm"
            class="bg-accent text-black font-medium px-4 py-2 rounded-lg hover:bg-accent-dark transition-colors"
          >
            {{ showCreateForm ? 'Cancel' : '+ Create Event' }}
          </button>
        </div>

        <!-- Create Event Form -->
        <AdminCreateEventForm
          v-if="showCreateForm"
          class="mb-6"
          @created="showCreateForm = false"
          @cancel="showCreateForm = false"
        />

        <!-- Loading -->
        <div v-if="eventsLoading" class="flex items-center justify-center py-12">
          <LoadingSpinner size="md" />
        </div>

        <!-- Error -->
        <div
          v-else-if="eventsError"
          class="flex flex-col items-center justify-center text-center min-h-[60vh]"
        >
          <h2 class="text-2xl font-bold text-primary mb-2">Something went wrong</h2>
          <p class="text-base text-gray-500 mb-6 max-w-sm">{{ eventsError }}</p>
          <button
            @click="adminStore.fetchEvents()"
            class="bg-accent text-black font-semibold px-8 py-3 rounded-lg text-base hover:bg-accent-dark transition-colors"
          >
            Try again
          </button>
        </div>

        <!-- Empty -->
        <div
          v-else-if="events.length === 0"
          class="bg-white rounded-lg shadow-sm p-8 text-center"
        >
          <p class="text-gray-500">No events yet. Create your first event!</p>
        </div>

        <!-- Split View (desktop keeps both panels visible) -->
        <div v-else class="flex flex-col lg:flex-row gap-6 lg:items-start">
          <!-- Sidebar / Event List -->
          <div
            class="w-full lg:w-[340px] lg:sticky lg:top-24 lg:max-h-[calc(100vh-7rem)] shrink-0"
          >
            <AdminEventSidebar
              :events="events"
              :selected-event-id="selectedEventId"
              @select="selectEvent"
            />
          </div>

          <!-- Desktop Detail Panel (hidden on mobile, shown via lg:block) -->
          <div class="hidden lg:block flex-1 min-w-0">
            <div v-if="selectedEvent" class="bg-white rounded-lg shadow-sm p-6">
              <AdminEventDetail
                :key="selectedEvent.id"
                :event="selectedEvent"
                @deleted="handleDeleted"
              />
            </div>

            <div v-else class="bg-white rounded-lg shadow-sm p-8 text-center">
              <p class="text-gray-400">Select an event to manage</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
