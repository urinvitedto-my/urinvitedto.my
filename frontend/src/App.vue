<script setup lang="ts">
import { ref, computed, nextTick, onMounted, onUnmounted } from "vue"
import { useRoute, useRouter } from "vue-router"
import { useAuthStore } from "@/stores/auth"
import { usePageLoading } from "@/composables/usePageLoading"
import Navbar from "@/components/Navbar.vue"
import Footer from "@/components/Footer.vue"
import ToastContainer from "@/components/ToastContainer.vue"
import EnvelopeLoader from "@/components/EnvelopeLoader.vue"

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const { pageLoading } = usePageLoading()
const showNavbar = computed(() => !route.meta.hideNavbar)
const showFooter = computed(() => !route.meta.hideFooter)

/** True while a lazy-loaded route chunk is being downloaded. */
const routeLoading = ref(true)

router.beforeEach(() => {
  routeLoading.value = true
})

router.afterEach(() => {
  // nextTick gives the view's setup() a chance to call startPageLoading()
  // before routeLoading flips to false, preventing a one-frame gap.
  nextTick(() => {
    routeLoading.value = false
  })
})

router.onError(() => {
  routeLoading.value = false
})

onMounted(() => {
  authStore.init()
})

onUnmounted(() => {
  authStore.cleanup()
})
</script>

<template>
  <div class="min-h-screen flex flex-col bg-surface">
    <Transition name="loader-fade">
      <EnvelopeLoader v-if="routeLoading || pageLoading" />
    </Transition>
    <Navbar v-if="showNavbar" />
    <main class="flex-1">
      <RouterView v-slot="{ Component }">
        <component :is="Component" />
      </RouterView>
    </main>
    <Footer v-if="showFooter" />
    <ToastContainer />
  </div>
</template>

