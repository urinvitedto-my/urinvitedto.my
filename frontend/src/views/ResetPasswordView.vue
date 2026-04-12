<script setup lang="ts">
import { ref, onBeforeUnmount } from "vue"
import { useRouter } from "vue-router"
import { useAuthStore } from "@/stores/auth"
import { updatePassword } from "@/services/supabase"

const router = useRouter()
const authStore = useAuthStore()
const password = ref("")
const confirmPassword = ref("")
const showPassword = ref(false)
const loading = ref(false)
const error = ref("")
const success = ref(false)
const passwordUpdated = ref(false)

/** If the user leaves this page without updating their password, log them out. */
onBeforeUnmount(() => {
  if (!passwordUpdated.value) {
    authStore.logout()
  }
})

async function handleReset() {
  if (!password.value || !confirmPassword.value) {
    error.value = "Please fill in both fields"
    return
  }

  if (password.value.length < 6) {
    error.value = "Password must be at least 6 characters"
    return
  }

  if (password.value !== confirmPassword.value) {
    error.value = "Passwords do not match"
    return
  }

  loading.value = true
  error.value = ""

  try {
    await updatePassword(password.value)
    passwordUpdated.value = true
    success.value = true
    setTimeout(() => router.push("/host/dashboard"), 2000)
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : "Failed to update password"
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div
    class="min-h-screen flex items-center justify-center pt-24 pb-12 px-4"
  >
    <div class="max-w-md w-full">
      <div class="text-center mb-8">
        <h1 class="text-3xl font-bold text-primary">RESET PASSWORD</h1>
        <p class="text-gray-600 mt-2">Enter your new password</p>
      </div>

      <div
        v-if="success"
        class="bg-white rounded-lg shadow-sm p-8 text-center space-y-4"
      >
        <p class="text-green-600 font-medium">Password updated successfully!</p>
        <p class="text-gray-500 text-sm">Redirecting to dashboard...</p>
      </div>

      <form
        v-else
        @submit.prevent="handleReset"
        class="bg-white rounded-lg shadow-sm p-8 space-y-6"
      >
        <div>
          <label for="password" class="block text-sm font-medium text-gray-700 mb-2">
            New Password
          </label>
          <div class="relative">
            <input
              id="password"
              v-model="password"
              :type="showPassword ? 'text' : 'password'"
              required
              class="w-full px-4 py-3 pr-12 border border-muted rounded-lg focus:outline-none focus:ring-2 focus:ring-accent"
              :disabled="loading"
            />
            <button
              type="button"
              @click="showPassword = !showPassword"
              class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 transition-colors"
              tabindex="-1"
            >
              <svg
                v-if="!showPassword"
                class="w-5 h-5"
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z"
                />
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                />
              </svg>
              <svg
                v-else
                class="w-5 h-5"
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M3.98 8.223A10.477 10.477 0 001.934 12c1.292 4.338 5.31 7.5 10.066 7.5.993 0 1.953-.138 2.863-.395M6.228 6.228A10.45 10.45 0 0112 4.5c4.756 0 8.773 3.162 10.065 7.498a10.523 10.523 0 01-4.293 5.774M6.228 6.228L3 3m3.228 3.228l3.65 3.65m7.894 7.894L21 21m-3.228-3.228l-3.65-3.65m0 0a3 3 0 10-4.243-4.243m4.242 4.242L9.88 9.88"
                />
              </svg>
            </button>
          </div>
        </div>

        <div>
          <label for="confirm-password" class="block text-sm font-medium text-gray-700 mb-2">
            Confirm Password
          </label>
          <input
            id="confirm-password"
            v-model="confirmPassword"
            :type="showPassword ? 'text' : 'password'"
            required
            class="w-full px-4 py-3 border border-muted rounded-lg focus:outline-none focus:ring-2 focus:ring-accent"
            :disabled="loading"
          />
        </div>

        <p v-if="error" class="text-red-600 text-sm">{{ error }}</p>

        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-accent text-black font-semibold py-3 rounded-lg hover:bg-accent-dark transition-colors disabled:opacity-50"
        >
          {{ loading ? "Updating..." : "Update Password" }}
        </button>
      </form>
    </div>
  </div>
</template>
