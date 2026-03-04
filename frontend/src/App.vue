<script setup lang="ts">
import { computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Navbar from '@/components/Navbar.vue'
import Footer from '@/components/Footer.vue'
import ToastContainer from '@/components/ToastContainer.vue'

const route = useRoute()
const authStore = useAuthStore()
const showNavbar = computed(() => !route.meta.hideNavbar)
const showFooter = computed(() => !route.meta.hideFooter)

onMounted(() => {
  authStore.init()
})

onUnmounted(() => {
  authStore.cleanup()
})
</script>

<template>
  <div class="min-h-screen flex flex-col bg-surface">
    <Navbar v-if="showNavbar" />
    <main class="flex-1">
      <RouterView v-slot="{ Component }">
        <Transition name="fade" mode="out-in">
          <component :is="Component" />
        </Transition>
      </RouterView>
    </main>
    <Footer v-if="showFooter" />
    <ToastContainer />
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
