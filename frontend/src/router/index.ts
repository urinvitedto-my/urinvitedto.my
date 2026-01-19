import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
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
    },
    {
      path: '/:type(wedding|birthday|party)/:slug/guest',
      name: 'guest',
      component: () => import('@/views/GuestView.vue'),
      props: true,
    },
    {
      path: '/host/login',
      name: 'host-login',
      component: () => import('@/views/HostLoginView.vue'),
    },
    {
      path: '/host/dashboard',
      name: 'host-dashboard',
      component: () => import('@/views/HostDashboardView.vue'),
    },
    {
      path: '/admin',
      name: 'admin',
      component: () => import('@/views/AdminView.vue'),
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('@/views/NotFoundView.vue'),
    },
  ],
})

export default router
