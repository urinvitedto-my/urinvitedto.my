<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAdminStore } from '@/stores/admin'
import { useToast } from '@/composables/useToast'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import type { AdminFAQ } from '@/types'

const props = defineProps<{
  eventId: string
  collapsed: boolean
}>()

const emit = defineEmits<{ toggle: [] }>()

const adminStore = useAdminStore()
const toast = useToast()

const items = computed(() => adminStore.getFAQs(props.eventId))
const loading = computed(() => adminStore.isSubLoading('faqs', props.eventId))
const error = computed(() => adminStore.getSubError('faqs', props.eventId))

const showCreateForm = ref(false)
const createForm = ref({ question: '', answer: '' })
const createLoading = ref(false)

const editingItemId = ref<string | null>(null)
const editForm = ref({ question: '', answer: '', orderIndex: 0 })
const editLoading = ref(false)

onMounted(() => adminStore.fetchFAQs(props.eventId))

/** Creates a new FAQ. */
async function handleCreate() {
  if (!createForm.value.question.trim() || !createForm.value.answer.trim()) return
  createLoading.value = true
  try {
    await adminStore.createFAQItem(props.eventId, {
      question: createForm.value.question.trim(),
      answer: createForm.value.answer.trim(),
    })
    showCreateForm.value = false
    createForm.value = { question: '', answer: '' }
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : 'Failed to create FAQ')
  } finally {
    createLoading.value = false
  }
}

/** Opens edit mode for a FAQ. */
function startEdit(item: AdminFAQ) {
  editingItemId.value = item.id
  editForm.value = {
    question: item.question,
    answer: item.answer,
    orderIndex: item.orderIndex,
  }
}

/** Saves edits to a FAQ. */
async function handleUpdate() {
  if (!editingItemId.value || !editForm.value.question.trim() || !editForm.value.answer.trim()) return
  editLoading.value = true
  try {
    await adminStore.updateFAQItem(props.eventId, editingItemId.value, {
      question: editForm.value.question.trim(),
      answer: editForm.value.answer.trim(),
      orderIndex: editForm.value.orderIndex,
    })
    editingItemId.value = null
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : 'Failed to update FAQ')
  } finally {
    editLoading.value = false
  }
}

/** Deletes a FAQ. */
async function handleDelete(itemId: string) {
  if (!(await toast.confirm('Delete this FAQ?'))) return
  try {
    await adminStore.deleteFAQItem(props.eventId, itemId)
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : 'Failed to delete FAQ')
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
        return adminStore.updateFAQItem(props.eventId, id, {
          question: item.question,
          answer: item.answer,
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
        FAQs
        <span v-if="items.length > 0" class="text-gray-400 font-normal">({{ items.length }})</span>
      </button>
      <button
        v-if="!collapsed"
        @click="showCreateForm = !showCreateForm"
        class="text-sm text-primary hover:underline"
      >
        {{ showCreateForm ? 'Cancel' : '+ Add FAQ' }}
      </button>
    </div>

    <template v-if="!collapsed">
    <!-- Create Form -->
    <div v-if="showCreateForm" class="bg-gray-50 p-4 rounded-lg mb-3">
      <form @submit.prevent="handleCreate" class="space-y-3">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Question</label>
          <input
            v-model="createForm.question"
            type="text"
            placeholder="e.g., What should I wear?"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Answer</label>
          <textarea
            v-model="createForm.answer"
            placeholder="Your answer here..."
            required
            rows="2"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
          ></textarea>
        </div>
        <button
          type="submit"
          :disabled="createLoading"
          class="bg-primary text-white font-medium px-4 py-2 rounded-lg hover:bg-primary-dark transition-colors disabled:opacity-50"
        >
          {{ createLoading ? 'Adding...' : 'Add FAQ' }}
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
    <p v-else-if="items.length === 0" class="text-sm text-gray-400 mb-3">No FAQs yet</p>

    <!-- Item List -->
    <div v-else class="space-y-2">
      <div
        v-for="(item, index) in items"
        :key="item.id"
        class="bg-gray-50 rounded-lg px-4 py-3"
      >
        <!-- View mode -->
        <template v-if="editingItemId !== item.id">
          <div class="flex items-start justify-between">
            <div class="flex items-start gap-3 min-w-0">
              <div class="flex flex-col gap-0.5 shrink-0 pt-0.5">
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
                <p class="text-sm font-medium text-primary">{{ item.question }}</p>
                <p class="text-xs text-gray-600 mt-1">{{ item.answer }}</p>
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
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Question</label>
              <input
                v-model="editForm.question"
                type="text"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
              />
            </div>
            <div class="grid md:grid-cols-2 gap-3">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Answer</label>
                <textarea
                  v-model="editForm.answer"
                  required
                  rows="2"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
                ></textarea>
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
