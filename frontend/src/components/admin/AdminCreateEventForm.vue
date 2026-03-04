<script setup lang="ts">
import { ref } from 'vue'
import { useAdminStore } from '@/stores/admin'
import { toISOOrUndefined } from '@/utils/date'
import type { EventType } from '@/types'

const emit = defineEmits<{ created: []; cancel: [] }>()

const adminStore = useAdminStore()

const form = ref<{
  type: EventType
  slug: string
  title: string
  isPublic: boolean
  startsAt: string
  location: string
}>({
  type: 'wedding',
  slug: '',
  title: '',
  isPublic: false,
  startsAt: '',
  location: '',
})
const loading = ref(false)
const error = ref('')

/** Creates a new event from the form data. */
async function handleSubmit() {
  loading.value = true
  error.value = ''

  try {
    await adminStore.createEvent({
      type: form.value.type,
      slug: form.value.slug,
      title: form.value.title,
      isPublic: form.value.isPublic,
      startsAt: toISOOrUndefined(form.value.startsAt),
      location: form.value.location || undefined,
    })
    emit('created')
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : 'Failed to create event'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="bg-white rounded-lg shadow-sm p-6 mb-6">
    <h2 class="text-lg font-semibold text-primary mb-4">Create New Event</h2>
    <form @submit.prevent="handleSubmit" class="space-y-4">
      <div class="grid md:grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Type</label>
          <select
            v-model="form.type"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
          >
            <option value="wedding">Wedding</option>
            <option value="birthday">Birthday</option>
            <option value="party">Party</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Slug (URL path)</label>
          <input
            v-model="form.slug"
            type="text"
            placeholder="john-jane-2024"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
          />
        </div>
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">Title</label>
        <input
          v-model="form.title"
          type="text"
          placeholder="John & Jane's Wedding"
          required
          class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
        />
      </div>
      <div class="grid md:grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Date & Time</label>
          <input
            v-model="form.startsAt"
            type="datetime-local"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Venue</label>
          <input
            v-model="form.location"
            type="text"
            placeholder="The Garden Venue"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
          />
        </div>
      </div>
      <div class="flex items-center gap-2">
        <input
          v-model="form.isPublic"
          type="checkbox"
          id="isPublic"
          class="rounded"
        />
        <label for="isPublic" class="text-sm text-gray-700">Public event (no invite code required)</label>
      </div>
      <p v-if="error" class="text-red-600 text-sm">{{ error }}</p>
      <div class="flex gap-3">
        <button
          type="submit"
          :disabled="loading"
          class="bg-primary text-white font-medium px-4 py-2 rounded-lg hover:bg-primary-dark transition-colors disabled:opacity-50"
        >
          {{ loading ? 'Creating...' : 'Create Event' }}
        </button>
        <button
          type="button"
          @click="emit('cancel')"
          class="text-gray-600 hover:text-gray-800 px-4 py-2"
        >
          Cancel
        </button>
      </div>
    </form>
  </div>
</template>
