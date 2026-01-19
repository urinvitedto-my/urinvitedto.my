<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { supabase, getUser, signOut } from '@/services/supabase'

const router = useRouter()
const loading = ref(true)
const isAdmin = ref(false)
const error = ref('')
const email = ref('')
const password = ref('')
const loginLoading = ref(false)

onMounted(async () => {
  await checkAdmin()
})

/**
 * Checks if user is logged in and is an admin.
 */
async function checkAdmin() {
  loading.value = true

  try {
    const user = await getUser()
    if (!user) {
      isAdmin.value = false
      loading.value = false
      return
    }

    // check if user email is in admins table
    const { data, error: fetchError } = await supabase
      .from('admins')
      .select('email')
      .eq('email', user.email)
      .single()

    if (fetchError || !data) {
      isAdmin.value = false
    } else {
      isAdmin.value = true
    }
  } catch (e: any) {
    // don't show auth errors on initial load - just show login form
    console.error('Admin check error:', e)
    isAdmin.value = false
  } finally {
    loading.value = false
  }
}

/**
 * Handles admin login.
 */
async function handleLogin() {
  if (!email.value || !password.value) {
    error.value = 'Please enter email and password'
    return
  }

  loginLoading.value = true
  error.value = ''

  try {
    await supabase.auth.signInWithPassword({
      email: email.value,
      password: password.value,
    })
    await checkAdmin()
  } catch (e: any) {
    error.value = e.message || 'Login failed'
  } finally {
    loginLoading.value = false
  }
}

/**
 * Handles logout.
 */
async function handleLogout() {
  await signOut()
  isAdmin.value = false
}
</script>

<template>
  <div class="admin-view min-h-screen py-8 px-4">
    <div class="max-w-6xl mx-auto">
      <!-- Loading -->
      <div v-if="loading" class="flex items-center justify-center py-20">
        <div class="animate-spin rounded-full h-12 w-12 border-4 border-[#fca311] border-t-transparent"></div>
      </div>

      <!-- Login Form -->
      <div v-else-if="!isAdmin" class="max-w-md mx-auto">
        <div class="text-center mb-8">
          <h1 class="text-3xl font-bold text-[#14213d]">Admin Login</h1>
          <p class="text-gray-600 mt-2">Sign in with your admin account</p>
        </div>

        <form @submit.prevent="handleLogin" class="bg-white rounded-lg shadow-sm p-8 space-y-6">
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
              Email
            </label>
            <input
              id="email"
              v-model="email"
              type="email"
              required
              class="w-full px-4 py-3 border border-[#e5e5e5] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#fca311]"
              :disabled="loginLoading"
            />
          </div>

          <div>
            <label for="password" class="block text-sm font-medium text-gray-700 mb-2">
              Password
            </label>
            <input
              id="password"
              v-model="password"
              type="password"
              required
              class="w-full px-4 py-3 border border-[#e5e5e5] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#fca311]"
              :disabled="loginLoading"
            />
          </div>

          <p v-if="error" class="text-red-600 text-sm">{{ error }}</p>

          <button
            type="submit"
            :disabled="loginLoading"
            class="w-full bg-[#14213d] text-white font-semibold py-3 rounded-lg hover:bg-[#1a2a4d] transition-colors disabled:opacity-50"
          >
            {{ loginLoading ? 'Signing in...' : 'Sign In' }}
          </button>
        </form>
      </div>

      <!-- Admin Dashboard -->
      <template v-else>
        <div class="flex items-center justify-between mb-8">
          <h1 class="text-2xl font-bold text-[#14213d]">Admin Dashboard</h1>
          <button
            @click="handleLogout"
            class="text-gray-600 hover:text-[#14213d] transition-colors"
          >
            Sign Out
          </button>
        </div>

        <div class="bg-white rounded-lg shadow-sm p-8">
          <p class="text-gray-600 text-center py-12">
            Admin CRUD interface coming soon. Use Supabase dashboard for direct data management.
          </p>
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped></style>
