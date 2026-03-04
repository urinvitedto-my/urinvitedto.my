<script setup lang="ts">
import { ref } from 'vue'
import { useAdminStore } from '@/stores/admin'
import { useToast } from '@/composables/useToast'
import type { AdminHost } from '@/types'

const props = defineProps<{
  eventId: string
  hosts: AdminHost[]
  collapsed: boolean
}>()

const emit = defineEmits<{ toggle: [] }>()

const adminStore = useAdminStore()
const toast = useToast()

const showForm = ref(false)
const form = ref({ email: '', displayName: '' })
const loading = ref(false)
const error = ref('')

/** Adds a host to the event. */
async function handleAdd() {
  loading.value = true
  error.value = ''

  try {
    await adminStore.addHost(props.eventId, {
      email: form.value.email,
      displayName: form.value.displayName,
    })
    form.value = { email: '', displayName: '' }
    showForm.value = false
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : 'Failed to add host'
  } finally {
    loading.value = false
  }
}

/** Removes a host after confirmation. */
async function handleDelete(hostId: string) {
  if (!(await toast.confirm('Remove this host?'))) return

  try {
    await adminStore.deleteHost(props.eventId, hostId)
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : 'Failed to remove host')
  }
}
</script>

<template>
  <div class="border-t border-gray-100 pt-4">
    <div class="flex items-center justify-between mb-3">
      <button
        @click="emit('toggle')"
        class="flex items-center gap-1.5 text-sm font-medium text-gray-700 hover:text-primary transition-colors"
      >
        <span
          class="inline-block transition-transform duration-200"
          :class="collapsed ? '' : 'rotate-90'"
        >▶</span>
        Hosts
        <span v-if="hosts.length > 0" class="text-gray-400 font-normal">({{ hosts.length }})</span>
      </button>
      <button
        v-if="!collapsed"
        @click="showForm = !showForm"
        class="text-sm text-primary hover:underline"
      >
        {{ showForm ? 'Cancel' : '+ Add Host' }}
      </button>
    </div>

    <template v-if="!collapsed">
      <!-- Host List -->
      <div v-if="hosts.length > 0" class="space-y-2 mb-3">
        <div
          v-for="host in hosts"
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
            @click="handleDelete(host.id)"
            class="text-red-500 hover:text-red-700 text-sm"
          >
            Remove
          </button>
        </div>
      </div>
      <p v-else class="text-sm text-gray-400 mb-3">No hosts added yet</p>

      <!-- Add Host Form -->
      <div v-if="showForm" class="bg-gray-50 p-4 rounded-lg">
        <form @submit.prevent="handleAdd" class="space-y-3">
          <div class="grid md:grid-cols-2 gap-3">
            <input
              v-model="form.email"
              type="email"
              placeholder="Host email"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
            />
            <input
              v-model="form.displayName"
              type="text"
              placeholder="Display name"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
            />
          </div>
          <p v-if="error" class="text-red-600 text-sm">{{ error }}</p>
          <button
            type="submit"
            :disabled="loading"
            class="bg-primary text-white font-medium px-4 py-2 rounded-lg hover:bg-primary-dark transition-colors disabled:opacity-50"
          >
            {{ loading ? 'Adding...' : 'Add Host' }}
          </button>
        </form>
      </div>
    </template>
  </div>
</template>
