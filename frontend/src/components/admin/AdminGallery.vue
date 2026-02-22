<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { supabase } from '@/services/supabase'
import {
  adminListGallery,
  adminCreateGalleryItem,
  adminUpdateGalleryItem,
  adminDeleteGalleryItem,
} from '@/services/api'
import type { AdminGalleryItem } from '@/types'

const props = defineProps<{
  eventId: string
  collapsed: boolean
}>()

const emit = defineEmits<{ toggle: [] }>()

const items = ref<AdminGalleryItem[]>([])
const loading = ref(false)
const error = ref('')

const uploading = ref(false)
const uploadError = ref('')
const fileInput = ref<HTMLInputElement | null>(null)

const editingItemId = ref<string | null>(null)
const editForm = ref({ caption: '', orderIndex: 0 })
const editLoading = ref(false)

const BUCKET = 'event-media'

onMounted(() => loadGallery())

/**
 * Loads all gallery items for this event.
 */
async function loadGallery() {
  loading.value = true
  error.value = ''

  try {
    const data = await adminListGallery(props.eventId)
    items.value = data.items
  } catch (e: any) {
    error.value = e.message || 'Failed to load gallery'
  } finally {
    loading.value = false
  }
}

/**
 * Detects media type from file MIME type.
 */
function detectMediaType(file: File): 'photo' | 'video' {
  return file.type.startsWith('video/') ? 'video' : 'photo'
}

/**
 * Generates a unique storage path for the file.
 */
function buildStoragePath(file: File): string {
  const ext = file.name.split('.').pop() || 'bin'
  const timestamp = Date.now()
  const random = Math.random().toString(36).slice(2, 8)
  return `${props.eventId}/${timestamp}-${random}.${ext}`
}

/**
 * Extracts storage path from a Supabase public URL for deletion.
 */
function extractStoragePath(mediaUrl: string): string | null {
  const marker = `/storage/v1/object/public/${BUCKET}/`
  const idx = mediaUrl.indexOf(marker)
  if (idx === -1) return null
  return mediaUrl.slice(idx + marker.length)
}

/**
 * Handles file selection and upload.
 */
async function handleFileUpload(event: Event) {
  const input = event.target as HTMLInputElement
  const files = input.files
  if (!files || files.length === 0) return

  uploading.value = true
  uploadError.value = ''

  try {
    for (const file of Array.from(files)) {
      const path = buildStoragePath(file)
      const mediaType = detectMediaType(file)

      const { error: uploadErr } = await supabase.storage
        .from(BUCKET)
        .upload(path, file, { contentType: file.type })

      if (uploadErr) throw new Error(uploadErr.message)

      const { data: urlData } = supabase.storage
        .from(BUCKET)
        .getPublicUrl(path)

      const newItem = await adminCreateGalleryItem(props.eventId, {
        mediaType,
        mediaUrl: urlData.publicUrl,
      })

      items.value.push(newItem)
    }

    items.value.sort((a, b) => a.orderIndex - b.orderIndex)
  } catch (e: any) {
    uploadError.value = e.message || 'Failed to upload'
  } finally {
    uploading.value = false
    if (fileInput.value) fileInput.value.value = ''
  }
}

/**
 * Opens edit mode for a gallery item.
 */
function startEdit(item: AdminGalleryItem) {
  editingItemId.value = item.id
  editForm.value = {
    caption: item.caption || '',
    orderIndex: item.orderIndex,
  }
}

/**
 * Saves edits to a gallery item.
 */
async function handleUpdate() {
  if (!editingItemId.value) return

  editLoading.value = true

  try {
    const updated = await adminUpdateGalleryItem(props.eventId, editingItemId.value, {
      caption: editForm.value.caption.trim() || null,
      orderIndex: editForm.value.orderIndex,
    })

    const idx = items.value.findIndex((i) => i.id === editingItemId.value)
    if (idx !== -1) items.value[idx] = updated
    items.value.sort((a, b) => a.orderIndex - b.orderIndex)
    editingItemId.value = null
  } catch (e: any) {
    alert(e.message || 'Failed to update gallery item')
  } finally {
    editLoading.value = false
  }
}

/**
 * Deletes a gallery item and its file from storage.
 */
async function handleDelete(item: AdminGalleryItem) {
  if (!confirm('Delete this gallery item?')) return

  try {
    await adminDeleteGalleryItem(props.eventId, item.id)

    const storagePath = extractStoragePath(item.mediaUrl)
    if (storagePath) {
      await supabase.storage.from(BUCKET).remove([storagePath])
    }

    items.value = items.value.filter((i) => i.id !== item.id)
  } catch (e: any) {
    alert(e.message || 'Failed to delete gallery item')
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
      adminUpdateGalleryItem(props.eventId, current.id, {
        caption: current.caption || null,
        orderIndex: swap.orderIndex,
      }),
      adminUpdateGalleryItem(props.eventId, swap.id, {
        caption: swap.caption || null,
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
        Gallery
        <span v-if="items.length > 0" class="text-gray-400 font-normal">({{ items.length }})</span>
      </button>
      <div v-if="!collapsed" class="flex items-center gap-2">
        <label
          class="text-sm text-[#14213d] hover:underline cursor-pointer"
          :class="{ 'opacity-50 pointer-events-none': uploading }"
        >
          {{ uploading ? 'Uploading...' : '+ Upload' }}
          <input
            ref="fileInput"
            type="file"
            accept="image/*,video/*"
            multiple
            class="hidden"
            :disabled="uploading"
            @change="handleFileUpload"
          />
        </label>
      </div>
    </div>

    <template v-if="!collapsed">
      <!-- Upload Error -->
      <p v-if="uploadError" class="text-red-600 text-sm mb-3">{{ uploadError }}</p>

      <!-- Loading -->
      <div v-if="loading" class="flex items-center justify-center py-4">
        <div class="animate-spin rounded-full h-6 w-6 border-2 border-[#fca311] border-t-transparent"></div>
      </div>

      <!-- Error -->
      <p v-else-if="error" class="text-red-600 text-sm mb-3">{{ error }}</p>

      <!-- Empty -->
      <p v-else-if="items.length === 0" class="text-sm text-gray-400 mb-3">No gallery items yet</p>

      <!-- Gallery Grid -->
      <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3">
        <div
          v-for="(item, index) in items"
          :key="item.id"
          class="bg-gray-50 rounded-lg overflow-hidden"
        >
          <!-- Edit mode -->
          <template v-if="editingItemId === item.id">
            <form @submit.prevent="handleUpdate" class="p-3 space-y-2">
              <div>
                <label class="block text-xs font-medium text-gray-700 mb-1">Caption</label>
                <input
                  v-model="editForm.caption"
                  type="text"
                  placeholder="Optional caption"
                  class="w-full px-2 py-1.5 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#fca311] focus:outline-none"
                />
              </div>
              <div>
                <label class="block text-xs font-medium text-gray-700 mb-1">Order</label>
                <input
                  v-model.number="editForm.orderIndex"
                  type="number"
                  min="0"
                  class="w-full px-2 py-1.5 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#fca311] focus:outline-none"
                />
              </div>
              <div class="flex gap-2">
                <button
                  type="submit"
                  :disabled="editLoading"
                  class="text-xs bg-[#14213d] text-white px-2.5 py-1 rounded-lg hover:bg-[#1a2a4d] disabled:opacity-50"
                >
                  {{ editLoading ? 'Saving...' : 'Save' }}
                </button>
                <button
                  type="button"
                  @click="editingItemId = null"
                  class="text-xs text-gray-500 hover:text-gray-700"
                >
                  Cancel
                </button>
              </div>
            </form>
          </template>

          <!-- View mode -->
          <template v-else>
            <!-- Thumbnail -->
            <div class="aspect-square relative group">
              <img
                v-if="item.mediaType === 'photo'"
                :src="item.mediaUrl"
                :alt="item.caption || 'Gallery image'"
                class="w-full h-full object-cover"
                loading="lazy"
              />
              <div
                v-else
                class="w-full h-full bg-[#14213d] flex items-center justify-center"
              >
                <svg class="w-8 h-8 text-white" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M8 5v14l11-7z" />
                </svg>
              </div>
              <span
                class="absolute top-1 left-1 text-[10px] px-1.5 py-0.5 rounded capitalize bg-black/60 text-white"
              >
                {{ item.mediaType }}
              </span>
            </div>

            <!-- Info & Actions -->
            <div class="px-2 py-2">
              <p v-if="item.caption" class="text-xs text-gray-600 truncate mb-1">{{ item.caption }}</p>
              <div class="flex items-center justify-between">
                <div class="flex gap-0.5">
                  <button
                    @click="moveItem(item.id, 'up')"
                    :disabled="index === 0"
                    class="text-gray-400 hover:text-gray-600 disabled:opacity-30 text-xs"
                    title="Move up"
                  >▲</button>
                  <button
                    @click="moveItem(item.id, 'down')"
                    :disabled="index === items.length - 1"
                    class="text-gray-400 hover:text-gray-600 disabled:opacity-30 text-xs"
                    title="Move down"
                  >▼</button>
                </div>
                <div class="flex items-center gap-1.5">
                  <button @click="startEdit(item)" class="text-xs text-[#14213d] hover:underline">Edit</button>
                  <button @click="handleDelete(item)" class="text-xs text-red-500 hover:text-red-700">Delete</button>
                </div>
              </div>
            </div>
          </template>
        </div>
      </div>
    </template>
  </div>
</template>
