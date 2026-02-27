<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import Navbar from '@/components/Navbar.vue'
import Footer from '@/components/Footer.vue'

const route = useRoute()
const showFooter = computed(() => !route.meta.hideFooter)
</script>

<template>
  <div class="min-h-screen flex flex-col bg-[#ececec]">
    <Navbar />
    <main class="flex-1">
      <RouterView />
    </main>
    <Footer v-if="showFooter" />
  </div>
</template>

<style scoped></style>
