import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  supabase,
  signIn as authSignIn,
  signOut as authSignOut,
  getSession,
  onAuthStateChange,
} from '@/services/supabase'
import type { User, Session } from '@supabase/supabase-js'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const session = ref<Session | null>(null)
  const isAdmin = ref(false)
  const initialized = ref(false)

  const isLoggedIn = computed(() => !!session.value)
  const userEmail = computed(() => user.value?.email ?? '')

  /**
   * Queries the admins table to check if the email has admin access.
   */
  async function checkAdmin(email: string) {
    try {
      const { data, error } = await supabase
        .from('admins')
        .select('email')
        .eq('email', email)
        .maybeSingle()
      isAdmin.value = !error && !!data
    } catch {
      isAdmin.value = false
    }
  }

  let initPromise: Promise<void> | null = null

  /**
   * Subscribes to auth changes and sets initial state. Call once in App.vue.
   * Safe to call from multiple places -- subsequent calls await the same promise.
   */
  function init() {
    if (!initPromise) {
      initPromise = performInit()
    }
    return initPromise
  }

  async function performInit() {
    const currentSession = await getSession()
    session.value = currentSession
    user.value = currentSession?.user ?? null

    if (currentSession?.user?.email) {
      await checkAdmin(currentSession.user.email)
    }

    onAuthStateChange(async (_event, newSession) => {
      session.value = newSession
      user.value = newSession?.user ?? null

      if (newSession?.user?.email) {
        await checkAdmin(newSession.user.email)
      } else {
        isAdmin.value = false
      }
    })

    initialized.value = true
  }

  /**
   * Signs in with email/password via Supabase auth.
   */
  async function login(email: string, password: string) {
    const data = await authSignIn(email, password)
    session.value = data.session
    user.value = data.user

    if (data.user?.email) {
      await checkAdmin(data.user.email)
    }
  }

  /**
   * Signs out the current user and resets state.
   */
  async function logout() {
    await authSignOut()
    user.value = null
    session.value = null
    isAdmin.value = false
  }

  return {
    user,
    session,
    isAdmin,
    initialized,
    isLoggedIn,
    userEmail,
    init,
    login,
    logout,
  }
})
