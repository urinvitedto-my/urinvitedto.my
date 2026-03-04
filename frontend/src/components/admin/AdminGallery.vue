<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { supabase } from '@/services/supabase'
import { useAdminStore } from '@/stores/admin'
import { useToast } from '@/composables/useToast'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import type { AdminGalleryItem } from '@/types'

const props = defineProps<{
  eventId: string
  collapsed: boolean
}>()

const emit = defineEmits<{ toggle: [] }>()

const adminStore = useAdminStore()
const toast = useToast()

const items = computed(() => adminStore.getGallery(props.eventId))
const loading = computed(() => adminStore.isSubLoading('gallery', props.eventId))
const error = computed(() => adminStore.getSubError('gallery', props.eventId))

const uploading = ref(false)
const uploadError = ref('')
const fileInput = ref<HTMLInputElement | null>(null)

const editingItemId = ref<string | null>(null)
const editForm = ref({ caption: '', orderIndex: 0 })
const editLoading = ref(false)

const BUCKET = 'event-media'

onMounted(() => adminStore.fetchGallery(props.eventId))

/** Detects media type from file MIME type. */
function detectMediaType(file: File): 'photo' | 'video' {
  return file.type.startsWith('video/') ? 'video' : 'photo'
}

/** Generates a unique storage path for the file. */
function buildStoragePath(file: File): string {
  const ext = file.name.split('.').pop() || 'bin'
  const timestamp = Date.now()
  const random = Math.random().toString(36).slice(2, 8)
  return `${props.eventId}/${timestamp}-${random}.${ext}`
}

/** Extracts storage path from a Supabase public URL for deletion. */
function extractStoragePath(mediaUrl: string): string | null {
  const marker = `/storage/v1/object/public/${BUCKET}/`
  const idx = mediaUrl.indexOf(marker)
  if (idx === -1) return null
  return mediaUrl.slice(idx + marker.length)
}

/** Handles file selection and upload. */
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

      await adminStore.createGalleryItem(props.eventId, {
        mediaType,
        mediaUrl: urlData.publicUrl,
      })
    }

    toast.success('Upload complete')
  } catch (e: unknown) {
    uploadError.value = e instanceof Error ? e.message : 'Failed to upload'
  } finally {
    uploading.value = false
    if (fileInput.value) fileInput.value.value = ''
  }
}

/** Opens edit mode for a gallery item. */
function startEdit(item: AdminGalleryItem) {
  editingItemId.value = item.id
  editForm.value = {
    caption: item.caption || '',
    orderIndex: item.orderIndex,
  }
}

/** Saves edits to a gallery item. */
async function handleUpdate() {
  if (!editingItemId.value) return

  editLoading.value = true

  try {
    await adminStore.updateGalleryItem(props.eventId, editingItemId.value, {
      caption: editForm.value.caption.trim() || null,
      orderIndex: editForm.value.orderIndex,
    })
    editingItemId.value = null
    toast.success('Gallery item updated')
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : 'Failed to update gallery item')
  } finally {
    editLoading.value = false
  }
}

/** Deletes a gallery item and its file from storage. */
async function handleDelete(item: AdminGalleryItem) {
  const confirmed = await toast.confirm('Delete this gallery item?')
  if (!confirmed) return

  try {
    await adminStore.deleteGalleryItemFromStore(props.eventId, item.id)

    const storagePath = extractStoragePath(item.mediaUrl)
    if (storagePath) {
      await supabase.storage.from(BUCKET).remove([storagePath])
    }

    toast.success('Gallery item deleted')
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : 'Failed to delete gallery item')
  }
}

/** Moves an item up or down in the order. */
async function moveItem(itemId: string, direction: 'up' | 'down') {
  try {
    await adminStore.swapOrder(
      adminStore.getGallery(props.eventId),
      itemId,
      direction,
      async (id, orderIndex) => adminStore.updateGalleryItem(props.eventId, id, { orderIndex }),
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
        Gallery
        <span v-if="items.length > 0" class="text-gray-400 font-normal">({{ items.length }})</span>
      </button>
      <div v-if="!collapsed" class="flex items-center gap-2">
        <label
          class="text-sm text-primary hover:underline cursor-pointer"
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
      <p v-if="uploadError" class="text-red-600 text-sm mb-3">{{ uploadError }}</p>

      <div v-if="loading" class="flex items-center justify-center py-4">
        <LoadingSpinner size="sm" />
      </div>

      <p v-else-if="error" class="text-red-600 text-sm mb-3">{{ error }}</p>

      <p v-else-if="items.length === 0" class="text-sm text-gray-400 mb-3">No gallery items yet</p>

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
                  class="w-full px-2 py-1.5 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
                />
              </div>
              <div>
                <label class="block text-xs font-medium text-gray-700 mb-1">Order</label>
                <input
                  v-model.number="editForm.orderIndex"
                  type="number"
                  min="0"
                  class="w-full px-2 py-1.5 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
                />
              </div>
              <div class="flex gap-2">
                <button
                  type="submit"
                  :disabled="editLoading"
                  class="text-xs bg-primary text-white px-2.5 py-1 rounded-lg hover:bg-primary-dark disabled:opacity-50"
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
                class="w-full h-full bg-primary flex items-center justify-center"
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
                  <button @click="startEdit(item)" class="text-xs text-primary hover:underline">Edit</button>
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
