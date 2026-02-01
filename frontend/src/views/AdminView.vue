<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { supabase, getUser, onAuthStateChange } from '@/services/supabase'
import { adminListEvents, adminCreateEvent, adminAddHost, adminDeleteHost } from '@/services/api'
import type { AdminEvent } from '@/types'

const loading = ref(true)
const isAdmin = ref(false)
const error = ref('')
const email = ref('')
const password = ref('')
const loginLoading = ref(false)

// events data
const events = ref<AdminEvent[]>([])
const eventsLoading = ref(false)
const eventsError = ref('')

// create event form
const showCreateForm = ref(false)
const createForm = ref({
  type: 'wedding',
  slug: '',
  title: '',
  isPublic: false,
  startsAt: '',
  location: '',
})
const createLoading = ref(false)
const createError = ref('')

// add host form
const selectedEventId = ref<string | null>(null)
const hostForm = ref({ email: '', displayName: '' })
const hostLoading = ref(false)
const hostError = ref('')

let authSubscription: { unsubscribe: () => void } | null = null

onMounted(async () => {
  await checkAdmin()

  const { data } = onAuthStateChange(async (_event, session) => {
    if (!session) {
      isAdmin.value = false
    }
  })
  authSubscription = data.subscription
})

onUnmounted(() => {
  authSubscription?.unsubscribe()
})

/**
 * Checks if user is logged in and is an admin.
 */
async function checkAdmin() {
  loading.value = true

  try {
    const user = await getUser()
    if (!user) {
      isAdmin.value = false
      loading.value = false
      return
    }

    const { data, error: fetchError } = await supabase
      .from('admins')
      .select('email')
      .eq('email', user.email)
      .single()

    if (fetchError || !data) {
      isAdmin.value = false
    } else {
      isAdmin.value = true
      await loadEvents()
    }
  } catch (e: any) {
    console.error('Admin check error:', e)
    isAdmin.value = false
  } finally {
    loading.value = false
  }
}

/**
 * Handles admin login.
 */
async function handleLogin() {
  if (!email.value || !password.value) {
    error.value = 'Please enter email and password'
    return
  }

  loginLoading.value = true
  error.value = ''

  try {
    await supabase.auth.signInWithPassword({
      email: email.value,
      password: password.value,
    })
    await checkAdmin()
  } catch (e: any) {
    error.value = e.message || 'Login failed'
  } finally {
    loginLoading.value = false
  }
}

/**
 * Loads all events.
 */
async function loadEvents() {
  eventsLoading.value = true
  eventsError.value = ''

  try {
    const data = await adminListEvents()
    events.value = data.events
  } catch (e: any) {
    eventsError.value = e.message || 'Failed to load events'
  } finally {
    eventsLoading.value = false
  }
}

/**
 * Converts datetime-local value to ISO 8601 format.
 */
function formatDateTimeForAPI(value: string): string | undefined {
  if (!value) return undefined
  // datetime-local gives "2024-06-15T14:00", need "2024-06-15T14:00:00Z"
  const date = new Date(value)
  return date.toISOString()
}

/**
 * Creates a new event.
 */
async function handleCreateEvent() {
  createLoading.value = true
  createError.value = ''

  try {
    const newEvent = await adminCreateEvent({
      type: createForm.value.type,
      slug: createForm.value.slug,
      title: createForm.value.title,
      isPublic: createForm.value.isPublic,
      startsAt: formatDateTimeForAPI(createForm.value.startsAt),
      location: createForm.value.location || undefined,
    })
    events.value.unshift(newEvent)
    showCreateForm.value = false
    resetCreateForm()
  } catch (e: any) {
    createError.value = e.message || 'Failed to create event'
  } finally {
    createLoading.value = false
  }
}

/**
 * Resets the create event form.
 */
function resetCreateForm() {
  createForm.value = {
    type: 'wedding',
    slug: '',
    title: '',
    isPublic: false,
    startsAt: '',
    location: '',
  }
  createError.value = ''
}

/**
 * Adds a host to an event.
 */
async function handleAddHost() {
  if (!selectedEventId.value) return

  hostLoading.value = true
  hostError.value = ''

  try {
    const newHost = await adminAddHost(selectedEventId.value, {
      email: hostForm.value.email,
      displayName: hostForm.value.displayName,
    })

    // update local state
    const event = events.value.find((e) => e.id === selectedEventId.value)
    if (event) {
      event.hosts.push(newHost)
    }

    hostForm.value = { email: '', displayName: '' }
  } catch (e: any) {
    hostError.value = e.message || 'Failed to add host'
  } finally {
    hostLoading.value = false
  }
}

/**
 * Removes a host from an event.
 */
async function handleDeleteHost(eventId: string, hostId: string) {
  if (!confirm('Remove this host?')) return

  try {
    await adminDeleteHost(eventId, hostId)
    const event = events.value.find((e) => e.id === eventId)
    if (event) {
      event.hosts = event.hosts.filter((h) => h.id !== hostId)
    }
  } catch (e: any) {
    alert(e.message || 'Failed to remove host')
  }
}

/**
 * Formats date for display.
 */
function formatDate(dateStr?: string): string {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

/**
 * Generates event URL.
 */
function getEventUrl(event: AdminEvent): string {
  return `/${event.type}/${event.slug}`
}
</script>

<template>
  <div class="admin-view min-h-screen pt-24 pb-8 px-4">
    <div class="max-w-6xl mx-auto">
      <!-- Loading -->
      <div v-if="loading" class="flex items-center justify-center py-20">
        <div class="animate-spin rounded-full h-12 w-12 border-4 border-[#fca311] border-t-transparent"></div>
      </div>

      <!-- Login Form -->
      <div v-else-if="!isAdmin" class="max-w-md mx-auto">
        <div class="text-center mb-8">
          <h1 class="text-3xl font-bold text-[#14213d]">Admin Login</h1>
          <p class="text-gray-600 mt-2">Sign in with your admin account</p>
        </div>

        <form @submit.prevent="handleLogin" class="bg-white rounded-lg shadow-sm p-8 space-y-6">
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
              Email
            </label>
            <input
              id="email"
              v-model="email"
              type="email"
              required
              class="w-full px-4 py-3 border border-[#e5e5e5] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#fca311]"
              :disabled="loginLoading"
            />
          </div>

          <div>
            <label for="password" class="block text-sm font-medium text-gray-700 mb-2">
              Password
            </label>
            <input
              id="password"
              v-model="password"
              type="password"
              required
              class="w-full px-4 py-3 border border-[#e5e5e5] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#fca311]"
              :disabled="loginLoading"
            />
          </div>

          <p v-if="error" class="text-red-600 text-sm">{{ error }}</p>

          <button
            type="submit"
            :disabled="loginLoading"
            class="w-full bg-[#14213d] text-white font-semibold py-3 rounded-lg hover:bg-[#1a2a4d] transition-colors disabled:opacity-50"
          >
            {{ loginLoading ? 'Signing in...' : 'Sign In' }}
          </button>
        </form>
      </div>

      <!-- Admin Dashboard -->
      <template v-else>
        <div class="flex items-center justify-between mb-8">
          <h1 class="text-2xl font-bold text-[#14213d]">Admin Dashboard</h1>
          <button
            @click="showCreateForm = true"
            class="bg-[#fca311] text-black font-medium px-4 py-2 rounded-lg hover:bg-[#e5930f] transition-colors"
          >
            + Create Event
          </button>
        </div>

        <!-- Create Event Form -->
        <div v-if="showCreateForm" class="bg-white rounded-lg shadow-sm p-6 mb-6">
          <h2 class="text-lg font-semibold text-[#14213d] mb-4">Create New Event</h2>
          <form @submit.prevent="handleCreateEvent" class="space-y-4">
            <div class="grid md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Type</label>
                <select
                  v-model="createForm.type"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#fca311] focus:outline-none"
                >
                  <option value="wedding">Wedding</option>
                  <option value="birthday">Birthday</option>
                  <option value="party">Party</option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Slug (URL path)</label>
                <input
                  v-model="createForm.slug"
                  type="text"
                  placeholder="john-jane-2024"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#fca311] focus:outline-none"
                />
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Title</label>
              <input
                v-model="createForm.title"
                type="text"
                placeholder="John & Jane's Wedding"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#fca311] focus:outline-none"
              />
            </div>
            <div class="grid md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Date & Time</label>
                <input
                  v-model="createForm.startsAt"
                  type="datetime-local"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#fca311] focus:outline-none"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Location</label>
                <input
                  v-model="createForm.location"
                  type="text"
                  placeholder="The Garden Venue"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#fca311] focus:outline-none"
                />
              </div>
            </div>
            <div class="flex items-center gap-2">
              <input
                v-model="createForm.isPublic"
                type="checkbox"
                id="isPublic"
                class="rounded"
              />
              <label for="isPublic" class="text-sm text-gray-700">Public event (no invite code required)</label>
            </div>
            <p v-if="createError" class="text-red-600 text-sm">{{ createError }}</p>
            <div class="flex gap-3">
              <button
                type="submit"
                :disabled="createLoading"
                class="bg-[#14213d] text-white font-medium px-4 py-2 rounded-lg hover:bg-[#1a2a4d] transition-colors disabled:opacity-50"
              >
                {{ createLoading ? 'Creating...' : 'Create Event' }}
              </button>
              <button
                type="button"
                @click="showCreateForm = false; resetCreateForm()"
                class="text-gray-600 hover:text-gray-800 px-4 py-2"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>

        <!-- Events Loading -->
        <div v-if="eventsLoading" class="flex items-center justify-center py-12">
          <div class="animate-spin rounded-full h-8 w-8 border-4 border-[#fca311] border-t-transparent"></div>
        </div>

        <!-- Events Error -->
        <div v-else-if="eventsError" class="text-center py-12">
          <p class="text-red-600 mb-4">{{ eventsError }}</p>
          <button @click="loadEvents" class="text-[#fca311] hover:underline">Try again</button>
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
            <div class="flex items-start justify-between mb-4">
              <div>
                <div class="flex items-center gap-2 mb-1">
                  <span class="inline-block px-2 py-0.5 bg-[#fca311] text-black text-xs font-medium rounded capitalize">
                    {{ event.type }}
                  </span>
                  <span v-if="event.isPublic" class="text-xs text-gray-500">Public</span>
                  <span v-else class="text-xs text-gray-500">Private</span>
                </div>
                <h3 class="text-lg font-semibold text-[#14213d]">{{ event.title }}</h3>
                <p class="text-sm text-gray-500">
                  {{ formatDate(event.startsAt) }} · {{ event.location || 'No location' }}
                </p>
              </div>
              <a
                :href="getEventUrl(event)"
                target="_blank"
                class="text-sm text-[#fca311] hover:underline"
              >
                {{ getEventUrl(event) }} →
              </a>
            </div>

            <!-- Hosts Section -->
            <div class="border-t border-gray-100 pt-4">
              <div class="flex items-center justify-between mb-3">
                <h4 class="text-sm font-medium text-gray-700">Hosts</h4>
                <button
                  @click="selectedEventId = selectedEventId === event.id ? null : event.id"
                  class="text-sm text-[#14213d] hover:underline"
                >
                  {{ selectedEventId === event.id ? 'Cancel' : '+ Add Host' }}
                </button>
              </div>

              <!-- Host List -->
              <div v-if="event.hosts.length > 0" class="space-y-2 mb-3">
                <div
                  v-for="host in event.hosts"
                  :key="host.id"
                  class="flex items-center justify-between bg-gray-50 px-3 py-2 rounded"
                >
                  <div>
                    <span class="font-medium">{{ host.displayName }}</span>
                    <span class="text-sm text-gray-500 ml-2">{{ host.contactEmail }}</span>
                    <span
                      v-if="host.authUserId"
                      class="text-xs text-green-600 ml-2"
                      title="Account linked"
                    >✓ Linked</span>
                    <span
                      v-else
                      class="text-xs text-orange-500 ml-2"
                      title="No auth account yet"
                    >⚠ Not linked</span>
                  </div>
                  <button
                    @click="handleDeleteHost(event.id, host.id)"
                    class="text-red-500 hover:text-red-700 text-sm"
                  >
                    Remove
                  </button>
                </div>
              </div>
              <p v-else class="text-sm text-gray-400 mb-3">No hosts added yet</p>

              <!-- Add Host Form -->
              <div v-if="selectedEventId === event.id" class="bg-gray-50 p-4 rounded-lg">
                <form @submit.prevent="handleAddHost" class="space-y-3">
                  <div class="grid md:grid-cols-2 gap-3">
                    <input
                      v-model="hostForm.email"
                      type="email"
                      placeholder="Host email"
                      required
                      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#fca311] focus:outline-none"
                    />
                    <input
                      v-model="hostForm.displayName"
                      type="text"
                      placeholder="Display name"
                      required
                      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#fca311] focus:outline-none"
                    />
                  </div>
                  <p v-if="hostError" class="text-red-600 text-sm">{{ hostError }}</p>
                  <button
                    type="submit"
                    :disabled="hostLoading"
                    class="bg-[#14213d] text-white font-medium px-4 py-2 rounded-lg hover:bg-[#1a2a4d] transition-colors disabled:opacity-50"
                  >
                    {{ hostLoading ? 'Adding...' : 'Add Host' }}
                  </button>
                </form>
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped></style>
