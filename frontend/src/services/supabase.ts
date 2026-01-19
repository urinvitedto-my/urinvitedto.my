// Supabase client for auth and direct DB access (admin/host dashboards)
import { createClient } from '@supabase/supabase-js'

const supabaseUrl = import.meta.env.VITE_SUPABASE_URL
const supabaseAnonKey = import.meta.env.VITE_SUPABASE_ANON_KEY

if (!supabaseUrl || !supabaseAnonKey) {
  console.warn('Supabase env vars not set. Auth features will not work.')
}

export const supabase = createClient(supabaseUrl || '', supabaseAnonKey || '')

/**
 * Signs in with email and password.
 */
export async function signIn(email: string, password: string) {
  const { data, error } = await supabase.auth.signInWithPassword({ email, password })
  if (error) throw error
  return data
}

/**
 * Signs out the current user.
 */
export async function signOut() {
  const { error } = await supabase.auth.signOut()
  if (error) throw error
}

/**
 * Gets the current session.
 */
export async function getSession() {
  const { data, error } = await supabase.auth.getSession()
  if (error) throw error
  return data.session
}

/**
 * Gets the current user. Returns null if not logged in.
 */
export async function getUser() {
  const { data, error } = await supabase.auth.getUser()
  // don't throw on "no session" - just return null
  if (error) {
    if (error.message?.includes('session') || error.name === 'AuthSessionMissingError') {
      return null
    }
    throw error
  }
  return data.user
}

/**
 * Listens for auth state changes.
 */
export function onAuthStateChange(callback: (event: string, session: any) => void) {
  return supabase.auth.onAuthStateChange(callback)
}
