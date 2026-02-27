<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const menuOpen = ref(false)
const navbarVisible = ref(true)
let lastScrollY = 0

const isLoggedIn = computed(() => authStore.isLoggedIn)
const userEmail = computed(() => authStore.userEmail)
const isAdmin = computed(() => authStore.isAdmin)

/**
 * Pages that need light (white) nav text over dark backgrounds.
 */
const useLightNav = computed(() =>
  ['home', 'event-landing', 'guest'].includes(route.name as string),
)

/**
 * Handles scroll to hide/show navbar. Only reappears near the top
 * to avoid floating over content on transparent-nav pages.
 */
function handleScroll() {
  const currentScrollY = window.scrollY

  if (currentScrollY > lastScrollY && currentScrollY > 50) {
    navbarVisible.value = false
    menuOpen.value = false
  } else if (currentScrollY <= 50) {
    navbarVisible.value = true
  }

  lastScrollY = currentScrollY
}

/**
 * Closes the menu when clicking outside of it.
 */
function handleClickOutside(e: MouseEvent) {
  const target = e.target as HTMLElement
  if (!target.closest('[data-navbar-menu]')) {
    menuOpen.value = false
  }
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll, { passive: true })
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
  document.removeEventListener('click', handleClickOutside)
})

/**
 * Handles user logout. Redirects to home unless on a public route
 * (event landing or guest page), where we stay put.
 */
async function handleLogout() {
  try {
    await authStore.logout()
  } finally {
    menuOpen.value = false
    const isPublicRoute = ['event-landing', 'guest'].includes(route.name as string)
    if (!isPublicRoute) {
      router.push('/')
    }
  }
}
</script>

<template>
  <nav
    :class="[
      'fixed top-0 left-0 right-0 z-40 transition-transform duration-300',
      navbarVisible ? 'translate-y-0' : '-translate-y-full',
    ]"
  >
    <div class="max-w-6xl mx-auto px-4">
      <div class="flex items-center justify-between h-16">
        <!-- Logo -->
        <RouterLink
          to="/"
          :class="[
            'text-2xl font-bold transition-colors',
            useLightNav ? 'text-white' : 'text-[#14213d]',
          ]"
        >
          urinvitedto.my
        </RouterLink>

        <!-- Desktop: LOGIN link when not logged in -->
        <RouterLink
          v-if="!isLoggedIn"
          to="/host/login"
          :class="[
            'hidden md:block font-bold uppercase tracking-wide transition-colors',
            useLightNav
              ? 'text-white hover:text-white/80'
              : 'text-[#14213d] hover:text-[#14213d]/80',
          ]"
        >
          LOGIN
        </RouterLink>

        <!-- Burger button + desktop dropdown -->
        <div data-navbar-menu :class="['relative', isLoggedIn ? '' : 'md:hidden']">
          <button
            @click.stop="menuOpen = !menuOpen"
            :class="[
              'cursor-pointer p-1.5 rounded-md transition-colors',
              useLightNav
                ? 'text-white hover:bg-white/10'
                : 'text-gray-600 hover:bg-gray-100',
            ]"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M4 6h16M4 12h16M4 18h16"
              />
            </svg>
          </button>

          <!-- Desktop dropdown -->
          <div
            v-if="menuOpen"
            class="hidden md:block absolute right-0 top-full mt-2 min-w-[240px] bg-white rounded-xl py-2 shadow-xl ring-1 ring-black/5"
          >
            <template v-if="isLoggedIn">
              <div class="px-4 py-2.5 border-b border-gray-100">
                <p class="text-xs text-gray-400 leading-none">Signed in as</p>
                <p class="text-sm text-[#14213d] font-medium truncate mt-1">
                  {{ userEmail }}
                </p>
              </div>
              <div class="py-1">
                <RouterLink
                  to="/host/dashboard"
                  class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors"
                  @click="menuOpen = false"
                >
                  Dashboard
                </RouterLink>
                <RouterLink
                  v-if="isAdmin"
                  to="/admin"
                  class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors"
                  @click="menuOpen = false"
                >
                  Admin Dashboard
                </RouterLink>
              </div>
              <div class="border-t border-gray-100 py-1">
                <button
                  @click.stop="handleLogout"
                  class="w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50 cursor-pointer transition-colors"
                >
                  Logout
                </button>
              </div>
            </template>
            <RouterLink
              v-else
              to="/host/login"
              class="block px-4 py-2.5 text-[#14213d] font-bold uppercase transition-colors hover:bg-gray-50"
              @click="menuOpen = false"
            >
              LOGIN
            </RouterLink>
          </div>
        </div>
      </div>

      <!-- Mobile menu -->
      <div
        v-if="menuOpen"
        data-navbar-menu
        class="md:hidden bg-white rounded-xl mt-2 py-3 shadow-xl ring-1 ring-black/5"
      >
        <template v-if="isLoggedIn">
          <div class="px-5 py-3 border-b border-gray-100 text-center">
            <p class="text-xs text-gray-400 leading-none">Signed in as</p>
            <p class="text-base text-[#14213d] font-medium truncate mt-1">
              {{ userEmail }}
            </p>
          </div>
          <div class="py-1">
            <RouterLink
              to="/host/dashboard"
              class="block px-5 py-3 text-base text-center text-gray-700 hover:bg-gray-50 transition-colors"
              @click="menuOpen = false"
            >
              Dashboard
            </RouterLink>
            <RouterLink
              v-if="isAdmin"
              to="/admin"
              class="block px-5 py-3 text-base text-center text-gray-700 hover:bg-gray-50 transition-colors"
              @click="menuOpen = false"
            >
              Admin Dashboard
            </RouterLink>
          </div>
          <div class="border-t border-gray-100 py-1">
            <button
              @click.stop="handleLogout"
              class="w-full px-5 py-3 text-base text-center text-red-600 hover:bg-red-50 cursor-pointer transition-colors"
            >
              Logout
            </button>
          </div>
        </template>
        <RouterLink
          v-else
          to="/host/login"
          class="block px-5 py-3 text-base text-center text-[#14213d] font-bold uppercase transition-colors hover:bg-gray-50"
          @click="menuOpen = false"
        >
          LOGIN
        </RouterLink>
      </div>
    </div>
  </nav>
</template>

<style scoped></style>
