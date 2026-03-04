<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAdminStore } from '@/stores/admin'
import { useToast } from '@/composables/useToast'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import { toISO, toDatetimeLocal, formatTime } from '@/utils/date'
import type { AdminScheduleItem } from '@/types'

const props = defineProps<{
  eventId: string
  collapsed: boolean
}>()

const emit = defineEmits<{ toggle: [] }>()

const adminStore = useAdminStore()
const toast = useToast()

const items = computed(() => adminStore.getSchedule(props.eventId))
const loading = computed(() => adminStore.isSubLoading('schedule', props.eventId))
const error = computed(() => adminStore.getSubError('schedule', props.eventId))

const showCreateForm = ref(false)
const createForm = ref({ time: '', title: '', description: '' })
const createLoading = ref(false)

const editingItemId = ref<string | null>(null)
const editForm = ref({ time: '', title: '', description: '', orderIndex: 0 })
const editLoading = ref(false)

onMounted(() => adminStore.fetchSchedule(props.eventId))

/** Creates a new schedule item. */
async function handleCreate() {
  if (!createForm.value.title.trim() || !createForm.value.time) return
  createLoading.value = true
  try {
    await adminStore.createScheduleItem(props.eventId, {
      time: toISO(createForm.value.time),
      title: createForm.value.title.trim(),
      description: createForm.value.description.trim() || null,
    })
    showCreateForm.value = false
    createForm.value = { time: '', title: '', description: '' }
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : 'Failed to create schedule item')
  } finally {
    createLoading.value = false
  }
}

/** Opens edit mode for a schedule item. */
function startEdit(item: AdminScheduleItem) {
  editingItemId.value = item.id
  editForm.value = {
    time: toDatetimeLocal(item.time),
    title: item.title,
    description: item.description || '',
    orderIndex: item.orderIndex,
  }
}

/** Saves edits to a schedule item. */
async function handleUpdate() {
  if (!editingItemId.value || !editForm.value.title.trim() || !editForm.value.time) return
  editLoading.value = true
  try {
    await adminStore.updateScheduleItem(props.eventId, editingItemId.value, {
      time: toISO(editForm.value.time),
      title: editForm.value.title.trim(),
      description: editForm.value.description.trim() || null,
      orderIndex: editForm.value.orderIndex,
    })
    editingItemId.value = null
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : 'Failed to update schedule item')
  } finally {
    editLoading.value = false
  }
}

/** Deletes a schedule item. */
async function handleDelete(itemId: string) {
  if (!(await toast.confirm('Delete this schedule item?'))) return
  try {
    await adminStore.deleteScheduleItem(props.eventId, itemId)
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : 'Failed to delete schedule item')
  }
}

/** Moves an item up or down in the order. */
async function moveItem(itemId: string, direction: 'up' | 'down') {
  try {
    await adminStore.swapOrder(
      items.value,
      itemId,
      direction,
      (id, orderIndex) => {
        const item = items.value.find((i) => i.id === id)!
        return adminStore.updateScheduleItem(props.eventId, id, {
          time: item.time,
          title: item.title,
          description: item.description || null,
          orderIndex,
        })
      },
    )
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : 'Failed to reorder')
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
        Schedule
        <span v-if="items.length > 0" class="text-gray-400 font-normal">({{ items.length }})</span>
      </button>
      <button
        v-if="!collapsed"
        @click="showCreateForm = !showCreateForm"
        class="text-sm text-primary hover:underline"
      >
        {{ showCreateForm ? 'Cancel' : '+ Add Item' }}
      </button>
    </div>

    <template v-if="!collapsed">
    <!-- Create Form -->
    <div v-if="showCreateForm" class="bg-gray-50 p-4 rounded-lg mb-3">
      <form @submit.prevent="handleCreate" class="space-y-3">
        <div class="grid md:grid-cols-2 gap-3">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Time</label>
            <input
              v-model="createForm.time"
              type="datetime-local"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Title</label>
            <input
              v-model="createForm.title"
              type="text"
              placeholder="e.g., Ceremony"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
            />
          </div>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Description (optional)</label>
          <input
            v-model="createForm.description"
            type="text"
            placeholder="Brief description"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
          />
        </div>
        <button
          type="submit"
          :disabled="createLoading"
          class="bg-primary text-white font-medium px-4 py-2 rounded-lg hover:bg-primary-dark transition-colors disabled:opacity-50"
        >
          {{ createLoading ? 'Adding...' : 'Add to Schedule' }}
        </button>
      </form>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-4">
      <LoadingSpinner size="sm" />
    </div>

    <!-- Error -->
    <p v-else-if="error" class="text-red-600 text-sm mb-3">{{ error }}</p>

    <!-- Empty -->
    <p v-else-if="items.length === 0" class="text-sm text-gray-400 mb-3">No schedule items yet</p>

    <!-- Item List -->
    <div v-else class="space-y-2">
      <div
        v-for="(item, index) in items"
        :key="item.id"
        class="bg-gray-50 rounded-lg px-4 py-3"
      >
        <!-- View mode -->
        <template v-if="editingItemId !== item.id">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3 min-w-0">
              <div class="flex flex-col gap-0.5 shrink-0">
                <button
                  @click="moveItem(item.id, 'up')"
                  :disabled="index === 0"
                  class="text-gray-400 hover:text-gray-600 disabled:opacity-30 text-xs leading-none"
                  title="Move up"
                >▲</button>
                <button
                  @click="moveItem(item.id, 'down')"
                  :disabled="index === items.length - 1"
                  class="text-gray-400 hover:text-gray-600 disabled:opacity-30 text-xs leading-none"
                  title="Move down"
                >▼</button>
              </div>
              <div class="min-w-0">
                <div class="flex items-center gap-2">
                  <span class="text-xs text-gray-500 shrink-0">{{ formatTime(item.time) }}</span>
                  <span class="text-sm font-medium text-primary">{{ item.title }}</span>
                </div>
                <p v-if="item.description" class="text-xs text-gray-500 mt-0.5 truncate">{{ item.description }}</p>
              </div>
            </div>
            <div class="flex items-center gap-2 shrink-0">
              <button @click="startEdit(item)" class="text-xs text-primary hover:underline">Edit</button>
              <button @click="handleDelete(item.id)" class="text-xs text-red-500 hover:text-red-700">Delete</button>
            </div>
          </div>
        </template>

        <!-- Edit mode -->
        <template v-else>
          <form @submit.prevent="handleUpdate" class="space-y-3">
            <div class="grid md:grid-cols-2 gap-3">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Time</label>
                <input
                  v-model="editForm.time"
                  type="datetime-local"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Title</label>
                <input
                  v-model="editForm.title"
                  type="text"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
                />
              </div>
            </div>
            <div class="grid md:grid-cols-2 gap-3">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Description (optional)</label>
                <input
                  v-model="editForm.description"
                  type="text"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Order</label>
                <input
                  v-model.number="editForm.orderIndex"
                  type="number"
                  min="0"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
                />
              </div>
            </div>
            <div class="flex gap-2">
              <button
                type="submit"
                :disabled="editLoading"
                class="text-sm bg-primary text-white px-3 py-1.5 rounded-lg hover:bg-primary-dark disabled:opacity-50"
              >
                {{ editLoading ? 'Saving...' : 'Save' }}
              </button>
              <button
                type="button"
                @click="editingItemId = null"
                class="text-sm text-gray-500 hover:text-gray-700"
              >
                Cancel
              </button>
            </div>
          </form>
        </template>
      </div>
    </div>
    </template>
  </div>
</template>
