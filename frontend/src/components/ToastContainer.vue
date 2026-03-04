<script setup lang="ts">
import { useToast } from '@/composables/useToast'

const { toasts, confirmState, dismiss, resolveConfirm } = useToast()
</script>

<template>
  <!-- Toast notifications -->
  <Teleport to="body">
    <div class="fixed top-4 right-4 z-100 flex flex-col gap-2 max-w-sm">
      <TransitionGroup name="toast">
        <div
          v-for="toast in toasts"
          :key="toast.id"
          :class="[
            'flex items-center gap-2 px-4 py-3 rounded-lg shadow-lg text-sm font-medium cursor-pointer',
            toast.type === 'success' && 'bg-success text-white',
            toast.type === 'error' && 'bg-error text-white',
            toast.type === 'info' && 'bg-primary text-white',
          ]"
          @click="dismiss(toast.id)"
        >
          <span v-if="toast.type === 'success'">&#10003;</span>
          <span v-else-if="toast.type === 'error'">&#10007;</span>
          <span>{{ toast.message }}</span>
        </div>
      </TransitionGroup>
    </div>

    <!-- Confirm dialog -->
    <Transition name="fade">
      <div
        v-if="confirmState.visible"
        class="fixed inset-0 z-101 flex items-center justify-center bg-black/40"
        @click.self="resolveConfirm(false)"
      >
        <div class="bg-white rounded-lg shadow-xl p-6 max-w-sm mx-4 w-full">
          <p class="text-sm text-gray-800 mb-5">{{ confirmState.message }}</p>
          <div class="flex justify-end gap-3">
            <button
              class="px-4 py-2 text-sm text-gray-600 hover:text-gray-800 transition-colors"
              @click="resolveConfirm(false)"
            >
              Cancel
            </button>
            <button
              class="px-4 py-2 text-sm font-medium bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
              @click="resolveConfirm(true)"
            >
              Confirm
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}
.toast-enter-from {
  opacity: 0;
  transform: translateX(100%);
}
.toast-leave-to {
  opacity: 0;
  transform: translateX(100%);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
