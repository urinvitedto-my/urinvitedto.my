<script setup lang="ts">
import { ref } from 'vue'
import { formatDate } from '@/utils/date'
import { useToast } from '@/composables/useToast'
import { useAdminStore } from '@/stores/admin'
import type { AdminEvent } from '@/types'

import AdminEditEventForm from './AdminEditEventForm.vue'
import AdminGuestOverview from './AdminGuestOverview.vue'
import AdminHosts from './AdminHosts.vue'
import AdminInvites from './AdminInvites.vue'
import AdminSchedule from './AdminSchedule.vue'
import AdminFAQs from './AdminFAQs.vue'
import AdminGifts from './AdminGifts.vue'
import AdminGallery from './AdminGallery.vue'
import AdminCustomContent from './AdminCustomContent.vue'
import AdminComponentOrder from './AdminComponentOrder.vue'

const props = defineProps<{ event: AdminEvent }>()
const emit = defineEmits<{ deleted: [] }>()

const adminStore = useAdminStore()
const toast = useToast()

const activeTab = ref('guests')
const isEditing = ref(false)

const tabs = [
  { key: 'guests', label: 'Guests' },
  { key: 'hosts', label: 'Hosts' },
  { key: 'invites', label: 'Invites' },
  { key: 'schedule', label: 'Schedule' },
  { key: 'faqs', label: 'FAQs' },
  { key: 'gifts', label: 'Gifts' },
  { key: 'gallery', label: 'Gallery' },
  { key: 'customContent', label: 'Content' },
  { key: 'componentOrder', label: 'Order' },
]

/** Generates the public URL for an event. */
function getEventUrl(event: AdminEvent): string {
  return `/${event.type}/${event.slug}`
}

/** Deletes an event after confirmation. */
async function handleDelete() {
  if (
    !(await toast.confirm(
      'Delete this event? This will remove all related data (hosts, invites, guests, etc.) and cannot be undone.',
    ))
  )
    return
  try {
    await adminStore.deleteEvent(props.event.id)
    emit('deleted')
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : 'Failed to delete event')
  }
}

function noop() {}
</script>

<template>
  <div>
    <!-- Header -->
    <template v-if="!isEditing">
      <div
        class="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-3 mb-6"
      >
        <div class="min-w-0">
          <div class="flex items-center gap-2 mb-1 flex-wrap">
            <span
              class="inline-block px-2 py-0.5 bg-accent text-black text-xs font-medium rounded capitalize"
            >
              {{ event.type }}
            </span>
            <span class="text-xs text-gray-500">{{
              event.isPublic ? 'Public' : 'Private'
            }}</span>
          </div>
          <h2 class="text-xl font-bold text-primary">{{ event.title }}</h2>
          <p v-if="event.description" class="text-sm text-gray-600 mt-1">
            {{ event.description }}
          </p>
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
            @click="isEditing = true"
            class="text-sm text-primary hover:underline"
          >
            Edit
          </button>
          <button @click="handleDelete" class="text-sm text-red-500 hover:text-red-700">
            Delete
          </button>
        </div>
      </div>
    </template>

    <!-- Edit Form -->
    <AdminEditEventForm
      v-else
      :event="event"
      @saved="isEditing = false"
      @cancel="isEditing = false"
    />

    <!-- Tab Bar -->
    <div class="border-b border-gray-200 mb-4 -mx-1 overflow-x-auto">
      <div class="flex min-w-max px-1">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          @click="activeTab = tab.key"
          class="px-3 py-2 text-sm font-medium whitespace-nowrap border-b-2 transition-colors"
          :class="
            activeTab === tab.key
              ? 'border-accent text-primary'
              : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
          "
        >
          {{ tab.label }}
        </button>
      </div>
    </div>

    <!-- Tab Content -->
    <div>
      <AdminGuestOverview v-if="activeTab === 'guests'" :event-id="event.id" />

      <AdminHosts
        v-else-if="activeTab === 'hosts'"
        :event-id="event.id"
        :hosts="event.hosts"
        :collapsed="false"
        @toggle="noop"
      />

      <AdminInvites
        v-else-if="activeTab === 'invites'"
        :event-id="event.id"
        :collapsed="false"
        @toggle="noop"
      />

      <AdminSchedule
        v-else-if="activeTab === 'schedule'"
        :event-id="event.id"
        :collapsed="false"
        @toggle="noop"
      />

      <AdminFAQs
        v-else-if="activeTab === 'faqs'"
        :event-id="event.id"
        :collapsed="false"
        @toggle="noop"
      />

      <AdminGifts
        v-else-if="activeTab === 'gifts'"
        :event-id="event.id"
        :collapsed="false"
        @toggle="noop"
      />

      <AdminGallery
        v-else-if="activeTab === 'gallery'"
        :event-id="event.id"
        :collapsed="false"
        @toggle="noop"
      />

      <AdminCustomContent
        v-else-if="activeTab === 'customContent'"
        :event-id="event.id"
        :collapsed="false"
        @toggle="noop"
      />

      <AdminComponentOrder
        v-else-if="activeTab === 'componentOrder'"
        :event-id="event.id"
        :collapsed="false"
        @toggle="noop"
      />
    </div>
  </div>
</template>
