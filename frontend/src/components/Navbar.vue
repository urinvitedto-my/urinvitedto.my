<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { onAuthStateChange, signOut } from '@/services/supabase'

const router = useRouter()
const route = useRoute()
const menuOpen = ref(false)
const isLoggedIn = ref(false)
const navbarVisible = ref(true)
let lastScrollY = 0

/**
 * Check if current page is the home page (needs light text on dark hero).
 */
const isHomePage = computed(() => route.path === '/')

let authSubscription: { unsubscribe: () => void } | null = null

/**
 * Handles scroll to hide/show navbar.
 */
function handleScroll() {
  const currentScrollY = window.scrollY

  if (currentScrollY > lastScrollY && currentScrollY > 80) {
    // scrolling down & past threshold - hide
    navbarVisible.value = false
    menuOpen.value = false
  } else {
    // scrolling up - show
    navbarVisible.value = true
  }

  lastScrollY = currentScrollY
}

onMounted(() => {
  // listen for auth state changes
  const { data } = onAuthStateChange((_event, session) => {
    isLoggedIn.value = !!session
  })
  authSubscription = data.subscription

  // listen for scroll
  window.addEventListener('scroll', handleScroll, { passive: true })
})

onUnmounted(() => {
  authSubscription?.unsubscribe()
  window.removeEventListener('scroll', handleScroll)
})

/**
 * Handles user logout and redirects to home.
 */
async function handleLogout() {
  await signOut()
  menuOpen.value = false
  router.push('/')
}
</script>

<template>
  <nav
    :class="[
      'fixed top-0 left-0 right-0 z-40 transition-transform duration-300',
      navbarVisible ? 'translate-y-0' : '-translate-y-full'
    ]"
  >
    <div class="max-w-6xl mx-auto px-4">
      <div class="flex items-center justify-between h-16">
        <!-- Logo -->
        <RouterLink
          to="/"
          :class="[
            'text-xl font-bold transition-colors',
            isHomePage ? 'text-white' : 'text-[#14213d]'
          ]"
        >
          urinvitedto.my
        </RouterLink>

        <!-- Desktop nav -->
        <div class="hidden md:flex items-center gap-6">
          <template v-if="isLoggedIn">
            <RouterLink
              to="/host/dashboard"
              :class="[
                'transition-colors font-medium',
                isHomePage ? 'text-white hover:text-white/80' : 'text-[#14213d] hover:text-[#14213d]/80'
              ]"
            >
              Dashboard
            </RouterLink>
            <span
              :class="[
                'h-4 w-px',
                isHomePage ? 'bg-white/30' : 'bg-gray-300'
              ]"
            ></span>
            <button
              @click="handleLogout"
              :class="[
                'transition-colors text-sm',
                isHomePage ? 'text-white/70 hover:text-white' : 'text-gray-500 hover:text-gray-700'
              ]"
            >
              Sign Out
            </button>
          </template>
          <RouterLink
            v-else
            to="/host/login"
            :class="[
              'transition-colors',
              isHomePage ? 'text-white/80 hover:text-white' : 'text-gray-600 hover:text-[#14213d]'
            ]"
          >
            Login
          </RouterLink>
        </div>

        <!-- Mobile menu button -->
        <button
          @click="menuOpen = !menuOpen"
          :class="['md:hidden', isHomePage ? 'text-white' : 'text-gray-600']"
        >
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
      <div v-if="menuOpen" class="md:hidden bg-white/95 backdrop-blur-sm rounded-b-lg py-4 px-4">
        <template v-if="isLoggedIn">
          <RouterLink
            to="/host/dashboard"
            class="block text-[#14213d] font-medium py-2 hover:text-[#14213d]/80 transition-colors"
            @click="menuOpen = false"
          >
            Dashboard
          </RouterLink>
          <hr class="my-3 border-gray-200" />
          <button
            @click="handleLogout"
            class="block text-gray-500 text-sm py-1 hover:text-gray-700 transition-colors"
          >
            Sign Out
          </button>
        </template>
        <RouterLink
          v-else
          to="/host/login"
          class="block text-gray-600 hover:text-[#14213d] transition-colors"
          @click="menuOpen = false"
        >
          Login
        </RouterLink>
      </div>
    </div>
  </nav>
</template>

<style scoped></style>
