<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAdminStore } from '@/stores/admin'
import { useToast } from '@/composables/useToast'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import type { CustomContent, CustomSection } from '@/types'

const props = defineProps<{
  eventId: string
  collapsed: boolean
}>()

const emit = defineEmits<{ toggle: [] }>()

const adminStore = useAdminStore()
const toast = useToast()

const loading = computed(() => adminStore.isSubLoading('customContent', props.eventId))
const error = computed(() => adminStore.getSubError('customContent', props.eventId))

const saving = ref(false)

const dressCodeEnabled = ref(false)
const dressCode = ref({ title: '', description: '', notes: '', examples: [] as string[] })
const newExample = ref('')

const locationEnabled = ref(false)
const location = ref({ parkingInfo: '', accessibilityNotes: '', mapEmbedUrl: '' })

const monetaryEnabled = ref(false)
const monetary = ref({ qrCodeUrl: '', instructions: '', accounts: [] as { method: string; number: string; name: string }[] })
const newAccount = ref({ method: '', number: '', name: '' })

const countdownEnabled = ref(false)
const countdown = ref({ customMessage: '' })

const customSections = ref<CustomSection[]>([])
const showAddSection = ref(false)
const DEFAULT_SECTION_BG = 'transparent'
const newSection = ref({ title: '', content: '', image: '', bgColor: DEFAULT_SECTION_BG })
const editingSectionId = ref<string | null>(null)
const editSection = ref({ title: '', content: '', image: '', bgColor: '', order: 0 })

onMounted(async () => {
  await adminStore.fetchCustomContent(props.eventId)
  const data = adminStore.getCustomContent(props.eventId)
  if (data) populateFromData(data)
})

/** Populates form refs from fetched data. */
function populateFromData(data: CustomContent) {
  if (data.dressCode) {
    dressCodeEnabled.value = true
    dressCode.value = {
      title: data.dressCode.title || '',
      description: data.dressCode.description || '',
      notes: data.dressCode.notes || '',
      examples: data.dressCode.examples || [],
    }
  }

  if (data.locationDetails) {
    locationEnabled.value = true
    location.value = {
      parkingInfo: data.locationDetails.parkingInfo || '',
      accessibilityNotes: data.locationDetails.accessibilityNotes || '',
      mapEmbedUrl: data.locationDetails.mapEmbedUrl || '',
    }
  }

  if (data.monetaryGifts) {
    monetaryEnabled.value = data.monetaryGifts.enabled
    monetary.value = {
      qrCodeUrl: data.monetaryGifts.qrCodeUrl || '',
      instructions: data.monetaryGifts.instructions || '',
      accounts: data.monetaryGifts.accounts || [],
    }
  }

  if (data.countdownTimer) {
    countdownEnabled.value = data.countdownTimer.enabled
    countdown.value = { customMessage: data.countdownTimer.customMessage || '' }
  }

  customSections.value = data.customSections || []
}

/** Builds the CustomContent object from form state. */
function buildPayload(): CustomContent {
  const payload: CustomContent = {}

  if (dressCodeEnabled.value) {
    payload.dressCode = {
      title: dressCode.value.title.trim(),
      description: dressCode.value.description.trim(),
      notes: dressCode.value.notes.trim() || undefined,
      examples: dressCode.value.examples.length > 0 ? dressCode.value.examples : undefined,
    }
  }

  if (locationEnabled.value) {
    payload.locationDetails = {
      parkingInfo: location.value.parkingInfo.trim() || undefined,
      accessibilityNotes: location.value.accessibilityNotes.trim() || undefined,
      mapEmbedUrl: location.value.mapEmbedUrl.trim() || undefined,
    }
  }

  payload.monetaryGifts = {
    enabled: monetaryEnabled.value,
    qrCodeUrl: monetary.value.qrCodeUrl.trim() || undefined,
    instructions: monetary.value.instructions.trim() || undefined,
    accounts: monetary.value.accounts.length > 0 ? monetary.value.accounts : undefined,
  }

  payload.countdownTimer = {
    enabled: countdownEnabled.value,
    customMessage: countdown.value.customMessage.trim() || undefined,
  }

  if (customSections.value.length > 0) {
    payload.customSections = customSections.value
  }

  return payload
}

/** Saves all custom content via the store. */
async function handleSave() {
  saving.value = true

  try {
    const saved = await adminStore.saveCustomContent(props.eventId, buildPayload())
    populateFromData(saved)
    toast.success('Custom content saved')
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : 'Failed to save')
  } finally {
    saving.value = false
  }
}

/** Adds an example to the dress code list. */
function addExample() {
  const val = newExample.value.trim()
  if (!val) return
  dressCode.value.examples.push(val)
  newExample.value = ''
}

/** Removes an example from the dress code list. */
function removeExample(index: number) {
  dressCode.value.examples.splice(index, 1)
}

/** Adds a monetary account. */
function addAccount() {
  if (!newAccount.value.method.trim() || !newAccount.value.number.trim()) return
  monetary.value.accounts.push({
    method: newAccount.value.method.trim(),
    number: newAccount.value.number.trim(),
    name: newAccount.value.name.trim(),
  })
  newAccount.value = { method: '', number: '', name: '' }
}

/** Removes a monetary account. */
function removeAccount(index: number) {
  monetary.value.accounts.splice(index, 1)
}

/** Adds a custom section. */
function addCustomSection() {
  if (!newSection.value.title.trim()) return
  customSections.value.push({
    id: crypto.randomUUID(),
    title: newSection.value.title.trim(),
    content: newSection.value.content.trim(),
    image: newSection.value.image.trim() || undefined,
    bgColor: newSection.value.bgColor || DEFAULT_SECTION_BG,
    order: customSections.value.length,
  })
  newSection.value = { title: '', content: '', image: '', bgColor: DEFAULT_SECTION_BG }
  showAddSection.value = false
}

/** Starts editing a custom section. */
function startEditSection(section: CustomSection) {
  editingSectionId.value = section.id
  editSection.value = {
    title: section.title,
    content: section.content,
    image: section.image || '',
    bgColor: section.bgColor || DEFAULT_SECTION_BG,
    order: section.order,
  }
}

/** Saves edits to a custom section. */
function saveEditSection() {
  if (!editingSectionId.value) return
  const idx = customSections.value.findIndex((s) => s.id === editingSectionId.value)
  if (idx !== -1) {
    customSections.value[idx] = {
      id: editingSectionId.value,
      title: editSection.value.title.trim(),
      content: editSection.value.content.trim(),
      image: editSection.value.image.trim() || undefined,
      bgColor: editSection.value.bgColor || DEFAULT_SECTION_BG,
      order: editSection.value.order,
    }
  }
  editingSectionId.value = null
}

/** Removes a custom section. */
function removeSection(id: string) {
  customSections.value = customSections.value.filter((s) => s.id !== id)
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
        Custom Content
      </button>
    </div>

    <template v-if="!collapsed">
      <div v-if="loading" class="flex items-center justify-center py-4">
        <LoadingSpinner size="sm" />
      </div>

      <p v-else-if="error" class="text-red-600 text-sm mb-3">{{ error }}</p>

      <div v-else class="space-y-4">
        <!-- Dress Code -->
        <div class="bg-gray-50 rounded-lg p-4">
          <div class="flex items-center justify-between mb-3">
            <h4 class="text-sm font-medium text-primary">Dress Code</h4>
            <label class="flex items-center gap-2 cursor-pointer">
              <input v-model="dressCodeEnabled" type="checkbox" class="rounded" />
              <span class="text-xs text-gray-500">{{ dressCodeEnabled ? 'Enabled' : 'Disabled' }}</span>
            </label>
          </div>
          <template v-if="dressCodeEnabled">
            <div class="space-y-2">
              <input
                v-model="dressCode.title"
                type="text"
                placeholder="e.g., Semi-Formal"
                class="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
              />
              <textarea
                v-model="dressCode.description"
                rows="2"
                placeholder="Description"
                class="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
              ></textarea>
              <input
                v-model="dressCode.notes"
                type="text"
                placeholder="Notes (optional)"
                class="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
              />
              <div>
                <p class="text-xs text-gray-500 mb-1">Examples</p>
                <div class="flex flex-wrap gap-1 mb-2">
                  <span
                    v-for="(ex, i) in dressCode.examples"
                    :key="i"
                    class="inline-flex items-center gap-1 px-2 py-0.5 bg-white rounded text-xs text-gray-700 border"
                  >
                    {{ ex }}
                    <button @click="removeExample(i)" class="text-red-400 hover:text-red-600">&times;</button>
                  </span>
                </div>
                <div class="flex gap-2">
                  <input
                    v-model="newExample"
                    type="text"
                    placeholder="Add example"
                    class="flex-1 px-3 py-1.5 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
                    @keyup.enter="addExample"
                  />
                  <button @click="addExample" class="text-xs text-primary hover:underline">Add</button>
                </div>
              </div>
            </div>
          </template>
        </div>

        <!-- Location Details -->
        <div class="bg-gray-50 rounded-lg p-4">
          <div class="flex items-center justify-between mb-3">
            <h4 class="text-sm font-medium text-primary">Location Details</h4>
            <label class="flex items-center gap-2 cursor-pointer">
              <input v-model="locationEnabled" type="checkbox" class="rounded" />
              <span class="text-xs text-gray-500">{{ locationEnabled ? 'Enabled' : 'Disabled' }}</span>
            </label>
          </div>
          <template v-if="locationEnabled">
            <div class="space-y-2">
              <input
                v-model="location.mapEmbedUrl"
                type="url"
                placeholder="Google Maps embed URL"
                class="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
              />
              <textarea
                v-model="location.parkingInfo"
                rows="2"
                placeholder="Parking information"
                class="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
              ></textarea>
              <textarea
                v-model="location.accessibilityNotes"
                rows="2"
                placeholder="Accessibility notes"
                class="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
              ></textarea>
            </div>
          </template>
        </div>

        <!-- Monetary Gifts -->
        <div class="bg-gray-50 rounded-lg p-4">
          <div class="flex items-center justify-between mb-3">
            <h4 class="text-sm font-medium text-primary">Monetary Gifts</h4>
            <label class="flex items-center gap-2 cursor-pointer">
              <input v-model="monetaryEnabled" type="checkbox" class="rounded" />
              <span class="text-xs text-gray-500">{{ monetaryEnabled ? 'Enabled' : 'Disabled' }}</span>
            </label>
          </div>
          <template v-if="monetaryEnabled">
            <div class="space-y-2">
              <input
                v-model="monetary.qrCodeUrl"
                type="url"
                placeholder="QR code image URL"
                class="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
              />
              <textarea
                v-model="monetary.instructions"
                rows="2"
                placeholder="Payment instructions"
                class="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
              ></textarea>
              <div>
                <p class="text-xs text-gray-500 mb-1">Payment Accounts</p>
                <div v-if="monetary.accounts.length > 0" class="space-y-1 mb-2">
                  <div
                    v-for="(acc, i) in monetary.accounts"
                    :key="i"
                    class="flex items-center justify-between bg-white px-3 py-1.5 rounded text-sm border"
                  >
                    <span>{{ acc.method }} &middot; {{ acc.number }} &middot; {{ acc.name }}</span>
                    <button @click="removeAccount(i)" class="text-red-400 hover:text-red-600 text-xs">Remove</button>
                  </div>
                </div>
                <div class="grid grid-cols-3 gap-2">
                  <input
                    v-model="newAccount.method"
                    type="text"
                    placeholder="Method (e.g., GCash)"
                    class="px-2 py-1.5 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
                  />
                  <input
                    v-model="newAccount.number"
                    type="text"
                    placeholder="Account number"
                    class="px-2 py-1.5 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
                  />
                  <div class="flex gap-1">
                    <input
                      v-model="newAccount.name"
                      type="text"
                      placeholder="Name"
                      class="flex-1 min-w-0 px-2 py-1.5 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
                    />
                    <button @click="addAccount" class="text-xs text-primary hover:underline shrink-0">Add</button>
                  </div>
                </div>
              </div>
            </div>
          </template>
        </div>

        <!-- Countdown Timer -->
        <div class="bg-gray-50 rounded-lg p-4">
          <div class="flex items-center justify-between mb-3">
            <h4 class="text-sm font-medium text-primary">Countdown Timer</h4>
            <label class="flex items-center gap-2 cursor-pointer">
              <input v-model="countdownEnabled" type="checkbox" class="rounded" />
              <span class="text-xs text-gray-500">{{ countdownEnabled ? 'Enabled' : 'Disabled' }}</span>
            </label>
          </div>
          <template v-if="countdownEnabled">
            <input
              v-model="countdown.customMessage"
              type="text"
              placeholder="Custom message above the countdown"
              class="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
            />
          </template>
        </div>

        <!-- Custom Sections -->
        <div class="bg-gray-50 rounded-lg p-4">
          <div class="flex items-center justify-between mb-3">
            <h4 class="text-sm font-medium text-primary">
              Custom Sections
              <span v-if="customSections.length > 0" class="text-gray-400 font-normal">({{ customSections.length }})</span>
            </h4>
            <button
              @click="showAddSection = !showAddSection"
              class="text-xs text-primary hover:underline"
            >
              {{ showAddSection ? 'Cancel' : '+ Add Section' }}
            </button>
          </div>

          <!-- Add section form -->
          <div v-if="showAddSection" class="bg-white p-3 rounded-lg mb-3 border space-y-2">
            <input
              v-model="newSection.title"
              type="text"
              placeholder="Section title"
              class="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
            />
            <textarea
              v-model="newSection.content"
              rows="3"
              placeholder="Content (HTML supported)"
              class="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
            ></textarea>
            <div class="grid grid-cols-2 gap-2">
              <input
                v-model="newSection.image"
                type="url"
                placeholder="Image URL (optional)"
                class="px-3 py-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
              />
              <div class="flex items-center gap-2">
                <label class="text-xs text-gray-500 shrink-0">Background</label>
                <input v-model="newSection.bgColor" type="color" class="w-8 h-8 rounded cursor-pointer" />
                <span class="text-xs text-gray-400">{{ newSection.bgColor }}</span>
              </div>
            </div>
            <button
              @click="addCustomSection"
              class="text-xs bg-primary text-white px-3 py-1.5 rounded-lg hover:bg-primary-dark"
            >
              Add Section
            </button>
          </div>

          <!-- Sections list -->
          <div v-if="customSections.length > 0" class="space-y-2">
            <div
              v-for="section in customSections"
              :key="section.id"
              class="bg-white p-3 rounded-lg border"
            >
              <template v-if="editingSectionId !== section.id">
                <div class="flex items-center justify-between">
                  <div class="min-w-0">
                    <span class="text-sm font-medium text-primary">{{ section.title }}</span>
                    <p v-if="section.content" class="text-xs text-gray-500 truncate mt-0.5">{{ section.content.slice(0, 80) }}</p>
                  </div>
                  <div class="flex items-center gap-2 shrink-0 ml-3">
                    <button @click="startEditSection(section)" class="text-xs text-primary hover:underline">Edit</button>
                    <button @click="removeSection(section.id)" class="text-xs text-red-500 hover:text-red-700">Delete</button>
                  </div>
                </div>
              </template>
              <template v-else>
                <div class="space-y-2">
                  <input
                    v-model="editSection.title"
                    type="text"
                    class="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
                  />
                  <textarea
                    v-model="editSection.content"
                    rows="3"
                    class="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
                  ></textarea>
                  <div class="grid grid-cols-3 gap-2">
                    <input
                      v-model="editSection.image"
                      type="url"
                      placeholder="Image URL"
                      class="px-3 py-2 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
                    />
                    <div class="flex items-center gap-2">
                      <label class="text-xs text-gray-500 shrink-0">BG</label>
                      <input v-model="editSection.bgColor" type="color" class="w-8 h-8 rounded cursor-pointer" />
                    </div>
                    <div class="flex items-center gap-1">
                      <label class="text-xs text-gray-500 shrink-0">Order</label>
                      <input
                        v-model.number="editSection.order"
                        type="number"
                        min="0"
                        class="w-full px-2 py-1.5 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-accent focus:outline-none"
                      />
                    </div>
                  </div>
                  <div class="flex gap-2">
                    <button
                      @click="saveEditSection"
                      class="text-xs bg-primary text-white px-3 py-1.5 rounded-lg hover:bg-primary-dark"
                    >
                      Save
                    </button>
                    <button
                      @click="editingSectionId = null"
                      class="text-xs text-gray-500 hover:text-gray-700"
                    >
                      Cancel
                    </button>
                  </div>
                </div>
              </template>
            </div>
          </div>
          <p v-else class="text-xs text-gray-400">No custom sections</p>
        </div>

        <!-- Save All -->
        <div class="flex items-center gap-3">
          <button
            @click="handleSave"
            :disabled="saving"
            class="bg-primary text-white font-medium px-4 py-2 rounded-lg hover:bg-primary-dark transition-colors disabled:opacity-50 text-sm"
          >
            {{ saving ? 'Saving...' : 'Save Custom Content' }}
          </button>
        </div>
      </div>
    </template>
  </div>
</template>
