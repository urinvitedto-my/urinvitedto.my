<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { signIn } from '@/services/supabase'

const router = useRouter()
const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

/**
 * Handles login form submission.
 */
async function handleLogin() {
  if (!email.value || !password.value) {
    error.value = 'Please enter email and password'
    return
  }

  loading.value = true
  error.value = ''

  try {
    await signIn(email.value, password.value)
    router.push('/host/dashboard')
  } catch (e: any) {
    error.value = e.message || 'Login failed'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="host-login-view min-h-screen flex items-center justify-center pt-24 pb-12 px-4">
    <div class="max-w-md w-full">
      <div class="text-center mb-8">
        <h1 class="text-3xl font-bold text-[#14213d]">Host Login</h1>
        <p class="text-gray-600 mt-2">Sign in to view your event dashboard</p>
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
            :disabled="loading"
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
            :disabled="loading"
          />
        </div>

        <p v-if="error" class="text-red-600 text-sm">{{ error }}</p>

        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-[#fca311] text-black font-semibold py-3 rounded-lg hover:bg-[#e5930f] transition-colors disabled:opacity-50"
        >
          {{ loading ? 'Signing in...' : 'Sign In' }}
        </button>
      </form>
    </div>
  </div>
</template>

<style scoped></style>
