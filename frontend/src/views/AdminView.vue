<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import { useAdminStore } from '@/stores/admin'
import { useToast } from '@/composables/useToast'
import { formatDate } from '@/utils/date'
import type { AdminEvent } from '@/types'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import AdminCreateEventForm from '@/components/admin/AdminCreateEventForm.vue'
import AdminEditEventForm from '@/components/admin/AdminEditEventForm.vue'
import AdminHosts from '@/components/admin/AdminHosts.vue'
import AdminInvites from '@/components/admin/AdminInvites.vue'
import AdminSchedule from '@/components/admin/AdminSchedule.vue'
import AdminFAQs from '@/components/admin/AdminFAQs.vue'
import AdminGifts from '@/components/admin/AdminGifts.vue'
import AdminGallery from '@/components/admin/AdminGallery.vue'
import AdminCustomContent from '@/components/admin/AdminCustomContent.vue'
import AdminComponentOrder from '@/components/admin/AdminComponentOrder.vue'

const authStore = useAuthStore()
const adminStore = useAdminStore()
const toast = useToast()

const { isAdmin } = storeToRefs(authStore)
const { events, loading: eventsLoading, error: eventsError } = storeToRefs(adminStore)

const showCreateForm = ref(false)
const editingEventId = ref<string | null>(null)
const collapsedSections = ref<Record<string, Record<string, boolean>>>({})

/** Checks if a section is collapsed for a given event. */
function isSectionCollapsed(eventId: string, section: string): boolean {
  return !(collapsedSections.value[eventId]?.[section] ?? false)
}

/** Toggles the collapsed state of a section for a given event. */
function toggleSection(eventId: string, section: string) {
  if (!collapsedSections.value[eventId]) {
    collapsedSections.value[eventId] = {}
  }
  collapsedSections.value[eventId][section] = !collapsedSections.value[eventId][section]
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

/** Deletes an event after confirmation. */
async function handleDeleteEvent(eventId: string) {
  if (!(await toast.confirm('Delete this event? This will remove all related data (hosts, invites, guests, etc.) and cannot be undone.'))) return

  try {
    await adminStore.deleteEvent(eventId)
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : 'Failed to delete event')
  }
}

/** Generates the public URL for an event. */
function getEventUrl(event: AdminEvent): string {
  return `/${event.type}/${event.slug}`
}
</script>

<template>
  <div class="admin-view">
    <div class="min-h-screen pt-24 pb-8 px-4">
      <div class="max-w-6xl mx-auto">
        <div class="flex items-center justify-between mb-8">
          <h1 class="text-2xl font-bold text-primary">Admin Dashboard</h1>
          <button
            @click="showCreateForm = true"
            class="bg-accent text-black font-medium px-4 py-2 rounded-lg hover:bg-accent-dark transition-colors"
          >
            + Create Event
          </button>
        </div>

        <!-- Create Event Form -->
        <AdminCreateEventForm
          v-if="showCreateForm"
          @created="showCreateForm = false"
          @cancel="showCreateForm = false"
        />

        <!-- Events Loading -->
        <div v-if="eventsLoading" class="flex items-center justify-center py-12">
          <LoadingSpinner size="md" />
        </div>

        <!-- Events Error -->
        <div v-else-if="eventsError" class="text-center py-12">
          <p class="text-red-600 mb-4">{{ eventsError }}</p>
          <button @click="adminStore.fetchEvents()" class="text-accent hover:underline">Try again</button>
        </div>

        <!-- Events List -->
        <div v-else-if="events.length === 0" class="bg-white rounded-lg shadow-sm p-8 text-center">
          <p class="text-gray-500">No events yet. Create your first event!</p>
        </div>

        <div v-else class="space-y-4">
          <div
            v-for="event in events"
            :key="event.id"
            class="bg-white rounded-lg shadow-sm p-6"
          >
            <!-- View Mode -->
            <template v-if="editingEventId !== event.id">
              <div class="flex items-start justify-between mb-4">
                <div>
                  <div class="flex items-center gap-2 mb-1">
                    <span class="inline-block px-2 py-0.5 bg-accent text-black text-xs font-medium rounded capitalize">
                      {{ event.type }}
                    </span>
                    <span v-if="event.isPublic" class="text-xs text-gray-500">Public</span>
                    <span v-else class="text-xs text-gray-500">Private</span>
                  </div>
                  <h3 class="text-lg font-semibold text-primary">{{ event.title }}</h3>
                  <p v-if="event.description" class="text-sm text-gray-600 mt-1">{{ event.description }}</p>
                  <p class="text-sm text-gray-500 mt-1">
                    {{ formatDate(event.startsAt) }} · {{ event.location || 'No location' }}
                  </p>
                </div>
                <div class="flex items-center gap-3 shrink-0">
                  <a
                    :href="getEventUrl(event)"
                    target="_blank"
                    class="text-sm text-accent hover:underline"
                  >
                    {{ getEventUrl(event) }} →
                  </a>
                  <button
                    @click="editingEventId = event.id"
                    class="text-sm text-primary hover:underline"
                  >
                    Edit
                  </button>
                  <button
                    @click="handleDeleteEvent(event.id)"
                    class="text-sm text-red-500 hover:text-red-700"
                  >
                    Delete
                  </button>
                </div>
              </div>
            </template>

            <!-- Edit Mode -->
            <AdminEditEventForm
              v-else
              :event="event"
              @saved="editingEventId = null"
              @cancel="editingEventId = null"
            />

            <!-- Hosts Section -->
            <AdminHosts
              :event-id="event.id"
              :hosts="event.hosts"
              :collapsed="isSectionCollapsed(event.id, 'hosts')"
              @toggle="toggleSection(event.id, 'hosts')"
            />

            <!-- Invites & Guests Section -->
            <AdminInvites
              :event-id="event.id"
              :collapsed="isSectionCollapsed(event.id, 'invites')"
              @toggle="toggleSection(event.id, 'invites')"
            />

            <!-- Schedule Section -->
            <AdminSchedule
              :event-id="event.id"
              :collapsed="isSectionCollapsed(event.id, 'schedule')"
              @toggle="toggleSection(event.id, 'schedule')"
            />

            <!-- FAQs Section -->
            <AdminFAQs
              :event-id="event.id"
              :collapsed="isSectionCollapsed(event.id, 'faqs')"
              @toggle="toggleSection(event.id, 'faqs')"
            />

            <!-- Gifts Section -->
            <AdminGifts
              :event-id="event.id"
              :collapsed="isSectionCollapsed(event.id, 'gifts')"
              @toggle="toggleSection(event.id, 'gifts')"
            />

            <!-- Gallery Section -->
            <AdminGallery
              :event-id="event.id"
              :collapsed="isSectionCollapsed(event.id, 'gallery')"
              @toggle="toggleSection(event.id, 'gallery')"
            />

            <!-- Custom Content Section -->
            <AdminCustomContent
              :event-id="event.id"
              :collapsed="isSectionCollapsed(event.id, 'customContent')"
              @toggle="toggleSection(event.id, 'customContent')"
            />

            <!-- Component Order Section -->
            <AdminComponentOrder
              :event-id="event.id"
              :collapsed="isSectionCollapsed(event.id, 'componentOrder')"
              @toggle="toggleSection(event.id, 'componentOrder')"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
