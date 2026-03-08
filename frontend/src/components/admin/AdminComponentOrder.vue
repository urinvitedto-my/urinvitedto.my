<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAdminStore } from '@/stores/admin'
import { useToast } from '@/composables/useToast'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import type { ComponentConfig } from '@/types'

const props = defineProps<{
  eventId: string
  collapsed: boolean
}>()

const emit = defineEmits<{ toggle: [] }>()

const adminStore = useAdminStore()
const toast = useToast()

const loading = computed(() => adminStore.isSubLoading('enabledComponents', props.eventId))
const error = computed(() => adminStore.getSubError('enabledComponents', props.eventId))

const saving = ref(false)
const components = ref<ComponentConfig[]>([])

const DEFAULT_COMPONENTS: ComponentConfig[] = [
  { name: 'EventDetails', enabled: true, order: 1 },
  { name: 'LocationPhoto', enabled: true, order: 2 },
  { name: 'CountdownTimer', enabled: true, order: 3 },
  { name: 'EventMap', enabled: true, order: 4 },
  { name: 'EventSchedule', enabled: true, order: 5 },
  { name: 'EventGallery', enabled: true, order: 6 },
  { name: 'AttireGuide', enabled: true, order: 7 },
  { name: 'EventFAQ', enabled: true, order: 8 },
  { name: 'MonetaryGifts', enabled: true, order: 9 },
  { name: 'GiftGuide', enabled: true, order: 10 },
]

const DISPLAY_NAMES: Record<string, string> = {
  EventDetails: 'Event Details',
  LocationPhoto: 'Location Photo',
  CountdownTimer: 'Countdown Timer',
  EventMap: 'Map & Directions',
  EventSchedule: 'Schedule',
  EventGallery: 'Gallery',
  AttireGuide: 'Attire Guide',
  EventFAQ: 'FAQs',
  MonetaryGifts: 'Monetary Gifts',
  GiftGuide: 'Gift Guide',
}

onMounted(async () => {
  await Promise.all([
    adminStore.fetchEnabledComponents(props.eventId),
    adminStore.fetchCustomContent(props.eventId),
  ])

  const storeData = adminStore.getEnabledComponents(props.eventId)
  let comps = storeData.length > 0 ? [...storeData] : [...DEFAULT_COMPONENTS]

  comps = syncCustomSectionEntries(comps)
  components.value = comps
})

/**
 * Syncs component list with actual custom sections:
 * expands legacy "CustomSections" into individual entries,
 * adds new sections, and removes stale ones.
 */
function syncCustomSectionEntries(comps: ComponentConfig[]): ComponentConfig[] {
  const sections = adminStore.getCustomContent(props.eventId)?.customSections ?? []
  const sectionIds = new Set(sections.map((s) => s.id))
  const result = [...comps]

  const legacyIdx = result.findIndex((c) => c.name === 'CustomSections')
  if (legacyIdx !== -1) {
    const legacyComp = result[legacyIdx]!
    const expanded = sections.map((s, i) => ({
      name: `CustomSection:${s.id}`,
      enabled: legacyComp.enabled,
      order: legacyComp.order + i * 0.001,
    }))
    result.splice(legacyIdx, 1, ...expanded)
  }

  const existingNames = new Set(result.map((c) => c.name))
  const maxOrder = Math.max(...result.map((c) => c.order), 0)
  sections.forEach((s, i) => {
    const name = `CustomSection:${s.id}`
    if (!existingNames.has(name)) {
      result.push({ name, enabled: true, order: maxOrder + 1 + i })
    }
  })

  const filtered = result.filter((c) => {
    if (c.name.startsWith('CustomSection:')) {
      return sectionIds.has(c.name.slice(14))
    }
    return true
  })

  filtered.sort((a, b) => a.order - b.order)
  filtered.forEach((c, i) => { c.order = i + 1 })

  return filtered
}

/** Saves enabled components via the store. */
async function handleSave() {
  saving.value = true

  try {
    const saved = await adminStore.saveEnabledComponents(props.eventId, {
      components: components.value,
    })
    if (saved.components?.length) {
      components.value = saved.components.sort((a, b) => a.order - b.order)
    }
    toast.success('Component order saved')
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : 'Failed to save')
  } finally {
    saving.value = false
  }
}

/** Resets to default component configuration. */
function resetToDefaults() {
  components.value = syncCustomSectionEntries([...DEFAULT_COMPONENTS])
}

/** Moves a component up or down in order. */
function moveComponent(index: number, direction: 'up' | 'down') {
  const swapIndex = direction === 'up' ? index - 1 : index + 1
  if (swapIndex < 0 || swapIndex >= components.value.length) return

  const temp = components.value[index]!.order
  components.value[index]!.order = components.value[swapIndex]!.order
  components.value[swapIndex]!.order = temp

  components.value.sort((a, b) => a.order - b.order)
}

/** Returns a readable name for a component key. */
function displayName(name: string): string {
  if (name.startsWith('CustomSection:')) {
    const sectionId = name.slice(14)
    const sections = adminStore.getCustomContent(props.eventId)?.customSections ?? []
    const section = sections.find((s) => s.id === sectionId)
    return section ? `Custom: ${section.title}` : 'Custom Section'
  }
  return DISPLAY_NAMES[name] || name
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
        Component Order
      </button>
    </div>

    <template v-if="!collapsed">
      <div v-if="loading" class="flex items-center justify-center py-4">
        <LoadingSpinner size="sm" />
      </div>

      <p v-else-if="error" class="text-red-600 text-sm mb-3">{{ error }}</p>

      <div v-else class="space-y-2">
        <p class="text-xs text-gray-500 mb-2">
          Toggle visibility and reorder event page sections. Disabled sections won't appear on the event page.
        </p>

        <div
          v-for="(comp, index) in components"
          :key="comp.name"
          class="flex items-center justify-between bg-gray-50 rounded-lg px-4 py-2.5"
          :class="{ 'opacity-50': !comp.enabled }"
        >
          <div class="flex items-center gap-3">
            <div class="flex flex-col gap-0.5 shrink-0">
              <button
                @click="moveComponent(index, 'up')"
                :disabled="index === 0"
                class="text-gray-400 hover:text-gray-600 disabled:opacity-30 text-xs leading-none"
                title="Move up"
              >▲</button>
              <button
                @click="moveComponent(index, 'down')"
                :disabled="index === components.length - 1"
                class="text-gray-400 hover:text-gray-600 disabled:opacity-30 text-xs leading-none"
                title="Move down"
              >▼</button>
            </div>
            <span class="text-sm text-primary">{{ displayName(comp.name) }}</span>
          </div>
          <label class="relative inline-flex items-center cursor-pointer">
            <input v-model="comp.enabled" type="checkbox" class="sr-only peer" />
            <div class="w-9 h-5 bg-gray-300 rounded-full peer peer-checked:bg-accent transition-colors after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:after:translate-x-full"></div>
          </label>
        </div>

        <div class="flex items-center gap-3 pt-2">
          <button
            @click="handleSave"
            :disabled="saving"
            class="bg-primary text-white font-medium px-4 py-2 rounded-lg hover:bg-primary-dark transition-colors disabled:opacity-50 text-sm"
          >
            {{ saving ? 'Saving...' : 'Save Component Order' }}
          </button>
          <button
            @click="resetToDefaults"
            class="text-sm text-gray-500 hover:text-gray-700"
          >
            Reset to defaults
          </button>
        </div>
      </div>
    </template>
  </div>
</template>
