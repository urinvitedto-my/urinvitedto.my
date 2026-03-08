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
    <div class="max-w-md mx-auto">
      <h2 class="text-4xl font-bold text-primary-dark text-center mb-2 font-kaushan text-outline">Your Invitation</h2>
      <p v-if="invite.label" class="text-muted text-xl text-center mb-8">
        {{ invite.label }}
      </p>

      <div class="divide-y divide-muted/50">
        <div
          v-for="guest in invite.guests"
          :key="guest.id"
          class="py-4 first:pt-0"
        >
          <!-- Submitted state -->
          <div v-if="getState(guest.id).submitted" class="flex items-center justify-between">
            <span class="font-medium text-guest-bg">{{ guest.displayName }}</span>
            <span
              :class="[
                'text-xs font-semibold px-4 py-1.5 rounded-full',
                getState(guest.id).status === 'yes'
                  ? 'bg-guest-bg text-green-500'
                  : 'bg-guest-bg text-red-500',
              ]"
            >
              {{ getState(guest.id).status === 'yes' ? 'Attending' : 'Not Attending' }}
            </span>
          </div>
          <p v-if="getState(guest.id).submitted && getState(guest.id).message" class="text-muted text-sm italic mt-1">
            "{{ getState(guest.id).message }}"
          </p>

          <!-- Pending state -->
          <template v-if="!getState(guest.id).submitted">
            <div class="flex items-center justify-between">
              <span class="font-medium text-guest-bg">{{ guest.displayName }}</span>
              <div class="flex gap-2">
                <button
                  @click="selectStatus(guest.id, 'yes')"
                  :class="[
                    'px-4 py-1.5 rounded-full text-xs font-semibold transition-colors',
                    getState(guest.id).status === 'yes'
                      ? 'bg-green-500 text-white'
                      : 'bg-surface text-gray-600 hover:bg-muted',
                  ]"
                  :disabled="getState(guest.id).submitting"
                >
                  Accept
                </button>
                <button
                  @click="selectStatus(guest.id, 'no')"
                  :class="[
                    'px-4 py-1.5 rounded-full text-xs font-semibold transition-colors',
                    getState(guest.id).status === 'no'
                      ? 'bg-red-500 text-white'
                      : 'bg-surface text-gray-600 hover:bg-muted',
                  ]"
                  :disabled="getState(guest.id).submitting"
                >
                  Decline
                </button>
              </div>
            </div>

            <template v-if="getState(guest.id).status">
              <textarea
                v-model="getState(guest.id).message"
                placeholder="Leave a message (optional)"
                rows="2"
                class="w-full mt-3 px-3 py-2 text-sm border border-muted rounded-lg focus:outline-none focus:ring-2 focus:ring-accent resize-none"
                :disabled="getState(guest.id).submitting"
              ></textarea>

              <p v-if="getState(guest.id).error" class="text-red-600 text-sm mt-2">
                {{ getState(guest.id).error }}
              </p>

              <button
                @click="handleSubmit(guest)"
                :disabled="getState(guest.id).submitting"
                class="w-full mt-2 bg-accent text-black text-sm font-semibold py-2 rounded-lg hover:bg-accent-dark transition-colors disabled:opacity-50"
              >
                {{ getState(guest.id).submitting ? 'Submitting...' : 'Confirm RSVP' }}
              </button>
            </template>
          </template>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.text-outline {
  text-shadow:
    -1px -1px 0 rgba(255, 255, 255, 0.4),
     1px -1px 0 rgba(255, 255, 255, 0.4),
    -1px  1px 0 rgba(255, 255, 255, 0.4),
     1px  1px 0 rgba(255, 255, 255, 0.4);
}
</style>
