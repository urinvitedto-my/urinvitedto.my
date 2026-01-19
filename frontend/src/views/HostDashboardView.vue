<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { supabase, getUser, signOut } from '@/services/supabase'

interface EventData {
  id: string
  type: string
  slug: string
  title: string
  starts_at: string
  is_public: boolean
}

interface GuestData {
  id: string
  display_name: string
  rsvp_status: string
  rsvp_message: string | null
  rsvp_at: string | null
}

const router = useRouter()
const loading = ref(true)
const events = ref<EventData[]>([])
const selectedEvent = ref<EventData | null>(null)
const guests = ref<GuestData[]>([])
const showAllGuests = ref(false)
const error = ref('')

onMounted(async () => {
  await checkAuthAndLoadData()
})

/**
 * Checks auth and loads host events.
 */
async function checkAuthAndLoadData() {
  loading.value = true
  error.value = ''

  try {
    const user = await getUser()
    if (!user) {
      router.push('/host/login')
      return
    }

    // fetch events where user is a host
    const { data, error: fetchError } = await supabase
      .from('events')
      .select(`
        id, type, slug, title, starts_at, is_public,
        hosts!inner(auth_user_id)
      `)
      .eq('hosts.auth_user_id', user.id)
      .order('starts_at', { ascending: false })

    if (fetchError) throw fetchError
    events.value = data || []
  } catch (e: any) {
    error.value = e.message || 'Failed to load events'
  } finally {
    loading.value = false
  }
}

/**
 * Selects an event and loads its guests.
 */
async function selectEvent(event: EventData) {
  selectedEvent.value = event
  loading.value = true

  try {
    const { data, error: fetchError } = await supabase
      .from('guests')
      .select('id, display_name, rsvp_status, rsvp_message, rsvp_at')
      .eq('event_id', event.id)
      .order('display_name')

    if (fetchError) throw fetchError
    guests.value = data || []
  } catch (e: any) {
    error.value = e.message || 'Failed to load guests'
  } finally {
    loading.value = false
  }
}

/**
 * Filters guests based on toggle.
 */
function filteredGuests(): GuestData[] {
  if (showAllGuests.value) return guests.value
  return guests.value.filter((g) => g.rsvp_status === 'yes')
}

/**
 * Handles logout.
 */
async function handleLogout() {
  await signOut()
  router.push('/host/login')
}

/**
 * Formats date for display.
 */
function formatDate(dateStr: string | null): string {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: 'numeric',
    minute: '2-digit',
  })
}
</script>

<template>
  <div class="host-dashboard-view min-h-screen py-8 px-4">
    <div class="max-w-6xl mx-auto">
      <!-- Header -->
      <div class="flex items-center justify-between mb-8">
        <h1 class="text-2xl font-bold text-[#14213d]">Host Dashboard</h1>
        <button
          @click="handleLogout"
          class="text-gray-600 hover:text-[#14213d] transition-colors"
        >
          Sign Out
        </button>
      </div>

      <!-- Loading -->
      <div v-if="loading && !selectedEvent" class="flex items-center justify-center py-20">
        <div class="animate-spin rounded-full h-12 w-12 border-4 border-[#fca311] border-t-transparent"></div>
      </div>

      <!-- Error -->
      <div v-else-if="error" class="text-center py-20">
        <p class="text-red-600 mb-4">{{ error }}</p>
        <button @click="checkAuthAndLoadData" class="text-[#fca311] hover:underline">
          Try again
        </button>
      </div>

      <!-- Content -->
      <div v-else class="grid md:grid-cols-3 gap-8">
        <!-- Events List -->
        <div class="bg-white rounded-lg shadow-sm p-6">
          <h2 class="text-lg font-semibold text-[#14213d] mb-4">Your Events</h2>
          <div v-if="events.length === 0" class="text-gray-500 text-center py-8">
            No events found
          </div>
          <ul v-else class="space-y-2">
            <li v-for="event in events" :key="event.id">
              <button
                @click="selectEvent(event)"
                :class="[
                  'w-full text-left px-4 py-3 rounded-lg transition-colors',
                  selectedEvent?.id === event.id
                    ? 'bg-[#fca311] text-black'
                    : 'bg-[#ececec] hover:bg-[#e5e5e5]',
                ]"
              >
                <div class="font-medium">{{ event.title }}</div>
                <div class="text-sm opacity-75">{{ formatDate(event.starts_at) }}</div>
              </button>
            </li>
          </ul>
        </div>

        <!-- Guests List -->
        <div class="md:col-span-2 bg-white rounded-lg shadow-sm p-6">
          <div v-if="!selectedEvent" class="text-gray-500 text-center py-20">
            Select an event to view guests
          </div>
          <template v-else>
            <div class="flex items-center justify-between mb-6">
              <h2 class="text-lg font-semibold text-[#14213d]">
                {{ selectedEvent.title }} - Guests
              </h2>
              <label class="flex items-center gap-2 text-sm">
                <input v-model="showAllGuests" type="checkbox" class="rounded" />
                Show all guests
              </label>
            </div>

            <div v-if="loading" class="flex items-center justify-center py-12">
              <div class="animate-spin rounded-full h-8 w-8 border-4 border-[#fca311] border-t-transparent"></div>
            </div>

            <div v-else-if="filteredGuests().length === 0" class="text-gray-500 text-center py-12">
              {{ showAllGuests ? 'No guests yet' : 'No confirmed guests yet' }}
            </div>

            <div v-else class="overflow-x-auto">
              <table class="w-full">
                <thead>
                  <tr class="border-b border-[#e5e5e5]">
                    <th class="text-left py-3 px-4 font-medium text-gray-600">Name</th>
                    <th class="text-left py-3 px-4 font-medium text-gray-600">Status</th>
                    <th class="text-left py-3 px-4 font-medium text-gray-600">Message</th>
                    <th class="text-left py-3 px-4 font-medium text-gray-600">RSVP Date</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="guest in filteredGuests()"
                    :key="guest.id"
                    class="border-b border-[#ececec]"
                  >
                    <td class="py-3 px-4">{{ guest.display_name }}</td>
                    <td class="py-3 px-4">
                      <span
                        :class="[
                          'inline-block px-2 py-1 rounded text-xs font-medium',
                          guest.rsvp_status === 'yes'
                            ? 'bg-green-100 text-green-800'
                            : guest.rsvp_status === 'no'
                              ? 'bg-red-100 text-red-800'
                              : 'bg-gray-100 text-gray-800',
                        ]"
                      >
                        {{ guest.rsvp_status }}
                      </span>
                    </td>
                    <td class="py-3 px-4 text-gray-600">
                      {{ guest.rsvp_message || '-' }}
                    </td>
                    <td class="py-3 px-4 text-gray-600">
                      {{ formatDate(guest.rsvp_at) }}
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
