<script setup lang="ts">
import { ref } from 'vue'
import { submitRSVP } from '@/services/api'
import type { Invite, Guest } from '@/types'

const props = defineProps<{
  invite: Invite
  type: string
  slug: string
  inviteCode: string
}>()

const emit = defineEmits<{
  rsvpUpdated: []
}>()

// track RSVP state per guest
const guestStates = ref<Record<string, {
  status: 'yes' | 'no' | null
  message: string
  submitting: boolean
  submitted: boolean
  error: string
}>>({})

// initialize state for each guest
props.invite.guests.forEach((guest) => {
  guestStates.value[guest.id] = {
    status: guest.rsvpStatus === 'pending' ? null : (guest.rsvpStatus as 'yes' | 'no'),
    message: guest.rsvpMessage || '',
    submitting: false,
    submitted: guest.rsvpStatus !== 'pending',
    error: '',
  }
})

/**
 * Selects RSVP status for a guest.
 */
function selectStatus(guestId: string, status: 'yes' | 'no') {
  guestStates.value[guestId].status = status
  guestStates.value[guestId].error = ''
}

/**
 * Submits RSVP for a guest.
 */
async function handleSubmit(guest: Guest) {
  const state = guestStates.value[guest.id]

  if (!state.status) {
    state.error = 'Please select Yes or No'
    return
  }

  state.submitting = true
  state.error = ''

  try {
    await submitRSVP(props.type, props.slug, {
      inviteCode: props.inviteCode,
      guestId: guest.id,
      status: state.status,
      message: state.message || undefined,
    })
    state.submitted = true
    emit('rsvpUpdated')
  } catch (e: any) {
    state.error = e.message || 'Failed to submit RSVP'
  } finally {
    state.submitting = false
  }
}
</script>

<template>
  <section class="invite-rsvp bg-[#14213d] py-12 px-4">
    <div class="max-w-3xl mx-auto">
      <h2 class="text-2xl font-bold text-white text-center mb-2">Your Invitation</h2>
      <p v-if="invite.label" class="text-[#e5e5e5] text-center mb-8">
        {{ invite.label }}
      </p>

      <div class="space-y-6">
        <div
          v-for="guest in invite.guests"
          :key="guest.id"
          class="bg-white rounded-lg p-6"
        >
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold text-[#14213d]">{{ guest.displayName }}</h3>
            <span
              v-if="guestStates[guest.id].submitted"
              :class="[
                'px-3 py-1 rounded-full text-sm font-medium',
                guestStates[guest.id].status === 'yes'
                  ? 'bg-green-100 text-green-800'
                  : 'bg-red-100 text-red-800',
              ]"
            >
              {{ guestStates[guest.id].status === 'yes' ? 'Attending' : 'Not Attending' }}
            </span>
          </div>

          <template v-if="!guestStates[guest.id].submitted">
            <!-- Status buttons -->
            <div class="flex gap-4 mb-4">
              <button
                @click="selectStatus(guest.id, 'yes')"
                :class="[
                  'flex-1 py-3 rounded-lg font-semibold transition-colors',
                  guestStates[guest.id].status === 'yes'
                    ? 'bg-green-500 text-white'
                    : 'bg-[#ececec] text-gray-700 hover:bg-[#e5e5e5]',
                ]"
                :disabled="guestStates[guest.id].submitting"
              >
                Yes, I'll be there
              </button>
              <button
                @click="selectStatus(guest.id, 'no')"
                :class="[
                  'flex-1 py-3 rounded-lg font-semibold transition-colors',
                  guestStates[guest.id].status === 'no'
                    ? 'bg-red-500 text-white'
                    : 'bg-[#ececec] text-gray-700 hover:bg-[#e5e5e5]',
                ]"
                :disabled="guestStates[guest.id].submitting"
              >
                Sorry, can't make it
              </button>
            </div>

            <!-- Message input -->
            <div v-if="guestStates[guest.id].status" class="mb-4">
              <textarea
                v-model="guestStates[guest.id].message"
                placeholder="Leave a message (optional)"
                rows="2"
                class="w-full px-4 py-3 border border-[#e5e5e5] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#fca311] resize-none"
                :disabled="guestStates[guest.id].submitting"
              ></textarea>
            </div>

            <!-- Error -->
            <p v-if="guestStates[guest.id].error" class="text-red-600 text-sm mb-4">
              {{ guestStates[guest.id].error }}
            </p>

            <!-- Submit button -->
            <button
              v-if="guestStates[guest.id].status"
              @click="handleSubmit(guest)"
              :disabled="guestStates[guest.id].submitting"
              class="w-full bg-[#fca311] text-black font-semibold py-3 rounded-lg hover:bg-[#e5930f] transition-colors disabled:opacity-50"
            >
              {{ guestStates[guest.id].submitting ? 'Submitting...' : 'Confirm RSVP' }}
            </button>
          </template>

          <!-- Show message if already submitted -->
          <p v-else-if="guestStates[guest.id].message" class="text-gray-600 italic">
            "{{ guestStates[guest.id].message }}"
          </p>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped></style>
