<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Navbar from '@/components/Navbar.vue'
import Footer from '@/components/Footer.vue'

const route = useRoute()
const authStore = useAuthStore()
const showNavbar = computed(() => !route.meta.hideNavbar)
const showFooter = computed(() => !route.meta.hideFooter)

onMounted(() => {
  authStore.init()
})
</script>

<template>
  <div class="min-h-screen flex flex-col bg-[#ececec]">
    <Navbar v-if="showNavbar" />
    <main class="flex-1">
      <RouterView />
    </main>
    <Footer v-if="showFooter" />
  </div>
</template>

<style scoped></style>
