import { ref } from "vue"

/**
 * Global page-loading flag shared between App.vue and route views.
 * App.vue shows the EnvelopeLoader overlay when this (or routeLoading) is true.
 * Views call start/stop to keep the overlay visible during data fetches.
 */
const loading = ref(false)

export function usePageLoading() {
  return {
    pageLoading: loading,
    startPageLoading() {
      loading.value = true
    },
    stopPageLoading() {
      loading.value = false
    },
  }
}
