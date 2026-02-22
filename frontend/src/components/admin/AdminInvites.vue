<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {
  adminListInvites,
  adminCreateInvite,
  adminDeleteInvite,
  adminAddGuest,
  adminUpdateGuest,
  adminDeleteGuest,
} from '@/services/api'
import type { AdminInvite, AdminGuest } from '@/types'

const props = defineProps<{ eventId: string }>()

const invites = ref<AdminInvite[]>([])
const loading = ref(false)
const error = ref('')

// create invite form
const showCreateForm = ref(false)
const createLabel = ref('')
const createLoading = ref(false)

// add guest form (per invite)
const addingGuestInviteId = ref<string | null>(null)
const guestName = ref('')
const guestLoading = ref(false)

// edit guest
const editingGuestId = ref<string | null>(null)
const editGuestForm = ref({ displayName: '', rsvpStatus: 'pending' })
const editGuestLoading = ref(false)

onMounted(() => loadInvites())

/**
 * Loads all invites with guests for this event.
 */
async function loadInvites() {
  loading.value = true
  error.value = ''

  try {
    const data = await adminListInvites(props.eventId)
    invites.value = data.invites
  } catch (e: any) {
    error.value = e.message || 'Failed to load invites'
  } finally {
    loading.value = false
  }
}

/**
 * Creates a new invite with auto-generated code.
 */
async function handleCreateInvite() {
  createLoading.value = true

  try {
    const newInvite = await adminCreateInvite(props.eventId, {
      label: createLabel.value || null,
    })
    invites.value.unshift(newInvite)
    showCreateForm.value = false
    createLabel.value = ''
  } catch (e: any) {
    alert(e.message || 'Failed to create invite')
  } finally {
    createLoading.value = false
  }
}

/**
 * Deletes an invite and all its guests.
 */
async function handleDeleteInvite(inviteId: string) {
  if (!confirm('Delete this invite and all its guests?')) return

  try {
    await adminDeleteInvite(props.eventId, inviteId)
    invites.value = invites.value.filter((i) => i.id !== inviteId)
  } catch (e: any) {
    alert(e.message || 'Failed to delete invite')
  }
}

/**
 * Adds a guest to an invite.
 */
async function handleAddGuest(inviteId: string) {
  if (!guestName.value.trim()) return

  guestLoading.value = true

  try {
    const newGuest = await adminAddGuest(props.eventId, inviteId, {
      displayName: guestName.value.trim(),
    })
    const invite = invites.value.find((i) => i.id === inviteId)
    if (invite) invite.guests.push(newGuest)
    guestName.value = ''
  } catch (e: any) {
    alert(e.message || 'Failed to add guest')
  } finally {
    guestLoading.value = false
  }
}

/**
 * Opens edit mode for a guest.
 */
function startEditGuest(guest: AdminGuest) {
  editingGuestId.value = guest.id
  editGuestForm.value = {
    displayName: guest.displayName,
    rsvpStatus: guest.rsvpStatus,
  }
}

/**
 * Saves guest edits.
 */
async function handleUpdateGuest(inviteId: string) {
  if (!editingGuestId.value) return

  editGuestLoading.value = true

  try {
    const updated = await adminUpdateGuest(props.eventId, editingGuestId.value, {
      displayName: editGuestForm.value.displayName,
      rsvpStatus: editGuestForm.value.rsvpStatus,
    })

    const invite = invites.value.find((i) => i.id === inviteId)
    if (invite) {
      const idx = invite.guests.findIndex((g) => g.id === editingGuestId.value)
      if (idx !== -1) invite.guests[idx] = updated
    }
    editingGuestId.value = null
  } catch (e: any) {
    alert(e.message || 'Failed to update guest')
  } finally {
    editGuestLoading.value = false
  }
}

/**
 * Deletes a guest.
 */
async function handleDeleteGuest(inviteId: string, guestId: string) {
  if (!confirm('Remove this guest?')) return

  try {
    await adminDeleteGuest(props.eventId, guestId)
    const invite = invites.value.find((i) => i.id === inviteId)
    if (invite) {
      invite.guests = invite.guests.filter((g) => g.id !== guestId)
    }
  } catch (e: any) {
    alert(e.message || 'Failed to delete guest')
  }
}

/**
 * Returns a CSS class for RSVP status badges.
 */
function rsvpClass(status: string): string {
  switch (status) {
    case 'yes': return 'bg-green-100 text-green-700'
    case 'no': return 'bg-red-100 text-red-700'
    default: return 'bg-gray-100 text-gray-600'
  }
}
</script>

<template>
  <div class="border-t border-gray-100 pt-4">
    <div class="flex items-center justify-between mb-3">
      <h4 class="text-sm font-medium text-gray-700">
        Invites
        <span v-if="invites.length > 0" class="text-gray-400 font-normal">({{ invites.length }})</span>
      </h4>
      <button
        @click="showCreateForm = !showCreateForm"
        class="text-sm text-[#14213d] hover:underline"
      >
        {{ showCreateForm ? 'Cancel' : '+ Add Invite' }}
      </button>
    </div>

    <!-- Create Invite Form -->
    <div v-if="showCreateForm" class="bg-gray-50 p-4 rounded-lg mb-3">
      <form @submit.prevent="handleCreateInvite" class="flex items-end gap-3">
        <div class="flex-1">
          <label class="block text-sm font-medium text-gray-700 mb-1">Label (optional)</label>
          <input
            v-model="createLabel"
            type="text"
            placeholder='e.g., "Smith Family"'
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#fca311] focus:outline-none"
          />
        </div>
        <button
          type="submit"
          :disabled="createLoading"
          class="bg-[#14213d] text-white font-medium px-4 py-2 rounded-lg hover:bg-[#1a2a4d] transition-colors disabled:opacity-50 shrink-0"
        >
          {{ createLoading ? 'Creating...' : 'Create Invite' }}
        </button>
      </form>
      <p class="text-xs text-gray-500 mt-2">A unique 6-character invite code will be auto-generated.</p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-4">
      <div class="animate-spin rounded-full h-6 w-6 border-2 border-[#fca311] border-t-transparent"></div>
    </div>

    <!-- Error -->
    <p v-else-if="error" class="text-red-600 text-sm mb-3">{{ error }}</p>

    <!-- Empty -->
    <p v-else-if="invites.length === 0" class="text-sm text-gray-400 mb-3">No invites yet</p>

    <!-- Invite List -->
    <div v-else class="space-y-3">
      <div
        v-for="invite in invites"
        :key="invite.id"
        class="bg-gray-50 rounded-lg p-4"
      >
        <!-- Invite Header -->
        <div class="flex items-center justify-between mb-2">
          <div class="flex items-center gap-3">
            <code class="bg-white px-2 py-1 rounded text-sm font-mono font-bold tracking-wider border border-gray-200">
              {{ invite.inviteCode }}
            </code>
            <span v-if="invite.label" class="text-sm text-gray-600">{{ invite.label }}</span>
            <span class="text-xs text-gray-400">
              {{ invite.guests.length }} guest{{ invite.guests.length !== 1 ? 's' : '' }}
            </span>
          </div>
          <button
            @click="handleDeleteInvite(invite.id)"
            class="text-red-500 hover:text-red-700 text-sm"
          >
            Delete
          </button>
        </div>

        <!-- Guest List -->
        <div v-if="invite.guests.length > 0" class="space-y-1 mb-2">
          <div
            v-for="guest in invite.guests"
            :key="guest.id"
            class="flex items-center justify-between bg-white px-3 py-2 rounded border border-gray-100"
          >
            <!-- View mode -->
            <template v-if="editingGuestId !== guest.id">
              <div class="flex items-center gap-2">
                <span class="text-sm font-medium">{{ guest.displayName }}</span>
                <span
                  class="text-xs px-1.5 py-0.5 rounded capitalize"
                  :class="rsvpClass(guest.rsvpStatus)"
                >
                  {{ guest.rsvpStatus }}
                </span>
                <span v-if="guest.rsvpMessage" class="text-xs text-gray-400 italic truncate max-w-[150px]" :title="guest.rsvpMessage">
                  "{{ guest.rsvpMessage }}"
                </span>
              </div>
              <div class="flex items-center gap-2 shrink-0">
                <button @click="startEditGuest(guest)" class="text-xs text-[#14213d] hover:underline">Edit</button>
                <button @click="handleDeleteGuest(invite.id, guest.id)" class="text-xs text-red-500 hover:text-red-700">Remove</button>
              </div>
            </template>

            <!-- Edit mode -->
            <template v-else>
              <form @submit.prevent="handleUpdateGuest(invite.id)" class="flex items-center gap-2 w-full">
                <input
                  v-model="editGuestForm.displayName"
                  type="text"
                  required
                  class="flex-1 px-2 py-1 text-sm border border-gray-300 rounded focus:ring-2 focus:ring-[#fca311] focus:outline-none"
                />
                <select
                  v-model="editGuestForm.rsvpStatus"
                  class="px-2 py-1 text-sm border border-gray-300 rounded focus:ring-2 focus:ring-[#fca311] focus:outline-none"
                >
                  <option value="pending">Pending</option>
                  <option value="yes">Yes</option>
                  <option value="no">No</option>
                </select>
                <button
                  type="submit"
                  :disabled="editGuestLoading"
                  class="text-xs bg-[#14213d] text-white px-2 py-1 rounded hover:bg-[#1a2a4d] disabled:opacity-50"
                >
                  Save
                </button>
                <button
                  type="button"
                  @click="editingGuestId = null"
                  class="text-xs text-gray-500 hover:text-gray-700"
                >
                  Cancel
                </button>
              </form>
            </template>
          </div>
        </div>

        <!-- Add Guest Form -->
        <div class="flex items-center gap-2 mt-2">
          <template v-if="addingGuestInviteId === invite.id">
            <form @submit.prevent="handleAddGuest(invite.id)" class="flex items-center gap-2 w-full">
              <input
                v-model="guestName"
                type="text"
                placeholder="Guest name"
                required
                class="flex-1 px-3 py-1.5 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#fca311] focus:outline-none"
              />
              <button
                type="submit"
                :disabled="guestLoading"
                class="text-sm bg-[#14213d] text-white px-3 py-1.5 rounded-lg hover:bg-[#1a2a4d] disabled:opacity-50 shrink-0"
              >
                {{ guestLoading ? 'Adding...' : 'Add' }}
              </button>
              <button
                type="button"
                @click="addingGuestInviteId = null; guestName = ''"
                class="text-sm text-gray-500 hover:text-gray-700"
              >
                Cancel
              </button>
            </form>
          </template>
          <button
            v-else
            @click="addingGuestInviteId = invite.id"
            class="text-xs text-[#14213d] hover:underline"
          >
            + Add Guest
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
