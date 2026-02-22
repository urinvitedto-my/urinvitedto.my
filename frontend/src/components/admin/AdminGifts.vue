<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {
  adminListGifts,
  adminCreateGift,
  adminUpdateGift,
  adminDeleteGift,
} from '@/services/api'
import type { AdminGift } from '@/types'

const props = defineProps<{
  eventId: string
  collapsed: boolean
}>()

const emit = defineEmits<{ toggle: [] }>()

const items = ref<AdminGift[]>([])
const loading = ref(false)
const error = ref('')

const showCreateForm = ref(false)
const createForm = ref({ giftType: 'physical', title: '', description: '', link: '' })
const createLoading = ref(false)

const editingItemId = ref<string | null>(null)
const editForm = ref({ giftType: 'physical', title: '', description: '', link: '', orderIndex: 0 })
const editLoading = ref(false)

onMounted(() => loadGifts())

/**
 * Loads all gifts for this event.
 */
async function loadGifts() {
  loading.value = true
  error.value = ''

  try {
    const data = await adminListGifts(props.eventId)
    items.value = data.items
  } catch (e: any) {
    error.value = e.message || 'Failed to load gifts'
  } finally {
    loading.value = false
  }
}

/**
 * Creates a new gift.
 */
async function handleCreate() {
  if (!createForm.value.title.trim()) return

  createLoading.value = true

  try {
    const newItem = await adminCreateGift(props.eventId, {
      giftType: createForm.value.giftType,
      title: createForm.value.title.trim(),
      description: createForm.value.description.trim() || null,
      link: createForm.value.link.trim() || null,
    })

    items.value.push(newItem)
    items.value.sort((a, b) => a.orderIndex - b.orderIndex)
    showCreateForm.value = false
    createForm.value = { giftType: 'physical', title: '', description: '', link: '' }
  } catch (e: any) {
    alert(e.message || 'Failed to create gift')
  } finally {
    createLoading.value = false
  }
}

/**
 * Opens edit mode for a gift.
 */
function startEdit(item: AdminGift) {
  editingItemId.value = item.id
  editForm.value = {
    giftType: item.giftType,
    title: item.title,
    description: item.description || '',
    link: item.link || '',
    orderIndex: item.orderIndex,
  }
}

/**
 * Saves edits to a gift.
 */
async function handleUpdate() {
  if (!editingItemId.value || !editForm.value.title.trim()) return

  editLoading.value = true

  try {
    const updated = await adminUpdateGift(props.eventId, editingItemId.value, {
      giftType: editForm.value.giftType,
      title: editForm.value.title.trim(),
      description: editForm.value.description.trim() || null,
      link: editForm.value.link.trim() || null,
      orderIndex: editForm.value.orderIndex,
    })

    const idx = items.value.findIndex((i) => i.id === editingItemId.value)
    if (idx !== -1) items.value[idx] = updated
    items.value.sort((a, b) => a.orderIndex - b.orderIndex)
    editingItemId.value = null
  } catch (e: any) {
    alert(e.message || 'Failed to update gift')
  } finally {
    editLoading.value = false
  }
}

/**
 * Deletes a gift.
 */
async function handleDelete(itemId: string) {
  if (!confirm('Delete this gift?')) return

  try {
    await adminDeleteGift(props.eventId, itemId)
    items.value = items.value.filter((i) => i.id !== itemId)
  } catch (e: any) {
    alert(e.message || 'Failed to delete gift')
  }
}

/**
 * Moves an item up or down in the order.
 */
async function moveItem(itemId: string, direction: 'up' | 'down') {
  const idx = items.value.findIndex((i) => i.id === itemId)
  if (idx === -1) return

  const swapIdx = direction === 'up' ? idx - 1 : idx + 1
  if (swapIdx < 0 || swapIdx >= items.value.length) return

  const current = items.value[idx]
  const swap = items.value[swapIdx]
  if (!current || !swap) return

  try {
    const [updatedCurrent, updatedSwap] = await Promise.all([
      adminUpdateGift(props.eventId, current.id, {
        giftType: current.giftType,
        title: current.title,
        description: current.description || null,
        link: current.link || null,
        orderIndex: swap.orderIndex,
      }),
      adminUpdateGift(props.eventId, swap.id, {
        giftType: swap.giftType,
        title: swap.title,
        description: swap.description || null,
        link: swap.link || null,
        orderIndex: current.orderIndex,
      }),
    ])

    items.value[idx] = updatedCurrent
    items.value[swapIdx] = updatedSwap
    items.value.sort((a, b) => a.orderIndex - b.orderIndex)
  } catch (e: any) {
    alert(e.message || 'Failed to reorder')
  }
}

/**
 * Returns badge class for gift type.
 */
function giftTypeClass(type: string): string {
  return type === 'monetary'
    ? 'bg-yellow-100 text-yellow-700'
    : 'bg-blue-100 text-blue-700'
}
</script>

<template>
  <div class="border-t border-gray-100 pt-4">
    <div class="flex items-center justify-between mb-3">
      <button
        @click="emit('toggle')"
        class="flex items-center gap-1.5 text-sm font-medium text-gray-700 hover:text-[#14213d] transition-colors"
      >
        <span
          class="inline-block transition-transform duration-200"
          :class="collapsed ? '' : 'rotate-90'"
        >▶</span>
        Gifts
        <span v-if="items.length > 0" class="text-gray-400 font-normal">({{ items.length }})</span>
      </button>
      <button
        v-if="!collapsed"
        @click="showCreateForm = !showCreateForm"
        class="text-sm text-[#14213d] hover:underline"
      >
        {{ showCreateForm ? 'Cancel' : '+ Add Gift' }}
      </button>
    </div>

    <template v-if="!collapsed">
    <!-- Create Form -->
    <div v-if="showCreateForm" class="bg-gray-50 p-4 rounded-lg mb-3">
      <form @submit.prevent="handleCreate" class="space-y-3">
        <div class="grid md:grid-cols-2 gap-3">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Type</label>
            <select
              v-model="createForm.giftType"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#fca311] focus:outline-none"
            >
              <option value="physical">Physical</option>
              <option value="monetary">Monetary</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Title</label>
            <input
              v-model="createForm.title"
              type="text"
              placeholder="e.g., Kitchen Set"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#fca311] focus:outline-none"
            />
          </div>
        </div>
        <div class="grid md:grid-cols-2 gap-3">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Description (optional)</label>
            <input
              v-model="createForm.description"
              type="text"
              placeholder="Brief description"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#fca311] focus:outline-none"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Link (optional)</label>
            <input
              v-model="createForm.link"
              type="url"
              placeholder="https://..."
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#fca311] focus:outline-none"
            />
          </div>
        </div>
        <button
          type="submit"
          :disabled="createLoading"
          class="bg-[#14213d] text-white font-medium px-4 py-2 rounded-lg hover:bg-[#1a2a4d] transition-colors disabled:opacity-50"
        >
          {{ createLoading ? 'Adding...' : 'Add Gift' }}
        </button>
      </form>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-4">
      <div class="animate-spin rounded-full h-6 w-6 border-2 border-[#fca311] border-t-transparent"></div>
    </div>

    <!-- Error -->
    <p v-else-if="error" class="text-red-600 text-sm mb-3">{{ error }}</p>

    <!-- Empty -->
    <p v-else-if="items.length === 0" class="text-sm text-gray-400 mb-3">No gifts yet</p>

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
                  <span
                    class="text-xs px-1.5 py-0.5 rounded capitalize"
                    :class="giftTypeClass(item.giftType)"
                  >
                    {{ item.giftType }}
                  </span>
                  <span class="text-sm font-medium text-[#14213d]">{{ item.title }}</span>
                </div>
                <div class="flex items-center gap-2 mt-0.5">
                  <p v-if="item.description" class="text-xs text-gray-500 truncate">{{ item.description }}</p>
                  <a
                    v-if="item.link"
                    :href="item.link"
                    target="_blank"
                    class="text-xs text-[#fca311] hover:underline shrink-0"
                  >Link</a>
                </div>
              </div>
            </div>
            <div class="flex items-center gap-2 shrink-0 ml-3">
              <button @click="startEdit(item)" class="text-xs text-[#14213d] hover:underline">Edit</button>
              <button @click="handleDelete(item.id)" class="text-xs text-red-500 hover:text-red-700">Delete</button>
            </div>
          </div>
        </template>

        <!-- Edit mode -->
        <template v-else>
          <form @submit.prevent="handleUpdate" class="space-y-3">
            <div class="grid md:grid-cols-2 gap-3">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Type</label>
                <select
                  v-model="editForm.giftType"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#fca311] focus:outline-none"
                >
                  <option value="physical">Physical</option>
                  <option value="monetary">Monetary</option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Title</label>
                <input
                  v-model="editForm.title"
                  type="text"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#fca311] focus:outline-none"
                />
              </div>
            </div>
            <div class="grid md:grid-cols-3 gap-3">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Description (optional)</label>
                <input
                  v-model="editForm.description"
                  type="text"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#fca311] focus:outline-none"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Link (optional)</label>
                <input
                  v-model="editForm.link"
                  type="url"
                  placeholder="https://..."
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#fca311] focus:outline-none"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Order</label>
                <input
                  v-model.number="editForm.orderIndex"
                  type="number"
                  min="0"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#fca311] focus:outline-none"
                />
              </div>
            </div>
            <div class="flex gap-2">
              <button
                type="submit"
                :disabled="editLoading"
                class="text-sm bg-[#14213d] text-white px-3 py-1.5 rounded-lg hover:bg-[#1a2a4d] disabled:opacity-50"
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
