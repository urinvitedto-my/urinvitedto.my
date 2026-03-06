<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAdminStore } from '@/stores/admin'
import { useToast } from '@/composables/useToast'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import type { AdminGift } from '@/types'

const props = defineProps<{
  eventId: string
  collapsed: boolean
}>()

const emit = defineEmits<{ toggle: [] }>()

const adminStore = useAdminStore()
const toast = useToast()

const items = computed(() => adminStore.getGifts(props.eventId))
const loading = computed(() => adminStore.isSubLoading('gifts', props.eventId))
const error = computed(() => adminStore.getSubError('gifts', props.eventId))

type GiftType = 'physical' | 'monetary'

const showCreateForm = ref(false)
const createForm = ref<{ giftType: GiftType; title: string; description: string; link: string }>({
  giftType: 'physical', title: '', description: '', link: '',
})
const createLoading = ref(false)

const editingItemId = ref<string | null>(null)
const editForm = ref<{ giftType: GiftType; title: string; description: string; link: string; orderIndex: number }>({
  giftType: 'physical', title: '', description: '', link: '', orderIndex: 0,
})
const editLoading = ref(false)

onMounted(() => adminStore.fetchGifts(props.eventId))

/** Creates a new gift. */
async function handleCreate() {
  if (!createForm.value.title.trim()) return
  createLoading.value = true
  try {
    await adminStore.createGiftItem(props.eventId, {
      giftType: createForm.value.giftType,
      title: createForm.value.title.trim(),
      description: createForm.value.description.trim() || null,
      link: createForm.value.link.trim() || null,
    })
    showCreateForm.value = false
    createForm.value = { giftType: 'physical', title: '', description: '', link: '' }
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : 'Failed to create gift')
  } finally {
    createLoading.value = false
  }
}

/** Opens edit mode for a gift. */
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

/** Saves edits to a gift. */
async function handleUpdate() {
  if (!editingItemId.value || !editForm.value.title.trim()) return
  editLoading.value = true
  try {
    await adminStore.updateGiftItem(props.eventId, editingItemId.value, {
      giftType: editForm.value.giftType,
      title: editForm.value.title.trim(),
      description: editForm.value.description.trim() || null,
      link: editForm.value.link.trim() || null,
      orderIndex: editForm.value.orderIndex,
    })
    editingItemId.value = null
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : 'Failed to update gift')
  } finally {
    editLoading.value = false
  }
}

/** Deletes a gift. */
async function handleDelete(itemId: string) {
  if (!(await toast.confirm('Delete this gift?'))) return
  try {
    await adminStore.deleteGiftItem(props.eventId, itemId)
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : 'Failed to delete gift')
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
        return adminStore.updateGiftItem(props.eventId, id, {
          giftType: item.giftType,
          title: item.title,
          description: item.description || null,
          link: item.link || null,
          orderIndex,
        })
      },
    )
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : 'Failed to reorder')
  }
}

/** Returns badge class for gift type. */
function giftTypeClass(type: 'physical' | 'monetary'): string {
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
        class="flex items-center gap-1.5 text-sm font-medium text-gray-700 hover:text-primary transition-colors"
      >
        <span
          class="inline-block transition-transform duration-200"
          :class="collapsed ? '' : 'rotate-90'"
        >▶</span>
        Gift Guide
        <span v-if="items.length > 0" class="text-gray-400 font-normal">({{ items.length }})</span>
      </button>
      <button
        v-if="!collapsed"
        @click="showCreateForm = !showCreateForm"
        class="text-sm text-primary hover:underline"
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
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
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
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
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
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Link (optional)</label>
            <input
              v-model="createForm.link"
              type="url"
              placeholder="https://..."
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
            />
          </div>
        </div>
        <button
          type="submit"
          :disabled="createLoading"
          class="bg-primary text-white font-medium px-4 py-2 rounded-lg hover:bg-primary-dark transition-colors disabled:opacity-50"
        >
          {{ createLoading ? 'Adding...' : 'Add Gift' }}
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
                  <span class="text-sm font-medium text-primary">{{ item.title }}</span>
                </div>
                <div class="flex items-center gap-2 mt-0.5">
                  <p v-if="item.description" class="text-xs text-gray-500 truncate">{{ item.description }}</p>
                  <a
                    v-if="item.link"
                    :href="item.link"
                    target="_blank"
                    class="text-xs text-accent hover:underline shrink-0"
                  >Link</a>
                </div>
              </div>
            </div>
            <div class="flex items-center gap-2 shrink-0 ml-3">
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
                <label class="block text-sm font-medium text-gray-700 mb-1">Type</label>
                <select
                  v-model="editForm.giftType"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
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
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
                />
              </div>
            </div>
            <div class="grid md:grid-cols-3 gap-3">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Description (optional)</label>
                <input
                  v-model="editForm.description"
                  type="text"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Link (optional)</label>
                <input
                  v-model="editForm.link"
                  type="url"
                  placeholder="https://..."
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
