<script setup lang="ts">
import { ref, watch } from 'vue'
import { useAdminStore } from '@/stores/admin'
import { toISOOrUndefined, toDatetimeLocal } from '@/utils/date'
import type { AdminEvent } from '@/types'

const props = defineProps<{ event: AdminEvent }>()
const emit = defineEmits<{ saved: []; cancel: [] }>()

const adminStore = useAdminStore()

const form = ref({
  type: props.event.type,
  slug: props.event.slug,
  title: props.event.title,
  description: props.event.description || '',
  isPublic: props.event.isPublic,
  startsAt: toDatetimeLocal(props.event.startsAt),
  location: props.event.location || '',
  coverImageUrl: props.event.coverImageUrl || '',
  locationPhotoUrl: props.event.locationPhotoUrl || '',
  musicUrl: props.event.musicUrl || '',
})
const loading = ref(false)
const error = ref('')

watch(() => props.event, (ev) => {
  form.value = {
    type: ev.type,
    slug: ev.slug,
    title: ev.title,
    description: ev.description || '',
    isPublic: ev.isPublic,
    startsAt: toDatetimeLocal(ev.startsAt),
    location: ev.location || '',
    coverImageUrl: ev.coverImageUrl || '',
    locationPhotoUrl: ev.locationPhotoUrl || '',
    musicUrl: ev.musicUrl || '',
  }
})

/** Saves edits to the event. */
async function handleSubmit() {
  loading.value = true
  error.value = ''

  try {
    await adminStore.updateEvent(props.event.id, {
      type: form.value.type,
      slug: form.value.slug,
      title: form.value.title,
      description: form.value.description || null,
      isPublic: form.value.isPublic,
      startsAt: toISOOrUndefined(form.value.startsAt) ?? null,
      location: form.value.location || null,
      coverImageUrl: form.value.coverImageUrl || null,
      locationPhotoUrl: form.value.locationPhotoUrl || null,
      musicUrl: form.value.musicUrl || null,
    })
    emit('saved')
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : 'Failed to update event'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <form @submit.prevent="handleSubmit" class="space-y-4 mb-4">
    <div class="flex items-center justify-between mb-2">
      <h3 class="text-lg font-semibold text-primary">Edit Event</h3>
      <button
        type="button"
        @click="emit('cancel')"
        class="text-sm text-gray-500 hover:text-gray-700"
      >
        Cancel
      </button>
    </div>
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
        required
        class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
      />
    </div>
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
      <textarea
        v-model="form.description"
        rows="3"
        class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
      ></textarea>
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
          class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
        />
      </div>
    </div>
    <div class="grid md:grid-cols-2 gap-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">Cover Image URL</label>
        <input
          v-model="form.coverImageUrl"
          type="url"
          placeholder="https://..."
          class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
        />
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">Venue Photo URL</label>
        <input
          v-model="form.locationPhotoUrl"
          type="url"
          placeholder="https://..."
          class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
        />
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">Background Music URL</label>
        <input
          v-model="form.musicUrl"
          type="url"
          placeholder="https://..."
          class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
        />
      </div>
    </div>
    <div class="flex items-center gap-2">
      <input
        v-model="form.isPublic"
        type="checkbox"
        id="editIsPublic"
        class="rounded"
      />
      <label for="editIsPublic" class="text-sm text-gray-700">Public event (no invite code required)</label>
    </div>
    <p v-if="error" class="text-red-600 text-sm">{{ error }}</p>
    <div class="flex gap-3">
      <button
        type="submit"
        :disabled="loading"
        class="bg-primary text-white font-medium px-4 py-2 rounded-lg hover:bg-primary-dark transition-colors disabled:opacity-50"
      >
        {{ loading ? 'Saving...' : 'Save Changes' }}
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
</template>
