<script setup lang="ts">
import { ref, watch } from 'vue'
import { useEventStore } from '@/stores/event'
import type { EventType, Invite, Guest } from '@/types'

const eventStore = useEventStore()

const props = defineProps<{
  invite: Invite
  type: EventType
  slug: string
  inviteCode: string
}>()

const emit = defineEmits<{
  rsvpUpdated: []
}>()

interface GuestState {
  status: 'yes' | 'no' | null
  message: string
  submitting: boolean
  submitted: boolean
  error: string
}

// track RSVP state per guest
const guestStates = ref<Record<string, GuestState>>({})

/** Initializes RSVP state for a guest if not already tracked. */
function initGuest(guest: Guest) {
  if (guestStates.value[guest.id]) return
  guestStates.value[guest.id] = {
    status: guest.rsvpStatus === 'pending' ? null : (guest.rsvpStatus as 'yes' | 'no'),
    message: guest.rsvpMessage || '',
    submitting: false,
    submitted: guest.rsvpStatus !== 'pending',
    error: '',
  }
}

/** Returns the RSVP state for a guest, lazily initializing if missing. */
function getState(guestId: string): GuestState {
  if (!guestStates.value[guestId]) {
    guestStates.value[guestId] = {
      status: null, message: '', submitting: false, submitted: false, error: '',
    }
  }
  return guestStates.value[guestId]
}

props.invite.guests.forEach(initGuest)

watch(() => props.invite.guests, (guests) => guests.forEach(initGuest))

/**
 * Selects RSVP status for a guest.
 */
function selectStatus(guestId: string, status: 'yes' | 'no') {
  const state = getState(guestId)
  state.status = status
  state.error = ''
}

/**
 * Submits RSVP for a guest.
 */
async function handleSubmit(guest: Guest) {
  const state = getState(guest.id)

  if (!state.status) {
    state.error = 'Please select Yes or No'
    return
  }

  state.submitting = true
  state.error = ''

  try {
    await eventStore.submitRSVP(props.type, props.slug, {
      inviteCode: props.inviteCode,
      guestId: guest.id,
      status: state.status,
      message: state.message || undefined,
    })
    state.submitted = true
    emit('rsvpUpdated')
  } catch (e: unknown) {
    state.error = e instanceof Error ? e.message : 'Failed to submit RSVP'
  } finally {
    state.submitting = false
  }
}
</script>

<template>
  <section class="invite-rsvp py-16 px-4">
    <div class="max-w-3xl mx-auto">
      <h2 class="text-2xl font-bold text-heading text-center mb-2">Your Invitation</h2>
      <p v-if="invite.label" class="text-gray-500 text-center mb-8">
        {{ invite.label }}
      </p>

      <div class="space-y-6">
        <div
          v-for="guest in invite.guests"
          :key="guest.id"
          class="bg-white/80 backdrop-blur border border-muted/50 shadow-sm rounded-xl p-6"
        >
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold text-heading">{{ guest.displayName }}</h3>
            <span
              v-if="getState(guest.id).submitted"
              :class="[
                'px-3 py-1 rounded-full text-sm font-medium',
                getState(guest.id).status === 'yes'
                  ? 'bg-green-100 text-green-800'
                  : 'bg-red-100 text-red-800',
              ]"
            >
              {{ getState(guest.id).status === 'yes' ? 'Attending' : 'Not Attending' }}
            </span>
          </div>

          <template v-if="!getState(guest.id).submitted">
            <!-- Status buttons -->
            <div class="flex gap-4 mb-4">
              <button
                @click="selectStatus(guest.id, 'yes')"
                :class="[
                  'flex-1 py-3 rounded-lg font-semibold transition-colors',
                  getState(guest.id).status === 'yes'
                    ? 'bg-green-500 text-white'
                    : 'bg-surface text-gray-700 hover:bg-muted',
                ]"
                :disabled="getState(guest.id).submitting"
              >
                Yes, I'll be there
              </button>
              <button
                @click="selectStatus(guest.id, 'no')"
                :class="[
                  'flex-1 py-3 rounded-lg font-semibold transition-colors',
                  getState(guest.id).status === 'no'
                    ? 'bg-red-500 text-white'
                    : 'bg-surface text-gray-700 hover:bg-muted',
                ]"
                :disabled="getState(guest.id).submitting"
              >
                Sorry, can't make it
              </button>
            </div>

            <!-- Message input -->
            <div v-if="getState(guest.id).status" class="mb-4">
              <textarea
                v-model="getState(guest.id).message"
                placeholder="Leave a message (optional)"
                rows="2"
                class="w-full px-4 py-3 border border-muted rounded-lg focus:outline-none focus:ring-2 focus:ring-accent resize-none"
                :disabled="getState(guest.id).submitting"
              ></textarea>
            </div>

            <!-- Error -->
            <p v-if="getState(guest.id).error" class="text-red-600 text-sm mb-4">
              {{ getState(guest.id).error }}
            </p>

            <!-- Submit button -->
            <button
              v-if="getState(guest.id).status"
              @click="handleSubmit(guest)"
              :disabled="getState(guest.id).submitting"
              class="w-full bg-accent text-black font-semibold py-3 rounded-lg hover:bg-accent-dark transition-colors disabled:opacity-50"
            >
              {{ getState(guest.id).submitting ? 'Submitting...' : 'Confirm RSVP' }}
            </button>
          </template>

          <!-- Show message if already submitted -->
          <p v-else-if="getState(guest.id).message" class="text-gray-600 italic">
            "{{ getState(guest.id).message }}"
          </p>
        </div>
      </div>
    </div>
  </section>
</template>
