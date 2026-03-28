<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { ref } from 'vue'
import { useAdminStore } from '@/stores/admin'
import { formatDate } from '@/utils/date'
import LoadingSpinner from '@/components/LoadingSpinner.vue'

const props = defineProps<{ eventId: string }>()

const adminStore = useAdminStore()

const invites = computed(() => adminStore.getInvites(props.eventId))
const loading = computed(() => adminStore.isSubLoading('invites', props.eventId))
const error = computed(() => adminStore.getSubError('invites', props.eventId))

type StatusFilter = 'all' | 'yes' | 'no' | 'pending'
const statusFilter = ref<StatusFilter>('all')

interface FlatGuest {
  id: string
  displayName: string
  rsvpStatus: 'pending' | 'yes' | 'no'
  rsvpMessage?: string
  rsvpAt?: string
  inviteLabel?: string
  inviteCode: string
}

/** Flattens all guests across all invites into a single sorted list. */
const flatGuests = computed<FlatGuest[]>(() => {
  const guests: FlatGuest[] = []
  for (const invite of invites.value) {
    for (const g of invite.guests) {
      guests.push({
        id: g.id,
        displayName: g.displayName,
        rsvpStatus: g.rsvpStatus,
        rsvpMessage: g.rsvpMessage,
        rsvpAt: g.rsvpAt,
        inviteLabel: invite.label ?? undefined,
        inviteCode: invite.inviteCode,
      })
    }
  }
  return guests.sort((a, b) => a.displayName.localeCompare(b.displayName))
})

const filteredGuests = computed(() => {
  if (statusFilter.value === 'all') return flatGuests.value
  return flatGuests.value.filter((g) => g.rsvpStatus === statusFilter.value)
})

const yesCount = computed(() => flatGuests.value.filter((g) => g.rsvpStatus === 'yes').length)
const noCount = computed(() => flatGuests.value.filter((g) => g.rsvpStatus === 'no').length)
const pendingCount = computed(() => flatGuests.value.filter((g) => g.rsvpStatus === 'pending').length)

function rsvpBadgeClass(status: string): string {
  switch (status) {
    case 'yes': return 'bg-green-100 text-green-700'
    case 'no': return 'bg-red-100 text-red-700'
    default: return 'bg-gray-100 text-gray-600'
  }
}

const filterButtons: { key: StatusFilter; label: string; countFn: () => number }[] = [
  { key: 'all', label: 'All', countFn: () => flatGuests.value.length },
  { key: 'yes', label: 'Confirmed', countFn: () => yesCount.value },
  { key: 'no', label: 'Declined', countFn: () => noCount.value },
  { key: 'pending', label: 'Pending', countFn: () => pendingCount.value },
]

onMounted(() => {
  if (invites.value.length === 0) {
    adminStore.fetchInvites(props.eventId)
  }
})
</script>

<template>
  <div>
    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-8">
      <LoadingSpinner size="sm" />
    </div>

    <!-- Error -->
    <p v-else-if="error" class="text-red-600 text-sm py-4">{{ error }}</p>

    <template v-else>
      <!-- Summary -->
      <div class="mb-4">
        <p class="text-sm text-gray-600">
          <span class="font-semibold text-primary">{{ flatGuests.length }}</span> total guests
          <span class="mx-1">—</span>
          <span class="text-green-600 font-medium">{{ yesCount }} confirmed</span>,
          <span class="text-red-500 font-medium">{{ noCount }} declined</span>,
          <span class="text-gray-500 font-medium">{{ pendingCount }} pending</span>
        </p>
      </div>

      <!-- Filter pills -->
      <div class="flex flex-wrap gap-2 mb-4">
        <button
          v-for="btn in filterButtons"
          :key="btn.key"
          @click="statusFilter = btn.key"
          class="px-3 py-1.5 rounded-full text-xs font-medium transition-colors"
          :class="statusFilter === btn.key
            ? 'bg-primary text-white'
            : 'bg-gray-100 text-gray-600 hover:bg-gray-200'"
        >
          {{ btn.label }} ({{ btn.countFn() }})
        </button>
      </div>

      <!-- Empty state -->
      <p v-if="flatGuests.length === 0" class="text-sm text-gray-400 py-4">
        No guests yet. Add invites and guests in the Invites tab.
      </p>

      <p v-else-if="filteredGuests.length === 0" class="text-sm text-gray-400 py-4">
        No guests match this filter.
      </p>

      <!-- Guest cards -->
      <div v-else class="space-y-2">
        <div
          v-for="guest in filteredGuests"
          :key="guest.id"
          class="bg-gray-50 rounded-lg px-4 py-3 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-1"
        >
          <div class="min-w-0 flex-1">
            <div class="flex items-center gap-2 flex-wrap">
              <span class="text-sm font-semibold text-primary">{{ guest.displayName }}</span>
              <span
                class="text-xs px-1.5 py-0.5 rounded capitalize"
                :class="rsvpBadgeClass(guest.rsvpStatus)"
              >
                {{ guest.rsvpStatus }}
              </span>
            </div>
            <p v-if="guest.rsvpMessage" class="text-xs text-gray-500 italic mt-0.5 wrap-break-word">
              "{{ guest.rsvpMessage }}"
            </p>
          </div>

          <div class="flex items-center gap-3 text-xs text-gray-400 shrink-0">
            <span v-if="guest.inviteLabel">{{ guest.inviteLabel }}</span>
            <code class="bg-white px-1.5 py-0.5 rounded border border-gray-200 font-mono text-[10px]">
              {{ guest.inviteCode }}
            </code>
            <span v-if="guest.rsvpAt">{{ formatDate(guest.rsvpAt, true) }}</span>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>
