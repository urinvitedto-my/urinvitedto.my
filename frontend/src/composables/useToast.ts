/** Lightweight toast notification composable (replaces alert/confirm). */
import { ref } from 'vue'

export type ToastType = 'success' | 'error' | 'info'

export interface Toast {
  id: number
  message: string
  type: ToastType
}

export interface ConfirmState {
  visible: boolean
  message: string
  resolve: ((value: boolean) => void) | null
}

const toasts = ref<Toast[]>([])
const confirmState = ref<ConfirmState>({ visible: false, message: '', resolve: null })
let nextId = 0

export function useToast() {
  /** Shows a toast notification. */
  function show(message: string, type: ToastType = 'info', duration = 3000) {
    const id = nextId++
    toasts.value.push({ id, message, type })
    setTimeout(() => dismiss(id), duration)
  }

  function dismiss(id: number) {
    toasts.value = toasts.value.filter((t) => t.id !== id)
  }

  function success(message: string) {
    show(message, 'success')
  }

  function error(message: string) {
    show(message, 'error', 5000)
  }

  /** Replaces native confirm() -- returns a promise that resolves to true/false. */
  function confirm(message: string): Promise<boolean> {
    return new Promise((resolve) => {
      confirmState.value = { visible: true, message, resolve }
    })
  }

  function resolveConfirm(result: boolean) {
    confirmState.value.resolve?.(result)
    confirmState.value = { visible: false, message: '', resolve: null }
  }

  return { toasts, confirmState, show, dismiss, success, error, confirm, resolveConfirm }
}
