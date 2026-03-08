<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const { isLoggedIn, userEmail, isAdmin } = storeToRefs(authStore)
const menuOpen = ref(false)
const navbarVisible = ref(true)
let lastScrollY = 0

/**
 * Pages that need light (white) nav text over dark backgrounds.
 */
const useLightNav = computed(() => {
  const name = route.name
  return typeof name === 'string' && ['event-landing', 'guest'].includes(name)
})

/**
 * Handles scroll to hide/show navbar. Only reappears near the top
 * to avoid floating over content on transparent-nav pages.
 */
function handleScroll() {
  const currentScrollY = window.scrollY

  if (currentScrollY > lastScrollY && currentScrollY > 15) {
    navbarVisible.value = false
    menuOpen.value = false
  } else if (currentScrollY <= 15) {
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
    const name = route.name
    const isPublicRoute =
      typeof name === 'string' && ['event-landing', 'guest'].includes(name)
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
    <div class="max-w-full mx-6 md:mx-16">
      <div class="flex items-center justify-between h-16">
        <!-- Logo -->
        <RouterLink to="/" class="block h-10 w-40 overflow-hidden rounded">
          <img
            src="/nav_logo.png"
            alt="urinvitedto.my"
            class="h-full w-full object-cover object-center"
          />
        </RouterLink>

        <!-- Desktop: LOGIN link when not logged in -->
        <RouterLink
          v-if="authStore.initialized && !isLoggedIn"
          to="/host/login"
          :class="[
            'hidden md:block font-bold uppercase tracking-wide transition-colors',
            useLightNav
              ? 'text-white hover:text-white/80'
              : 'text-primary hover:text-primary/80',
          ]"
        >
          LOGIN
        </RouterLink>

        <!-- Burger button + desktop dropdown -->
        <div
          data-navbar-menu
          :class="['relative', authStore.initialized && isLoggedIn ? '' : 'md:hidden']"
        >
          <button
            @click.stop="menuOpen = !menuOpen"
            :aria-label="menuOpen ? 'Close menu' : 'Open menu'"
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
            <template v-if="authStore.initialized && isLoggedIn">
              <div class="px-4 py-2.5 border-b border-gray-300">
                <div class="flex items-center gap-2">
                  <span
                    class="w-5 h-5 shrink-0 rounded-full border border-primary bg-white inline-flex items-center justify-center"
                  >
                    <svg
                      class="w-3 h-3 text-primary"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M15.75 6a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0ZM4.501 20.118a7.5 7.5 0 0 1 14.998 0A17.933 17.933 0 0 1 12 21.75c-2.676 0-5.216-.584-7.499-1.632Z"
                      />
                    </svg>
                  </span>
                  <p class="text-sm text-primary font-medium truncate mt-1">
                    {{ userEmail }}
                  </p>
                </div>
              </div>
              <div class="py-1">
                <RouterLink
                  to="/host/dashboard"
                  class="flex items-center gap-2 px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors"
                  @click="menuOpen = false"
                >
                  <svg
                    class="w-4 h-4 shrink-0"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M3.75 6A2.25 2.25 0 0 1 6 3.75h2.25A2.25 2.25 0 0 1 10.5 6v2.25a2.25 2.25 0 0 1-2.25 2.25H6a2.25 2.25 0 0 1-2.25-2.25V6ZM3.75 15.75A2.25 2.25 0 0 1 6 13.5h2.25a2.25 2.25 0 0 1 2.25 2.25V18a2.25 2.25 0 0 1-2.25 2.25H6A2.25 2.25 0 0 1 3.75 18v-2.25ZM13.5 6a2.25 2.25 0 0 1 2.25-2.25H18A2.25 2.25 0 0 1 20.25 6v2.25A2.25 2.25 0 0 1 18 10.5h-2.25a2.25 2.25 0 0 1-2.25-2.25V6ZM13.5 15.75a2.25 2.25 0 0 1 2.25-2.25H18a2.25 2.25 0 0 1 2.25 2.25V18A2.25 2.25 0 0 1 18 20.25h-2.25A2.25 2.25 0 0 1 13.5 18v-2.25Z"
                    />
                  </svg>
                  Dashboard
                </RouterLink>
                <RouterLink
                  v-if="isAdmin"
                  to="/admin"
                  class="flex items-center gap-2 px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors"
                  @click="menuOpen = false"
                >
                  <svg
                    class="w-4 h-4 shrink-0"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M9 12.75 11.25 15 15 9.75m-3-7.036A11.959 11.959 0 0 1 3.598 6 11.99 11.99 0 0 0 3 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285Z"
                    />
                  </svg>
                  Admin Dashboard
                </RouterLink>
              </div>
              <div class="border-t border-gray-100 py-1">
                <button
                  @click.stop="handleLogout"
                  class="w-full flex items-center gap-2 px-4 py-2 text-sm text-red-600 hover:bg-red-50 cursor-pointer transition-colors"
                >
                  <svg
                    class="w-4 h-4 shrink-0"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M8.25 9V5.25A2.25 2.25 0 0 1 10.5 3h6a2.25 2.25 0 0 1 2.25 2.25v13.5A2.25 2.25 0 0 1 16.5 21h-6a2.25 2.25 0 0 1-2.25-2.25V15m-3 0-3-3m0 0 3-3m-3 3H15"
                    />
                  </svg>
                  Logout
                </button>
              </div>
            </template>
            <RouterLink
              v-else-if="authStore.initialized"
              to="/host/login"
              class="block px-4 py-2.5 text-primary font-bold uppercase transition-colors hover:bg-gray-50"
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
        <template v-if="authStore.initialized && isLoggedIn">
          <div class="px-5 py-3 border-b border-gray-300 text-center">
            <div class="flex items-center justify-center gap-2 mt-1">
              <span
                class="w-5 h-5 shrink-0 rounded-full border border-primary bg-white inline-flex items-center justify-center"
              >
                <svg
                  class="w-3 h-3 text-primary"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M15.75 6a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0ZM4.501 20.118a7.5 7.5 0 0 1 14.998 0A17.933 17.933 0 0 1 12 21.75c-2.676 0-5.216-.584-7.499-1.632Z"
                  />
                </svg>
              </span>
              <p class="text-base text-primary font-medium truncate">
                {{ userEmail }}
              </p>
            </div>
          </div>
          <div class="py-1">
            <RouterLink
              to="/host/dashboard"
              class="flex items-center justify-center gap-2 px-5 py-3 text-base text-gray-700 hover:bg-gray-50 transition-colors"
              @click="menuOpen = false"
            >
              <svg
                class="w-4 h-4 shrink-0"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M3.75 6A2.25 2.25 0 0 1 6 3.75h2.25A2.25 2.25 0 0 1 10.5 6v2.25a2.25 2.25 0 0 1-2.25 2.25H6a2.25 2.25 0 0 1-2.25-2.25V6ZM3.75 15.75A2.25 2.25 0 0 1 6 13.5h2.25a2.25 2.25 0 0 1 2.25 2.25V18a2.25 2.25 0 0 1-2.25 2.25H6A2.25 2.25 0 0 1 3.75 18v-2.25ZM13.5 6a2.25 2.25 0 0 1 2.25-2.25H18A2.25 2.25 0 0 1 20.25 6v2.25A2.25 2.25 0 0 1 18 10.5h-2.25a2.25 2.25 0 0 1-2.25-2.25V6ZM13.5 15.75a2.25 2.25 0 0 1 2.25-2.25H18a2.25 2.25 0 0 1 2.25 2.25V18A2.25 2.25 0 0 1 18 20.25h-2.25A2.25 2.25 0 0 1 13.5 18v-2.25Z"
                />
              </svg>
              Dashboard
            </RouterLink>
            <RouterLink
              v-if="isAdmin"
              to="/admin"
              class="flex items-center justify-center gap-2 px-5 py-3 text-base text-gray-700 hover:bg-gray-50 transition-colors"
              @click="menuOpen = false"
            >
              <svg
                class="w-4 h-4 shrink-0"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M9 12.75 11.25 15 15 9.75m-3-7.036A11.959 11.959 0 0 1 3.598 6 11.99 11.99 0 0 0 3 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285Z"
                />
              </svg>
              Admin Dashboard
            </RouterLink>
          </div>
          <div class="border-t border-gray-100 py-1">
            <button
              @click.stop="handleLogout"
              class="w-full flex items-center justify-center gap-2 px-5 py-3 text-base text-red-600 hover:bg-red-50 cursor-pointer transition-colors"
            >
              <svg
                class="w-4 h-4 shrink-0"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M8.25 9V5.25A2.25 2.25 0 0 1 10.5 3h6a2.25 2.25 0 0 1 2.25 2.25v13.5A2.25 2.25 0 0 1 16.5 21h-6a2.25 2.25 0 0 1-2.25-2.25V15m-3 0-3-3m0 0 3-3m-3 3H15"
                />
              </svg>
              Logout
            </button>
          </div>
        </template>
        <RouterLink
          v-else-if="authStore.initialized"
          to="/host/login"
          class="block px-5 py-3 text-base text-center text-primary font-bold uppercase transition-colors hover:bg-gray-50"
          @click="menuOpen = false"
        >
          LOGIN
        </RouterLink>
      </div>
    </div>
  </nav>
</template>
