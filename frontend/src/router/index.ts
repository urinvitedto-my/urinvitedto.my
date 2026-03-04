import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  scrollBehavior(_to, _from, savedPosition) {
    return savedPosition || { top: 0 }
  },
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/LandingView.vue'),
    },
    {
      path: '/:type(wedding|birthday|party)/:slug',
      name: 'event-landing',
      component: () => import('@/views/EventLandingView.vue'),
      props: true,
      meta: { hideFooter: true, hideNavbar: true },
    },
    {
      path: '/:type(wedding|birthday|party)/:slug/guest',
      name: 'guest',
      component: () => import('@/views/GuestView.vue'),
      props: true,
      meta: { hideFooter: true, hideNavbar: true },
    },
    {
      path: '/host/login',
      name: 'host-login',
      component: () => import('@/views/HostLoginView.vue'),
      meta: { guestOnly: true },
    },
    {
      path: '/host/dashboard',
      name: 'host-dashboard',
      component: () => import('@/views/HostDashboardView.vue'),
      meta: { hideFooter: true, requiresAuth: true },
    },
    {
      path: '/admin',
      name: 'admin',
      component: () => import('@/views/AdminView.vue'),
      meta: { hideFooter: true, requiresAuth: true, requiresAdmin: true },
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('@/views/NotFoundView.vue'),
    },
  ],
})

router.beforeEach(async (to) => {
  // Lazy-import to avoid circular dep at module level
  const { useAuthStore } = await import('@/stores/auth')
  const authStore = useAuthStore()

  if (!authStore.initialized) {
    try {
      await authStore.init()
    } catch {
      return { name: 'host-login' }
    }
  }

  // Redirect logged-in users away from login page
  if (to.meta.guestOnly && authStore.isLoggedIn) {
    return { name: 'host-dashboard' }
  }

  // Protected routes: must be logged in
  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    return { name: 'host-login' }
  }

  // Admin routes: must be admin
  if (to.meta.requiresAdmin && !authStore.isAdmin) {
    return { name: 'host-dashboard' }
  }
})

export default router
