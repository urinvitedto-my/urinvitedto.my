<script setup lang="ts">
import { ref } from 'vue'

const menuOpen = ref(false)
</script>

<template>
  <nav class="bg-white shadow-sm sticky top-0 z-40">
    <div class="max-w-6xl mx-auto px-4">
      <div class="flex items-center justify-between h-16">
        <!-- Logo -->
        <RouterLink to="/" class="text-xl font-bold text-[#14213d]">
          urinvitedto.my
        </RouterLink>

        <!-- Desktop nav -->
        <div class="hidden md:flex items-center gap-6">
          <RouterLink
            to="/host/login"
            class="text-gray-600 hover:text-[#14213d] transition-colors"
          >
            Login
          </RouterLink>
        </div>

        <!-- Mobile menu button -->
        <button @click="menuOpen = !menuOpen" class="md:hidden text-gray-600">
          <svg
            v-if="!menuOpen"
            class="w-6 h-6"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M4 6h16M4 12h16M4 18h16"
            />
          </svg>
          <svg
            v-else
            class="w-6 h-6"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </button>
      </div>

      <!-- Mobile menu -->
      <div v-if="menuOpen" class="md:hidden border-t border-[#e5e5e5] py-4 space-y-4">
        <RouterLink
          to="/host/login"
          class="block text-gray-600 hover:text-[#14213d] transition-colors"
          @click="menuOpen = false"
        >
          Host Login
        </RouterLink>
      </div>
    </div>
  </nav>
</template>

<style scoped></style>
